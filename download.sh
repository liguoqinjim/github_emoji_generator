#!/bin/bash

#echo "测试api rate"
#curl -i "https://api.github.com/emojis?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}"

# https://${GITHUB_TOKEN}@github.com/liguoqinjim/github_emoji.git
echo "下载github_emoji"
curl -o emojis.json "https://api.github.com/emojis?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}"

#cat emojis.json

# unicode_emoji
echo "下载unicode_emoji"
curl -o full-emoji-list.html https://unicode.org/emoji/charts/full-emoji-list.html