package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

const filename = "contacts.json"

func main() {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			log.Panic(err)
		}

		err = json.NewEncoder(file).Encode([]Contact{})
		file.Close()

		if err != nil {
			log.Panic(err)
		}
	}

	r := gin.Default()

	r.POST("/contact/add", func(c *gin.Context) {
		var contacts []Contact

		file, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			return
		}

		err = json.NewDecoder(file).Decode(&contacts)
		file.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read contacts"})
			return
		}

		var req Contact

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		contacts = append(contacts, req)

		file, err = os.Create(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			return
		}

		err = json.NewEncoder(file).Encode(contacts)
		file.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save contact"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "contact added successfully"})
	})

	r.GET("/contacts/get", func(c *gin.Context) {
		name := c.Query("name")
		email := c.Query("email")

		if name == "" && email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "either name or email required",
			})
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open contacts"})
			return
		}

		var contacts []Contact

		err = json.NewDecoder(file).Decode(&contacts)
		file.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read contacts"})
			return
		}

		var results []Contact

		for _, contact := range contacts {
			nameMatches := name == "" || contact.Name == name
			emailMatches := email == "" || contact.Email == email

			if nameMatches && emailMatches {
				results = append(results, contact)
			}
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"contacts": results,
		})
	})

	r.GET("/contacts", func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open contacts"})
			return
		}

		var contacts []Contact

		err = json.NewDecoder(file).Decode(&contacts)
		file.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read contacts"})
			return
		}

		total := len(contacts)
		start := (page - 1) * limit
		end := start + limit

		if start >= total {
			c.JSON(http.StatusOK, gin.H{
				"page":     page,
				"limit":    limit,
				"total":    total,
				"contacts": []Contact{},
			})
			return
		}

		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"contacts": contacts[start:end],
		})
	})

	r.Run()
}
