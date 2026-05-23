# Database Views

This file documents the views defined in `schema.sql`.

Boolean values are stored as `INTEGER` with `0`/`1`. Date/time values are stored as `TEXT` in the `DBDateTime` format. Room fields come from a `LEFT JOIN rooms`, so they can be `NULL` when no room is assigned.

---

# Database View: `v_get_user_billettholder`

**Purpose:** Maps application users to billettholder records through `relation_billettholdere_users`.

**Columns:**

| Column | Type | Description |
| --- | --- | --- |
| `user_id` | `INTEGER` | Internal `users.id`. |
| `external_id` | `TEXT` | External auth/provider user id. |
| `user_email` | `TEXT` | User email. |
| `user_is_admin` | `INTEGER` | Admin flag, `0` or `1`. |
| `user_inserted_at` | `TEXT` | User insert timestamp. |
| `billettholder_id` | `INTEGER` | Linked billettholder id. |
| `billettholder_user_id` | `INTEGER` | Linked `users.id` from the relation table. |
| `billettholder_user_inserted_at` | `TEXT` | Relation insert timestamp. |

**Go usage example:**

```go
rows, err := db.Query("SELECT billettholder_id FROM v_get_user_billettholder WHERE external_id = $1", externalID)
```

---

# Database View: `v_events_by_pulje_active`

**Purpose:** Approved events joined with active pulje metadata and optional room metadata for event listings.

**Filters:**

- `events.status = 'Godkjent'`
- `relation_event_puljer.is_in_pulje = 1`

**Columns:**

| Column | Type | Description |
| --- | --- | --- |
| `id` | `TEXT` | Event id. |
| `title` | `TEXT` | Event title. |
| `intro` | `TEXT` | Event intro. |
| `description` | `TEXT` | Event description. |
| `system` | `TEXT` | Game system. |
| `event_type` | `TEXT` | Event type. |
| `age_group` | `TEXT` | Age group. |
| `event_runtime` | `TEXT` | Runtime category. |
| `host_name` | `TEXT` | Display host name. |
| `user_id` | `INTEGER` | Owning user id. |
| `email` | `TEXT` | Contact email. |
| `phone_number` | `TEXT` | Contact phone number. |
| `max_players` | `INTEGER` | Max players. |
| `beginner_friendly` | `INTEGER` | Beginner-friendly flag, `0` or `1`. |
| `can_be_run_in_english` | `INTEGER` | English-capable flag, `0` or `1`. |
| `notes` | `TEXT` | Event notes. |
| `status` | `TEXT` | Event status. |
| `created_at` | `TEXT` | Event creation timestamp. |
| `is_published` | `INTEGER` | Pulje publish flag for this event, `0` or `1`. |
| `pulje_id` | `TEXT` | Pulje id. |
| `room_id` | `INTEGER` | Assigned room id, or `NULL`. |
| `room_number` | `TEXT` | Room number, or `NULL`. |
| `room_name` | `TEXT` | Room name, or `NULL`. |
| `room_floor` | `INTEGER` | Room floor, or `NULL`. |
| `room_max_concurrent_games` | `INTEGER` | Room capacity, or `NULL`. |
| `room_notes` | `TEXT` | Room notes, or `NULL`. |
| `room_is_disabled` | `INTEGER` | Room disabled flag, `0`/`1`, or `NULL`. |
| `pulje_name` | `TEXT` | Pulje display name. |
| `pulje_start_at` | `TEXT` | Pulje start timestamp. |
| `pulje_end_at` | `TEXT` | Pulje end timestamp. |

**Go usage example:**

```go
rows, err := db.Query("SELECT id, title, pulje_id, pulje_start_at, room_number, room_name, room_is_disabled FROM v_events_by_pulje_active WHERE is_published = $1", 1)
```

---

# Database View: `v_event_summary`

**Purpose:** Compact projection of commonly used event metadata for admin/summary screens.

**Columns:**

