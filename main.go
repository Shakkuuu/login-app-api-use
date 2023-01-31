package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/Shakkuuu/login-app-api-use/entity"
	"github.com/Shakkuuu/login-app-api-use/gachagame"
	"github.com/Shakkuuu/login-app-api-use/memo"
	"github.com/Shakkuuu/login-app-api-use/minigame"
	"github.com/Shakkuuu/login-app-api-use/ticketandcoin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type User entity.User
type Memo entity.Memo
type Coin entity.Coin

// UserのJSONをUser構造体にデコード
func (u *User) UnmarshalJSON(body []byte) error {
	// 自分で新しく定義した構造体
	u2 := &struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Ticket   int    `json:"ticket"`
		Coin     int    `json:"coin"`
	}{}
	err := json.Unmarshal(body, u2)
	if err != nil {
		panic(err)
	}
	// 新しく定義した構造体の結果をもとのmに詰める
	u.ID = u2.ID
	u.Name = u2.Name
	u.Password = u2.Password
	u.Ticket = u2.Ticket
	u.Coin = u2.Coin

	return err
}

// MemoのJSONをMemo構造体にデコード
func (m *Memo) UnmarshalJSON(body []byte) error {
	// 自分で新しく定義した構造体
	m2 := &struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}
	err := json.Unmarshal(body, m2)
	if err != nil {
		panic(err)
	}
	// 新しく定義した構造体の結果をもとのmに詰める
	m.ID = m2.ID
	m.Name = m2.Name
	m.Title = m2.Title
	m.Text = m2.Text

	return err
}

// MemoのJSONをMemo構造体にデコード
func (co *Coin) UnmarshalJSON(body []byte) error {
	// 自分で新しく定義した構造体
	co2 := &struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Qty       float32 `json:"qty"`
		Speed     float32 `json:"speed"`
		Speedneed float32 `json:"speedneed"`
	}{}
	err := json.Unmarshal(body, co2)
	if err != nil {
		panic(err)
	}
	// 新しく定義した構造体の結果をもとのmに詰める
	co.ID = co2.ID
	co.Name = co2.Name
	co.Qty = co2.Qty
	co.Speed = co2.Speed
	co.Speedneed = co2.Speedneed

	return err
}

func main() {
	fmt.Println("start")

	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// ログインページ
	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/signup", signup)
	r.GET("/logout", logout)

	r.POST("/login", loginuser)
	r.POST("/signup", signupuser)

	// ログイン後のトップページ
	ms := memo.Memoservice{}
	gg := gachagame.Gachaservice{}
	if err := gg.GachaCreate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tc := ticketandcoin.TandC{}
	mg := minigame.MG{}

	menu := r.Group("/menu")
	menu.GET("/top", top)
	menu.GET("/memo", ms.Memo)

	menu.POST("/memo", ms.MemoCreate)

	game := menu.Group("/game")
	game.GET("/gachagame", gg.GachaGame)
	game.GET("/tandc", tc.TicketandCoin)

	game.POST("/draw", gg.DrawGacha)
	game.POST("/tadd", tc.TicketAdd)
	// game.POST("/cadd", tc.CoinAdd)

	game.GET("/minigame", mg.Minigamemain)

	game.POST("/addcoin", mg.Addcoin)
	game.POST("/addspeedup", mg.AddSpeedUp)

	// user設定ページ
	settings := menu.Group("/settings")
	settings.GET("/deleteuser", deleteusercheck)
	settings.GET("/renameuser", renameusercheck)

	settings.POST("/deleteuser", deleteuser)
	settings.POST("/renameuser", renameuser)

	// 8082ポートで起動
	r.Run(":8082")
}

// ログイン後のトップページ
func top(c *gin.Context) {
	// ログインしているユーザー名を取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		uname = "ゲスト"
		msg := "現在ゲストで使用しています。ログインまたはアカウント登録しましょう。"
		c.HTML(200, "top.html", gin.H{"user": uname, "message": msg})
		return
	}

	gg := gachagame.Gachaservice{}
	gg.TCSet(uname)

	c.HTML(200, "top.html", gin.H{"user": uname})
}

