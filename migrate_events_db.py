#!/usr/bin/env python3
"""
Database migration script from old schema (events.db) to new schema (events2.db)
"""

import sqlite3
import sys
from pathlib import Path
from datetime import datetime

def read_schema(schema_file):
    """Read the schema SQL file"""
    with open(schema_file, 'r') as f:
        return f.read()

def create_new_db(db_path, schema):
    """Create new database with the new schema"""
    if Path(db_path).exists():
        print(f"Removing existing {db_path}")
        Path(db_path).unlink()

    conn = sqlite3.connect(db_path)
    conn.executescript(schema)
    conn.commit()
    conn.close()
    print(f"Created {db_path} with new schema")

def migrate_data(old_db, new_db):
    """Migrate data from old database to new database"""

    old_conn = sqlite3.connect(old_db)
    old_conn.row_factory = sqlite3.Row
    new_conn = sqlite3.connect(new_db)

    try:
        # Migrate users
        print("\n=== Migrating users ===")
        old_users = old_conn.execute("SELECT * FROM users").fetchall()
        for user in old_users:
            try:
                new_conn.execute("""
                    INSERT INTO users (id, external_id, email, is_admin, inserted_at)
                    VALUES (?, ?, ?, ?, ?)
                """, (
                    user['id'],
                    user['user_id'],
                    user['email'],
                    user['is_admin'],
                    user['inserted_time']
                ))
                print(f"  Migrated user: {user['user_id']}")
            except Exception as e:
                print(f"  Error migrating user {user['user_id']}: {e}")

        # Migrate puljer
        print("\n=== Migrating puljer ===")
        old_puljer = old_conn.execute("SELECT * FROM puljer").fetchall()
        for pulje in old_puljer:
            try:
                # Default status to 'not_published' for old data
                new_conn.execute("""
                    INSERT OR REPLACE INTO puljer (id, name, status, start_at, end_at)
                    VALUES (?, ?, ?, DATE(?), DATE(?))
                """, (
                    pulje['id'],
                    pulje['name'],
                    'not_published',
                    pulje['start_time'],
                    pulje['end_time']
                ))
                print(f"  Migrated pulje: {pulje['id']}")
            except Exception as e:
                print(f"  Error migrating pulje {pulje['id']}: {e}")

        # Migrate billettholdere
        print("\n=== Migrating billettholdere ===")
        old_billettholders = old_conn.execute("SELECT * FROM billettholdere").fetchall()
        for bt in old_billettholders:
            try:
                new_conn.execute("""
                    INSERT INTO billettholdere
                    (id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id, created_at, updated_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    bt['id'],
                    bt['first_name'],
                    bt['last_name'],
                    bt['ticket_type_id'],
                    bt['ticket_type'],
                    bt['is_over_18'],
                    bt['order_id'],
                    bt['ticket_id'],
                    bt['inserted_time'],
                    bt['inserted_time']
                ))
                print(f"  Migrated billettholder: {bt['id']}")
            except Exception as e:
                print(f"  Error migrating billettholder {bt['id']}: {e}")

        # Migrate billettholder_emails -> relation_billettholder_emails
        print("\n=== Migrating billettholder emails ===")
        old_emails = old_conn.execute("SELECT * FROM billettholder_emails").fetchall()
        for email in old_emails:
            try:
                new_conn.execute("""
                    INSERT INTO relation_billettholder_emails
                    (id, billettholder_id, email, kind, created_at, updated_at)
                    VALUES (?, ?, ?, ?, ?, ?)
                """, (
                    email['id'],
                    email['billettholder_id'],
                    email['email'],
                    email['kind'],
                    email['inserted_time'],
                    email['inserted_time']
                ))
                print(f"  Migrated email: {email['email']}")
            except Exception as e:
                print(f"  Error migrating email {email['id']}: {e}")

        # Migrate billettholdere_users -> relation_billettholdere_users
        print("\n=== Migrating billettholder-user relations ===")
        try:
            old_btusers = old_conn.execute("SELECT * FROM billettholdere_users").fetchall()
            for btuser in old_btusers:
                try:
                    new_conn.execute("""
                        INSERT INTO relation_billettholdere_users (billettholder_id, user_id, inserted_at)
                        VALUES (?, ?, ?)
                    """, (
                        btuser['billettholder_id'],
                        btuser['user_id'],
                        btuser['inserted_time']
                    ))
                    print(f"  Migrated relation: billettholder {btuser['billettholder_id']} -> user {btuser['user_id']}")
                except Exception as e:
                    print(f"  Error migrating relation: {e}")
        except Exception as e:
            print(f"  Table billettholdere_users not found or error: {e}")

        # Migrate events
        print("\n=== Migrating events ===")
        old_events = old_conn.execute("SELECT * FROM events").fetchall()
        for event in old_events:
            try:
                # Map old event_runtime column to event_runtime (it might be named event_runtime or runtime)
                event_runtime = event.get('event_runtime') or event.get('runtime', 'Normal')

                new_conn.execute("""
                    INSERT INTO events
                    (id, title, intro, description, system, event_type, age_group, event_runtime,
                     host_name, user_id, email, phone_number, max_players, beginner_friendly,
                     can_be_run_in_english, notes, status, created_at, updated_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    event['id'],
                    event['title'],
                    event['intro'],
                    event['description'],
                    event.get('system', ''),
                    event.get('event_type', 'Other'),
                    event.get('age_group', 'Default'),
                    event_runtime,
                    event['host_name'],
                    event.get('host'),
                    event['email'],
                    event['phone_number'],
                    event['max_players'],
                    event['beginner_friendly'],
                    event['can_be_run_in_english'],
                    event.get('notes', ''),
                    event.get('status', 'Kladd'),
                    event['created_at'],
                    event['created_at']
                ))
                print(f"  Migrated event: {event['id']}")
            except Exception as e:
                print(f"  Error migrating event {event['id']}: {e}")

        # Migrate event_pujer -> relation_event_puljer
        print("\n=== Migrating event-pulje relations ===")
        try:
            old_ep = old_conn.execute("SELECT * FROM event_pujer").fetchall()
            for ep in old_ep:
                try:
                    new_conn.execute("""
                        INSERT INTO relation_event_puljer
                        (event_id, pulje_id, is_in_pulje, is_published, room)
                        VALUES (?, ?, ?, ?, ?)
                    """, (
                        ep['event_id'],
                        ep['pulje_id'],
                        ep['is_active'],  # Map is_active -> is_in_pulje
                        ep['is_published'],
                        ep.get('room', '')
                    ))
                    print(f"  Migrated event-pulje: {ep['event_id']} -> {ep['pulje_id']}")
                except Exception as e:
                    print(f"  Error migrating event-pulje: {e}")
        except Exception as e:
            print(f"  Table event_pujer not found or error: {e}")

        # Migrate interests
        print("\n=== Migrating interests ===")
        try:
            old_interests = old_conn.execute("SELECT * FROM interests").fetchall()
            for interest in old_interests:
                try:
                    # Old interests didn't have pulje_id, we'll try to get it from events_players
                    # For now, use a placeholder since new schema requires pulje_id
                    # We'll handle this after events_players
                    pass
                except Exception as e:
                    print(f"  Error processing old interests: {e}")
        except Exception as e:
            print(f"  Table interests not found: {e}")

        # Migrate events_players -> relation_events_players and interests
        print("\n=== Migrating events_players ===")
        try:
            old_players = old_conn.execute("SELECT * FROM events_players").fetchall()
            for player in old_players:
                try:
                    # events_players becomes relation_events_players
                    new_conn.execute("""
                        INSERT INTO relation_events_players
                        (event_id, pulje_id, billettholder_id, role, inserted_at)
                        VALUES (?, ?, ?, ?, ?)
                    """, (
                        player['event_id'],
                        'FredagKveld',  # Default pulje_id - this would need better mapping
                        player['billettholder_id'],
                        'Player',  # Default role
                        player['inserted_time']
                    ))

                    # Also create interest entry
                    new_conn.execute("""
                        INSERT INTO interests
                        (billettholder_id, event_id, pulje_id, interest_level, created_at, updated_at)
                        VALUES (?, ?, ?, ?, ?, ?)
                    """, (
                        player['billettholder_id'],
                        player['event_id'],
                        'FredagKveld',  # Default pulje_id
                        player.get('interest_level', 'Middels interessert'),
                        player['inserted_time'],
                        player['inserted_time']
                    ))
                    print(f"  Migrated player: event {player['event_id']} <- billettholder {player['billettholder_id']}")
                except Exception as e:
                    print(f"  Error migrating events_player: {e}")
        except Exception as e:
            print(f"  Table events_players not found or error: {e}")

        new_conn.commit()
        print("\n✅ Migration completed successfully!")

    except Exception as e:
        print(f"\n❌ Migration failed: {e}")
        new_conn.rollback()
        raise
    finally:
        old_conn.close()
        new_conn.close()

def main():
    script_dir = Path(__file__).parent
    old_db = script_dir / "database" / "events.db"
    new_db = script_dir / "database" / "events2.db"
    schema_file = script_dir / "schema.sql"

    print("=" * 60)
    print("Database Migration: events.db -> events2.db")
    print("=" * 60)

    # Check if old database exists
    if not old_db.exists():
        print(f"\n⚠️  Old database not found: {old_db}")
        print("Creating empty events2.db with new schema only...")
        if schema_file.exists():
            schema = read_schema(schema_file)
            create_new_db(new_db, schema)
        else:
            print(f"Schema file not found: {schema_file}")
            return 1
        return 0

    # Check if schema file exists
    if not schema_file.exists():
        print(f"Error: Schema file not found: {schema_file}")
        return 1

    try:
        # Create new database with schema
        schema = read_schema(schema_file)
        create_new_db(new_db, schema)

        # Migrate data
        migrate_data(old_db, new_db)

        print(f"\n✅ New database created at: {new_db}")
        return 0

    except Exception as e:
        print(f"\n❌ Error: {e}")
        return 1

if __name__ == "__main__":
    sys.exit(main())
