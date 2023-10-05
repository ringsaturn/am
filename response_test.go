package am_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	am "github.com/ringsaturn/am"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/place_results_response_1.json
	expectPlaceResultsResponse1 []byte

	//go:embed testdata/search_response_1.json
	expectSearchResponse1 []byte

	//go:embed testdata/direction_response_1.json
	expectDirectionResponse1 []byte
)

func TestPlaceResultsResponseUnmarshal(t *testing.T) {
	expect := &am.PlaceResults{}
	err := json.Unmarshal(expectPlaceResultsResponse1, expect)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(expect.Results))
	assert.Equal(t, "Apple Park Way", expect.Results[0].Name)
}

func TestSearchResponseUnmarshal(t *testing.T) {
	expect := &am.SearchResponse{}
	err := json.Unmarshal(expectSearchResponse1, expect)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(expect.Results))
	assert.Equal(t, "Eiffel Tower", expect.Results[0].Name)
}

func TestDirectionResponseUnmarshal(t *testing.T) {
	expect := &am.DirectionsResponse{}
	err := json.Unmarshal(expectDirectionResponse1, expect)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(expect.Routes))
	assert.Equal(t, 3, len(expect.StepPaths[1]))

}
