package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pubsubTopic       = "concept-list"
	publishInterval   = 1 * time.Minute
	peerCheckInterval = 5 * time.Minute
	GUID2CIDPath      = "/ccn/GUID-CID.json"
	peerListPath      = "/ccn/peer-list.json"
	ownerGUIDPath     = "/ccn/owner-guid.json"
)

type PeerMessage struct {
	PeerID    PeerID `json:"peerId"`
	OwnerGUID GUID   `json:"ownerGuid"`
	CIDs      []CID  `json:"cids"`
}

var (
	node Node_i

	conceptMap map[GUID]Concept_i
	GUID2CID   map[GUID]CID
	conceptMu  sync.RWMutex

	peerMap   PeerMap
	peerMapMu sync.RWMutex
	peerID    PeerID

	ownerGUID GUID
	ownerMu   sync.RWMutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for this example
		},
	}
)

type PeerMap map[PeerID]Peer_i

func (pm *PeerMap) UnmarshalJSON(data []byte) error {
	var rawMap map[PeerID]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	*pm = make(PeerMap)
	for peerID, raw := range rawMap {
		var p Peer
		if err := json.Unmarshal(raw, &p); err != nil {
			return err
		}
		(*pm)[peerID] = &p
	}
	return nil
}

func main() {
	node = NewIPFSShell("localhost:5001")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initializeLists(ctx)

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
	GUID2CID = make(map[GUID]CID)
	peerMap = make(PeerMap)

	if err := node.Bootstrap(ctx); err != nil {
		log.Fatalf("Failed to bootstrap IPFS: %v", err)
	}

	var err error
	peerID, err = node.ID(ctx)
	if err != nil {
		log.Fatalf("Failed to get peer ID: %v", err)
	}

	if err := node.Load(ctx, GUID2CIDPath, &GUID2CID); err != nil {
		log.Printf("Failed to load concept list: %v\n", err)
	}
	if err := node.Load(ctx, peerListPath, &peerMap); err != nil {
		log.Printf("Failed to load peer list: %v\n", err)
	}
	for id, peer := range peerMap {
		if peer.GetOwnerGUID() == "" {
			delete(peerMap, id)
		}
	}
	peerMap[peerID] = &Peer{
		ID:        peerID,
		Timestamp: time.Now(),
		CIDs:      make(map[CID]bool),
	}
	loadOrCreateOwner(ctx)
	peerMap[peerID].(*Peer).OwnerGUID = ownerGUID
	for _, cid := range peerMap[peerID].GetCIDs() {
		conceptReader, err := node.Get(context.Background(), cid)
		if err != nil {
			log.Fatalf("Unable to get Concept: %s: %v", cid, err)
		}
		var c Concept
		err = json.NewDecoder(conceptReader).Decode(&c)
		if err != nil {
			log.Fatalf("Unable to parse Concept: %s: %v", cid, err)
		}
		c.CID = cid
		conceptMap[c.GUID] = &c
		GUID2CID[c.GUID] = cid
	}
}

func loadOrCreateOwner(ctx context.Context) {
	var guid GUID
	err := node.Load(ctx, ownerGUIDPath, &guid)
	if err != nil {
		log.Printf("Failed to load owner GUID from IPFS: %v", err)
		log.Println("Generating new owner GUID...")
		guid = GUID(uuid.New().String())
		if err := node.Save(ctx, ownerGUIDPath, guid); err != nil {
			log.Fatalf("Failed to save new owner GUID: %v", err)
		}
	}

	ownerMu.Lock()
	ownerGUID = guid
	ownerMu.Unlock()

	log.Printf("Owner GUID: %s", ownerGUID)
	cid, ok := GUID2CID[ownerGUID]
	if !ok {
		ownerConcept := Concept{
			GUID:        guid,
			Name:        "Owner",
			Description: "Owner",
			Type:        "Owner",
			Timestamp:   time.Now(),
		}
		addOrUpdateConcept(context.Background(), &ownerConcept)
		cid = ownerConcept.CID
	}
	peerMap[peerID].AddCID(cid)
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

func (c *Concept) Update(ctx context.Context) error {
	conceptJSON, _ := json.Marshal(c)
	cid, err := node.Add(ctx, strings.NewReader(string(conceptJSON)))
	if err != nil {
		return err
	}
	c.CID = cid
	return nil
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

func addOrUpdateConcept(ctx context.Context, concept Concept_i) {
	conceptMu.Lock()
	defer conceptMu.Unlock()

	if err := concept.Update(context.Background()); err != nil {
		log.Printf("Failed to update concept: %v", err)
	}
	conceptMap[concept.GetGUID()] = concept
	GUID2CID[concept.GetGUID()] = concept.GetCID()
	log.Printf("Added/Updated concept: GUID=%s, Name=%s, CID=%s\n", concept.GetGUID(), concept.GetName(), concept.GetCID())

	if err := node.Save(ctx, GUID2CIDPath, GUID2CID); err != nil {
		log.Printf("Failed to save concept list: %v", err)
	}
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
		GUID:        GUID(uuid.New().String()),
		Name:        newConcept.Name,
		Description: newConcept.Description,
		Type:        newConcept.Type,
		Timestamp:   time.Now(),
	}
	concept.Update(c.Request.Context())

	addNewConcept(concept)
	c.JSON(http.StatusOK, gin.H{
		"guid": concept.GUID,
		"cid":  string(concept.CID),
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
	ch, err := node.Subscribe(ctx, pubsubTopic)
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

	if err := node.Publish(ctx, pubsubTopic, data); err != nil {
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

	if err := node.Save(context.Background(), peerListPath, peerMap); err != nil {
		log.Printf("Failed to save peerMap: %v", err)
	}
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
	peers, err := node.ListPeers(ctx)
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

	if err := node.Save(context.Background(), peerListPath, peerMap); err != nil {
		log.Printf("Failed to save peerMap: %v", err)
	}
}

func addNewConcept(concept Concept_i) {
	conceptMu.Lock()
	conceptMap[concept.GetGUID()] = concept
	GUID2CID[concept.GetGUID()] = concept.GetCID()
	conceptMu.Unlock()

	peerMap[peerID].AddCID(concept.GetCID())

	log.Printf("Added new concept: GUID=%s, CID=%s, Name=%s", concept.GetGUID(), concept.GetCID(), concept.GetName())

	go publishPeerMessage(context.Background())
}
