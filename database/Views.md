# Database Views

This file documents database views. Each entry contains a short purpose and a Go usage example.

---

# Database View: `v_get_user_billettholder`

**Purpose:** Join between `users` and `billettholdere_users` to map application users to ticket-holder records.

**Go usage example:**

```go
rows, err := db.Query("SELECT billettholder_id FROM v_get_user_billettholder WHERE user_id = $1", userID)
```

---

# Database View: `v_events_by_pulje_active`

**Purpose:** Pre-joined view of approved events with their active pulje metadata for listings.

**Go usage example:**

```go
rows, err := db.Query("SELECT id, title, pulje_name FROM v_events_by_pulje_active WHERE is_published = $1", 1)
```

---

# Database View: `v_event_summary`

**Purpose:** Compact projection of commonly used event metadata for admin/summary screens.

**Go usage example:**

```go
rows, err := db.Query("SELECT id, title, status FROM v_event_summary WHERE status = $1", "Godkjent")
```

---

# Database View: `v_billettholder_emails`

**Purpose:** Joins `billettholdere` with their emails for email-based lookups.

**Go usage example:**

```go
rows, err := db.Query("SELECT billettholder_id, first_name, last_name FROM v_billettholder_emails WHERE email = $1", email)
```

---

# Database View: `v_event_puljer_active`

**Purpose:** Active and published puljer for an event.

**Go usage example:**

```go
rows, err := db.Query("SELECT pulje_id, pulje_name FROM v_event_puljer_active WHERE event_id = $1", eventID)
```

