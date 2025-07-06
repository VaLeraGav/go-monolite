#!/bin/sh

PGHOST=${PGHOST:-postgres}
PGUSER=${PGUSER:-user}
PGDATABASE=${PGDATABASE:-mydb}
PGPASSWORD=${PGPASSWORD:-password}
BACKUP_DIR=${BACKUP_DIR:-/backups}

LATEST_BACKUP=$(ls -t "$BACKUP_DIR"/*.sql 2>/dev/null | head -1)

if [ -z "$LATEST_BACKUP" ]; then
  echo "No backup files found in $BACKUP_DIR"
  exit 1
fi

echo "Restoring from $LATEST_BACKUP"

PGPASSWORD=$PGPASSWORD psql -h "$PGHOST" -U "$PGUSER" -d "$PGDATABASE" -f "$LATEST_BACKUP"

if [ $? -eq 0 ]; then
  echo "Restore completed successfully."
else
  echo "Restore failed!"
  exit 1
fi
