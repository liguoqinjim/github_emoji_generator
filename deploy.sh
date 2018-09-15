#!/bin/bash

# 调试用set -e
# set -e

echo "发布到github_emoji"

git clone https://${GITHUB_TOKEN}@github.com/liguoqinjim/github_emoji.git ./files2
cd files2
ls

git rm -rf .
cp -R ../files/* .
cp ../deployFiles/README.md .
ls
echo "add"
git add -f --ignore-errors --all
echo "commit"
git -c user.name='liguoqinjim' -c user.email='liguoqinjim23@gmail.com' commit -m "deploy by travis"
git push -f -q https://${GITHUB_TOKEN}@github.com/liguoqinjim/github_emoji.git master
