package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	locationCacheJson := os.TempDir() + "/location.json"
	_, err := os.Stat(locationCacheJson)

	if os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		if err := CacheLocation(reader, locationCacheJson); err != nil {
			fmt.Println("Error caching location: ", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error checking file: ", err)
	}

	athanCacheJson := os.TempDir() + "/athan.json"
	athanCacheCheck, err := os.Stat(athanCacheJson)

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

	// allFlag := flag.Bool("all", false, "Print all athan times")
	// forceFlag := flag.Bool("force", false, "force cache update (use if cache is outdated or bugging)")
	// setLocationFlag := flag.Bool("set-location", false, "set location manually")

	// flag.Parse()

	// switch {
	// case *setLocationFlag:
	// case *allFlag:
	// 	AllAthanTimes(athanCacheJson)
	// case *forceFlag:
	// 	reader := bufio.NewReader(os.Stdin)
	// 	if err := CacheLocation(reader); err != nil {
	// 		fmt.Println("Error caching location: ", err)
	// 		return
	// 	}
	// 	CacheAthanTimes(locationCacheJson, athanCacheJson)
	// 	fmt.Println("Cache updated")
	// 	GetNextAthan(athanCacheJson)
	// default:
	// 	GetNextAthan(athanCacheJson)
	// }
}
