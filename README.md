# github_emoji_generator[![Build Status](https://travis-ci.org/liguoqinjim/github_emoji_generator.svg?branch=master)](https://travis-ci.org/liguoqinjim/github_emoji_generator)
生成 (https://github.com/liguoqinjim/github_emoji) 内容的工具

创建三种
1.按照github的api生成
2.按照unicode的分类，但是和github对照
3.所有github特殊的

有几个问题
1. github有的emoji的unicode和标准的unicode不一样，如github的asterisk，标准unicode有3位
2. github有unicode相同的emoji，但是emoji的name不一样，如collision和boom
3. unicode里面`U+1F487 U+200D U+2642 U+FE0F`，但是github里面没有`U+200D`,`U+FE0F`
4. github里面的man的emoji都加了2642，但是unicode里面是没有的
5. github里面有一些emoji是一样的unicode，但是名字不一样。

## api
1. https://api.github.com/emojis
2. https://unicode.org/emoji/charts/full-emoji-list.html


```
{{ range . }}{{if .Match}}{{else}}{{if .Spec}}{{else}}|{{.Key}}|:{{.Key}}:|{{.Value}}|
{{end}}{{end}}{{ end }}
```

```
{{ range . }}{{if .Match}}{{else}}|{{.Key}}|:{{.Key}}:|{{.Value}}|
{{end}}{{ end }}
```
