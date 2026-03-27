# Block Attribute Syntax Proposal

**Date:** 2026-03-27  
**Status:** Draft  
**Author:** AI Assistant  
**Related:** Shortcode System, Block System

---

## 1. Executive Summary

**Problem:** Current block references (`[[BLOCK_id]]`) are static and cannot accept runtime attributes. Content editors must create duplicate blocks for minor variations (e.g., "Menu Sidebar" vs "Menu Footer" vs "Menu Mobile").

**Solution:** Introduce attribute-based syntax for blocks: `<block id="..." attr="value" />`. This maintains backward compatibility with `[[BLOCK_id]]` while adding runtime configurability.

**Key Benefits:**
- One block, multiple presentations via runtime attributes
- Clean XML-like syntax with angle brackets (inspired by shortcodes but distinct)
- No breaking changes - both syntaxes coexist
- Simplified content management (fewer block duplicates)

---

## 2. Current State Analysis

### 2.1 Existing Block Reference Syntax

```html
<!-- Current static reference -->
[[BLOCK_block_abc123]]
[[BLOCK_block_xyz789]]
```

**Limitations:**
- Block renders identically everywhere
- No way to pass runtime configuration
- Must create separate blocks for variations:
  - `block_menu_sidebar` (depth=2, style=vertical)
  - `block_menu_footer` (depth=1, style=horizontal)
  - `block_menu_mobile` (depth=1, style=collapsed)

### 2.2 Comparison: Shortcode Syntax

```html
<!-- Current shortcode with attributes -->
<product_list category="electronics" limit="5" sort="price_desc" />
<gallery folder="vacation-2024" layout="masonry" thumbnail-size="300" />
```

**Advantages:**
- Runtime attribute parsing
- Context-aware rendering (access to `*http.Request`)
- Single shortcode, infinite variations

### 2.3 The Gap

| Feature | `[[BLOCK_id]]` | `<shortcode />` |
|---------|----------------|-----------------|
| Stored in DB | ✅ | ❌ |
| Editable via Admin | ✅ | ❌ |
| Runtime attributes | ❌ | ✅ |
| Context-aware | ❌ | ✅ |
| Versioning/history | ✅ | ❌ |

**Goal:** Bring runtime configurability to blocks while maintaining their DB-backed nature.

---

## 3. Proposed Solution

### 3.1 New Syntax Options

Two syntax variants are proposed. **Option A (Generic)** is the recommended approach.

#### Option A: Generic `<block />` Tag (Recommended)

**Primary Syntax (Angle Brackets):**
```html
<!-- Basic reference (backward compatible behavior) -->
<block id="block_abc123" />

<!-- With runtime attributes -->
<block id="block_abc123" depth="2" style="sidebar" highlight-current="true" />

<!-- With multiple attributes -->
<block id="block_abc123" width="300px" show-excerpt="true" />
```

**Alternative Syntax (Square Brackets - for HTML attribute contexts):**
```html
<!-- Use when embedding in HTML attributes to avoid syntax highlighting issues -->
[[block id='menu_main' depth='2']]

<!-- Example: In data attributes -->
<div data-content="[[block id='cta_button']]">

<!-- Example: In JavaScript template strings -->
<script>
  const content = `[[block id='dynamic_content']]`;
</script>
```

**Why Two Syntaxes?**
- **Angle brackets** (`<block />`) - Primary syntax, clean and familiar
- **Square brackets** (`[[block ...]]`) - Alternative for contexts where `<` and `"` cause issues:
  - Inside HTML attributes: `<div data-view="[[block id='x']]">`
  - Inside JavaScript strings: `var x = "[[block id='y']]"`
  - Avoids syntax highlighter confusion
  - Uses single quotes internally to avoid escaping

**Note:** This is a block reference syntax enhancement, not a shortcode. While it uses similar angle bracket notation, it operates on DB-stored blocks rather than code-based shortcodes. The block's type is determined by its stored `Type()` field - no runtime type override is supported.

**Advantages:**
- Single parser implementation
- Self-documenting (type attribute explicit)
- Easy to extend with new attributes
- Works with all block types without code changes

#### Option B: Type-Specific `<block-TYPE />` Tags

```html
<!-- HTML block -->
<block-html id="block_abc123" wrap="div" class="highlight" />

<!-- Menu block -->
<block-menu id="block_xyz789" depth="2" style="vertical" />

<!-- Custom gallery block -->
<block-gallery id="block_gal001" layout="masonry" lightbox="true" />
```

**Advantages:**
- IDE-friendly (type known from tag name)
- Validation at parse time
- Clear separation of concerns

