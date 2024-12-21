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

func convertTime(athanTime string) string {
	newTime, err := time.Parse("15:04", athanTime)
	if err != nil {
		fmt.Println("Error parsing time: ", err)
	}
	return newTime.Format("15:04")
}

func getTodaysAthanTimes(athanCacheJson string) AthanTimes {
	athanJson, err := os.ReadFile(athanCacheJson)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	var times []AthanTimes
	err = json.Unmarshal(athanJson, &times)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
	}

	_, _, day := time.Now().Date()
	return times[day-1]
}

func cacheAthan(locationCacheJson string, athanCacheJson string) {

	// getting relevant data from location.json
	locationJson, err := os.ReadFile(locationCacheJson)
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(locationJson, &data)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
	}

	location := strings.Split(data["loc"].(string), ",")
	latitude := location[0]
	longitude := location[1]

	// getting athan times from api
	const baseUrl = "http://api.aladhan.com/v1/calendar"
	params := url.Values{
		"latitude":  {latitude},
		"longitude": {longitude},
		"method":    {"3"},
		"month":     {strconv.Itoa(int(time.Now().Month()))},
		"year":      {strconv.Itoa(time.Now().Year())},
	}

	resp, err := http.Get(baseUrl + "?" + params.Encode())
	if err != nil {
		fmt.Println("Request error: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
		return
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
		fmt.Println("Error marshalling athan times: ", err)
		return
	}

	err = os.WriteFile(athanCacheJson, athanTimesJson, 0644)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}
}

func nextAthanString(athanCacheJson string) {
	currentHours, currentMinutes := time.Now().Hour(), time.Now().Minute()
	currentCombined := currentHours*60 + currentMinutes
	todaysTimes := getTodaysAthanTimes(athanCacheJson)

	val := reflect.ValueOf(todaysTimes)
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name != "Sunrise" {
			athanTime, err := time.Parse("15:04", val.Field(i).Interface().(string))
			if err != nil {
				fmt.Println("Error parsing time: ", err)
			}
			hours, minutes := athanTime.Hour(), athanTime.Minute()
			athanCombined := hours*60 + minutes
			if athanCombined > currentCombined {
				timeRemaining := time.Duration(athanCombined-currentCombined) * time.Minute
				hours := int(timeRemaining.Hours())
				minutes := int(timeRemaining.Minutes()) % 60
				athanName := val.Type().Field(i).Name

				if hours == 0 && minutes == 0 {
					fmt.Printf("%s is now \n", athanName)
				}

				parts := []string{}
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
				result := fmt.Sprintf("%s is in %s", athanName, strings.Join(parts, " and "))
				fmt.Println(result)
				break
			}
		}
	}
}

func allAthanTimes(athanCacheJson string) {
	values := reflect.ValueOf(getTodaysAthanTimes(athanCacheJson))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Prayer", "Time"})

	for i := 0; i < values.NumField(); i++ {
		table.Append([]string{values.Type().Field(i).Name, values.Field(i).Interface().(string)})
	}

	table.Render()
}