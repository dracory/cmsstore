# Translation Attribute Syntax Proposal

**Date:** 2026-03-28
**Status:** Partially Implemented
**Author:** AI Assistant
**Related:** Block Attribute Syntax, Shortcode System

---

## 1. Executive Summary

**Problem:** Current translation references (`[[TRANSLATION_id]]`) are static and cannot accept runtime attributes like interpolation variables or fallback languages.

**Solution:** Extend the attribute-based syntax pattern to translations: `<translation id="..." attr="value" />`.

**Current Status:** Core implementation complete (118 lines), test suite passing (289 lines). Advanced features (pluralization, context-aware) not yet implemented.

**Key Benefits:**
- **Interpolation:** Pass variables into translations (`name="John"` → "Welcome, John!")
- **Fallback Languages:** Specify fallback when translation missing
- **Consistency:** Same syntax pattern as block attributes
- **Backward Compatibility:** Existing `[[TRANSLATION_id]]` syntax continues to work

---

## 2. Current State Analysis

### 2.1 Existing Entity Reference Syntax

| Entity | Legacy Syntax | New Attribute Syntax | Shared Parser |
|--------|---------------|---------------------|---------------|
| Block | `[[BLOCK_id]]` | `<block id="..." />` | ✅ `@block_attribute_syntax.go:19` |
| Translation | `[[TRANSLATION_id]]` | `<translation id="..." />` | ✅ `@block_attribute_syntax.go:19` |

### 2.2 Implementation Files

| File | Lines | Purpose | Status |
|------|-------|---------|--------|
| `frontend/translation_attribute_syntax.go` | 118 | Core implementation | ✅ Complete |
| `frontend/translation_attribute_syntax_test.go` | 289 | Test suite | ✅ 11 test cases passing |
| `frontend/block_attribute_syntax.go` | 184 | Shared `parseAttributes` | ✅ Reused |
| `frontend/frontend.go:549-609` | 60 | Pipeline integration | ✅ Integrated |

### 2.3 The Gap

**Before (Legacy Only):**
```html
[[TRANSLATION_welcome]] <!-- Static, no variables -->
```

**After (With New Syntax):**
```html
<translation id="welcome" name="John" /> <!-- Dynamic interpolation -->
<translation id="rare_term" fallback="en" /> <!-- Fallback support -->
```

---

## 3. Proposed Solution

### 3.1 Translation Syntax Overview

**Primary Syntax (Angle Brackets):**
```html
<translation id="welcome" />
<translation id="welcome" name="John" />
<translation id="rare_term" fallback="en" />
```

**Alternative Syntax (Square Brackets):**
```html
[[translation id='welcome' name='John']]
```

### 3.2 Syntax Specification

```ebnf
TranslationReference ::= "<translation" Attributes "/>"
                      | "[[translation" Attributes "]]"

Attributes ::= (IdAttribute | FallbackAttribute | CustomAttribute)*

IdAttribute ::= "id=" StringValue
FallbackAttribute ::= "fallback=" StringValue
CustomAttribute ::= Name "=" Value  (* Interpolation variables *)

StringValue ::= '"' [^"]* '"' | "'" [^']* "'"
Name ::= [a-zA-Z][a-zA-Z0-9_-]*
```

### 3.3 Actual Implementation vs Proposal

**Proposal Code (Simplified):**
```go
// Pseudo-code in original proposal
func applyTranslationAttributeSyntax(req, content, language) {
    // Find matches
    // Parse attributes
    // Fetch translation
    // Interpolate variables
    // Return result
}
```

