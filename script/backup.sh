#!/bin/sh

PGHOST=${PGHOST:-postgres}
PGUSER=${PGUSER:-user}
PGDATABASE=${PGDATABASE:-mydb}
PGPASSWORD=${PGPASSWORD:-password}
BACKUP_DIR=/backups

mkdir -p "$BACKUP_DIR"

BACKUP_FILE="$BACKUP_DIR/backup_$(date +%Y-%m-%d_%H-%M-%S).sql"

echo "[$(date)] Starting backup..."

PGPASSWORD=$PGPASSWORD pg_dump -h "$PGHOST" -U "$PGUSER" -d "$PGDATABASE" -F p -f "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "[$(date)] Backup completed successfully: $BACKUP_FILE"
    
    # Keep only backups from the last 7 days
    echo "[$(date)] Cleaning old backups..."
    find "$BACKUP_DIR" -type f -name '*.sql' -mtime +7 -exec rm {} \;
    echo "[$(date)] Old backups cleaned."
else
    echo "[$(date)] ERROR: Backup failed!"
    exit 1
fi