package cmsstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dracory/entitystore"
)

// CustomEntityStore wraps entitystore to provide CMS-specific custom entity functionality.
type CustomEntityStore struct {
	inner       entitystore.StoreInterface
	definitions map[string]CustomEntityDefinition // Entity type -> definition mapping
}

// NewCustomEntityStore creates a new custom entity store wrapper.
func NewCustomEntityStore(db *sql.DB, options CustomEntityStoreOptions) (*CustomEntityStore, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	// Set defaults
	if options.EntityTableName == "" {
		options.EntityTableName = "cms_custom_entity"
	}
	if options.AttributeTableName == "" {
		options.AttributeTableName = "cms_custom_attribute"
	}
	if options.EntityTrashTableName == "" {
		options.EntityTrashTableName = "cms_custom_entity_trash"
	}
	if options.AttributeTrashTableName == "" {
		options.AttributeTrashTableName = "cms_custom_attribute_trash"
	}

	// Create entitystore with optional features
	entityStoreOpts := entitystore.NewStoreOptions{
		DB:                      db,
		EntityTableName:         options.EntityTableName,
		AttributeTableName:      options.AttributeTableName,
		EntityTrashTableName:    options.EntityTrashTableName,
		AttributeTrashTableName: options.AttributeTrashTableName,
		AutomigrateEnabled:      options.AutomigrateEnabled,
		RelationshipsEnabled:    options.RelationshipsEnabled,
		TaxonomiesEnabled:       options.TaxonomiesEnabled,
	}

	if options.RelationshipsEnabled {
		if options.RelationshipTableName == "" {
			options.RelationshipTableName = "cms_custom_relationship"
		}
		if options.RelationshipTrashTableName == "" {
			options.RelationshipTrashTableName = "cms_custom_relationship_trash"
		}
		entityStoreOpts.RelationshipTableName = options.RelationshipTableName
		entityStoreOpts.RelationshipTrashTableName = options.RelationshipTrashTableName
	}

	if options.TaxonomiesEnabled {
		if options.TaxonomyTableName == "" {
			options.TaxonomyTableName = "cms_custom_taxonomy"
		}
		if options.TaxonomyTermTableName == "" {
			options.TaxonomyTermTableName = "cms_custom_taxonomy_term"
		}
		if options.EntityTaxonomyTableName == "" {
			options.EntityTaxonomyTableName = "cms_custom_entity_taxonomy"
		}
		entityStoreOpts.TaxonomyTableName = options.TaxonomyTableName
		entityStoreOpts.TaxonomyTermTableName = options.TaxonomyTermTableName
		entityStoreOpts.EntityTaxonomyTableName = options.EntityTaxonomyTableName
		entityStoreOpts.TaxonomyTrashTableName = options.TaxonomyTrashTableName
		entityStoreOpts.TaxonomyTermTrashTableName = options.TaxonomyTermTrashTableName
	}

	inner, err := entitystore.NewStore(entityStoreOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create entitystore: %w", err)
	}

	return &CustomEntityStore{
		inner:       inner,
		definitions: make(map[string]CustomEntityDefinition),
	}, nil
}

// RegisterEntityType registers a custom entity type definition.
func (s *CustomEntityStore) RegisterEntityType(def CustomEntityDefinition) error {
	if def.Type == "" {
		return errors.New("entity type is required")
	}
	if def.TypeLabel == "" {
		return errors.New("entity type label is required")
	}

	s.definitions[def.Type] = def
	return nil
}

// GetEntityDefinition retrieves the definition for an entity type.
func (s *CustomEntityStore) GetEntityDefinition(entityType string) (CustomEntityDefinition, bool) {
	def, ok := s.definitions[entityType]
	return def, ok
}

// GetAllDefinitions returns all registered entity type definitions.
func (s *CustomEntityStore) GetAllDefinitions() []CustomEntityDefinition {
	defs := make([]CustomEntityDefinition, 0, len(s.definitions))
	for _, def := range s.definitions {
		defs = append(defs, def)
	}
	return defs
}

