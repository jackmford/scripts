#!/bin/bash

# Replace this with your actual Pi IP or hostname if set
PI_USER=jackmf
PI_HOST=backup-pi.local
PI_DEST="/mnt/backupdrive/backups/laptop"
RSYNC_IGNORE="/Users/jackfordyce/Lab/git-repos/scripts/pi-backup/rsync-ignore.txt"

echo "Starting backup..."

echo "Backing up... ~/Lab"
rsync -avh --delete --exclude-from=$RSYNC_IGNORE ~/Lab $PI_USER@$PI_HOST:$PI_DEST/

echo "Backing up... Obsidian Vault"
rsync -avh --delete --exclude-from=$RSYNC_IGNORE ~/Library/Mobile\ Documents/iCloud~md~obsidian/Documents/Obsidian\ Vault $PI_USER@$PI_HOST:$PI_DEST/

echo "Backup complete."

