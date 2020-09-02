package tfl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiURL             string = "https://api.tfl.gov.uk"
	journeyResultsPath string = "Journey/JourneyResults"
	toPath             string = "to"
	stopPointPath      string = "StopPoint"
	searchPath         string = "Search"
	fareToPath         string = "FareTo"
)

// Option is a functional option for configuring the API client
type Option func(*TflClient) error

// WithBaseURL allows overriding of API client baseURL for testing
func WithBaseURL(baseURL string) Option {
	return func(c *TflClient) error {
		parsedURL, err := url.Parse(baseURL)
		c.baseURL = parsedURL
		return err
	}
}

// WithAppID sets the appID for the API
func WithAppID(appID string) Option {
	return func(c *TflClient) error {
		c.appID = appID
		return nil
	}
}

// WithAppKey sets the appKey for the API
func WithAppKey(appKey string) Option {
	return func(c *TflClient) error {
		c.appKey = appKey
		return nil
	}
}

func (c *TflClient) parseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}
	return nil
}

type Api interface {
	SearchStopPoints(string) (*[]EntityMatchedStop, error)
	SearchStopPointsWithModes(string, []string) (*[]EntityMatchedStop, error)
	GetStopPointForID(string) (*StopPointAPIResponse, error)
	GetJourneyPlannerItinerary(JourneyPlannerQuery) (*JourneyPlannerItineraryResult, error)
	SingleFareFinder(SingleFareFinderInput) (*[]FaresSection, error)
}

// Client holds information necessary to make a request to your API
type TflClient struct {
	Client  *http.Client
	baseURL *url.URL
	appID   string
	appKey  string
}

// New returns a new instance of the Client
func New(opts ...Option) (*TflClient, error) {
	parsedURL, _ := url.Parse(apiURL)

	c := &TflClient{
		baseURL: parsedURL,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	if err := c.parseOptions(opts...); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *TflClient) buildURL(pathParams []string) string {
	queryParams := map[string]string{}
	return c.buildURLWithQueryParams(pathParams, &queryParams)
}

func (c *TflClient) buildURLWithQueryParams(pathParams []string, queryParams *map[string]string) string {

	builtURL := c.baseURL
	builtURL.Path = strings.Join(pathParams, "/")

	params := url.Values{}
	for key, val := range *queryParams {
		params.Add(key, val)
	}
	params.Add("app_id", c.appID)
	params.Add("app_key", c.appKey)

	builtURL.RawQuery = params.Encode()
	return builtURL.String()
}

// getJSON wraps around the client to execute the GET request and maps the result to the provided interface type
// Also handles a non-OK response from the API and extracts the error if so
func (c *TflClient) getJSON(url string, respObj interface{}) error {

	fmt.Printf("GET - %s\n", url)
	resp, err := c.Client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errObj := APIErrorResponse{}
		if err := serialiseResponse(resp, &errObj); err != nil {
			return err
		}
		return errors.New(errObj.Message)
	}

	return serialiseResponse(resp, &respObj)
}

// serialiseResponse takes in a response and attempts to serialise it to the provided interface
func serialiseResponse(resp *http.Response, obj interface{}) error {

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, &obj)
}

// GetStopPointForID retrieves the StopPoint information for a given ID
// It queries the endpoint /StopPoint/{id}
func (c *TflClient) GetStopPointForID(id string) (*StopPointAPIResponse, error) {

	pathParams := []string{stopPointPath, id}
	url := c.buildURL(pathParams)

	resp := StopPointAPIResponse{}
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SearchStopPoints retrieves MatchedStops for a given search term
// It queries the endpoint /StopPoint/Search/{searchTerm}
func (c *TflClient) SearchStopPoints(searchTerm string) (*[]EntityMatchedStop, error) {
	return c.SearchStopPointsWithModes(searchTerm, []string{})
}

// SearchStopPointsWithModes retrieves MatchedStops for a given search term, filtered against StopPoint mode
// It queries the endpoint /StopPoint/Search/{searchTerm}
func (c *TflClient) SearchStopPointsWithModes(searchTerm string, modes []string) (*[]EntityMatchedStop, error) {

	// TODO validate query:
	// - searchTerm mustn't be bad
	// - modes must be from valid list
	pathParams := []string{stopPointPath, searchPath, searchTerm}
	var url string

	if modes != nil && len(modes) > 0 {
		queryParams := &map[string]string{
			"modes": strings.Join(modes, ","),
		}
		url = c.buildURLWithQueryParams(pathParams, queryParams)
	} else {
		url = c.buildURL(pathParams)
	}

	resp := EntitySearchResponse{}
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}

	return &resp.Matches, nil
}

// GetJourneyPlannerItinerary retrieves MatchedStops for a given search term
// It queries the endpoint /Journey/JourneyResult/{from}/to/{to}
func (c *TflClient) GetJourneyPlannerItinerary(query JourneyPlannerQuery) (*JourneyPlannerItineraryResult, error) {

	pathParams := []string{journeyResultsPath, query.From, toPath, query.To}
	// TODO validate query:
	// - date and time are mandatory
	// - modes must be from valid list
	queryParams := &map[string]string{
		"date": query.Date,
		"time": query.Time,
		// "mode": strings.Join(query.modes, ","),
	}
	if query.Modes != nil && len(query.Modes) > 0 {
		(*queryParams)["mode"] = strings.Join(query.Modes, ",")
	}
	url := c.buildURLWithQueryParams(pathParams, queryParams)

	resp := JourneyPlannerItineraryResult{}
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SingleFareFinder retrieves a single fare cost between two stations
// It queries the endpoint /StopPoint/{from}/FareTo/{to}
func (c *TflClient) SingleFareFinder(input SingleFareFinderInput) (*[]FaresSection, error) {

	pathParams := []string{stopPointPath, input.From, fareToPath, input.To}
	queryParams := &map[string]string{}
	url := c.buildURLWithQueryParams(pathParams, queryParams)

	resp := []FaresSection{}
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
