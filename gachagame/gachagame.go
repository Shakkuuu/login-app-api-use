package gachagame

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Shakkuuu/gacha-golang/gacha"
	"github.com/Shakkuuu/login-app-api-use/entity"
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

type Gachaservice struct{}

func (gg Gachaservice) GachaCreate() error {
	db, err = sql.Open(sqlite.DriverName, "results.db")
	if err != nil {
		return fmt.Errorf("データベースのOpen:%w", err)
	}

	if err := CreateTable(db); err != nil {
		return err
	}

	tickets := 20
	coin := 100
	// チケット数とコイン数の設定
	p = gacha.NewPlayer(tickets, coin)

	return nil
}

func (gg Gachaservice) GachaGame(c *gin.Context) {
	play = gacha.NewPlay(p)

	results, err := getResults(db, 200)

	ti, co := p.Maisu()
	kai := p.DrawableNum()
	fmt.Printf("チケット:%d コイン:%d 引ける回数:%d \n", ti, co, kai)

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
		Tickets: ti,
		Coins:   co,
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
	num, err := strconv.Atoi(c.PostForm("num"))
	kai := p.DrawableNum()
	if kai == 0 {
		msg = "チケットあるいはコインがありません"
		c.Redirect(303, "/menu/gachagame")
		return
	}
	if num > kai {
		// fmt.Println("引ける回数を超えてます")
		msg = "引ける回数を超えてます"
		c.Redirect(303, "/menu/gachagame")
		return
	}
	// numnum = num
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < num; i++ {
		if !play.Draw() {
			break
		}

		if err := saveResult(db, play.Result()); err != nil {
			fmt.Println(err)
			return
		}

		onere = append(onere, play.Result().String())
	}

	if err := play.Err(); err != nil {
		fmt.Println(err)
		return
	}

	c.Redirect(303, "/menu/gachagame")
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
