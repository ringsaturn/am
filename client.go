package am

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	V1_TOKEN               = "https://maps-api.apple.com/v1/token"              // https://developer.apple.com/documentation/applemapsserverapi/generate_a_maps_access_token
	V1_GEO_CODE            = "https://maps-api.apple.com/v1/geocode"            // https://developer.apple.com/documentation/applemapsserverapi/geocode_an_address
	V1_REVERSE_GEO_CODE    = "https://maps-api.apple.com/v1/reverseGeocode"     // https://developer.apple.com/documentation/applemapsserverapi/reverse_geocode_a_location
	V1_SEARCH              = "https://maps-api.apple.com/v1/search"             // https://developer.apple.com/documentation/applemapsserverapi/search_for_places_that_match_specific_criteria
	V1_SEARCH_AUTOCOMPLETE = "https://maps-api.apple.com/v1/searchAutocomplete" // https://developer.apple.com/documentation/applemapsserverapi/search_for_places_that_meet_specific_criteria_to_autocomplete_a_place_search
	V1_DIRECTIONS          = "https://maps-api.apple.com/v1/directions"         // https://developer.apple.com/documentation/applemapsserverapi/search_for_directions_and_estimated_travel_time_between_locations
	V1_ETAS                = "https://maps-api.apple.com/v1/etas"               // https://developer.apple.com/documentation/applemapsserverapi/determine_estimated_arrival_times_and_distances_to_one_or_more_destinations
)

// AccessTokenSaver is an interface to save and get access token.
//
// Please implement this interface if you want to save access token in Redis or other places.
type AccessTokenSaver interface {
	GetAccessToken(context.Context) (string, int64, error)
	SetAccessToken(context.Context, string, int64) error
}

type memorySaver struct {
	mapAuthToken       string
	mapAccessToken     string
	mapAccessTokenExp  int64
	mapAccessTokenLock sync.RWMutex
}

func (s *memorySaver) GetAccessToken(ctx context.Context) (string, int64, error) {
	s.mapAccessTokenLock.Lock()
	defer s.mapAccessTokenLock.Unlock()
	return s.mapAccessToken, s.mapAccessTokenExp, nil
}

func (s *memorySaver) SetAccessToken(ctx context.Context, accessToken string, exp int64) error {
	s.mapAccessTokenLock.RLock()
	s.mapAccessToken = accessToken
	s.mapAccessTokenExp = exp
	s.mapAccessTokenLock.RUnlock()
	return nil
}

// AutoRefresh is a function to refresh access token based on token expire time.
type AutoRefresh func(ctx context.Context, client Client) (string, int64, error)

func newAutoRefresh() AutoRefresh {
	mutex := &sync.Mutex{}
	return func(ctx context.Context, client Client) (string, int64, error) {
		token, exp, err := client.GetAccessToken(ctx)
		if err != nil {
			return "", 0, err
		}
		// If lock failed, it means another goroutine is refreshing token.
		// So we just return current token.
		locked := mutex.TryLock()
		if !locked {
			return token, exp, nil
		}
		defer mutex.Unlock()
		now := time.Now().Unix()
		if exp-now > 60 {
			return "", 0, nil
		}
		resp, err := client.GetNewAccessToken(ctx)
		if err != nil {
			return "", 0, err
		}
		err = client.SetAccessToken(ctx, resp.AccessToken, resp.ExpiresInSeconds)
		if err != nil {
			return "", 0, err
		}
		return resp.AccessToken, resp.ExpiresInSeconds, nil
	}
}

