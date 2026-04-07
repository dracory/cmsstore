# Custom Variables Example

This example demonstrates how to use custom variables in blocks to expose dynamic data that can be referenced anywhere in page or template content.

## Overview

Custom variables allow blocks to set arbitrary key-value pairs that are automatically replaced in the final rendered content. This is useful for:

- Blog posts exposing title, author, date, etc.
- Products exposing price, SKU, availability
- Events exposing date, location, organizer
- User profiles exposing user data
- Any dynamic content that needs to be referenced elsewhere

## How It Works

### 1. In Your Block's Render Method

```go
func (b *BlogBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
    // Set custom variables
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog_title", "My Blog Post")
        vars.Set("author_name", "John Doe")
        vars.Set("publish_date", "2026-04-07")
    }
    
    return html, nil
}
```

### 2. In Your Page/Template Content

```html
<!DOCTYPE html>
<html>
<head>
    <title>[[PageTitle]] - [[blog_title]]</title>
    <meta property="og:title" content="[[blog_title]]">
    <meta property="article:author" content="[[author_name]]">
</head>
<body>
    <h1>[[blog_title]]</h1>
    <p>By [[author_name]] on [[publish_date]]</p>
    
    [[BLOCK_blog_content]]
</body>
</html>
```

### 3. Rendered Output

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Site - My Blog Post</title>
    <meta property="og:title" content="My Blog Post">
    <meta property="article:author" content="John Doe">
</head>
<body>
    <h1>My Blog Post</h1>
    <p>By John Doe on 2026-04-07</p>
    
    <article>...</article>
</body>
</html>
```

## Variable Naming Conventions

You have complete freedom in naming variables. Choose what works for your project:

- **snake_case**: `blog_title`, `user_name`, `product_price`
- **camelCase**: `blogTitle`, `userName`, `productPrice`
- **PascalCase**: `BlogTitle`, `UserName`, `ProductPrice`
- **namespaced**: `blog:title`, `product:price`, `event:date`
- **prefixed**: `$blogTitle`, `@userName`, `#productPrice`
- **dotted**: `blog.title`, `user.name`, `product.price`

## Multiple Blocks

Multiple blocks can set different variables without conflicts:

```go
// Blog block
vars.Set("blog_title", "My Post")
vars.Set("blog_author", "John")

// Product block
vars.Set("product_name", "Widget")
vars.Set("product_price", "$99")

// User block
vars.Set("user_name", "Jane")
vars.Set("user_role", "Admin")
```

All variables are available for substitution in the final content.

## Variable Collisions

If multiple blocks set the same variable name, the last block (in content order) wins. To avoid collisions:

1. Use unique variable names
2. Use namespacing (e.g., `blog:title`, `product:title`)
3. Use prefixes (e.g., `blog_title`, `product_title`)

## Running the Example

```bash
cd examples/custom-variables
go run main.go
```

This will display information about how to use custom variables in your blocks.

## Use Cases

### Blog Post SEO

```go
vars.Set("blog_title", post.Title)
vars.Set("blog_excerpt", post.Excerpt)
vars.Set("blog_author", post.Author)
vars.Set("blog_date", post.Date.Format("2006-01-02"))
vars.Set("blog_image", post.FeaturedImage)
```

Then in your template:

```html
<head>
    <title>[[blog_title]] | My Blog</title>
    <meta name="description" content="[[blog_excerpt]]">
    <meta property="og:title" content="[[blog_title]]">
    <meta property="og:description" content="[[blog_excerpt]]">
    <meta property="og:image" content="[[blog_image]]">
    <meta property="article:author" content="[[blog_author]]">
    <meta property="article:published_time" content="[[blog_date]]">
</head>
```

### E-commerce Product

```go
vars.Set("product_name", product.Name)
vars.Set("product_price", product.Price)
vars.Set("product_sku", product.SKU)
vars.Set("product_brand", product.Brand)
vars.Set("product_availability", product.Availability)
```

Then in your template:

```html
<head>
    <title>[[product_name]] - [[product_brand]]</title>
    <script type="application/ld+json">
    {
        "@context": "https://schema.org/",
        "@type": "Product",
        "name": "[[product_name]]",
        "brand": "[[product_brand]]",
        "sku": "[[product_sku]]",
        "offers": {
            "@type": "Offer",
            "price": "[[product_price]]",
            "availability": "[[product_availability]]"
        }
    }
    </script>
</head>
```

### Event Information

```go
vars.Set("event_name", event.Name)
vars.Set("event_date", event.Date.Format("2006-01-02"))
vars.Set("event_location", event.Location)
vars.Set("event_organizer", event.Organizer)
```

Then in your template:

```html
<head>
    <title>[[event_name]] - [[event_date]]</title>
    <meta property="og:title" content="[[event_name]]">
    <meta property="event:start_time" content="[[event_date]]">
    <meta property="event:location" content="[[event_location]]">
</head>
<body>
    <h1>[[event_name]]</h1>
    <p>Date: [[event_date]]</p>
    <p>Location: [[event_location]]</p>
    <p>Organized by: [[event_organizer]]</p>
</body>
```

## Best Practices

1. **Check for nil**: Always check if VarsContext is available before using it
   ```go
   if vars := cmsstore.VarsFromContext(ctx); vars != nil {
       vars.Set("key", "value")
   }
   ```

2. **Use descriptive names**: Choose variable names that clearly indicate their purpose
   ```go
   vars.Set("blog_post_title", title)  // Good
   vars.Set("t", title)                // Bad
   ```

3. **Avoid collisions**: Use namespacing or prefixes to prevent conflicts
   ```go
   vars.Set("blog:title", title)       // Namespaced
   vars.Set("blog_title", title)       // Prefixed
   ```

4. **Document your variables**: Document which variables your block sets
   ```go
   // BlogBlockType sets the following variables:
   // - blog_title: The blog post title
   // - blog_author: The author name
   // - blog_date: Publication date in YYYY-MM-DD format
   ```

5. **Handle missing data**: Provide defaults or empty strings for missing data
   ```go
   author := post.Author
   if author == "" {
       author = "Anonymous"
   }
   vars.Set("blog_author", author)
   ```
