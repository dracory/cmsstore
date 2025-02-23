# Cms Store <a href="https://gitpod.io/#https://github.com/gouniverse/cmsstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

<img src="https://opengraph.githubassets.com/5b92c81c05d64a82c3fb4ba95739403a2d38cbad61f260a0701b3366b3d10327/gouniverse/cmsstore" />

[![Tests Status](https://github.com/gouniverse/cmsstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/cmsstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/cmsstore)](https://goreportcard.com/report/github.com/gouniverse/cmsstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/cmsstore)](https://pkg.go.dev/github.com/gouniverse/cmsstore)

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
go get -u github.com/gouniverse/cmsstore
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
