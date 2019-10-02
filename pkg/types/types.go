package types

type (
	// Odds - holds the home/draw/away odds and lay odds if exchange
	Odds struct {
		H2h    []float64 `json:"h2h"`
		H2hlay []float64 `json:"h2h_lay"`
	}

	// Bookmaker - holds the information about the bookmaker and the odds being offered for a fixture
	Bookmaker struct {
		ID       string `json:"site_key"`
		Name     string `json:"site_name"`
		LastTime int64  `json:"last_update"`
		Offers   Odds   `json:"odds"`
	}

	// Fixture - holds the information about the teams playing and the bookmaker offers
	Fixture struct {
		LeagueID   string      `json:"sport_key"`
		LeagueName string      `json:"sport_nice"`
		Teams      []string    `json:"teams"`
		Kickoff    int64       `json:"commence_time"`
		Bookies    []Bookmaker `json:"sites"`
	}

	// Fixtures - holds the informationn about a leagues fixtures
	Fixtures struct {
		RequestOutcome bool      `json:"success"`
		League         []Fixture `json:"data"`
	}

	// OddsInformation - holds the information for different exchanges about a team
	OddsInformation struct {
		ComparativeOdds float64
		QualifyingLoss  float64
		BackStake       float64
		BackOdds        float64
		LayStake        float64
		LayOdds         float64
		Profit          float64
		Exchange        string
		Team            string
	}

	// OddsReport - holds the information about fixtures that have met the filter criteria
	OddsReport struct {
		ReportInformation []OddsInformation
		Fixture           string
		Date              string
		Added             bool
	}

	// PMFeed - holds the json return from the http request
	PMFeed struct {
		Events      []PMEvent `json:"data"`
		RecordCount int       `json:"data_lenght"`
	}

	// PMEvent - holds the information about a match from the PM feed
	PMEvent struct {
		Exchange   string `json:"exchange"`
		Fixture    string `json:"event_name"`
		Home       string `json:"team1"`
		Away       string `json:"team2"`
		LeagueName string `json:"competition"`
		Date       string `json:"date"`
		Time       string `json:"time"`
		Team       string `json:"runner_name"`
		Bookmaker  string `json:"bookmaker_name"`
		BackString string `json:"bookmaker_odds"`
		LayString  string `json:"betfair_odds"`
		BackOdds   float64
		LayOdds    float64
	}
)
