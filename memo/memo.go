package memo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Shakkuuu/login-app-api-use/entity"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Memoservice struct{}
type Memo entity.Memo

// メモ機能

func (ms Memoservice) Memo(c *gin.Context) {
	// ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// そのユーザーのメモを取得
	m := memoGet(uname)

	c.HTML(200, "memo.html", gin.H{"memos": m})
}

// メモ取得
func memoGet(uname string) []Memo {
	url := "http://localhost:8081/memos/showname/" + uname
	// apiでユーザー名からメモの一覧を取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 取得したjsonのオープン
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// openしたjsonを構造体にデコード
	var m []Memo
	if err := json.Unmarshal(body, &m); err != nil {
		log.Fatal(err)
	}

	return m
}

// メモ登録
func (ms Memoservice) MemoCreate(c *gin.Context) {
	// ログイン中ユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	url := "http://localhost:8081/memos"
	username := uname
	// 登録するメモのタイトルと本文を取得
	title := c.PostForm("title")
	text := c.PostForm("text")

	// 未入力か確認
	if username == "" || title == "" || text == "" {
		c.Redirect(303, "/menu/memo")
		return
	}

	// jsonに変換
	jsonStr := `{"Name":"` + username + `","Title":"` + title + `","Text":"` + text + `"}`

	// apiにjsonを送信して登録
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	c.Redirect(303, "/menu/memo")
}
