package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

func login() error {
	var loginStatus LoginStatus
	oauthKey, err := ShowQRcode()
	if err != nil {
		return err
	}
	textbackground(0x04)
	fmt.Println("请扫码登陆")
	resettextbackground()

	for {
		loginStatus, err = getLoginInfo(oauthKey)
		if err != nil {
			oauthKey, err = ShowQRcode()
			if err != nil {
				return err
			}
			textbackground(0x04)
			fmt.Println("请扫码登陆")
			resettextbackground()
			continue
		}
		if loginStatus.Status {
			break
		} else if loginStatus.Data == -4.0 || loginStatus.Data == -5.0 {
			time.Sleep(1 * time.Second)
		} else {
			textbackground(0x04)
			fmt.Println("二维码已失效")
			resettextbackground()
			oauthKey, err = ShowQRcode()
			if err != nil {
				return err
			}
			textbackground(0x04)
			fmt.Println("请扫码登陆")
			resettextbackground()
		}
	}

	return nil
}

type LoginUrl struct {
	Data struct {
		Url      string `json:"url"`
		OauthKey string `json:"oauthKey"`
	} `json:"data"`
}

func ShowQRcode() (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://passport.bilibili.com/qrcode/getLoginUrl",
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var loginUrl LoginUrl
	{
		loginUrlJson, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(loginUrlJson, &loginUrl)
		if err != nil {
			return "", err
		}
	}
	// log.Println(loginUrl)

	qrc, err := qrcode.New(loginUrl.Data.Url, qrcode.Low)
	if err != nil {
		return "", err
	}
	binary := qrc.Bitmap()
	for _, row := range binary[3 : len(binary)-3] {
		for _, bit := range row[3 : len(binary)-3] {
			if bit {
				textbackground(0x00)
				fmt.Print("　")
			} else {
				textbackground(0xFF)
				fmt.Print("　")
				// fmt.Print("■")
			}
		}
		resettextbackground()
		fmt.Println()
	}
	resettextbackground()

	return loginUrl.Data.OauthKey, nil
}

type LoginStatus struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func getLoginInfo(oauthKey string) (loginStatus LoginStatus, err error) {
	loginStatus.Status = false
	req, err := http.NewRequest(
		"POST",
		"https://passport.bilibili.com/qrcode/getLoginInfo",
		strings.NewReader(
			fmt.Sprintf(`oauthKey=%s`, oauthKey),
		),
	)
	if err != nil {
		return loginStatus, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return loginStatus, err
	}
	{
		loginStatusJson, err := io.ReadAll(resp.Body)
		if err != nil {
			return loginStatus, err
		}
		err = json.Unmarshal(loginStatusJson, &loginStatus)
		if err != nil {
			return loginStatus, err
		}
	}
	return loginStatus, nil
}
