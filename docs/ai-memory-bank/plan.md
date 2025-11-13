# CMS Store Improvement Plan

Based on the analysis of the project structure, database schema, store initialization, and query implementations, I have identified several potential improvements:

1.  **Automated Schema Management:** Enhance the existing `AutoMigrate` function to automatically detect and apply schema changes.
2.  **Middleware Management:** Use a more structured approach for managing middlewares, instead of storing them as serialized text.
3.  **Code Generation:** Use code generation to reduce boilerplate code in the query structs.

## [2025-11-13] Site order_by validation error
- [ ] Investigate where `siteQuery` is created without an order_by value
- [ ] Implement validation or defaulting to avoid empty `order_by` values
- [ ] Verify the behaviour via tests or reasoning

## [2025-11-13] Translation implementation tests
- [X] Determine behaviours to cover for translation implementation
- [X] Add unit tests for translation implementation setters/getters
- [X] Document or share results

## [2025-11-13] REST API tests warnings
- [ ] Review warnings about response usage before error checks
- [ ] Refactor rest tests to assert on request/response errors first
- [ ] Confirm tests compile and run cleanly

## [2025-11-13] Add domain button bug
- [ ] Reproduce add domain button failure and capture console/network errors
- [ ] Trace frontend logic to ensure SweetAlert dependency loads before usage
- [ ] Verify backend site update controller handles domain addition request