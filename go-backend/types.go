package main

import (
	"sync"
	"time"
)

type PeerMessage struct {
	PeerID        PeerID          `json:"peerId"`
	OwnerGUID     GUID            `json:"ownerGuid"`
	CIDs          []CID           `json:"cids"`
	Relationships RelationshipMap `json:"relationships"`
}

type PeerMap map[PeerID]Peer_i

type ConceptFilter struct {
	CID            string
	GUID           GUID
	Name           string
	Description    string
	Type           string
	TimestampAfter *time.Time
}

var (
	conceptMap map[GUID]*Concept
	GUID2CID   map[GUID]CID
	conceptMu  sync.RWMutex

	relationshipMap RelationshipMap

	peerMap   PeerMap
	peerMapMu sync.RWMutex
	peerID    PeerID

	ownerGUID GUID
	ownerMu   sync.RWMutex
)
