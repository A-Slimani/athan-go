package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	locationCacheJson := os.TempDir() + "/location.json"
	locationCacheCheck, err := os.Stat(locationCacheJson)

	if os.IsNotExist(err) {
		cacheLocation(locationCacheJson)
	} else if locationCacheCheck.ModTime().AddDate(0, 0, 1).Before(time.Now()) {
		cacheLocation(locationCacheJson)
	} else if err != nil {
		fmt.Println("Error checking file: ", err)
	}

	athanCacheJson := os.TempDir() + "/athan.json"
	athanCacheCheck, err := os.Stat(athanCacheJson)

	if err != nil {
		if os.IsNotExist(err) {
			cacheAthan(locationCacheJson, athanCacheJson)
		} else {
			fmt.Println("Error checking file: ", err)
		}
	} else {
		_, _, newMonthCheck := athanCacheCheck.ModTime().Date()
		if newMonthCheck == 1 {
			cacheAthan(locationCacheJson, athanCacheJson)
		}
	}

	allFlag := flag.Bool("all", false, "Print all athan times")
	forceFlag := flag.Bool("force", false, "force cache update")

	flag.Parse()

	if *allFlag {
		allAthanTimes(athanCacheJson)
	} else if *forceFlag {
		cacheLocation(locationCacheJson)
		cacheAthan(locationCacheJson, athanCacheJson)
		getNextAthan(athanCacheJson)
	} else {
		getNextAthan(athanCacheJson)
	}
}
