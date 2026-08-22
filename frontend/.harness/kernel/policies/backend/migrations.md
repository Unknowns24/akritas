# Backend Database Migrations Policy

## Purpose

Relational database schema and persistent data changes must be explicit, versioned, ordered and reviewable.

For Go backends using GORM, migrations use `github.com/go-gormigrate/gormigrate/v2` and live in the database adapter. Runtime model synchronization must not replace migration history.

## Required structure

For each relational database adapter, use the following structure unless an existing project already has an equivalent migration layout:

```text
internal/adapter/db/<technology>/
├── migrations/
│   ├── migrate.go
│   ├── schema/
│   └── data/
├── repository/
└── ...
```

Responsibilities:

- `migrations/schema/`: schema evolution such as tables, columns, indexes, constraints and renames.
- `migrations/data/`: seeds, backfills, normalization and other persistent data transformations.
- `migrations/migrate.go`: the single ordered registry of migrations executed by the application migration runner.

Do not place migration logic in repositories, handlers, usecases, domain entities or application services.

## Migration engine

When GORM is the persistence ORM:

- Use `github.com/go-gormigrate/gormigrate/v2`.
- Every migration must be represented by a `*gormigrate.Migration`.
- Every migration must have an immutable, unique `ID`.
- Migrations must be registered explicitly in `migrations/migrate.go` in execution order.
- The migration runner must return errors to the caller. Migration failures must never be silently ignored.

Example registry shape:

```go
func Run(db *gorm.DB) error {
    m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        schema.SCHEMA_20260820_01_AddProjects(),
        data.DATA_20260820_10_BackfillProjects(),
    })

    return m.Migrate()
}
```

## Naming and ordering

Migration IDs and filenames must encode their chronological order.

Use:

```text
YYYYMMDD_NN_<snake_case_description>
```

Examples:

```text
20260820_01_add_projects
20260820_02_add_project_indexes
20260820_10_seed_default_projects
20260820_11_backfill_project_owners
```

Conventions:

- `01`-`09`: schema migrations for that date.
- `10` and above: data migrations for that date.
- A data migration that depends on a schema migration must be registered after that schema migration.
- The Go constructor should make the migration kind, date and order obvious, for example:
  - `SCHEMA_20260820_01_AddProjects()`
  - `DATA_20260820_10_BackfillProjects()`
- The migration `ID` is permanent after it has been applied anywhere outside an ephemeral local database.

Never edit, rename, reorder or reuse the ID of an already-applied migration to represent a new change. Create a new migration instead.

## Schema migrations

Schema changes must be implemented inside `migrations/schema/`.

Allowed operations include GORM migrator operations, explicit SQL when required, and `AutoMigrate` scoped to the versioned migration.

Example:

```go
func SCHEMA_20260820_01_AddProjects() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "20260820_01_add_projects",
        Migrate: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&domain.Project{})
        },
        Rollback: func(tx *gorm.DB) error {
            if tx.Migrator().HasTable(&domain.Project{}) {
                return tx.Migrator().DropTable(&domain.Project{})
            }
            return nil
        },
    }
}
```

Rules:

- Prefer explicit checks such as `HasTable`, `HasColumn` and `HasIndex` when they make a migration safer across real existing databases.
- Preserve existing data when renaming or replacing columns unless destructive behavior is an explicit approved requirement.
- Destructive operations require special care and must be visible in the implementation brief/TDD plan.
- Do not rely on current model shape alone when a historical migration needs precise SQL or compatibility behavior.

## Data migrations

Persistent data transformations must be implemented inside `migrations/data/`.

Use data migrations for:

- deterministic seed data required by the application;
- backfilling newly introduced fields or tables;
- normalizing existing records;
- migrating data from a legacy representation to a new one.

Rules:

- A data migration must operate through the transaction/database handle supplied by gormigrate.
- It must not call HTTP APIs, message brokers or other non-database side effects.
- It must be deterministic for the database state it is designed to migrate.
- It must not depend on REST handlers, usecases or concrete repositories merely to transform persisted data.
- For irreversible migrations, `Rollback` may return `nil`, but irreversibility must be intentional and obvious in code/review.

## Rollback

Every migration must declare a `Rollback` function.

The rollback must reverse the migration when doing so is safe and meaningful.

If a data transformation cannot be safely reconstructed, an explicit no-op rollback is acceptable:

```go
Rollback: func(tx *gorm.DB) error { return nil },
```

Do not implement a misleading rollback that can corrupt or discard valid production data.

## AutoMigrate policy

Global or startup-wide GORM `AutoMigrate` is forbidden as a schema management strategy.

Forbidden:

```go
db.AutoMigrate(
    &domain.User{},
    &domain.Project{},
    &domain.Order{},
)
```

The database connection initializer must open/configure the connection only. Schema evolution belongs to ordered migrations.

`AutoMigrate` is allowed only inside a specific versioned schema migration when its behavior is appropriate for that change.

This guarantees that a deployment cannot acquire untracked schema changes merely because a domain struct changed.

## Domain and adapter boundaries

Migration code belongs to the DB adapter and may depend on GORM/gormigrate.

Domain entities must remain independent from migration infrastructure:

- `internal/core/**` must never import GORM or gormigrate.
- Migration packages may import domain models when GORM metadata on those models is the established project convention.
- Migration-only helper structs may be defined inside the migration package when historical schema shape differs from the current domain model.

## Testing requirements

Migration tests are required when the migration contains non-trivial behavior or risk, including:

- data backfills;
- data normalization;
- column renames/copies;
- constraint/index transitions;
- compatibility with legacy schema states;
- destructive or conditionally destructive operations.

Tests should prove the relevant before/after state, not merely that the migration returns `nil`.

Where practical, test rollback behavior for reversible schema migrations.

A trivial additive migration may rely on the project's integration/migration test coverage if a dedicated test adds no meaningful confidence.

## Existing projects

Before adding a migration:

1. Inspect the existing migration registry and naming convention.
2. Determine the latest migration IDs and ordering.
3. Reuse the project's established DB technology and adapter location.
4. Add the new migration without rewriting historical migrations.
5. Register it in the central migration registry.
6. Verify that application startup/deployment actually invokes the migration runner through the existing wiring.

If the project currently mixes global `AutoMigrate` with versioned migrations, do not silently redesign the entire database lifecycle as part of an unrelated feature. New schema changes must still be versioned; migration of the legacy bootstrap behavior should be proposed explicitly if it is outside task scope.

## Review checklist

A backend change that modifies persistence is incomplete if any applicable item fails:

- schema/data change has a versioned migration;
- migration is in the correct `schema/` or `data/` package;
- filename and immutable ID follow ordering conventions;
- migration is explicitly registered in `migrate.go`;
- schema migration precedes dependent data migration;
- rollback is present and safe/intentional;
- no new global `AutoMigrate` was introduced;
- risky migration behavior has tests;
- migration errors propagate to startup/deployment instead of being ignored.
