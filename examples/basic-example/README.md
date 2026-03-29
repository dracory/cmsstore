# Basic Example - Getting Started with Custom Entities

This is the simplest possible example to get started with custom entities in cmsstore.

## What This Example Shows

- Setting up a CMS store with custom entities enabled
- Defining a simple entity type (`note`)
- Creating, reading, updating, and deleting entities
- Listing and counting entities
- Soft delete functionality

## Run the Example

```bash
cd examples/basic-example
go run main.go
```

## Expected Output

```
=== Initializing Database ===
✓ Database initialized

=== Creating CMS Store ===
✓ CMS store created
✓ Custom entity type 'note' registered

=== Creating a Note ===
✓ Created note with ID: [id]

=== Retrieving the Note ===
✓ Found note:
  ID: [id]
  Type: note
  Created: [timestamp]

=== Updating the Note ===
✓ Note updated successfully

=== Creating More Notes ===
✓ Created: Shopping List
✓ Created: Meeting Notes
✓ Created: Ideas

=== Listing All Notes ===
Total notes: 4
1. Note ID: [id] (Created: [timestamp])
2. Note ID: [id] (Created: [timestamp])
3. Note ID: [id] (Created: [timestamp])
4. Note ID: [id] (Created: [timestamp])

=== Counting Notes ===
Total count: 4

=== Deleting a Note ===
✓ Note deleted (moved to trash)

Remaining notes: 3
```

## Code Walkthrough

### 1. Define Entity Type

```go
CustomEntityDefinitions: []cmsstore.CustomEntityDefinition{
    {
        Type:      "note",
        TypeLabel: "Note",
        Group:     "Personal",
        Attributes: []cmsstore.CustomAttributeDefinition{
            {Name: "title", Type: "string", Label: "Title", Required: true},
            {Name: "content", Type: "string", Label: "Content"},
            {Name: "priority", Type: "int", Label: "Priority"},
        },
    },
}
```

### 2. Create Entity

```go
noteID, err := customStore.Create(ctx, "note", map[string]interface{}{
    "title":    "My First Note",
    "content":  "This is a simple note",
    "priority": 1,
}, nil, nil)
```

### 3. Retrieve Entity

```go
note, err := customStore.FindByID(ctx, noteID)
```

### 4. Update Entity

```go
err = customStore.Update(ctx, note, map[string]interface{}{
    "priority": 5,
})
```

### 5. List Entities

```go
allNotes, err := customStore.List(ctx, entitystore.EntityQueryOptions{
    EntityType: "note",
})
```

### 6. Delete Entity

```go
err = customStore.Delete(ctx, noteID)
```

## Key Concepts

- **No Migrations**: Add entity types without database schema changes
- **Flexible Attributes**: Define any attributes you need
- **Type Safety**: Attributes are type-checked (string, int, float, bool)
- **Soft Delete**: Deleted entities are moved to trash, not permanently removed

## Next Steps

Once you understand this basic example:

1. **Advanced Features**: See `../customentities/` for:
   - Relationships between entities
   - Taxonomy and categorization
   - More complex use cases

2. **CMS Integration**: See `../custom-block-types/` for:
   - Extending CMS blocks with custom types
   - Integration with existing CMS features

## Learn More

- [Custom Entities Documentation](../../docs/CUSTOM_ENTITIES.md)
- [Custom Entities Proposal](../../docs/proposals/2026-03-17-custom-entities-support.md)
- [entitystore Library](https://github.com/dracory/entitystore)
