package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ComicDetail struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Ep_list []Chapter `json:"ep_list"`
	} `json:"data"`
}

func getChaptersInfo(comicId int) (chaptersInfo []Chapter, err error) {
	req, err := http.NewRequest(
		"POST",
		"https://manga.bilibili.com/twirp/comic.v1.Comic/ComicDetail?device=pc&platform=web",
		strings.NewReader(
			fmt.Sprintf(`{"comic_id":%d}`, comicId),
		),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	comicDetailJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var comicDetail ComicDetail
	err = json.Unmarshal(comicDetailJson, &comicDetail)
	if err != nil {
		return nil, err
	}
	if comicDetail.Code != 0 {
		return nil, errors.New(comicDetail.Msg)
	}

	return comicDetail.Data.Ep_list, nil
}
