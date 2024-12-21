package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func cacheLocation(locationCacheJson string) {

	resp, err := http.Get("https://ipinfo.io")
	if err != nil {
		fmt.Println("Request error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
	}

	f, err := os.Create(locationCacheJson)
	if err != nil {
		fmt.Println("Error creating file: ", err)
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}
