package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Memo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (u *User) UnmarshalJSON(body []byte) error {
	// 自分で新しく定義した構造体
	u2 := &struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	err := json.Unmarshal(body, u2)
	if err != nil {
		panic(err)
	}
	// 新しく定義した構造体の結果をもとのiに詰める
	u.ID = u2.ID
	u.Name = u2.Name
	u.Password = u2.Password

	return err
}

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
	// 新しく定義した構造体の結果をもとのiに詰める
	m.ID = m2.ID
	m.Name = m2.Name
	m.Title = m2.Title
	m.Text = m2.Text

	return err
}

func main() {
	fmt.Println("start")

	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", index)
	r.GET("/login", login)
	r.GET("/signup", signup)
	r.GET("/logout", logout)

	r.POST("/login", loginuser)
	r.POST("/signup", signupuser)

	menu := r.Group("/menu")
	menu.GET("/top", top)
	menu.GET("/memo", memo)
	menu.POST("/memo", memocreate)

	settings := menu.Group("/settings")
	settings.GET("/deleteuser", deleteusercheck)
	settings.GET("/renameuser", renameusercheck)

	settings.POST("/deleteuser", deleteuser)
	settings.POST("/renameuser", renameuser)

	r.Run(":8082")
}

func top(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	// if err != false {
	// 	msg := "ログイン失敗?"
	// 	c.HTML(200, "index.html", gin.H{"message": msg})
	// }
	if uname == "" {
		uname = "ゲスト"
		msg := "現在ゲストで使用しています。ログインしましょう。"
		c.HTML(200, "top.html", gin.H{"user": uname, "message": msg})
		return
	}
	c.HTML(200, "top.html", gin.H{"user": uname})
}

func index(c *gin.Context) {
	u := UserGet()
	c.HTML(200, "index.html", gin.H{"users": u})
}

func login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func signup(c *gin.Context) {
	c.HTML(200, "signup.html", nil)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(303, "/login")
}

func deleteusercheck(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	c.HTML(200, "deleteuser.html", gin.H{"username": uname})
}

func renameusercheck(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	c.HTML(200, "renameuser.html", gin.H{"username": uname})
}

func UserGet() []string {
	url := "http://localhost:8081/users"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(body))

	var d []User

	if err := json.Unmarshal(body, &d); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(d)

	var userslice []string

	for _, v := range d {
		userslice = append(userslice, v.Name)
	}

	return userslice

}

func loginuser(c *gin.Context) {
	session := sessions.Default(c)

	// url := "http://localhost:8081/users"
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		msg := "入力されてない項目があるよ"
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"message": msg})
		return
	}

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

	session.Set("uname", username)
	session.Save()

	// result := username + "でログインしました。"

	// c.HTML(200, "index.html", gin.H{"result": result})
	c.Redirect(303, "/menu/top")
}

func LoginCheck(username string, password string) string {
	url := "http://localhost:8081/users/showname/" + username

	b, _ := exec.Command("curl", url, "-X", "GET").Output()

	if len(b) == 2 {
		// fmt.Println(b)
		// fmt.Println(err)
		fmt.Println("そのuserいない")
		msg := "inai"
		return msg
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// msg := "id見つからないよ"
		// c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
		log.Fatal(err)
		// return
	}

	fmt.Println(string(body))

	var d User

	if err := json.Unmarshal(body, &d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)

	checkpass := d.Password
	if checkpass != password {
		msg := "no"
		return msg
	}

	msg := "ok"
	return msg
}

func signupuser(c *gin.Context) {
	session := sessions.Default(c)

	url := "http://localhost:8081/users"
	username := c.PostForm("username")
	password := c.PostForm("password")
	checkpass := c.PostForm("checkpassword")

	if username == "" || password == "" || checkpass == "" {
		msg := "入力されてない項目があるよ"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	if password != checkpass {
		msg := "パスワードが一致していないよ"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	// aaa := AlreadyName(username)
	// fmt.Println(aaa)
	if m := AlreadyName(username); m == "aru" {
		// msg := m
		msg := "その名前は既にあります"
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"message": msg})
		return
	}

	jsonStr := `{"Name":"` + username + `","Password":"` + password + `"}`

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

	session.Set("uname", username)
	session.Save()

	c.Redirect(303, "/menu/top")

	// result := username + "を登録しました。"

	// c.HTML(200, "index.html", gin.H{"result": result})
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

func AlreadyName(username string) string {
	url := "http://localhost:8081/users/showname/" + username

	b, _ := exec.Command("curl", url, "-X", "GET").Output()

	if len(b) == 2 {
		// fmt.Println(b)
		// fmt.Println(err)
		fmt.Println("まだない")
		msg := "no"
		return msg
	}

	fmt.Println("既にある")
	msg := "aru"
	return msg
}

func deleteuser(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	url1 := "http://localhost:8081/users/showname/" + uname

	resp1, err := http.Get(url1)
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		// msg := "id見つからないよ"
		// c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
		log.Fatal(err)
		// return
	}

	var d User

	if err := json.Unmarshal(body1, &d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)

	id := strconv.Itoa(d.ID)
	url2 := "http://localhost:8081/users/" + id

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

	session.Clear()
	session.Save()

	result := "userを削除しました。"

	uname = "ゲスト"
	msg := "現在ゲストで使用しています。ログインしましょう。"
	c.HTML(200, "top.html", gin.H{"user": uname, "message": msg, "result": result})

	// c.HTML(200, "index.html", gin.H{"result": result})
}

