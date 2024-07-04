package main

import (
	"context"
	"log"
	"time"
)

// CoreConcept represents a fundamental concept in the system
type CoreConcept struct {
	GUID        GUID
	Name        string
	Description string
	Type        string
}

// CoreRelationship represents a fundamental relationship type in the system
type CoreRelationship struct {
	GUID        GUID
	Name        string
	Description string
}

var (
	// Core Concept Types
	ConceptType = CoreConcept{
		GUID:        "00000000-0000-0000-0000-000000000001",
		Name:        "Concept",
		Description: "A fundamental unit of knowledge or idea",
		Type:        "ConceptType",
	}
	RelationshipType = CoreConcept{
		GUID:        "00000000-0000-0000-0000-000000000002",
		Name:        "Relationship",
		Description: "A connection or association between concepts",
		Type:        "ConceptType",
	}
	InteractionType = CoreConcept{
		GUID:        "00000000-0000-0000-0000-000000000003",
		Name:        "Interaction",
		Description: "A type of action that can occur between concepts or within relationships",
		Type:        "ConceptType",
	}

	// Core Relationship Types
	IsA = CoreRelationship{
		GUID:        "10000000-0000-0000-0000-000000000001",
		Name:        "Is A",
		Description: "Indicates that one concept is a type or instance of another",
	}
	HasA = CoreRelationship{
		GUID:        "10000000-0000-0000-0000-000000000002",
		Name:        "Has A",
		Description: "Indicates that one concept possesses or includes another",
	}
	RelatedTo = CoreRelationship{
		GUID:        "10000000-0000-0000-0000-000000000003",
		Name:        "Related To",
		Description: "Indicates a general relationship between concepts",
	}

	// Core Interaction Types
	MusicInteraction = CoreConcept{
		GUID:        "20000000-0000-0000-0000-000000000001",
		Name:        "Music",
		Description: "Interaction through musical elements",
		Type:        "InteractionType",
	}
	MeditationInteraction = CoreConcept{
		GUID:        "20000000-0000-0000-0000-000000000002",
		Name:        "Meditation",
		Description: "Interaction through meditative practices",
		Type:        "InteractionType",
	}
	FlowStateInteraction = CoreConcept{
		GUID:        "20000000-0000-0000-0000-000000000003",
		Name:        "Flow State",
		Description: "Interaction in a state of optimal experience and focus",
		Type:        "InteractionType",
	}
)

// BootstrapCoreConceptsAndRelationships initializes the system with core concepts and relationships
func BootstrapCoreConceptsAndRelationships(ctx context.Context) error {
	log.Println("Bootstrapping core concepts and relationships...")

	// Add core concept types
	if err := addCoreConcept(ctx, ConceptType); err != nil {
		return err
	}
	if err := addCoreConcept(ctx, RelationshipType); err != nil {
		return err
	}
	if err := addCoreConcept(ctx, InteractionType); err != nil {
		return err
	}

	// Add core relationship types
	if err := addCoreRelationship(ctx, IsA); err != nil {
		return err
	}
	if err := addCoreRelationship(ctx, HasA); err != nil {
		return err
	}
	if err := addCoreRelationship(ctx, RelatedTo); err != nil {
		return err
	}

	// Add core interaction types
	if err := addCoreConcept(ctx, MusicInteraction); err != nil {
		return err
	}
	if err := addCoreConcept(ctx, MeditationInteraction); err != nil {
		return err
	}
	if err := addCoreConcept(ctx, FlowStateInteraction); err != nil {
		return err
	}

	// Create relationships between core concepts
	if err := createCoreRelationship(ctx, RelationshipType.GUID, IsA.GUID, ConceptType.GUID); err != nil {
		return err
	}
	if err := createCoreRelationship(ctx, InteractionType.GUID, IsA.GUID, ConceptType.GUID); err != nil {
		return err
	}

	log.Println("Core concepts and relationships bootstrapped successfully")
	return nil
}

func addCoreConcept(ctx context.Context, cc CoreConcept) error {
	concept := &Concept{
		GUID:        GUID(cc.GUID),
		Name:        cc.Name,
		Description: cc.Description,
		Type:        cc.Type,
		Timestamp:   time.Now(),
	}
	return addOrUpdateConcept(ctx, concept)
}

func addCoreRelationship(ctx context.Context, cr CoreRelationship) error {
	concept := &Concept{
		GUID:        GUID(cr.GUID),
		Name:        cr.Name,
		Description: cr.Description,
		Type:        "RelationshipType",
		Timestamp:   time.Now(),
	}
	return addOrUpdateConcept(ctx, concept)
}

func createCoreRelationship(ctx context.Context, sourceID, relationshipTypeID, targetID GUID) error {
	relationship := CreateRelationship(sourceID, targetID, relationshipTypeID)
	relationshipMap[relationship.ID] = relationship
	return saveRelationships(ctx)
}

// Call this function in your main.go after initializing the IPFS node
func InitializeSystem(ctx context.Context) error {
	if err := BootstrapCoreConceptsAndRelationships(ctx); err != nil {
		log.Printf("Error during boostrapping core concepts: %\n", err)
		return err
	}
	// Add any other initialization steps here
	return nil
}
