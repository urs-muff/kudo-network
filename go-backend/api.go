package main

import (
	"context"
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
	GetGUID() GUID
	GetName() string
	GetDescription() string
	GetType() string
	GetCID() CID
	GetContent() string
	GetTimestamp() time.Time
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

	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, data []byte) error

	// Subscribe subscribes to a topic and returns a channel for receiving messages
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)

	// Connect connects to a peer
	Connect(ctx context.Context, peerID PeerID) error

	// Disconnect disconnects from a peer
	Disconnect(ctx context.Context, peerID PeerID) error

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
	Guid        GUID
	Name        string
	Description string
	Type        string
	Cid         CID
	Content     string
	Timestamp   time.Time
}

func (c Concept) GetGUID() GUID           { return c.Guid }
func (c Concept) GetName() string         { return c.Name }
func (c Concept) GetDescription() string  { return c.Description }
func (c Concept) GetType() string         { return c.Type }
func (c Concept) GetCID() CID             { return c.Cid }
func (c Concept) GetContent() string      { return c.Content }
func (c Concept) GetTimestamp() time.Time { return c.Timestamp }

// ConcretePeer implements the Peer_i interface
type Peer struct {
	ID        PeerID
	OwnerGUID GUID
	CIDs      []CID
	Timestamp time.Time
}

func (p Peer) GetID() PeerID           { return p.ID }
func (p Peer) GetOwnerGUID() GUID      { return p.OwnerGUID }
func (p Peer) GetCIDs() []CID          { return p.CIDs }
func (p Peer) GetTimestamp() time.Time { return p.Timestamp }
func (p *Peer) AddCID(cid CID)         { p.CIDs = append(p.CIDs, cid) }
func (p *Peer) RemoveCID(cid CID) {
	for i, c := range p.CIDs {
		if c == cid {
			p.CIDs = append(p.CIDs[:i], p.CIDs[i+1:]...)
			break
		}
	}
}
