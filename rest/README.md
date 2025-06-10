# CMS Store REST API

This package provides a REST API for the CMS Store, allowing web applications to interact with the CMS functionality through standard HTTP requests. It's designed to be easily integrated into any existing Go application.

## Features

- Page management (create, read, update, delete)
- Menu management (create, read, list)
- Simple integration with any existing Go HTTP server
- JSON responses for all endpoints

## Getting Started

### Prerequisites

- Go 1.24 or later
- A running instance of the CMS Store

### Installation

```bash
go get github.com/gouniverse/cmsstore/rest
```

### Basic Usage

The REST API is designed to be easily integrated into any existing Go HTTP server:

```go
package main

import (
	"log"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/rest"
)

func main() {
	// Initialize your CMS store
	store := cmsstore.NewStore(db) // Your database connection

	// Create the REST API
	api := rest.NewRestAPI(store)
	
	// Register the API handler with your existing router/mux
	http.HandleFunc("/api/", api.Handler())
	
	// Start your server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

## API Reference

### Page Endpoints

#### Create a Page

**Request:**
```
POST /api/pages
Content-Type: application/json

{
  "title": "My New Page",
  "content": "<p>Page content goes here</p>",
  "status": "published",
  "site_id": "site-123"
}
```

**Response:**
```json
{
  "success": true,
  "id": "page_123",
  "title": "My New Page",
  "content": "<p>Page content goes here</p>",
  "status": "published"
}
```

#### Get a Page

**Request:**
```
GET /api/pages/{page_id}
```

**Response:**
```json
{
  "success": true,
  "id": "page_123",
  "title": "My New Page",
  "content": "<p>Page content goes here</p>",
  "status": "published"
}
```

#### List All Pages

**Request:**
```
GET /api/pages
```

**Response:**
```json
{
  "success": true,
  "pages": [
    {
      "id": "page_123",
      "title": "My New Page",
      "content": "<p>Page content goes here</p>",
      "status": "published"
    },
    {
      "id": "page_124",
      "title": "Another Page",
      "content": "<p>More content</p>",
      "status": "draft"
    }
  ]
}
```

#### Update a Page

**Request:**
```
PUT /api/pages/{page_id}
Content-Type: application/json

{
  "title": "Updated Page Title",
  "content": "<p>Updated content</p>"
}
```

**Response:**
```json
{
  "success": true,
  "id": "page_123",
  "title": "Updated Page Title",
  "content": "<p>Updated content</p>",
  "status": "published"
}
```

#### Delete a Page

**Request:**
```
DELETE /api/pages/{page_id}
```

**Response:**
```json
{
  "success": true,
  "message": "Page deleted successfully"
}
```

### Menu Endpoints

#### Create a Menu

**Request:**
```
POST /api/menus
Content-Type: application/json

{
  "title": "Main Menu",
  "site_id": "site-123"
}
```

**Response:**
```json
{
  "success": true,
  "id": "menu_123",
  "title": "Main Menu"
}
```

#### Get a Menu

**Request:**
```
GET /api/menus/{menu_id}
```

**Response:**
```json
{
  "success": true,
  "id": "menu_123",
  "title": "Main Menu"
}
```

#### List All Menus

**Request:**
```
GET /api/menus
```

**Response:**
```json
{
  "success": true,
  "menus": [
    {
      "id": "menu_123",
      "title": "Main Menu"
    },
    {
      "id": "menu_124",
      "title": "Footer Menu"
    }
  ]
}
```

## Error Handling

Errors are returned with appropriate HTTP status codes and JSON bodies:

```json
{
  "success": false,
  "error": "Detailed error message here"
}
```

## Extending the API

To add new endpoints or functionality to the REST API, you can extend the `RestAPI` struct in `rest.go`. Follow the pattern of the existing handlers to maintain consistency.

## License

This package is part of the CMS Store and is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).
