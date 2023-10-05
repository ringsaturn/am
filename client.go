package am

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	ENDPOINT               = "https://maps.apple.com"
	V1_TOKEN               = "/v1/token"              // https://developer.apple.com/documentation/applemapsserverapi/generate_a_maps_access_token
	V1_GEO_CODE            = "/v1/geocode"            // https://developer.apple.com/documentation/applemapsserverapi/geocode_an_address
	V1_REVERSE_GEO_CODE    = "/v1/reverseGeocode"     // https://developer.apple.com/documentation/applemapsserverapi/reverse_geocode_a_location
	V1_SEARCH              = "/v1/search"             // https://developer.apple.com/documentation/applemapsserverapi/search_for_places_that_match_specific_criteria
	V1_SEARCH_AUTOCOMPLETE = "/v1/searchAutocomplete" // https://developer.apple.com/documentation/applemapsserverapi/search_for_places_that_meet_specific_criteria_to_autocomplete_a_place_search
	V1_DIRECTIONS          = "/v1/directions"         // https://developer.apple.com/documentation/applemapsserverapi/search_for_directions_and_estimated_travel_time_between_locations
	V1_ETAS                = "/v1/etas"               // https://developer.apple.com/documentation/applemapsserverapi/determine_estimated_arrival_times_and_distances_to_one_or_more_destinations
)

type Client interface {
	Geocode(ctx context.Context, req *GeocodeRequest) (*PlaceResults, error)
	ReverseGeocode(ctx context.Context, req *ReverseRequest) (*PlaceResults, error)
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
	SearchAutoComplete(ctx context.Context, req *SearchAutoCompleteRequest) (*SearchAutocompleteResponse, error)
}

type BaseClient struct {
	mapAuthToken       string
	mapAccessToken     string
	mapAccessTokenExp  int64
	mapAccessTokenLock sync.RWMutex

	HTTPClient *http.Client

	Endpoint          string
	TokenAPI          string
	GeoCodeAPI        string
	ReverseGeocodeAPI string
	SearchAPI         string
	SearchAutoAPI     string
	DirectionsAPI     string
	EtasAPI           string
}

type Option func(*BaseClient)

// Use this func to set mapAuthToken.
func WithMapAuthToken(token string) Option {
	return func(c *BaseClient) {
		c.mapAuthToken = token
	}
}

func WithMapAccessToken(token string, expire int64) Option {
	return func(c *BaseClient) {
		c.mapAccessToken = token
		c.mapAccessTokenExp = expire
	}
}

func NewClient(opts ...Option) Client {
	c := &BaseClient{
		HTTPClient:        http.DefaultClient,
		Endpoint:          ENDPOINT,
		TokenAPI:          V1_TOKEN,
		GeoCodeAPI:        V1_GEO_CODE,
		ReverseGeocodeAPI: V1_REVERSE_GEO_CODE,
		SearchAPI:         V1_SEARCH,
		SearchAutoAPI:     V1_SEARCH_AUTOCOMPLETE,
		DirectionsAPI:     V1_DIRECTIONS,
		EtasAPI:           V1_ETAS,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func handleErr(httpStatusCode int, bodyBytes []byte) error {
	err := &ErrorFromAPI{
		StatusCode: httpStatusCode,
		RawBody:    bodyBytes,
	}
	resp := &ErrorResponse{}
	if _err := json.Unmarshal(bodyBytes, resp); _err != nil {
		// Server return invalid json.
		// Just return the raw body.
		return err
	}
	err.Response = resp
	return err
}

func (c *BaseClient) RefreshMapAccessToken(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Endpoint+c.TokenAPI, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+c.mapAuthToken)
	httpResponse, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	bodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return handleErr(httpResponse.StatusCode, bodyBytes)
	}
	resp := &AccessTokenResponse{}
	if err := json.Unmarshal(bodyBytes, resp); err != nil {
		return err
	}
	c.mapAccessTokenLock.Lock()
	c.mapAccessToken = resp.AccessToken
	c.mapAccessTokenExp = resp.ExpiresInSeconds
	c.mapAccessTokenLock.Unlock()
	return nil
}

type query interface {
	URLValues() (url.Values, error)
}

func do[expect any](ctx context.Context, c *BaseClient, req query) (*expect, error) {
	q, err := req.URLValues()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Endpoint+c.GeoCodeAPI, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.mapAccessToken)
	request.URL.RawQuery = q.Encode()
	httpResponse, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	bodyBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return nil, handleErr(httpResponse.StatusCode, bodyBytes)
	}
	resp := new(expect)
	if err := json.Unmarshal(bodyBytes, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *BaseClient) Geocode(ctx context.Context, req *GeocodeRequest) (*PlaceResults, error) {
	return do[PlaceResults](ctx, c, req)
}

func (c *BaseClient) ReverseGeocode(ctx context.Context, req *ReverseRequest) (*PlaceResults, error) {
	return do[PlaceResults](ctx, c, req)
}

func (c *BaseClient) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return do[SearchResponse](ctx, c, req)
}

func (c *BaseClient) SearchAutoComplete(ctx context.Context, req *SearchAutoCompleteRequest) (*SearchAutocompleteResponse, error) {
	return do[SearchAutocompleteResponse](ctx, c, req)
}

func (c *BaseClient) Directions(ctx context.Context, req *DirectionsRequest) (*DirectionsResponse, error) {
	return do[DirectionsResponse](ctx, c, req)
}

func (c *BaseClient) Eta(ctx context.Context, req *EtaRequest) (*EtaResponse, error) {
	return do[EtaResponse](ctx, c, req)
}
