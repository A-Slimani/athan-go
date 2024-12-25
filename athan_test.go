package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func Test_buildAthanString(t *testing.T) {
	testCases := []struct {
		hours     int
		minutes   int
		athanName string
		expected  string
	}{
		{3, 12, "Fajr", "Fajr in 3 hours and 12 minutes"},
		{3, 1, "Fajr", "Fajr in 3 hours and 1 minute"},
		{1, 2, "Fajr", "Fajr in 1 hour and 2 minutes"},
		{1, 1, "Fajr", "Fajr in 1 hour and 1 minute"},
		{1, 0, "Fajr", "Fajr in 1 hour"},
		{0, 1, "Fajr", "Fajr in 1 minute"},
		{0, 0, "Fajr", "Fajr is now \n"},
	}

	for _, tc := range testCases {
		result := buildAthanString(tc.hours, tc.minutes, tc.athanName)
		if result != tc.expected {
			t.Errorf("Expected: %s, returned: %s ", tc.expected, result)
		}
	}
}

func Test_CacheAthanTimes(t *testing.T) {
	locationCacheJson := os.TempDir() + "/location_test.json"
	athanCacheJson := os.TempDir() + "/athan_test.json"

	// setup to build location.json to run CacheAthanTimes
	location := Location{
		City:    "Sydney",
		Country: "AU",
	}
	locationJson, _ := json.Marshal(location)
	err := os.WriteFile(locationCacheJson, locationJson, 0644)
	if err != nil {
		t.Fatalf("Error writing location.json: %v", err)
	}

	err = CacheAthanTimes(locationCacheJson, athanCacheJson)
	if err != nil {
		t.Fatalf("Error caching athan times: %v", err)
	}

	if _, err := os.Stat(athanCacheJson); err != nil {
		t.Fatalf("Error checking file: %v", err)
	}

	os.Remove(locationCacheJson)
	os.Remove(athanCacheJson)
}

// TODO: testing the amount of times(object) in the jsonTODO: Do this test later
func Test_getAthanTimesForDay(t *testing.T) {
	athanCacheJson := "./testing_files/athan_test.json" // TODO: craft an example file to test with
	tests := []struct {
		name       string
		day        int
		wantErr    bool
		athanTimes AthanTimes
	}{
		{
			"Valid day lowest bound",
			1,
			false,
			AthanTimes{"03:56", "05:37", "12:44", "16:29", "19:51", "21:26"},
		},
		{
			"Valid day uppest bound",
			31,
			false,
			AthanTimes{"04:03", "05:47", "12:58", "16:43", "20:09", "21:46"},
		},
		{"Invalid day", 0, true, AthanTimes{}},
		{"Invalid day", 32, true, AthanTimes{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getAthanTimesForDay(athanCacheJson, tt.day)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("getAthanTimesForDay() error = %v", err)
			}
		})
	}
}

func Test_convertTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"12 hour", "08:34 (AEST)", "08:34"},
		{"24 hour", "20:34 (AEST)", "20:34"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTime(tt.input); got != tt.expected {
				t.Errorf("convertTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_GetNextAthan(t *testing.T) {
	testAthanCache := "./testing_files/athan_test.json"

	tests := []struct {
		name            string
		athanTestCache  string
		mockCurrentTime time.Time
		expected        string
		wantErr         bool
	}{
		{
			"Valid",
			testAthanCache,
			time.Date(2021, 10, 1, 5, 0, 0, 0, time.UTC),
			"Fajr in 3 hours and 56 minutes",
			false,
		},
		{
			"Loop to next day",
			testAthanCache,
			time.Date(2021, 10, 1, 5, 23, 0, 0, time.UTC),
			"Fajr in 4 hours and 55 minutes",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if athanString, err := GetNextAthan(tt.athanTestCache, tt.mockCurrentTime); (err != nil) != tt.wantErr && tt.expected != *athanString {
				t.Errorf("GetNextAthan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