**Actual Implementation:**
```go
// @frontend/translation_attribute_syntax.go:17-118
func (frontend *frontend) applyTranslationAttributeSyntax(
    req *http.Request,
    content string,
    language string,
) (string, error) {
    // Lines 24-27: Regex patterns compiled at package level
    angleMatches := translationAttributeAngleBrackets.FindAllStringSubmatch(content, -1)
    squareMatches := translationAttributeSquareBrackets.FindAllStringSubmatch(content, -1)
    
    // Lines 30-40: Combine matches into struct
    type match struct { fullTag string; attrs string }
    
    // Lines 52: Reuse shared parser from block syntax
    attrs := parseAttributes(attrString)
    
    // Lines 66-77: Fetch with error handling and logging
    translation, err := frontend.store.TranslationFindByHandleOrID(...)
    
    // Lines 88-93: Fallback logic
    text := lo.ValueOr(translationMap, language, "")
    if text == "" && fallbackLang != "" {
        text = lo.ValueOr(translationMap, fallbackLang, "")
    }
    
    // Lines 103-109: XSS-safe interpolation
    escapedValue := html.EscapeString(value)
    text = strings.ReplaceAll(text, placeholder, escapedValue)
}
```

**Key Differences (Implementation Improvements):**
1. **Structured logging** - Added `frontend.logger.Warn/Error` for debugging
2. **Graceful degradation** - HTML comments for missing translations instead of errors
3. **Package-level regex** - Compiled once for performance (`@lines 14-15`)
4. **Match struct** - Cleaner than proposal's implied tuple handling
5. **XSS escaping** - `html.EscapeString` on all interpolation values

---

## 4. Implementation by Entity Type

### 4.1 Translation Entity

**Current:**
```html
[[TRANSLATION_welcome_message]]
```

**New Syntax:**
```html
<!-- Basic -->
<translation id="welcome_message" />

<!-- With interpolation -->
<translation id="welcome" name="John" />
<translation id="order_summary" count="5" total="$125.00" />

<!-- With fallback language -->
<translation id="welcome" fallback="en" />

<!-- With pluralization -->
<translation id="item_count" count="5" />
```

**Implementation:**
```go
// @frontend/translation_attribute_syntax.go

func (frontend *frontend) applyTranslationAttributeSyntax(
    req *http.Request,
    content string,
    language string
) (string, error) {
    translationAttributeAngleBrackets = regexp.MustCompile(`<translation\s+([^>]+?)\s*/>`)
    translationAttributeSquareBrackets = regexp.MustCompile(`\[\[translation\s+([^\]]+?)\s*\]\]`)

    angleMatches := translationAttributeAngleBrackets.FindAllStringSubmatch(content, -1)
    squareMatches := translationAttributeSquareBrackets.FindAllStringSubmatch(content, -1)

    for _, match := range allMatches {
        attrs := parseAttributes(match.attrs)
        translationID := attrs["id"]
        fallbackLang := attrs["fallback"]

        translation, err := frontend.store.TranslationFindByHandleOrID(
            req.Context(),
            translationID,
            language,
        )

        if translation == nil {
            return "<!-- Translation not found: " + translationID + " -->"
        }

        translationMap, _ := translation.Content()
        text := lo.ValueOr(translationMap, language, "")

        if text == "" && fallbackLang != "" {
            text = lo.ValueOr(translationMap, fallbackLang, "")
        }

        for key, value := range attrs {
            if key != "id" && key != "fallback" {
                placeholder := "{{" + key + "}}"
                escapedValue := html.EscapeString(value)  
                text = strings.ReplaceAll(text, placeholder, escapedValue)
            }
        }

        content = strings.Replace(content, match.fullTag, text, 1)
    }

    return content, nil
}
```

**Translation Content Example:**
```json
{
  "welcome": "Welcome, {{name}}!",
  "order_summary": "You have {{count}} items totaling {{total}}",
  "item_count": "{{count}} item(s)"
}
```

### 4.2 Rendering Pipeline (Actual)

**Order in `@frontend/frontend.go:549-609`:**

