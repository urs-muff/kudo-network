package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	shell "github.com/ipfs/go-ipfs-api"
)

const (
	pubsubTopic       = "concept-list"
	publishInterval   = 1 * time.Minute
	peerCheckInterval = 5 * time.Minute
)

// IPFSShell implements the Node_i interface using go-ipfs-api
type IPFSShell struct {
	sh *shell.Shell
}

func NewIPFSShell(url string) *IPFSShell {
	return &IPFSShell{sh: shell.NewShell(url)}
}

func (i *IPFSShell) Add(ctx context.Context, content io.Reader) (CID, error) {
	cid, err := i.sh.Add(content)
	return CID(cid), err
}

func (i *IPFSShell) Get(ctx context.Context, cid CID) (io.ReadCloser, error) {
	return i.sh.Cat(string(cid))
}

func (i *IPFSShell) Remove(ctx context.Context, cid CID) error {
	// Note: IPFS doesn't have a direct "remove" function. This is a placeholder.
	return fmt.Errorf("remove operation not supported")
}

func (i *IPFSShell) List(ctx context.Context) ([]CID, error) {
	// This is a placeholder. You might need to implement this using IPFS pinning or a custom index.
	return nil, fmt.Errorf("list operation not implemented")
}

func (i *IPFSShell) Publish(ctx context.Context, topic string, data []byte) error {
	return i.sh.PubSubPublish(topic, string(data))
}

func (i *IPFSShell) Subscribe(ctx context.Context, topic string) (<-chan []byte, error) {
	sub, err := i.sh.PubSubSubscribe(topic)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)
	go func() {
		defer close(ch)
		for {
			msg, err := sub.Next()
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("Error receiving message: %v", err)
				continue
			}
			ch <- msg.Data
		}
	}()

	return ch, nil
}

func (i *IPFSShell) Connect(ctx context.Context, peerID PeerID) error {
	return i.sh.SwarmConnect(ctx, string(peerID))
}

func (i *IPFSShell) Disconnect(ctx context.Context, peerID PeerID) error {
	// Note: go-ipfs-api doesn't provide a direct disconnect method. This is a placeholder.
	return fmt.Errorf("disconnect operation not supported")
}

func (i *IPFSShell) ListPeers(ctx context.Context) ([]Peer_i, error) {
	swarmPeers, err := i.sh.SwarmPeers(ctx)
	if err != nil {
		return nil, err
	}

	peers := make([]Peer_i, len(swarmPeers.Peers))
	for j, p := range swarmPeers.Peers {
		peers[j] = &ConcretePeer{
			ID:        PeerID(p.Peer),
			Timestamp: time.Now(),
		}
	}

	return peers, nil
}

func (i *IPFSShell) Bootstrap(ctx context.Context) error {
	bootstrapNodes := []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	}

	for _, addr := range bootstrapNodes {
		if err := i.sh.SwarmConnect(ctx, addr); err != nil {
			log.Printf("Failed to connect to bootstrap node %s: %v", addr, err)
		} else {
			log.Printf("Connected to bootstrap node: %s", addr)
		}
	}

	return nil
}

func (i *IPFSShell) ID(ctx context.Context) (PeerID, error) {
	info, err := i.sh.ID()
	if err != nil {
		return "", err
	}
	return PeerID(info.ID), nil
}

var (
	ipfs        Node_i
	conceptList []Concept_i
	conceptMu   sync.RWMutex
	peerList    map[PeerID]Peer_i
	peerListMu  sync.RWMutex
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for this example
		},
	}
)

func main() {
	ipfs = NewIPFSShell("localhost:5001")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Bootstrap the IPFS node
	if err := ipfs.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap IPFS: %v", err)
	}

	// Start IPFS routines
	go publishRoutine(ctx)
	go peerDiscoveryRoutine(ctx)
	go subscribeRoutine(ctx)

	// Set up Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	r.POST("/concept", addConcept)
	r.GET("/ws", handleWebSocket)
	r.GET("/ws/peers", handlePeerWebSocket)

	// Start server
	log.Fatal(r.Run(":9090"))
}