| Column | Type | Description |
| --- | --- | --- |
| `id` | `TEXT` | Event id. |
| `title` | `TEXT` | Event title. |
| `intro` | `TEXT` | Event intro. |
| `status` | `TEXT` | Event status. |
| `system` | `TEXT` | Game system. |
| `host_name` | `TEXT` | Display host name. |
| `beginner_friendly` | `INTEGER` | Beginner-friendly flag, `0` or `1`. |
| `event_type` | `TEXT` | Event type. |
| `age_group` | `TEXT` | Age group. |
| `event_runtime` | `TEXT` | Runtime category. |
| `can_be_run_in_english` | `INTEGER` | English-capable flag, `0` or `1`. |
| `created_at` | `TEXT` | Event creation timestamp. |
| `updated_at` | `TEXT` | Event update timestamp. |

**Go usage example:**

```go
rows, err := db.Query("SELECT id, title, status FROM v_event_summary WHERE status = $1", "Godkjent")
```

---

# Database View: `v_billettholder_emails`

**Purpose:** Billettholdere joined with their email rows for email-based lookup and display.

**Columns:**

| Column | Type | Description |
| --- | --- | --- |
| `billettholder_id` | `INTEGER` | Billettholder id. |
| `first_name` | `TEXT` | Billettholder first name. |
| `last_name` | `TEXT` | Billettholder last name. |
| `ticket_type_id` | `INTEGER` | CheckIn ticket type id. |
| `ticket_type` | `TEXT` | CheckIn ticket type name. |
| `is_over_18` | `INTEGER` | Over-18 flag, `0` or `1`. |
| `order_id` | `INTEGER` | CheckIn order id. |
| `ticket_id` | `INTEGER` | CheckIn ticket id. |
| `billettholder_created_at` | `TEXT` | Billettholder creation timestamp. |
| `billettholder_updated_at` | `TEXT` | Billettholder update timestamp. |
| `email_id` | `INTEGER` | Email row id, or `NULL`. |
| `email` | `TEXT` | Email address, or `NULL`. |
| `kind` | `TEXT` | Email kind, or `NULL`. |
| `email_created_at` | `TEXT` | Email creation timestamp, or `NULL`. |
| `email_updated_at` | `TEXT` | Email update timestamp, or `NULL`. |

**Go usage example:**

```go
rows, err := db.Query("SELECT billettholder_id, first_name, last_name FROM v_billettholder_emails WHERE email = $1", email)
```

---

# Database View: `v_event_puljer_active`

**Purpose:** Active and published pulje rows for events, including optional room metadata.

**Filters:**

- `relation_event_puljer.is_in_pulje = 1`
- `relation_event_puljer.is_published = 1`

**Columns:**

| Column | Type | Description |
| --- | --- | --- |
| `event_id` | `TEXT` | Event id. |
| `pulje_id` | `TEXT` | Pulje id. |
| `room_id` | `INTEGER` | Assigned room id, or `NULL`. |
| `room_number` | `TEXT` | Room number, or `NULL`. |
| `room_name` | `TEXT` | Room name, or `NULL`. |
| `room_floor` | `INTEGER` | Room floor, or `NULL`. |
| `room_max_concurrent_games` | `INTEGER` | Room capacity, or `NULL`. |
| `room_notes` | `TEXT` | Room notes, or `NULL`. |
| `room_is_disabled` | `INTEGER` | Room disabled flag, `0`/`1`, or `NULL`. |
| `pulje_name` | `TEXT` | Pulje display name. |
| `pulje_start_at` | `TEXT` | Pulje start timestamp. |
| `pulje_end_at` | `TEXT` | Pulje end timestamp. |
| `is_in_pulje` | `INTEGER` | Active-in-pulje flag, always `1` in this view. |
| `is_published` | `INTEGER` | Published-in-pulje flag, always `1` in this view. |

**Go usage example:**

```go
rows, err := db.Query("SELECT pulje_id, pulje_start_at, pulje_end_at, room_number, room_name, room_notes FROM v_event_puljer_active WHERE event_id = $1", eventID)
```
