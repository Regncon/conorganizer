#!/bin/bash
set -e

# Validate config file
if [ -f /etc/litestream.yml ]; then
	echo "Configuration file found"
else
	echo "No configuration found, using cli arguments "
fi

# Restore the database if it does not already exist.
if [ -f /var/lib/regnconDB ]; then
	echo "Database already exists, skipping restore"
else
	echo "No database found, restoring from replica if exists"
	litestream restore -if-replica-exists /var/lib/regncon/events.db
fi

# Run litestream with your app as the subprocess.
exec litestream replicate -exec "/usr/local/bin/regncon -dbp /var/lib/regncon/events.db"
