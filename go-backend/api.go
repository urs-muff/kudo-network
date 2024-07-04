package main

import (
	"context"
	"encoding/json"
	"io"
	"time"
)

// CID represents a Content Identifier in IPFS
type CID string

// PeerID represents a unique identifier for a peer in the network
type PeerID string

// GUID represents a globally unique identifier
type GUID string

// Concept_i represents a concept stored in the network
type Concept_i interface {
	GetCID() CID

	GetGUID() GUID
	GetName() string
	GetDescription() string
	GetType() string
	GetTimestamp() time.Time

	Update(ctx context.Context) error
}

// Peer_i represents a peer in the network
type Peer_i interface {
	GetID() PeerID
	GetOwnerGUID() GUID
	GetCIDs() []CID
	GetTimestamp() time.Time
	AddCID(cid CID)
	RemoveCID(cid CID)
}

// Network_i defines the interface for interacting with the network
type Network_i interface {
	// Add adds content to the network and returns its CID
	Add(ctx context.Context, content io.Reader) (CID, error)

	// Get retrieves content from the network by its CID
	Get(ctx context.Context, cid CID) (io.ReadCloser, error)

	// Remove removes content from the network by its CID
	Remove(ctx context.Context, cid CID) error

	// List returns a list of all CIDs stored by this node
	List(ctx context.Context) ([]CID, error)

	// Load loads data from a given path in the network
	Load(ctx context.Context, path string, target interface{}) error

	// Save saves data to a given path in the network
	Save(ctx context.Context, path string, data interface{}) error

	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, data []byte) error

	// Subscribe subscribes to a topic and returns a channel for receiving messages
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)

	// Connect connects to a peer
	Connect(ctx context.Context, peerID PeerID) error

	// ListPeers returns a list of connected peers
	ListPeers(ctx context.Context) ([]Peer_i, error)
}

// Node_i represents a node in the network
type Node_i interface {
	Network_i

	// Bootstrap connects to bootstrap nodes
	Bootstrap(ctx context.Context) error

	// ID returns the ID of this node
	ID(ctx context.Context) (PeerID, error)
}

// Now let's define some concrete implementations of these interfaces

// Concept implements the Concept_i interface
type Concept struct {
	CID         CID `json:"-"`
	GUID        GUID
	Name        string
	Description string
	Type        string
	Timestamp   time.Time
}

func (c Concept) GetCID() CID             { return c.CID }
func (c Concept) GetGUID() GUID           { return c.GUID }
func (c Concept) GetName() string         { return c.Name }
func (c Concept) GetDescription() string  { return c.Description }
func (c Concept) GetType() string         { return c.Type }
func (c Concept) GetTimestamp() time.Time { return c.Timestamp }

// ConcretePeer implements the Peer_i interface
type Peer struct {
	ID        PeerID
	OwnerGUID GUID
	CIDs      map[CID]bool
	Timestamp time.Time
}

func (p Peer) GetID() PeerID      { return p.ID }
func (p Peer) GetOwnerGUID() GUID { return p.OwnerGUID }
func (p Peer) GetCIDs() []CID {
	ret := make([]CID, 0)
	for cid := range p.CIDs {
		ret = append(ret, cid)
	}
	return ret
}
func (p Peer) GetTimestamp() time.Time { return p.Timestamp }
func (p *Peer) AddCID(cid CID)         { p.CIDs[cid] = true }
func (p *Peer) RemoveCID(cid CID)      { delete(p.CIDs, cid) }

func (p *Peer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        PeerID
		OwnerGUID GUID
		CIDs      []CID
		Timestamp time.Time
	}{
		ID:        p.ID,
		OwnerGUID: p.OwnerGUID,
		CIDs:      p.GetCIDs(),
		Timestamp: p.Timestamp,
	})
}

func (p *Peer) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID        PeerID    `json:"id"`
		OwnerGUID GUID      `json:"ownerGuid"`
		CIDs      []CID     `json:"cids"`
		Timestamp time.Time `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	p.ID = temp.ID
	p.OwnerGUID = temp.OwnerGUID
	p.Timestamp = temp.Timestamp
	p.CIDs = make(map[CID]bool)

	for _, cid := range temp.CIDs {
		p.CIDs[cid] = true
	}

	return nil
}
