package main

import (
	"context"
	"time"

	"github.com/ringsaturn/am"
	"golang.org/x/text/language"
)

func main() {
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
}
