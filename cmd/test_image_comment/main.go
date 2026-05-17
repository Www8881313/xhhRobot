package main

import (
	"fmt"
	"os"

	"xhhrobot/config"
	"xhhrobot/db"
	"xhhrobot/loger"
	"xhhrobot/xhh"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: test_image_comment LINK_ID TEXT IMAGE_URL")
		os.Exit(1)
	}

	linkID := os.Args[1]
	text := os.Args[2]
	imageURL := os.Args[3]

	loger.InitLog()
	config.InitConfig()
	db.Init()
	xhh.Init()

	if xhh.CommentPostImage(text, linkID, imageURL) {
		fmt.Println("comment image ok")
		return
	}

	fmt.Println("comment image failed")
	os.Exit(1)
}
