package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type AthanTimes struct {
	Fajr    string `json:"Fajr"`
	Sunrise string `json:"Sunrise"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type Timings struct {
	Timing AthanTimes `json:"timings"`
}

type Response struct {
	Data []Timings `json:"data"`
}

// HELPER FUNCTIONS

func convertTime(athanTime string) string {
	athanTime = athanTime[:5]
	newTime, err := time.Parse("15:04", athanTime)
	if err != nil {
		fmt.Println("Error parsing time: ", err)
	}
	return newTime.Format("15:04")
}

func buildAthanString(hours int, minutes int, athanName string, athanTime string) string {
	parts := []string{}
	if hours == 0 && minutes == 0 {
		return fmt.Sprintf("%s is now", athanName) + " at " + athanTime
	}
	if hours > 0 {
		part := fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			part += "s"
		}
		parts = append(parts, part)
	}
	if minutes > 0 {
		part := fmt.Sprintf("%d minute", minutes)
		if minutes > 1 {
			part += "s"
		}
		parts = append(parts, part)
	}
	return fmt.Sprintf("%s in %s", athanName, strings.Join(parts, " and ")) + " at " + athanTime
}

// MAIN FUNCTIONS

func getAthanTimesForDay(athanCacheJson string, day int) (*AthanTimes, error) {
	d := day - 1
	if d < 0 || d > 30 {
		return nil, fmt.Errorf("invalid day outside of month range")
	}

	athanJson, err := os.ReadFile(athanCacheJson)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var times []AthanTimes
	err = json.Unmarshal(athanJson, &times)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	currentTime := &times[d]
	return currentTime, nil
}

func CacheAthanTimes(locationCacheJson string, athanCacheJson string) error {
	locationJson, err := os.ReadFile(locationCacheJson)
	if err != nil {
		return fmt.Errorf("error reading file: %d", err)
	}

	var location Location
	err = json.Unmarshal(locationJson, &location)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %d", err)
	}

	const baseUrl = "https://api.aladhan.com/v1/calendarByCity"
	params := url.Values{
		"city":    {location.City},
		"country": {location.Country},
		"method":  {"3"},
		"month":   {strconv.Itoa(int(time.Now().Month()))},
		"year":    {strconv.Itoa(time.Now().Year())},
	}

	resp, err := http.Get(baseUrl + "?" + params.Encode())
	if err != nil {
		return fmt.Errorf("request error: %d", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %d", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %d", err)
	}

	var transformedData []AthanTimes
	// find a better way to loop through all key value pairs
	for i := range response.Data {
		response.Data[i].Timing.Fajr = convertTime(response.Data[i].Timing.Fajr)
		response.Data[i].Timing.Sunrise = convertTime(response.Data[i].Timing.Sunrise)
		response.Data[i].Timing.Dhuhr = convertTime(response.Data[i].Timing.Dhuhr)
		response.Data[i].Timing.Asr = convertTime(response.Data[i].Timing.Asr)
		response.Data[i].Timing.Maghrib = convertTime(response.Data[i].Timing.Maghrib)
		response.Data[i].Timing.Isha = convertTime(response.Data[i].Timing.Isha)
		transformedData = append(transformedData, response.Data[i].Timing)
	}

	athanTimesJson, err := json.MarshalIndent(transformedData, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling athan times: %d", err)
	}

	err = os.WriteFile(athanCacheJson, athanTimesJson, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %d", err)
	}
	return nil
}

func GetNextAthan(athanCacheJson string, currentTime time.Time) (*string, error) {
	currentHours, currentMinutes := currentTime.Hour(), currentTime.Minute()
	currentTimeCombined := currentHours*60 + currentMinutes
	todaysTimes, err := getAthanTimesForDay(athanCacheJson, time.Now().Day()-1)
	if err != nil || todaysTimes == nil {
		return nil, fmt.Errorf("error getting athan times: %w", err)
	}

	val := reflect.ValueOf(*todaysTimes)
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name != "Sunrise" {
			timeToCheck, err := time.Parse("15:04", val.Field(i).Interface().(string))
			if err != nil {
				return nil, fmt.Errorf("error parsing time: %w", err)
			}
			hours, minutes := timeToCheck.Hour(), timeToCheck.Minute()
			athanCombined := hours*60 + minutes
			if athanCombined > currentTimeCombined {
				timeRemaining := time.Duration(athanCombined-currentTimeCombined) * time.Minute
				hours := int(timeRemaining.Hours())
				minutes := int(timeRemaining.Minutes()) % 60
				athanName := val.Type().Field(i).Name
				athanTime := val.Field(i).Interface().(string)

				returnStr := buildAthanString(hours, minutes, athanName, athanTime)
				return &returnStr, nil
			}
		}
	}
	tomorrowsTimes, _ := getAthanTimesForDay(athanCacheJson, time.Now().Day())
	timeToCheck, err := time.Parse("15:04", tomorrowsTimes.Fajr)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %w", err)
	}

	hours, minutes := timeToCheck.Hour(), timeToCheck.Minute()
	athanTimeCombined := hours*60 + minutes
	maxTime := 24 * 60
	timeRemaining := time.Duration((maxTime-currentTimeCombined)+athanTimeCombined) * time.Minute
	hours = int(timeRemaining.Hours())
	minutes = int(timeRemaining.Minutes()) % 60

	athanName := "Fajr"
	athanTime := val.Field(0).Interface().(string)

	returnStr := buildAthanString(hours, minutes, athanName, athanTime)
	return &returnStr, nil
}

func AllAthanTimes(athanCacheJson string, locationCacheJson string, today int) error {
	d := today - 1
	if d < 0 || d > 30 {
		return fmt.Errorf("invalid day outside of month range")
	}

	times, _ := getAthanTimesForDay(athanCacheJson, today)
	values := reflect.ValueOf(*times)

	table := tablewriter.NewWriter(os.Stdout)
	location, err := os.ReadFile(locationCacheJson)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var locationData Location
	err = json.Unmarshal(location, &locationData)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	fmt.Printf("Athan times for %s, %s\n", locationData.City, locationData.Country)
	table.SetHeader([]string{"Prayer", "Time"})

	for i := 0; i < values.NumField(); i++ {
		table.Append([]string{values.Type().Field(i).Name, values.Field(i).Interface().(string)})
	}

	table.Render()

	return nil
}
