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
	event struct {
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

	latlong struct {
		X float32
		Y float32
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
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.BodyLimit("16K"))

	e.POST("/api/event", createEvent)
	e.PUT("/api/event/:id", updateEvent)
	e.GET("/api/event/:id", getEvent)
	e.GET("/api/events", getEvents)
	e.DELETE("/api/event:/id", deleteEvent)

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

func getEvents(c echo.Context) error {

	var evnt []event
	err := db.Select(&evnt, "SELECT * FROM events LIMIT 10")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, evnt)
}

func getEvent(c echo.Context) error {
	id := c.Param("id")

	var evnt event
	err := db.Get(&evnt, "SELECT * FROM events WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, "")
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	// defer stmt.Close()

	return c.JSON(http.StatusOK, evnt)
}

func createEvent(c echo.Context) error {
	evnt := new(event)
	if err := c.Bind(evnt); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	query := `INSERT into events(
			name, "imageUrl", "dateStart", "dateEnd", "appDeadline",
			latlng, location, description, price, "maxPlayers",
			owner, type
		)
		VALUES(
			:name, :imageUrl, :dateStart, :dateEnd, :appDeadline,
			:latlng, :location, :description, :price, :maxPlayers,
			:owner, :type
		)
		RETURNING id`

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.Get(&id, evnt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, id)
}

func deleteEvent(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "")
}

func updateEvent(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, "")
}
