package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

var suffixPattern *regexp.Regexp

var chaptersSelector *regexp.Regexp

func init() {
	var err error
	suffixPattern, err = regexp.Compile(
		`(\.[(a-z)|(A-Z)|(0-9)]*)\?`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	chaptersSelector, err = regexp.Compile(
		`^257$`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func downloadCommic(chaptersInfo []Chapter) error {
	err := login()
	if err != nil {
		return err
	}

	fail := make([]string, 0)
	for _, chapterInfo := range chaptersInfo {
		matches := chaptersSelector.FindStringSubmatch(chapterInfo.ShortTitle)
		if len(matches) == 0 {
			continue
		}
		var ImagesUrl []string
		err = nil
		for i := 0; i < 5; i++ {
			ImagesUrl, err = getImgUrl(chapterInfo.Id)
			if err == nil {
				break
			}
		}
		title := chapterInfo.ShortTitle + " " + chapterInfo.Title
		if err != nil {
			fail = append(fail, title)
			log.Println("下载" + title + "失败：" + err.Error())
			continue
		}
		_, err := os.Stat(title)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(title, 0777)
				if err != nil {
					fail = append(fail, title)
					log.Println("下载" + title + "失败：" + err.Error())
					continue
				}
			} else {
				fail = append(fail, title)
				log.Println("下载" + title + "失败：" + err.Error())
				continue
			}
		}

		failed := false
		for index, ImageUrl := range ImagesUrl {
			var i int
			for i = 0; i < 5; i++ {
				err := saveImg(title, ImageUrl, index)
				if err == nil {
					break
				}
			}
			if i >= 5 {
				failed = true
				break
			}
			// time.Sleep(1 * time.Second)
		}
		if failed {
			fail = append(fail, title)
			log.Println("下载" + title + "失败：" + err.Error())
		}
	}

	if len(fail) == 0 {
		log.Println("全部下载完成")
	} else {
		log.Println("以下章节下载失败：")
		for _, f := range fail {
			log.Println(f)
		}
	}
	return nil
}

func saveImg(title string, url string, index int) error {
	suffix := ""
	matches := suffixPattern.FindAllStringSubmatch(url, -1)
	if len(matches) > 0 {
		suffix = matches[len(matches)-1][1]
	}
	fileName := fmt.Sprintf("%s/%02d%s", title, index, suffix)

	log.Println("正在下载：" + fileName)
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		log.Println("下载" + fileName + "失败：" + err.Error())
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("下载" + fileName + "失败：" + err.Error())
		return err
	}

	imgFile, err := os.OpenFile(
		fileName,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Println("下载" + fileName + "失败：" + err.Error())
		return err
	}
	defer imgFile.Close()
	_, err = io.Copy(imgFile, resp.Body)
	if err != nil {
		log.Println("下载" + fileName + "失败：" + err.Error())
		return err
	}
	return nil
}
