package am

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorFromAPI struct {
	StatusCode int            `json:"statusCode"`
	Response   *ErrorResponse `json:"response"`
	RawBody    []byte         `json:"body"`
	Header     http.Header    `json:"headers"`
}

func (e *ErrorFromAPI) Error() string {
	if e.Response == nil {
		return fmt.Sprintf(
			"am: call API failed(%d/%s) with raw body=`%s`",
			e.StatusCode,
			http.StatusText(e.StatusCode),
			string(e.RawBody),
		)
	}
	if len(e.Response.Error.Details) == 0 {
		return fmt.Sprintf(
			"am: call API failed(%d/%s) with message=%s",
			e.StatusCode,
			http.StatusText(e.StatusCode),
			e.Response.Error.Message,
		)
	}
	return fmt.Sprintf(
		"am: call API failed(%d/%s) with message=%s and details=`%v`",
		e.StatusCode,
		http.StatusText(e.StatusCode),
		e.Response.Error.Message,
		strings.Join(e.Response.Error.Details, ", "),
	)
}

type ErrorResponseError struct {
	Details []string `json:"details"`
	Message string   `json:"message"`
}

// Original response from API.
//
// https://developer.apple.com/documentation/applemapsserverapi/errorresponse
type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

// https://developer.apple.com/documentation/applemapsserverapi/generate_a_maps_access_token
type AccessTokenResponse struct {
	AccessToken      string `json:"accessToken"`
	ExpiresInSeconds int64  `json:"expiresInSeconds"`
}

type MapRegion = Region

// StructuredAddress is 'An object that describes the detailed address components of a place.'
type StructuredAddress struct {
	AdministrativeArea     string   `json:"administrativeArea"`     // The state or province of the place.
	AdministrativeAreaCode string   `json:"administrativeAreaCode"` // The short code for the state or area.
	AreasOfInterest        []string `json:"areasOfInterest"`        // Common names of the area in which the place resides.
	DependentLocalities    []string `json:"dependentLocalities"`    // Common names for the local area or neighborhood of the place.
	FullThoroughfare       string   `json:"fullThoroughfare"`       // A combination of thoroughfare and subthoroughfare.
	Locality               string   `json:"locality"`               // The city of the place.
	PostCode               string   `json:"postCode"`               // The postal code of the place.
	SubLocality            string   `json:"subLocality"`            // The name of the area within the locality.
	SubThoroughfare        string   `json:"subThoroughfare"`        // The number on the street at the place.
	Thoroughfare           string   `json:"thoroughfare"`           // The street name at the place.
}

// https://developer.apple.com/documentation/applemapsserverapi/place
type Place struct {
	// The country or region of the place.
	Country string `json:"country"`

	// The 2-letter country code of the place.
	CountryCode string `json:"countryCode"`

	// The geographic region associated with the place.
	//
	// This is a rectangular region on a map expressed as south-west and
	// north-east points. Specifically south latitude, west longitude,
	// north latitude, and east longitude.
	DisplayMapRegion MapRegion `json:"displayMapRegion"`

	// The address of the place, formatted using its conventions of its country
	// or region.
	FormattedAddressLines []string `json:"formattedAddressLines"`

	// A place name that you can use for display purposes.
	Name string `json:"name"`

	// The latitude and longitude of this place.
	Coordinate Location `json:"coordinate"`

	// A StructuredAddress object that describes details of the place’s address.
	StructuredAddress StructuredAddress `json:"structuredAddress"`
}

type PlaceResults struct {
	Results []Place `json:"results"`
}

type SearchResponse struct {
	DisplayMapRegion MapRegion `json:"displayMapRegion"`
	Results          []Place   `json:"results"`
}

// https://developer.apple.com/documentation/applemapsserverapi/autocompleteresult
type AutocompleteResult struct {
	// The relative URI to the search endpoint to use to fetch more details
	// pertaining to the result. If available, the framework encodes opaque data
	// about the autocomplete result in the completion URL’s metadata parameter.
	// If clients need to fetch the search result in a certain language, they’re
	// responsible for specifying the lang parameter in the request.
	CompletionURL string `json:"completionUrl"`

	// A JSON string array to use to create a long form of display text for the
	// completion result.
	DisplayLines []string `json:"displayLines"`

	// A Location object that specifies the location for the request in terms of
	// its latitude and longitude.
	Location Location `json:"location"`

	// A StructuredAddress object that describes the detailed address components
	// of a place.
	StructuredAddress StructuredAddress `json:"structuredAddress"`
}