**Disadvantages:**
- Requires registration per type
- Parser complexity (pattern matching or multiple handlers)
- Less flexible (can't dynamically switch types)

### 3.2 Syntax Specification

**EBNF Grammar:**
```ebnf
BlockReference ::= "<block" Attributes "/>"

Attributes ::= (IdAttribute | TypeAttribute | CustomAttribute)*

IdAttribute ::= "id=" StringValue
              (* Required: the block identifier *)

CustomAttribute ::= Name "=" Value
                  (* Runtime attributes passed to renderer *)

StringValue ::= '"' [^"]* '"'
               | "'" [^']* "'"

Name ::= [a-zA-Z][a-zA-Z0-9_-]*

Value ::= StringValue | [^\s>]*
```

**Examples:**
```html
<!-- Angle bracket syntax (primary) -->
<block id="menu_main" />
<block id="menu_main" depth="2" />
<block id="content_hero" bg-color="#f0f0f0" padding="20" />

<!-- Square bracket syntax (alternative for HTML attributes) -->
[[block id='menu_main']]
[[block id='menu_main' depth='2']]
[[block id='content_hero' bg-color='#f0f0f0' padding='20']]
```

---

## 4. Implementation Design

### 4.1 Processing Pipeline

The rendering pipeline will process content in this sequence:

```
1. Placeholder replacement (PageTitle, PageContent, etc.)
   → "[[PageTitle]]" → "My Page"

2. Legacy block references (backward compatibility)
   → "[[BLOCK_id]]" → block content (static)

3. NEW: Block attribute syntax
   → "<block id=... />" or "[[block id=...]]" → block content (dynamic)

4. Other shortcodes
   → "<product_list ... />" → rendered HTML

5. Translation rendering
   → "[[TRANSLATION_id]]" → translated text
```

### 4.2 BlockType Interface Extension

**Breaking Change (Acceptable for Important Feature)**

The `BlockType` interface is extended with variadic options for backward compatibility.

**Current Interface:**
```go
// @cmsstore/block_type.go:68-90
type BlockType interface {
    TypeKey() string
    TypeLabel() string
    Render(ctx context.Context, block BlockInterface) (string, error)
    GetAdminFields(block BlockInterface, r *http.Request) interface{}
    SaveAdminFields(r *http.Request, block BlockInterface) error
}
```

**New Interface (Breaking Change):**
```go
// BlockType defines a complete block type with frontend rendering and admin UI.
type BlockType interface {
    TypeKey() string
    TypeLabel() string
    
    // Render renders the block for frontend display.
    // Options can include runtime attributes via WithAttributes(attrs)
    Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error)
    
    GetAdminFields(block BlockInterface, r *http.Request) interface{}
    SaveAdminFields(r *http.Request, block BlockInterface) error
}

// Note: GetAttributeDefinitions was removed as it's redundant.
// Runtime attributes are implicitly defined by what GetAdminFields() exposes
// and what the Render() method chooses to read from RenderOptions.Attributes.

// RenderOption configures block rendering
type RenderOption func(*RenderOptions)

// RenderOptions holds rendering configuration
type RenderOptions struct {
    Attributes map[string]string
}

// WithAttributes passes runtime attributes to the renderer
func WithAttributes(attrs map[string]string) RenderOption {
    return func(opts *RenderOptions) {
        opts.Attributes = attrs
    }
}
```

**Usage Pattern:**
```go
// Without attributes (backward compatible call)
html, err := blockType.Render(ctx, block)

// With attributes (new syntax)
html, err := blockType.Render(ctx, block, WithAttributes(attrs))
```

**Implementation Example:**
```go
func (t *MenuBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Parse options
    options := &RenderOptions{}
    for _, opt := range opts {
        opt(options)
    }
    
    // Get attributes (empty map if not provided)
    attrs := options.Attributes
    
    // Use attributes or fallback to defaults
    depth := cast.ToInt(attrs["depth"])
    if depth == 0 {
        depth = cast.ToInt(block.Meta("default_depth"))
    }
    if depth == 0 {
        depth = 2 // Block type default
    }
    
    // ... render with depth
}
```

### 4.3 Migration Path for Existing Block Types

**Step 1: Update Signature**
```go
// Before
func (t *HTMLBlockType) Render(ctx context.Context, block BlockInterface) (string, error) {
    return block.Content(), nil
}

// After (add variadic opts)
func (t *HTMLBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    return block.Content(), nil // Attributes ignored for simple blocks
}
```

**Step 2: (Optional) Support Attributes**
```go
func (t *HTMLBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Parse options
    options := &RenderOptions{}
    for _, opt := range opts {
        opt(options)
    }
    
    content := block.Content()
    
    // Apply wrapper if specified
    if wrapper := options.Attributes["wrap"]; wrapper != "" {
        return fmt.Sprintf("<%s>%s</%s>", wrapper, content, wrapper), nil
    }
    
    return content, nil
}

// Attributes are implicitly documented by GetAdminFields()
// No separate GetAttributeDefinitions needed
```

### 4.4 Frontend Rendering Implementation

**New Method: `applyBlockAttributeSyntax()`**

```go
// @frontend/frontend.go (new method)

// Package-level compiled regex for performance
// Matches: <block ...attributes... /> OR [[block ...attributes...]]
// Edge cases handled:
// - Optional whitespace before />
// - Attributes with quotes (single/double)
// - Self-closing tag only (no container blocks)
// - Alternative square bracket syntax for HTML attribute contexts
var blockAttributeAngleBrackets = regexp.MustCompile(`<block\s+([^>]+?)\s*/>`)
var blockAttributeSquareBrackets = regexp.MustCompile(`\[\[block\s+([^\]]+?)\s*\]\]`)

// applyBlockAttributeSyntax processes block references with attributes.
// Supports two syntaxes:
// 1. <block id="..." /> - Primary syntax
// 2. [[block id='...']] - Alternative for HTML attribute contexts
// This is called after legacy [[BLOCK_id]] processing.
func (frontend *frontend) applyBlockAttributeSyntax(req *http.Request, content string) (string, error) {
    // Find all <block ... /> tags (angle bracket syntax)
    angleMatches := blockAttributeAngleBrackets.FindAllStringSubmatch(content, -1)
    
    // Find all [[block ...]] tags (square bracket syntax)
    squareMatches := blockAttributeSquareBrackets.FindAllStringSubmatch(content, -1)
    
    // Combine matches
    allMatches := make([][2]string, 0, len(angleMatches)+len(squareMatches))
    for _, m := range angleMatches {
        allMatches = append(allMatches, [2]string{m[0], m[1]}) // [fullTag, attributes]
    }
    for _, m := range squareMatches {
        allMatches = append(allMatches, [2]string{m[0], m[1]}) // [fullTag, attributes]
    }
    
    if len(allMatches) == 0 {
        return content, nil // No block references found
    }
    
    // Extract all block IDs for batch fetching
    blockIDs := make([]string, 0, len(allMatches))
    for _, match := range allMatches {
        attrs := parseAttributes(match[1])
        if blockID := attrs["id"]; blockID != "" {
            blockIDs = append(blockIDs, blockID)
        }
    }
    
    // Batch fetch all blocks (performance optimization)
    blocks, err := frontend.store.BlockFindByIDs(req.Context(), blockIDs)
    if err != nil {
        return "", fmt.Errorf("batch fetch blocks: %w", err)
    }
    
    // Create lookup map
    blockMap := make(map[string]cmsstore.BlockInterface, len(blocks))
    for _, block := range blocks {
        blockMap[block.ID()] = block
    }
    
    // Process each match
    for _, match := range allMatches {
        fullTag := match[0]      // "<block id=... />" or "[[block id=...]]]"
        attrString := match[1]  // "id=... attr=..."
        
        // Parse attributes
        attrs := parseAttributes(attrString)
        blockID := attrs["id"]
        if blockID == "" {
            continue // Skip invalid tags
        }
        
        // Get block from pre-fetched map
        block, exists := blockMap[blockID]
        if !exists {
            frontend.logger.Warn("Block attribute syntax: not found", "id", blockID)
            content = strings.Replace(content, fullTag, "<!-- Block not found: "+blockID+" -->", 1)
            continue
        }
        
        // Security: Check if block is active
        if !block.IsActive() {
            frontend.logger.Warn("Block attribute syntax: inactive block", "id", blockID)
            content = strings.Replace(content, fullTag, "<!-- Block inactive: "+blockID+" -->", 1)
            continue
        }
        
        // Get block type from stored value
        blockTypeKey := block.Type()
        
        // Get block type (from global registry)
        blockType := cmsstore.GetBlockType(blockTypeKey)
        
        // Remove system attrs and sanitize before passing to renderer
        runtimeAttrs := filterAndSanitizeAttrs(attrs) // remove "id", "type", sanitize values
        
        // Render with attributes
        var html string
        if blockType != nil {
            // Render with attributes (or without if empty)
            // Block types validate attributes internally if needed
            if len(runtimeAttrs) > 0 {
                html, err = blockType.Render(req.Context(), block, cmsstore.WithAttributes(runtimeAttrs))
            } else {
                html, err = blockType.Render(req.Context(), block)
            }
        } else {
            // Fallback to local renderer registry
            renderer := frontend.blockRenderers.GetRenderer(blockTypeKey)
            html, err = renderer.Render(req.Context(), block)
        }
        
        if err != nil {
            frontend.logger.Error("Block attribute syntax: render error", "id", blockID, "error", err)
            html = "<!-- Block render error: " + blockID + " -->"
        }
        
        // Replace in content
        content = strings.Replace(content, fullTag, html, 1)
    }
    
    return content, nil
}

// Package-level compiled regex for attribute parsing (performance)
var attributePattern = regexp.MustCompile(`(\w+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s/>]*)))?`)

// parseAttributes parses "key=value key2='value2' key3=unquoted" into map
// Handles edge cases:
// - Double quotes: key="value with spaces"
// - Single quotes: key='value with spaces'
// - Unquoted: key=value (no spaces allowed)
// - Boolean flags: key (no value, empty string)
func parseAttributes(s string) map[string]string {
    if s == "" {
        return map[string]string{}
    }
    
    attrs := make(map[string]string)
    matches := attributePattern.FindAllStringSubmatch(s, -1)
    
    for _, m := range matches {
        if len(m) < 2 || m[1] == "" {
            continue // Skip invalid matches
        }
        
        key := strings.TrimSpace(m[1])
        
        // Value could be in group 2 (double quotes), 3 (single quotes), or 4 (unquoted)
        val := ""
        if len(m) > 2 && m[2] != "" {
            val = m[2] // Double quoted
        } else if len(m) > 3 && m[3] != "" {
            val = m[3] // Single quoted
        } else if len(m) > 4 && m[4] != "" {
            val = m[4] // Unquoted
        }
        // If no value found, it's a boolean flag (empty string)
        
        attrs[key] = val
    }
    
    return attrs
}

// filterAndSanitizeAttrs removes system-reserved attributes and sanitizes values
func filterAndSanitizeAttrs(attrs map[string]string) map[string]string {
    systemAttrs := map[string]bool{"id": true} // Only 'id' is system-reserved
    filtered := make(map[string]string)
    for k, v := range attrs {
        if !systemAttrs[k] {
            // Security: HTML escape attribute values to prevent XSS
            filtered[k] = html.EscapeString(v)
        }
    }
    return filtered
}

// hashAttributes creates a stable hash of attributes for caching
func hashAttributes(attrs map[string]string) string {
    if len(attrs) == 0 {
        return "none"
    }
    
    // Sort keys for consistent hashing
    keys := make([]string, 0, len(attrs))
    for k := range attrs {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    
    // Build sorted key=value string
    var parts []string
    for _, k := range keys {
        parts = append(parts, fmt.Sprintf("%s=%s", k, attrs[k]))
    }
    
    // Hash the result
    h := sha256.Sum256([]byte(strings.Join(parts, "&")))
    return hex.EncodeToString(h[:8]) // Use first 8 bytes for shorter keys
}
```

