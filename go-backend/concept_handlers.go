package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func addConcept(c *gin.Context) {
	var newConcept struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	if err := c.BindJSON(&newConcept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	concept := &Concept{
		GUID:          GUID(uuid.New().String()),
		Name:          newConcept.Name,
		Description:   newConcept.Description,
		Type:          newConcept.Type,
		Timestamp:     time.Now(),
		Relationships: []GUID{},
	}
	concept.Update(c.Request.Context())

	addNewConcept(concept)
	c.JSON(http.StatusOK, gin.H{
		"guid": concept.GUID,
		"cid":  string(concept.CID),
	})
}

func getConcept(c *gin.Context) {
	guid := GUID(c.Param("guid"))

	conceptMu.RLock()
	concept, exists := conceptMap[guid]
	conceptMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Concept not found"})
		return
	}
	c.JSON(http.StatusOK, concept)
}

func updateOwner(c *gin.Context) {
	var ownerConcept Concept
	if err := c.BindJSON(&ownerConcept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner data"})
		return
	}

	ownerConcept.Type = "Owner"
	ownerMu.RLock()
	ownerConcept.GUID = ownerGUID
	ownerMu.RUnlock()
	ownerConcept.Timestamp = time.Now()

	addOrUpdateConcept(c.Request.Context(), &ownerConcept)

	if err := node.Save(context.Background(), peerListPath, peerMap); err != nil {
		log.Printf("Failed to save peerMap: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Owner updated successfully", "guid": ownerConcept.GUID})
}

func getOwner(c *gin.Context) {
	conceptMu.RLock()
	ownerConcept, exists := conceptMap[ownerGUID]
	conceptMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Owner not found"})
		return
	}

	c.JSON(http.StatusOK, ownerConcept)
}

func deleteConcept(c *gin.Context) {
	guid := GUID(c.Param("guid"))

	conceptMu.Lock()
	defer conceptMu.Unlock()

	concept, exists := conceptMap[guid]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Concept not found"})
		return
	}

	if err := node.Remove(context.Background(), concept.GetCID()); err != nil {
		log.Printf("Failed to remove concept: %v", err)
	}
	delete(conceptMap, guid)
	delete(GUID2CID, guid)
	if err := node.Save(context.Background(), GUID2CIDPath, GUID2CID); err != nil {
		log.Printf("Failed to save concept list: %v", err)
	}

	c.Status(http.StatusNoContent)
}

func queryConcepts(c *gin.Context) {
	filter := ConceptFilter{
		CID:         c.Query("cid"),
		GUID:        GUID(c.Query("guid")),
		Name:        c.Query("name"),
		Description: c.Query("description"),
		Type:        c.Query("type"),
	}

	if timestamp := c.Query("timestamp"); timestamp != "" {
		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
			return
		}
		filter.TimestampAfter = &t
	}

	concepts := filterConcepts(filter)
	c.JSON(http.StatusOK, concepts)
}
