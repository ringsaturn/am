# `am`: Apple Maps Server API SDK for Go [![Go Reference](https://pkg.go.dev/badge/github.com/ringsaturn/am.svg)](https://pkg.go.dev/github.com/ringsaturn/am)

This Go library provides a comprehensive interface for interacting with various
Apple Maps APIs. It allows users to perform operations like geocoding, reverse
geocoding, searching for places, autocomplete search, getting directions, and
estimating arrival times.

```bash
go install github.com/ringsaturn/am
```

## Features

- **Token Management**: Easily manage your Apple Maps API tokens with built-in
  support for token refreshing and storage.
- **Geocoding**: Convert addresses into geographic coordinates.
- **Reverse Geocoding**: Convert geographic coordinates into human-readable
  addresses.
- **Search**: Search for places that match specific criteria.
- **Autocomplete Search**: Autocomplete search for places based on partial
  input.
- **Directions**: Get directions and estimated travel times between locations.
- **ETA**: Determine estimated arrival times and distances to destinations.

## Installation

To use this library, first ensure you have Go installed on your system. Then,
you can install it using `go get`:

```bash
go get -u github.com/ringsaturn/am
```

## Usage

```go
client := am.NewClient("your_auth_token")
client.Directions(context.Background(), &am.DirectionsRequest{
  Origin: am.OneOfLoc{
    Location: &am.Location{
      Latitude:  37.331871,
      Longitude: -122.029626,
    },
  },
  Destination: am.OneOfLoc{
    Address: "1 Infinite Loop, Cupertino, CA 95014",
  },
  ArrivalDate:             time.Now().Add(time.Hour * 2),
  DepartureDate:           time.Now(),
  Avoid:                   []am.DirectionsAvoid{am.DirectionsAvoidTolls},
  Lang:                    language.AmericanEnglish,
  RequestsAlternateRoutes: true,
  SearchLocation: am.Location{
    Latitude:  37.331871,
    Longitude: -122.029626,
  },
  SearchRegion: am.Region{
    EastLongitude: -122.029626,
    NorthLatitude: 37.331871,
    SouthLatitude: 37.331871,
    WestLongitude: -122.029626,
  },
  UserLocation: am.Location{
    Latitude:  37.331871,
    Longitude: -122.029626,
  },
})
```

For more usage examples, see the
[`client_exmaple_test.go`](./client_exmaple_test.go).

## License

This library is distributed under the [MIT](./LICENSE), see LICENSE for more
details.