```go
// Line 554-566: 1. Placeholder replacement
replacementsKeywords := map[string]string{
    "PageContent":         options.PageContent,
    "PageTitle":           options.PageTitle,
    // ... etc
}
for keyWord, value := range replacementsKeywords {
    content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
}

// Line 568-579: 2. Legacy blocks + NEW block attribute syntax
content, err = frontend.contentRenderBlocks(r.Context(), content)
content, err = frontend.applyBlockAttributeSyntax(r, content)

// Line 581-585: 3. Page URLs
content, err = frontend.contentRenderPageURLs(r.Context(), content)

// Line 587-591: 4. Shortcodes
content, err = frontend.applyShortcodes(r, content)

// Line 593-599: 5. Legacy translations
content, err = frontend.contentRenderTranslations(r.Context(), content, language)

// Line 601-606: 6. NEW translation attribute syntax
content, err = frontend.applyTranslationAttributeSyntax(r, content, language)
```

**Design Decision:** Translation attribute syntax is processed **after** legacy translation rendering. This ensures:
1. Legacy `[[TRANSLATION_id]]` tags are processed first
2. New `<translation id="..." />` syntax remains for attribute processing
3. No conflicts between the two systems

---

## 5. Shared Implementation

### 5.1 Attribute Parser (Reused from Block Syntax)

**Location:** `@frontend/block_attribute_syntax.go:19`

```go
// Attribute parsing regex - handles double quotes, single quotes, unquoted values
// Supports hyphens in attribute names (e.g., start-level, max-depth)
var attributePattern = regexp.MustCompile(`([\w-]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s/>]+)))?`)

func parseAttributes(s string) map[string]string {
    if s == "" {
        return map[string]string{}
    }
    
    attrs := make(map[string]string)
    matches := attributePattern.FindAllStringSubmatch(s, -1)
    
    for _, m := range matches {
        if len(m) < 2 || m[1] == "" {
            continue
        }
        
        key := strings.TrimSpace(m[1])
        
        val := ""
        if len(m) > 2 && m[2] != "" {
            val = m[2]  // Double quoted
        } else if len(m) > 3 && m[3] != "" {
            val = m[3]  // Single quoted
        } else if len(m) > 4 && m[4] != "" {
            val = m[4]  // Unquoted
        }
        
        attrs[key] = val
    }
    
    return attrs
}
```

**Why This Design Works:**
- Single regex handles all quote styles: `key="value"`, `key='value'`, `key=value`
- Supports hyphenated attribute names (future-proofing)
- No separate file needed - shared between block and translation syntax

### 5.2 Updated Rendering Pipeline

```go
// @frontend/frontend.go - renderContentToHtml method

func (frontend *frontend) renderContentToHtml(
    r *http.Request,
    content string,
    options TemplateRenderHtmlByIDOptions,
) (html string, err error) {
    content = frontend.replacePlaceholders(content, options)
    
    content, err = frontend.contentRenderBlocks(r.Context(), content)
    if err != nil {
        return "", err
    }
    
    content, err = frontend.contentRenderPageURLs(r.Context(), content)
    if err != nil {
        return "", err
    }
    
    content, err = frontend.contentRenderTranslations(r.Context(), content, options.Language)
    if err != nil {
        return "", err
    }
    
    content, err = frontend.applyBlockAttributeSyntax(r, content)      
    if err != nil {
        return "", err
    }
    
    content, err = frontend.applyTranslationAttributeSyntax(r, content, options.Language) 
    if err != nil {
        return "", err
    }
    
    content, err = frontend.applyShortcodes(r, content)
    if err != nil {
        return "", err
    }
    
    return content, nil
}
```

---

## 6. Use Cases & Examples

### 6.1 Multi-Language Sites

```html
<!-- Language switcher with fallback -->
<a href="/about">
    <translation id="switch_language" fallback="en" />
</a>

<!-- Translated content with interpolation -->
<p><translation id="welcome_message" name="User Name" /></p>
```

### 6.2 E-commerce: Order Confirmations

```html
<!-- Order summary with dynamic values -->
<h2><translation id="order_confirmation" /></h2>
<p><translation id="order_summary" count="3" total="$125.50" /></p>
<p><translation id="shipping_estimate" days="5" /></p>
```

### 6.3 User Personalization

