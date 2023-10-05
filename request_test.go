package am_test

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/biter777/countries"
	am "github.com/ringsaturn/am"
	"golang.org/x/text/language"
)

func TestGeocodeRequest_URLValues(t *testing.T) {
	type fields struct {
		Query            string
		LimitToCountries []countries.CountryCode
		Lang             language.Tag
		SearchLocation   am.Location
		SearchRegion     am.Region
		UserLocation     am.Location
	}
	tests := []struct {
		name    string
		fields  fields
		want    url.Values
		wantErr bool
	}{
		{
			name: "hello",
			fields: fields{
				Query: "hello",
			},
			want: url.Values{
				"q": []string{"hello"},
			},
			wantErr: false,
		},
		{
			name: "US,CA",
			fields: fields{
				Query: "hello",
				LimitToCountries: []countries.CountryCode{
					countries.USA,
					countries.Canada,
				},
			},
			want: url.Values{
				"q":                []string{"hello"},
				"limitToCountries": []string{"US,CA"},
			},
		},
		{
			name: "Shinjuku",
			fields: fields{
				Query: "Shinjuku City",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.English,
			},
			want: url.Values{
				"q":                []string{"Shinjuku City"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en"},
			},
			wantErr: false,
		},
		{
			name: "KyoAni",
			fields: fields{
				Query: "KyoAni",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.AmericanEnglish,
				SearchLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
				SearchRegion: am.Region{
					NorthLatitude: 35.0219,
					EastLongitude: 135.8426,
					SouthLatitude: 34.8440,
					WestLongitude: 135.6215,
				},
				UserLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
			},
			want: url.Values{
				"q":                []string{"KyoAni"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en-US"},
				"searchLocation":   []string{"34.985849,135.7561864"},
				"searchRegion":     []string{"35.0219,135.8426,34.844,135.6215"},
				"userLocation":     []string{"34.985849,135.7561864"},
			},
			wantErr: false,
		},
		{
			name: "Empty",
			fields: fields{
				Query: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &am.GeocodeRequest{
				Query:            tt.fields.Query,
				LimitToCountries: tt.fields.LimitToCountries,
				Lang:             tt.fields.Lang,
				SearchLocation:   tt.fields.SearchLocation,
				SearchRegion:     tt.fields.SearchRegion,
				UserLocation:     tt.fields.UserLocation,
			}
			got, err := req.URLValues()
			if (err != nil) != tt.wantErr {
				t.Errorf("GeocodeRequest.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeocodeRequest.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseRequest_URLValues(t *testing.T) {
	type fields struct {
		Loc  am.Location
		Lang language.Tag
	}
	tests := []struct {
		name    string
		fields  fields
		want    url.Values
		wantErr bool
	}{
		{
			name: "hello",
			fields: fields{
				Loc: am.Location{
					Latitude:  37.33182,
					Longitude: -122.03118,
				},
				Lang: language.AmericanEnglish,
			},
			want: url.Values{
				"loc":  []string{"37.33182,-122.03118"},
				"lang": []string{"en-US"},
			},
			wantErr: false,
		},
		{
			name: "no lang",
			fields: fields{
				Loc: am.Location{
					Latitude:  37.33182,
					Longitude: -122.03118,
				},
			},
			want: url.Values{
				"loc": []string{"37.33182,-122.03118"},
			},
			wantErr: false,
		},
		{
			name: "no loc",
			fields: fields{
				Lang: language.AmericanEnglish,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &am.ReverseRequest{
				Loc:  tt.fields.Loc,
				Lang: tt.fields.Lang,
			}
			got, err := req.URLValues()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReverseRequest.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseRequest.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchRequest_URLValues(t *testing.T) {
	type fields struct {
		Query                string
		ExcludePoiCategories []am.PoiCategory
		IncludePoiCategories []am.PoiCategory
		LimitToCountries     []countries.CountryCode
		ResultTypeFilter     []string
		Lang                 language.Tag
		SearchLocation       am.Location
		SearchRegion         am.Region
		UserLocation         am.Location
	}
	tests := []struct {
		name    string
		fields  fields
		want    url.Values
		wantErr bool
	}{
		{
			name: "hello",
			fields: fields{
				Query: "hello",
			},
			want: url.Values{
				"q": []string{"hello"},
			},
			wantErr: false,
		},
		{
			name: "US,CA",
			fields: fields{
				Query: "hello",
				LimitToCountries: []countries.CountryCode{
					countries.USA,
					countries.Canada,
				},
			},
			want: url.Values{
				"q":                []string{"hello"},
				"limitToCountries": []string{"US,CA"},
			},
		},
		{
			name: "Shinjuku",
			fields: fields{
				Query: "Shinjuku City",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.English,
			},
			want: url.Values{
				"q":                []string{"Shinjuku City"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en"},
			},
			wantErr: false,
		},
		{
			name: "KyoAni",
			fields: fields{
				Query: "KyoAni",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.AmericanEnglish,
				SearchLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
				SearchRegion: am.Region{
					NorthLatitude: 35.0219,
					EastLongitude: 135.8426,
					SouthLatitude: 34.8440,
					WestLongitude: 135.6215,
				},
				UserLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
			},
			want: url.Values{
				"q":                []string{"KyoAni"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en-US"},
				"searchLocation":   []string{"34.985849,135.7561864"},
				"searchRegion":     []string{"35.0219,135.8426,34.844,135.6215"},
				"userLocation":     []string{"34.985849,135.7561864"},
			},
			wantErr: false,
		},
		{
			name: "Empty",
			fields: fields{
				Query: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &am.SearchRequest{
				Query:                tt.fields.Query,
				ExcludePoiCategories: tt.fields.ExcludePoiCategories,
				IncludePoiCategories: tt.fields.IncludePoiCategories,
				LimitToCountries:     tt.fields.LimitToCountries,
				ResultTypeFilter:     tt.fields.ResultTypeFilter,
				Lang:                 tt.fields.Lang,
				SearchLocation:       tt.fields.SearchLocation,
				SearchRegion:         tt.fields.SearchRegion,
				UserLocation:         tt.fields.UserLocation,
			}
			got, err := req.URLValues()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchRequest.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchRequest.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchAutoCompleteRequest_URLValues(t *testing.T) {
	type fields struct {
		Query                string
		ExcludePoiCategories []am.PoiCategory
		IncludePoiCategories []am.PoiCategory
		LimitToCountries     []countries.CountryCode
		ResultTypeFilter     []string
		Lang                 language.Tag
		SearchLocation       am.Location
		SearchRegion         am.Region
		UserLocation         am.Location
	}
	tests := []struct {
		name    string
		fields  fields
		want    url.Values
		wantErr bool
	}{
		{
			name: "hello",
			fields: fields{
				Query: "hello",
			},
			want: url.Values{
				"q": []string{"hello"},
			},
			wantErr: false,
		},
		{
			name: "US,CA",
			fields: fields{
				Query: "hello",
				LimitToCountries: []countries.CountryCode{
					countries.USA,
					countries.Canada,
				},
			},
			want: url.Values{
				"q":                []string{"hello"},
				"limitToCountries": []string{"US,CA"},
			},
		},
		{
			name: "Shinjuku",
			fields: fields{
				Query: "Shinjuku City",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.English,
			},
			want: url.Values{
				"q":                []string{"Shinjuku City"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en"},
			},
			wantErr: false,
		},
		{
			name: "KyoAni",
			fields: fields{
				Query: "KyoAni",
				LimitToCountries: []countries.CountryCode{
					countries.Japan,
				},
				Lang: language.AmericanEnglish,
				SearchLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
				SearchRegion: am.Region{
					NorthLatitude: 35.0219,
					EastLongitude: 135.8426,
					SouthLatitude: 34.8440,
					WestLongitude: 135.6215,
				},
				UserLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
			},
			want: url.Values{
				"q":                []string{"KyoAni"},
				"limitToCountries": []string{"JP"},
				"lang":             []string{"en-US"},
				"searchLocation":   []string{"34.985849,135.7561864"},
				"searchRegion":     []string{"35.0219,135.8426,34.844,135.6215"},
				"userLocation":     []string{"34.985849,135.7561864"},
			},
			wantErr: false,
		},
		{
			name: "Empty",
			fields: fields{
				Query: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &am.SearchAutoCompleteRequest{
				Query:                tt.fields.Query,
				ExcludePoiCategories: tt.fields.ExcludePoiCategories,
				IncludePoiCategories: tt.fields.IncludePoiCategories,
				LimitToCountries:     tt.fields.LimitToCountries,
				ResultTypeFilter:     tt.fields.ResultTypeFilter,
				Lang:                 tt.fields.Lang,
				SearchLocation:       tt.fields.SearchLocation,
				SearchRegion:         tt.fields.SearchRegion,
				UserLocation:         tt.fields.UserLocation,
			}
			got, err := req.URLValues()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchRequest.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchRequest.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectionsRequest_URLValues(t *testing.T) {
	type fields struct {
		Origin                  *am.OneOfLoc
		Destination             *am.OneOfLoc
		ArrivalDate             time.Time
		Avoid                   []am.DirectionsAvoid
		DepartureDate           time.Time
		Lang                    language.Tag
		RequestsAlternateRoutes bool
		SearchLocation          am.Location
		SearchRegion            am.Region
		TransportType           am.TransportType
		UserLocation            am.Location
	}
	tests := []struct {
		name    string
		fields  fields
		want    url.Values
		wantErr bool
	}{
		{
			name: "hello",
			fields: fields{
				Origin: &am.OneOfLoc{
					Location: &am.Location{
						Latitude:  37.33182,
						Longitude: -122.03118,
					},
				},
				Destination: &am.OneOfLoc{
					Address: "NYC",
				},
				ArrivalDate:   time.Unix(1696484859, 0),
				Avoid:         []am.DirectionsAvoid{am.DirectionsAvoidTolls},
				DepartureDate: time.Unix(1696484859, 0).Add(time.Hour * 2),
				Lang:          language.AmericanEnglish,
				SearchLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
				SearchRegion: am.Region{
					NorthLatitude: 35.0219,
					EastLongitude: 135.8426,
					SouthLatitude: 34.8440,
					WestLongitude: 135.6215,
				},
				TransportType: am.TransportTypeAutomobile,
				UserLocation: am.Location{
					Latitude:  34.985849,
					Longitude: 135.7561864,
				},
			},
			want: url.Values{
				"origin":         []string{"37.33182,-122.03118"},
				"destination":    []string{"NYC"},
				"arrivalDate":    []string{"2023-10-05T05:47:39Z"},
				"avoid":          []string{"Tolls"},
				"departureDate":  []string{"2023-10-05T07:47:39Z"},
				"lang":           []string{"en-US"},
				"searchLocation": []string{"34.985849,135.7561864"},
				"searchRegion":   []string{"35.0219,135.8426,34.844,135.6215"},
				"transportType":  []string{"Automobile"},
				"userLocation":   []string{"34.985849,135.7561864"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &am.DirectionsRequest{
				Origin:                  tt.fields.Origin,
				Destination:             tt.fields.Destination,
				ArrivalDate:             tt.fields.ArrivalDate,
				Avoid:                   tt.fields.Avoid,
				DepartureDate:           tt.fields.DepartureDate,
				Lang:                    tt.fields.Lang,
				RequestsAlternateRoutes: tt.fields.RequestsAlternateRoutes,
				SearchLocation:          tt.fields.SearchLocation,
				SearchRegion:            tt.fields.SearchRegion,
				TransportType:           tt.fields.TransportType,
				UserLocation:            tt.fields.UserLocation,
			}
			got, err := req.URLValues()
			if (err != nil) != tt.wantErr {
				t.Errorf("DirectionsRequest.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DirectionsRequest.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
