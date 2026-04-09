package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dracory/cmsstore"
)

// BlogBlockType demonstrates how to use custom variables in blocks
type BlogBlockType struct{}

func (b *BlogBlockType) Type() string {
	return "blog_post"
}

func (b *BlogBlockType) Label() string {
	return "Blog Post"
}

func (b *BlogBlockType) Description() string {
	return "A blog post block that exposes variables for use in page/template content"
}

func (b *BlogBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{
		{Name: "blog_title", Description: "The blog post title"},
		{Name: "blog_author", Description: "The post author name"},
		{Name: "blog_date", Description: "Publication date in YYYY-MM-DD format"},
		{Name: "blog_category", Description: "The post category"},
		{Name: "blog_reading_time", Description: "Estimated reading time"},
		{Name: "blog_excerpt", Description: "Short summary of the post"},
	}
}

func (b *BlogBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
	// Simulate fetching blog post data
	post := struct {
		Title       string
		Author      string
		Date        time.Time
		Category    string
		ReadingTime string
		Excerpt     string
	}{
		Title:       "Understanding Custom Variables in CMS",
		Author:      "Jane Developer",
		Date:        time.Now(),
		Category:    "Development",
		ReadingTime: "5 min read",
		Excerpt:     "Learn how to use custom variables to expose block data to your pages and templates.",
	}

	// Set custom variables that can be used anywhere in the page/template
	if vars := cmsstore.VarsFromContext(ctx); vars != nil {
		// Using snake_case naming
		vars.Set("blog_title", post.Title)
		vars.Set("blog_author", post.Author)
		vars.Set("blog_date", post.Date.Format("2006-01-02"))
		vars.Set("blog_category", post.Category)
		vars.Set("blog_reading_time", post.ReadingTime)
		vars.Set("blog_excerpt", post.Excerpt)

		// You can also use other naming conventions
		vars.Set("BlogTitle", post.Title)  // PascalCase
		vars.Set("blog:title", post.Title) // Namespaced
		vars.Set("$blogTitle", post.Title) // Prefixed
		vars.Set("post.title", post.Title) // Dotted
		vars.Set("POST_TITLE", post.Title) // Upper case
	}

	// Return the HTML for the block itself
	html := fmt.Sprintf(`
		<article class="blog-post">
			<h2>%s</h2>
			<div class="meta">
				<span class="author">By %s</span>
				<span class="date">%s</span>
				<span class="category">%s</span>
				<span class="reading-time">%s</span>
			</div>
			<div class="excerpt">
				<p>%s</p>
			</div>
		</article>
	`, post.Title, post.Author, post.Date.Format("January 2, 2006"), post.Category, post.ReadingTime, post.Excerpt)

	return html, nil
}

// ProductBlockType demonstrates using custom variables for e-commerce
type ProductBlockType struct{}

func (p *ProductBlockType) Type() string {
	return "product"
}

func (p *ProductBlockType) Label() string {
	return "Product"
}

func (p *ProductBlockType) Description() string {
	return "A product block that exposes product data as variables"
}

func (p *ProductBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{
		{Name: "product_name", Description: "The product name"},
		{Name: "product_price", Description: "The product price with currency symbol"},
		{Name: "product_sku", Description: "The product SKU/model number"},
		{Name: "product_availability", Description: "Stock availability status"},
		{Name: "product_brand", Description: "The product brand name"},
	}
}

func (p *ProductBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
	// Simulate product data
	product := struct {
		Name         string
		Price        string
		SKU          string
		Availability string
		Brand        string
	}{
		Name:         "Wireless Headphones",
		Price:        "$149.99",
		SKU:          "WH-1000XM4",
		Availability: "In Stock",
		Brand:        "TechBrand",
	}

	// Set product variables
	if vars := cmsstore.VarsFromContext(ctx); vars != nil {
		vars.Set("product_name", product.Name)
		vars.Set("product_price", product.Price)
		vars.Set("product_sku", product.SKU)
		vars.Set("product_availability", product.Availability)
		vars.Set("product_brand", product.Brand)
	}

	html := fmt.Sprintf(`
		<div class="product">
			<h3>%s</h3>
			<p class="price">%s</p>
			<p class="sku">SKU: %s</p>
			<p class="availability">%s</p>
			<p class="brand">Brand: %s</p>
		</div>
	`, product.Name, product.Price, product.SKU, product.Availability, product.Brand)

	return html, nil
}

func main() {
	fmt.Println("Custom Variables Example")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("This example demonstrates how to use custom variables in blocks.")
	fmt.Println()
	fmt.Println("1. Register your block type:")
	fmt.Println("   cmsstore.RegisterBlockType(&BlogBlockType{})")
	fmt.Println()
	fmt.Println("2. In your block's Render method, set variables:")
	fmt.Println("   if vars := cmsstore.VarsFromContext(ctx); vars != nil {")
	fmt.Println("       vars.Set(\"blog_title\", \"My Post\")")
	fmt.Println("       vars.Set(\"author_name\", \"John Doe\")")
	fmt.Println("   }")
	fmt.Println()
	fmt.Println("3. Use variables in your page/template content:")
	fmt.Println("   <h1>[[blog_title]]</h1>")
	fmt.Println("   <p>By [[author_name]]</p>")
	fmt.Println()
	fmt.Println("Example page content:")
	fmt.Println("---------------------")
	fmt.Println(`<!DOCTYPE html>
<html>
<head>
    <title>[[PageTitle]] - [[blog_title]]</title>
    <meta name="description" content="[[blog_excerpt]]">
    <meta property="og:title" content="[[blog_title]]">
    <meta property="article:author" content="[[blog_author]]">
    <meta property="article:published_time" content="[[blog_date]]">
</head>
<body>
    <header>
        <h1>[[blog_title]]</h1>
        <p>By [[blog_author]] on [[blog_date]] in [[blog_category]]</p>
        <p>[[blog_reading_time]]</p>
    </header>
    
    <main>
        [[BLOCK_blog_content]]
    </main>
    
    <footer>
        <p>Author: [[blog_author]]</p>
    </footer>
</body>
</html>`)
	fmt.Println()
	fmt.Println("Variable Naming Conventions:")
	fmt.Println("---------------------------")
	fmt.Println("You can use any naming convention:")
	fmt.Println("  - snake_case: blog_title, user_name")
	fmt.Println("  - camelCase: blogTitle, userName")
	fmt.Println("  - PascalCase: BlogTitle, UserName")
	fmt.Println("  - namespaced: blog:title, user:name")
	fmt.Println("  - prefixed: $blogTitle, @userName")
	fmt.Println("  - dotted: blog.title, user.name")
}
