package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type ConceptStructure struct {
	Concepts      []ConceptNode      `yaml:"concepts"`
	Relationships []RelationshipNode `yaml:"relationships"`
}

type ConceptNode struct {
	Name          string             `yaml:"name"`
	Description   string             `yaml:"description"`
	Type          string             `yaml:"type"`
	Children      []ConceptNode      `yaml:"children,omitempty"`
	Relationships []RelationshipType `yaml:"relationships,omitempty"`
}

type RelationshipType struct {
	Type   string `yaml:"type"`
	Target string `yaml:"target"`
}

type RelationshipNode struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

var guidMap = make(map[string]GUID)

func generateGUID(name string) GUID {
	if guid, exists := guidMap[name]; exists {
		return guid
	}
	hash := sha256.Sum256([]byte(name))
	guid := GUID(hex.EncodeToString(hash[:16])) // Use first 16 bytes for GUID
	guidMap[name] = guid
	return guid
}

func parseConceptStructure(filename string) (*ConceptStructure, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var structure ConceptStructure
	err = yaml.Unmarshal(data, &structure)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return &structure, nil
}

func createConcepts(ctx context.Context, node ConceptNode, parentGUID GUID) (*Concept, error) {
	guid := generateGUID(node.Name)
	concept := &Concept{
		GUID:        guid,
		Name:        node.Name,
		Description: node.Description,
		Type:        node.Type,
		Timestamp:   time.Now(),
	}

	if parentGUID != "" {
		err := createCoreRelationship(ctx, parentGUID, generateGUID("Is A"), guid)
		if err != nil {
			return nil, fmt.Errorf("failed to create 'Is A' relationship for %s: %v", node.Name, err)
		}
	}

	for _, child := range node.Children {
		childConcept, err := createConcepts(ctx, child, guid)
		if err != nil {
			return nil, err
		}
		concept.Relationships = append(concept.Relationships, childConcept.GUID)
	}

	err := addOrUpdateConcept(ctx, concept)
	if err != nil {
		return nil, fmt.Errorf("failed to add or update concept %s: %v", node.Name, err)
	}

	return concept, nil
}

func createRelationships(ctx context.Context, node ConceptNode) error {
	sourceGUID := generateGUID(node.Name)
	for _, rel := range node.Relationships {
		targetGUID := generateGUID(rel.Target)
		relTypeGUID := generateGUID(rel.Type)
		err := createCoreRelationship(ctx, sourceGUID, relTypeGUID, targetGUID)
		if err != nil {
			return fmt.Errorf("failed to create relationship %s -> %s -> %s: %v", node.Name, rel.Type, rel.Target, err)
		}
	}

	for _, child := range node.Children {
		err := createRelationships(ctx, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func BootstrapFromStructure(ctx context.Context, filename string) error {
	structure, err := parseConceptStructure(filename)
	if err != nil {
		return fmt.Errorf("failed to parse concept structure: %v", err)
	}

	// Create relationship types
	for _, rel := range structure.Relationships {
		relationship := &Concept{
			GUID:        generateGUID(rel.Name),
			Name:        rel.Name,
			Description: rel.Description,
			Type:        "RelationshipType",
			Timestamp:   time.Now(),
		}
		err := addOrUpdateConcept(ctx, relationship)
		if err != nil {
			return fmt.Errorf("failed to add relationship type %s: %v", rel.Name, err)
		}
	}

	// First pass: create all concepts
	for _, node := range structure.Concepts {
		_, err := createConcepts(ctx, node, "")
		if err != nil {
			return fmt.Errorf("failed to create concept %s: %v", node.Name, err)
		}
	}

	// Second pass: create relationships
	for _, node := range structure.Concepts {
		err := createRelationships(ctx, node)
		if err != nil {
			return fmt.Errorf("failed to create relationships for concept %s: %v", node.Name, err)
		}
	}

	log.Printf("Bootstrapped %d concepts and %d relationship types", len(guidMap)-len(structure.Relationships), len(structure.Relationships))
	return nil
}
