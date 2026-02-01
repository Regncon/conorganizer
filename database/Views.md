# Database View: `get_billettholder_id`

This view is defined in the migration file `20260201161247_views.sql`. It combines data from the `users` table and the `billettholdere_users` table, making it easier to query information about users who are also billettholdere (ticket holders).

## Definition

The view is created with the following SQL:

```sql
CREATE VIEW IF NOT EXISTS
    get_billettholder_id AS
SELECT
    u.*,
    bu.*
FROM
    billettholdere_users AS bu
    LEFT JOIN users AS u ON u.id = bu.user_id;
```

- This view selects all columns from both `users` (`u.*`) and `billettholdere_users` (`bu.*`).
- It performs a `LEFT JOIN` from `billettholdere_users` to `users` using the `user_id` field.

## Purpose

- To simplify queries that need both user information and billettholder (ticket holder) details.
- To avoid writing repetitive join statements in your application code.
- To make it easy to fetch the billettholder id for a given user.

## How to Use

You can use this view in your SQL queries just like a regular table.

### Example Query

```sql
SELECT id
FROM get_billettholder_id
WHERE user_id = '42DGW23DAcd257ad';
```

This will return the `id` from the `billettholdere_users` table for the user with `user_id` '42DGW23DAcd257ad'.
Use this when you need the billettholder id for a specific user.

### In Application Code

You can query the view as you would any table. For example, in Go:

````go
rows, userQueryErr := db.Query("SELECT id FROM get_billettholder_id WHERE user_id = $1", userID)
````
