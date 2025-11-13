# Installation Guide

## Prerequisites

- Go 1.23.3 or higher
- One of the following SQL databases:
  - SQLite (using modernc.org/sqlite)
  - MySQL 5.7+ (using go-sql-driver/mysql)
  - PostgreSQL 9.6+ (using lib/pq)
- Git

## Installation Steps

1. Clone the repository:
   ```bash
   git clone [repository-url]
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure your database settings in the environment file (.env):
   ```env
   # Database Configuration
   DB_DRIVER=sqlite    # Options: sqlite, mysql, postgres
   DB_HOST=localhost   # Database host (not needed for SQLite)
   DB_PORT=3306       # Database port (not needed for SQLite)
   DB_DATABASE=cms    # Database name (for SQLite, this is the file path)
   DB_USERNAME=user   # Database username (not needed for SQLite)
   DB_PASSWORD=pass   # Database password (not needed for SQLite)

   # Server Configuration
   SERVER_HOST=localhost
   SERVER_PORT=8080
   APP_URL=http://localhost:8080
   ```

4. Run the CMS:
   ```bash
   go run .
   ```

## Configuration

The CMS uses standard Go database/sql with the goqu query builder. Configuration is done through environment variables in a `.env` file:

### Database Settings
- `DB_DRIVER` - Database driver (sqlite, mysql, postgres)
- `DB_HOST` - Database host (for MySQL/PostgreSQL)
- `DB_PORT` - Database port (for MySQL/PostgreSQL)
- `DB_DATABASE` - Database name or file path for SQLite
- `DB_USERNAME` - Database username (for MySQL/PostgreSQL)
- `DB_PASSWORD` - Database password (for MySQL/PostgreSQL)

### Server Settings
- `SERVER_HOST` - Server host address
- `SERVER_PORT` - Server port number
- `APP_URL` - Full application URL

### Database Connection Examples

1. SQLite:
   ```env
   DB_DRIVER=sqlite
   DB_DATABASE=./data.db
   ```

2. MySQL:
   ```env
   DB_DRIVER=mysql
   DB_HOST=localhost
   DB_PORT=3306
   DB_DATABASE=cms
   DB_USERNAME=user
   DB_PASSWORD=pass
   ```

3. PostgreSQL:
   ```env
   DB_DRIVER=postgres
   DB_HOST=localhost
   DB_PORT=5432
   DB_DATABASE=cms
   DB_USERNAME=user
   DB_PASSWORD=pass
   ```

## Next Steps

- [Architecture Overview](../architecture/overview.md)
- [Quick Start Guide](./quickstart.md)
- [Basic Usage](../guides/basic-usage.md) 