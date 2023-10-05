package am

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/biter777/countries"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"golang.org/x/text/language"
)

// Ensure impl query
var (
	_ query = (*GeocodeRequest)(nil)
	_ query = (*ReverseRequest)(nil)
	_ query = (*SearchRequest)(nil)
	_ query = (*SearchAutoCompleteRequest)(nil)
	_ query = (*DirectionsRequest)(nil)
	_ query = (*EtaRequest)(nil)
)

func poiCategoriesToString(pois []PoiCategory) string {
	items := []string{}
	for _, poi := range pois {
		items = append(items, string(poi))
	}
	return strings.Join(items, ",")
}

func countriesToString(countries []countries.CountryCode) string {
	items := []string{}
	for _, c := range countries {
		items = append(items, c.Alpha2())
	}
	return strings.Join(items, ",")
}

func avoidToString(avoid []DirectionsAvoid) string {
	items := []string{}
	for _, a := range avoid {
		items = append(items, string(a))
	}
	return strings.Join(items, ",")
}

// https://developer.apple.com/documentation/applemapsserverapi/geocode_an_address
type GeocodeRequest struct {
	// (Required) The address to geocode. For example: q=1 Apple Park, Cupertino, CA
	Query string `query:"q" vd:"$!=''"`

	// A comma-separated list of two-letter ISO 3166-1 codes to limit the results to.
	// For example: limitToCountries=US,CA.
	//
	// If you specify two or more countries, the results reflect the best
	// available results for some or all of the countries rather than everything
	// related to the query for those countries.
	LimitToCountries []countries.CountryCode `query:"limitToCountries,omitempty"`

	// The language the server should use when returning the response, specified
	// using a BCP 47 language code. For example, for English use lang=en-US.
	//
	// Default: en-US
	Lang language.Tag `query:"lang"`

	// A location defined by the application as a hint. Specify the location as
	// a comma-separated string containing the latitude and longitude.
	// For example, searchLocation=37.78,-122.42.
	SearchLocation Location `query:"searchLocation,omitempty"`

	// A region the app defines as a hint. Specify the region specified as a
	// comma-separated string that describes the region in the form
	// north-latitude, east-longitude, south-latitude, west-longitude.
	// For example, searchRegion=38,-122.1,37.5,-122.5.
	SearchRegion Region `query:"searchRegion,omitempty"`

	// The location of the user, specified as a comma-separated string that
	// contains the latitude and longitude.
	// For example: userLocation=37.78,-122.42.
	//
	// Certain APIs, such as Searching, may opt to
	// use the userLocation, if specified, as a fallback for the searchLocation.
	UserLocation Location `query:"userLocation,omitempty"`
}

func (req *GeocodeRequest) Validate() error { return vd.Validate(req) }

func (req *GeocodeRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	q := make(url.Values)
	q.Set("q", req.Query)

	if len(req.LimitToCountries) > 0 {
		q.Set("limitToCountries", countriesToString(req.LimitToCountries))
	}

	if !req.Lang.IsRoot() {
		q.Set("lang", req.Lang.String())
	}
	if !req.SearchLocation.IsEmpty() {
		q.Set("searchLocation", req.SearchLocation.QueryString())
	}
	if !req.SearchRegion.IsEmpty() {
		q.Set("searchRegion", req.SearchRegion.QueryString())
	}
	if !req.UserLocation.IsEmpty() {
		q.Set("userLocation", req.UserLocation.QueryString())
	}
	return q, nil
}

type ReverseRequest struct {
	// (Required) The coordinate to reverse geocode as a comma-separated string
	// that contains the latitude and longitude. For example: loc=37.3316851,-122.0300674.
	Loc Location `query:"loc"`

	// The language the server uses when returning the response, specified using
	// a BCP 47 language code. For example, for English, use lang=en-US.
	// Default: en-US
	Lang language.Tag `query:"lang"`
}

func (req *ReverseRequest) Validate() error {
	if req.Loc.IsEmpty() {
		return errors.New("am: loc is required")
	}
	return vd.Validate(req)
}

func (req *ReverseRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Set("loc", req.Loc.QueryString())
	if !req.Lang.IsRoot() {
		q.Set("lang", req.Lang.String())
	}
	return q, nil
}

