# Database Migration: Old Schema → New Schema

## Overview

The application has been upgraded from the old database schema to a new schema with improved structure and audit tracking capabilities.

### What Changed

**New Features in Schema:**
- Added audit tracking fields to all tables: `created_by_id`, `updated_by_id`
- Status tracking for events: `status_changed_by_id`, `status_changed_at`, `status_changed_action`
- Status field for puljer with CHECK constraint
- Better foreign key relationships and constraints
- Renamed tables and columns for consistency:
  - `billettholder_emails` → `relation_billettholder_emails`
  - `billettholdere_users` → `relation_billettholdere_users`
  - `event_pujer` → `relation_event_puljer`
  - `events_players` → `relation_events_players`
  - `user_id` (in users table) → `external_id`
  - `inserted_time` → `created_at`, `updated_at`
  - `start_time`/`end_time` → `start_at`/`end_at`

## Migration Instructions

### Step 1: Create New Database with New Schema

The new database (`events2.db`) has been created with the new schema defined in `schema.sql`.

### Step 2: Migrate Data (If You Have Existing Data)

If you have an existing `database/events.db` file:

```bash
python3 migrate_events_db.py
```

This script will:
1. Create `events2.db` with the new schema
2. Automatically migrate data from `events.db` to `events2.db`
3. Handle mapping between old and new column names
4. Set default values for new columns that don't have mappings in the old schema

### Step 3: Update Configuration

Update your application configuration to use the new database:

```go
// In main.go or your config
dsn := flag.String("dbp", "database/events2.db", "absolute path to database file")
```

### Step 4: Add Reference Data

If migrating from scratch, you can populate reference tables by running:

```bash
python3 << 'EOF'
import sqlite3

conn = sqlite3.connect("database/events2.db")
cursor = conn.cursor()

# Insert event statuses
statuses = [('Kladd',), ('Innsendt',), ('Godkjent',), ('Forkastet',)]
cursor.executemany("INSERT OR IGNORE INTO event_statuses (status) VALUES (?)", statuses)

# Insert event types
types = [('Roleplay',), ('Boardgame',), ('Cardgame',), ('Other',)]
cursor.executemany("INSERT OR IGNORE INTO events_types (event_type) VALUES (?)", types)

# Insert age groups
ages = [('Default',), ('ChildFriendly',), ('AdultsOnly',)]
cursor.executemany("INSERT OR IGNORE INTO age_groups (age_group) VALUES (?)", ages)

# Insert runtimes
runtimes = [('Normal',), ('ShortRunning',), ('LongRunning',)]
cursor.executemany("INSERT OR IGNORE INTO event_runtimes (runtime) VALUES (?)", runtimes)

# Insert interest levels
levels = [('Litt interessert',), ('Middels interessert',), ('Veldig interessert',)]
cursor.executemany("INSERT OR IGNORE INTO interest_levels (interest_level) VALUES (?)", levels)

# Insert pulje statuses
pulje_statuses = [('not_published',), ('published',), ('locked',), ('completed',)]
cursor.executemany("INSERT OR IGNORE INTO pulje_statuses (status) VALUES (?)", pulje_statuses)

conn.commit()
conn.close()
print("✅ Reference data inserted")
EOF
```

## Migration Script Details

### What Gets Migrated

1. **Users** - All users are migrated with `user_id` mapped to `external_id`
2. **Billettholders** - All billing holders with audit timestamps
3. **Billettholder Emails** - Email associations with audit tracking
4. **Puljer** - All pulse/schedule entries (status defaults to `not_published`)
5. **Events** - All events (old `host` becomes `user_id`)
6. **Event-Pulje Relations** - Old `event_pujer` becomes `relation_event_puljer`
7. **Events Players** - Old data becomes `relation_events_players` with role `Player`
8. **Interests** - Event interests with pulje tracking (pulje_id mapped from events_players)

### Default Values for New Fields

Fields with no data in the old schema get defaults:

| Field | Default | Reason |
|-------|---------|--------|
| `created_by_id` | NULL | No user tracking in old system |
| `updated_by_id` | NULL | No user tracking in old system |
| `status_changed_by_id` | NULL | No user tracking in old system |
| `status_changed_at` | NULL | No tracking in old system |
| `status_changed_action` | NULL | No action tracking in old system |
| Pulje `status` | `'not_published'` | Default safe state |
| `event_runtime` | `'Normal'` | Default runtime type |
| `age_group` | `'Default'` | Default age group |
| `event_type` | `'Other'` | Default event type |

## Verification

To verify the migration:

```bash
python3 << 'EOF'
import sqlite3

conn = sqlite3.connect("database/events2.db")
cursor = conn.cursor()

# Count all tables
cursor.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;")
tables = cursor.fetchall()

print("Database tables and row counts:")
for table in tables:
    cursor.execute(f"SELECT COUNT(*) FROM {table[0]}")
    count = cursor.fetchone()[0]
    print(f"  {table[0]}: {count} rows")

conn.close()
EOF
```

## Rollback (If Needed)

If something goes wrong, simply:

1. Delete `database/events2.db`
2. Your original `database/events.db` is untouched
3. Adjust the code and try again

## Application Code Updates

The application code has been updated to use the new models:

- ✅ Billettholder models now include `created_by_id` and `updated_by_id`
- ✅ BillettholderEmail models include new audit fields
- ✅ Event models include status tracking fields
- ✅ PuljeRow includes name and status fields
- ✅ All database queries updated to handle new fields

## Next Steps

1. **Test the application** with the new database
2. **Verify all queries** work with the new schema
3. **Deploy** once testing is complete
4. **Backup** your old database before complete migration

## Support

If you encounter issues during migration:

1. Check the migration script output for specific errors
2. Verify both databases exist in `database/` directory
3. Ensure schema.sql is in the repository root
4. Run migration script in verbose mode for debugging