### 4.5 Updated Rendering Pipeline

```go
// @frontend/frontend.go - renderContentToHtml method

func (frontend *frontend) renderContentToHtml(
    r *http.Request,
    content string,
    options TemplateRenderHtmlByIDOptions,
) (html string, err error) {
    // 1. Placeholder replacement
    content = frontend.replacePlaceholders(content, options)
    
    // 2. Legacy block references (maintain backward compatibility)
    content, err = frontend.contentRenderBlocks(r.Context(), content)
    if err != nil {
        return "", err
    }
    
    // 3. NEW: Block attribute syntax
    content, err = frontend.applyBlockAttributeSyntax(r, content)
    if err != nil {
        return "", err
    }
    
    // 4. Page URL placeholders
    content, err = frontend.contentRenderPageURLs(r.Context(), content)
    if err != nil {
        return "", err
    }
    
    // 5. Other shortcodes (product_list, gallery, etc.)
    content, err = frontend.applyShortcodes(r, content)
    if err != nil {
        return "", err
    }
    
    // 6. Translations
    content, err = frontend.contentRenderTranslations(r.Context(), content, options.Language)
    if err != nil {
        return "", err
    }
    
    return content, nil
}
```

---

## 5. Block Type Implementation Examples

### 5.1 Menu Block with Runtime Attributes