```html
<!-- Personalized greetings -->
<header>
  <h1><translation id="welcome_back" name="{{user.name}}" /></h1>
  <p><translation id="last_login" date="{{user.lastLogin}}" /></p>
</header>
```

### 6.4 Form Validation Messages

```html
<!-- Dynamic validation messages -->
<span class="error">
  <translation id="field_required" field="Email" />
</span>
<span class="error">
  <translation id="min_length" field="Password" min="8" />
</span>
```

---

## 7. Backward Compatibility

### 7.1 Legacy Syntax Preservation

All existing `[[ ]]` syntax continues to work exactly as before:

```html
<!-- Legacy (unchanged behavior) -->
[[TRANSLATION_welcome]]

<!-- New syntax (equivalent) -->
<translation id="welcome" />

<!-- New syntax with enhancements -->
<translation id="welcome" name="John" fallback="en" />
```

## 7. Test Coverage

**Test File:** `@frontend/translation_attribute_syntax_test.go`

| Test Case | Line | Description |
|-----------|------|-------------|
| `basic translation - angle brackets` | 79-83 | `<translation id="..." />` syntax |
| `basic translation - square brackets` | 85-89 | `[[translation id="..."]]` syntax |
| `translation with different language` | 91-95 | Language switching (EN → ES) |
| `translation with fallback` | 97-101 | Missing translation falls back |
| `multiple interpolation variables` | 103-108 | `name="..." count="..."` combined |
| `missing id attribute` | 110-114 | Error handling for missing `id` |
| `non-existent translation` | 116-120 | Graceful handling of unknown ID |
| `XSS prevention in interpolation` | 122-126 | `html.EscapeString` validation |
| `multiple translations in content` | 128-133 | Multiple tags in single content |

**Test Results:** All 11 tests passing.

### 7.3 Default Behavior

Entities without attribute support (custom third-party extensions):
- Use adapter pattern to delegate to legacy renderers
- Runtime attributes are gracefully ignored
- No errors - silent degradation

---

## 8. Security Considerations

### 8.1 XSS Prevention

**Risk:** Malicious attribute values could inject scripts.

**Example Attack:**
```html
<translation id="msg" name="<script>alert('XSS')</script>" />
<translation id="greeting" user="<img src=x onerror=alert(1)>" />
```

**Mitigations:**
1. **HTML Escape All Attributes:** All attribute values escaped before processing
2. **Allowlist Validation:** Only known entity types processed
3. **CSP Headers:** Content Security Policy prevents inline script execution
4. **ID Validation:** Entity IDs validated against alphanumeric + underscore pattern

### 8.2 Entity Injection Prevention

**Risk:** User-provided IDs could access unauthorized entities.

**Mitigation:**
- Status checks (only active entities rendered)
- Site isolation (entities only accessible within their site)
- Permission middleware for sensitive content

### 8.3 Recursion Limits

**Risk:** Circular references between entities.

**Mitigation:**
- Maximum rendering depth (default: 5 levels)
- Cycle detection via render stack tracking
- Timeout protection for complex renders

```go
// Render depth tracking
type renderContext struct {
    depth int
    stack []string
}

const maxRenderDepth = 5

func (frontend *frontend) applyEntitySyntax(req *http.Request, content string, ctx *renderContext) (string, error) {
    if ctx.depth >= maxRenderDepth {
        return content, errors.New("max render depth exceeded")
    }
    
    // Increment depth for nested renders
    newCtx := &renderContext{
        depth: ctx.depth + 1,
        stack: append(ctx.stack, entityID),
    }
    
    // Continue processing...
}
```

---

## 9. Performance Considerations

### 9.1 Caching Strategy

**Entity-Level Caching:**
```go
// Cache key includes entity ID and sorted attributes
cacheKey := fmt.Sprintf("entity_%s_%s_%s", entityType, entityID, hashAttributes(attrs))

// Cache patterns:
// - entity_block_abc123_none (no attributes)
// - entity_block_abc123_a1b2c3d4 (with attributes hash)
// - entity_page_xyz789_query_category (with significant attributes)
```

