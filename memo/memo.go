package memo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ログイン後のトップページ
func top(c *gin.Context) {
	// ログインしているユーザー名を取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		uname = "ゲスト"
		msg := "現在ゲストで使用しています。ログインしましょう。"
		c.HTML(200, "top.html", gin.H{"user": uname, "message": msg})
		return
	}

	c.HTML(200, "top.html", gin.H{"user": uname})
}

// メモ機能

func memo(c *gin.Context) {
	// ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// そのユーザーのメモを取得
	m := MemoGet(uname)

	c.HTML(200, "memo.html", gin.H{"memos": m})
}

// メモ取得
func MemoGet(uname string) []Memo {
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
func memocreate(c *gin.Context) {
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
