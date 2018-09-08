package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type GithubEmoji struct {
	Key   string
	Value string
	Codes []string

	Match bool
	Spec  bool //github专属emoji
}

type UnicodeEmoji struct {
	Id         int
	Codes      []string
	ShortName  string
	BigHead    string
	MediumHead string

	Match       bool
	GithubEmoji *GithubEmoji
}

var githubEmojis = make([]*GithubEmoji, 0)

func main() {
	//读取github emojis
	f, err := os.Open("emojis.json")
	if err != nil {
		log.Fatalf("os.Open error:%v", err)
	}
	defer f.Close()

	//reg
	r, err := regexp.Compile(`(\")(.*?)(\")`)
	if err != nil {
		log.Fatalf("regexp.Compile error:%v", err)
	}
	r2, err := regexp.Compile(`(unicode/)(.*?)(.png)`)
	if err != nil {
		log.Fatalf("regexp.Compile error:%v", err)
	}

	//按行读取
	buf := bufio.NewReaderSize(f, 0)
	for {
		line, err := buf.ReadBytes('\n')

		if !strings.Contains(string(line), ":") {
			if err == io.EOF {
				break
			}
			continue
		}

		vs := r.FindAllStringSubmatch(string(line), -1)
		if len(vs) != 2 {
			log.Fatalf("vs length error:%d", len(vs))
		}

		//key
		key := vs[0][2]

		//value
		value := vs[1][2]

		//code
		if strings.Contains(value, "unicode") {
			cs := r2.FindStringSubmatch(value)
			if len(cs) != 4 {
				log.Fatalf("cs length error:%d,%s", len(cs), value)
			}
			codes := strings.Split(cs[2], "-")
			githubEmojis = append(githubEmojis, &GithubEmoji{Key: key, Value: value, Codes: codes})
		} else {
			githubEmojis = append(githubEmojis, &GithubEmoji{Key: key, Value: value, Spec: true})
		}

		if err == io.EOF {
			break
		}
	}
	log.Println("github emoji length=", len(githubEmojis))

	f2, err := os.Open("full-emoji-list.html")
	if err != nil {
		log.Fatalf("os.Open error:%v", err)
	}
	defer f2.Close()

	doc, err := goquery.NewDocumentFromReader(f2)
	if err != nil {
		log.Fatalf("goquery.NewDocumentFromReader error:%v", err)
	}

	unicodeEmojisSlice := make([]*UnicodeEmoji, 0)
	bigHead := ""
	mediumHead := ""

	bigHeads := make([]string, 0)
	allMediumHeads := make([][]string, 0)
	mediumHeads := make([]string, 0)

	unicodeEmojis := make([]*UnicodeEmoji, 0)
	unicodeEmojisMedium := make([][]*UnicodeEmoji, 0)
	unicodeEmojisBig := make([][][]*UnicodeEmoji, 0)

	initBig := false
	initMedium := false
	doc.Find("body > div.main > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		trType := getTrType(s)
		if trType == TR_TYPE_ERROR {
			log.Fatalf("trType error:%d", trType)
		}

		if trType == TR_TYPE_BIG_HEAD {
			h := s.Find("th > a").Text()
			if h == "" {
				log.Fatalf("bigHead error")
			}
			bigHead = h
			bigHeads = append(bigHeads, bigHead)

			if !initBig {
				initBig = true
			} else {
				allMediumHeads = append(allMediumHeads, mediumHeads)
				unicodeEmojisMedium = append(unicodeEmojisMedium, unicodeEmojis)
				unicodeEmojisBig = append(unicodeEmojisBig, unicodeEmojisMedium)

				mediumHeads = make([]string, 0)
				unicodeEmojis = make([]*UnicodeEmoji, 0)
				unicodeEmojisMedium = make([][]*UnicodeEmoji, 0)
			}

			log.Println("BigHead=", bigHead)
		} else if trType == TR_TYPE_MEDIUM_HEAD {
			h := s.Find("th > a").Text()
			if h == "" {
				log.Fatalf("bigHead error")
			}
			mediumHead = h
			mediumHeads = append(mediumHeads, mediumHead)

			if !initMedium {
				initMedium = true
			} else {
				unicodeEmojisMedium = append(unicodeEmojisMedium, unicodeEmojis)
				unicodeEmojis = make([]*UnicodeEmoji, 0)
			}

			log.Println("MediumHead=", mediumHead)
		} else if trType == TR_TYPE_EMOJI {
			id, err := strconv.Atoi(s.Find("td.rchars").Text())
			if err != nil {
				log.Fatalf("emoji id error:%v", err)
			}

			c, exists := s.Find("td.code > a").Attr("name")
			if !exists {
				log.Fatalf("emoji name error:%v", err)
			}
			cs := strings.Split(c, "_")

			sname := s.Find("td.name").Text()
			if sname == "" {
				log.Fatalf("emoji short name error")
			}

			ue := &UnicodeEmoji{Id: id, Codes: cs, ShortName: sname, BigHead: bigHead, MediumHead: mediumHead}

			unicodeEmojisSlice = append(unicodeEmojisSlice, ue)
			unicodeEmojis = append(unicodeEmojis, ue)
		}
	})
	allMediumHeads = append(allMediumHeads, mediumHeads)
	unicodeEmojisMedium = append(unicodeEmojisMedium, unicodeEmojis)
	unicodeEmojisBig = append(unicodeEmojisBig, unicodeEmojisMedium)

	log.Println("unicode emojis length=", len(unicodeEmojisSlice))

	//github和unicode作对比
	for _, v := range unicodeEmojisSlice {
		v.MatchGithubEmoji()
	}

	dir, _ := os.Getwd() //当前的目录

	//创建all_github_emojis
	tmplPaths := []string{"tmpls/all_github_emoji.tmpl", "tmpls/unicode_emoji.tmpl"}
	t, err := template.New("all_github").ParseFiles(tmplPaths...)
	if err != nil {
		log.Fatalf("t1 New error:%v", err)
	}
	if err := os.MkdirAll(dir+"/files/all_github_emojis", os.ModePerm); err != nil {
		log.Fatalf("os.MkdirAll error:%v", err)
	}
	if f_all, err := os.Create("files/all_github_emojis/README.md"); err != nil {
		log.Fatalf("os.Create error:%v", err)
	} else {
		err := t.ExecuteTemplate(f_all, "all_github_emoji.tmpl", githubEmojis)
		if err != nil {
			log.Fatalf("t1.Execute error:%v", err)
		}

		f_all.Close()
	}

	for n, bigHead := range bigHeads {
		err := os.MkdirAll(dir+"/files/unicode/"+bigHead, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			log.Fatalf("os.Mkdir error:%v", err)
		}

		for n2, mediumHead := range allMediumHeads[n] {
			tmplData := &TmplData{FileName: fmt.Sprintf("%s|%s", bigHead, mediumHead), UnicodeEmojis: unicodeEmojisBig[n][n2]}

			f, err := os.Create(fmt.Sprintf("files/unicode/%s/%s.md", bigHead, mediumHead))
			if err != nil {
				log.Fatalf("os.Create error:%v", err)
			}

			err = t.ExecuteTemplate(f, "unicode_emoji.tmpl", tmplData)
			if err != nil {
				log.Fatalf("t.Execute error:%v", err)
			}

			f.Close()
		}
	}
}

