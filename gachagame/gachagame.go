package gachagame

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Shakkuuu/gacha-golang/gacha"
	"github.com/Shakkuuu/login-app-api-use/entity"
	"github.com/Shakkuuu/login-app-api-use/ticketandcoin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tenntenn/sqlite"
)

type TmpResults entity.TmpResults

var db *sql.DB
var err error

var p *gacha.Player
var play *gacha.Play

// var numnum int
var onere []string

// var rere []string
var msg string

var ticket int
var coin int
var kai int

type Gachaservice struct{}

func (gg Gachaservice) GachaCreate() error {
	db, err = sql.Open(sqlite.DriverName, "results.db")
	if err != nil {
		return fmt.Errorf("データベースのOpen:%w", err)
	}

	if err := CreateTable(db); err != nil {
		return err
	}

	return nil
}

func (gg Gachaservice) TCSet(uname string) {
	tc := ticketandcoin.TandC{}
	ticket, coin = tc.TicketandCoinGet(uname)
	// チケット数とコイン数の設定
	p = gacha.NewPlayer(ticket, coin)
}

func (gg Gachaservice) GachaGame(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	play = gacha.NewPlay(p)

	results, err := getResults(db, 200)

	// ti, co := p.Maisu()
	// kai := p.DrawableNum()
	tc := ticketandcoin.TandC{}
	ticket, coin = tc.TicketandCoinGet(uname)
	kai = ticket + coin/10
	fmt.Printf("チケット:%d コイン:%d 引ける回数:%d \n", ticket, coin, kai)

	reamap := map[gacha.Rarity]int{}
	for _, reav := range results {
		reamap[reav.Rarity]++
	}
	var rea []string
	for rarity, count := range reamap {
		countStr := strconv.Itoa(count)
		rea = append(rea, rarity.String()+":"+countStr)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	rr := TmpResults{
		DB:      results,
		One:     onere,
		Msg:     msg,
		Tickets: ticket,
		Coins:   coin,
		Kaisu:   kai,
		Rari:    rea,
	}

	fmt.Println(rr)

	c.HTML(200, "gacha.html", gin.H{"DB": rr.DB, "One": rr.One, "Msg": rr.Msg, "Tickets": rr.Tickets, "Coins": rr.Coins, "Kaisu": rr.Kaisu, "Rari": rr.Rari})
	// c.HTML(200, "gacha.html", nil)

	msg = ""
	onere = nil
}

func (gg Gachaservice) DrawGacha(c *gin.Context) {
	//ログイン中のユーザー名取得
	session := sessions.Default(c)
	uname, _ := session.Get("uname").(string)

	num, err := strconv.Atoi(c.PostForm("num"))
	// kai = p.DrawableNum()
	if kai == 0 {
		msg = "チケットあるいはコインがありません"
		c.Redirect(303, "/menu/game/gachagame")
		return
	}
	if num > kai {
		// fmt.Println("引ける回数を超えてます")
		msg = "引ける回数を超えてます"
		c.Redirect(303, "/menu/game/gachagame")
		return
	}
	// numnum = num
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < num; i++ {
		if !play.Draw() {
			if err := saveResult(db, play.Result()); err != nil {
				fmt.Println(err)
				return
			}

			subdraw(uname)

			onere = append(onere, play.Result().String())
			break
		}

		if err := saveResult(db, play.Result()); err != nil {
			fmt.Println(err)
			return
		}

		subdraw(uname)

		onere = append(onere, play.Result().String())
	}

	if err := play.Err(); err != nil {
		fmt.Println(err)
		return
	}

	c.Redirect(303, "/menu/game/gachagame")
}

func subdraw(uname string) {
	tc := ticketandcoin.TandC{}
	tic, _ := tc.TicketandCoinGet(uname)
	if tic > 0 {
		tc.TicketSub(uname)
		return
	}
	for range make([]struct{}, 10) {
		tc.CoinSub(uname) // 1回あたり10枚消費する
	}
}

func CreateTable(db *sql.DB) error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS results(
		id        INTEGER PRIMARY KEY,
		rarity	  TEXT NOT NULL,
		name      TEXT NOT NULL
	);`

	_, err := db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("テーブル作成:%w", err)
	}

	return nil
}

func saveResult(db *sql.DB, card *gacha.Card) error {
	const sqlStr = `INSERT INTO results(rarity, name) VALUES (?,?);`

	_, err := db.Exec(sqlStr, card.Rarity.String(), card.Name)
	if err != nil {
		return err
	}
	return nil
}

func getResults(db *sql.DB, limit int) ([]*gacha.Card, error) {
	const sqlStr = `SELECT rarity, name FROM results LIMIT ?`
	rows, err := db.Query(sqlStr, limit)
	if err != nil {
		return nil, fmt.Errorf("%qの実行:%w", sqlStr, err)
	}
	defer rows.Close()

	var results []*gacha.Card
	for rows.Next() {
		var card gacha.Card
		err := rows.Scan(&card.Rarity, &card.Name)
		if err != nil {
			return nil, fmt.Errorf("Scan:%w", err)
		}
		results = append(results, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("結果の取得:%w", err)
	}

	return results, nil
}
