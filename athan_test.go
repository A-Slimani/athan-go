package main

import (
	"encoding/json"
	"os"
	"testing"
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

	// check if the file athan file is correct

	os.Remove(locationCacheJson)
	os.Remove(athanCacheJson)
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
