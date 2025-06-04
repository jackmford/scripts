#!/bin/bash

PI_USER=jackmf
PI_HOST=backup-pi.local
PI_DEST="/mnt/backupdrive/backups/laptop"
RSYNC_IGNORE="/Users/jackfordyce/Lab/git-repos/scripts/pi-backup/rsync-ignore.txt"

echo "Starting backup..."

echo "Backing up... ~/Lab"
rsync -avh --delete --partial --append-verify --exclude-from=$RSYNC_IGNORE ~/Lab $PI_USER@$PI_HOST:$PI_DEST/

echo "Backing up... Obsidian Vault"
rsync -avh --delete --partial --append-verify --exclude-from=$RSYNC_IGNORE ~/Library/Mobile\ Documents/iCloud~md~obsidian/Documents/Obsidian\ Vault $PI_USER@$PI_HOST:$PI_DEST/

echo "Backing up... iCloud Drive"
rsync -avh --delete --partial --append-verify --exclude-from=$RSYNC_IGNORE ~/Library/Mobile\ Documents/com~apple~CloudDocs $PI_USER@$PI_HOST:$PI_DEST/

echo "Backup complete."