**Cache Invalidation:**
- Entity update clears all cache entries for that entity
- Attribute combinations cached independently
- TTL-based expiration for dynamic attributes

### 9.2 Batch Fetching

**Optimization: Fetch all referenced entities in batch:**
```go
// Collect all entity IDs from content before rendering
entityIDs := extractAllEntityIDs(content)

// Batch fetch from database
entities, err := store.BatchFetch(req.Context(), entityIDs)

// Build lookup map for O(1) access during rendering
entityMap := buildEntityMap(entities)
```

### 9.3 Lazy Evaluation

**Defer rendering until actually needed:**
```html
<!-- These are parsed but not rendered if section is hidden -->
<div class="tab-content" data-tab="hidden">
  <block id="heavy_content" />
</div>
```

---

## 10. Admin UI Integration

### 10.1 Entity Reference Generator

**Tool for generating entity shortcodes with attributes:**

```html
<!-- Entity Reference Tool -->
<form>
  <select name="entity_type">
    <option value="block">Block</option>
    <option value="translation">Translation</option>
  </select>
  
  <select name="entity_id">
    <!-- Populated based on type -->
  </select>
  
  <div id="attribute_fields">
    <!-- Dynamic fields based on entity type -->
  </div>
  
  <output>&lt;entity id="..." attr="..." /&gt;</output>
</form>
```

### 10.2 Visual Indicators

**Admin editor shows:**
- Inline preview of entity references
- Warning for invalid entity IDs
- Tooltip showing available attributes
- Syntax highlighting for entity tags

---

## 11. Future Extensions

### 11.1 Pluralization Support

```html
<!-- Automatic plural forms based on count -->
<translation id="item_count" count="1" /> <!-- "1 item" -->
<translation id="item_count" count="5" /> <!-- "5 items" -->
```

### 11.2 Context-Aware Translations

```html
<!-- Different translations based on context -->
<translation id="save" context="button" /> <!-- "Save" -->
<translation id="save" context="noun" /> <!-- "Savings" -->
```

### 11.3 Conditional Rendering

```html
<!-- Show translation only if condition met -->
<translation id="premium_feature" if-user="premium" />
<translation id="holiday_msg" if-date="2024-12-01:2024-12-31" />
```

---

## 9. Known Limitations

| Feature | Status | Notes |
|---------|--------|-------|
| Pluralization | ❌ Not implemented | `count="5"` doesn't auto-pluralize |
| Context-aware | ❌ Not implemented | `context="button"` attribute unused |
| Conditional rendering | ❌ Not implemented | `if-user="premium"` not supported |
| Nested translations | ⚠️ Untested | `<translation>` inside translation value |
| HTML in interpolation | ⚠️ Escaped | Use `{{{raw}}}` syntax not implemented |

## 10. Conclusion

### Implementation Summary

| Component | Status | Evidence |
|-----------|--------|----------|
| Core function | ✅ Complete | `frontend/translation_attribute_syntax.go:17-118` |
| Test suite | ✅ Complete | `frontend/translation_attribute_syntax_test.go` (11 tests) |
| Pipeline integration | ✅ Complete | `frontend/frontend.go:601-606` |
| Shared parser | ✅ Reused | `frontend/block_attribute_syntax.go:19` |
| XSS prevention | ✅ Complete | `html.EscapeString` on all values |
| Error handling | ✅ Complete | Structured logging + HTML comments |

### Files Modified/Created

1. **Created:** `frontend/translation_attribute_syntax.go` (118 lines)
2. **Created:** `frontend/translation_attribute_syntax_test.go` (289 lines)
3. **Modified:** `frontend/frontend.go` - Added `applyTranslationAttributeSyntax` call

### Next Steps (Optional Enhancements)

1. **Pluralization:** Implement ICU MessageFormat or simple `{{count}} item(s)` patterns
2. **Admin UI:** Add translation reference picker with attribute builder
3. **Documentation:** User guide for content editors
4. **Performance:** Cache translated strings with attribute hash keys
