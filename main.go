package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const COMIC_ID int = 26009

type Chapter struct {
	Id         int    `json:"id"`
	ShortTitle string `json:"short_title"`
	Title      string `json:"title"`
}

var client *http.Client

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
}

func main() {
	chaptersInfo, err := getChaptersInfo(COMIC_ID)
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
