# Architecture Overview

## System Components

The CMS Store is built with a modular architecture, consisting of the following main components:

### Core Components

1. **Blocks**
   - Manages content blocks
   - Handles block rendering and storage
   - Supports dynamic content insertion

2. **Menus**
   - Menu management system
   - Hierarchical menu structure
   - Dynamic menu item handling

3. **Pages**
   - Page content management
   - URL routing and handling
   - Template integration

4. **Templates**
   - Template management
   - Theme support
   - Layout handling

5. **Translations**
   - Multi-language support
   - Content localization
   - Language switching

## Entity Relationships

### Core Entity Structure
```mermaid
erDiagram
    %% Site relationships - Site is the central entity
    Site ||--o{ Page : "can have many"
    Site ||--o{ Template : "can have many"
    Site ||--o{ Block : "can have many"
    Site ||--o{ Menu : "can have many"
    Site ||--o{ Translation : "can have many"
    
    %% Page relationships
    Page o|--o{ Block : "can own many"
    Page o|--o| Template : "can use one"
    Page }|--|| Site : "must belong to one"
    
    %% Menu relationships
    Menu ||--o{ MenuItem : "can have many"
    MenuItem o|--o| Page : "can link to one"
    MenuItem o|--o| MenuItem : "can have parent"
    Menu }|--|| Site : "must belong to one"
    
    %% Template relationships
    Template }|--|| Site : "must belong to one"
    
    %% Block relationships
    Block }|--|| Site : "must belong to one"
    
    %% Translation relationships
    Translation }|--|| Site : "must belong to one"
```

### Entity Relationships Explained

1. **Site-Centric Architecture**
   - Sites are the top-level entities
   - Each site can have zero or many:
     * Pages
     * Templates
     * Blocks
     * Menus
     * Translations
   - All entities must belong to exactly one site

2. **Page Relationships**
   - Each page must belong to exactly one site
   - Each page can optionally use one template
   - Each page can own zero or many blocks
   - Pages can be optionally linked from zero or many menu items

3. **Menu Structure**
   - Each menu must belong to exactly one site
   - Each menu can have zero or many menu items
   - Each menu item:
     * Can optionally link to one page
     * Can optionally have one parent menu item
     * Can have zero or many child menu items

4. **Content Management**
   - Each block:
     * Must belong to exactly one site
     * Can either belong to one specific page or be global (available to all pages)
   - Each template must belong to exactly one site
   - Each translation must belong to exactly one site

## Database Schema

The system uses a relational database with a modular schema. While some tables are required for core CMS functionality, others are optional and only used when specific features are enabled.

### Core Entities (Required)
- `site` - The root entity for multi-site management; stores site configuration and domain mapping.
- `page` - Stores page content, URL aliases, SEO metadata, and template associations.
- `block` - Stores reusable content blocks (can be global or page-specific).
- `template` - Defines layouts and design structures for pages.

### Optional Entities (Feature-Dependent)
- `menu` - Stores navigation menu structures (Enabled via `MenusEnabled`).
- `menu_item` - Stores individual links and hierarchy within menus (Enabled via `MenusEnabled`).
- `translation` - Stores localized content for multi-language support (Enabled via `TranslationsEnabled`).
- `versioning` - Stores historical snapshots of all the above entities for audit and rollback (Enabled via `VersioningEnabled`).

Each table includes standard fields for:
- Unique identifiers (Short IDs)
- Timestamps (`created_at`, `updated_at`)
- Soft delete support (`soft_deleted_at`)
- Status tracking (`active`, `inactive`, `draft`)
- Extensible metadata storage (JSON)

### Transactional Versioning
The system provides strong consistency between entity data and its version history. All write operations on versioned entities are performed within a single database transaction, ensuring that a version snapshot is never missed if an update succeeds, and vice versa.

## Middleware System

The CMS implements a middleware system for:
- Authentication
- Authorization
- Caching
- Request/Response handling

## Future Improvements

- Automated schema management
- Enhanced middleware management
- Code generation for query structs
- Foreign key constraints implementation
- Domain name normalization 