```go
// @blocks/menu/menu_block_type.go

type MenuBlockType struct {
    store StoreInterface
    logger *slog.Logger
}

func (t *MenuBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    // No runtime attributes - use stored defaults
    return t.renderWithAttributes(ctx, block, nil)
}

func (t *MenuBlockType) RenderWithAttributes(ctx context.Context, block cmsstore.BlockInterface, attrs map[string]string) (string, error) {
    // Get menu ID from block content or metas
    menuID := block.Content()
    
    // Runtime overrides (if provided via <block id="..." depth="2" />)
    depth := cast.ToInt(attrs["depth"])
    if depth == 0 {
        depth = cast.ToInt(block.Meta("default_depth")) // Fallback to stored
    }
    
    style := attrs["style"]
    if style == "" {
        style = block.Meta("default_style") // Fallback to stored
    }
    
    highlightCurrent := attrs["highlight-current"] == "true"
    
    // Fetch menu items
    menuItems, err := t.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
        SetMenuID(menuID).
        SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE))
    if err != nil {
        return "", err
    }
    
    // Render with runtime configuration
    renderer := NewMenuRenderer(t)
    return renderer.RenderMenuHTML(ctx, menuItems, style, "", 0, depth)
}

func (t *MenuBlockType) ValidateAttributes(attrs map[string]string) error {
    validAttrs := []string{"depth", "style", "highlight-current", "css-class"}
    validSet := make(map[string]bool)
    for _, a := range validAttrs {
        validSet[a] = true
    }
    
    for k := range attrs {
        if !validSet[k] {
            return fmt.Errorf("invalid attribute: %s", k)
        }
    }
    return nil
}

func (t *MenuBlockType) GetAttributeDefinitions() []cmsstore.BlockAttributeDefinition {
    return []cmsstore.BlockAttributeDefinition{
        {
            Name:        "depth",
            Type:        "int",
            Required:    false,
            Default:     2,
            Description: "Maximum nesting depth for menu items",
            Validation:  "range:1,10",
        },
        {
            Name:        "style",
            Type:        "enum",
            Required:    false,
            Default:     "vertical",
            Description: "Menu presentation style",
            EnumValues:  []string{"vertical", "horizontal", "dropdown", "sidebar"},
        },
        {
            Name:        "highlight-current",
            Type:        "bool",
            Required:    false,
            Default:     false,
            Description: "Highlight current page in menu",
        },
    }
}
```

**Usage:**
```html
<!-- One menu block, three presentations -->
<block id="menu_main" depth="2" style="sidebar" />
<block id="menu_main" depth="1" style="horizontal" />
<block id="menu_main" depth="3" style="dropdown" highlight-current="true" />
```

### 5.2 HTML Block with Wrapper Attributes

```go
// @blocks/html/html_block_type.go

func (t *HTMLBlockType) RenderWithAttributes(ctx context.Context, block cmsstore.BlockInterface, attrs map[string]string) (string, error) {
    content := block.Content()
    
    // Runtime wrapper configuration
    wrapper := attrs["wrap"] // div, section, article, etc.
    cssClass := attrs["class"]
    cssId := attrs["css-id"] // 'id' is reserved, so use 'css-id'
    
    if wrapper != "" {
        html := fmt.Sprintf("<%s", wrapper)
        if cssClass != "" {
            html += fmt.Sprintf(` class="%s"`, cssClass)
        }
        if cssId != "" {
            html += fmt.Sprintf(` id="%s"`, cssId)
        }
        html += ">" + content + fmt.Sprintf("</%s>", wrapper)
        return html, nil
    }
    
    // No wrapper - return raw content
    return content, nil
}
```

**Usage:**
```html
<!-- Same HTML block, different wrappers -->
<block id="cta_content" wrap="section" class="hero-section" css-id="main-cta" />
<block id="cta_content" wrap="aside" class="sidebar-cta" />
<block id="cta_content" /> <!-- Raw output, no wrapper -->
```

---

## 6. Backward Compatibility

### 6.1 Legacy Syntax Preservation

The existing `[[BLOCK_id]]` syntax continues to work exactly as before:

```html
<!-- Legacy (unchanged behavior) -->
[[BLOCK_block_abc123]]

<!-- New syntax (equivalent) -->
<block id="block_abc123" />

<!-- New syntax with enhancements -->
<block id="block_abc123" depth="2" />
```

### 6.2 Migration Path

**Phase 1: Implementation (No breaking changes)**
1. Add `RenderWithAttributes` to `BlockType` interface
2. Implement in all built-in block types
3. Add `applyBlockShortcodes` to frontend
4. Legacy blocks unaffected

**Phase 2: Gradual Adoption**
1. Content editors can optionally use new syntax
2. Documentation updated with examples
3. Admin UI shows available attributes per block type

**Phase 3: Deprecation (Future)**
- `[[BLOCK_id]]` deprecated (not removed) in favor of `<block id="..." />`
- Migration tool provided to convert existing references

### 6.3 Default Behavior

Block types without `RenderWithAttributes` implementation:
- Use adapter pattern to delegate to `Render()`
- Runtime attributes are ignored
- No error - graceful degradation

---

## 7. Use Cases & Examples

### 7.1 E-commerce: Product Lists

```html
<!-- One "Featured Products" block, multiple contexts -->
<block id="products_featured" category="electronics" limit="4" />
<block id="products_featured" category="clothing" limit="8" layout="grid" />
<block id="products_featured" category="sale" limit="6" sort="discount_desc" />
```

### 7.2 Navigation: Multi-context Menus

```html
<!-- Main menu in different contexts -->
<block id="menu_main" depth="1" style="horizontal" />
<block id="menu_main" depth="2" style="sidebar" highlight-current="true" />
<block id="menu_main" depth="1" style="collapsed" /> <!-- Mobile -->
```

### 7.3 Content: Reusable CTAs

```html
<!-- Call-to-action block with contextual styling -->
<block id="cta_subscribe" wrap="section" class="subscribe-hero" />
<block id="cta_subscribe" wrap="aside" class="subscribe-sidebar" />
<block id="cta_subscribe" wrap="div" class="subscribe-inline" />
```

### 7.4 Media: Dynamic Galleries

```html
<!-- Gallery with runtime configuration -->
<block id="gallery_portfolio" folder="web-design" layout="masonry" lightbox="true" />
<block id="gallery_portfolio" folder="photography" layout="slider" autoplay="5" />
```

### 7.5 Data: Dynamic Tables

```html
<!-- Data table with filters -->
<block id="table_orders" status="pending" sort="date_desc" limit="20" />
<block id="table_orders" status="completed" sort="total_desc" limit="10" />
```

---

## 8. Security Considerations

### 8.1 XSS Prevention

**Risk:** Malicious attribute values could inject scripts.

