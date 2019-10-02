package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/guysports/oddsreader/pkg/types"
)

const (
	encodedbaseurl = "aHR0cHM6Ly9hcGkudGhlLW9kZHMtYXBpLmNvbS92My9vZGRzLz8="
	encodeduri     = "aHR0cHM6Ly9wcm9maXRtYXhpbWlzZXIuY28udWsvb2Rkcy1zZXJ2ZXIvZ2V0ZGF0YS5waHA/ZHJhdz0xJnN0YXJ0PTAmbGVuZ3RoPTIwMCZzZWFyY2glNUJ2YWx1ZSU1RD0mc2VhcmNoJTVCcmVnZXglNUQ9ZmFsc2Umc2VhcmNoVGV4dD1FbmdsYW5kJTNBK0NhcmFiYW8rQ3VwJTJDRW5nbGFuZCslMkYrUHJlbWllcitMZWFndWUrJTNFK1JlZ3VsYXIrU2Vhc29uLTE5JTJDRW5nbGFuZCslMkYrQ2hhbXBpb25zaGlwKyUzRStSZWd1bGFyK1NlYXNvbi0xOSUyQ1NwYWluKyUyRitQcmltZXJhK0RpdmlzaW4rJTNFK1JlZ3VsYXIrU2Vhc29uLTE5JTJDSXRhbHkrJTJGK1NlcmllK0ErJTNFK1JlZ3VsYXIrU2Vhc29uLTE5JTJDR2VybWFueSslMkYrMS4rQnVuZGVzbGlnYSslM0UrUmVndWxhcitTZWFzb24tMTklMkNGQkVDVVArJTJGK1VFRkErQ2hhbXBpb25zK0xlYWd1ZSUyQ0ZCRUNVUCslMkYrVUVGQStFdXJvcGErTGVhZ3VlJTJGK0ZBK0N1cCZzcG9ydExpc3Q9MiZleGNoYW5nZUxpc3Q9MSUyQzMmbWFya2V0TGlzdD0yJm1pbkxpcXVpZGl0eT0xJm1pblJhdGluZz04MCZtYXhSYXRpbmc9MjAwJm1pbk9kZHM9MSZtYXhPZGRzPTUwJmJvb2tpZU5hbWU9QmV0KzM2NSZwZXJpb2RGcm9tPTEwJnJldGVudGlvbj1mYWxzZSZfPTE1Njk3NTMzMTc2NCZwZXJpb2RUbz0="
	encodedh1      = "aHR0cHM6Ly9wcm9maXRtYXhpbWlzZXIuY28udWsvbWFzdGVybWluZC9wcm9maXRtYXhpbWlzZXIvdXNlci9vZGRzX3NlcnZlcj8="
	encodedh2      = "cHJvZml0bWF4aW1pc2VyLmNvLnVr"
	encodedc1      = "bWFzdGVybWluZA=="
	encodedc2      = "YmFzZV9zaXRl"
	encodedv1      = "MQ=="
	encodedv2      = "bWFzdGVybWluZA=="
)

type (
	// OddsAPI - API interface to REST calls to retrieve league odds
	OddsAPI struct {
		apikey string
		uri    string
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
	json.Unmarshal(byteValue, &f)

	return &f, nil
}

// RetrievePM - Get the fixtures from PM
func (o *OddsAPI) RetrievePM(interval int) (*types.PMFeed, error) {

	// Build the request
	url := uri(encodeduri) + strconv.Itoa(interval)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Referer", uri(encodedh1))
	req.Header.Set("Host", uri(encodedh2))
	req.AddCookie(&http.Cookie{
		Name:  uri(encodedc1),
		Value: uri(encodedv1),
	})
	req.AddCookie(&http.Cookie{
		Name:  uri(encodedc2),
		Value: uri(encodedv2),
	})

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
	json.Unmarshal(byteValue, &f)

	return &f, nil
}

func uri(enc string) string {
	data, _ := base64.StdEncoding.DecodeString(enc)
	return string(data)
}