type SearchRequest struct {
	// (Required) The place to search for. For example, q=eiffel tower.
	Query string `query:"q" vd:"$!=''"`

	// A comma-separated list of strings that describes the points of interest
	// to exclude from the search results.
	// For example, excludePoiCategories=Restaurant,Cafe.
	//
	// See PoiCategory for a complete list of possible values.
	ExcludePoiCategories []PoiCategory `query:"excludePoiCategories"`

	// A comma-separated list of strings that describes the points of interest
	// to include in the search results.
	// For example, includePoiCategories=Restaurant,Cafe.
	//
	// See PoiCategory for a complete list of possible values.
	IncludePoiCategories []PoiCategory `query:"includePoiCategories"`

	// A comma-separated list of two-letter ISO 3166-1 codes of the countries to
	// limit the results to.
	// For example, limitToCountries=US,CA limits the search to the United
	// States and Canada.
	//
	// If you specify two or more countries, the results reflect the best
	// available results for some or all of the countries rather than everything
	// related to the query for those countries.
	LimitToCountries []countries.CountryCode `query:"limitToCountries,omitempty"`

	// A comma-separated list of strings that describes the kind of result types
	// to include in the response. For example, resultTypeFilter=Poi.
	//
	// Possible Values: Poi, Address
	ResultTypeFilter []string `query:"resultTypeFilter"`

	// The language the server should use when returning the response, specified
	// using a BCP 47 language code.
	// For example, for English use lang=en-US. Defaults to en-US.
	//
	// Default: en-US
	Lang language.Tag `query:"lang"`

	// A location defined by the application as a hint. Specify the location as
	// a comma-separated string containing the latitude and longitude.
	// For example, searchLocation=37.78,-122.42.
	SearchLocation Location `query:"searchLocation"`

	// A region the app defines as a hint. Specify the region specified as a
	// comma-separated string that describes the region in the form
	// north-latitude,east-longitude,south-latitude,west-longitude.
	// For example, searchRegion=38,-122.1,37.5,-122.5.
	SearchRegion Region `query:"searchRegion"`

	// The location of the user, specified as a comma-separated string that
	// contains the latitude and longitude.
	// For example, userLocation=37.78,-122.42.
	// Search may opt to use the userLocation, if specified, as a fallback for
	// the searchLocation.
	UserLocation Location `query:"userLocation"`
}

func (req *SearchRequest) Validate() error { return vd.Validate(req) }

func (req *SearchRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Add("q", req.Query)
	if len(req.ExcludePoiCategories) > 0 {
		q.Add("excludePoiCategories", poiCategoriesToString(req.ExcludePoiCategories))
	}
	if len(req.IncludePoiCategories) > 0 {
		q.Add("includePoiCategories", poiCategoriesToString(req.IncludePoiCategories))
	}
	if len(req.LimitToCountries) > 0 {
		q.Add("limitToCountries", countriesToString(req.LimitToCountries))
	}
	if len(req.ResultTypeFilter) > 0 {
		q.Add("resultTypeFilter", strings.Join(req.ResultTypeFilter, ","))
	}
	if !req.Lang.IsRoot() {
		q.Set("lang", req.Lang.String())
	}
	if !req.SearchLocation.IsEmpty() {
		q.Set("searchLocation", req.SearchLocation.QueryString())
	}
	if !req.SearchRegion.IsEmpty() {
		q.Set("searchRegion", req.SearchRegion.QueryString())
	}
	if !req.UserLocation.IsEmpty() {
		q.Set("userLocation", req.UserLocation.QueryString())
	}
	return q, nil
}

type SearchAutoCompleteRequest struct {
	// (Required) The query to autocomplete. For example, q=eiffel.
	Query string `query:"q" vd:"$!=''"`

	// A comma-separated list of strings that describes the points of interest
	// to exclude from the search results.
	// For example, excludePoiCategories=Restaurant,Cafe.
	//
	// See PoiCategory for a complete list of possible values.
	ExcludePoiCategories []PoiCategory `query:"excludePoiCategories"`

	// A comma-separated list of strings that describes the points of interest
	// to include in the search results.
	// For example, includePoiCategories=Restaurant,Cafe.
	//
	// See PoiCategory for a complete list of possible values.
	IncludePoiCategories []PoiCategory `query:"includePoiCategories"`

	// A comma-separated list of two-letter ISO 3166-1 codes of the countries to
	// limit the results to.
	// For example, limitToCountries=US,CA limits the search to the United
	// States and Canada.
	//
	// If you specify two or more countries, the results reflect the best
	// available results for some or all of the countries rather than everything
	// related to the query for those countries.
	LimitToCountries []countries.CountryCode `query:"limitToCountries,omitempty"`

	// A comma-separated list of strings that describes the kind of result types
	// to include in the response. For example, resultTypeFilter=Poi.
	//
	// Possible Values: Address, Poi, Query
	ResultTypeFilter []string `query:"resultTypeFilter"`

	// The language the server should use when returning the response, specified
	// using a BCP 47 language code.
	// For example, for English use lang=en-US. Defaults to en-US.
	//
	// Default: en-US
	Lang language.Tag `query:"lang"`

	// A location defined by the application as a hint. Specify the location as
	// a comma-separated string containing the latitude and longitude.
	// For example, searchLocation=37.78,-122.42.
	SearchLocation Location `query:"searchLocation"`

	// A region the app defines as a hint. Specify the region specified as a
	// comma-separated string that describes the region in the form
	// north-latitude,east-longitude,south-latitude,west-longitude.
	// For example, searchRegion=38,-122.1,37.5,-122.5.
	SearchRegion Region `query:"searchRegion"`

	// The location of the user, specified as a comma-separated string that
	// contains the latitude and longitude.
	// For example, userLocation=37.78,-122.42.
	// Search may opt to use the userLocation, if specified, as a fallback for
	// the searchLocation.
	UserLocation Location `query:"userLocation"`
}

