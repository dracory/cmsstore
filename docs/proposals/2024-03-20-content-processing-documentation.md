# [Implemented] Content Processing Pipeline Documentation

## Status
**✅ IMPLEMENTED** - Documentation exists and is comprehensive

## Summary
- **Problem**: Content processing pipeline needed clear documentation
- **Solution**: Created comprehensive documentation covering the entire pipeline

## Where to Find the Documentation

### 1. Frontend Overview (`docs/frontend/overview.md`)
Complete documentation of the content processing pipeline including:

**Request Flow Diagram:**
- HTTP Request → Domain Validation → Site Resolution → Page Resolution
- Template Handling → Content Processing → Final HTML

**Content Processing Sequence:**
```
Template Merge (if exists) 
  → Replace Placeholders (PageTitle, PageContent, etc.)
  → Render Blocks ([[BLOCK_id]])
  → Render Page URLs ([[PAGE_URL_id]])
  → Apply Shortcodes (<shortcode>)
  → Process Translations ([[TRANSLATION_id]])
  → Apply Middlewares
  → Final HTML
```

**Key Components Documented:**
- **Blocks**: `[[BLOCK_blockID]]` syntax, caching, type-based rendering
- **Templates**: Layout definition, placeholder substitution
- **Shortcodes**: Custom content generators with `<shortcode>` syntax
- **Translations**: Multi-language support via `[[TRANSLATION_id]]`
- **URL Patterns**: Dynamic routing (`:any`, `:num`, `:alpha`, etc.)
- **Caching**: TTL-based caching system
- **Middleware**: Request/response processing

### 2. Block System Architecture (`docs/BLOCK_SYSTEM_ARCHITECTURE.md`)
Three-layer architecture documentation:
- **`blocks/`** - Unified built-in types (HTML, Menu, Navbar, Breadcrumbs)
- **`frontend/blocks/`** - Frontend renderers
- **`admin/blocks/`** - Admin UI providers

### 3. Actual Implementation Code (`frontend/frontend.go`)
The rendering pipeline is clearly documented in code comments:

```go
// renderContentToHtml - Processing sequence (order is important):
// 1. Replace placeholders with values
// 2. Render blocks
// 3. Apply block attribute syntax (<block id="..." />)
// 4. Render page URLs
// 5. Apply shortcodes
// 6. Render translations
// 7. Return final HTML
```

**Key Functions:**
- `PageRenderHtmlBySiteAndAlias()` - Main entry point
- `pageOrTemplateContent()` - Template merging logic
- `renderContentToHtml()` - Core processing pipeline
- `contentRenderBlocks()` - Block rendering
- `applyShortcodes()` - Shortcode processing
- `contentRenderTranslations()` - Translation handling
- `applyMiddlewares()` - Middleware application

### 4. Code-Level Documentation

**Template Merging:**
- Happens in `pageOrTemplateContent()` BEFORE content processing
- If page has `TemplateID()`, template content is loaded
- Page content becomes available via `[[PageContent]]` placeholder

**Block Rendering:**
- `fetchBlockContent()` - Loads and renders blocks with caching
- `renderBlockByType()` - Dispatches to type-specific renderers
- Registry pattern: `BlockRendererRegistry`

**Error Handling:**
- Missing templates → Returns page content as-is
- Invalid blocks → Logged, empty content returned
- Missing translations → Empty string replacement
- Middleware errors → Logged, content preserved

## What Was Documented

✅ **Template handling sequence** - Template merge happens BEFORE content processing
✅ **Exact processing order** - Placeholders → Blocks → Page URLs → Shortcodes → Translations → Middlewares
✅ **Block rendering flow** - Type-based dispatch with caching
✅ **Placeholder syntax** - `[[Keyword]]` and `[[ Keyword ]]` both supported
✅ **URL pattern matching** - `:any`, `:num`, `:all`, `:string`, `:number`, `:numeric`, `:alpha`
✅ **Caching strategy** - TTL-based with selective invalidation
✅ **Edge cases** - Missing content handling documented

## Files Referenced
- `docs/frontend/overview.md` - Main frontend documentation
- `docs/BLOCK_SYSTEM_ARCHITECTURE.md` - Block system architecture
- `docs/BLOCK_EXTENSIBILITY.md` - Block extensibility guide
- `docs/UNIFIED_BLOCK_TYPES.md` - Unified block type system
- `frontend/frontend.go` - Implementation with inline documentation
- `frontend/block_renderer.go` - Block rendering registry