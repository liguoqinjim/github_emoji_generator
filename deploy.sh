#!/bin/bash

# 调试用set -e
# set -e

echo "发布到github_emoji"

echo "github token"
echo $GITHUB_TOKEN

#touch ~/.git-credentials
#echo "https://136542728%40qq.com:wangxiu@git.coding.net" > ~/.git-credentials
#cat ~/.git-credentials
#git version
#git config --global credential.helper store

#git https://github.com/liguoqinjim/github_emoji.git file2
#cd file2
#git rm -rf .
#cp -R ../file/* .
#echo "add"
#git add -f --ignore-errors --all
#echo "commit"
#git -c user.name='travis' -c user.email='travis' commit -m init
#git "push"
#git push -f -q https://git.coding.net/liguoqinjim/liguoqinjim.coding.me.git master
#
#echo "生成liguoqinjim.com"
#cd ..
#cp config_com.toml config.toml
#hugo