// ログイン前のindexページ
func index(c *gin.Context) {
	// 登録されているユーザーの一覧取得
	u := UserGet()

	c.HTML(200, "index.html", gin.H{"users": u})
}

// ログイン サインアップ ログアウトのページ表示

func login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func signup(c *gin.Context) {
	c.HTML(200, "signup.html", nil)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	// ログアウトによるセッション削除
	session.Clear()
	session.Save()

	c.Redirect(303, "/login")
}

// ユーザー削除とリネームのページ表示

func deleteusercheck(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		// uname = "ゲスト"
		// msg := "現在ゲストで使用しています。ログインまたはアカウント登録しましょう。"
		c.HTML(200, "error.html", nil)
		return
	}
	c.HTML(200, "deleteuser.html", gin.H{"username": uname})
}

func renameusercheck(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		// uname = "ゲスト"
		// msg := "現在ゲストで使用しています。ログインまたはアカウント登録しましょう。"
		c.HTML(200, "error.html", nil)
		return
	}
	c.HTML(200, "renameuser.html", gin.H{"username": uname})
}

// ユーザー一覧をapiから取得しusernameを配列にする
func UserGet() []string {
	// apiからユーザー一覧データの取得
	url := "http://localhost:8081/users"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 送られてきたjsonのopen
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// openしたjsonを構造体にデコード
	var d []User
	if err := json.Unmarshal(body, &d); err != nil {
		log.Fatal(err)
	}

	// デコードしたユーザーデータのNameを配列に入れる
	var userslice []string
	for _, v := range d {
		userslice = append(userslice, v.Name)
	}

	return userslice
}

// ログイン処理
func loginuser(c *gin.Context) {
	session := sessions.Default(c)

	// 入力内容の取得
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 未入力かチェック
	if username == "" || password == "" {
		msg := "入力されてない項目があるよ"
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"message": msg})
		return
	}

	// データベースにユーザーが登録されているか、パスワードが合っているかチェック
	m := LoginCheck(username, password)
	if m == "inai" {
		msg := "そのusernameはいません"
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"message": msg})
		return
	} else if m == "no" {
		msg := "パスワードが間違っています"
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"message": msg})
		return
	}

	// ログインしたユーザーでセッション確立
	session.Set("uname", username)
	session.Save()

	c.Redirect(303, "/menu/top")
}

// ログインチェック
func LoginCheck(username string, password string) string {
	url := "http://localhost:8081/users/showname/" + username

	// ユーザーが存在するかチェック
	b, _ := exec.Command("curl", url, "-X", "GET").Output()
	if len(b) == 2 {
		fmt.Println("そのuserいない")
		msg := "inai"
		return msg
	}

	// apiからユーザーで検索して取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 送られてきたjsonのオープン
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(body))

	// openしたjsonを構造体にデコード
	var d User
	if err := json.Unmarshal(body, &d); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(d)

	// パスワードが一致するか確認
	checkpass := d.Password
	if checkpass != password {
		msg := "no"
		return msg
	}

	msg := "ok"
	return msg
}

