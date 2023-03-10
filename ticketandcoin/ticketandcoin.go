package ticketandcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/Shakkuuu/login-app-api-use/entity"
	"github.com/Shakkuuu/login-app-api-use/minigame"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type User entity.User

type TandC struct{}

func (tc TandC) TicketandCoin(c *gin.Context) {
	// ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// ログインせずにアクセスされた場合のゲストモード
	if uname == "" {
		// uname = "ゲスト"
		// msg := "現在ゲストで使用しています。ログインまたはアカウント登録しましょう。"
		c.HTML(200, "error.html", nil)
		return
	}

	// そのユーザーのコインを取得
	ti, _ := tc.TicketandCoinGet(uname)
	mg := minigame.MG{}
	co := mg.ApiCoinGet(uname)

	c.HTML(200, "ticketandcoin.html", gin.H{"Tickets": ti, "Coins": int(co.Qty)})
}

func (tc TandC) TicketandCoinGet(uname string) (int, int) {
	url := "http://localhost:8081/users/showname/" + uname

	// ユーザーが存在するかチェック
	b, _ := exec.Command("curl", url, "-X", "GET").Output()
	if len(b) == 2 {
		fmt.Println("そのuserいない")
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

	return d.Ticket, d.Coin
}

func (tc TandC) TicketAdd(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/tiadd/" + uname
	http.Get(url)

	c.Redirect(303, "/menu/game/tandc")
}

// func (tc TandC) CoinAdd(c *gin.Context) {
// 	//ログイン中のユーザー名取得
// 	session := sessions.Default(c)
// 	uname, _ := session.Get("uname").(string)

// 	// apiでユーザー名からidの取得
// 	url := "http://localhost:8081/users/coadd/" + uname
// 	http.Get(url)

// 	c.Redirect(303, "/menu/game/tandc")
// }

func (tc TandC) TicketSub(uname string) {

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/tisub/" + uname
	http.Get(url)

}

// func (tc TandC) CoinSub(uname string) {

// 	// apiでユーザー名からidの取得
// 	url := "http://localhost:8081/users/cosub/" + uname
// 	http.Get(url)

// }

// コイン削除
func (tc TandC) CoinDelete(username string) error {
	url := "http://localhost:8081/gamecoin/" + username

	// apiでユーザーの削除
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	return nil
}

// コインuserアップデート
func (tc TandC) CoinUserPUT(username string, rename string) error {
	// 現在のコイン情報の取得
	mg := minigame.MG{}
	coi := mg.ApiCoinGet(username)
	// コイン情報の修正
	coi2 := entity.Coin{
		ID:        coi.ID,
		Name:      rename,
		Qty:       coi.Qty,
		Speed:     coi.Speed,
		Speedneed: coi.Speedneed,
	}

	url := "http://localhost:8081/gamecoin/" + username

	jsonData, err := json.Marshal(coi2)
	if err != nil {
		return err
	}

	// apiでユーザー名の更新
	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
