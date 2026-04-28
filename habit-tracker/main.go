package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Habit struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type HabitCompletion struct {
	ID            int    `json:"id"`
	HabitID       int    `json:"habit_id"`
	CompletedDate string `json:"completed_date"`
	CreatedAt     string `json:"created_at"`
}

type DB struct {
	Habits      []Habit           `json:"habits"`
	Completions []HabitCompletion `json:"completions"`
}

const dbFile = "db.json"

func loadDB() (DB, error) {
	var db DB

	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		db = DB{
			Habits:      []Habit{},
			Completions: []HabitCompletion{},
		}
		return db, saveDB(db)
	}

	file, err := os.Open(dbFile)
	if err != nil {
		return db, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&db)
	return db, err
}

func saveDB(db DB) error {
	file, err := os.Create(dbFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(db)
}

func nextHabitID(habits []Habit) int {
	maxID := 0
	for _, habit := range habits {
		if habit.ID > maxID {
			maxID = habit.ID
		}
	}
	return maxID + 1
}

func nextCompletionID(completions []HabitCompletion) int {
	maxID := 0
	for _, completion := range completions {
		if completion.ID > maxID {
			maxID = completion.ID
		}
	}
	return maxID + 1
}

func findHabit(habits []Habit, id int) bool {
	for _, habit := range habits {
		if habit.ID == id {
			return true
		}
	}
	return false
}

func main() {
	r := gin.Default()

	r.POST("/habits", func(c *gin.Context) {
		var body struct {
			Name string `json:"name"`
		}

		if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		db, err := loadDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load database"})
			return
		}

		habit := Habit{
			ID:        nextHabitID(db.Habits),
			Name:      body.Name,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		db.Habits = append(db.Habits, habit)

		if err := saveDB(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save database"})
			return
		}

		c.JSON(http.StatusCreated, habit)
	})

	r.GET("/habits", func(c *gin.Context) {
		db, err := loadDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load database"})
			return
		}

		c.JSON(http.StatusOK, db.Habits)
	})

	r.POST("/habits/:id/complete", func(c *gin.Context) {
		habitID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid habit id"})
			return
		}

		db, err := loadDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load database"})
			return
		}

		if !findHabit(db.Habits, habitID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
			return
		}

		today := time.Now().Format("2006-01-02")

		for _, completion := range db.Completions {
			if completion.HabitID == habitID && completion.CompletedDate == today {
				c.JSON(http.StatusConflict, gin.H{"error": "habit already completed today"})
				return
			}
		}

		completion := HabitCompletion{
			ID:            nextCompletionID(db.Completions),
			HabitID:       habitID,
			CompletedDate: today,
			CreatedAt:     time.Now().Format(time.RFC3339),
		}

		db.Completions = append(db.Completions, completion)

		if err := saveDB(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save database"})
			return
		}

		c.JSON(http.StatusCreated, completion)
	})

	r.GET("/habits/:id/streak", func(c *gin.Context) {
		habitID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid habit id"})
			return
		}

		db, err := loadDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load database"})
			return
		}

		if !findHabit(db.Habits, habitID) {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
			return
		}

		completedDates := map[string]bool{}

		for _, completion := range db.Completions {
			if completion.HabitID == habitID {
				completedDates[completion.CompletedDate] = true
			}
		}

		streak := 0
		currentDate := time.Now()

		for {
			date := currentDate.Format("2006-01-02")

			if !completedDates[date] {
				break
			}

			streak++
			currentDate = currentDate.AddDate(0, 0, -1)
		}

		c.JSON(http.StatusOK, gin.H{
			"habit_id": habitID,
			"streak":   streak,
		})
	})
	r.Run(":8080")
}
