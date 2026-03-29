# Custom Block Types Example

This example demonstrates how to extend the CMS with custom block types using custom entities. Perfect for building flexible content management systems with specialized block types.

## What This Example Shows

- Defining multiple custom block types (Hero, Feature, Testimonial, FAQ)
- Creating blocks with different attribute structures
- Organizing blocks by type and category
- Building a flexible content block system
- Practical CMS use cases

## Run the Example

```bash
cd examples/custom-block-types
go run main.go
```

## Block Types Defined

### 1. Hero Block
Large banner with image and call-to-action button.

**Attributes:**
- `title` (string, required) - Main headline
- `subtitle` (string) - Supporting text
- `background_image` (string) - Background image URL
- `cta_text` (string) - Call-to-action button text
- `cta_link` (string) - Button destination URL
- `alignment` (string) - Text alignment (left, center, right)

**Use Case:** Landing page headers, promotional banners

### 2. Feature Block
Icon-based feature highlights.

**Attributes:**
- `icon` (string, required) - Icon identifier
- `title` (string, required) - Feature title
- `description` (string) - Feature description
- `link` (string) - Learn more link

**Use Case:** Product features, service highlights, benefits

### 3. Testimonial Block
Customer testimonials with ratings.

**Attributes:**
- `quote` (string, required) - Testimonial text
- `author_name` (string, required) - Customer name
- `author_title` (string) - Job title/company
- `author_image` (string) - Profile image URL
- `rating` (int) - Star rating (1-5)

**Use Case:** Social proof, customer reviews, case studies

### 4. FAQ Block
Frequently asked questions.

**Attributes:**
- `question` (string, required) - Question text
- `answer` (string, required) - Answer text
- `category` (string) - FAQ category
- `order` (int) - Display order

**Use Case:** Help documentation, support pages, knowledge base

## Example Output

```
=== Custom Block Types Example ===
Demonstrating how to extend CMS with custom block types

=== Creating Hero Block ===
✓ Created hero block: [id]

=== Creating Feature Blocks ===
✓ Created feature: Fast Performance
✓ Created feature: Secure by Default
✓ Created feature: Developer Friendly

=== Creating Testimonial Blocks ===
✓ Created testimonial from: Sarah Johnson
✓ Created testimonial from: Mike Chen

=== Creating FAQ Blocks ===
✓ Created FAQ: How do custom blocks work?
✓ Created FAQ: Can I create my own block types?
✓ Created FAQ: Are custom blocks searchable?

=== Listing Blocks by Type ===
hero_block          : 1 block(s)
feature_block       : 3 block(s)
testimonial_block   : 2 block(s)
faq_block           : 3 block(s)

=== Total Custom Blocks ===
Total custom blocks created: 9
```

## How to Use in Your CMS

### 1. Define Block Types

```go
CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
    {
        Type:      "hero_block",
        TypeLabel: "Hero Block",
        Group:     "Content Blocks",
        Attributes: []cmsstore.CustomAttributeDefinition{
            {Name: "title", Type: "string", Label: "Title", Required: true},
            {Name: "subtitle", Type: "string", Label: "Subtitle"},
            // ... more attributes
        },
    },
}
```

### 2. Create Blocks

```go
heroID, err := customStore.Create(ctx, "hero_block", map[string]interface{}{
    "title":    "Welcome",
    "subtitle": "Get started today",
    "cta_text": "Sign Up",
}, nil, nil)
```

### 3. Retrieve Blocks for a Page

```go
heroBlocks, err := customStore.List(ctx, entitystore.EntityQueryOptions{
    EntityType: "hero_block",
})
```

### 4. Render in Templates

```go
for _, block := range heroBlocks {
    // Get attributes and render HTML
    // Use block.ID() to identify the block
}
```

## Real-World Use Cases

### Landing Page Builder
```
Hero Block → Feature Blocks → Testimonial Blocks → FAQ Blocks → CTA Block
```

### Marketing Website
- Product pages with feature blocks
- Customer stories with testimonial blocks
- Support pages with FAQ blocks
- Campaign landing pages with hero blocks

### Documentation Site
- Getting started guides with FAQ blocks
- Feature documentation with feature blocks
- User testimonials for social proof

### E-commerce Site
- Product highlights with feature blocks
- Customer reviews with testimonial blocks
- Product FAQs with FAQ blocks

## Extending This Example

### Add More Block Types

```go
// Gallery Block
{
    Type:      "gallery_block",
    TypeLabel: "Gallery Block",
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "images", Type: "string", Label: "Image URLs (JSON)"},
        {Name: "layout", Type: "string", Label: "Layout Type"},
    },
}

// Video Block
{
    Type:      "video_block",
    TypeLabel: "Video Block",
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "video_url", Type: "string", Label: "Video URL", Required: true},
        {Name: "thumbnail", Type: "string", Label: "Thumbnail URL"},
        {Name: "autoplay", Type: "bool", Label: "Autoplay"},
    },
}

// Pricing Block
{
    Type:      "pricing_block",
    TypeLabel: "Pricing Block",
    Attributes: []cmsstore.CustomAttributeDefinition{
        {Name: "plan_name", Type: "string", Label: "Plan Name", Required: true},
        {Name: "price", Type: "float", Label: "Price", Required: true},
        {Name: "features", Type: "string", Label: "Features (JSON)"},
    },
}
```

### Add Relationships

Link blocks to pages:

```go
heroID, err := customStore.Create(ctx, "hero_block", attrs,
    []cmsstore.RelationshipDefinition{
        {
            Type:     "belongs_to",
            TargetID: pageID,
            Metadata: map[string]interface{}{
                "position": "header",
                "order":    1,
            },
        },
    }, nil)
```

### Add Taxonomies

Categorize blocks:

```go
// Create category taxonomy
categoryTaxonomy, _ := innerStore.TaxonomyCreateByOptions(ctx, 
    entitystore.TaxonomyOptions{
        Name: "Block Categories",
        Slug: "block-categories",
    })

// Create terms
marketingTerm, _ := innerStore.TaxonomyTermCreateByOptions(ctx,
    entitystore.TaxonomyTermOptions{
        TaxonomyID: categoryTaxonomy.ID(),
        Name:       "Marketing",
        Slug:       "marketing",
    })

// Assign to blocks
heroID, _ := customStore.Create(ctx, "hero_block", attrs, nil, 
    []string{marketingTerm.ID()})
```

## Benefits

✓ **No Code Changes** - Add new block types without modifying code
✓ **Flexible Structure** - Each block type has its own attributes
✓ **Reusable** - Use the same block across multiple pages
✓ **Searchable** - All block content is indexed and searchable
✓ **Version Control** - Track changes to blocks over time
✓ **Type Safe** - Attributes are validated by type

## Integration Tips

1. **Admin UI**: Build forms dynamically based on attribute definitions
2. **Preview**: Render blocks in real-time as editors create them
3. **Templates**: Create reusable templates for each block type
4. **Validation**: Add custom validation rules for specific block types
5. **Import/Export**: Bulk import/export blocks for content migration

## Next Steps

- See `../basic-example` for fundamental concepts
- See `../customentities` for advanced features (relationships, taxonomies)
- Check the [Custom Entities Documentation](../../docs/CUSTOM_ENTITIES.md)

## Learn More

- [Custom Entities Proposal](../../docs/proposals/2026-03-17-custom-entities-support.md)
- [entitystore Library](https://github.com/dracory/entitystore)
- [CMS Store Documentation](../../docs/)
