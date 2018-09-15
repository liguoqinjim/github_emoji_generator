#!/bin/bash

# test，测试API的使用量
curl -i https://api.github.com/users/whatever?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}

# github_api
# https://${GITHUB_TOKEN}@github.com/liguoqinjim/github_emoji.git
curl -o emojis.json https://api.github.com/users/whatever?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}

cat emojis.json

# unicode_emoji
curl -o full-emoji-list.html https://unicode.org/emoji/charts/full-emoji-list.html