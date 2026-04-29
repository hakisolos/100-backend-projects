package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Card struct {
	ID    int64  `json:"id"`
	Front string `json:"front"`
	Back  string `json:"back"`
	Known bool   `json:"known"`
}

type Deck struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Cards []*Card `json:"cards"`
}

var (
	lock       sync.Mutex
	decks            = map[int64]*Deck{}
	nextDeckID int64 = 1
	nextCardID int64 = 1
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "api running"})
	})
	app.GET("/decks", listDecks)
	app.POST("/decks", createDeck)
	app.GET("/decks/:id", getDeck)
	app.PUT("/decks/:id", updateDeck)
	app.DELETE("/decks/:id", deleteDeck)
	app.GET("/decks/:id/cards", listCards)
	app.POST("/decks/:id/cards", createCard)
	app.GET("/decks/:id/cards/:cardId", getCard)
	app.PUT("/decks/:id/cards/:cardId", updateCard)
	app.DELETE("/decks/:id/cards/:cardId", deleteCard)
	app.GET("/decks/:id/cards/random", randomCard)
	app.POST("/decks/:id/cards/:cardId/known", markKnown)
	app.POST("/decks/:id/cards/:cardId/unknown", markUnknown)
	app.Run(":8080")
}

func listDecks(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	result := make([]*Deck, 0, len(decks))
	for _, deck := range decks {
		result = append(result, deck)
	}
	c.JSON(http.StatusOK, result)
}

func createDeck(c *gin.Context) {
	var payload struct {
		Name string `json:"name"`
	}
	if c.BindJSON(&payload) != nil || payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	lock.Lock()
	deck := &Deck{ID: nextDeckID, Name: payload.Name, Cards: []*Card{}}
	decks[nextDeckID] = deck
	nextDeckID++
	lock.Unlock()
	c.JSON(http.StatusCreated, deck)
}

func getDeck(c *gin.Context) {
	deck, ok := fetchDeck(c.Param("id"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	c.JSON(http.StatusOK, deck)
}

func updateDeck(c *gin.Context) {
	var payload struct {
		Name string `json:"name"`
	}
	if c.BindJSON(&payload) != nil || payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	lock.Lock()
	deck, ok := decks[toID(c.Param("id"))]
	if !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	deck.Name = payload.Name
	lock.Unlock()
	c.JSON(http.StatusOK, deck)
}

func deleteDeck(c *gin.Context) {
	lock.Lock()
	id := toID(c.Param("id"))
	if _, ok := decks[id]; !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	delete(decks, id)
	lock.Unlock()
	c.Status(http.StatusNoContent)
}

func listCards(c *gin.Context) {
	deck, ok := fetchDeck(c.Param("id"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	c.JSON(http.StatusOK, deck.Cards)
}

func createCard(c *gin.Context) {
	var payload struct {
		Front string `json:"front"`
		Back  string `json:"back"`
	}
	if c.BindJSON(&payload) != nil || payload.Front == "" || payload.Back == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	lock.Lock()
	deck, ok := decks[toID(c.Param("id"))]
	if !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	card := &Card{ID: nextCardID, Front: payload.Front, Back: payload.Back}
	deck.Cards = append(deck.Cards, card)
	nextCardID++
	lock.Unlock()
	c.JSON(http.StatusCreated, card)
}

func getCard(c *gin.Context) {
	card, ok := fetchCard(c)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}
	c.JSON(http.StatusOK, card)
}

func updateCard(c *gin.Context) {
	var payload struct {
		Front string `json:"front"`
		Back  string `json:"back"`
	}
	if c.BindJSON(&payload) != nil || payload.Front == "" || payload.Back == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	lock.Lock()
	card, ok := fetchCard(c)
	if !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}
	card.Front = payload.Front
	card.Back = payload.Back
	lock.Unlock()
	c.JSON(http.StatusOK, card)
}

func deleteCard(c *gin.Context) {
	lock.Lock()
	deck, ok := decks[toID(c.Param("id"))]
	if !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	cardID := toID(c.Param("cardId"))
	for i, card := range deck.Cards {
		if card.ID == cardID {
			deck.Cards = append(deck.Cards[:i], deck.Cards[i+1:]...)
			lock.Unlock()
			c.Status(http.StatusNoContent)
			return
		}
	}
	lock.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
}

func randomCard(c *gin.Context) {
	deck, ok := fetchDeck(c.Param("id"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}
	if len(deck.Cards) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no cards in deck"})
		return
	}
	c.JSON(http.StatusOK, deck.Cards[rand.Intn(len(deck.Cards))])
}

func markKnown(c *gin.Context) {
	setKnownState(c, true)
}

func markUnknown(c *gin.Context) {
	setKnownState(c, false)
}

func setKnownState(c *gin.Context, state bool) {
	lock.Lock()
	card, ok := fetchCard(c)
	if !ok {
		lock.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}
	card.Known = state
	lock.Unlock()
	c.JSON(http.StatusOK, card)
}

func fetchDeck(idParam string) (*Deck, bool) {
	lock.Lock()
	defer lock.Unlock()
	deck, ok := decks[toID(idParam)]
	return deck, ok
}

func fetchCard(c *gin.Context) (*Card, bool) {
	deck, ok := fetchDeck(c.Param("id"))
	if !ok {
		return nil, false
	}
	cardID := toID(c.Param("cardId"))
	for _, card := range deck.Cards {
		if card.ID == cardID {
			return card, true
		}
	}
	return nil, false
}

func toID(param string) int64 {
	id, _ := strconv.ParseInt(param, 10, 64)
	return id
}
