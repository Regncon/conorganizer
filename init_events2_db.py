#!/usr/bin/env python3
"""
Initialize events2.db with reference data
"""

import sqlite3
import sys
from pathlib import Path

def init_reference_data(db_path):
    """Initialize all reference data tables"""

    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    try:
        print(f"Initializing reference data in {db_path}")

        # Insert event statuses
        print("  - Adding event statuses...")
        statuses = [
            ('Kladd',),
            ('Innsendt',),
            ('Godkjent',),
            ('Forkastet',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO event_statuses (status) VALUES (?)", statuses)

        # Insert event types
        print("  - Adding event types...")
        types = [
            ('Roleplay',),
            ('Boardgame',),
            ('Cardgame',),
            ('Other',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO events_types (event_type) VALUES (?)", types)

        # Insert age groups
        print("  - Adding age groups...")
        ages = [
            ('Default',),
            ('ChildFriendly',),
            ('AdultsOnly',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO age_groups (age_group) VALUES (?)", ages)

        # Insert runtimes
        print("  - Adding event runtimes...")
        runtimes = [
            ('Normal',),
            ('ShortRunning',),
            ('LongRunning',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO event_runtimes (runtime) VALUES (?)", runtimes)

        # Insert interest levels
        print("  - Adding interest levels...")
        levels = [
            ('Litt interessert',),
            ('Middels interessert',),
            ('Veldig interessert',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO interest_levels (interest_level) VALUES (?)", levels)

        # Insert pulje statuses
        print("  - Adding pulje statuses...")
        pulje_statuses = [
            ('not_published',),
            ('published',),
            ('locked',),
            ('completed',)
        ]
        cursor.executemany("INSERT OR IGNORE INTO pulje_statuses (status) VALUES (?)", pulje_statuses)

        conn.commit()
        print("\n✅ Reference data initialized successfully!")
        return 0

    except Exception as e:
        print(f"\n❌ Error initializing reference data: {e}")
        conn.rollback()
        return 1
    finally:
        conn.close()

def main():
    script_dir = Path(__file__).parent
    db_path = script_dir / "database" / "events2.db"

    if not db_path.exists():
        print(f"Error: Database not found at {db_path}")
        print("Please run migrate_events_db.py first to create the database")
        return 1

    return init_reference_data(db_path)

if __name__ == "__main__":
    sys.exit(main())
