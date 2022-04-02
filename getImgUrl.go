package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ImageIndex struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Host   string `json:"host"`
		Images []struct {
			Path string `json:"path"`
		} `json:"images"`
	} `json:"data"`
}

type ImageToken struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	} `json:"data"`
}

func getImgUrl(id int) (urls []string, err error) {
	// log.Println(id)
	imageIndex, err := getImageIndex(id)
	if err != nil {
		return nil, err
	}
	// log.Println(imageIndex)

	urls = make([]string, len(imageIndex.Data.Images))
	for i, url := range imageIndex.Data.Images {
		urls[i] = url.Path
	}
	imagesToken, err := getImageToken(urls)
	if err != nil {
		return nil, err
	}
	// log.Println(imageToken)
	for i, imageToken := range imagesToken.Data {
		urls[i] = imageToken.Url + "?token=" + imageToken.Token
	}

	return urls, nil
}

func getImageIndex(id int) (imageIndex ImageIndex, err error) {
	req, err := http.NewRequest(
		"POST",
		"https://manga.bilibili.com/twirp/comic.v1.Comic/GetImageIndex?device=pc&platform=web",
		strings.NewReader(
			fmt.Sprintf(`{"ep_id":%d}`, id),
		),
	)
	if err != nil {
		return imageIndex, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return imageIndex, err
	}
	imageIndexJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return imageIndex, err
	}
	// ImageIndexFile, err := os.OpenFile(
	// 	"ImageIndex.json",
	// 	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	// 	0666,
	// )
	// if err == nil {
	// 	ImageIndexFile.Write(images)
	// 	ImageIndexFile.Close()
	// }
	err = json.Unmarshal(imageIndexJson, &imageIndex)
	if err != nil {
		return imageIndex, err
	}
	if imageIndex.Code != 0 {
		return imageIndex, errors.New(imageIndex.Msg)
	}
	return imageIndex, nil
}

func getImageToken(urls []string) (imagesToken ImageToken, err error) {
	urlsJsonBytes, err := json.Marshal(urls)
	if err != nil {
		return imagesToken, err
	}
	urlsJson := strings.ReplaceAll(string(urlsJsonBytes), `"`, `\"`)
	req, err := http.NewRequest(
		"POST",
		"https://manga.bilibili.com/twirp/comic.v1.Comic/ImageToken?device=pc&platform=web",
		strings.NewReader(
			fmt.Sprintf(`{"urls":"%s"}`, string(urlsJson)),
		),
	)
	if err != nil {
		return imagesToken, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return imagesToken, err
	}
	imageTokenJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return imagesToken, err
	}
	// log.Println(urlsJson)
	// log.Println(string(imageTokenJson))
	err = json.Unmarshal(imageTokenJson, &imagesToken)
	if err != nil {
		return imagesToken, err
	}
	if imagesToken.Code != 0 {
		return imagesToken, errors.New(imagesToken.Msg)
	}
	return imagesToken, nil
}