**Example Attack:**
```html
<block id="menu" class="<script>alert('XSS')</script>" />
<block id="content" onclick="malicious()" />
```

**Mitigations:**
1. **HTML Escape All Attributes:** All attribute values are escaped before passing to renderers
2. **Attribute Validation:** Block types validate attribute names and values
3. **CSP Headers:** Content Security Policy prevents inline script execution
4. **Allowlist Approach:** Only known attributes are processed

**Implementation:**
```go
func filterAndSanitizeAttrs(attrs map[string]string) map[string]string {
    // HTML escape all values
    for k, v := range attrs {
        attrs[k] = html.EscapeString(v)
    }
    return attrs
}
```

### 8.2 SQL Injection Prevention

**Risk:** Block IDs from user input used in queries.

**Mitigation:**
- Use parameterized queries (already implemented in store)
- Validate block ID format (alphanumeric + underscore only)

### 8.3 Denial of Service

**Risk:** Deeply nested or recursive block references.

**Example:**
```html
<!-- Block A references Block B, Block B references Block A -->
<block id="block_a" /> <!-- Contains <block id="block_b" /> -->
```

**Mitigation:**
- Maximum recursion depth limit (default: 3 levels)
- Cycle detection
- Timeout for block rendering

---

## 9. Attribute Resolution Hierarchy

### 9.1 Resolution Order

When a block is rendered with `<block id="..." attr="value" />`, attributes are resolved in this priority order:

1. **Runtime attributes** (highest priority) - from shortcode syntax
2. **Block metas** - stored in block's `metas` field
3. **Site-wide defaults** - from site configuration
4. **Block type defaults** - from `BlockAttributeDefinition.Default`

**Example:**
```go
func (t *MenuBlockType) RenderWithAttributes(ctx context.Context, block BlockInterface, attrs map[string]string) (string, error) {
    // 1. Runtime attribute (highest priority)
    depth := attrs["depth"]
    
    // 2. Fallback to block meta
    if depth == "" {
        depth = block.Meta("default_depth")
    }
    
    // 3. Fallback to site-wide default (from context or config)
    if depth == "" {
        if site := SiteFromContext(ctx); site != nil {
            depth = site.Meta("menu_default_depth")
        }
    }
    
    // 4. Fallback to block type default
    if depth == "" {
        depth = "2" // Block type default
    }
    
    depthInt := cast.ToInt(depth)
    // ... render with depthInt
}
```

### 9.2 Attribute Inheritance

Blocks can inherit attributes from parent blocks or templates:

```html
<!-- Template sets default style -->
<div data-menu-style="sidebar">
    <!-- Block inherits style="sidebar" unless overridden -->
    <block id="menu_main" />
    
    <!-- Override inherited style -->
    <block id="menu_main" style="horizontal" />
</div>
```

**Implementation:** Use context to pass inherited attributes down the rendering chain.

---

## 10. Nested Blocks and Recursion

### 10.1 Nested Block Support

Blocks can contain other block references:

```html
<block id="wrapper_section">
    <h2>Navigation</h2>
    <block id="menu_main" depth="2" />
    <block id="menu_footer" depth="1" />
</block>
```

### 10.2 Recursion Limits

**Problem:** Infinite recursion if blocks reference each other.

**Solution:** Track rendering depth in context:

```go
type blockRenderDepthKey struct{}

func (frontend *frontend) applyBlockAttributeSyntax(req *http.Request, content string) (string, error) {
    // Get current depth from context
    depth := 0
    if d, ok := req.Context().Value(blockRenderDepthKey{}).(int); ok {
        depth = d
    }
    
    // Check maximum depth
    const maxDepth = 3
    if depth >= maxDepth {
        frontend.logger.Warn("Block attribute syntax: max recursion depth reached", "depth", depth)
        return content, nil // Stop processing
    }
    
    // Increment depth for nested renders
    ctx := context.WithValue(req.Context(), blockRenderDepthKey{}, depth+1)
    req = req.WithContext(ctx)
    
    // ... continue processing
}
```

### 10.3 Cycle Detection

Detect circular references:

```go
type blockRenderStackKey struct{}

func (frontend *frontend) applyBlockAttributeSyntax(req *http.Request, content string) (string, error) {
    // Get current render stack
    stack := []string{}
    if s, ok := req.Context().Value(blockRenderStackKey{}).([]string); ok {
        stack = s
    }
    
    // Check for cycles
    for _, match := range matches {
        blockID := parseAttributes(match[1])["id"]
        
        // Check if already in stack
        for _, id := range stack {
            if id == blockID {
                return "", fmt.Errorf("circular block reference detected: %s", blockID)
            }
        }
        
        // Add to stack
        newStack := append(stack, blockID)
        ctx := context.WithValue(req.Context(), blockRenderStackKey{}, newStack)
        // ... render with new context
    }
}
```

---

## 11. Admin UI Integration

### 8.1 Attribute Documentation

Each block type exposes its supported attributes:

```go
// Block edit page shows:
// "This block supports runtime attributes:"
// - depth (int): Maximum nesting depth
// - style (enum): vertical | horizontal | dropdown
// - highlight-current (bool): Highlight current page
```

### 8.2 Reference Generator

Admin UI tool for generating block shortcodes:

```html
<!-- Block: Menu Main -->
<form>
  <select name="style">
    <option value="vertical">Vertical</option>
    <option value="horizontal">Horizontal</option>
  </select>
  <input type="number" name="depth" value="2" min="1" max="5" />
  <label><input type="checkbox" name="highlight-current" /> Highlight Current</label>
  
  <!-- Generated reference -->
  <output>&lt;block id="menu_main" style="vertical" depth="2" highlight-current="true" /&gt;</output>
</form>
```

---

## 12. Performance Considerations

### 9.1 Caching Strategy

Block shortcodes should be cache-aware:

```go
func (frontend *frontend) applyBlockAttributeSyntax(req *http.Request, content string) (string, error) {
    // Cache key includes block ID + sorted attributes
    // <block id="menu_main" depth="2" style="sidebar" />
    // Cache key: "block_attr:menu_main:depth=2:style=sidebar"
    
    // If cached, return immediately
    // If not, render and cache with appropriate TTL
}
```

### 9.2 Attribute Normalization

