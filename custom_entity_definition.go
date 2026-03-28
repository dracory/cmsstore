package cmsstore

// CustomAttributeDefinition defines the structure of a custom entity attribute.
type CustomAttributeDefinition struct {
	Name        string // Attribute name (e.g., "title", "price")
	Type        string // Data type: "string", "int", "float", "bool", "json"
	Label       string // Human-readable label for admin UI
	Required    bool   // Whether the attribute is required
	DefaultValue interface{} // Default value for the attribute
	Validation  string // Validation rules (optional)
	Help        string // Help text for admin UI
}

// CustomEntityDefinition defines a custom entity type configuration.
type CustomEntityDefinition struct {
	Type       string                      // Entity type identifier (e.g., "shop_product", "blog_post")
	TypeLabel  string                      // Human-readable label (e.g., "Product", "Blog Post")
	Group      string                      // Group for admin navigation (e.g., "Shop", "Blog")
	Icon       string                      // Icon for admin UI (optional)
	Attributes []CustomAttributeDefinition // Attribute definitions
	
	// Relationship configuration
	AllowRelationships bool     // Whether this entity type supports relationships
	AllowedRelationTypes []string // Allowed relationship types (e.g., "belongs_to", "has_many")
	
	// Taxonomy configuration
	AllowTaxonomies bool     // Whether this entity type supports taxonomies
	TaxonomyIDs     []string // IDs of taxonomies that can be assigned to this entity type
}

// RelationshipDefinition defines a relationship to be created with an entity.
type RelationshipDefinition struct {
	TargetID         string                 // ID of the related entity
	Type             string                 // Relationship type (e.g., "belongs_to", "has_many")
	Metadata         map[string]interface{} // Additional metadata for the relationship
}