// Create creates a new custom entity with attributes, relationships, and taxonomy assignments.
func (s *CustomEntityStore) Create(
	ctx context.Context,
	entityType string,
	attrs map[string]interface{},
	relationships []RelationshipDefinition,
	taxonomyTermIDs []string,
) (string, error) {
	// Validate entity type is registered
	def, ok := s.definitions[entityType]
	if !ok {
		return "", fmt.Errorf("entity type '%s' is not registered", entityType)
	}

	// Validate required attributes
	if err := s.validateAttributes(def, attrs); err != nil {
		return "", err
	}

	// Create entity via entitystore
	entity := entitystore.NewEntity()
	entity.SetType(entityType)

	if err := s.inner.EntityCreate(ctx, entity); err != nil {
		return "", fmt.Errorf("failed to create entity: %w", err)
	}

	// Set attributes after entity creation
	for key, value := range attrs {
		if err := s.setAttributeValue(ctx, entity.ID(), key, value); err != nil {
			return "", fmt.Errorf("failed to set attribute '%s': %w", key, err)
		}
	}

	// Add relationships if enabled and provided
	if def.AllowRelationships && len(relationships) > 0 {
		for _, rel := range relationships {
			if err := s.validateRelationshipType(def, rel.Type); err != nil {
				return "", err
			}

			// Serialize metadata to JSON string
			metadataStr := ""
			if rel.Metadata != nil {
				metadataBytes, _ := json.Marshal(rel.Metadata)
				metadataStr = string(metadataBytes)
			}

			_, err := s.inner.RelationshipCreateByOptions(ctx, entitystore.RelationshipOptions{
				EntityID:         entity.ID(),
				RelatedEntityID:  rel.TargetID,
				RelationshipType: rel.Type,
				Metadata:         metadataStr,
			})
			if err != nil {
				return "", fmt.Errorf("failed to create relationship: %w", err)
			}
		}
	}

	// Assign taxonomy terms if enabled and provided
	if def.AllowTaxonomies && len(taxonomyTermIDs) > 0 {
		for _, termID := range taxonomyTermIDs {
			// Get term to find its taxonomy
			term, err := s.inner.TaxonomyTermFind(ctx, termID)
			if err != nil {
				return "", fmt.Errorf("failed to find taxonomy term '%s': %w", termID, err)
			}
			if term == nil {
				return "", fmt.Errorf("taxonomy term not found: %s", termID)
			}

			// Validate taxonomy is allowed for this entity type
			if !s.isTaxonomyAllowed(def, term.GetTaxonomyID()) {
				return "", fmt.Errorf("taxonomy '%s' is not allowed for entity type '%s'", term.GetTaxonomyID(), entityType)
			}

			if err := s.inner.EntityTaxonomyAssign(ctx, entity.ID(), term.GetTaxonomyID(), termID); err != nil {
				return "", fmt.Errorf("failed to assign taxonomy term: %w", err)
			}
		}
	}

	return entity.ID(), nil
}

// FindByID retrieves a custom entity by ID.
func (s *CustomEntityStore) FindByID(ctx context.Context, entityID string) (entitystore.EntityInterface, error) {
	return s.inner.EntityFindByID(ctx, entityID)
}

// List retrieves custom entities with optional filtering.
func (s *CustomEntityStore) List(ctx context.Context, options entitystore.EntityQueryOptions) ([]entitystore.EntityInterface, error) {
	return s.inner.EntityList(ctx, options)
}

// Update updates a custom entity's attributes.
func (s *CustomEntityStore) Update(ctx context.Context, entity entitystore.EntityInterface, attrs map[string]interface{}) error {
	// Validate entity type is registered
	def, ok := s.definitions[entity.GetType()]
	if !ok {
		return fmt.Errorf("entity type '%s' is not registered", entity.GetType())
	}

	// Validate attributes
	if err := s.validateAttributes(def, attrs); err != nil {
		return err
	}

	// Update attributes
	for key, value := range attrs {
		if err := s.setAttributeValue(ctx, entity.ID(), key, value); err != nil {
			return fmt.Errorf("failed to set attribute '%s': %w", key, err)
		}
	}

	return s.inner.EntityUpdate(ctx, entity)
}

// Delete soft-deletes a custom entity.
func (s *CustomEntityStore) Delete(ctx context.Context, entityID string) error {
	_, err := s.inner.EntityTrash(ctx, entityID)
	return err
}

// Count counts custom entities matching the query options.
func (s *CustomEntityStore) Count(ctx context.Context, options entitystore.EntityQueryOptions) (int64, error) {
	return s.inner.EntityCount(ctx, options)
}

