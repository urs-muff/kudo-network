package main

import (
	"bytes"
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
	conceptListCID    = "concept-list-cid"
	peerListCID       = "peer-list-cid"
	ownerGUIDCID      = "owner-guid-cid"
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
		peers[j] = &Peer{
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

type PeerMessage struct {
	PeerID    PeerID `json:"peerId"`
	OwnerGUID GUID   `json:"ownerGuid"`
	CIDs      []CID  `json:"cids"`
}

var (
	ipfs       Node_i
	conceptMap map[GUID]Concept_i
	conceptMu  sync.RWMutex
	peerMap    map[PeerID]Peer_i
	peerMapMu  sync.RWMutex
	ownerGUID  GUID
	ownerMu    sync.RWMutex
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for this example
		},
	}
)

func main() {
	ipfs = NewIPFSShell("localhost:5001")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initializeLists(ctx)
	loadOwnerGUID(ctx)

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

func initializeLists(ctx context.Context) {
	conceptMap = make(map[GUID]Concept_i)
	peerMap = make(map[PeerID]Peer_i)

	if err := ipfs.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap IPFS: %v", err)
	}

	loadFromIPFS(ctx, CID(conceptListCID), &conceptMap)
	loadFromIPFS(ctx, CID(peerListCID), &peerMap)
}

func loadOwnerGUID(ctx context.Context) {
	data, err := ipfs.Get(ctx, CID(ownerGUIDCID))
	if err != nil {
		log.Printf("Failed to load owner GUID from IPFS: %v", err)
		log.Println("Generating new owner GUID...")
		ownerMu.Lock()
		ownerGUID = GUID(uuid.New().String())
		ownerMu.Unlock()
		if err := saveOwnerGUID(ctx); err != nil {
			log.Printf("Failed to save new owner GUID: %v", err)
		}
		return
	}
	defer data.Close()

	var guid GUID
	if err := json.NewDecoder(data).Decode(&guid); err != nil {
		log.Printf("Failed to decode owner GUID: %v", err)
		return
	}

	ownerMu.Lock()
	ownerGUID = guid
	ownerMu.Unlock()

	log.Printf("Loaded owner GUID from IPFS: %s", ownerGUID)
}

