package minigame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/Shakkuuu/login-app-api-use/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var Coin entity.Coin

type MG struct{}

func (mg MG) Minigamemain(c *gin.Context) {
	// ログインしているユーザー名を取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		// uname = "ゲスト"
		// msg := "現在ゲストで使用しています。ログインまたはアカウント登録しましょう。"
		c.HTML(200, "error.html", nil)
		return
	}

	url := "http://localhost:8081/gamecoin/showname/" + uname
	// ユーザーコインを持っているかチェック
	b, _ := exec.Command("curl", url, "-X", "GET").Output()
	if len(b) == 2 {
		fmt.Println("そのuser coinない")
		apicoincreate(uname)
	}

	// 現在のコイン枚数の取得
	coi := mg.ApiCoinGet(uname)
	// fmt.Println(coi)
	c.HTML(200, "minigame.html", gin.H{"Qty": int(coi.Qty), "Speed": coi.Speed, "Speedneed": int(coi.Speedneed)})
}

func (mg MG) Addcoin(c *gin.Context) {
	// ログインしているユーザー名を取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// 現在のコイン枚数の取得
	coi := mg.ApiCoinGet(uname)
	// コインの追加
	qty := coi.Qty + coi.Speed
	coi2 := entity.Coin{
		ID:        coi.ID,
		Name:      coi.Name,
		Qty:       qty,
		Speed:     coi.Speed,
		Speedneed: coi.Speedneed,
	}
	// fmt.Println(coi.Qty)
	// fmt.Println(coi2)
	mg.ApiCoinAdd(uname, coi2)
	c.Redirect(303, "/menu/game/minigame")
}

func (mg MG) AddSpeedUp(c *gin.Context) {
	// ログインしているユーザー名を取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// 現在のコイン枚数の取得
	coi := mg.ApiCoinGet(uname)

	// コイン枚数が必要なコイン枚数を超えているか
	if int(coi.Qty) < int(coi.Speedneed) {
		c.Redirect(303, "/menu/game/minigame")
		return
	}

	// コインを消費して生産速度のアップと必要コイン枚数の増加
	qty := coi.Qty - float32(int(coi.Speedneed))
	sp := coi.Speed + 0.4
	spn := coi.Speedneed * 1.7
	coi2 := entity.Coin{
		ID:        coi.ID,
		Name:      coi.Name,
		Qty:       qty,
		Speed:     sp,
		Speedneed: spn,
	}
	apiaddspeedup(uname, coi2)
	c.Redirect(303, "/menu/game/minigame")
}

func (mg MG) ApiCoinGet(username string) entity.Coin {
	url := "http://localhost:8081/gamecoin/showname/" + username
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
	var co entity.Coin
	if err := json.Unmarshal(body, &co); err != nil {
		log.Fatal(err)
	}

	return co
}

func apicoincreate(username string) {
	url := "http://localhost:8081/gamecoin"

	qty := 0
	speed := 1
	speedneed := 20

	coindata := entity.Coin{
		Name:      username,
		Qty:       float32(qty),
		Speed:     float32(speed),
		Speedneed: float32(speedneed),
	}
	jsonData, err := json.Marshal(coindata)
	if err != nil {
		fmt.Println(err)
		return
	}

	// apiにjsonを送信して登録
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
}

func (mg MG) ApiCoinAdd(uname string, coi2 entity.Coin) {
	url := "http://localhost:8081/gamecoin/" + uname

	jsonData, err := json.Marshal(coi2)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jsonData))

	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(jsonData),
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
}

func apiaddspeedup(uname string, coi2 entity.Coin) {
	url := "http://localhost:8081/gamecoin/" + uname

	jsonData, err := json.Marshal(coi2)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(jsonData))

	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(jsonData),
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
}
