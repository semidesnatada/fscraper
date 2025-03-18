#!/bin/bash
set -e

# Set variables for database connection
PGUSER=seanlowery
PGDATABASE=fscraped

# Set the path where you want to store the backup files
BACKUP_DIR=db_backups

# Get current date and time
datestamp=$(date +'%Y-%m-%d')
timestamp=$(date +'%H%M')

# Execute pg_dump command to dump the database
# pg_dump -U "$PGUSER" -d "$PGDATABASE" > "$BACKUP_DIR/$PGDATABASE"_"$datestamp"_"$timestamp".sql
pg_dump -U "$PGUSER" -d "$PGDATABASE" -F tar -f "$BACKUP_DIR/$PGDATABASE"_"$datestamp"_"$timestamp".tar

echo "Backup successfully saved to : "$BACKUP_DIR/$PGDATABASE"_"$datestamp"_"$timestamp".tar"
echo "To restore, add your desired restore filename to: ./restore_db.sh, and run it from command line."