func saveOwnerGUID(ctx context.Context) error {
	ownerMu.RLock()
	guid := ownerGUID
	ownerMu.RUnlock()

	data, err := json.Marshal(guid)
	if err != nil {
		return fmt.Errorf("failed to marshal owner GUID: %v", err)
	}

	cid, err := ipfs.Add(ctx, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to save owner GUID to IPFS: %v", err)
	}

	log.Printf("Saved owner GUID to IPFS with CID: %s", cid)
	return nil
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

func loadFromIPFS(ctx context.Context, cid CID, target interface{}) {
	data, err := ipfs.Get(ctx, cid)
	if err != nil {
		log.Printf("Failed to load data from IPFS (CID: %s): %v", cid, err)
		return
	}
	defer data.Close()

	if err := json.NewDecoder(data).Decode(target); err != nil {
		log.Printf("Failed to decode data from IPFS (CID: %s): %v", cid, err)
		return
	}

	log.Printf("Loaded data from IPFS (CID: %s)", cid)
}

func saveToIPFS(ctx context.Context, data interface{}) (CID, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	cid, err := ipfs.Add(ctx, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to save data to IPFS: %v", err)
	}

	log.Printf("Saved data to IPFS with CID: %s", cid)
	return cid, nil
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

func updateOwner(c *gin.Context) {
	var ownerConcept Concept
	if err := c.BindJSON(&ownerConcept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner data"})
		return
	}

	ownerConcept.Type = "Owner"
	ownerMu.RLock()
	ownerConcept.Guid = ownerGUID
	ownerMu.RUnlock()
	ownerConcept.Timestamp = time.Now()

	peerID, err := ipfs.ID(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get peer ID"})
		return
	}

	addOrUpdateConcept(&ownerConcept)

	peerMapMu.Lock()
	peerMap[peerID] = &Peer{
		ID:        peerID,
		OwnerGUID: ownerConcept.Guid,
		Timestamp: time.Now(),
	}
	peerMapMu.Unlock()

	saveToIPFS(context.Background(), peerMap)

	c.JSON(http.StatusOK, gin.H{"message": "Owner updated successfully", "guid": ownerConcept.Guid})
}

func addOrUpdateConcept(concept Concept_i) {
	conceptMu.Lock()
	defer conceptMu.Unlock()

	conceptMap[concept.GetGUID()] = concept
	log.Printf("Added/Updated concept: GUID=%s, Name=%s", concept.GetGUID(), concept.GetName())

	saveToIPFS(context.Background(), conceptMap)
}

func getOwner(c *gin.Context) {
	ownerMu.RLock()
	ownerGUID := ownerGUID
	ownerMu.RUnlock()

	conceptMu.RLock()
	ownerConcept, exists := conceptMap[ownerGUID]
	conceptMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Owner not found"})
		return
	}

	c.JSON(http.StatusOK, ownerConcept)
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

	concept := &Concept{
		Guid:        GUID(uuid.New().String()),
		Name:        newConcept.Name,
		Description: newConcept.Description,
		Type:        newConcept.Type,
		Cid:         cid,
		Content:     newConcept.Content,
		Timestamp:   time.Now(),
	}

	addNewConcept(concept)

	log.Printf("Added new concept: GUID=%s, Name=%s, CID=%s", concept.Guid, concept.Name, concept.Cid)

	c.JSON(http.StatusOK, gin.H{
		"guid": concept.Guid,
		"cid":  string(concept.Cid),
	})
}

func handleWebSocket(c *gin.Context) {
	handleWebSocketConnection(c, sendConceptList)
}

func handlePeerWebSocket(c *gin.Context) {
	handleWebSocketConnection(c, sendPeerList)
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

func deleteConcept(c *gin.Context) {
	guid := GUID(c.Param("guid"))

	conceptMu.Lock()
	defer conceptMu.Unlock()

	if _, exists := conceptMap[guid]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Concept not found"})
		return
	}

	delete(conceptMap, guid)
	saveToIPFS(context.Background(), conceptMap)

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

func listPeers(c *gin.Context) {
	peerMapMu.RLock()
	defer peerMapMu.RUnlock()

	filteredPeerMap := make(map[PeerID]Peer_i)
	for peerID, peer := range peerMap {
		if peer.GetOwnerGUID() != "" {
			filteredPeerMap[peerID] = peer
		}
	}

	c.JSON(http.StatusOK, filteredPeerMap)
}

type ConceptFilter struct {
	CID            string
	GUID           GUID
	Name           string
	Description    string
	Type           string
	TimestampAfter *time.Time
}

func filterConcepts(filter ConceptFilter) []Concept_i {
	conceptMu.RLock()
	defer conceptMu.RUnlock()

	if isEmptyFilter(filter) {
		concepts := make([]Concept_i, 0, len(conceptMap))
		for _, concept := range conceptMap {
			concepts = append(concepts, concept)
		}
		return concepts
	}

	var filteredConcepts []Concept_i
	for _, concept := range conceptMap {
		if matchesConcept(concept, filter) {
			filteredConcepts = append(filteredConcepts, concept)
		}
	}
	return filteredConcepts
}

func isEmptyFilter(filter ConceptFilter) bool {
	return filter.CID == "" && filter.GUID == "" && filter.Name == "" &&
		filter.Description == "" && filter.Type == "" && filter.TimestampAfter == nil
}

func matchesConcept(concept Concept_i, filter ConceptFilter) bool {
	if filter.CID != "" && string(concept.GetCID()) != filter.CID {
		return false
	}
	if filter.GUID != "" && concept.GetGUID() != filter.GUID {
		return false
	}
	if filter.Name != "" && !strings.Contains(strings.ToLower(concept.GetName()), strings.ToLower(filter.Name)) {
		return false
	}
	if filter.Description != "" && !strings.Contains(strings.ToLower(concept.GetDescription()), strings.ToLower(filter.Description)) {
		return false
	}
	if filter.Type != "" && concept.GetType() != filter.Type {
		return false
	}
	if filter.TimestampAfter != nil && !concept.GetTimestamp().After(*filter.TimestampAfter) {
		return false
	}
	return true
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
		case msg := <-ch:
			handleReceivedMessage(msg)
		}
	}
}

func handleWebSocketConnection(c *gin.Context, sendFunc func(*websocket.Conn)) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection established")

	sendFunc(conn)

	go periodicSend(conn, sendFunc)

	keepAlive(conn)

	log.Printf("WebSocket connection closed")
}

func periodicSend(conn *websocket.Conn, sendFunc func(*websocket.Conn)) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		sendFunc(conn)
	}
}

func keepAlive(conn *websocket.Conn) {
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
	}
}

