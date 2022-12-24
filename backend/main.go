package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

type payload struct {
	Id    int
	Type  string `json:"_type"`
	T     string
	Batt  int
	Lat   float64
	Lon   float64
	Tid   string
	Tst   int
	Topic string
}

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./data/owntracks.db")
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	createLocationsTable := `
	CREATE TABLE IF NOT EXISTS locations 
	(
	  id INTEGER PRIMARY KEY,
	  t  TEXT NOT NULL,                -- trigger
	  batt INTEGER,
	  topic TEXT, 
	  tid TEXT NOT NULL,                -- tracker id, identifies the device
	  lat REAL NOT NULL,                -- latitude
	  lon REAL NOT NULL,                -- longitude
	  tst INTEGER NOT NULL UNIQUE       -- timestamp, only one datapoint per tst allowed
	);`

	if _, err := db.Exec(createLocationsTable); err != nil {
		log.Fatal("Failed to create table: ", err)
	}
}

func postLocation(ctx echo.Context) error {
	var p payload

	err := json.NewDecoder(ctx.Request().Body).Decode(&p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if p.Type != "location" {
		return echo.NewHTTPError(http.StatusBadRequest, "Only location type requests allowed")
	}

	log.Println(p)

	res, err := db.Exec("INSERT INTO locations (t, batt, topic, lat, lon, tst, tid) VALUES (?, ?, ?, ?, ?, ?, ?);", p.T, p.Batt, p.Topic, p.Lat, p.Lon, p.Tst, p.Tid)
	if err != nil {
		log.Println("Insert error: ", err)
		return ctx.String(http.StatusOK, "")
	}

	id, _ := res.LastInsertId()
	return ctx.String(http.StatusOK, strconv.Itoa(int(id)))
}

func getLocations(ctx echo.Context) error {
	start := ctx.QueryParam("start")
	end := ctx.QueryParam("end")

	query := "SELECT * from locations"
	if start != "" && end != "" {
		startTime, _ := time.Parse(time.RFC3339, start)
		endTime, _ := time.Parse(time.RFC3339, end)
		// TODO errors

		log.Println(startTime.Unix(), endTime.Unix())

		query += fmt.Sprintf(" where tst > %d and tst < %d", startTime.Unix(), endTime.Unix())
	}

	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	data := []payload{}
	for rows.Next() {
		i := payload{}
		err = rows.Scan(&i.Id, &i.T, &i.Batt, &i.Topic, &i.Tid, &i.Lat, &i.Lon, &i.Tst)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		data = append(data, i)
	}

	return ctx.JSON(http.StatusOK, data)
}

func main() {
	initDatabase()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/locations", getLocations)
	e.POST("/locations", postLocation)

	e.Logger.Fatal(e.Start(":8091"))
}
