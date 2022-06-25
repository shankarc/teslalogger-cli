package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"teslalogger/cli/utils"
)

func main() {

	costPerKwh := utils.GetCostPerKwh()

	db, err := sql.Open("mysql", utils.GetDSN())

	if err != nil {
		log.Fatal("Could not open mysql db: ", err.Error())
	}

	defer db.Close()

	update := fmt.Sprintf("UPDATE chargingstate set cost_per_kwh = %f,"+
		"cost_currency ='USD',"+
		"cost_total = %f * charge_energy_added WHERE cost_total is NULL AND "+
		"fast_charger_type != 'Tesla';", costPerKwh, costPerKwh)

	res, err := db.Exec(update)

	if err != nil {
		log.Fatal("Could not Query the database: ", err.Error())
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Updated %d rows\n", affectedRows)
}
