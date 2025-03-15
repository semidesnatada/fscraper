#!/bin/bash
set -e

# Set variables for database connection
PGUSER=seanlowery
PGDATABASE=fscraped

# Add own backup link to the below variable
DB_BACKUP_LINK=db_backups/fscraped_2025-03-15_2001.tar

pg_restore -U "$PGUSER" -d "$PGDATABASE" "$DB_BACKUP_LINK" --clean

# WARNING - BEFORE RUNNING THIS SCRIPT, 
# DROP AND RECREATE THE DATABASE VIA PSQL, 
# THEN RUN GOOSE UP MIGRATIONS. THEN RESTORE.