For consistent caching:

```go
// Normalize attribute order: sort keys alphabetically
// <block id="x" z="1" a="2" /> → cache key uses "a=2:z=1"
```

### 9.3 Database Query Optimization

```go
// Batch fetch blocks for all shortcodes in content
// Single query: SELECT * FROM blocks WHERE id IN (?, ?, ?)
// Instead of N queries
```

---

## 13. Error Handling

### 10.1 Graceful Degradation

| Error Condition | Behavior |
|-----------------|----------|
| Block not found | HTML comment: `<!-- Block not found: block_abc123 -->` |
| Invalid attributes | Log warning, ignore invalid attrs, render with valid ones |
| Type mismatch | Log error, render with stored type |
| Render failure | HTML comment: `<!-- Block render error: block_abc123 -->` |

### 10.2 Validation Levels

1. **Parse-time:** Valid XML-like syntax
2. **Fetch-time:** Block exists in DB
3. **Render-time:** Attributes valid for block type
4. **Output-time:** HTML generation succeeds

---

## 14. Edge Cases and Parsing Considerations

### 14.1 Attribute Parsing Edge Cases

Based on learnings from the shortcode package implementation, the following edge cases must be handled:

**1. Quote Variations**
```html
<!-- Double quotes -->
<block id="menu_main" style="sidebar" />

<!-- Single quotes -->
<block id='menu_main' style='sidebar' />

<!-- Mixed quotes -->
<block id="menu_main" style='sidebar' />

<!-- Unquoted (no spaces) -->
<block id=menu_main depth=2 />
```

**2. Boolean Flags (No Value)**
```html
<!-- Flag without value = empty string -->
<block id="menu" featured />
<!-- Equivalent to: featured="" -->
```

**3. Whitespace Handling**
```html
<!-- Extra spaces around = -->
<block id = "menu" style = "sidebar" />

<!-- Trailing whitespace before /> -->
<block id="menu"    />

<!-- No space before /> -->
<block id="menu"/>
```

**4. Special Characters in Values**
```html
<!-- URLs with special chars -->
<block id="image" src="https://example.com/img?size=large&format=webp" />

<!-- HTML entities (will be escaped) -->
<block id="content" class="<script>alert('xss')</script>" />
```

**5. Empty Attributes**
```html
<!-- Empty string value -->
<block id="menu" title="" />

<!-- No attributes except id -->
<block id="menu" />
```

**6. Invalid Syntax (Should Be Ignored)**
```html
<!-- Missing id -->
<block style="sidebar" /> <!-- Ignored: no id -->

<!-- Malformed quotes -->
<block id="menu style="sidebar" /> <!-- Parsing error -->

<!-- Not self-closing -->
<block id="menu"> <!-- Not matched by regex -->
```

### 14.2 Content Matching Edge Cases

**1. Nested Angle Brackets**
```html
<!-- Block containing HTML -->
<block id="wrapper" />
<!-- If wrapper contains: <div><block id="inner" /></div> -->
<!-- Inner block should be processed recursively -->
```

**2. Multiple Blocks on Same Line**
```html
<block id="menu" /> <block id="footer" />
```

**3. Blocks in Different Contexts**
```html
<!-- In HTML attributes - USE SQUARE BRACKET SYNTAX -->
<div data-content="[[block id='menu']]">
<!-- Angle brackets would break syntax highlighting: -->
<!-- <div data-content="<block id='menu' />"> ❌ -->

<!-- In JavaScript strings - USE SQUARE BRACKET SYNTAX -->
<script>var x = "[[block id='menu']]";</script>
<!-- Angle brackets would break syntax highlighting: -->
<!-- <script>var x = "<block id='menu' />";</script> ❌ -->

<!-- In HTML comments (should NOT be processed) -->
<!-- <block id="menu" /> -->
<!-- [[block id='menu']] -->
```

**Implementation:** Both syntaxes are processed. Square bracket syntax specifically designed for HTML attribute and string contexts.

### 14.3 Regex Pattern Improvements

**Current Patterns:**
```go
// Angle bracket syntax (primary)
var blockAttributeAngleBrackets = regexp.MustCompile(`<block\s+([^>]+?)\s*/>`)

// Square bracket syntax (alternative for HTML attributes)
var blockAttributeSquareBrackets = regexp.MustCompile(`\[\[block\s+([^\]]+?)\s*\]\]`)
```

**Handles:**
- Self-closing tags only (both syntaxes)
- Optional whitespace before `/>` or `]]`
- Non-greedy attribute capture
- Two syntax variants for different contexts

**Does NOT Handle:**
- Container blocks: `<block id="x">content</block>` or `[[block id='x']]content[[/block]]`
- Context-aware parsing (inside HTML comments)
- Escaped brackets: `\<block\>` or `\[\[block\]\]`

**Recommendation:** Current pattern is sufficient for initial implementation. Add container block support in Phase 2 if needed.

### 14.4 Sequential Replacement Considerations

The implementation uses sequential replacement (one at a time):
```go
content = strings.Replace(content, fullTag, html, 1)
```

**Why Sequential:**
- Allows nested blocks to be processed in subsequent passes
- Prevents offset issues when replacing multiple matches
- Matches shortcode package behavior

**Trade-off:** Multiple passes through content for nested blocks (acceptable for max depth of 3).

---

## 15. Testing Strategy

### 15.1 Unit Tests

