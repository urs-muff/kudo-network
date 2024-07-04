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

	// Fundamental Concepts
	TechnologyConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000001",
		Name:        "Technology",
		Description: "Tools and knowledge used to solve problems or improve conditions",
		Type:        "FundamentalConcept",
	}
	SelfConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000002",
		Name:        "Self",
		Description: "The essential qualities that make a person distinct from all others",
		Type:        "FundamentalConcept",
	}
	EgoConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000003",
		Name:        "Ego",
		Description: "The part of the mind that mediates between the conscious and the unconscious",
		Type:        "FundamentalConcept",
	}
	MathConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000004",
		Name:        "Mathematics",
		Description: "The abstract science of number, quantity, and space",
		Type:        "FundamentalConcept",
	}
	KnowledgeConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000005",
		Name:        "Knowledge",
		Description: "Facts, information, and skills acquired through experience or education",
		Type:        "FundamentalConcept",
	}
	WisdomConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000006",
		Name:        "Wisdom",
		Description: "The quality of having experience, knowledge, and good judgment",
		Type:        "FundamentalConcept",
	}
	ExperienceConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000007",
		Name:        "Experience",
		Description: "Practical contact with and observation of facts or events",
		Type:        "FundamentalConcept",
	}
	PropertyConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000008",
		Name:        "Property",
		Description: "A quality or characteristic belonging to or representative of something",
		Type:        "FundamentalConcept",
	}
	SoulConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000009",
		Name:        "Soul",
		Description: "The spiritual or immaterial part of a human being or animal",
		Type:        "FundamentalConcept",
	}
	BodyConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000010",
		Name:        "Body",
		Description: "The physical structure of a person or animal",
		Type:        "FundamentalConcept",
	}
	MindConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000011",
		Name:        "Mind",
		Description: "The element of a person that enables them to be aware of the world and their experiences",
		Type:        "FundamentalConcept",
	}
	BoundaryConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000012",
		Name:        "Boundary",
		Description: "A line that marks the limits of an area or conceptual division",
		Type:        "FundamentalConcept",
	}
	DualismConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000013",
		Name:        "Dualism",
		Description: "The division of something conceptually into two opposed or contrasted aspects",
		Type:        "FundamentalConcept",
	}
	SpaceConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000014",
		Name:        "Space",
		Description: "A continuous area or expanse that is free, available, or unoccupied",
		Type:        "FundamentalConcept",
	}
	TimeConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000015",
		Name:        "Time",
		Description: "The indefinite continued progress of existence and events in the past, present, and future",
		Type:        "FundamentalConcept",
	}
	IndividualExpressionConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000016",
		Name:        "Individual Expression",
		Description: "The unique way an individual communicates or manifests their personality or feelings",
		Type:        "FundamentalConcept",
	}
	SocietyConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000017",
		Name:        "Society",
		Description: "The aggregate of people living together in a more or less ordered community",
		Type:        "FundamentalConcept",
	}
	ZenConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000018",
		Name:        "Zen",
		Description: "A school of Mahayana Buddhism emphasizing the value of meditation and intuition",
		Type:        "FundamentalConcept",
	}
	PlatosCaveConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000019",
		Name:        "Plato's Cave",
		Description: "An allegory used to illustrate the way in which reality may be perceived",
		Type:        "FundamentalConcept",
	}
	UnityConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000020",
		Name:        "Unity",
		Description: "The state of being united or joined as a whole",
		Type:        "FundamentalConcept",
	}
	CoherenceConcept = CoreConcept{
		GUID:        "30000000-0000-0000-0000-000000000021",
		Name:        "Coherence",
		Description: "The quality of forming a unified whole; logical interconnection",
		Type:        "FundamentalConcept",
	}
	// Building Block Concepts
	InformationConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000001",
		Name:        "Information",
		Description: "Data, facts, or knowledge that can be communicated or received",
		Type:        "BuildingBlockConcept",
	}
	EnergyConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000002",
		Name:        "Energy",
		Description: "The capacity to do work or cause change in a system",
		Type:        "BuildingBlockConcept",
	}
	ConsciousnessConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000003",
		Name:        "Consciousness",
		Description: "The state of being aware of and responsive to one's surroundings",
		Type:        "BuildingBlockConcept",
	}
	PatternConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000004",
		Name:        "Pattern",
		Description: "A recurring theme, structure, or behavior in a system",
		Type:        "BuildingBlockConcept",
	}
	ProcessConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000005",
		Name:        "Process",
		Description: "A series of actions or steps taken to achieve a particular end",
		Type:        "BuildingBlockConcept",
	}
	IntentionConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000006",
		Name:        "Intention",
		Description: "A purposeful plan or aim guiding one's actions",
		Type:        "BuildingBlockConcept",
	}
	PerceptionConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000007",
		Name:        "Perception",
		Description: "The ability to see, hear, or become aware of something through the senses",
		Type:        "BuildingBlockConcept",
	}
	CommunicationConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000008",
		Name:        "Communication",
		Description: "The exchange of information between entities",
		Type:        "BuildingBlockConcept",
	}
	TransformationConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000009",
		Name:        "Transformation",
		Description: "A thorough or dramatic change in form, appearance, or character",
		Type:        "BuildingBlockConcept",
	}
	EmergenceConcept = CoreConcept{
		GUID:        "50000000-0000-0000-0000-000000000010",
		Name:        "Emergence",
		Description: "The process of coming into existence or prominence",
		Type:        "BuildingBlockConcept",
	}
	ComplexityConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000001",
		Name:        "Complexity",
		Description: "The state of having many interconnected parts or intricate details",
		Type:        "FundamentalConcept",
	}
	OrderConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000002",
		Name:        "Order",
		Description: "The arrangement or disposition of people or things in relation to each other according to a particular sequence, pattern, or method",
		Type:        "FundamentalConcept",
	}
	ChaosConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000003",
		Name:        "Chaos",
		Description: "Complete disorder and confusion",
		Type:        "FundamentalConcept",
	}
	SymmetryConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000004",
		Name:        "Symmetry",
		Description: "The quality of being made up of exactly similar parts facing each other or around an axis",
		Type:        "FundamentalConcept",
	}
	EntropyConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000005",
		Name:        "Entropy",
		Description: "A measure of the disorder or randomness in a system",
		Type:        "FundamentalConcept",
	}
	SynergyConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000006",
		Name:        "Synergy",
		Description: "The interaction or cooperation of two or more organizations, substances, or other agents to produce a combined effect greater than the sum of their separate effects",
		Type:        "FundamentalConcept",
	}
	ResonanceConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000007",
		Name:        "Resonance",
		Description: "The reinforcement or prolongation of sound by reflection from a surface or by the synchronous vibration of a neighboring object",
		Type:        "FundamentalConcept",
	}
	FractalConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000008",
		Name:        "Fractal",
		Description: "A never-ending pattern that repeats itself at different scales",
		Type:        "FundamentalConcept",
	}
	DynamicsConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000009",
		Name:        "Dynamics",
		Description: "The forces or properties which stimulate growth, development, or change within a system or process",
		Type:        "FundamentalConcept",
	}
	PotentialConcept = CoreConcept{
		GUID:        "60000000-0000-0000-0000-000000000010",
		Name:        "Potential",
		Description: "Latent qualities or abilities that may be developed and lead to future success or usefulness",
		Type:        "FundamentalConcept",
	}
	NetworkConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000001",
		Name:        "Network",
		Description: "A system of interconnected people, things, or ideas",
		Type:        "FundamentalConcept",
	}
	TribeConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000002",
		Name:        "Tribe",
		Description: "A social group connected by social, economic, religious, or blood ties",
		Type:        "FundamentalConcept",
	}
	CommunityConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000003",
		Name:        "Community",
		Description: "A group of people living in the same place or having characteristics in common",
		Type:        "FundamentalConcept",
	}
	PhysicalConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000004",
		Name:        "Physical",
		Description: "Relating to the body or material existence, as opposed to the mind or spirit",
		Type:        "FundamentalConcept",
	}
	MediumConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000005",
		Name:        "Medium",
		Description: "An intervening substance through which impressions are conveyed or forces act",
		Type:        "BuildingBlockConcept",
	}
	FormConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000006",
		Name:        "Form",
		Description: "The visible shape or configuration of something",
		Type:        "BuildingBlockConcept",
	}
	AbstractionConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000007",
		Name:        "Abstraction",
		Description: "The quality of dealing with ideas rather than events or concrete objects",
		Type:        "BuildingBlockConcept",
	}
	LogicConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000008",
		Name:        "Logic",
		Description: "Reasoning conducted according to strict principles of validity",
		Type:        "BuildingBlockConcept",
	}
	ReasoningConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000009",
		Name:        "Reasoning",
		Description: "The action of thinking about something in a logical way",
		Type:        "BuildingBlockConcept",
	}
	BasicPerceptionsConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000010",
		Name:        "Basic Perceptions",
		Description: "Fundamental sensory experiences that form the basis of more complex cognition",
		Type:        "BuildingBlockConcept",
	}
	ThoughtConcept = CoreConcept{
		GUID:        "80000000-0000-0000-0000-000000000011",
		Name:        "Thought",
		Description: "An idea or mental image formed by thinking",
		Type:        "FundamentalConcept",
	}
	// Coherence Investment and Reward System Concepts
	InvestmentConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000001",
		Name:        "Investment",
		Description: "The act of dedicating resources with the expectation of future benefits",
		Type:        "FundamentalConcept",
	}
	RewardConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000002",
		Name:        "Reward",
		Description: "A benefit given in recognition of effort, achievement, or contribution",
		Type:        "FundamentalConcept",
	}
	CoherenceScoreConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000003",
		Name:        "Coherence Score",
		Description: "A measure of alignment and harmony within a system or network",
		Type:        "FundamentalConcept",
	}
	ContributionConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000004",
		Name:        "Contribution",
		Description: "An action or input that adds value or supports the goals of a system",
		Type:        "FundamentalConcept",
	}
	CoherenceInvestmentSystemConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000005",
		Name:        "Coherence Investment System",
		Description: "A system that incentivizes and rewards actions that increase overall coherence",
		Type:        "FundamentalConcept",
	}
	MotivationConcept = CoreConcept{
		GUID:        "A0000000-0000-0000-0000-000000000006",
		Name:        "Motivation",
		Description: "The reason or reasons one has for acting or behaving in a particular way",
		Type:        "FundamentalConcept",
	}

	// New Relationship Types
	InfluencesRelationship = CoreRelationship{
		GUID:        "40000000-0000-0000-0000-000000000001",
		Name:        "Influences",
		Description: "Indicates that one concept has an effect on another",
	}
	ComposedOfRelationship = CoreRelationship{
		GUID:        "40000000-0000-0000-0000-000000000002",
		Name:        "Composed Of",
		Description: "Indicates that one concept is made up of or includes other concepts",
	}
	ManifestsAsRelationship = CoreRelationship{
		GUID:        "40000000-0000-0000-0000-000000000003",
		Name:        "Manifests As",
		Description: "Indicates that one concept is a concrete expression or instance of another",
	}
	EmergentFromRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000001",
		Name:        "Emergent From",
		Description: "Indicates that one concept arises as a result of complex interactions in another",
	}
	SymbioticWithRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000002",
		Name:        "Symbiotic With",
		Description: "Indicates a mutually beneficial relationship between concepts",
	}
	TransformsIntoRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000003",
		Name:        "Transforms Into",
		Description: "Indicates that one concept can change or evolve into another",
	}
	ResonatesWithRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000004",
		Name:        "Resonates With",
		Description: "Indicates a harmonious or synchronous relationship between concepts",
	}
	CatalyzesRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000005",
		Name:        "Catalyzes",
		Description: "Indicates that one concept initiates or accelerates the development or progress of another",
	}
	OpposedToRelationship = CoreRelationship{
		GUID:        "70000000-0000-0000-0000-000000000006",
		Name:        "Opposed To",
		Description: "Indicates that one concept is in opposition or contrast to another",
	}
	PartOfRelationship = CoreRelationship{
		GUID:        "90000000-0000-0000-0000-000000000001",
		Name:        "Part Of",
		Description: "Indicates that one concept is a component or subset of another",
	}
	ContrastsWith = CoreRelationship{
		GUID:        "90000000-0000-0000-0000-000000000002",
		Name:        "Contrasts With",
		Description: "Indicates that one concept is notably different from another in a specific aspect",
	}
	FacilitatesRelationship = CoreRelationship{
		GUID:        "90000000-0000-0000-0000-000000000003",
		Name:        "Facilitates",
		Description: "Indicates that one concept makes another concept easier or more likely to occur",
	}
	IncreasesRelationship = CoreRelationship{
		GUID:        "B0000000-0000-0000-0000-000000000001",
		Name:        "Increases",
		Description: "Indicates that one concept leads to an increase in another",
	}
	GeneratesRelationship = CoreRelationship{
		GUID:        "B0000000-0000-0000-0000-000000000002",
		Name:        "Generates",
		Description: "Indicates that one concept produces or creates another",
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

// BootstrapExpandedCoreConceptsAndRelationships initializes the system with the expanded set of core concepts and relationships
func BootstrapExpandedCoreConceptsAndRelationships(ctx context.Context) error {
	log.Println("Bootstrapping expanded core concepts and relationships...")

	// Add fundamental concepts
	fundamentalConcepts := []CoreConcept{
		TechnologyConcept, SelfConcept, EgoConcept, MathConcept, KnowledgeConcept,
		WisdomConcept, ExperienceConcept, PropertyConcept, SoulConcept, BodyConcept,
		MindConcept, BoundaryConcept, DualismConcept, SpaceConcept, TimeConcept,
		IndividualExpressionConcept, SocietyConcept,
		ZenConcept, PlatosCaveConcept, UnityConcept, CoherenceConcept,
		InformationConcept, EnergyConcept, ConsciousnessConcept, PatternConcept,
		ProcessConcept, IntentionConcept, PerceptionConcept, CommunicationConcept,
		TransformationConcept, EmergenceConcept,
		ComplexityConcept, OrderConcept, ChaosConcept, SymmetryConcept,
		EntropyConcept, SynergyConcept, ResonanceConcept, FractalConcept,
		DynamicsConcept, PotentialConcept,
		NetworkConcept, TribeConcept, CommunityConcept, PhysicalConcept,
		MediumConcept, FormConcept, AbstractionConcept, LogicConcept,
		ReasoningConcept, BasicPerceptionsConcept, ThoughtConcept,
		InvestmentConcept, RewardConcept, CoherenceScoreConcept,
		ContributionConcept, CoherenceInvestmentSystemConcept,
		MotivationConcept,
	}

	for _, concept := range fundamentalConcepts {
		if err := addCoreConcept(ctx, concept); err != nil {
			return err
		}
	}

	// Add new relationship types
	newRelationships := []CoreRelationship{
		InfluencesRelationship, ComposedOfRelationship, ManifestsAsRelationship,
		EmergentFromRelationship, SymbioticWithRelationship, TransformsIntoRelationship,
		ResonatesWithRelationship, CatalyzesRelationship, OpposedToRelationship,
		PartOfRelationship, ContrastsWith, FacilitatesRelationship,
		IncreasesRelationship, GeneratesRelationship,
	}

	for _, relationship := range newRelationships {
		if err := addCoreRelationship(ctx, relationship); err != nil {
			return err
		}
	}

	// Create some example relationships between concepts
	relationships := []struct {
		source  GUID
		relType GUID
		target  GUID
	}{
		{TechnologyConcept.GUID, InfluencesRelationship.GUID, SocietyConcept.GUID},
		{MindConcept.GUID, ComposedOfRelationship.GUID, EgoConcept.GUID},
		{SelfConcept.GUID, ManifestsAsRelationship.GUID, IndividualExpressionConcept.GUID},
		{KnowledgeConcept.GUID, InfluencesRelationship.GUID, WisdomConcept.GUID},
		{ExperienceConcept.GUID, InfluencesRelationship.GUID, KnowledgeConcept.GUID},
		{BodyConcept.GUID, RelatedTo.GUID, MindConcept.GUID},
		{SpaceConcept.GUID, RelatedTo.GUID, TimeConcept.GUID},
		{DualismConcept.GUID, RelatedTo.GUID, BoundaryConcept.GUID},
		{ZenConcept.GUID, RelatedTo.GUID, MeditationInteraction.GUID},
		{ZenConcept.GUID, InfluencesRelationship.GUID, WisdomConcept.GUID},
		{PlatosCaveConcept.GUID, RelatedTo.GUID, KnowledgeConcept.GUID},
		{PlatosCaveConcept.GUID, InfluencesRelationship.GUID, ExperienceConcept.GUID},
		{UnityConcept.GUID, RelatedTo.GUID, CoherenceConcept.GUID},
		{UnityConcept.GUID, InfluencesRelationship.GUID, SocietyConcept.GUID},
		{CoherenceConcept.GUID, InfluencesRelationship.GUID, FlowStateInteraction.GUID},
		{CoherenceConcept.GUID, RelatedTo.GUID, WisdomConcept.GUID},
		{KnowledgeConcept.GUID, ComposedOfRelationship.GUID, InformationConcept.GUID},
		{ExperienceConcept.GUID, ComposedOfRelationship.GUID, PerceptionConcept.GUID},
		{TechnologyConcept.GUID, ComposedOfRelationship.GUID, ProcessConcept.GUID},
		{ZenConcept.GUID, RelatedTo.GUID, ConsciousnessConcept.GUID},
		{CoherenceConcept.GUID, RelatedTo.GUID, PatternConcept.GUID},
		{FlowStateInteraction.GUID, RelatedTo.GUID, EnergyConcept.GUID},
		{SelfConcept.GUID, ComposedOfRelationship.GUID, IntentionConcept.GUID},
		{SocietyConcept.GUID, ComposedOfRelationship.GUID, CommunicationConcept.GUID},
		{WisdomConcept.GUID, RelatedTo.GUID, TransformationConcept.GUID},
		{UnityConcept.GUID, RelatedTo.GUID, EmergenceConcept.GUID},
		{ComplexityConcept.GUID, EmergentFromRelationship.GUID, PatternConcept.GUID},
		{OrderConcept.GUID, InfluencesRelationship.GUID, CoherenceConcept.GUID},
		{OrderConcept.GUID, OpposedToRelationship.GUID, ChaosConcept.GUID},
		{ChaosConcept.GUID, TransformsIntoRelationship.GUID, OrderConcept.GUID},
		{SymmetryConcept.GUID, RelatedTo.GUID, PatternConcept.GUID},
		{EntropyConcept.GUID, OpposedToRelationship.GUID, OrderConcept.GUID},
		{SynergyConcept.GUID, EmergentFromRelationship.GUID, RelationshipType.GUID},
		{ResonanceConcept.GUID, ResonatesWithRelationship.GUID, FlowStateInteraction.GUID},
		{FractalConcept.GUID, ManifestsAsRelationship.GUID, PatternConcept.GUID},
		{DynamicsConcept.GUID, InfluencesRelationship.GUID, ProcessConcept.GUID},
		{PotentialConcept.GUID, TransformsIntoRelationship.GUID, EnergyConcept.GUID},
		{ZenConcept.GUID, CatalyzesRelationship.GUID, CoherenceConcept.GUID},
		{MeditationInteraction.GUID, CatalyzesRelationship.GUID, ConsciousnessConcept.GUID},
		{NetworkConcept.GUID, RelatedTo.GUID, CommunityConcept.GUID},
		{TribeConcept.GUID, IsA.GUID, CommunityConcept.GUID},
		{CommunityConcept.GUID, PartOfRelationship.GUID, SocietyConcept.GUID},
		{PhysicalConcept.GUID, ContrastsWith.GUID, AbstractionConcept.GUID},
		{MediumConcept.GUID, RelatedTo.GUID, CommunicationConcept.GUID},
		{FormConcept.GUID, ManifestsAsRelationship.GUID, PhysicalConcept.GUID},
		{AbstractionConcept.GUID, RelatedTo.GUID, MathConcept.GUID},
		{LogicConcept.GUID, PartOfRelationship.GUID, ReasoningConcept.GUID},
		{ReasoningConcept.GUID, InfluencesRelationship.GUID, KnowledgeConcept.GUID},
		{BasicPerceptionsConcept.GUID, InfluencesRelationship.GUID, ExperienceConcept.GUID},
		{NetworkConcept.GUID, ManifestsAsRelationship.GUID, SocietyConcept.GUID},
		{TribeConcept.GUID, PartOfRelationship.GUID, NetworkConcept.GUID},
		{PhysicalConcept.GUID, RelatedTo.GUID, BodyConcept.GUID},
		{MediumConcept.GUID, FacilitatesRelationship.GUID, InformationConcept.GUID},
		{FormConcept.GUID, RelatedTo.GUID, PatternConcept.GUID},
		{AbstractionConcept.GUID, InfluencesRelationship.GUID, ThoughtConcept.GUID},
		{LogicConcept.GUID, InfluencesRelationship.GUID, MathConcept.GUID},
		{ReasoningConcept.GUID, RelatedTo.GUID, WisdomConcept.GUID},
		{BasicPerceptionsConcept.GUID, ComposedOfRelationship.GUID, PerceptionConcept.GUID},
		{InvestmentConcept.GUID, IncreasesRelationship.GUID, CoherenceScoreConcept.GUID},
		{ContributionConcept.GUID, GeneratesRelationship.GUID, RewardConcept.GUID},
		{CoherenceScoreConcept.GUID, InfluencesRelationship.GUID, RewardConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, ComposedOfRelationship.GUID, InvestmentConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, ComposedOfRelationship.GUID, RewardConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, ComposedOfRelationship.GUID, CoherenceScoreConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, ComposedOfRelationship.GUID, ContributionConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, IncreasesRelationship.GUID, CoherenceConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, FacilitatesRelationship.GUID, SynergyConcept.GUID},
		{ContributionConcept.GUID, IncreasesRelationship.GUID, CoherenceScoreConcept.GUID},
		{RewardConcept.GUID, IncreasesRelationship.GUID, MotivationConcept.GUID},
		{MotivationConcept.GUID, InfluencesRelationship.GUID, ContributionConcept.GUID},
		{MotivationConcept.GUID, RelatedTo.GUID, IntentionConcept.GUID},
		{MotivationConcept.GUID, InfluencesRelationship.GUID, InvestmentConcept.GUID},
		{CoherenceInvestmentSystemConcept.GUID, InfluencesRelationship.GUID, MotivationConcept.GUID},
		{MotivationConcept.GUID, FacilitatesRelationship.GUID, FlowStateInteraction.GUID},
	}

	for _, rel := range relationships {
		if err := createCoreRelationship(ctx, rel.source, rel.relType, rel.target); err != nil {
			return err
		}
	}

	log.Println("Expanded core concepts and relationships bootstrapped successfully")
	return nil
}

// Call this function in your main.go after initializing the IPFS node
func InitializeSystem(ctx context.Context) error {
	if err := BootstrapCoreConceptsAndRelationships(ctx); err != nil {
		log.Printf("Error during boostrapping core concepts: %\n", err)
		return err
	}
	if err := BootstrapExpandedCoreConceptsAndRelationships(ctx); err != nil {
		return err
	}
	return nil
}