func renameuser(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	url1 := "http://localhost:8081/users/showname/" + uname

	resp1, err := http.Get(url1)
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer resp1.Body.Close()

	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		// msg := "id見つからないよ"
		// c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
		log.Fatal(err)
		// return
	}

	var d User

	if err := json.Unmarshal(body1, &d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)

	id := strconv.Itoa(d.ID)
	password := d.Password
	rename := c.PostForm("rename")
	url2 := "http://localhost:8081/users/" + id

	if rename == "" {
		// msg := "入力されてない項目があるよ"
		// c.HTML(200, "error.html", gin.H{"message": msg})
		c.Redirect(303, "/menu/settings/renameuser")
		return
	}

	jsonStr := `{"Name":"` + rename + `","Password":"` + password + `"}`
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

	result := "usernameを更新しました。"

	session.Set("uname", rename)
	session.Save()
	reuname, _ := session.Get("uname").(string)

	c.HTML(200, "top.html", gin.H{"user": reuname, "result": result})
}

func memo(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)
	m := MemoGet(uname)
	c.HTML(200, "memo.html", gin.H{"memos": m})
}

func MemoGet(uname string) []Memo {
	url := "http://localhost:8081/memos/showname/" + uname
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(body))

	var m []Memo

	if err := json.Unmarshal(body, &m); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(d)

	return m

	// var titleslice []string

	// for _, v := range m {
	// 	titleslice = append(titleslice, v.Title)
	// }

	// return titleslice

}

func memocreate(c *gin.Context) {
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	url := "http://localhost:8081/memos"
	username := uname
	title := c.PostForm("title")
	text := c.PostForm("text")

	if username == "" || title == "" || text == "" {
		// msg := "入力されてない項目があるよ"
		// c.HTML(http.StatusBadRequest, "memo.html", gin.H{"message": msg})
		// return
		c.Redirect(303, "/menu/memo")
		return
	}

	jsonStr := `{"Name":"` + username + `","Title":"` + title + `","Text":"` + text + `"}`

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

	// result := username + "を登録しました。"

	// c.HTML(200, "index.html", gin.H{"result": result})
}

// func idshow(c *gin.Context) {
// 	msg := "見つかりました"
// 	id := c.Query("id")
// 	if id == "" {
// 		msg = "idが入力されていないよ"
// 		c.HTML(200, "error.html", gin.H{"message": msg})
// 		return
// 	}

// 	url := "http://localhost:8080/bihins/showid/" + id
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		msg := "id見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		msg := "id見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}

// 	fmt.Println(string(body))

// 	var d User

// 	if err := json.Unmarshal(body, &d); err != nil {
// 		msg := "id見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}
// 	// fmt.Println(d.CreatedAt.Format())
// 	// fmt.Println(d)
// 	// fmt.Printf("%+v\n", d)
// 	// c.JSON(http.StatusOK, gin.H{"item": d})
// 	c.HTML(200, "list.html", gin.H{"idfind": d, "message": msg})
// }

// func danshow(c *gin.Context) {
// 	msg := "見つかりました"
// 	dan := c.Query("dantai")
// 	if dan == "" {
// 		msg = "danが入力されていないよ"
// 		c.HTML(200, "error.html", gin.H{"message": msg})
// 		return
// 	}

// 	url := "http://localhost:8080/bihins/showdan/" + dan
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		msg := "dan見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		msg := "id見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}

// 	fmt.Println(string(body))

// 	var d []User