**Attribute Parser Tests:**
```go
func TestParseAttributes(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected map[string]string
    }{
        {
            name:  "double quotes",
            input: `id="abc" style="sidebar"`,
            expected: map[string]string{"id": "abc", "style": "sidebar"},
        },
        {
            name:  "single quotes",
            input: `id='abc' style='sidebar'`,
            expected: map[string]string{"id": "abc", "style": "sidebar"},
        },
        {
            name:  "mixed quotes",
            input: `id="abc" style='sidebar'`,
            expected: map[string]string{"id": "abc", "style": "sidebar"},
        },
        {
            name:  "unquoted values",
            input: `id=abc depth=2`,
            expected: map[string]string{"id": "abc", "depth": "2"},
        },
        {
            name:  "boolean flag",
            input: `id="menu" featured`,
            expected: map[string]string{"id": "menu", "featured": ""},
        },
        {
            name:  "extra whitespace",
            input: `id = "abc"  style = "sidebar"`,
            expected: map[string]string{"id": "abc", "style": "sidebar"},
        },
        {
            name:  "empty value",
            input: `id="menu" title=""`,
            expected: map[string]string{"id": "menu", "title": ""},
        },
        {
            name:  "special chars in value",
            input: `id="img" src="https://example.com/img?size=large&format=webp"`,
            expected: map[string]string{"id": "img", "src": "https://example.com/img?size=large&format=webp"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseAttributes(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestBlockAttributePatternMatching(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        shouldMatch bool
        expectedID  string
    }{
        {
            name:      "angle bracket - basic self-closing",
            input:     `<block id="abc" />`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "square bracket - basic",
            input:     `[[block id='abc']]`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "square bracket - with attributes",
            input:     `[[block id='abc' depth='2' style='sidebar']]`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "no space before />",
            input:     `<block id="abc"/>`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "extra whitespace",
            input:     `<block id="abc"    />`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "multiple attributes",
            input:     `<block id="abc" depth="2" style="sidebar" />`,
            shouldMatch: true,
            expectedID:  "abc",
        },
        {
            name:      "missing id (should still match pattern)",
            input:     `<block style="sidebar" />`,
            shouldMatch: true,
            expectedID:  "", // Will be filtered out in processing
        },
        {
            name:      "not self-closing",
            input:     `<block id="abc">`,
            shouldMatch: false,
        },
        {
            name:      "container block",
            input:     `<block id="abc">content</block>`,
            shouldMatch: false, // Not supported in v1
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            matches := blockAttributePattern.FindAllStringSubmatch(tt.input, -1)
            if tt.shouldMatch {
                assert.NotEmpty(t, matches)
                if tt.expectedID != "" {
                    attrs := parseAttributes(matches[0][1])
                    assert.Equal(t, tt.expectedID, attrs["id"])
                }
            } else {
                assert.Empty(t, matches)
            }
        })
    }
}

func TestMenuBlockWithAttributes(t *testing.T) {
    block := cmsstore.NewBlock()
    block.SetContent("menu_main")
    
    attrs := map[string]string{"depth": "3", "style": "sidebar"}
    
    html, err := menuBlockType.RenderWithAttributes(ctx, block, attrs)
    
    assert.NoError(t, err)
    assert.Contains(t, html, `class="sidebar"`)
}
```

### 15.2 Integration Tests

```go
func TestFrontendBlockAttributeSyntaxPipeline(t *testing.T) {
    content := `<block id="test_menu" depth="2" style="vertical" />`
    
    result := frontend.PageRenderHtmlBySiteAndAlias(..., content, ...)
    
    assert.Contains(t, result, `<ul class="menu-vertical">`)
    assert.NotContains(t, result, `<block id=`) // Fully rendered
}
```

### 15.3 Security Tests

```go
func TestAttributeSanitization(t *testing.T) {
    tests := []struct {
        name     string
        input    map[string]string
        expected map[string]string
    }{
        {
            name: "XSS attempt in attribute",
            input: map[string]string{
                "id":    "menu",
                "class": "<script>alert('xss')</script>",
            },
            expected: map[string]string{
                "class": "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
            },
        },
        {
            name: "HTML entities",
            input: map[string]string{
                "id":    "content",
                "title": "A & B > C < D",
            },
            expected: map[string]string{
                "title": "A &amp; B &gt; C &lt; D",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := filterAndSanitizeAttrs(tt.input)
            for k, expectedVal := range tt.expected {
                assert.Equal(t, expectedVal, result[k])
            }
        })
    }
}

func TestRecursionLimits(t *testing.T) {
    // Create blocks that reference each other
    blockA := cmsstore.NewBlock().SetID("block_a").SetContent(`<block id="block_b" />`)
    blockB := cmsstore.NewBlock().SetID("block_b").SetContent(`<block id="block_a" />`)
    
    // Store blocks
    store.BlockCreate(ctx, blockA)
    store.BlockCreate(ctx, blockB)
    
    // Attempt to render - should hit depth limit
    content := `<block id="block_a" />`
    result := frontend.applyBlockAttributeSyntax(req, content)
    
    // Should not hang, should log warning about max depth
    assert.NotContains(t, result, "<block") // All blocks processed or stopped
}
```

---

## 16. Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Breaking existing blocks | High | Low | Optional interface pattern; no breaking changes |
| Performance regression | Medium | Medium | Batch queries, caching, compiled regex |
| XSS via attributes | High | Medium | HTML escape all values, attribute validation |
| Infinite recursion | Medium | Low | Depth limits, cycle detection |
| Cache key collisions | Low | Low | Hash-based keys with sorted attributes |
| Over-complexity | Medium | Medium | Optional adoption, clear documentation, examples |

---

## 17. Implementation Phases

### Phase 1: Core Infrastructure ✅ COMPLETED
- [X] Extend `BlockType` interface with variadic `opts ...RenderOption`
- [X] Create `RenderOption` and `WithAttributes()` helper
- [X] Remove `GetAttributeDefinitions()` - redundant with `GetAdminFields()`
- [X] Implement `applyBlockAttributeSyntax` in frontend
- [X] Support both angle bracket and square bracket syntaxes
- [X] Add attribute parsing utilities with HTML escaping
- [ ] Add batch query support: `BlockFindByIDs` (optional performance optimization)

### Phase 2: Update Built-in Block Types ✅ COMPLETED
- [X] Update all block type signatures: `Render(ctx, block, opts ...RenderOption)`
- [X] Implement attribute support in `HTMLBlockType` (wrap)
- [X] Implement attribute support in `MenuBlockType` (depth, style, mode, class, id, start-level)
- [X] Implement attribute support in `NavbarBlockType` (style, mode, class, id)
- [X] Implement attribute support in `BreadcrumbsBlockType` (style, mode, class, id, separator)
- [X] Add comprehensive tests for attribute parsing and rendering

### Phase 3: Admin UI ✅ COMPLETED
- [X] Display block attribute syntax on block edit page
- [ ] Create block reference generator tool (optional enhancement)
- [ ] Add validation UI for runtime attributes (optional enhancement)

