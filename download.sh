#!/bin/bash

# github_api
curl -o emojis.json https://api.github.com/emojis
cat emojis.json

# unicode_emoji
curl -o full-emoji-list.html https://unicode.org/emoji/charts/full-emoji-list.html