func sendConceptList(conn *websocket.Conn) {
	sendJSONList(conn, conceptMap, &conceptMu, "concept map")
}

func sendPeerList(conn *websocket.Conn) {
	peerMapMu.RLock()
	defer peerMapMu.RUnlock()

	filteredPeerMap := make(map[PeerID]Peer_i)
	for peerID, peer := range peerMap {
		if peer.GetOwnerGUID() != "" {
			filteredPeerMap[peerID] = peer
		}
	}

	if err := conn.WriteJSON(filteredPeerMap); err != nil {
		log.Printf("Failed to send peer list: %v", err)
	} else {
		log.Printf("Sent filtered peer list with %d peers", len(filteredPeerMap))
	}
}

func sendJSONList(conn *websocket.Conn, list interface{}, mu sync.Locker, itemType string) {
	mu.Lock()
	defer mu.Unlock()

	if err := conn.WriteJSON(list); err != nil {
		log.Printf("Failed to send %s: %v", itemType, err)
	} else {
		log.Printf("Sent %s", itemType)
	}
}

func publishPeerMessage(ctx context.Context) {
	peerID, err := ipfs.ID(ctx)
	if err != nil {
		log.Printf("Error getting peer ID: %v", err)
		return
	}

	peerMapMu.RLock()
	peer, exists := peerMap[peerID]
	peerMapMu.RUnlock()

	if !exists {
		log.Printf("Peer information not set for this peer")
		return
	}

	conceptMu.RLock()
	cids := make([]CID, 0, len(conceptMap))
	for _, concept := range conceptMap {
		cids = append(cids, concept.GetCID())
	}
	conceptMu.RUnlock()

	message := PeerMessage{
		PeerID:    peerID,
		OwnerGUID: peer.GetOwnerGUID(),
		CIDs:      cids,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling peer message: %v", err)
		return
	}

	if err := ipfs.Publish(ctx, pubsubTopic, data); err != nil {
		log.Printf("Error publishing peer message: %v", err)
	} else {
		log.Printf("Published peer message with %d CIDs", len(cids))
	}
}

func handleReceivedMessage(data []byte) {
	var message PeerMessage
	if err := json.Unmarshal(data, &message); err != nil {
		log.Printf("Error unmarshaling received message: %v", err)
		return
	}

	log.Printf("Received message from peer: %s", message.PeerID)

	// Add or update the sender in the peer list
	addOrUpdatePeer(message.PeerID, message.OwnerGUID)

	// Update the CIDs for this peer
	updatePeerCIDs(message.PeerID, message.CIDs)
}

func addOrUpdatePeer(peerID PeerID, ownerGUID GUID) {
	peerMapMu.Lock()
	defer peerMapMu.Unlock()

	peerMap[peerID] = &Peer{
		ID:        peerID,
		OwnerGUID: ownerGUID,
		Timestamp: time.Now(),
	}
	log.Printf("Updated peer: %s", peerID)

	saveToIPFS(context.Background(), peerMap)
}

func updatePeerCIDs(peerID PeerID, cids []CID) {
	conceptMu.Lock()
	defer conceptMu.Unlock()

	for _, cid := range cids {
		found := false
		for _, concept := range conceptMap {
			if concept.GetCID() == cid {
				found = true
				break
			}
		}
		if !found {
			// If the concept is not in our list, we might want to fetch it
			// This is left as an exercise, as it depends on how you want to handle this case
			log.Printf("Found new CID from peer %s: %s", peerID, cid)
		}
	}
}

func discoverPeers(ctx context.Context) {
	peers, err := ipfs.ListPeers(ctx)
	if err != nil {
		log.Printf("Error discovering peers: %v", err)
		return
	}

	peerMapMu.Lock()
	defer peerMapMu.Unlock()

	for _, peer := range peers {
		peerID := peer.GetID()
		if _, exists := peerMap[peerID]; !exists {
			peerMap[peerID] = peer
			log.Printf("Discovered new peer: %s", peerID)
		}
	}

	log.Printf("Discovered %d peers", len(peerMap))

	saveToIPFS(ctx, peerMap)
}

func addNewConcept(concept Concept_i) {
	conceptMu.Lock()
	conceptMap[concept.GetGUID()] = concept
	conceptMu.Unlock()

	log.Printf("Added new concept: GUID=%s, Name=%s", concept.GetGUID(), concept.GetName())

	go publishPeerMessage(context.Background())
}
