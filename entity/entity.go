package entity

import (
	"github.com/Shakkuuu/gacha-golang/gacha"
)

// userの構造体
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Ticket   int    `json:"ticket"`
	Coin     int    `json:"coin"`
}

// Coin
type Coin struct {
	Name      string  `json:"name"`
	Qty       float32 `json:"qty"`
	Speed     float32 `json:"speed"`
	Speedneed float32 `json:"speedneed"`
}

// memoの構造体
type Memo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type TmpResults struct {
	DB      []*gacha.Card
	One     []string
	Msg     string
	Tickets int
	Coins   int
	Kaisu   int
	Rari    []string
}