// ユーザー新規登録
func signupuser(c *gin.Context) {
	session := sessions.Default(c)

	url := "http://localhost:8081/users"

	// 入力内容の取得
	username := c.PostForm("username")
	password := c.PostForm("password")
	checkpass := c.PostForm("checkpassword")

	// 未入力の確認
	if username == "" || password == "" || checkpass == "" {
		msg := "入力されてない項目があるよ"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	// パスワードと確認用の再入力パスワードが一致するか確認
	if password != checkpass {
		msg := "パスワードが一致していないよ"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	// 入力されたユーザー名が既に登録されているか確認
	if m := AlreadyName(username); m == "aru" {
		msg := "その名前は既にあります"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	// apiへの送信用のjson設定
	jsonStr := `{"Name":"` + username + `","Password":"` + password + `"}`

	// apiへのユーザー情報送信
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

	// サインインしたユーザーでセッション確立
	session.Set("uname", username)
	session.Save()

	c.Redirect(303, "/menu/top")
}

// func aaa(c *gin.Context) {
// 	bb := c.Param("name")
// 	url := "http://localhost:8081/users/showname/" + bb
// 	d, _ := exec.Command("curl", url, "-X", "GET").Output()
// 	// d, _ := http.Get(url)
// 	fmt.Println(d)
// 	if len(d) == 2 {
// 		fmt.Println("gg")
// 	}
// }

// 入力されたユーザーが既に登録されているか確認
func AlreadyName(username string) string {
	url := "http://localhost:8081/users/showname/" + username

	// ユーザーが存在するかチェック
	b, _ := exec.Command("curl", url, "-X", "GET").Output()
	if len(b) == 2 {
		fmt.Println("まだない")
		msg := "no"
		return msg
	}

	fmt.Println("既にある")
	msg := "aru"
	return msg
}

// ユーザーの削除
func deleteuser(c *gin.Context) {
	// ログイン中のユーザー取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	url1 := "http://localhost:8081/users/showname/" + uname

	// ユーザー名からapiでid取得
	resp1, err := http.Get(url1)
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()

	// 取得したjsonのオープン
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal(err)
	}

	// openしたjsonを構造体にデコード
	var d User
	if err := json.Unmarshal(body1, &d); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(d)

	id := strconv.Itoa(d.ID)
	url2 := "http://localhost:8081/users/" + id

	// apiでユーザーの削除
	req2, err := http.NewRequest("DELETE", url2, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp2, err := client.Do(req2)
	if err != nil {
		log.Fatal(err)
	}
	defer resp2.Body.Close()

	// ユーザーのメモとコインの削除
	var derr error
	ms := memo.Memoservice{}
	derr = ms.MemoDelete(uname)
	if derr != nil {
		log.Fatal(err)
	}
	tc := ticketandcoin.TandC{}
	derr = tc.CoinDelete(uname)
	if derr != nil {
		log.Fatal(err)
	}

	// セッションの解除
	session.Clear()
	session.Save()

	// 削除メッセージとゲストに切り替え
	result := "userを削除しました。"
	uname = "ゲスト"
	msg := "現在ゲストで使用しています。ログインしましょう。"
	c.HTML(200, "top.html", gin.H{"user": uname, "message": msg, "result": result})
}

// ユーザー名の変更
func renameuser(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// apiでユーザー名からidの取得
	url1 := "http://localhost:8081/users/showname/" + uname
	resp1, err := http.Get(url1)
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()

	// 取得したjsonのオープン
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		log.Fatal(err)
	}

	// openしたjsonを構造体にデコード
	var d User
	if err := json.Unmarshal(body1, &d); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(d)

	id := strconv.Itoa(d.ID)
	password := d.Password
	// 変更後のユーザー名取得
	rename := c.PostForm("rename")
	url2 := "http://localhost:8081/users/" + id

	// 未入力か確認
	if rename == "" {
		c.Redirect(303, "/menu/settings/renameuser")
		return
	}

	// 登録内容をjsonに設定
	jsonStr := `{"Name":"` + rename + `","Password":"` + password + `"}`

	// apiでユーザー名の更新
	req, err := http.NewRequest(
		"PUT",
		url2,
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

	// ユーザーのコインのname更新
	var perr error
	tc := ticketandcoin.TandC{}
	perr = tc.CoinUserPUT(uname, rename)
	if perr != nil {
		log.Fatal(err)
	}

	// 更新メッセージと更新後のユーザー名でセッション確立
	result := "usernameを更新しました。"
	session.Set("uname", rename)
	session.Save()
	reuname, _ := session.Get("uname").(string)

	c.HTML(200, "top.html", gin.H{"user": reuname, "result": result})
}
