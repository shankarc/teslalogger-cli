package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init(){
	var envFile = os.Getenv("HOME") + "/.env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Could not load %s, Err: %s", envFile, err)
	}
}

func GetDSN() string {
	var username = os.Getenv("MYSQL_USER")
	var password = os.Getenv("MYSQL_PASSWORD")
	var host = os.Getenv("MYSQL_HOST")
	var port = os.Getenv("MYSQL_PORT")
	var dbname = os.Getenv("MYSQL_DATABASE")
	var hostname = fmt.Sprintf("%s:%s", host, port)

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func GetSQLHost() string {
	return os.Getenv("MYSQL_HOST")
}

func GetCostPerKwh() float64 {
	cost, err := strconv.ParseFloat(os.Getenv("PRICE_PER_KWH"), 64)
	if err != nil {
		log.Fatal("PRICE_PER_KWH is not set in .env file: ", err)
	}
	return cost
}
