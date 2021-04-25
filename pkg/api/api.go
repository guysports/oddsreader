package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/guysports/oddsreader/pkg/types"
)

const (
	encodedbaseurl = "aHR0cHM6Ly9hcGkudGhlLW9kZHMtYXBpLmNvbS92My9vZGRzLz8="
)

type (
	// OddsAPI - API interface to REST calls to retrieve league odds
	OddsAPI struct {
		apikey string
		client *http.Client
	}
)

// NewOddsAPI - Create the API interface
func NewOddsAPI(apikey string) *OddsAPI {

	return &OddsAPI{
		apikey: apikey,
		client: &http.Client{},
	}
}

// RetrieveLeague - Get the fixtures for a league
func (o *OddsAPI) RetrieveLeague(region, league, market string) (*types.Fixtures, error) {

	// Build the request
	url := fmt.Sprintf("%sapiKey=%s&sport=%s&region=%s&market=%s", uri(encodedbaseurl), o.apikey, league, region, market)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error constructing HTTP request: %s", err.Error())
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Unmarshall responnse into struct
	defer resp.Body.Close()

	// Make sure we get the correct response.
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response from API got: %d\n", resp.StatusCode)
	}
	//remainingCalls := resp.Header.Get()

	// Define a new fixtures object
	f := types.Fixtures{}

	// Decode body in interface{} here, because the Body.Close() means decode doesn't work after stream is closed.
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		fmt.Printf("Invalid JSON in HTTP Response from API %s\n", err.Error())
	}

	return &f, nil
}

// RetrieveLeagueFromFile - umarshall into structs from file
func (o *OddsAPI) RetrieveLeagueFromFile(file string) (*types.Fixtures, error) {

	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	f := types.Fixtures{}
	_ = json.Unmarshal(byteValue, &f)

	return &f, nil
}

// RetrievePM - Get the fixtures from PM
func (o *OddsAPI) RetrievePM(queryDetail *types.Query, interval int) (*types.PMFeed, error) {

	// Build the request
	timeMillis := time.Now().UnixNano() / int64(time.Millisecond)

	url := queryDetail.BaseUri + strconv.Itoa(interval) + "&_=" + strconv.FormatInt(timeMillis, 10)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Referer", queryDetail.Referer)
	req.Header.Set("Host", queryDetail.Host)
	req.AddCookie(&queryDetail.Cookies[0])
	req.AddCookie(&queryDetail.Cookies[1])

	if err != nil {
		return nil, fmt.Errorf("Error constructing HTTP request: %s", err.Error())
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Unmarshall responnse into struct
	defer resp.Body.Close()

	// Make sure we get the correct response.
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response from API got: %d\n", resp.StatusCode)
	}

	// Define a new fixtures object
	f := types.PMFeed{}

	// Decode body in interface{} here, because the Body.Close() means decode doesn't work after stream is closed.
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		fmt.Printf("Invalid JSON in HTTP Response from API %s\n", err.Error())
	}

	return &f, nil
}

// RetrievePMFromFile - unmarshall into structs from file
func (o *OddsAPI) RetrievePMFromFile(file string) (*types.PMFeed, error) {

	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	f := types.PMFeed{}
	_ = json.Unmarshal(byteValue, &f)

	return &f, nil
}

func uri(enc string) string {
	data, _ := base64.StdEncoding.DecodeString(enc)
	return string(data)
}