func (req *SearchAutoCompleteRequest) Validate() error { return vd.Validate(req) }

func (req *SearchAutoCompleteRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Add("q", req.Query)
	if len(req.ExcludePoiCategories) > 0 {
		q.Add("excludePoiCategories", poiCategoriesToString(req.ExcludePoiCategories))
	}
	if len(req.IncludePoiCategories) > 0 {
		q.Add("includePoiCategories", poiCategoriesToString(req.IncludePoiCategories))
	}
	if len(req.LimitToCountries) > 0 {
		q.Add("limitToCountries", countriesToString(req.LimitToCountries))
	}
	if len(req.ResultTypeFilter) > 0 {
		q.Add("resultTypeFilter", strings.Join(req.ResultTypeFilter, ","))
	}
	if !req.Lang.IsRoot() {
		q.Set("lang", req.Lang.String())
	}
	if !req.SearchLocation.IsEmpty() {
		q.Set("searchLocation", req.SearchLocation.QueryString())
	}
	if !req.SearchRegion.IsEmpty() {
		q.Set("searchRegion", req.SearchRegion.QueryString())
	}
	if !req.UserLocation.IsEmpty() {
		q.Set("userLocation", req.UserLocation.QueryString())
	}
	return q, nil
}

type OneOfLoc struct {
	Address  string
	Location *Location
}

func (loc *OneOfLoc) IsEmpty() bool {
	return loc.Address == "" &&
		(loc.Location == nil || loc.Location.IsEmpty())
}

func (loc *OneOfLoc) QueryString() string {
	if loc.Address != "" {
		return loc.Address
	}
	return loc.Location.QueryString()
}

type DirectionsRequest struct {
	// (Required) The starting location as an address, or coordinates you
	// specify as latitude, longitude. For example, origin=37.7857,-122.4011
	Origin *OneOfLoc `query:"origin" vd:"$!=''"`

	// (Required) The destination as an address, or coordinates you specify as
	// latitude, longitude. For example, destination=San Francisco City Hall, CA
	Destination *OneOfLoc `query:"destination" vd:"$!=''"`

	// The date and time to arrive at the destination in ISO 8601 format in UTC
	// time. For example, 2023-04-15T16:42:00Z.
	//
	// You can specify only arrivalDate or departureDate. If you don’t specify
	// either option, the departureDate defaults to now, which the server
	// interprets as the current time.
	ArrivalDate time.Time `query:"arrivalDate"`

	// A comma-separated list of the features to avoid when calculating
	// direction routes. For example, avoid=Tolls.
	// See DirectionsAvoid for a complete list of possible values.
	Avoid []DirectionsAvoid `query:"avoid"`

	// The date and time to depart from the origin in ISO 8601 format in UTC time.
	// For example, 2023-04-15T16:42:00Z.
	// You can only specify arrivalDate or departureDate. If you don’t specify
	// either option, the departureDate defaults to now, which the server
	// interprets as the current time.
	DepartureDate time.Time `query:"departureDate"`

	// The language the server uses when returning the response, specified using
	// a BCP 47 language code. For example, for English, use lang=en-US.
	// Default: en-US
	Lang language.Tag `query:"lang"`

	// When you set this to true, the server returns additional routes, when
	// available. For example, requestsAlternateRoutes=true.
	// Default: false
	RequestsAlternateRoutes bool `query:"requestsAlternateRoutes"`

	// A searchLocation the app defines as a hint for the query input for origin
	// or destination. Specify the location as a comma-separated string that
	// contains the latitude and longitude. For example, 37.7857,-122.4011.
	// If you don’t provide a searchLocation, the server uses userLocation and
	// searchLocation as fallback hints.
	SearchLocation Location `query:"searchLocation"`

	// A region the app defines as a hint for the query input for origin or
	// destination. Specify the region as a comma-separated string that
	// describes the region in the form of a north-latitude, east-longitude,
	// south-latitude, west-longitude string. For example, 38,-122.1,37.5,-122.5.
	// If you don’t provide a searchLocation, the server uses userLocation and
	// searchRegion as fallback hints.
	SearchRegion Region `query:"searchRegion"`

	// The mode of transportation the server returns directions for.
	// Default: Automobile
	// Possible Values: Automobile, Walking
	TransportType TransportType `query:"transportType"`

	// The location of the user, specified as a comma-separated string that
	// contains the latitude and longitude. For example, userLocation=37.78,-122.42.
	// If you don’t provide a searchLocation, the server uses userLocation and
	// searchRegion as fallback hints.
	UserLocation Location `query:"userLocation"`
}

