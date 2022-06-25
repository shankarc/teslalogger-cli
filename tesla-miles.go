package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"teslalogger/cli/utils"
)

type Date struct {
	year    int
	month   int
	day     int
	weekday string
}

type Miles struct {
	startOdometer    int
	previousOdometer int
	endOdometer      int
	maxSpeed         int
	totalMile        int
	date             Date
}

func (m *Miles) Print() {
	var missingMiles string

	if m.previousOdometer != m.startOdometer {
		missingMiles = fmt.Sprintf("missing %d miles", (m.startOdometer - m.previousOdometer))
	}

	yearMonthDay := fmt.Sprintf("%4d/%02d/%02d %s",
		m.date.year,
		m.date.month,
		m.date.day,
		m.date.weekday)

	fmt.Printf("%s, %d, %d, %d mph, %d miles %s\n",
		yearMonthDay, m.startOdometer,
		m.endOdometer,
		m.maxSpeed,
		m.totalMile,
		missingMiles)
}

func main() {

	db, err := sql.Open("mysql", utils.GetDSN())

	if err != nil {
		log.Fatal("Could not open mysql db: ", err)
	}

	defer db.Close()

	var (
		startOdometer, endOdometer, maxSpeed, totalMile int
		previousEndMile, odometer, speed                int
		date                                            string
		currentDate                                     Date
	)

	previousDate := Date{year: 1969, month: 12, day: 31}

	query := `SELECT DATE_FORMAT(pos.datum, '%Y-%m-%dT%H:%i:%sZ') as date, 
			CAST(odometer * 0.6213712 as UNSIGNED) as odometer,
			CAST(speed/1.609 as UNSIGNED)as mph FROM pos`

	result, err := db.Query(query)

	if err != nil {
		log.Fatal("Could not Query the database: ", err)
	}

	for result.Next() {

		err = result.Scan(&date, &odometer, &speed)
		if err != nil {
			log.Fatal(err)
		}

		date_time, err := time.Parse(time.RFC3339, date)
		if err != nil {
			log.Fatal(err)
		}

		currentDate = Date{day: int(date_time.Day()),
			month:   int(date_time.Month()),
			year:    int(date_time.Year()),
			weekday: date_time.Weekday().String()}

		if currentDate != previousDate {
			if odometer == 0 {
				odometer = endOdometer
			}
			totalMile = endOdometer - startOdometer
			if totalMile > 0 {
				milesForTheDay := Miles{
					previousOdometer: previousEndMile,
					startOdometer:    startOdometer,
					endOdometer:      endOdometer,
					maxSpeed:         maxSpeed,
					totalMile:        totalMile,
					date:             previousDate,
				}

				milesForTheDay.Print()
				previousEndMile = endOdometer
			}
			startOdometer = odometer
			maxSpeed = 0

		} else {
			endOdometer = odometer
			if speed > maxSpeed {
				maxSpeed = speed
			}
		}
		previousDate = currentDate
	}

	// residual
	totalMile = endOdometer - startOdometer
	if totalMile > 0 {
		milesForTheDay := Miles{
			previousOdometer: previousEndMile,
			startOdometer:    startOdometer,
			endOdometer:      endOdometer,
			maxSpeed:         maxSpeed,
			totalMile:        totalMile,
			date:             previousDate,
		}
		milesForTheDay.Print()
	}

}
