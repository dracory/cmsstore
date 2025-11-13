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
- [ ] Determine behaviours to cover for translation implementation
- [ ] Add unit tests for translation implementation setters/getters
- [ ] Document or share results