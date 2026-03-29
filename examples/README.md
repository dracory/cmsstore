# CMS Store Examples

This directory contains comprehensive examples demonstrating various features of the cmsstore package.

## Available Examples

### 1. [basic-example](./basic-example/) - Getting Started ⭐

**Perfect for beginners!** The simplest possible example to understand custom entities.

**What you'll learn:**
- Setting up custom entities
- Creating, reading, updating, deleting entities
- Basic CRUD operations
- Soft delete functionality

**Run:**
```bash
cd basic-example
go run main.go
```

---

### 2. [customentities](./customentities/) - Advanced Features

**Comprehensive examples** covering all custom entity features.

**Examples included:**
- **01_basic_usage.go** - Detailed CRUD operations
- **02_relationships.go** - Linking entities together
- **03_taxonomy.go** - Categorization and tagging

**What you'll learn:**
- Advanced entity operations
- Relationships between entities (Author → Post → Comment)
- Taxonomy and categorization
- Complex data structures

**Run:**
```bash
cd customentities
go run 01_basic_usage.go
go run 02_relationships.go
go run 03_taxonomy.go
```

---

### 3. [custom-block-types](./custom-block-types/) - CMS Integration

**Real-world CMS use case** showing how to extend the CMS with custom block types.

**Block types demonstrated:**
- Hero blocks (banners with CTAs)
- Feature blocks (icon-based highlights)
- Testimonial blocks (customer reviews)
- FAQ blocks (questions and answers)

**What you'll learn:**
- Building flexible content management systems
- Creating reusable content blocks
- Organizing blocks by type
- Practical CMS patterns

**Run:**
```bash
cd custom-block-types
go run main.go
```

---

## Quick Start Guide

### Choose Your Path

**New to custom entities?**
→ Start with [basic-example](./basic-example/)

**Want to see all features?**
→ Explore [customentities](./customentities/)

**Building a CMS?**
→ Check out [custom-block-types](./custom-block-types/)

### Prerequisites

```bash
# Ensure you have Go installed
go version

# Navigate to any example folder
cd basic-example

# Run the example
go run main.go
```

## Feature Comparison

| Feature | basic-example | customentities | custom-block-types |
|---------|---------------|----------------|-------------------|
| Basic CRUD | ✅ | ✅ | ✅ |
| Entity Definitions | ✅ | ✅ | ✅ |
| Relationships | ❌ | ✅ | ⚠️ (shown in README) |
| Taxonomies | ❌ | ✅ | ⚠️ (shown in README) |
| Real-world Use Case | ❌ | ❌ | ✅ |
| Complexity | Simple | Advanced | Intermediate |
| Lines of Code | ~150 | ~700 | ~250 |

## Learning Path

```
1. basic-example/
   ↓ (Understand fundamentals)
   
2. customentities/01_basic_usage.go
   ↓ (Learn detailed CRUD)
   
3. customentities/02_relationships.go
   ↓ (Link entities together)
   
4. customentities/03_taxonomy.go
   ↓ (Categorize and tag)
   
5. custom-block-types/
   ↓ (Apply to real CMS)
   
6. Build your own!
```

## Common Use Cases

### Content Management System
→ See [custom-block-types](./custom-block-types/)
- Hero blocks, features, testimonials, FAQs
- Flexible page builders
- Reusable content components

### Blog Platform
→ See [customentities/02_relationships.go](./customentities/02_relationships.go)
- Authors, posts, comments
- Relationship management
- Content hierarchy

### E-commerce
→ Combine examples
- Products with categories (taxonomy)
- Product relationships (related products)
- Custom product types (variants)

### Knowledge Base
→ See [custom-block-types](./custom-block-types/) + [customentities/03_taxonomy.go](./customentities/03_taxonomy.go)
- FAQ blocks
- Article categorization
- Searchable content

## Key Concepts

### Custom Entities
Schema-less entities that can be defined without database migrations.

```go
CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
    {
        Type:      "product",
        TypeLabel: "Product",
        Attributes: []cmsstore.CustomAttributeDefinition{
            {Name: "title", Type: "string", Required: true},
            {Name: "price", Type: "float", Required: true},
        },
    },
}
```

### Relationships
Link entities together (belongs_to, has_many).

```go
postID, _ := customStore.Create(ctx, "post", attrs,
    []cmsstore.RelationshipDefinition{
        {Type: "belongs_to", TargetID: authorID},
    }, nil)
```

### Taxonomies
Categorize and tag entities.

```go
productID, _ := customStore.Create(ctx, "product", attrs, nil,
    []string{categoryTermID, tagTermID})
```

## Documentation

- [Custom Entities Documentation](../docs/CUSTOM_ENTITIES.md)
- [Custom Entities Proposal](../docs/proposals/2026-03-17-custom-entities-support.md)
- [entitystore Library](https://github.com/dracory/entitystore)

## Example Structure

Each example folder contains:
- `main.go` - Runnable example code
- `README.md` - Detailed documentation
- Comments explaining each step

## Tips for Learning

1. **Start Simple** - Begin with basic-example
2. **Read Comments** - Code is heavily commented
3. **Experiment** - Modify examples to see what happens
4. **Check Output** - Run examples to see results
5. **Read READMEs** - Each folder has detailed docs

## Troubleshooting

### Import errors
Make sure you're in the correct directory and have run `go mod tidy`.

### Database errors
Examples use in-memory SQLite, no setup required.

### "main redeclared" errors
These are expected - each example is a separate program. Run them individually.

## Contributing

Found an issue or want to add an example? Contributions welcome!

## License

Same as parent project.