func addConcept(c *gin.Context) {
	var newConcept struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Content     string `json:"content"`
	}

	if err := c.BindJSON(&newConcept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	cid, err := ipfs.Add(c.Request.Context(), strings.NewReader(newConcept.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add content to IPFS"})
		return
	}

	concept := &ConcreteConcept{
		Guid:        GUID(uuid.New().String()),
		Name:        newConcept.Name,
		Description: newConcept.Description,
		Type:        newConcept.Type,
		Cid:         cid,
		Content:     newConcept.Content,
		Timestamp:   time.Now(),
	}

	addNewConcept(concept)

	c.JSON(http.StatusOK, gin.H{
		"guid": concept.Guid,
		"cid":  string(concept.Cid),
	})
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Send initial concept list
	sendConceptList(conn)

	// Start goroutine to continuously send updates
	go func() {
		for {
			time.Sleep(5 * time.Second)
			sendConceptList(conn)
		}
	}()

	// Keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
	}
}

func handlePeerWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade peer connection: %v", err)
		return
	}
	defer conn.Close()

	// Send initial peer list
	sendPeerList(conn)

	// Start goroutine to continuously send updates
	go func() {
		for {
			time.Sleep(5 * time.Second)
			sendPeerList(conn)
		}
	}()

	// Keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Peer WebSocket read error: %v", err)
			break
		}
	}
}

func sendConceptList(conn *websocket.Conn) {
	conceptMu.RLock()
	defer conceptMu.RUnlock()

	for _, concept := range conceptList {
		err := conn.WriteJSON(concept)
		if err != nil {
			log.Printf("Failed to send concept: %v", err)
		}
	}
}

func sendPeerList(conn *websocket.Conn) {
	peerListMu.RLock()
	defer peerListMu.RUnlock()

	err := conn.WriteJSON(peerList)
	if err != nil {
		log.Printf("Failed to send peer list: %v", err)
	}
}

func publishRoutine(ctx context.Context) {
	ticker := time.NewTicker(publishInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			publishConceptList(ctx)
		}
	}
}

func publishConceptList(ctx context.Context) {
	conceptMu.RLock()
	data, err := json.Marshal(conceptList)
	conceptMu.RUnlock()

	if err != nil {
		log.Printf("Error marshaling concept list: %v", err)
		return
	}

	err = ipfs.Publish(ctx, pubsubTopic, data)
	if err != nil {
		log.Printf("Error publishing concept list: %v", err)
	} else {
		log.Println("Published concept list")
	}
}

func peerDiscoveryRoutine(ctx context.Context) {
	ticker := time.NewTicker(peerCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			discoverPeers(ctx)
		}
	}
}

func discoverPeers(ctx context.Context) {
	peers, err := ipfs.ListPeers(ctx)
	if err != nil {
		log.Printf("Error discovering peers: %v", err)
		return
	}

	peerListMu.Lock()
	defer peerListMu.Unlock()

	for _, peer := range peers {
		peerID := peer.GetID()
		if _, exists := peerList[peerID]; !exists {
			peerList[peerID] = peer
		}
	}

	log.Printf("Discovered %d peers", len(peerList))
}

func subscribeRoutine(ctx context.Context) {
	ch, err := ipfs.Subscribe(ctx, pubsubTopic)
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	log.Printf("Subscribed to topic: %s", pubsubTopic)

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			var receivedConcepts []ConcreteConcept
			err := json.Unmarshal(data, &receivedConcepts)
			if err != nil {
				log.Printf("Error unmarshaling received concepts: %v", err)
				continue
			}

			updateConceptList(receivedConcepts)
			log.Printf("Received concept list with %d concepts", len(receivedConcepts))
		}
	}
}

func updateConceptList(newConcepts []ConcreteConcept) {
	conceptMu.Lock()
	defer conceptMu.Unlock()

	// This is a simple update mechanism. In a real-world scenario, you might want to
	// implement a more sophisticated merging strategy.
	for _, newConcept := range newConcepts {
		found := false
		for i, existingConcept := range conceptList {
			if existingConcept.GetGUID() == newConcept.GetGUID() {
				conceptList[i] = &newConcept
				found = true
				break
			}
		}
		if !found {
			conceptList = append(conceptList, &newConcept)
		}
	}
}

func addNewConcept(concept Concept_i) {
	conceptMu.Lock()
	conceptList = append(conceptList, concept)
	conceptMu.Unlock()

	log.Printf("Added new concept: %s", concept.GetGUID())

	// Trigger immediate publish
	go publishConceptList(context.Background())
}
