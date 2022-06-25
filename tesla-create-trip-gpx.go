package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"teslalogger/cli/utils"
)

func getArgs() (time.Time, time.Time, string) {

	from := flag.String("s", "", "The start date  e.g 2022-05-23 inclusive.")
	to := flag.String("e", "", "The end date  e.g 2022-06-01 inclusive.")
	output := flag.String("o", "", "The file name to save gpx file.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Will create a gpx (GPS Exchange Format) file for the trip taken, given the start and end date.\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *from == "" || *to == "" || *output == "" {
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

	// func (t Time) AddDate(years int, months int, days int) Time
	// To make toDate inclusive
	return fromDate, toDate.AddDate(0, 0, 1), *output
}

func saveGpx(endpoint, filePath string) {

	response, err := http.Get(endpoint)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Could not get GPX data from teslalogger: ", err)
	}

	f, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666,
	)

	defer f.Close()

	_, err = f.Write(responseData)

	if err != nil {
		log.Fatal("Could not write GPX data to the file: ", err)
	}
}

func getPosIDs(fromDate, toDate time.Time) (int, int) {
	var (
		date                 string
		posid                int
		startPosid, endPosid int
	)

	db, err := sql.Open("mysql", utils.GetDSN())

	if err != nil {
		log.Fatal("Could not open mysql db: ", err)
	}

	defer db.Close()

	query := `SELECT id, DATE_FORMAT(Datum,'%Y-%m-%dT%H:%i:%sZ') as date  FROM teslalogger.pos;`
	result, err := db.Query(query)

	if err != nil {
		log.Fatal("Could not Query the database: ", err)
	}

	for result.Next() {
		err = result.Scan(&posid, &date)

		if err != nil {
			log.Fatal(err)
		}

		// Convert datetime to time structure
		date_time, err := time.Parse(time.RFC3339, date)
		if err != nil {
			log.Fatal(err)
		}

		// If date is after fromDate and not after toDate:
		if date_time.After(fromDate) {
			if date_time.After(toDate) {
				endPosid = posid
				break
			}
			if startPosid == 0 {
				startPosid = posid
			}
		}

	}
	return startPosid, endPosid
}

func main() {
	startDate, endDate, fileName := getArgs()
	start, end := getPosIDs(startDate, endDate)
	url := utils.GetSQLHost()
	url = fmt.Sprintf("http://%s:5010/export/trip?from=%d&to=%d&carID=1", url, start, end)
	saveGpx(url, fileName)
}
