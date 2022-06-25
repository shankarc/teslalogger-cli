package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"teslalogger/cli/utils"
)

func getArgs() (time.Time, time.Time, bool) {

	from := flag.String("s", "", "The date to start calculating super charger cost e.g 2022-05-23")
	to := flag.String("e", "", "The end date e.g 2022-06-01 inclusive.")
	verbose := flag.Bool("v", false, "Verbose")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Will print the super charger cost, the charging cost has to be entered manually in teslalogger.\n")
		flag.PrintDefaults()

	}

	flag.Parse()

	if *from == "" || *to == "" {
		flag.Usage()
		os.Exit(0)
	}

	fromDate, err := time.Parse("2006-01-02", *from)

	if err != nil {
		log.Fatal("Could not parse -s argument: ", err)
	}

	toDate, err := time.Parse("2006-01-02", *to)

	if err != nil {
		log.Fatal("Could not parse -e argument: ", err)
	}

	return fromDate, toDate.AddDate(0, 0, 1), *verbose
}

func main() {

	var (
		energyAdded, cost, totalCost float64
		sdate, edate                 string
		charging                     int
	)

	fromDate, toDate, verbose := getArgs()

	db, err := sql.Open("mysql", utils.GetDSN())

	if err != nil {
		log.Fatal("Could not open mysql db: ", err)
	}

	defer db.Close()

	query := `SELECT DATE_FORMAT(StartDate,'%Y-%m-%dT%H:%i:%sZ') as sdate,
			DATE_FORMAT(EndDate,'%Y-%m-%dT%H:%i:%sZ') as edate,
			charge_energy_added, cost_total FROM teslalogger.chargingstate 
			WHERE fast_charger_brand='Tesla';`

	result, err := db.Query(query)

	if err != nil {
		log.Fatal("Could not Query the database: ", err)
	}

	for result.Next() {
		err = result.Scan(&sdate, &edate, &energyAdded, &cost)

		if err != nil {
			log.Fatal(err)
		}
		// Convert datetime to time structure
		start_date_time, err := time.Parse(time.RFC3339, sdate)
		if err != nil {
			log.Fatal(err)
		}

		end_date_time, err := time.Parse(time.RFC3339, edate)
		if err != nil {
			log.Fatal(err)
		}

		if start_date_time.After(fromDate) {

			if start_date_time.After(toDate) {
				break
			}

			if verbose {
				start_time := start_date_time.Format("2006-01-02 03:04:05 PM")
				duration := end_date_time.Sub(start_date_time).Seconds() / 60
				fmt.Printf("%s, Duration: %.2f minutes Energy Added: %.2f Kwh, Cost: $%.2f\n",
					start_time, duration, energyAdded, cost)
			}

			totalCost += cost
			charging += 1
		}

	}

	fmt.Printf("From %s upto %s, Number of charges:%d, Total Cost:$%.2f\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"), charging, totalCost)
}
