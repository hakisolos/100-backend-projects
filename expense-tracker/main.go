package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpenseLogs struct {
	ID     string    `json:"id"`
	Type   string    `json:"type"`
	Amount int       `json:"amount"`
	Note   string    `json:"note"`
	Date   time.Time `json:"-"`
}

func main() {
	var err error
	_, err = os.Stat("logs.json")
	if err != nil {
		file, err := os.Create("logs.json")
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(file).Encode([]ExpenseLogs{})
		file.Close()
	}

	app := gin.Default()

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api running",
		})
	})

	app.POST("/logs/add", func(c *gin.Context) {
		var req ExpenseLogs
		err := c.ShouldBindJSON(&req)
		req.Date = time.Now()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
		}
		err = AddExpense(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "an error occured while adding transaction",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "log added",
		})
	})
	app.Run()

}

func AddExpense(exp ExpenseLogs) error {
	var logs []ExpenseLogs
	file, err := os.Open("logs.json")
	if err != nil {
		return err
	}
	json.NewDecoder(file).Decode(&logs)
	file.Close()

	logs = append(logs, exp)
	file, err = os.Create("logs.json")
	if err != nil {
		return err
	}
	json.NewEncoder(file).Encode(&logs)
	file.Close()
	return nil
}
