# CMS Store - Project Overview

## Project Purpose
A modular, embeddable CMS package for Go applications that provides content management capabilities without requiring a full-stack framework.

## Technical Stack
- **Language**: Go (Golang)
- **Database**: SQL (database-agnostic, uses database/sql)
- **ORM**: neat (github.com/dracory/neat)
- **Dependency Management**: Go modules
- **Testing**: Native Go testing package

## Architecture
- **Modular Design**: Separate components for blocks, menus, pages, templates, sites, and translations
- **Store Pattern**: Central store interface for database operations
- **ORM Pattern**: Type-safe database operations using neat
- **Middleware Support**: Extensible through middleware interface
- **Versioning**: Built-in support for content versioning

## Key Components
1. **Store**: Core database operations and state management
2. **Query Builders**: Type-safe query construction for each entity
3. **Models**: Page, Block, Menu, Site, Template, and Translation
4. **Interfaces**: Well-defined contracts for all major components

## Project Structure
- `/admin`: Admin interface components
- `/frontend`: Frontend assets and templates
- `store.go`: Core store implementations with schema migrations
- `store_*.go`: Entity-specific store implementations
- `*_query.go`: Query builder implementations
- `*_test.go`: Test files

## Dependencies
- `github.com/dracory/neat`: ORM for database operations
- `github.com/dromara/carbon/v2`: Date/time handling
- `github.com/dracory/database`: Database utilities
- `github.com/samber/lo`: Lo-Dash like Go utilities

## Development Status
- **Stable**: Core functionality implemented and tested
- **Documentation**: In progress, with focus on query interfaces
- **Testing**: Comprehensive test coverage for store operations

## License
GNU Affero General Public License v3.0 (AGPL-3.0)
