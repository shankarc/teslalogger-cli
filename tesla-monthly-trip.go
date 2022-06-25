package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"teslalogger/cli/utils"
)

type Miles struct {
	miles, startmile, endmile, day int
}

func PrintTrip(m *map[string][]Miles) {
	// Go does not range keys in sorted order
	// or in the order appended
	var dates []string
	for date, _ := range *m {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	for _, yearMonth := range dates {
		//  dereference the map and access (yyyy/mm)
		miles := (*m)[yearMonth]
		var odometerArray []int
		var milesArray []int
		var days, prevDayCounted int
		for _, mile := range miles {
			odometerArray = append(odometerArray, mile.startmile, mile.endmile)
			milesArray = append(milesArray, mile.miles)
			// Count the day once
			if mile.day != prevDayCounted {
				days += 1
				prevDayCounted = mile.day
			}
		}
		sort.Ints(odometerArray)
		min, max := odometerArray[0], odometerArray[len(odometerArray)-1]
		totalMiles := max - min
		fmt.Printf("%s, %d, %d, %d days, %d miles\n", yearMonth, min, max, days, totalMiles)
	}

}

func main() {

	var (
		startOdometer, endOdometer int
		miles                      int
		date                       string
	)

	type Date struct {
		year, month, day int
	}

	montlyMiles := make(map[string][]Miles)

	db, err := sql.Open("mysql", utils.GetDSN())

	if err != nil {
		log.Fatal("Could not open mysql db: ", err)
	}

	defer db.Close()

	query := `SELECT DATE_FORMAT(startdate, '%Y-%m-%dT%H:%i:%sZ') as date,
		CAST(km_diff * 0.6213712 AS UNSIGNED) AS miles,
        CAST(StartKm * 0.6213712 AS UNSIGNED) AS start_miles,
        CAST(EndKm * 0.6213712 AS UNSIGNED) AS end_miles FROM trip`

	result, err := db.Query(query)

	if err != nil {
		log.Fatal("Could not Query the database: ", err)
	}

	for result.Next() {

		err = result.Scan(&date, &miles, &startOdometer, &endOdometer)
		if err != nil {
			log.Fatal(err)
		}

		dateTime, err := time.Parse(time.RFC3339, date)
		if err != nil {
			log.Fatal(err)
		}

		currentDate := Date{day: int(dateTime.Day()),
			month: int(dateTime.Month()),
			year:  int(dateTime.Year())}

		if miles > 1000 {
			// Anomaly in db
			continue
		}
		yearMonth := fmt.Sprintf("%4d/%02d",
			currentDate.year,
			currentDate.month)

		var miles Miles = Miles{
			miles:     miles,
			startmile: startOdometer,
			endmile:   endOdometer,
			day:       currentDate.day,
		}
		montlyMiles[yearMonth] = append(montlyMiles[yearMonth], miles)
	}

	PrintTrip(&montlyMiles)
}
