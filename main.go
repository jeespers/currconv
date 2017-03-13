package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

// Payload is a struct for the response data
type Payload struct {
	Amount    float64
	Currency  string
	Converted map[string]float64
}

// Upload is a struct for the fetched data
type Upload struct {
	Currency  string             `json:"base"`
	Converted map[string]float64 `json:"rates"`
}

// P is the global Payload
var P Payload

// U is the global Upload var
var U Upload

func main() {
	http.HandleFunc("/convert", convertionToJSONHandler)
	http.ListenAndServe(":8080", nil)
}

// Parses the query from the url
func requestInputParser(w http.ResponseWriter, r *http.Request) (float64, string, error) {
	a, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	c := r.FormValue("currency")
	return a, c, err
}

// fetches and unmarshals the json-data from Fixer.io
func jsonFetcher(url string) error {
	// finds it
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	// reads it
	read, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Checks it
	if len(read) < 30 {
		return err
	}

	// Writes it to U Upload according to `json:"<attr>"`
	err = json.Unmarshal(read, &U)
	if err != nil {
		return err
	}
	return nil
}

func convertionToJSONHandler(w http.ResponseWriter, r *http.Request) {
	// parse header request
	a, c, err := requestInputParser(w, r)
	if err != nil {
		http.Error(w, "ERROR 400 BAD SERVER ERROR", 500)
	}

	url := "http://api.fixer.io/latest?base=" + c

	err = jsonFetcher(url)
	if err != nil {
		http.Error(w, "ERROR 500 BAD SERVER ERROR", 500)
	}

	// Calculate conversions
	for key, value := range U.Converted {
		U.Converted[key] = RoundPlus(value*a, 2)
	}

	// Fill Payload with data
	P := Payload{a, U.Currency, U.Converted}
	result, err := json.MarshalIndent(P, "", "  ")
	if err != nil {
		http.Error(w, "ERROR 500 INTERNAL SERVER ERROR", 500)
	}
	// update response header and print data to body
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(result))
}

// Round
func Round(f float64) float64 {
	return math.Floor(f + .5)
}

// RoundPlus implements Round to handle trunkation
func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}
