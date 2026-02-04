# CMS Store - Active Development Plan

## Current Tasks

### [ ] Implement Shortened IDs for All CMS Entities
Migrate from long numeric IDs (32-char HumanUid) to shortened TimestampMicro IDs (9-char) for all CMS entities while maintaining backward compatibility.

**Affected Entities:**
- Pages
- Blocks
- Menus
- Menu Items
- Sites
- Templates
- Translations

---

## Problem Analysis

### Current State
- IDs generated using `uid.HumanUid()` produce 32-character numeric strings (e.g., `20260116055547619570214289007495`)
- Database column `id` has `Length: 40` (VARCHAR(40))
- IDs stored and transmitted as strings throughout the system
- JavaScript clients receiving these IDs via API may truncate them when parsing as numbers (exceeds `Number.MAX_SAFE_INTEGER`)

### Target State (Dual System)
- **New entities:** Generate with TimestampMicro → shortened to 9 chars (stored in DB as 9 chars)
- **Existing entities:** Keep original 32-char IDs in database (no migration of stored values)
- **API/MCP responses:** Return ALL IDs as shortened (both new 9-char and old 32-char shortened to 21-char)
- **API/MCP inputs:** Accept both short and long formats for lookups
- Support both formats in all API endpoints and database queries
- No breaking changes for existing clients

---

## Recommended ID Generation Method

**TimestampMicro + Crockford Base32 (lowercase)**

- **Base ID:** `uid.TimestampMicro()` produces 16-digit microsecond timestamps (e.g., `1768543534819239`)
- **Shortened Length:** 9 characters (vs 32 for HumanUid)
- **Human-friendly:** Case-insensitive, excludes ambiguous characters (I, L, O, U)
- **Chronologically sortable:** IDs naturally sort by creation time
- **URL-safe:** Lowercase format is optimal for URLs
- **Example:** `86ccrtsgx` (from TimestampMicro `1768543534819239`)
- **Collision-resistant:** Microsecond precision provides sufficient uniqueness for CMS content

### Why TimestampMicro over HumanUid?
- **Shorter:** 16 digits vs 32 digits → 9 chars vs 21 chars when shortened
- **Sufficient precision:** Microseconds are more than adequate for CMS content creation rates
- **Time-ordered:** Natural chronological sorting in databases
- **Efficient:** Smaller storage footprint and faster comparisons

**Implementation Note:** Use `strings.ToLower(uid.ShortenCrockford(uid.TimestampMicro()))` to ensure consistent lowercase output.

---

## Dual System Architecture

### Database Layer
- **New entities:** Store 9-char shortened TimestampMicro IDs directly
- **Existing entities:** Keep original 32-char HumanUid IDs unchanged (no data migration)
- Both ID formats coexist in the same `id` column (VARCHAR(40) accommodates both)

### API/MCP Layer
- **Responses:** Always return shortened IDs
  - New entities: Return 9-char ID as-is
  - Old entities: Dynamically shorten 32-char to 21-char on-the-fly
- **Requests:** Accept both formats
  - Short IDs: Lookup directly
  - Long IDs: Try direct lookup, then try unshortening

### Lookup Logic
```
EntityFindByID(id):
  1. Try direct database lookup (handles both 9-char and 32-char)
  2. If not found and ID looks shortened:
     - Try unshortening and lookup again
  3. Return entity with ID shortened for API response
```

### Benefits
- No database migration needed for existing records
- Gradual transition to shorter IDs
- API consumers always get consistent short IDs
- Backward compatible with old long IDs in requests

---

## Benefits of TimestampMicro

### Storage Efficiency
- 72% reduction in ID length: 32 chars → 9 chars
- Smaller database indexes and faster lookups
- Reduced network payload size

### Performance
- Chronologically ordered IDs improve database B-tree performance
- Natural sorting by creation time (no need for separate created_at index for ordering)
- Faster string comparisons due to shorter length

### Developer Experience
- Human-readable timestamps (can be unshortened for debugging)
- Predictable ID format
- Easy to reason about ID generation

### Uniqueness
- Microsecond precision = 1,000,000 unique IDs per second
- More than sufficient for CMS content creation rates
- Collision probability: virtually zero for typical CMS workloads

---

## Migration Strategy

#### 1. Add ID Helper Functions
- Create centralized ID utilities for shortening/unshortening
- Add `normalizeID()` helper to handle both formats and case normalization

#### 2. Update ID Generation for All Entities
Update the following files to use TimestampMicro:
- `page_implementation.go` - Update `NewPage()`
- `block.go` - Update `NewBlock()`
- `menu_implementation.go` - Update `NewMenu()`
- `menu_item_implementation.go` - Update `NewMenuItem()`
- `site_implementation.go` - Update `NewSite()`
- `template_implementation.go` - Update `NewTemplate()`
- `translation_implementation.go` - Update `NewTranslation()`

Use: `strings.ToLower(uid.ShortenCrockford(uid.TimestampMicro()))` for new IDs (9 chars)

#### 3. Add ID Lookup Support
Update store methods to accept both long and short IDs:
- `store_pages.go` - Update `PageFindByID()`
- `store_blocks.go` - Update `BlockFindByID()`
- `store_menus.go` - Update `MenuFindByID()`
- `store_menu_items.go` - Update `MenuItemFindByID()`
- `store_sites.go` - Update `SiteFindByID()`
- `store_templates.go` - Update `TemplateFindByID()`
- `store_translations.go` - Update `TranslationFindByID()`

Logic:
- Normalize input IDs to lowercase for Crockford lookups
- Try direct lookup first, then try unshortening if not found