//unicode分类用的
type TmplData struct {
	FileName      string
	UnicodeEmojis []*UnicodeEmoji
}

const (
	TR_TYPE_BIG_HEAD = iota + 1
	TR_TYPE_MEDIUM_HEAD
	TR_TYPE_MEANLESS  //无意义行
	TR_TYPE_EMOJI

	TR_TYPE_ERROR
)

//得到这一行tr是什么类型的
func getTrType(s *goquery.Selection) int {
	if s.Find(".bighead").Size() == 1 {
		return TR_TYPE_BIG_HEAD
	}

	if s.Find(".mediumhead").Size() == 1 {
		return TR_TYPE_MEDIUM_HEAD
	}

	if s.Find("th").Size() == 15 {
		return TR_TYPE_MEANLESS
	}

	if s.Find("td.rchars").Size() == 1 {
		_, err := strconv.Atoi(s.Find("td.rchars").Text())
		if err != nil {
			return TR_TYPE_ERROR
		}

		return TR_TYPE_EMOJI
	}

	return TR_TYPE_ERROR
}

func (unicodeEmoji *UnicodeEmoji) MatchGithubEmoji() {
	for _, githubEmoji := range githubEmojis {

		for _, ucode := range unicodeEmoji.Codes {
			has := false
			for _, gcode := range githubEmoji.Codes {
				if strings.ToLower(ucode) == strings.ToLower(gcode) {
					has = true
					break
				}
			}
			if !has {
				break
			}

			if unicodeEmoji.Match {
				log.Fatalf("data error")
			}
			unicodeEmoji.Match = true
			githubEmoji.Match = true
			unicodeEmoji.GithubEmoji = githubEmoji
			return
		}
	}
}

//在当前的目录的files文件夹下创建目录，path格式  dir1/dir2
func NewDir(path string) {

}