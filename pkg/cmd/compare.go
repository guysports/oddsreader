package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/guysports/oddsreader/pkg/api"
	"github.com/guysports/oddsreader/pkg/types"
)

type (
	// Compare has the options to tailor the comparisonn
	Compare struct {
		Bookmaker string   `help:"Specify the bookmaker to compare the odds for"`
		Oddslimit float64  `help:"Specify the virtual odds of win / ql as a filter"`
		Exchanges []string `help:"Specify a comma separated list of exchanges to compare layodds with"`
		BackStake float64  `help:"Specify the amount to bet off the back bet on the fixture"`
	}
)

const (
	sevenDaysAsSeconds = 604800
)

// Run the compare command to obtain the odds and filter based on options
func (c *Compare) Run(globals *Globals) error {
	apic := api.NewOddsAPI(globals.APIKey)

	for _, league := range globals.Leagues {
		report := []types.OddsReport{}
		fixtures, err := apic.RetrieveLeague("uk", league, "h2h")
		//file := fmt.Sprintf("%s.json", league)
		//fixtures, err := apic.RetrieveLeagueFromFile(file)
		if err != nil {
			fmt.Printf("Error returned %v\n", err)
		}
		for _, fixture := range fixtures.League {
			fixrep := c.comparePrices(c.Bookmaker, c.Exchanges, fixture)
			if fixrep != nil {
				report = append(report, *fixrep)
			}
		}

		printReport(fixtures.League[0].LeagueName, report)

	}
	return nil
}

func (c *Compare) comparePrices(bookmaker string, exchanges []string, fixture types.Fixture) *types.OddsReport {

	// Only interested in fixtures within 7 days
	fixtureTime := time.Unix(0, fixture.Kickoff*1000*int64(time.Millisecond))
	if time.Until(fixtureTime) > 168*time.Hour {
		return nil
	}

	fixReport := types.OddsReport{
		ReportInformation: []types.OddsInformation{},
		Fixture:           fmt.Sprintf("%s v %s", fixture.Teams[0], fixture.Teams[1]),
		Date:              time.Unix(0, fixture.Kickoff*1000*int64(time.Millisecond)).Format("2 Jan 2006 15:04 BST"),
	}

	// Get bookmaker, betfair, smarkets and matchbook site objects
	fmt.Printf("Exchanges %v\n", exchanges)
	for _, exchange := range exchanges {

		var oddsInfo *types.OddsInformation
		bm := types.Bookmaker{}
		ex := types.Bookmaker{}

		for _, obj := range fixture.Bookies {
			// Get the required bookmaker and exchange objects
			if obj.ID == bookmaker {
				bm = obj
				continue
			}
			if obj.ID == exchange {
				ex = obj
				continue
			}
		}

		// Check we have both bookmaker and exchange odds available for comparison
		if bm.ID == "" || ex.ID == "" {
			continue
		}

		// Get commission for the exchange
		comm := getCommission(ex.ID)

		// Check for home team
		oddsInfo = c.validateCriteria(bm.Offers.H2h[0], ex.Offers.H2hlay[0], comm)
		if oddsInfo != nil {
			oddsInfo.Exchange = ex.ID
			oddsInfo.Team = fixture.Teams[0]
			fixReport.ReportInformation = append(fixReport.ReportInformation, *oddsInfo)

		}

		// Check for away team
		oddsInfo = c.validateCriteria(bm.Offers.H2h[1], ex.Offers.H2hlay[1], comm)
		if oddsInfo != nil {
			oddsInfo.Exchange = ex.ID
			oddsInfo.Team = fixture.Teams[1]
			fixReport.ReportInformation = append(fixReport.ReportInformation, *oddsInfo)
		}
	}

	// If there is a report then return it
	if len(fixReport.ReportInformation) > 0 {
		return &fixReport
	}

	return nil
}

func (c *Compare) validateCriteria(backodds, layodds, comm float64) *types.OddsInformation {

	laystake, ql := c.calculateQualifyingLoss(backodds, layodds, comm)
	profit, odds := c.calculateOutcomeOdds(backodds, laystake, ql, comm)

	if odds > c.Oddslimit {
		return &types.OddsInformation{
			ComparativeOdds: odds,
			QualifyingLoss:  ql,
			BackStake:       c.BackStake,
			BackOdds:        backodds,
			LayStake:        laystake,
			LayOdds:         layodds,
			Profit:          profit,
		}
	}

	return nil
}

func printReport(league string, report []types.OddsReport) {

	fmt.Printf("-------------------- %s --------------------\n", league)

	for _, fixture := range report {
		fmt.Printf("%s %s Qualifying Bets on Exchanges\n", fixture.Date, fixture.Fixture)
		for _, line := range fixture.ReportInformation {
			fmt.Printf("Exchange %s -- %s Backstake £%0.2f @ %0.2f, Laystake £%0.2f @ %0.2f, QL £%0.2f, Potential Profit £%0.2f at Odds %0.2f\n",
				line.Exchange,
				line.Team,
				line.BackStake,
				line.BackOdds,
				line.LayStake,
				line.LayOdds,
				line.QualifyingLoss,
				line.Profit,
				line.ComparativeOdds,
			)
		}

		fmt.Println("")
	}
}

func (c *Compare) calculateQualifyingLoss(backodds, layodds, commission float64) (float64, float64) {
	// Lay stake = stake * backodds / (layodds-commission)
	laystake := (c.BackStake * backodds) / (layodds - commission)

	// Qualifying loss = stake - laystake * (1-commission)
	ql := c.BackStake - laystake*(1-commission)

	return math.Round(laystake*100) / 100, math.Round(ql*100) / 100
}

func (c *Compare) calculateOutcomeOdds(backodds, laystake, ql, commission float64) (float64, float64) {
	// Winning in event of 2up
	profit := c.BackStake*(backodds-1) + laystake*(1-commission)
	odds := 1 + (profit / ql)
	if ql <= 0 {
		odds = 0
	}
	return math.Round(profit*100) / 100, math.Round(odds*100) / 100
}

func getCommission(exchange string) float64 {
	switch strings.ToLower(exchange) {
	case "betfair":
		return 0.02
	case "matchbook":
		return 0.02
	case "smarkets":
		return 0.0
	}

	// Default to 5%
	return 0.05
}
