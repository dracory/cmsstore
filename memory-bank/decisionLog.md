## [2025-02-23] - Constant Refactoring

**Context:** Hardcoded strings were used in multiple query files. This made the code less maintainable and harder to update.

**Decision:** Refactor the code to use constants defined in `consts.go`.

**Rationale:** Using constants improves code readability, maintainability, and reduces the risk of errors.  Changes are centralized in one location.

**Implementation:**  Moved all relevant strings to `consts.go` and updated the affected files (`site_query.go`, `page_query.go`, `menu_query.go`, `template_query.go`). Added comments to `consts.go` for better understanding.