type Client interface {
	AccessTokenSaver

	GetNewAccessToken(context.Context) (*AccessTokenResponse, error)
	Geocode(context.Context, *GeocodeRequest) (*PlaceResults, error)
	ReverseGeocode(context.Context, *ReverseRequest) (*PlaceResults, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	SearchAutoComplete(context.Context, *SearchAutoCompleteRequest) (*SearchAutocompleteResponse, error)
	Directions(context.Context, *DirectionsRequest) (*DirectionsResponse, error)
	Eta(context.Context, *EtaRequest) (*EtaResponse, error)
}

type baseClient struct {
	authToken     string
	tokenSaver    AccessTokenSaver
	client        *http.Client
	autoRefreshFn AutoRefresh
}

type Option func(*baseClient)

// Will use `memorySaver`(an internal function) default.
//
// Please implement your own `AccessTokenSaver` if you want to save access token
// in Redis or other places.
func WithTokenSaver(saver AccessTokenSaver) Option {
	return func(c *baseClient) {
		c.tokenSaver = saver
	}
}

// Will use `defaultAutoRefresh`(an internal function) by default.
// If you want to disable auto refresh, please set this option to nil.
//
// If you want to implement your own auto refresh function, please make sure the
// function is thread safe. Because the function could be called by multiple
// goroutines.
func WithAutoTokenRefresh(fn AutoRefresh) Option {
	return func(c *baseClient) {
		c.autoRefreshFn = fn
	}
}

// Will use [http.DefaultClient] by default.
func WithHTTPClient(client *http.Client) Option {
	return func(c *baseClient) {
		c.client = client
	}
}

func NewClient(authToken string, opts ...Option) Client {
	c := &baseClient{
		tokenSaver: &memorySaver{
			mapAuthToken: authToken,
		},
		client:        http.DefaultClient,
		autoRefreshFn: newAutoRefresh(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func handleErr(httpStatusCode int, header http.Header, bodyBytes []byte) error {
	err := &ErrorFromAPI{
		StatusCode: httpStatusCode,
		RawBody:    bodyBytes,
		Header:     header,
	}
	resp := &ErrorResponse{}
	if _err := json.Unmarshal(bodyBytes, resp); _err != nil {
		return err
	}
	err.Response = resp
	return err
}

type query interface {
	URLValues() (url.Values, error)
}

func do[expect any](
	ctx context.Context,
	httpClient *http.Client,
	api string,
	token string,
	req query,
) (*expect, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	if req != nil {
		q, err := req.URLValues()
		if err != nil {
			return nil, err
		}
		request.URL.RawQuery = q.Encode()
	}
	httpResponse, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	bodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	header := httpResponse.Header
	if httpResponse.StatusCode != http.StatusOK {
		return nil, handleErr(httpResponse.StatusCode, header, bodyBytes)
	}
	resp := new(expect)
	err = json.Unmarshal(bodyBytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *baseClient) GetNewAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	return do[AccessTokenResponse](ctx, http.DefaultClient, V1_TOKEN, c.authToken, nil)
}

func (c *baseClient) GetAccessToken(ctx context.Context) (string, int64, error) {
	return c.tokenSaver.GetAccessToken(ctx)
}

func (c *baseClient) SetAccessToken(ctx context.Context, accessToken string, exp int64) error {
	return c.tokenSaver.SetAccessToken(ctx, accessToken, exp)
}

func doWithReadAccessToken[expect any](
	ctx context.Context,
	c *baseClient,
	autoFresh AutoRefresh,
	api string,
	req query,
) (*expect, error) {
	var (
		accessToken string
		err         error
	)
	if autoFresh != nil {
		accessToken, _, err = autoFresh(ctx, c)
	} else {
		accessToken, _, err = c.GetAccessToken(ctx)
	}
	if err != nil {
		return nil, err
	}
	return do[expect](ctx, c.client, api, accessToken, req)
}

func (c *baseClient) Geocode(ctx context.Context, req *GeocodeRequest) (*PlaceResults, error) {
	return doWithReadAccessToken[PlaceResults](ctx, c, c.autoRefreshFn, V1_GEO_CODE, req)
}

func (c *baseClient) ReverseGeocode(ctx context.Context, req *ReverseRequest) (*PlaceResults, error) {
	return doWithReadAccessToken[PlaceResults](ctx, c, c.autoRefreshFn, V1_REVERSE_GEO_CODE, req)
}

func (c *baseClient) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return doWithReadAccessToken[SearchResponse](ctx, c, c.autoRefreshFn, V1_SEARCH, req)
}

func (c *baseClient) SearchAutoComplete(ctx context.Context, req *SearchAutoCompleteRequest) (*SearchAutocompleteResponse, error) {
	return doWithReadAccessToken[SearchAutocompleteResponse](ctx, c, c.autoRefreshFn, V1_SEARCH_AUTOCOMPLETE, req)
}

func (c *baseClient) Directions(ctx context.Context, req *DirectionsRequest) (*DirectionsResponse, error) {
	return doWithReadAccessToken[DirectionsResponse](ctx, c, c.autoRefreshFn, V1_DIRECTIONS, req)
}

func (c *baseClient) Eta(ctx context.Context, req *EtaRequest) (*EtaResponse, error) {
	return doWithReadAccessToken[EtaResponse](ctx, c, c.autoRefreshFn, V1_ETAS, req)
}