// GetRelationships retrieves relationships for an entity.
func (s *CustomEntityStore) GetRelationships(ctx context.Context, entityID string) ([]entitystore.RelationshipInterface, error) {
	return s.inner.RelationshipList(ctx, entitystore.RelationshipQueryOptions{
		EntityID: entityID,
	})
}

// GetTaxonomyAssignments retrieves taxonomy assignments for an entity.
func (s *CustomEntityStore) GetTaxonomyAssignments(ctx context.Context, entityID string) ([]entitystore.EntityTaxonomyInterface, error) {
	return s.inner.EntityTaxonomyList(ctx, entitystore.EntityTaxonomyQueryOptions{
		EntityID: entityID,
	})
}

// Inner returns the underlying entitystore for advanced operations.
func (s *CustomEntityStore) Inner() entitystore.StoreInterface {
	return s.inner
}

// validateAttributes validates that required attributes are present and types are correct.
func (s *CustomEntityStore) validateAttributes(def CustomEntityDefinition, attrs map[string]interface{}) error {
	// Check required attributes
	for _, attrDef := range def.Attributes {
		if attrDef.Required {
			if _, ok := attrs[attrDef.Name]; !ok {
				return fmt.Errorf("required attribute '%s' is missing", attrDef.Name)
			}
		}
	}

	return nil
}

// setAttributeValue sets an attribute value using entitystore's attribute system.
func (s *CustomEntityStore) setAttributeValue(ctx context.Context, entityID string, key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		return s.inner.AttributeSetString(ctx, entityID, key, v)
	case int:
		return s.inner.AttributeSetInt(ctx, entityID, key, int64(v))
	case int64:
		return s.inner.AttributeSetInt(ctx, entityID, key, v)
	case float64:
		return s.inner.AttributeSetFloat(ctx, entityID, key, v)
	case float32:
		return s.inner.AttributeSetFloat(ctx, entityID, key, float64(v))
	case bool:
		if v {
			return s.inner.AttributeSetInt(ctx, entityID, key, 1)
		}
		return s.inner.AttributeSetInt(ctx, entityID, key, 0)
	default:
		// For complex types, convert to string representation
		return s.inner.AttributeSetString(ctx, entityID, key, fmt.Sprintf("%v", v))
	}
}

// Note: Helper method for retrieving attribute values removed as it's not currently used.
// If needed, use: customStore.Inner().AttributeFind(ctx, entityID, key)

// validateRelationshipType validates that a relationship type is allowed for the entity.
func (s *CustomEntityStore) validateRelationshipType(def CustomEntityDefinition, relType string) error {
	if !def.AllowRelationships {
		return fmt.Errorf("relationships are not allowed for entity type '%s'", def.Type)
	}

	if len(def.AllowedRelationTypes) == 0 {
		return nil // All relationship types allowed
	}

	for _, allowed := range def.AllowedRelationTypes {
		if allowed == relType {
			return nil
		}
	}

	return fmt.Errorf("relationship type '%s' is not allowed for entity type '%s'", relType, def.Type)
}

// isTaxonomyAllowed checks if a taxonomy is allowed for the entity type.
func (s *CustomEntityStore) isTaxonomyAllowed(def CustomEntityDefinition, taxonomyID string) bool {
	if !def.AllowTaxonomies {
		return false
	}

	if len(def.TaxonomyIDs) == 0 {
		return true // All taxonomies allowed
	}

	for _, allowed := range def.TaxonomyIDs {
		if allowed == taxonomyID {
			return true
		}
	}

	return false
}

// CustomEntityStoreOptions holds configuration options for the custom entity store.
type CustomEntityStoreOptions struct {
	EntityTableName         string
	AttributeTableName      string
	EntityTrashTableName    string
	AttributeTrashTableName string
	AutomigrateEnabled      bool

	// Relationship options
	RelationshipsEnabled       bool
	RelationshipTableName      string
	RelationshipTrashTableName string

	// Taxonomy options
	TaxonomiesEnabled          bool
	TaxonomyTableName          string
	TaxonomyTermTableName      string
	EntityTaxonomyTableName    string
	TaxonomyTrashTableName     string
	TaxonomyTermTrashTableName string
}
