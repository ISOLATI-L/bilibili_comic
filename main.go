package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Chapter struct {
	Id         int    `json:"id"`
	ShortTitle string `json:"short_title"`
	Title      string `json:"title"`
}

var client *http.Client
var CFG config

func init() {
	var err error
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	client = &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不进入重定向
		},
		Jar: jar,
	}
	CFG = loadConfig()
	initChaptersSelector()
}

func main() {
	chaptersInfo, err := getChaptersInfo(CFG.id)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = downloadCommic(chaptersInfo)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// length := len(chaptersInfo)
	// ImagesUrl, err := getImgUrl(chaptersInfo[length-1].Id)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// for _, ImageUrl := range ImagesUrl {
	// 	log.Println(ImageUrl)
	// }

	// err = login()
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	// ImagesUrl, err = getImgUrl(656591)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// for _, ImageUrl := range ImagesUrl {
	// 	log.Println(ImageUrl)
	// }
}
