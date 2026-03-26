# Navbar Block

A dedicated navbar block type for rendering navigation bars with Bootstrap 5 support.

## Features

### 🎯 **Purpose-Built for Navigation**
- Specifically designed for navigation bars (not generic menus)
- Better naming and clearer purpose than the generic "menu" block
- Optimized for Bootstrap 5 navbar components

### 🎨 **Dual Rendering Modes**
- **Bootstrap 5 (Default)**: Full Bootstrap 5 navbar with proper classes and structure
- **Plain**: Simple navbar without framework dependencies

### ⚙️ **Enhanced Configuration**
- **Brand Text & URL**: Add your brand/logo to the navbar
- **Fixed Positioning**: Option to fix navbar to top of page
- **Dark Theme**: Support for both light and dark themes
- **Multiple Styles**: Default, centered, and bottom alignment
- **CSS Customization**: CSS ID and class support

### 📱 **Mobile Responsive**
- Bootstrap 5 includes collapsible mobile menu
- Proper responsive behavior out of the box

## Usage

### Basic Bootstrap 5 Navbar

```go
// Create navbar block with Bootstrap 5 rendering
navbarBlock := navbar.NewNavbarBlockType(store)

// Configure with meta fields
block.SetMeta(cmsstore.BLOCK_META_MENU_ID, "main-menu")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE, "bootstrap5")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT, "My Website")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_BRAND_URL, "/")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_FIXED, "true")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_DARK, "true")
```

### Plain Navbar

```go
// Configure with plain rendering
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE, "plain")
block.SetMeta(cmsstore.BLOCK_META_NAVBAR_STYLE, "centered")
```

## Generated HTML

### Bootstrap 5 Output
```html
<nav class="navbar navbar-expand-lg navbar-dark bg-dark fixed-top">
  <a class="navbar-brand" href="/">My Website</a>
  <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarContent">
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbarContent">
    <ul class="navbar-nav">
      <li class="nav-item">
        <a class="nav-link" href="/home">Home</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="/about">About</a>
      </li>
    </ul>
  </div>
</nav>
```

### Plain Output
```html
<nav class="navbar navbar-style-default">
  <a class="navbar-brand" href="/">My Website</a>
  <ul class="navbar-menu">
    <li class="navbar-item">
      <a class="navbar-link" href="/home">Home</a>
    </li>
    <li class="navbar-item">
      <a class="navbar-link" href="/about">About</a>
    </li>
  </ul>
</nav>
```

## Configuration Options

| Meta Key | Description | Default |
|----------|-------------|---------|
| `menu_id` | Menu ID to display | Required |
| `navbar_style` | Layout style (default, centered, bottom) | `default` |
| `navbar_rendering_mode` | Rendering mode (bootstrap5, plain) | `bootstrap5` |
| `navbar_brand_text` | Brand text to display | Empty |
| `navbar_brand_url` | Brand link URL | `/` |
| `navbar_fixed` | Fixed positioning (true/false) | `false` |
| `navbar_dark` | Dark theme (true/false) | `false` |
| `navbar_css_class` | Additional CSS classes | Empty |
| `navbar_css_id` | CSS ID for the navbar | Empty |

## Comparison with Menu Block

| Feature | Menu Block | Navbar Block |
|---------|------------|--------------|
| **Purpose** | Generic menu lists | Navigation bars |
| **Structure** | Simple nav/ul lists | Bootstrap navbar components |
| **Brand Support** | No | Yes (text + URL) |
| **Positioning** | No | Fixed positioning option |
| **Mobile Support** | No | Collapsible mobile menu |
| **Theme Support** | No | Dark/light themes |
| **Default Rendering** | Plain | Bootstrap 5 |

## Files

- `navbar_block_type.go` - Main block implementation
- `renderer.go` - Rendering logic for both modes
- `README.md` - This documentation

## Dependencies

- `github.com/dracory/hb` - HTML builder library
- `github.com/dracory/cmsstore` - Core CMS interfaces
