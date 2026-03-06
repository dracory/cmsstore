# Core Store Coverage Improvement Plan

**Current Coverage**: 61.5%  
**Target Coverage**: 75%+  
**Required Improvement**: +13.5% or more

## Overview

This plan outlines a systematic approach to increase test coverage for the core CMS store functionality from 61.5% to at least 75%. The plan focuses on the most impactful areas first and provides a phased implementation strategy.

## Current State Analysis

Based on coverage analysis, the main gaps are:
- Many getter methods have 0% coverage
- Query validation functions have partial coverage (50-70%)
- Several store operations have 0% coverage
- Edge cases and error conditions are under-tested

## Implementation Plan

### Phase 1: Foundation Tests (Week 1-2) - Target: +5-7% Coverage

#### Entity Getter Methods (Priority: HIGH)

**Block Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `CreatedAt()` - 0% coverage
- [ ] `CreatedAtCarbon()` - 0% coverage
- [ ] `Content()` - 0% coverage
- [ ] `Editor()` - 0% coverage
- [ ] `Memo()` - 0% coverage
- [ ] `Name()` - 0% coverage
- [ ] `PageID()` - 0% coverage
- [ ] `ParentID()` - 0% coverage
- [ ] `Sequence()` - 0% coverage
- [ ] `SequenceInt()` - 0% coverage
- [ ] `SiteID()` - 0% coverage
- [ ] `TemplateID()` - 0% coverage
- [ ] `Type()` - 0% coverage
- [ ] `UpdatedAt()` - 0% coverage
- [ ] `UpdatedAtCarbon()` - 0% coverage

**Menu Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `Memo()` - 0% coverage
- [ ] `Name()` - 0% coverage
- [ ] `SiteID()` - 0% coverage

**Menu Item Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `Content()` - 0% coverage
- [ ] `Editor()` - 0% coverage
- [ ] `Handle()` - 0% coverage
- [ ] `Memo()` - 0% coverage
- [ ] `PageID()` - 0% coverage
- [ ] `ParentID()` - 0% coverage
- [ ] `Target()` - 0% coverage
- [ ] `URL()` - 0% coverage

**Page Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `Memo()` - 0% coverage
- [ ] `Name()` - 0% coverage
- [ ] `SiteID()` - 0% coverage
- [ ] `TemplateID()` - 0% coverage

**Site Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `Memo()` - 0% coverage
- [ ] `Name()` - 0% coverage

**Template Implementation Tests**
- [ ] `IsActive()` - 0% coverage
- [ ] `IsInactive()` - 0% coverage
- [ ] `Name()` - 0% coverage
- [ ] `SiteID()` - 0% coverage

**Query Validation Functions**
- [ ] `BlockQuery.Validate()` - 52.6% coverage
- [ ] `MenuQuery.Validate()` - 52.2% coverage
- [ ] `MenuItemQuery.Validate()` - 52.4% coverage
- [ ] `PageQuery.Validate()` - 52.2% coverage
- [ ] `SiteQuery.Validate()` - 52.2% coverage
- [ ] `TemplateQuery.Validate()` - 52.2% coverage

### Phase 2: Query Coverage (Week 3-4) - Target: +4-6% Coverage

#### Query Parameter Methods

**Block Query Methods**
- [ ] `SetColumns()` - 0% coverage
- [ ] `CreatedAtGte()` - 0% coverage
- [ ] `CreatedAtLte()` - 0% coverage
- [ ] `IDIn()` - 0% coverage
- [ ] `NameLike()` - 0% coverage
- [ ] `Offset()` - 0% coverage
- [ ] `OrderBy()` - 0% coverage
- [ ] `PageID()` - 0% coverage
- [ ] `ParentID()` - 0% coverage
- [ ] `Sequence()` - 0% coverage
- [ ] `SiteID()` - 0% coverage
- [ ] `SortOrder()` - 0% coverage
- [ ] `Status()` - 0% coverage
- [ ] `StatusIn()` - 0% coverage
- [ ] `TemplateID()` - 0% coverage
- [ ] `SetCountOnly()` - 0% coverage
- [ ] `SetIDIn()` - 0% coverage
- [ ] `SetNameLike()` - 0% coverage
- [ ] `SetOffset()` - 0% coverage
- [ ] `SetOrderBy()` - 0% coverage
- [ ] `SetPageID()` - 0% coverage
- [ ] `SetParentID()` - 0% coverage
- [ ] `SetSequence()` - 0% coverage
- [ ] `SetSiteID()` - 0% coverage
- [ ] `SetSortOrder()` - 0% coverage
- [ ] `SetStatus()` - 0% coverage
- [ ] `SetStatusIn()` - 0% coverage
- [ ] `SetTemplateID()` - 0% coverage

**Similar coverage gaps exist for:**
- [ ] Menu Query methods
- [ ] Menu Item Query methods
- [ ] Page Query methods
- [ ] Site Query methods
- [ ] Template Query methods

### Phase 3: Store Operations (Week 5-6) - Target: +3-5% Coverage

#### Count Operations (0% coverage)
- [ ] `BlockCount()`
- [ ] `MenuCount()`
- [ ] `PageCount()`
- [ ] `SiteCount()`
- [ ] `TemplateCount()`

#### Delete Operations (0% coverage)
- [ ] `BlockDelete()`
- [ ] `MenuDelete()`
- [ ] `PageDelete()`
- [ ] `SiteDelete()`
- [ ] `TemplateDelete()`

#### Versioning Operations (0% coverage)
- [ ] `VersioningDelete()`
- [ ] `VersioningSoftDelete()`
- [ ] `VersioningUpdate()`

