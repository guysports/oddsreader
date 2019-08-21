package types

type (
	// Odds - holds the home/draw/away odds and lay odds if exchange
	Odds struct {
		H2h    []int `json:"h2h"`
		H2hlay []int `json:"h2h_lay"`
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
		LeagueName string      `json:"sport_name"`
		Teams      []string    `json:"teams"`
		Kickoff    int64       `json:"commence_time"`
		Bookies    []Bookmaker `json:"sites"`
	}
)
