# CMS Store Improvement Plan

Based on the analysis of the project structure, database schema, store initialization, and query implementations, I have identified several potential improvements:

1.  **Automated Schema Management:** Enhance the existing `AutoMigrate` function to automatically detect and apply schema changes.
2.  **Middleware Management:** Use a more structured approach for managing middlewares, instead of storing them as serialized text.
3.  **Code Generation:** Use code generation to reduce boilerplate code in the query structs.