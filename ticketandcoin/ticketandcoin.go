package ticketandcoin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/Shakkuuu/login-app-api-use/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type User entity.User

type TandC struct{}

func (tc TandC) TicketandCoin(c *gin.Context) {
	// ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// そのユーザーのメモを取得
	ti, co := tc.TicketandCoinGet(uname)

	c.HTML(200, "ticketandcoin.html", gin.H{"Tickets": ti, "Coins": co})
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

	fmt.Println(string(body))

	// openしたjsonを構造体にデコード
	var d User
	if err := json.Unmarshal(body, &d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)

	return d.Ticket, d.Coin
}

func (tc TandC) TicketAdd(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/tiadd/" + uname
	http.Get(url)

	c.Redirect(303, "/menu/tandc")
}

func (tc TandC) CoinAdd(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/coadd/" + uname
	http.Get(url)

	c.Redirect(303, "/menu/tandc")
}

func (tc TandC) TicketSub(uname string) {

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/tisub/" + uname
	http.Get(url)

}

func (tc TandC) CoinSub(uname string) {

	// apiでユーザー名からidの取得
	url := "http://localhost:8081/users/cosub/" + uname
	http.Get(url)

}