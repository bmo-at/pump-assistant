package main

import (
	"database/sql"
	"net/http"
	"time"

	//"log"
	"github.com/gin-gonic/gin"
	//"github.com/gin-contrib/cors"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase() *sql.DB {
	database, err := sql.Open("sqlite3", "./fuel_data.db")
	if err != nil {
		panic(err)
	}
	create, err := database.Prepare("CREATE TABLE IF NOT EXISTS fuel (id INTEGER PRIMARY KEY, datetime INTEGER, km_on_tank REAL, litres_filled_up REAL, price_payed_total REAL)")
	if err != nil {
		panic(err)
	}
	_, err = create.Exec()
	if err != nil {
		panic(err)
	}
	return database
}

var database = InitializeDatabase()

func main() {
	defer database.Close()
	router := SetupRouter()

	router.Run("0.0.0.0:8080")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/insert", Insert)
	router.GET("/values/all", GetAll)
	return router
}

type InsertionPayload struct {
	KM_ON_LAST_TANK   float64 `json:"km_on_last_tank", binding:"required"`
	LITRES_FILLED_UP  float64 `json:"litres_filled_up", binding:"required"`
	PRICE_PAYED_TOTAL float64 `json:"price_payed_total", binding:"required"`
}

func Insert(ctx *gin.Context) {
	var payload InsertionPayload

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Invalid JSON",
		})
		ctx.Abort()
		return
	}

	insert, err := database.Prepare("INSERT INTO fuel(datetime, km_on_tank, litres_filled_up, price_payed_total) VALUES (?, ?, ?, ?)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "Could not create statement",
		})
		ctx.Abort()
		return

	}
	result, err := insert.Exec(time.Now().Unix(), payload.KM_ON_LAST_TANK, payload.LITRES_FILLED_UP, payload.PRICE_PAYED_TOTAL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "Could not insert row",
		})
		ctx.Abort()
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

type Row struct {
	DateTime                time.Time
	Kilometers_On_Last_Tank float64
	Litres_Filled_Up        float64
	Price_Payed_Total       float64
}

func GetAll(ctx *gin.Context) {
	rows, _ := database.Query("SELECT * FROM fuel")

	var id int
	var unix_time int
	var km_on_tank float64
	var litres_filled_up float64
	var price_payed_total float64

	var result []Row = []Row{}

	for rows.Next() {
		rows.Scan(&id, &unix_time, &km_on_tank, &litres_filled_up, &price_payed_total)
		result = append(result, Row{
			DateTime:                time.Unix(int64(unix_time), int64(0)),
			Kilometers_On_Last_Tank: km_on_tank,
			Litres_Filled_Up:        litres_filled_up,
			Price_Payed_Total:       price_payed_total,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
