package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func createCoreRelationship(ctx context.Context, sourceGUID, relationshipTypeGUID, targetGUID GUID) error {
	relationshipID := generateGUID(fmt.Sprintf("%s-%s-%s", sourceGUID, relationshipTypeGUID, targetGUID))

	relationship := &Relationship{
		ID:              relationshipID,
		SourceID:        sourceGUID,
		TargetID:        targetGUID,
		Type:            relationshipTypeGUID,
		EnergyFlow:      1.0,
		FrequencySpec:   []float64{1.0},
		Amplitude:       1.0,
		Volume:          1.0,
		Depth:           1,
		Interactions:    0,
		LastInteraction: time.Now(),
		Timestamp:       time.Now(),
	}

	relationshipMap[relationshipID] = relationship

	// Update the relationships for the source and target concepts
	sourceConcept, exists := conceptMap[sourceGUID]
	if !exists {
		return fmt.Errorf("source concept with GUID %s not found", sourceGUID)
	}
	sourceConcept.Relationships = append(sourceConcept.Relationships, relationshipID)

	targetConcept, exists := conceptMap[targetGUID]
	if !exists {
		return fmt.Errorf("target concept with GUID %s not found", targetGUID)
	}
	targetConcept.Relationships = append(targetConcept.Relationships, relationshipID)

	return nil
}

// Call this function in your main.go after initializing the IPFS node
func InitializeSystem(ctx context.Context) error {
	/*
		if err := BootstrapCoreConceptsAndRelationships(ctx); err != nil {
			log.Printf("Error during boostrapping core concepts: %\n", err)
			return err
		}
		if err := BootstrapExpandedCoreConceptsAndRelationships(ctx); err != nil {
			return err
		}
	*/
	log.Println("Bootstrapping concepts and relationships...")

	if err := BootstrapFromStructure(ctx, "data/concepts_structure.yaml"); err != nil {
		log.Printf("Error during bootstrapping concepts: %v\n", err)
		return err
	}

	log.Println("Concepts and relationships bootstrapped successfully")
	return nil
}
