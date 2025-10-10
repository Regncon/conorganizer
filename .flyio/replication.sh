#!/bin/bash
set -e

# Validate config file
if [ -f /etc/litestream.yml ]; then
	echo "Configuration file found"
else
	echo "No configuration found, using cli arguments "
fi

# Restore the database if it does not already exist.
if [ -f /data/regncon/database/events.db ]; then
	echo "Database already exists, skipping restore"
else
	echo "No database found, restoring from replica if exists"
	litestream restore /data/regncon/database/events.db
fi

# Check if the image folder exists and we have write permissions to it
if [ -d /data/regncon/uploads ] && [ -w /data/regncon/uploads ]; then
    echo "Image upload folder exists and is writable"
else
    echo "Image upload folder does not exist or is not writable"
    exit 1
fi

# Run litestream with your app as the subprocess.
exec litestream replicate -exec "/usr/local/bin/regncon -dbp /data/regncon/database/events.db -image-path /data/regncon/uploads"
