package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GUID2CIDPath  = "/ccn/GUID-CID.json"
	peerListPath  = "/ccn/peer-list.json"
	ownerGUIDPath = "/ccn/owner-guid.json"
)

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

func (rm *RelationshipMap) UnmarshalJSON(data []byte) error {
	var rawMap map[PeerID]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	*rm = make(RelationshipMap)
	for guid, raw := range rawMap {
		var r Relationship
		if err := json.Unmarshal(raw, &r); err != nil {
			return err
		}
		(*rm)[GUID(guid)] = &r
	}
	return nil
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

func addOrUpdateConcept(ctx context.Context, concept *Concept) error {
	conceptMu.Lock()
	defer conceptMu.Unlock()

	if err := concept.Update(context.Background()); err != nil {
		log.Printf("Failed to update concept: %v", err)
		return err
	}
	conceptMap[concept.GetGUID()] = concept
	GUID2CID[concept.GetGUID()] = concept.GetCID()
	log.Printf("Added/Updated concept: GUID=%s, Name=%s, CID=%s\n", concept.GetGUID(), concept.GetName(), concept.GetCID())

	if err := node.Save(ctx, GUID2CIDPath, GUID2CID); err != nil {
		log.Printf("Failed to save concept list: %v", err)
		return err
	}
	return nil
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

func handleReceivedMessage(data []byte) {
	var message PeerMessage
	if err := json.Unmarshal(data, &message); err != nil {
		log.Printf("Error unmarshaling received message: %v", err)
		return
	}

	log.Printf("Received message from peer: %s", message.PeerID)

	// Add or update the sender in the peer list
	addOrUpdatePeer(message.PeerID, message.OwnerGUID)

	// Update local relationships with received ones
	for id, relationship := range message.Relationships {
		if _, exists := relationshipMap[id]; !exists {
			relationshipMap[id] = relationship
		}
	}
	saveRelationships(context.Background())

	// Update the CIDs for this peer
	updatePeerCIDs(message.PeerID, message.CIDs)
}

func addNewConcept(concept *Concept) {
	conceptMu.Lock()
	conceptMap[concept.GetGUID()] = concept
	GUID2CID[concept.GetGUID()] = concept.GetCID()
	conceptMu.Unlock()

	peerMap[peerID].AddCID(concept.GetCID())

	log.Printf("Added new concept: GUID=%s, CID=%s, Name=%s", concept.GetGUID(), concept.GetCID(), concept.GetName())

	go publishPeerMessage(context.Background())
}

// Modify the Interact method of Relationship
func (r *Relationship) Interact(interactionType GUID) {
	r.Interactions++
	r.Depth = int(math.Log2(float64(r.Interactions))) + 1
	r.LastInteraction = time.Now()

	// Get the interaction type concept
	conceptMu.RLock()
	interactionConcept, ok := conceptMap[interactionType]
	conceptMu.RUnlock()

	if !ok {
		log.Printf("Interaction type %s not found", interactionType)
		return
	}

	// Apply effects based on the interaction type
	switch interactionConcept.Name {
	case "Music":
		r.FrequencySpec = append(r.FrequencySpec, 440.0) // Add A4 note
		r.Amplitude *= 1.05
	case "Meditation":
		r.EnergyFlow *= 1.1
		r.Volume *= 0.95
	case "FlowState":
		r.EnergyFlow *= 1.2
		r.Amplitude *= 1.1
		r.Volume *= 1.05
	default:
		r.EnergyFlow *= 1.05
	}

	r.Timestamp = time.Now()
}
