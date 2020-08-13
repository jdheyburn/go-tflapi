package tfl

// APIErrorResponse represents Tfl.Api.Presentation.Entities.ApiError
type APIErrorResponse struct {
	TimestampUTC   string `json:"timestampUTC"`
	ExceptionType  string `json:"exceptionType"`
	HTTPStatusCode uint16 `json:"httpStatusCode"`
	HTTPStatus     string `json:"httpStatus"`
	RelativeURI    string `json:"relativeUri"`
	Message        string `json:"message"`
}

// StopPointAPIResponse represents Tfl.Api.Presentation.Entities.StopPoint
type StopPointAPIResponse struct {
	// RespType             string           `json:"$type"`
	NaptanID string   `json:"naptanId"`
	Modes    []string `json:"modes"`
	IcsCode  string   `json:"icsCode"`
	StopType string   `json:"stopType"`
	// Lines                []LineIdentifier `json:"lines"`          // TODO define Lines
	// LineGroup            string           `json:"lineGroup"`      // TODO define LineGroup
	// LineModeGroups       string           `json:"lineModeGroups"` // TODO define LineModeGroup
	Status               bool                   `json:"status"`
	ID                   string                 `json:"id"`
	StationNaptan        string                 `json:"stationNaptan"`
	HubNaptanCode        string                 `json:"hubNaptanCode"`
	CommonName           string                 `json:"commonName"`
	PlaceType            string                 `json:"placeType"`
	AdditionalProperties []AdditionalProperties `json:"additionalProperties"`
	Children             []StopPointAPIResponse `json:"children"`
	// Lat                  float64          `json:"lat"`
	// Lon                  float64          `json:"lon"`
}

type AdditionalProperties struct {
	RespType        string `json:"$type"`
	Category        string `json:"category"`
	Key             string `json:"key"`
	SourceSystemKey string `json:"sourceSystemKey"`
	Value           string `json:"value"`
}