// https://developer.apple.com/documentation/applemapsserverapi/searchautocompleteresponse
type SearchAutocompleteResponse struct {
	Results []AutocompleteResult `json:"results"`
}

type DirectionsResponseRoute struct {
	// Total distance that the route covers, in meters.
	DistanceMeters int64 `json:"distanceMeters"`
	// The estimated time to traverse this route in seconds. If you’ve specified
	// a departureDate or arrivalDate, then the estimated time includes traffic
	// conditions assuming user departs or arrives at that time.
	//
	// If you set neither departureDate or arrivalDate, then estimated time
	// represents current traffic conditions assuming user departs “now” from
	// the point of origin.
	DurationSeconds int64 `json:"durationSeconds"`

	// When true, this route has tolls; if false, this route has no tolls.
	// If the value isn’t defined (”undefined”), the route may or may not have tolls.
	HasTolls bool `json:"hasTolls"`

	// The route name that you can use for display purposes.
	Name string `json:"name"`

	// An array of integer values that you can use to determine the number steps
	// along this route. Each value in the array corresponds to an index into
	// the steps array.
	StepIndexes []int64 `json:"stepIndexes"`

	// A string that represents the mode of transportation the service used to
	// estimate the arrival time. Same as the input query param transportType or
	// Automobile if the input query didn’t specify a transportation type.
	// Possible Values: Automobile, Walking
	TransportType TransportType `json:"transportType"`
}

type DirectionsResponseStep struct {
	// Total distance covered by the step, in meters.
	DistanceMeters int64 `json:"distanceMeters"`

	// The estimated time to traverse this step, in seconds.
	DurationSeconds int64 `json:"durationSeconds"`

	// The localized instruction string for this step that you can use for display purposes.
	// You can specify the language to receive the response in using the lang parameter.
	Instructions string `json:"instructions"`

	// A pointer to this step’s path. The pointer is in the form of an index
	// into the stepPaths array contained in a Route.
	// Step paths are self-contained which implies that the last point of a
	// previous step path along a route is the same as the first point of the
	// next step path. Clients are responsible for avoiding duplication when
	// rendering the point.
	StepPathIndex int64 `json:"stepPathIndex"`

	// A string indicating the transport type for this step if it’s different
	// from the transportType in the route.
	// Possible Values: Automobile, Walking
	TransportType TransportType `json:"transportType"`
}

// https://developer.apple.com/documentation/applemapsserverapi/directionsresponse
type DirectionsResponse struct {
	// A Place object that describes the destination.
	Destination Place `json:"destination"`

	// A Place result that describes the origin.
	Origin Place `json:"origin"`

	// An array of routes. Each route references steps based on indexes into the steps array.
	Routes []DirectionsResponseRoute `json:"routes"`

	// An array of step paths across all steps across all routes. Each step path
	// is a single polyline represented as an array of points. You reference the
	// step paths by index into the array.
	StepPaths [][]Location `json:"stepPaths"`

	// An array of all steps across all routes. You reference the route steps by
	// index into this array. Each step in turn references its path based on
	// indexes into the stepPaths array.
	Steps []DirectionsResponseStep `json:"steps"`
}

// https://developer.apple.com/documentation/applemapsserverapi/etaresponse/eta
type EtaResponseEta struct {
	// The destination as a Location.
	Destination Location `json:"destination"`

	// The distance in meters to the destination.
	DistanceMeters int64 `json:"distanceMeters"`

	// The estimated travel time in seconds, including delays due to traffic.
	ExpectedTravelTimeSeconds int64 `json:"expectedTravelTimeSeconds"`

	// The expected travel time, in seconds, without traffic.
	StaticTravelTimeSeconds int64 `json:"staticTravelTimeSeconds"`

	// A string that represents the mode of transportation for this ETA, which is one of:
	// Automobile, Walking, Transit
	//
	// NOTE(ringsaturn): Apple's example use `AUTOMOBILE`, who knows why.
	TransportType TransportType `json:"transportType"`
}

// https://developer.apple.com/documentation/applemapsserverapi/etaresponse
type EtaResponse struct {
	// An array of one or more EtaResponse.Eta objects.
	Etas []EtaResponseEta `json:"etas"`
}
