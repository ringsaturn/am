package am_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	am "github.com/ringsaturn/am"
	"golang.org/x/text/language"
)

func ExampleNewClient() {
	simpleClient := am.NewClient("your_auth_token")
	fmt.Println(simpleClient)

	customSaverClient := am.NewClient("your_auth_token", am.WithTokenSaver(&FooTokenSaver{}))
	fmt.Println(customSaverClient)

	customHTTP := &http.Client{
		Timeout: 10,
	}
	customHTTPClient := am.NewClient("your_auth_token", am.WithHTTPClient(customHTTP))
	fmt.Println(customHTTPClient)
}

func ExampleClient_Geocode() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.GeocodeRequest{
		Query: "1600 Amphitheatre Parkway, Mountain View, CA",
	}
	resp, err := client.Geocode(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ExampleClient_ReverseGeocode() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.ReverseRequest{
		Loc: am.Location{
			Latitude:  40.714224,
			Longitude: -73.961452,
		},
		Lang: language.Spanish,
	}
	resp, err := client.ReverseGeocode(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ExampleClient_Search() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.SearchRequest{
		Query: "Eiffel Tower",
		Lang:  language.French,
	}
	resp, err := client.Search(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ExampleClient_SearchAutoComplete() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.SearchAutoCompleteRequest{
		Query: "Eiffel",
		Lang:  language.French,
	}
	resp, err := client.SearchAutoComplete(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ExampleClient_Directions() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.DirectionsRequest{
		Origin:      am.OneOfLoc{Address: "Disneyland"},
		Destination: am.OneOfLoc{Address: "Universal Studios Hollywood"},
		ArrivalDate: time.Now().Add(3 * time.Hour),
	}
	resp, err := client.Directions(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func ExampleClient_Eta() {
	client := am.NewClient("your_auth_token")
	ctx := context.Background()
	req := &am.EtaRequest{
		Origin: am.Location{
			Latitude:  40.714224,
			Longitude: -73.961452,
		},
		Destinations: []am.Location{
			{
				Latitude:  34.138116,
				Longitude: -118.353378,
			},
		},
		TransportType: am.EtasTransportTypeAutomobile,
		DepartureDate: time.Now().Add(1 * time.Hour),
	}
	resp, err := client.Eta(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
