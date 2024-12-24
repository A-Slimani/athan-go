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
	newTime, err := time.Parse("15:04", athanTime)
	if err != nil {
		fmt.Println("Error parsing time: ", err)
	}
	return newTime.Format("15:04")
}

func buildAthanString(hours int, minutes int, athanName string) string {
	parts := []string{}
	if hours == 0 && minutes == 0 {
		return fmt.Sprintf("%s is now \n", athanName)
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
	return fmt.Sprintf("%s in %s", athanName, strings.Join(parts, " and "))
}

// MAIN FUNCTIONS
func getAthanTimesForDay(athanCacheJson string, day int) AthanTimes {
	athanJson, err := os.ReadFile(athanCacheJson)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	var times []AthanTimes
	err = json.Unmarshal(athanJson, &times)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
	}

	return times[day-1]
}

func CacheAthanTimes(locationCacheJson string, athanCacheJson string) error {

	// getting relevant data from location.json
	locationJson, err := os.ReadFile(locationCacheJson)
	if err != nil {
		return fmt.Errorf("error reading file: %d", err)
	}

	var location Location
	err = json.Unmarshal(locationJson, &location)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %d", err)
	}

	// getting athan times from api
	const baseUrl = "https://api.aladhan.com/v1/calendar"
	// change this to always use the city / country combo
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
	for i := range response.Data {
		response.Data[i].Timing.Fajr = convertTime(response.Data[i].Timing.Fajr[:5])
		response.Data[i].Timing.Sunrise = convertTime(response.Data[i].Timing.Sunrise[:5])
		response.Data[i].Timing.Dhuhr = convertTime(response.Data[i].Timing.Dhuhr[:5])
		response.Data[i].Timing.Asr = convertTime(response.Data[i].Timing.Asr[:5])
		response.Data[i].Timing.Maghrib = convertTime(response.Data[i].Timing.Maghrib[:5])
		response.Data[i].Timing.Isha = convertTime(response.Data[i].Timing.Isha[:5])
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

func GetNextAthan(athanCacheJson string) {
	currentHours, currentMinutes := time.Now().Hour(), time.Now().Minute()
	currentTimeCombined := currentHours*60 + currentMinutes
	todaysTimes := getAthanTimesForDay(athanCacheJson, time.Now().Day()-1)

	val := reflect.ValueOf(todaysTimes)
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name != "Sunrise" {
			athanTime, err := time.Parse("15:04", val.Field(i).Interface().(string))
			if err != nil {
				fmt.Println("Error parsing time: ", err)
			}
			hours, minutes := athanTime.Hour(), athanTime.Minute()
			athanCombined := hours*60 + minutes
			if athanCombined > currentTimeCombined {
				timeRemaining := time.Duration(athanCombined-currentTimeCombined) * time.Minute
				hours := int(timeRemaining.Hours())
				minutes := int(timeRemaining.Minutes()) % 60
				athanName := val.Type().Field(i).Name

				returnStr := buildAthanString(hours, minutes, athanName)
				fmt.Println(returnStr)
				break
			}
		}
	}
	tomorrowsTimes := getAthanTimesForDay(athanCacheJson, time.Now().Day())
	athanTime, err := time.Parse("15:04", tomorrowsTimes.Fajr)
	if err != nil {
		fmt.Println("Error parsing time: ", err)
	}

	hours, minutes := athanTime.Hour(), athanTime.Minute()
	athanTimeCombined := hours*60 + minutes
	maxTime := 24 * 60
	timeRemaining := time.Duration((maxTime-currentTimeCombined)+athanTimeCombined) * time.Minute
	hours = int(timeRemaining.Hours())
	minutes = int(timeRemaining.Minutes()) % 60

	athanName := "Fajr"

	returnStr := buildAthanString(hours, minutes, athanName)
	fmt.Println(returnStr)
}

func AllAthanTimes(athanCacheJson string) {
	values := reflect.ValueOf(getAthanTimesForDay(athanCacheJson, time.Now().Day()-1))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Prayer", "Time"})

	for i := 0; i < values.NumField(); i++ {
		table.Append([]string{values.Type().Field(i).Name, values.Field(i).Interface().(string)})
	}

	table.Render()
}
