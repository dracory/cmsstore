# Breadcrumbs Block

A breadcrumbs block type for displaying navigation trails that show users their current location within a website's hierarchy.

## Features

### 🎯 **Purpose-Built for Navigation**
- Specifically designed for breadcrumb navigation trails
- Shows users their current location in the site hierarchy
- Provides clear navigation context and easy back-navigation

### 🎨 **Dual Rendering Modes**
- **Bootstrap 5 (Default)**: Full Bootstrap 5 breadcrumb components with proper accessibility
- **Plain**: Simple breadcrumb structure without framework dependencies

### ⚙️ **Enhanced Configuration**
- **Custom Separator**: Choose between `/`, `>`, `→`, or any custom separator
- **Home Link**: Configurable home text and URL
- **Multiple Styles**: Default, centered, and right-aligned breadcrumbs
- **CSS Customization**: CSS ID and class support for styling

### 📱 **Accessibility Ready**
- Bootstrap 5 includes proper ARIA attributes
- Semantic HTML structure for screen readers
- Clear visual hierarchy for all users

## Usage

### Basic Bootstrap 5 Breadcrumbs

```go
// Create breadcrumbs block with Bootstrap 5 rendering
breadcrumbsBlock := breadcrumbs.NewBreadcrumbsBlockType(store)

// Configure with meta fields
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE, "default")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE, "bootstrap5")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, "/")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT, "Home")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL, "/")
```

### Plain Breadcrumbs

```go
// Configure with plain rendering
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE, "plain")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE, "centered")
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, "→")
```

## Generated HTML

### Bootstrap 5 Output
```html
<nav class="breadcrumb" aria-label="breadcrumb">
  <li class="breadcrumb-item">
    <a href="/">Home</a>
  </li>
  <li class="breadcrumb-item">
    <a href="/products">Products</a>
  </li>
  <li class="breadcrumb-item active" aria-current="page">
    Electronics
  </li>
</nav>
```

### Plain Output
```html
<nav class="breadcrumbs breadcrumbs-style-default">
  <a class="breadcrumb-link" href="/">Home</a>
  <span class="breadcrumb-separator">/</span>
  <a class="breadcrumb-link" href="/products">Products</a>
  <span class="breadcrumb-separator">/</span>
  <span class="breadcrumb-item active">Electronics</span>
</nav>
```

## Configuration Options

| Meta Key | Description | Default |
|----------|-------------|---------|
| `breadcrumbs_style` | Layout style (default, centered, right) | `default` |
| `breadcrumbs_rendering_mode` | Rendering mode (bootstrap5, plain) | `bootstrap5` |
| `breadcrumbs_separator` | Separator between items | `/` |
| `breadcrumbs_home_text` | Home link text | `Home` |
| `breadcrumbs_home_url` | Home link URL | `/` |
| `breadcrumbs_css_class` | Additional CSS classes | Empty |
| `breadcrumbs_css_id` | CSS ID for the breadcrumbs | Empty |

## Breadcrumb Generation

The breadcrumbs block automatically generates breadcrumb items based on the current page context:

1. **Home**: Always included as the first breadcrumb (configurable)
2. **Current Page**: Added as the active breadcrumb (no link)
3. **Dynamic Path**: Future versions can support full page hierarchy

## Styling Examples

### Bootstrap 5 Variants
```html
<!-- Default left-aligned -->
<nav class="breadcrumb" aria-label="breadcrumb">...</nav>

<!-- Centered -->
<nav class="breadcrumb justify-content-center" aria-label="breadcrumb">...</nav>

<!-- Right-aligned -->
<nav class="breadcrumb justify-content-end" aria-label="breadcrumb">...</nav>
```

### Plain Variants
```html
<!-- Default -->
<nav class="breadcrumbs breadcrumbs-style-default">...</nav>

<!-- Centered -->
<nav class="breadcrumbs breadcrumbs-style-centered">...</nav>

<!-- Right-aligned -->
<nav class="breadcrumbs breadcrumbs-style-right">...</nav>
```

## Custom Separators

```go
// Arrow separator
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, "→")

// Greater than separator
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, ">")

// Custom text separator
block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, "•")
```

## Integration with Other Blocks

The breadcrumbs block works well with:
- **Navbar Block**: For main navigation
- **Menu Block**: For footer or sidebar navigation
- **Page Content**: Automatically detects current page context

## Files

- `breadcrumbs_block_type.go` - Main block implementation
- `renderer.go` - Rendering logic for both modes
- `README.md` - This documentation

## Dependencies

- `github.com/dracory/hb` - HTML builder library
- `github.com/dracory/cmsstore` - Core CMS interfaces

## Best Practices

1. **Keep breadcrumbs concise**: Limit to 3-5 levels for optimal UX
2. **Use meaningful labels**: Make breadcrumb text descriptive and concise
3. **Consistent placement**: Place breadcrumbs near the top of the page
4. **Proper hierarchy**: Follow logical site structure in breadcrumb order
