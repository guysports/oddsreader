package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/guysports/oddsreader/pkg/api"
	"github.com/guysports/oddsreader/pkg/types"
)

type (
	// PM has the options to tailor the comparisonn
	PM struct {
		Queryfile string   `help:"Specify location of file containing PM queries"`
		Bookmaker string   `help:"Specify the bookmaker to compare the odds for"`
		Oddslimit float64  `help:"Specify the virtual odds of win / ql as a filter"`
		QLLimit   float64  `help:"Specify the maximum qualifying loss wanted as a filter"`
		Exchanges []string `help:"Specify a comma separated list of exchanges to compare layodds with"`
		BackStake float64  `help:"Specify the amount to bet off the back bet on the fixture"`
		Interval  int      `help:"Specify the number of hours into the future to look for fixtures"`
	}
)

// Run the compare command to obtain the odds and filter based on options
func (c *PM) Run(globals *Globals) error {
	apic := api.NewOddsAPI(globals.APIKey)
	queryDetails, err := parseQuery(c.Queryfile)
	if err != nil {
		return err
	}

	rawreport := []types.OddsReport{}
	// fixtures, err := apic.RetrievePMFromFile("pm.json")
	fixtures, err := apic.RetrievePM(queryDetails, c.Interval)
	if err != nil {
		fmt.Printf("Error returned %v\n", err)
	}

	for _, event := range fixtures.Events {
		if event.Team == "Draw" {
			continue
		}
		if event.Exchange != "Betfair" {
			continue
		}
		event.BackOdds, _ = strconv.ParseFloat(event.BackString, 64)
		event.LayOdds, _ = strconv.ParseFloat(event.LayString, 64)
		fixrep := c.comparePrices(event)
		if fixrep != nil {
			rawreport = append(rawreport, *fixrep)
		}
	}

	// Sort the report into fixtures
	report := []types.OddsReport{}
	for idx1, rep1 := range rawreport {
		if rep1.Added {
			continue
		}
		for idx2, rep2 := range rawreport {
			// Ignore same report or report already added
			if idx1 >= idx2 {
				continue
			}

			// Add reports of the same fixture to rep1
			if rep2.Fixture == rep1.Fixture {

				rep1.ReportInformation = append(rep1.ReportInformation, rep2.ReportInformation...)
				rawreport[idx2].Added = true
			}
		}
		report = append(report, rep1)
	}

	printReport(fmt.Sprintf("PM Selections matching Odds Criteria of %0.0f", c.Oddslimit), report)

	return nil
}

func (c *PM) comparePrices(fixture types.PMEvent) *types.OddsReport {

	fixReport := types.OddsReport{
		ReportInformation: []types.OddsInformation{},
		Fixture:           fmt.Sprintf("%s v %s", fixture.Home, fixture.Away),
		Date:              fmt.Sprintf("%s %s", fixture.Date, fixture.Time),
	}

	var oddsInfo *types.OddsInformation

	// Get commission for the exchange
	comm := getCommission(fixture.Exchange)

	// Check for event
	oddsInfo = c.validateCriteria(fixture.BackOdds, fixture.LayOdds, comm)
	if oddsInfo != nil {
		oddsInfo.Exchange = fixture.Exchange
		oddsInfo.Team = fixture.Team
		fixReport.ReportInformation = append(fixReport.ReportInformation, *oddsInfo)
	}

	// If there is a report then return it
	if len(fixReport.ReportInformation) > 0 {
		return &fixReport
	}

	return nil
}

func (c *PM) validateCriteria(backodds, layodds, comm float64) *types.OddsInformation {

	laystake, ql := c.calculateQualifyingLoss(backodds, layodds, comm)
	profit, odds := c.calculateOutcomeOdds(backodds, laystake, ql, comm)

	if odds > c.Oddslimit && (ql <= c.QLLimit || c.QLLimit == 0) {
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

func (c *PM) calculateQualifyingLoss(backodds, layodds, commission float64) (float64, float64) {
	// Lay stake = stake * backodds / (layodds-commission)
	laystake := (c.BackStake * backodds) / (layodds - commission)

	// Qualifying loss = stake - laystake * (1-commission)
	ql := c.BackStake - laystake*(1-commission)

	return math.Round(laystake*100) / 100, math.Round(ql*100) / 100
}

func (c *PM) calculateOutcomeOdds(backodds, laystake, ql, commission float64) (float64, float64) {
	// Winning in event of 2up
	profit := c.BackStake*(backodds-1) + laystake*(1-commission)
	odds := 1 + (profit / ql)
	if ql <= 0 {
		odds = 0
	}
	return math.Round(profit*100) / 100, math.Round(odds*100) / 100
}

func parseQuery(queryFile string) (*types.Query, error) {
	queryBytes, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return nil, err
	}
	query := types.Query{}
	_ = json.Unmarshal(queryBytes, &query)
	fmt.Printf("Parsed query: %#v\n", query)
	return &query, nil
}