// 	if err := json.Unmarshal(body, &d); err != nil {
// 		msg := "id見つからないよ"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}
// 	fmt.Println(d)
// 	// fmt.Printf("%+v\n", d)
// 	// c.JSON(http.StatusOK, gin.H{"item": d})
// 	c.HTML(200, "list.html", gin.H{"danfind": d, "message": msg})
// }

// func listbihin(c *gin.Context) {
// 	url := "http://localhost:8080/bihins"
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		msg := "一覧の取得ができなかった"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		msg := "一覧の読み込みができなかった"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}

// 	fmt.Println(string(body))

// 	var d []User

// 	if err := json.Unmarshal(body, &d); err != nil {
// 		msg := "一覧の変換ができなかった"
// 		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
// 		// log.Fatal(err)
// 		return
// 	}

// 	fmt.Println(d)
// 	// fmt.Printf("%+v\n", d)
// 	// c.JSON(http.StatusOK, gin.H{"item": d})
// 	c.HTML(200, "list.html", gin.H{"bihins": d})
// }

// func deletecheck(c *gin.Context) {
// 	id := c.Param("id")
// 	// id := c.Param("id")

// 	url := "http://localhost:8080/bihins/showid/" + id
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(string(body))

// 	var d User

// 	if err := json.Unmarshal(body, &d); err != nil {
// 		if err, ok := err.(*json.SyntaxError); ok {
// 			fmt.Println(string(body[err.Offset-1:]))
// 		}
// 		// 2009/11/10 23:00:00 json: cannot unmarshal string into Go struct field .C of type int
// 		log.Fatal(err)
// 	}

// 	// fmt.Println(d)
// 	// fmt.Printf("%+v\n", d)
// 	// c.JSON(http.StatusOK, gin.H{"item": d})
// 	c.HTML(200, "delete.html", gin.H{"list": d})
// }

// func deletebihin(c *gin.Context) {
// id := c.Param("id")
// // id := c.Param("id")
// url := "http://localhost:8080/bihins/" + id

// req, err := http.NewRequest("DELETE", url, nil)
// if err != nil {
// 	log.Fatal(err)
// }

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	log.Fatal(err)
// }
// defer resp.Body.Close()

// result := "削除しました。"

// c.HTML(200, "index.html", gin.H{"result": result})
// }

// func koushincheck(c *gin.Context) {
// 	id := c.Param("id")
// 	// id := c.Param("id")

// 	url := "http://localhost:8080/bihins/showid/" + id
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(string(body))

// 	var d User

// 	if err := json.Unmarshal(body, &d); err != nil {
// 		if err, ok := err.(*json.SyntaxError); ok {
// 			fmt.Println(string(body[err.Offset-1:]))
// 		}
// 		log.Fatal(err)
// 	}
// 	// fmt.Println(d)
// 	// fmt.Printf("%+v\n", d)
// 	// c.JSON(http.StatusOK, gin.H{"item": d})
// 	c.HTML(200, "koushin.html", gin.H{"list": d})
// }

// func putbihin(c *gin.Context) {
// 	id := c.PostForm("id")
// 	url := "http://localhost:8080/bihins/" + id
// 	knrid := c.PostForm("knrid")
// 	bihin := c.PostForm("bihin")
// 	dantai := c.PostForm("dantai")
// 	place := c.PostForm("place")
// 	price := c.PostForm("price")
// 	qty := c.PostForm("qty")
// 	partnum := c.PostForm("partnum")
// 	note := c.PostForm("note")
// 	shishutsuid := c.PostForm("shishutsuid")

// if knrid == "" || bihin == "" || dantai == "" || place == "" || price == "" || qty == "" || partnum == "" || note == "" || shishutsuid == "" {
// 	// msg := "入力されてない項目があるよ"
// 	// c.HTML(200, "error.html", gin.H{"message": msg})
// 	c.Redirect(303, "/koushincheck/"+id)
// 	return
// }

// jsonStr := `{"KnrId":"` + knrid + `","Bihin":"` + bihin + `","Dantai":"` + dantai + `","Place":"` + place + `","Price":"` + price + `","Qty":"` + qty + `","PartNum":"` + partnum + `","Note":"` + note + `","Shishutsuid":"` + shishutsuid + `"}`

// req, err := http.NewRequest(
// 	"PUT",
// 	url,
// 	bytes.NewBuffer([]byte(jsonStr)),
// )
// if err != nil {
// 	log.Fatal(err)
// }

// // Content-Type 設定
// req.Header.Set("Content-Type", "application/json")

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	log.Fatal(err)
// }
// defer resp.Body.Close()

// result := "更新しました。"

// c.HTML(200, "index.html", gin.H{"result": result})
// }
