package main

import (
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

}

func Test_convertTime(t *testing.T) {
	type args struct {
		athanTime string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTime(tt.args.athanTime); got != tt.want {
				t.Errorf("convertTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
