# About
Command-line utilities uses [teslalogger](https://github.com/bassmaster187/TeslaLogger) database to extract Tesla driving information.
## Prerequisite 
It is assumed that teslalogger is running in your LAN and the following information is provided in
__$HOME/.env__ file  

These information is provided in teslaloggers README file.

```
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_DATABASE=
MYSQL_HOST=
MYSQL_PORT=
```
This is to set charging cost at home.
```
PRICE_PER_KWH=0.128
```
Tested on MacOs, Linux, and Windows(gitbash)

## Init
If you have to initialize
```
go mod init teslalogger/cli  
go mod tidy                                             
```

## Build
In unix like system.  
```
./build.sh 
```

builds the executable in _bin/_  dir, move them where they can be accessible ($PATH)

For windows.  
```
build.bat
```

You can also build individually.  
```
go build -o tesla-miles tesla-miles.go
```

## tesla-miles
Prints the miles covered till the day with start odometer, end odometer, max speed, and the miles travelled that day.  
```
tesla-miles  

2022/06/15 Wednesday, 60196, 60279, 84 mph, 73 miles
2022/06/16 Thursday, 60279, 60360, 88 mph, 71 miles
2022/06/17 Friday, 60360, 60442, 84 mph, 72 miles
```
## tesla-create-trip-gpx
Creates _trip.gpx_ file, which can be used to create a google map KML file.
```
tesla-create-trip-gpx  -s 2022-05-23 -e 2022-06-01 -o trip.gpx
```

## tesla-super-charger-cost
Prints the super charger cost, the charging cost has to be entered manually in teslalogger.   

```
tesla-super-charger-cost -s 2022-05-23  -e  2022-06-01 

From 2022-05-23 upto 2022-06-02, Number of charges:20, Total Cost:$237.32
```

argument _-v_ to print all the charger stops in between.

```
tesla-super-charger-cost  -s 2022-05-23 -e 2022-05-25 -v
2022-05-23 02:11:41 PM, Duration: 28.32 minutes Energy Added: 43.85 Kwh, Cost: $17.55
2022-05-23 04:26:51 PM, Duration: 10.05 minutes Energy Added: 22.08 Kwh, Cost: $10.12
2022-05-23 05:48:13 PM, Duration: 9.32 minutes Energy Added: 24.80 Kwh, Cost: $9.88
2022-05-23 08:30:04 PM, Duration: 20.15 minutes Energy Added: 39.57 Kwh, Cost: $15.12
2022-05-25 10:16:48 AM, Duration: 19.13 minutes Energy Added: 32.44 Kwh, Cost: $12.96
2022-05-25 01:57:04 PM, Duration: 10.37 minutes Energy Added: 19.88 Kwh, Cost: $7.14
2022-05-25 07:59:32 PM, Duration: 16.50 minutes Energy Added: 38.85 Kwh, Cost: $9.50
From 2022-05-23 upto 2022-05-26, Number of charges:7, Total Cost:$82.27
```

## tesla-monthly-trip
Prints the miles covered in a month with start odometer, end odometer, number of days and miles travelled.
```
$ tesla-monthly-trip 

2022/03, 51499, 53746, 24 days, 2247 miles
2022/04, 53746, 55736, 24 days, 1990 miles
2022/05, 55736, 58887, 26 days, 3151 miles
2022/06, 58887, 60442, 12 days, 1555 miles
```
## tesla-set-charging-cost
Sets the charging cost in the teslalogger database

```
tesla-set-charging-cost
```
__Notes:__ 

* It only sets for home charger, not Tesla Supercharger. 
* It uses _PRICE_PER_KWH_ variable is .env file.