// LineIdentifier represents
type LineIdentifier struct {
	RespType  string `json:"$type"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
	Type      string `json:"type"`
	Crowding  string `json:"crowding"` // TODO define this as a struct
	RouteType string `json:"routeType"`
	Status    string `json:"status"`
}

// EntityMatchedStop represents Tfl.Api.Presentation.Entities.MatchedStop
type EntityMatchedStop struct {
	Modes   []string `json:"modes"`
	IcsCode string   `json:"icsId"`
	Name    string   `json:"name"`
	Zone    string   `json:"zone"`
	ID      string   `json:"id"`
}

// EntitySearchResponse represents Tfl.Api.Presentation.Entities.SearchResponse
type EntitySearchResponse struct {
	Matches []EntityMatchedStop `json:"matches"`
}

// JourneyPlannerQuery is used to hold the data for querying JourneyPlannerItinerary
// TODO change dateTime from string to time.Time?
type JourneyPlannerQuery struct {
	From, To, Date, Time string
	Modes                []string
}

// JourneyPlannerItineraryResult represents Tfl.Api.Presentation.Entities.JourneyPlanner.ItineraryResult
type JourneyPlannerItineraryResult struct {
	Journeys []JourneyPlannerJourney `json:"journeys"`
}

// JourneyPlannerJourney represents Tfl.Api.Presentation.Entities.JourneyPlanner.Journey
type JourneyPlannerJourney struct {
	StartDateTime   string      `json:"startDateTime"`
	ArrivalDateTime string      `json:"arrivalDateTime"`
	Duration        uint16      `json:"duration"`
	Legs            []Leg       `json:"legs"`
	Fare            JourneyFare `json:"fare"`
}

// Leg represents Tfl.Api.Presentation.Entities.JourneyPlanner.Leg
type Leg struct {
	Duration       uint16 `json:"duration"`
	Instruction    `json:"instruction"`
	DepartureTime  string               `json:"departureTime"`
	ArrivalTime    string               `json:"arrivalTime"`
	DeparturePoint StopPointAPIResponse `json:"departurePoint"`
	ArrivalPoint   StopPointAPIResponse `json:"arrivalPoint"`
}

// Instruction represents Tfl.Api.Presentation.Entities.Instruction
type Instruction struct {
	Summary  string `json:"summary"`
	Detailed string `json:"detailed"`
}

// JourneyFare represents Tfl.Api.Presentation.Entities.JourneyPlanner.JourneyFare
type JourneyFare struct {
	TotalCost uint16 `json:"totalCost"`
	Fares     []Fare `json:"fares"`
}

// Fare represents Tfl.Api.Presentation.Entities.JourneyPlanner.Fare
type Fare struct {
	LowZone           uint8     `json:"lowZone"`
	HighZone          uint8     `json:"highZone"`
	Cost              uint16    `json:"cost"`
	ChargeProfileName string    `json:"chargeProfileName"`
	IsHopperFare      bool      `json:"isHopperFare"`
	PeakCost          uint16    `json:"peak"`
	OffPeakCost       uint16    `json:"offPeak"`
	Taps              []FareTap `json:"taps"`
}

// FareTap represents Tfl.Api.Presentation.Entities.JourneyPlanner.FareTap
type FareTap struct {
	AtcoCode   string         `json:"atcoCode"`
	TapDetails FareTapDetails `json:"tapDetails"`
}

// FareTapDetails represents Tfl.Api.Presentation.Entities.JourneyPlanner.FareTapDetails
type FareTapDetails struct {
	ModeType     string `json:"modeType"`
	TapTimestamp string `json:"tapTimestamp"`
}

// FaresSection represents Tfl.Api.Presentation.Entities.Fares.FaresSection
type FaresSection struct {
	Header   string        `json:"header"`
	Index    int        `json:"index"`
	Journey  FaresJourney  `json:"journey"`
	Rows     []FareDetails `json:"rows"`
	Messages []Message     `json:"messages"`
}

// FaresJourney represents Tfl.Api.Presentation.Entities.Fares.Journey
type FaresJourney struct {
	FromStation FareStation `json:"fromStation"`
	ToStation   FareStation `json:"toStation"`
}

// FareStation represents Tfl.Api.Presentation.Entities.Fares.Journey
type FareStation struct {
	AtcoCode     string `json:"atcoCode"`
	CommonName   string `json:"commonName"`
	FareCategory string `json:"fareCategory"`
}

// Message represents Tfl.Api.Presentation.Entities.Message
type Message struct {
	BulletOrder int    `json:"bulletOrder"`
	Text        string `json:"messageText"`
}

// FareDetails represents Tfl.Api.Presentation.Fares.FareDetails
type FareDetails struct {
	StartDate               string   `json:"startDate"`
	EndDate                 string   `json:"endDate"`
	PassengerType           string   `json:"passengerType"`
	ContactlessPaygOnlyFare bool     `json:"contactlessPAYGOnlyFare"`
	From                    string   `json:"from"`
	To                      string   `json:"to"`
	FromStation             string   `json:"fromStation"`
	ToStation               string   `json:"toStation"`
	DisplayName             string   `json:"displayName"`
	DisplayOrder            int      `json:"displayOrder"`
	RouteDescription        string   `json:"routeDescription"`
	SpecialFare             bool     `json:"specialFare"`
	ThroughFare             bool     `json:"throughFare"`
	IsTour                  bool     `json:"isTour"`
	TicketsAvailable        []Ticket `json:"ticketsAvailable"`
}

// Ticket represents Tfl.Api.Presentation.Entities.Fares.Ticket
type Ticket struct {
	PassengerType string     `json:"passengerType"`
	TicketType    TicketType `json:"ticketType"`
	TicketTime    TicketTime `json:"ticketTime"`
	Cost          string     `json:"cost"`
	Description   string     `json:"description"`
	Mode          string     `json:"mode"`
	DisplayOrder  int        `json:"displayOrder"`
}

// TicketType represents Tfl.Api.Presentation.Entities.Fares.TicketType
type TicketType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// TicketTime represents Tfl.Api.Presentation.Entities.Fares.TicketTime
type TicketTime struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// SingleFareFinderInput is used as the input object for SingleFareFinder
type SingleFareFinderInput struct {
	From, To string
}
