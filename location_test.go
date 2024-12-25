package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func Test_CacheLocation(t *testing.T) {
	locationCacheJson := "./testing_files/location_test.json"
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			"Automically set location",
			"y\n",
			false,
		},
		{
			"Manually set location",
			"n\nSydney, Australia\n",
			false,
		},
		{
			"Invalid input",
			"invalid\n",
			true,
		},
		{
			"Empty input",
			"\n",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			err := CacheLocation(reader, locationCacheJson)
			if (err != nil) != tt.wantErr {
				t.Errorf("CacheLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeLocationInfoToJson(t *testing.T) {
	locationCacheJson := "./testing_files/location_test.json"
	type args struct {
		location Location
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		checkFile    bool
		expectedFile string
	}{
		{
			"Valid location",
			args{Location{"Sydney", "Australia"}},
			false,
			true,
			locationCacheJson,
		},
		{
			"Valid location w code",
			args{Location{"Sydney", "AU"}},
			false,
			true,
			locationCacheJson,
		},
		{
			"Invalid location Missing country",
			args{Location{"Sydney", ""}},
			true,
			false,
			"",
		},
		{
			"Invalid location Missing city",
			args{Location{"", "Australia"}},
			true,
			false,
			"",
		},
		{
			"Invalid location spaces",
			args{Location{" ", " "}},
			true,
			false,
			"",
		},
		{
			"Invalid location empty strings",
			args{Location{"", ""}},
			true,
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeLocationInfoToJson(tt.args.location, locationCacheJson); (err != nil) != tt.wantErr {
				t.Errorf("writeLocationInfoToJson() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expectedFile != "" {
				if _, err := os.Stat(tt.expectedFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s to be created, but it does not exist", tt.expectedFile)
				} else if err != nil {
					t.Errorf("Error checking file %s: %v", tt.expectedFile, err)
				}
			}
		})
	}
}
