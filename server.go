package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var dingdingWebhook string

var TucaoType = map[string]string{
	"post.created":  "新增吐槽",
	"post.updated":  "修改吐槽",
	"reply.created": "回复吐槽",
	"reply.updated": "修改回复",
}

type TucaoPayload struct {
	Post struct {
		Content string `form:"content" json:"content"`
	}
	User struct {
		Username string `form:"username" json:"username"`
		OpenID   string `form:"openid" json:"openid"`
		IsAdmin  bool   `form:"isadmin" json:"isadmin"`
	}
}

type TucaoMessage struct {
	ID       string `form:"id" json:"id"`
	Type     string `form:"type" json:"type"`
	Payload  TucaoPayload
	CreateAt string `form:"crate_at" json:"create_at"`
}

type DingdingBody struct {
	MsgType  string           `json:"msgtype"`
	Markdown DingdingMarkdown `json:"markdown"`
}

type DingdingMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func tucaoHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("1111111")
	decoder := json.NewDecoder(req.Body)
	var msg TucaoMessage
	err := decoder.Decode(&msg)
	fmt.Println(msg, err)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		adminMsg := ""
		if msg.Payload.User.IsAdmin {
			adminMsg = "管理员"
		}
		resData := DingdingBody{
			MsgType: "markdown",
			Markdown: DingdingMarkdown{
				Title: TucaoType[msg.Type],
				Text:  fmt.Sprintf("用户：%s ID：%s %s", msg.Payload.User.Username, msg.Payload.User.OpenID, adminMsg),
			},
		}
		b, err := json.Marshal(resData)
		if err != nil {
			log.Fatal("json format error:", err)
		}
		body := bytes.NewBuffer(b)
		resp, err := http.Post(dingdingWebhook, "application/json;charset=utf-8", body)
		if err != nil {
			log.Fatal("Dingding webhook post error %s", err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read failed:", err)
			return
		}
		log.Println("content:", string(content))
		rw.WriteHeader(http.StatusOK)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dingdingWebhook = os.Getenv("DINGDING_WEBHOOK_URL")
	if dingdingWebhook == "" {
		log.Fatal("No dingding webhook url set")
	}
	http.HandleFunc("/tucao/webhook", tucaoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
