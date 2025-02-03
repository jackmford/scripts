#!/bin/bash

# Configuration
NOTES_DIR="/Users/jackfordyce/Library/Mobile Documents/iCloud~md~obsidian/Documents/Obsidian Vault/para/1 Projects/2025/dailys/"

# Cleanup script
cleanup_empty_notes() {
    local days_to_keep=7  # How many days of history to preserve regardless of content

    find "$NOTES_DIR" -name "*.md" -type f | while read -r note; do
        # Skip recent files based on days_to_keep
        if [[ $(find "$note" -mtime -"$days_to_keep") ]]; then
            continue
        fi

        # Check if file only contains template content or is nearly empty
        local meaningful_content=$(grep -Ev '^(#daily|Tasks:|Scratchpad:|💻 Log:|Journal:|Vim Tip:|$|- \[ \] (Make bed|Vitamins \+ Moisturize|Leetcode))' "$note" | wc -l)

        if  [[ $meaningful_content -gt 0 ]]; then
            echo "Removing empty note: $(basename "$note")"
            rm "$note"
        fi
        
    done
}

cleanup_empty_notes
