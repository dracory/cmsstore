## [2025-02-23] - Constant Refactoring
**Context:** Hardcoded strings were used in multiple query files. This made the code less maintainable and harder to update.

**Decision:** Refactor the code to use constants defined in `consts.go`.

**Rationale:** Using constants improves code readability, maintainability, and reduces the risk of errors. Changes are centralized in one location.

**Implementation:** Moved all relevant strings to `consts.go` and updated the affected files (`site_query.go`, `page_query.go`, `menu_query.go`, `template_query.go`). Added comments to `consts.go` for better understanding.

## [2025-02-23] - Shortcode Analysis
**Context:** User requested information on shortcodes.
**Decision:** Analyzed shortcode implementation and usage.
**Rationale:** To understand the system's architecture and functionality.
**Implementation:** No code changes were made. Findings documented in Memory Bank.

## [2025-02-23] - Multisite Architecture Investigation
**Context:** User requested investigation into how multisites work.

**Decision:** The application uses a multisite architecture.  The README was updated with a user-friendly description.

**Rationale:**  Analysis of `site.go`, `store_sites.go`, and `admin/sites` directory controllers confirmed the presence of a multisite implementation.  Each site is stored as a `Site` struct, allowing for multiple domains per site.  The README now contains a concise explanation of the multisite capabilities.

**Implementation:**  The system uses a `Site` struct and associated controllers for creating, updating, deleting, and managing sites.  A `SiteQuery` allows for filtering and sorting of sites.  The README now includes a user-friendly description of the multisite features.