### Phase 4: Edge Cases & Integration (Week 7-8) - Target: +1-3% Coverage

#### Error Handling and Edge Cases
- [ ] Test invalid input parameters
- [ ] Test boundary conditions
- [ ] Test error message accuracy
- [ ] Test concurrent access patterns
- [ ] Test database constraint violations
- [ ] Test transaction rollback scenarios

#### Integration Tests
- [ ] Complex multi-entity operations
- [ ] Cross-entity relationship tests
- [ ] Performance under load
- [ ] Memory usage validation

## Expected Coverage Impact

| Phase | Expected Coverage Increase | Cumulative Coverage |
|-------|---------------------------|-------------------|
| Phase 1 | +5-7% | 66.5-68.5% |
| Phase 2 | +4-6% | 70.5-74.5% |
| Phase 3 | +3-5% | 73.5-79.5% |
| Phase 4 | +1-3% | 74.5-82.5% |

**Target Achievement**: 75%+ coverage by end of Phase 3, with potential to reach 82.5% by end of Phase 4.

## Test File Organization

### Entity Tests
- `block_implementation_test.go`
- `menu_implementation_test.go`
- `menu_item_implementation_test.go`
- `page_implementation_test.go`
- `site_implementation_test.go`
- `template_implementation_test.go`
- `translation_implementation_test.go`

### Query Tests
- `block_query_test.go`
- `menu_query_test.go`
- `menu_item_query_test.go`
- `page_query_test.go`
- `site_query_test.go`
- `template_query_test.go`
- `translation_query_test.go`

### Store Tests
- `store_blocks_test.go`
- `store_menus_test.go`
- `store_menu_items_test.go`
- `store_pages_test.go`
- `store_sites_test.go`
- `store_templates_test.go`
- `store_translations_test.go`
- `store_versioning_test.go`

### Integration Tests
- `store_test.go`
- `integration_test.go`

## Testing Principles

### 1. Comprehensive Method Testing
- [ ] Test all public methods with positive cases
- [ ] Test all public methods with negative cases
- [ ] Test all parameter validation
- [ ] Test all error conditions

### 2. Edge Case Coverage
- [ ] Empty string inputs
- [ ] Null/nil values
- [ ] Boundary values
- [ ] Invalid enum values
- [ ] Very large inputs
- [ ] Very small inputs

### 3. Error Handling
- [ ] Verify error messages are accurate
- [ ] Test error code consistency
- [ ] Test error propagation
- [ ] Test error recovery

### 4. Concurrency Testing
- [ ] Test concurrent read operations
- [ ] Test concurrent write operations
- [ ] Test race conditions
- [ ] Test deadlocks

### 5. Integration Testing
- [ ] Test entity relationships
- [ ] Test transaction boundaries
- [ ] Test data consistency
- [ ] Test performance characteristics

## Implementation Guidelines

### 1. Test Structure
```go
func TestEntityMethod(t *testing.T) {
    tests := []struct {
        name     string
        setup    func(*Entity)
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 2. Database Setup
- Use the existing test utilities in `testutils/`
- Create isolated test databases
- Clean up after each test
- Use transactions for rollback capability

### 3. Test Data Management
- Use factory functions for test data creation
- Ensure test data is deterministic
- Clean up test data after tests
- Use realistic test data

### 4. Performance Considerations
- Keep individual tests fast (< 100ms)
- Use parallel test execution where possible
- Minimize database operations in tests
- Use in-memory databases for unit tests

## Progress Tracking

### Weekly Checkpoints
- [ ] Week 1: Complete Phase 1 foundation tests
- [ ] Week 2: Validate Phase 1 coverage improvement
- [ ] Week 3: Complete Phase 2 query coverage
- [ ] Week 4: Validate Phase 2 coverage improvement
- [ ] Week 5: Complete Phase 3 store operations
- [ ] Week 6: Validate Phase 3 coverage improvement
- [ ] Week 7: Complete Phase 4 edge cases
- [ ] Week 8: Final validation and optimization

### Coverage Validation
```bash
# Run coverage after each phase
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Success Criteria

### Minimum Success (75% coverage)
- [ ] Achieve 75%+ coverage by end of Phase 3
- [ ] All critical entity methods tested
- [ ] All query validation functions tested
- [ ] All store operations tested

### Optimal Success (80%+ coverage)
- [ ] Achieve 80%+ coverage by end of Phase 4
- [ ] Comprehensive edge case coverage
- [ ] Robust integration tests
- [ ] Performance and concurrency validation

### Quality Metrics
- [ ] All tests pass consistently
- [ ] Test execution time is reasonable
- [ ] Test code follows project standards
- [ ] Test documentation is complete

## Risk Mitigation

### Potential Challenges
1. **Complex Dependencies**: Use mocking and test doubles
2. **Database Complexity**: Use in-memory databases for unit tests
3. **Test Maintenance**: Keep tests simple and focused
4. **Performance Impact**: Optimize test execution time

### Mitigation Strategies
- [ ] Use table-driven tests for maintainability
- [ ] Implement test utilities for common patterns
- [ ] Use parallel test execution
- [ ] Regular test refactoring and cleanup
- [ ] Continuous coverage monitoring

## Next Steps

1. **Start with Phase 1**: Focus on entity getter methods
2. **Set up test infrastructure**: Ensure test utilities are ready
3. **Implement incrementally**: Add tests one method at a time
4. **Validate progress**: Run coverage after each batch of tests
5. **Adjust plan**: Modify approach based on findings and challenges

This plan provides a clear roadmap to achieve the 75% coverage target while maintaining code quality and test effectiveness.