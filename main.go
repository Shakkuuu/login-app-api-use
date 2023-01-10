package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
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

func main() {
	fmt.Println("start")

	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")

	r.GET("/", index)
	// r.GET("/idshow", idshow)
	// r.GET("/danshow", danshow)
	// r.GET("/listbihin", listbihin)
	// r.GET("/koushincheck/:id", koushincheck)
	// r.GET("/deletecheck/:id", deletecheck)
	r.POST("/postuser", postuser)
	// r.POST("/deletebihin/:id", deletebihin)
	// r.POST("/putbihin", putbihin)
	// r.GET("/aaa/:name", aaa)

	r.Run(":8082")
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
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

func postuser(c *gin.Context) {
	url := "http://localhost:8081/users"
	username := c.PostForm("username")
	password := c.PostForm("password")
	checkpass := c.PostForm("checkpassword")

	if username == "" || password == "" {
		msg := "入力されてない項目があるよ"
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"message": msg})
		return
	}

	if password != checkpass {
		msg := "パスワードが一致していないよ"
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"message": msg})
		return
	}

	// aaa := AlreadyName(username)
	// fmt.Println(aaa)
	if m := AlreadyName(username); m == "aru" {
		// msg := m
		msg := "その名前は既にあります"
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"message": msg})
		return
	}

	jsonStr := `{"Name":"` + username + `","Password":"` + password + `"}`

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		msg := "登録できなかった"
		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
		// log.Fatal(err)
		return
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := "登録できなかった"
		c.HTML(200, "error.html", gin.H{"err": err, "message": msg})
		// log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	result := username + "を登録しました。"

	c.HTML(200, "index.html", gin.H{"result": result})
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

	// fmt.Println(username)
	// b, err := http.Get(url)
	// fmt.Println(b)
	// fmt.Println(err)

	// if b == nil {
	// 	// fmt.Println(b)
	// 	fmt.Println(err)
	// 	fmt.Println("まだない")
	// 	// msg := "ng"
	// 	return err
	// }

	fmt.Println("既にある")
	msg := "aru"
	return msg
}

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
// 	id := c.Param("id")
// 	// id := c.Param("id")
// 	url := "http://localhost:8080/bihins/" + id

// 	req, err := http.NewRequest("DELETE", url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	result := "削除しました。"

// 	c.HTML(200, "index.html", gin.H{"result": result})
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

// 	if knrid == "" || bihin == "" || dantai == "" || place == "" || price == "" || qty == "" || partnum == "" || note == "" || shishutsuid == "" {
// 		// msg := "入力されてない項目があるよ"
// 		// c.HTML(200, "error.html", gin.H{"message": msg})
// 		c.Redirect(303, "/koushincheck/"+id)
// 		return
// 	}

// 	jsonStr := `{"KnrId":"` + knrid + `","Bihin":"` + bihin + `","Dantai":"` + dantai + `","Place":"` + place + `","Price":"` + price + `","Qty":"` + qty + `","PartNum":"` + partnum + `","Note":"` + note + `","Shishutsuid":"` + shishutsuid + `"}`

// 	req, err := http.NewRequest(
// 		"PUT",
// 		url,
// 		bytes.NewBuffer([]byte(jsonStr)),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Content-Type 設定
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	result := "更新しました。"

// 	c.HTML(200, "index.html", gin.H{"result": result})
// }
