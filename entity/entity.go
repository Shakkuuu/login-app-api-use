package entity

import (
	"github.com/Shakkuuu/gacha-golang/gacha"
)

// userの構造体
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
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
