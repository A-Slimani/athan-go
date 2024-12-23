package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	locationCacheJson := os.TempDir() + "/location.json"
	locationCacheCheck, err := os.Stat(locationCacheJson)

	athanCacheJson := os.TempDir() + "/athan.json"
	athanCacheCheck, err := os.Stat(athanCacheJson)

	allFlag := flag.Bool("all", false, "Print all athan times")
	forceFlag := flag.Bool("force", false, "force cache update (use if cache is outdated or bugging)")
	setLocationFlag := flag.Bool("set-location", false, "set location manually")

	flag.Parse()

	switch {
	case *setLocationFlag:
	case *allFlag:
		AllAthanTimes(athanCacheJson)
	case *forceFlag:
		CacheLocation(locationCacheJson)
		CacheAthanTimes(locationCacheJson, athanCacheJson)
		fmt.Println("Cache updated")
		GetNextAthan(athanCacheJson)
	default:
		GetNextAthan(athanCacheJson)
	}

	if os.IsNotExist(err) {
		// make this an option
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Location cache not found, would you like to set it automatically? (y/n): ")
		CacheLocation(locationCacheJson)
	} else if locationCacheCheck.ModTime().AddDate(0, 0, 1).Before(time.Now()) { // check this with unit tests
		CacheLocation(locationCacheJson)
	} else if err != nil {
		fmt.Println("Error checking file: ", err)
	}

	if err != nil {
		if os.IsNotExist(err) {
			CacheAthanTimes(locationCacheJson, athanCacheJson)
		} else {
			fmt.Println("Error checking file: ", err)
		}
	} else {
		newMonthCheck := athanCacheCheck.ModTime().Day()
		if newMonthCheck == 1 {
			CacheAthanTimes(locationCacheJson, athanCacheJson)
		}
	}
}
