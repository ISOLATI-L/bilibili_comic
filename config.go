package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
)

type config struct {
	id             int
	chapterPattern string
}

const (
	DEFAULT_COMIC_ID       int    = 26009
	DEFAULT_CHAPTERPATTERN string = ".*?"
)

func loadConfig() (cfg config) {
	var configFile *goconfig.ConfigFile
	var err error
	configFile, err = goconfig.LoadConfigFile("config.ini")
	for err != nil {
		if os.IsNotExist(err) {
			os.Create("config.ini")
			configFile, err = goconfig.LoadConfigFile("config.ini")
		} else {
			log.Fatalln(err.Error())
		}
	}
	defer func() {
		err = goconfig.SaveConfigFile(configFile, "config.ini")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	idStr, err1 := configFile.GetValue(goconfig.DEFAULT_SECTION, "id")
	id, err2 := strconv.Atoi(idStr)
	if err1 != nil || err2 != nil {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"id",
			fmt.Sprint(DEFAULT_COMIC_ID),
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"id",
			fmt.Sprintf(
				"id为漫画id，如《名侦探柯南》的链接为https://ac.qq.com/Comic/comicInfo/id/%d，id即为%d",
				DEFAULT_COMIC_ID, DEFAULT_COMIC_ID,
			),
		)
		id = DEFAULT_COMIC_ID
	}
	cfg.id = id

	chapterPattern, err := configFile.GetValue(goconfig.DEFAULT_SECTION, "chapterPattern")
	if err != nil || chapterPattern == "" {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			DEFAULT_CHAPTERPATTERN,
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			fmt.Sprintf(
				"chapterPattern为检索章节的正则表达式，默认为%s（即检索所有章节）",
				DEFAULT_CHAPTERPATTERN,
			),
		)
		chapterPattern = DEFAULT_CHAPTERPATTERN
	}
	cfg.chapterPattern = chapterPattern

	return cfg
}
