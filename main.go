package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

type (
	Event struct {
		ID          int64     `json:"id" db:"id"`
		Name        string    `json:"name,omitempty" db:"name"`
		ImageURL    string    `json:"imageURL,omitempty" db:"imageUrl"`
		DateStart   time.Time `json:"date_start,omitempty" db:"dateStart"`
		DateEnd     time.Time `json:"date_end,omitempty" db:"dateEnd"`
		AppDeadline time.Time `json:"app_deadline,omitempty" db:"appDeadline"`
		Latlong     string    `json:"latlong,omitempty" db:"latlng"`
		Location    string    `json:"location,omitempty" db:"location"`
		Description string    `json:"descr,omitempty" db:"description"`
		Price       string    `json:"price,omitempty" db:"price"`
		MaxPlayers  int64     `json:"max_players,omitempty" db:"maxPlayers"`
		OwnerID     int64     `json:"owner_id,omitempty" db:"owner"`
		EventTypeID int64     `json:"event_type_id,omitempty" db:"type"`
	}
)

func main() {
	initDB()
	initEcho()
}

func initEcho() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.BodyLimit("16K"))

	e.POST("/event", CreateEvent)
	e.PUT("/event/:id", UpdateEvent)
	e.GET("/event/:id", GetEvent)
	e.GET("/events", GetEvents)
	e.DELETE("/event:/id", DeleteEvent)

	e.Logger.Fatal(e.Start(":8080"))
}

func initDB() {
	var err error
	db, err = sqlx.Open("postgres", "user=calendar password=calendar host=192.168.1.4 dbname=calendar sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func GetEvents(c echo.Context) error {

	var evnt []Event
	err := db.Select(&evnt, "SELECT * FROM events")
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, evnt)
}

func GetEvent(c echo.Context) error {
	id := c.Param("id")

	var evnt Event
	err := db.Get(&evnt, "SELECT * FROM events WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, "")
		}
		log.Fatal(err)
	}

	// defer stmt.Close()

	return c.JSON(http.StatusOK, evnt)
}

func CreateEvent(c echo.Context) error {
	evnt := new(Event)
	if err := c.Bind(evnt); err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT into events VALUES RETURNING id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(evnt)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusCreated, result)
}

func DeleteEvent(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "")
}

func UpdateEvent(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "")
}
