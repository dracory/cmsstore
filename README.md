# Cms Store <a href="https://gitpod.io/#https://github.com/dracory/cmsstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/dracory/cmsstore" />

[![Tests Status](https://github.com/dracory/cmsstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/cmsstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/cmsstore)](https://goreportcard.com/report/github.com/dracory/cmsstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/cmsstore)](https://pkg.go.dev/github.com/dracory/cmsstore)

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Introduction

All of the existing GoLang CMSs require a full installations from scratch. 
Its impossible to just add them to an exiting Go application, and even when added
feel like you don't get what you hoped for.

This package allows to add a content management system as a module dependency,
which can be easily updated or removed as required to ANY Go app.
It is fully self contained, and does not require any additional packages
or dependencies. Removal is also a breeze just remove the module.

## Installation

```
go get -u github.com/dracory/cmsstore
```

## Features

- Multi-site
- Templates
- Pages
- Blocks
- Menus
- Translations
- Custom Entity Types
- Supports middleware
- Supports shortcodes

## Simplest Initialization

The simplest initialization involves providing a database instance. Note that this minimal setup has limited capabilities; features like database migrations are not automatically run.

```go
// Establish database connection. Replace with your actual database connection details.
db, err := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))
if err != nil {
	log.Panic("Database is NIL: " + err.Error())
	return
}
if db == nil {
	log.Panic("Database is NIL")
	return
}

// Configure store options. Customize table names as needed. All table names are prefixed with `cms_`.
opts := cmsstore.NewStoreOptions{
    DB: db,
    BlockTableName: "cms_blocks", // Name of the database table for blocks
    PageTableName: "cms_pages",   // Name of the database table for pages
    SiteTableName: "cms_sites",   // Name of the database table for sites
    TemplateTableName: "cms_templates", // Name of the database table for templates
    MenusEnabled: true,       // Enable menu functionality
    MenuTableName: "cms_menus",   // Name of the database table for menus
    MenuItemTableName: "cms_menu_items", // Name of the database table for menu items
    TranslationsEnabled: true, // Enable translation functionality
    TranslationTableName: "cms_translations", // Name of the database table for translations
    TranslationLanguageDefault: "en", // Default language code
    TranslationLanguages: map[string]string{"en": "English"}, // Supported languages
    VersioningEnabled: true,   // Enable versioning functionality
    VersioningTableName: "cms_versions", // Name of the database table for versions
}

// Create a new store instance with the specified options.
store, err := cmsstore.NewStore(opts)
if err != nil {
    log.Panic(err)
}

// Create a new CMS instance using the created store.
myCms, errCms := cms.NewCms(cms.Config{
	Store: store,
})
```

## Shortcodes

Shortcodes provide a powerful way to inject custom complex rendering logic
into your content. i.e database queries, showing lists of shop products,
blog posts, etc.

This allows for highly flexible and customizable content generation, without
the need to write complex templates.

Shortcodes are defined externally as part of your Go project and registered
with the store during initialization.

Another advantage of shortcodes is their versioning. Because shortcodes are
part of your project's standard code, they are easily version-controlled.

**Defining a Shortcode:**

Create a function that implements the `cmsstore.ShortcodeInterface`. This interface requires three methods:

- `Alias() string`: Returns a unique identifier for your shortcode.
- `Description() string`: Provides a brief description of the shortcode's purpose.
- `Render(r *http.Request, s string, m map[string]string) string`: This method performs the actual rendering. It takes the HTTP request, the content string, and a map of parameters as input and returns the rendered string.

**Registering a Shortcode:**

When creating a new store using `cmsstore.NewStore`, pass your custom shortcode functions in the `Shortcodes` field of the `cmsstore.NewStoreOptions` struct.

**Using a Shortcode:**

Place your shortcode within your content using the specified delimiters (configurable via `shortcode.NewShortcode`). The shortcode will be rendered dynamically during content processing.

**Example:**

Let's say you have a shortcode named "my-shortcode" that takes a "name" parameter. You would use it in your content like this: `<my-shortcode name="John Doe">`. The `Render` method of your shortcode implementation would then process this and generate the appropriate output.

## Multisite Functionality

This Content Management System (CMS) offers a robust multisite functionality,
enabling you to manage multiple websites from a single installation.

Each site has its own:
- Templates: layouts for your pages
- Pages: individual content pages for each site
- Blocks: reusable content blocks that can be used across multiple pages and sites, i.e. header, footer, sidebar, galleries, sliders, carousels, testimonials, FAQs, etc.
- Menus: customizable navigation menus for each site
- Translations: support for multiple languages

It alows you to:
- Manage multiple brands, regions, or languages from a single installation
- Create separate websites for specific campaigns or events
- Easily setup different themes for each site
- Rapidly add landing pages, tracking pages, sales funnel pages, etc.

## Pages

Pages are the core content units of the CMS.  They are structured using a flexible
data model that allows for rich content and metadata.  

Pages can be created, updated, deleted, and versioned through
a dedicated admin interface.

They support various features such as custom templates, SEO metadata 
(keywords, descriptions, robots directives), and middleware
for custom processing.

### Benefits

- **Organized Content:** Pages provide a structured way to organize and manage website content.
- **SEO Optimization:** Built-in support for SEO metadata helps improve search engine rankings.
- **Flexible Templating:**  The use of templates allows for consistent design and branding across the site.
- **Version Control:** Versioning capabilities allow for easy rollback to previous versions.
- **Customizable Workflows:** Middleware support enables custom processing and workflows.

### SEO Capabilities

Pages support various SEO features, including:

- **Meta Keywords:** Specify relevant keywords for search engines.
- **Meta Descriptions:** Provide concise descriptions for search results.
- **Meta Robots:** Control how search engines crawl and index the page.
- **Canonical URLs:** Specify the preferred URL for the page, preventing duplicate content issues.

### Editors

Pages can be edited through a user-friendly admin interface.
The admin interface provides tools for creating, updating, and deleting pages,
as well as managing page versions.
Page content can be easily managed using a variety of editors:
- **CodeMirror Editor:** Raw HTML editing for advanced customization.
- **WYSIWYG Editor:** Rich text editing for advanced formatting.
- **Markdown Editor:** Markdown syntax for simple text formatting.
TODO: Access control can be implemented to restrict editing permissions.

## Templates

Templates define the layout of pages.
They can be customized to create different page designs.
Templates can be associated with pages, allowing for flexible and reusable designs.
The system supports template versioning and management through a dedicated admin interface.

### Benefits

- **Consistent Design:** Templates ensure a consistent look and feel across the website.
- **Reusable Layouts:**  Templates allow for the creation of reusable layouts, reducing development time and effort.
- **Flexible Layouts:**  Templates can be customized to create various page layouts.
- **Version Control:** Versioning capabilities allow for easy rollback to previous versions.

## Translations

This CMS supports multilingual content through a robust translation system.  Translations are managed as individual entities, allowing for efficient management and updates.  Each translation is associated with a specific content item and language code.  The system supports multiple languages and allows for easy switching between languages.

### Benefits of Using Translations

- **Multilingual Support:** Easily create and manage content in multiple languages.
- **Global Reach:** Expand your reach to a wider audience by providing content in their native language.
- **Improved User Experience:** Provide a more user-friendly experience by offering content in the user's preferred language.
- **Efficient Management:** Manage translations efficiently through a dedicated admin interface.

### Using Translations

Translations are inserted into content using placeholders of the form `[[TRANSLATION_translationID]]`. The frontend code then replaces these placeholders with the appropriate translated text based on the user's selected language.  Administrators can manage translations through a dedicated admin interface, providing tools for creating, updating, and deleting translations.  The system automatically selects the appropriate translation based on the user's language preference.

## Blocks

Blocks are reusable content units that can be assembled to create pages.
They provide a modular approach to content creation, allowing for flexible
and dynamic page layouts.
Blocks can be managed through a dedicated admin interface, allowing for easy
creation, updating, and deletion.
The use of blocks promotes reusability, consistency, and maintainability
of website content.

### Benefits of Using Blocks

- **Reusability:** Create once, use many times across different pages and sites.
- **Consistency:** Maintain a consistent look and feel across your website.
- **Maintainability:** Easily update content in one place and have it reflected everywhere.
- **Flexibility:** Create dynamic and flexible page layouts by combining different blocks.
- **Efficiency:** Reduce development time and effort by reusing existing blocks.

### Using Blocks

Blocks are inserted into pages using placeholders of the form `[[BLOCK_blockID]]`.
The frontend code then replaces these placeholders with the rendered content
of the corresponding block.
This allows for dynamic content generation and flexible page layouts.
Administrators can manage blocks through a dedicated admin interface,
providing tools for creating, updating, and deleting blocks.

## Menus

This CMS provides a robust menu management system, allowing you to create and manage hierarchical menus for your website.  Menus are structured as trees, enabling you to organize your navigation in a clear and intuitive way.  The system supports creating, updating, deleting, and filtering menus through a dedicated admin interface.

### Benefits of Using Menus

- **Organized Navigation:** Create clear and intuitive navigation for your website visitors.
- **Improved User Experience:** Guide users easily through your website content.
- **Flexible Structure:** Create hierarchical menus to reflect the structure of your website.
- **Efficient Management:** Manage menus efficiently through a dedicated admin interface.

### Using Menus

Menus are created and managed through the admin interface, which provides tools for creating, updating, and deleting menu items.  The hierarchical structure of menus allows you to organize your navigation in a clear and intuitive way.  The system supports various menu types and allows for customization of menu items.

## CMS URL Patterns

The following URL patterns are supported:

- :any - ([^/]+)
- :num - ([0-9]+)
- :all - (.*)
- :string - ([a-zA-Z]+)
- :number - ([0-9]+)
- :numeric - ([0-9-.]+)
- :alpha - ([a-zA-Z0-9-_]+)

Example:
```
/blog/:num/:any
/shop/product/:num/:any
```

# Documentation

For more information, please refer to the [Documentation](./docs/README.md).
