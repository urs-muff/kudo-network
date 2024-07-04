package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	pubsubTopic       = "concept-list"
	publishInterval   = 1 * time.Minute
	peerCheckInterval = 5 * time.Minute
)

var (
	node Node_i
)

func main() {
	node = NewIPFSShell("localhost:5001")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initializeLists(ctx)
	InitializeSystem(ctx)

	// Start IPFS routines
	go runPeriodicTask(ctx, publishInterval, publishPeerMessage)
	go runPeriodicTask(ctx, peerCheckInterval, discoverPeers)
	go subscribeRoutine(ctx)

	// Set up Gin router
	r := gin.Default()
	setupRoutes(r)

	// Start server
	log.Fatal(r.Run(":9090"))
}

func setupRoutes(r *gin.Engine) {
	r.Use(corsMiddleware())
	r.POST("/concept", addConcept)
	r.GET("/concept/:guid", getConcept)
	r.POST("/owner", updateOwner)
	r.GET("/owner", getOwner)
	r.DELETE("/concept/:guid", deleteConcept)
	r.GET("/concepts", queryConcepts)
	r.GET("/peers", listPeers)
	r.GET("/ws", handleWebSocket)
	r.GET("/ws/peers", handlePeerWebSocket)
	r.POST("/relationship", addRelationship)
	r.PUT("/relationship/:id/deepen", deepenRelationship)
	r.GET("/relationship/:id", getRelationship)
	r.GET("/relationship-types", getRelationshipTypes)
	r.GET("/relationship-type/:type", getRelationshipsByType)
	r.GET("/interact/:id", interactWithRelationship)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func runPeriodicTask(ctx context.Context, interval time.Duration, task func(context.Context)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			task(ctx)
		}
	}
}