func (req *DirectionsRequest) Validate() error { return vd.Validate(req) }

func (req *DirectionsRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Origin == nil || req.Origin.IsEmpty() {
		return nil, errors.New("am: origin is required")
	}
	if req.Destination == nil || req.Destination.IsEmpty() {
		return nil, errors.New("am: destination is required")
	}
	q := make(url.Values)
	q.Add("origin", req.Origin.QueryString())
	q.Add("destination", req.Destination.QueryString())
	if !req.ArrivalDate.IsZero() {
		q.Set("arrivalDate", req.ArrivalDate.UTC().Format(time.RFC3339))
	}
	if len(req.Avoid) > 0 {
		q.Add("avoid", avoidToString(req.Avoid))
	}
	if !req.DepartureDate.IsZero() {
		q.Set("departureDate", req.DepartureDate.UTC().Format(time.RFC3339))
	}
	if !req.Lang.IsRoot() {
		q.Set("lang", req.Lang.String())
	}
	if req.RequestsAlternateRoutes {
		q.Set("requestsAlternateRoutes", "true")
	}
	if !req.SearchLocation.IsEmpty() {
		q.Set("searchLocation", req.SearchLocation.QueryString())
	}
	if !req.SearchRegion.IsEmpty() {
		q.Set("searchRegion", req.SearchRegion.QueryString())
	}
	if req.TransportType != "" {
		q.Set("transportType", string(req.TransportType))
	}
	if !req.UserLocation.IsEmpty() {
		q.Set("userLocation", req.UserLocation.QueryString())
	}
	return q, nil
}

// https://developer.apple.com/documentation/applemapsserverapi/determine_estimated_arrival_times_and_distances_to_one_or_more_destinations
type EtasTransportType string

const (
	EtasTransportTypeAutomobile EtasTransportType = "Automobile"
	EtasTransportTypeWalking    EtasTransportType = "Walking"

	// Ask Apple why this not in direction API
	EtasTransportTypeTransit EtasTransportType = "Transit"
)

type EtaRequest struct {
	// (Required) The starting point for estimated arrival time requests,
	// specified as a comma-separated string that contains the latitude and
	// longitude. For example, origin=37.331423,-122.030503.
	Origin Location `query:"origin"`

	// (Required) Destination coordinates represented as pairs of latitude and
	// longitude separated by a vertical bar character (”|”).
	// For example, destinations=37.32556561130194,-121.94635203581443|37.44176585512703,-122.17259315798667.
	// The parameter must specify at least one destination coordinate, but no
	// more than 10 destinations. Specify the location as a comma-separated
	// string that contains the latitude and longitude.
	Destinations []Location `query:"destinations"`

	// The mode of transportation to use when estimating arrival times.
	// Default: Automobile
	// Possible Values: Automobile, Transit, Walking
	TransportType EtasTransportType `query:"transportType"`

	// The time of departure to use in an estimated arrival time request, in ISO
	// 8601 format in UTC time.
	//
	// For example, departureDate=2020-09-15T16:42:00Z.
	// If you don’t specify a departure date, the server uses the current date
	// and time when you make the request.
	DepartureDate time.Time `query:"departureDate"`

	// The intended time of arrival in ISO 8601 format in UTC time.
	ArrivalDate time.Time `query:"arrivalDate"`
}

func (req *EtaRequest) Validate() error {
	if req.Origin.IsEmpty() {
		return errors.New("am: origin is required")
	}
	if len(req.Destinations) == 0 {
		return errors.New("am: destinations is required")
	}
	if len(req.Destinations) > 10 {
		return errors.New("am: destinations max length is 10")
	}
	return vd.Validate(req)
}

func (req *EtaRequest) URLValues() (url.Values, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	q := make(url.Values)
	q.Set("origin", req.Origin.QueryString())
	{
		destinations := []string{}
		for _, dest := range req.Destinations {
			destinations = append(destinations, dest.QueryString())
		}
		q.Set("destinations", strings.Join(destinations, "|"))
	}
	if req.TransportType != "" {
		q.Set("transportType", string(req.TransportType))
	}
	if !req.DepartureDate.IsZero() {
		q.Set("departureDate", req.DepartureDate.UTC().Format(time.RFC3339))
	}
	if !req.ArrivalDate.IsZero() {
		q.Set("arrivalDate", req.ArrivalDate.UTC().Format(time.RFC3339))
	}
	return q, nil
}
