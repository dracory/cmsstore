# Future Ideas

1.  **Foreign Key Constraints:** Add foreign key constraints to enforce referential integrity between tables (e.g., `site_id` in `block`, `page`, `template`, `menu`, and `translation` tables; `template_id` in `block` and `page` tables; `menu_id` in `menu_item` table; `page_id` in `block` and `menu_item` tables).
2.  **Project Structure:** Optimize the project structure by:
    *   Organizing files into separate packages (e.g., `models/`, `queries/`, `db/`).
    *   Extracting common components from the `admin/` directory into a shared package.
    *   Ensuring consistent file naming conventions.
    *   Reviewing and updating dependencies in `go.mod`.
3.  **Domain Names:** Normalize the `domain_names` column in the `site` table into a separate table.