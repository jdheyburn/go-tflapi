package tfl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	appID  = "APP_ID"
	appKey = "APP_KEY"
)

var (
	server *httptest.Server
	client *TflClient
)

func getTestDataFileContents(fname string) []byte {
	fpath := filepath.Join("testdata", fname)
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	return b
}

func TflAPIClientStub() func() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp []byte

		switch r.URL.RequestURI() {
		case fmt.Sprintf("/StopPoint/9100ECROYDN?app_id=%s&app_key=%s", appID, appKey):
			resp = getTestDataFileContents("Should_retrieve_StopPoint_given_valid_ID.json")
		case fmt.Sprintf("/StopPoint/INVALID?app_id=%s&app_key=%s", appID, appKey):
			resp = getTestDataFileContents("Should_handle_response_for_invalid_ID.json")
			w.WriteHeader(http.StatusNotFound)
		case fmt.Sprintf("/StopPoint/Search/%s?app_id=%s&app_key=%s", "London%20Bridge", appID, appKey):
			resp = getTestDataFileContents("Should_retrieve_Search_Reponses_given_valid_ID.json")
		case fmt.Sprintf("/StopPoint/Search/%s?app_id=%s&app_key=%s&modes=%s", "London%20Bridge", appID, appKey, "national-rail%2Ctube"):
			resp = getTestDataFileContents("Should_retrieve_filtered_Search_Reponses_given_valid_ID.json")
		case fmt.Sprintf("/StopPoint/Search/Nope?app_id=%s&app_key=%s", appID, appKey):
			resp = getTestDataFileContents("Should_retrieve_no_matches_for_invalid_searchTerm.json")
		case fmt.Sprintf("/Journey/JourneyResults/1001089/to/1000173?app_id=%s&app_key=%s&date=20190401&mode=%s&time=0715", appID, appKey, "national-rail%2Ctube"):
			resp = getTestDataFileContents("Should_retrieve_journey_planner_itinerary_for_valid_search.json")
		}

		w.Write(resp)
	}))

	client, _ = New(
		WithBaseURL(server.URL),
		WithAppID(appID),
		WithAppKey(appKey),
	)

	return func() {
		server.Close()
	}
}

func TestMain(m *testing.M) {
	teardown := TflAPIClientStub()
	defer teardown()
	os.Exit(m.Run())
}

func TestTflAPIClient_buildURL(t *testing.T) {
	type args struct {
		pathParams []string
	}
	tests := []struct {
		name string
		api  *TflClient
		args args
		want string
	}{
		{
			name: "Should build URL correctly encoded with no queryParams",
			api:  client,
			args: args{
				pathParams: []string{
					"StopPoint", "Search", "London Bridge",
				},
			},
			want: fmt.Sprintf("%s/StopPoint/Search/%s?app_id=%s&app_key=%s", server.URL, "London%20Bridge", appID, appKey),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.api.buildURL(tt.args.pathParams); got != tt.want {
				t.Errorf("TflAPIClient.buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTflAPIClient_buildURLWithQueryParams(t *testing.T) {
	type args struct {
		pathParams  []string
		queryParams *map[string]string
	}
	tests := []struct {
		name string
		api  *TflClient
		args args
		want string
	}{
		{
			name: "Should build URL correctly with queryParams",
			api:  client,
			args: args{
				pathParams: []string{
					"StopPoint", "Search",
				},
				queryParams: &map[string]string{
					"queryParam":   "one",
					"aBeforeAppID": "two",
				},
			},
			want: fmt.Sprintf("%s/StopPoint/Search?aBeforeAppID=two&app_id=%s&app_key=%s&queryParam=one",
				server.URL, appID, appKey),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.api.buildURLWithQueryParams(tt.args.pathParams, tt.args.queryParams); got != tt.want {
				t.Errorf("TflAPIClient.buildURLWithQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTflAPIClient_GetStopPointForID(t *testing.T) {

	expected := StopPointAPIResponse{}
	json.Unmarshal(getTestDataFileContents("Should_retrieve_StopPoint_given_valid_ID.json"), &expected)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		api     *TflClient
		args    args
		want    *StopPointAPIResponse
		wantErr error
	}{
		{
			name: "Should retrieve StopPoint given valid ID",
			api:  client,
			args: args{
				id: "9100ECROYDN",
			},
			want: &expected,
		},
		{
			name: "Should handle response for invalid ID",
			api:  client,
			args: args{
				id: "INVALID",
			},
			wantErr: errors.New("The following stop point is not recognised: INVALID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.api.GetStopPointForID(tt.args.id)
			if tt.wantErr != nil && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TflAPIClient.GetStopPointForID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TflAPIClient.GetStopPointForID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTflAPIClient_SearchStopPoints(t *testing.T) {

	response := EntitySearchResponse{}
	json.Unmarshal(getTestDataFileContents("Should_retrieve_Search_Reponses_given_valid_ID.json"), &response)
	expected := response.Matches

	type args struct {
		searchTerm string
	}
	tests := []struct {
		name    string
		api     *TflClient
		args    args
		want    *[]EntityMatchedStop
		wantErr error
	}{
		{
			name: "Should retrieve matches for valid searchTerm",
			api:  client,
			args: args{
				searchTerm: "London Bridge",
			},
			want: &expected,
		},
		{
			name: "Should retrieve no matches for invalid searchTerm",
			api:  client,
			args: args{
				searchTerm: "Nope",
			},
			want: &[]EntityMatchedStop{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.api.SearchStopPoints(tt.args.searchTerm)
			if tt.wantErr != nil && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TflAPIClient.SearchStopPoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TflAPIClient.SearchStopPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTflAPIClient_SearchStopPointsWithModes(t *testing.T) {

	response := EntitySearchResponse{}
	json.Unmarshal(getTestDataFileContents("Should_retrieve_filtered_Search_Reponses_given_valid_ID.json"), &response)
	expected := response.Matches

	type args struct {
		searchTerm string
		modes      []string
	}
	tests := []struct {
		name    string
		api     *TflClient
		args    args
		want    *[]EntityMatchedStop
		wantErr error
	}{
		{
			name: "Should retrieve matches for valid searchTerm and filter by modes",
			api:  client,
			args: args{
				searchTerm: "London Bridge",
				modes: []string{"national-rail", "tube"},
			},
			want: &expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.api.SearchStopPointsWithModes(tt.args.searchTerm, tt.args.modes)
			if tt.wantErr != nil && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TflAPIClient.SearchStopPointsWithModes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TflAPIClient.SearchStopPointsWithModes() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestTflAPIClient_GetJourneyPlannerItinerary(t *testing.T) {

	expected := JourneyPlannerItineraryResult{}
	json.Unmarshal(getTestDataFileContents("Should_retrieve_journey_planner_itinerary_for_valid_search.json"), &expected)

	type args struct {
		query JourneyPlannerQuery
	}
	tests := []struct {
		name    string
		api     *TflClient
		args    args
		want    *JourneyPlannerItineraryResult
		wantErr error
	}{
		{
			name: "Should retrieve journey planner itinerary for valid search",
			api:  client,
			args: args{
				query: JourneyPlannerQuery{
					From:  "1001089",
					To:    "1000173",
					Date:  "20190401",
					Time:  "0715",
					Modes: []string{"national-rail", "tube"},
				},
			},
			want: &expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.api.GetJourneyPlannerItinerary(tt.args.query)
			if tt.wantErr != nil && !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TflAPIClient.GetJourneyPlannerItinerary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TflAPIClient.GetJourneyPlannerItinerary() = %v, want %v", got, tt.want)
			}
		})
	}
}
