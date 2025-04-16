#!/bin/bash

# Replace this with your actual Pi IP or hostname if set
PI_USER=jackmf
PI_HOST=backup-pi.local
PI_DEST="/mnt/backupdrive/backups/laptop"

echo "Starting backup..."

rsync -avh --delete --exclude-from=rsync-ignore.txt ~/Lab $PI_USER@$PI_HOST:$PI_DEST/
rsync -avh --delete --exclude-from=rsync-ignore.txt ~/Lab $PI_USER@$PI_HOST:$PI_DEST/
rsync -avh --delete --exclude-from=rsync-ignore.txt ~/Library/Mobile\ Documents/iCloud~md~obsidian/Documents/Obsidian\ Vault $PI_USER@$PI_HOST:$PI_DEST/

echo "Backup complete."