### Phase 2: API Response Enhancement (Non-Breaking)

#### Update MCP Response Logic
Update MCP tools to always return shortened IDs:
- `mcp/mcp.go` - Update entity-to-map functions
  - For 9-char IDs (new): Return as-is
  - For 32-char IDs (old): Shorten to 21-char on-the-fly using `ShortenCrockford()`
  - Optionally add `id_original` field for debugging

#### Update MCP Tools
- Accept both short and long ID formats in all tool parameters
- Normalize IDs before database lookups
- Always return shortened IDs in responses

### Phase 3: Database Schema (No Changes Required)

#### Verify ID Column Capacity
All entity table creation files already use `Length: 40` for VARCHAR(40):
- Sufficient for both 9-char (new) and 32-char (old) IDs
- No schema changes needed
- No data migration required

### Phase 4: Testing & Validation

#### Add Comprehensive Tests
For each entity type, test:
- New entity creation with 9-char TimestampMicro IDs
- Lookup by short ID (9-char and 21-char)
- Lookup by long ID (32-char backward compatibility)
- API responses return shortened IDs for both old and new entities
- Versioning with both ID formats
- Edge cases (empty, invalid, malformed IDs)
- On-the-fly shortening of old 32-char IDs in responses

#### Update Documentation
- Document the dual ID system architecture
- Explain that new entities use 9-char IDs, old entities keep 32-char in DB
- Document that APIs always return shortened IDs
- Provide examples of both ID formats and lookup behavior

---

## Implementation Details

### Files to Modify

**Core Entity Files:**
- `page_implementation.go`
- `block.go`
- `menu_implementation.go`
- `menu_item_implementation.go`
- `site_implementation.go`
- `template_implementation.go`
- `translation_implementation.go`

**Store Files:**
- `store_pages.go`
- `store_blocks.go`
- `store_menus.go`
- `store_menu_items.go`
- `store_sites.go`
- `store_templates.go`
- `store_translations.go`

**MCP Files:**
- `mcp/mcp.go`

**Test Files:**
- Add tests for each entity type

### New Files to Create
- `id_helpers.go` - Centralized ID shortening/unshortening utilities

---

## Backward Compatibility Guarantees

✅ **Existing long IDs continue to work** - All lookup methods accept both 9-char and 32-char formats  
✅ **API responses are consistent** - All IDs returned as shortened (9-char for new, 21-char for old)  
✅ **Database unchanged** - No migration of existing 32-char IDs required  
✅ **No forced migration** - Old 32-char IDs remain in database indefinitely  
✅ **Gradual transition** - New entities use 9-char IDs, old entities keep 32-char in DB  
✅ **Transparent to clients** - API consumers see shortened IDs regardless of storage format

---

## Configuration Options

Add optional configuration to `NewStoreOptions`:

```go
UseShortIDs bool        // Enable shortened IDs for new records (default: true)
IDGenerator string      // "timestamp_micro", "nano_uid", "micro_uid", "human_uid" (default: "timestamp_micro")
ShortIDMethod string    // "crockford", "base58", "base62" (default: "crockford")
ShortIDLowercase bool   // Use lowercase for Crockford (default: true for URL-friendliness)
```

---

## Rollout Plan

1. **Development:** Implement Phase 1-3, deploy to dev environment
2. **Testing:** Run comprehensive tests, validate dual system and backward compatibility
3. **Staging:** Deploy to staging, test with real data (mix of old and new IDs)
4. **Production:** Deploy all phases (new entities use 9-char IDs, APIs return shortened IDs)
5. **Monitor:** Track API usage, verify both ID formats work correctly
6. **Optimize:** After validation period, optimize queries and consider future cleanup

---

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| ID collisions with TimestampMicro | Microsecond precision (1,000,000 IDs/second) is sufficient for CMS content; add collision detection if needed |
| Case sensitivity issues | Store IDs in lowercase, normalize input to lowercase for lookups |
| Performance impact of dual lookup | Optimize query logic, leverage existing indexes |
| Time-based ID predictability | Not a security concern for CMS content; use SecUid/NanoUid if security-sensitive |
| Client confusion with two IDs | Clear documentation, deprecation notices for long IDs |
| Versioning system compatibility | Test versioning with both ID formats, ensure entity tracking works |

---

## Success Criteria

✅ All new entities stored with 9-char TimestampMicro IDs in database  
✅ All existing entities keep 32-char IDs in database (no migration)  
✅ API responses return shortened IDs for ALL entities (9-char for new, 21-char for old)  
✅ API inputs accept both short and long ID formats  
✅ IDs are chronologically sortable (TimestampMicro ordering for new entities)  
✅ No JavaScript integer truncation issues  
✅ All tests pass (unit, integration, backward compatibility, dual system)  
✅ Zero downtime deployment  
✅ Performance metrics remain stable or improve

---

## Timeline Estimate

- **Phase 1:** 4-5 hours (core implementation - ID generation and lookup for all entities)
- **Phase 2:** 2-3 hours (API response shortening logic)
- **Phase 3:** 1 hour (schema verification - no changes needed)
- **Phase 4:** 4-5 hours (testing + documentation for all entities)
- **Total:** ~11-14 hours development + testing time

---

## Next Steps

- [ ] Review and approve this plan
- [ ] Create feature branch: `feature/shortened-ids`
- [ ] Implement Phase 1 (dual ID support for all entities)
- [ ] Implement Phase 2 (API response shortening)
- [ ] Write comprehensive tests
- [ ] Code review and iterate
- [ ] Deploy to staging for validation
- [ ] Production deployment with monitoring