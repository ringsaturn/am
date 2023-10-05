package am

import (
	"strconv"
	"strings"
)

type Location struct {
	Latitude  float64 `json:"latitude" vd:"$>=-90 && $<=90"`
	Longitude float64 `json:"longitude" vd:"$>=-180 && $<=180"`
}

func (Location Location) IsEmpty() bool {
	return Location.Latitude == 0 && Location.Longitude == 0
}

// searchLocation=37.78,-122.42.
func (location Location) QueryString() string {
	return strings.Join([]string{
		strconv.FormatFloat(location.Latitude, 'f', -1, 64),
		strconv.FormatFloat(location.Longitude, 'f', -1, 64),
	}, ",")
}

type Region struct {
	EastLongitude float64 `json:"eastLongitude" vd:"$>=-180 && $<=180"`
	NorthLatitude float64 `json:"northLatitude" vd:"$>=-90 && $<=90"`
	SouthLatitude float64 `json:"southLatitude" vd:"$>=-90 && $<=90"`
	WestLongitude float64 `json:"westLongitude" vd:"$>=-180 && $<=180"`
}

func (region Region) IsEmpty() bool {
	return (region.EastLongitude == 0 && region.NorthLatitude == 0 &&
		region.SouthLatitude == 0 && region.WestLongitude == 0)
}

// Region order as: north-latitude, east-longitude, south-latitude, west-longitude.
// For example, searchRegion=38,-122.1,37.5,-122.5.
func (region Region) QueryString() string {
	return strings.Join([]string{
		strconv.FormatFloat(region.NorthLatitude, 'f', -1, 64),
		strconv.FormatFloat(region.EastLongitude, 'f', -1, 64),
		strconv.FormatFloat(region.SouthLatitude, 'f', -1, 64),
		strconv.FormatFloat(region.WestLongitude, 'f', -1, 64),
	}, ",")
}

// PoiCategory: 'A string that describes a specific point of interest (POI) category.'
// https://developer.apple.com/documentation/applemapsserverapi/poicategory
type PoiCategory string

// Codes generated via ChatGPT
const (
	Airport         PoiCategory = "Airport"         // An airport.
	AirportGate     PoiCategory = "AirportGate"     // A specific gate at an airport.
	AirportTerminal PoiCategory = "AirportTerminal" // A specific named terminal at an airport.
	AmusementPark   PoiCategory = "AmusementPark"   // An amusement park.
	ATM             PoiCategory = "ATM"             // An automated teller machine.
	Aquarium        PoiCategory = "Aquarium"        // An aquarium.
	Bakery          PoiCategory = "Bakery"          // A bakery.
	Bank            PoiCategory = "Bank"            // A bank.
	Beach           PoiCategory = "Beach"           // A beach.
	Brewery         PoiCategory = "Brewery"         // A brewery.
	Cafe            PoiCategory = "Cafe"            // A cafe.
	Campground      PoiCategory = "Campground"      // A campground.
	CarRental       PoiCategory = "CarRental"       // A car rental location.
	EVCharger       PoiCategory = "EVCharger"       // An electric vehicle (EV) charger.
	FireStation     PoiCategory = "FireStation"     // A fire station.
	FitnessCenter   PoiCategory = "FitnessCenter"   // A fitness center.
	FoodMarket      PoiCategory = "FoodMarket"      // A food market.
	GasStation      PoiCategory = "GasStation"      // A gas station.
	Hospital        PoiCategory = "Hospital"        // A hospital.
	Hotel           PoiCategory = "Hotel"           // A hotel.
	Laundry         PoiCategory = "Laundry"         // A laundry.
	Library         PoiCategory = "Library"         // A library.
	Marina          PoiCategory = "Marina"          // A marina.
	MovieTheater    PoiCategory = "MovieTheater"    // A movie theater.
	Museum          PoiCategory = "Museum"          // A museum.
	NationalPark    PoiCategory = "NationalPark"    // A national park.
	Nightlife       PoiCategory = "Nightlife"       // A night life venue.
	Park            PoiCategory = "Park"            // A park.
	Parking         PoiCategory = "Parking"         // A parking location for an automobile.
	Pharmacy        PoiCategory = "Pharmacy"        // A pharmacy.
	Playground      PoiCategory = "Playground"      // A playground.
	Police          PoiCategory = "Police"          // A police station.
	PostOffice      PoiCategory = "PostOffice"      // A post office.
	PublicTransport PoiCategory = "PublicTransport" // A public transportation station.
	ReligiousSite   PoiCategory = "ReligiousSite"   // A religious site.
	Restaurant      PoiCategory = "Restaurant"      // A restaurant.
	Restroom        PoiCategory = "Restroom"        // A restroom.
	School          PoiCategory = "School"          // A school.
	Stadium         PoiCategory = "Stadium"         // A stadium.
	Store           PoiCategory = "Store"           // A store.
	Theater         PoiCategory = "Theater"         // A theater.
	University      PoiCategory = "University"      // A university.
	Winery          PoiCategory = "Winery"          // A winery.
	Zoo             PoiCategory = "Zoo"             // A zoo.
)

type DirectionsAvoid string

const (
	// When you set avoid=Tolls, routes without tolls are higher up in the list
	// of returned routes. Note that even when you set avoid=Tolls, the routes
	// the server returns may contain tolls (if no reasonable toll-free routes
	// exist). Ensure you check the value of routes[i].hasTolls in the response
	// to verify toll assumptions.
	DirectionsAvoidTolls DirectionsAvoid = "Tolls"
)

// The mode of transportation the server returns directions for.
type TransportType string

const (
	TransportTypeAutomobile TransportType = "Automobile"
	TransportTypeWalking    TransportType = "Walking"
)
