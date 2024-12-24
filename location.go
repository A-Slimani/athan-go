package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// HELPER FUNCTIONS

func writeLocationInfoToJson(location Location, locationCacheJson string) error {
	city := location.City
	country := location.Country

	pattern := "^[a-zA-Z]+$"
	match, err := regexp.MatchString(pattern, city)
	if !match || err != nil {
		return fmt.Errorf("invalid city input")
	}
	match, err = regexp.MatchString(pattern, country)
	if !match || err != nil {
		return fmt.Errorf("invalid city input")
	}

	locationJson, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("error marshalling json: %v", err)
	}

	f, err := os.Create(locationCacheJson)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(locationJson)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// MAIN FUNCTION

func CacheLocation(reader *bufio.Reader, locationCacheJson string) error {
	fmt.Print("Set location automatically? (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch strings.ToLower(input) {
	case "y", "yes":
		resp, err := http.Get("https://ipinfo.io")
		if err != nil {
			return fmt.Errorf("request error: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}

		var location Location
		err = json.Unmarshal(body, &location)
		if err != nil {
			return fmt.Errorf("error unmarshalling json: %v", err)
		}
		return writeLocationInfoToJson(location, locationCacheJson)

	case "n", "no":
		fmt.Print("Please enter your location in this format (city, country_code) e.g. Sydney, AU: ")
		input, _ = reader.ReadString('\n')

		inputSplit := strings.Split(input, ",")

		if len(inputSplit) != 2 {
			return fmt.Errorf("invalid input")
		}

		city := strings.TrimSpace(inputSplit[0])
		country := strings.TrimSpace(inputSplit[1])

		location := Location{
			City:    city,
			Country: country,
		}

		return writeLocationInfoToJson(location, locationCacheJson)

	default:
		return fmt.Errorf("invalid input")
	}
}