### Phase 4: Documentation & Examples (3-4 days)
- [ ] Update developer documentation
- [ ] Create use-case examples
- [ ] Migration guide from `[[BLOCK_id]]` to `<block id="..." />`

### Phase 5: Optimization (1 week)
- [ ] Implement caching with hash-based keys
- [ ] Add recursion depth tracking and cycle detection
- [ ] Performance benchmarking and profiling

**Total Timeline:** 4-5 weeks

---

## 18. Decision Log

| Decision | Rationale |
|----------|-----------|
| Generic `<block />` over `<block-TYPE />` | Single parser, simpler registration, cleaner syntax |
| Two syntax variants (angle + square brackets) | Square brackets avoid syntax highlighting issues in HTML attributes and JS strings |
| No type override support | Block type is determined by DB value; runtime override adds complexity and security risks |
| Variadic options pattern for Render() | Backward compatible calls without attributes, clean extension point |
| Breaking change to BlockType interface | Acceptable for important feature; variadic args minimize migration effort |
| No GetAttributeDefinitions() method | Redundant with GetAdminFields(); attributes implicitly defined by what block types read |
| Additive change (no removal of `[[BLOCK_id]]`) | Legacy syntax preserved, new syntax adds capabilities |

---

## 19. Open Questions

### 18.1 Attribute Namespacing

**Question:** Should we enforce attribute namespacing to avoid collisions?

**Options:**
- `data-*` prefix for custom attributes (HTML5 convention)
- `block-*` prefix for system attributes
- No prefix (current proposal)

**Recommendation:** Start without prefixes, add if collisions become an issue.

### 18.2 Attribute Type Coercion

**Question:** Should attributes be automatically type-coerced or remain strings?

**Current:** All attributes are strings, block types cast as needed

**Alternative:** Parse and type-coerce during attribute parsing:
```go
attrs := map[string]interface{}{
    "depth": 2,           // int
    "enabled": true,      // bool
    "style": "sidebar",   // string
}
```

**Recommendation:** Keep as strings for simplicity. Block types handle coercion.

### 18.3 Attribute Validation Timing

**Question:** When should attributes be validated?

**Options:**
1. Parse-time (before rendering)
2. Render-time (during block type render)
3. Both

**Current proposal:** Render-time validation via `ValidateAttributes()`

**Consideration:** Parse-time validation could fail fast but requires attribute definitions to be available globally.

### 18.4 Caching Granularity

**Question:** Should we cache at page level or block level?

**Options:**
- **Page-level:** Cache entire rendered page (current behavior)
- **Block-level:** Cache individual block renders with attributes
- **Hybrid:** Cache blocks, compose into pages

**Recommendation:** Start with block-level caching, measure performance impact.

---

## 20. Conclusion

This proposal bridges the gap between static block references and dynamic shortcodes, providing:

- **Content editors:** Fewer duplicate blocks, more flexible presentations
- **Developers:** Clean interface extension, backward compatibility
- **Performance:** Caching support, optimized queries
- **Future-proof:** Foundation for more advanced block features

The syntax `<block id="..." attr="value" />` is intuitive, uses familiar angle bracket notation (inspired by but distinct from shortcodes), and provides immediate value while maintaining full backward compatibility.

---

## Appendices

### Appendix A: Full Interface Definition

```go
package cmsstore

// BlockType defines a complete block type with frontend rendering and admin UI.
// Extended to support runtime attributes via shortcode syntax.
type BlockType interface {
    // TypeKey returns the unique identifier for this block type.
    TypeKey() string

    // TypeLabel returns the human-readable display name.
    TypeLabel() string

    // Render renders the block for frontend display.
    // Called for legacy [[BLOCK_id]] references.
    Render(ctx context.Context, block BlockInterface) (string, error)

    // GetAdminFields returns form fields for the admin content editing tab.
    GetAdminFields(block BlockInterface, r *http.Request) interface{}

    // SaveAdminFields processes form submission and updates the block.
    SaveAdminFields(r *http.Request, block BlockInterface) error

    // === ATTRIBUTE SUPPORT (NEW) ===

    // RenderWithAttributes renders the block with runtime attributes.
    // Called for <block id="..." attr="value" /> shortcode references.
    // Implementations should delegate to Render() if attrs is nil or empty.
    RenderWithAttributes(ctx context.Context, block BlockInterface, attrs map[string]string) (string, error)

    // ValidateAttributes checks if provided attributes are valid.
    // Return nil to accept all attributes (permissive).
    ValidateAttributes(attrs map[string]string) error

    // GetAttributeDefinitions returns metadata about supported attributes.
    // Used for documentation, admin UI, and validation.
    GetAttributeDefinitions() []BlockAttributeDefinition
}

// BlockAttributeDefinition describes a single runtime attribute.
type BlockAttributeDefinition struct {
    Name        string      // Attribute name
    Type        string      // "string", "int", "bool", "enum", "json"
    Required    bool        // Whether attribute is required
    Default     interface{} // Default value
    Description string      // Human-readable description
    EnumValues  []string    // Valid values for enum type
    Validation  string      // Validation rule (regex, range, etc.)
}

// BaseBlockType provides default implementations for backward compatibility.
type BaseBlockType struct{}

func (b *BaseBlockType) RenderWithAttributes(ctx context.Context, block BlockInterface, attrs map[string]string) (string, error) {
    return b.Render(ctx, block)
}

func (b *BaseBlockType) ValidateAttributes(attrs map[string]string) error {
    return nil // Permissive by default
}

func (b *BaseBlockType) GetAttributeDefinitions() []BlockAttributeDefinition {
    return nil
}
```

### Appendix B: Alternative Syntaxes Considered

| Syntax | Pros | Cons | Status |
|--------|------|------|--------|
| `<block id="..." />` | Clean, self-closing, XML-like | None | **Selected** |
| `[[BLOCK_id attr="val"]]` | Consistent with existing | Parser complexity, less intuitive | Rejected |
| `{block: id="..."}` | JSON-like | Unfamiliar, harder to parse | Rejected |
| `<cms:block id="..." />` | Namespaced | Verbose, unnecessary | Rejected |
| `<block-TYPE id="..." />` | Type clear from tag | Requires type registration, less flexible | Alternative |
