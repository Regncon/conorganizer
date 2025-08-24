CREATE TABLE IF NOT EXISTS backup_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    backup_type TEXT NOT NULL CHECK (backup_type IN ('hourly', 'daily', 'weekly', 'yearly', 'manually')),
    stage TEXT NOT NULL CHECK (stage IN ('starting', 'downloading', 'decompressing', 'validating', 'moving', 'completed')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'error', 'success')),
    file_path TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_backup_logs_type_status_created_at
ON backup_logs(backup_type, status, created_at DESC);

