# PAGE_URL Placeholder System

The CMS store supports a `[PAGE_URL_ID]` placeholder system for content references.

## PAGE_URL Placeholder System

**Purpose**: Replace `[PAGE_URL_ID]` with page alias in content
**Format**: `[[PAGE_URL_pageID]]` → `/about-us`
**Implementation**: Frontend content rendering pipeline

### How It Works

1. **Content contains**: `[[PAGE_URL_abc123def]]`
2. **Find page by ID**: `PageFindByID("abc123def")`
3. **Get page alias**: `page.Alias()` returns "about-us"
4. **Replace with**: `/about-us`

### Usage

```html
<!-- Content with PAGE_URL placeholder -->
<p>Visit our [[PAGE_URL_abc123def]] page for more information.</p>

<!-- Renders to -->
<p>Visit our /about-us page for more information.</p>
```

### Implementation Details

The system directly uses the page's `Alias()` method and formats it with a leading slash:

```go
// Frontend content rendering (following existing pattern)
pagePath := "/" + strings.TrimPrefix(page.Alias(), "/")
```

#### Frontend Content Processing
```go
// Part of content rendering pipeline
content, err = frontend.contentRenderPageURLs(r.Context(), content)
```

### Features

- **Simple**: Just replaces with page alias
- **Cached**: Page URLs cached for performance
- **Error handling**: Graceful fallback for missing pages
- **Follows BLOCK pattern**: Same approach as `[BLOCK_ID]`

### Admin UI

The page editor displays copy-paste shortcodes:

```
To link to this page permanently, use the following shortcode:

<!-- START: Page: About Us -->
[[PAGE_URL_abc123def]]
<!-- END: Page: About Us -->
```

### Benefits

- **Easy to use**: Simple placeholder syntax
- **Reliable**: Uses existing page lookup infrastructure
- **Performant**: Cached page lookups
- **Consistent**: Follows established CMS patterns

### Example Flow

1. **Page created**: ID="abc123def", Alias="about-us"
2. **Content written**: `[[PAGE_URL_abc123def]]`
3. **Content rendered**: `/about-us`
4. **User clicks**: Goes to `/about-us` (existing routing)

This system provides a simple way to reference pages in content that will work as long as the page exists, regardless of alias changes.
