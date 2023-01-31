# login-app-api-use

このアプリケーションは、github.com/Shakkuuu/login-app-api をAPIとして使用しています。  
APIをdockerで起動させ、localhost:8081/users にアクセスしてGETやPOST、PUT、DELETEを行っています。

## 使い方

1. localサーバを起動  
``` go run main.go ```
2. ブラウザで localhost:8082 にアクセス
3. 登録しているユーザーの一覧が表示される

## ログイン機能

* signupを押すとユーザー登録画面に移る
* username,password,checkpassword(パスワードの再確認用)を入力して登録を押す
* loginを押すとログイン画面に移る
* 既にユーザー登録しているユーザーのusername,passwordを入力してログインを押す
* ユーザー登録またはログインをすると、ユーザーごとのtop画面に移る

## top画面

* メモ登録:ユーザーごとにメモを登録できる
* ガチャ:ガチャを引ける
* チケットとコイン:ガチャで使うチケットとコインを追加する
* logout:ログアウトする
* userの削除:今ログインしているユーザーを削除する
* usernameの更新:今ログインしているユーザーのユーザー名を変える

## メモ機能

* ログイン後のtop画面でメモ登録を押すとメモページに移る
* されぞれのユーザーが登録したメモのタイトル一覧が表示される
* タイトルを押すと本文が表示される
* titleと本文を入力してメモ登録を押すと、メモが登録される

## ガチャ機能

* 上にログイン中のユーザーが持っているチケット数,コイン数と、ガチャが引ける回数が表示される
* 回数を入力してガチャを引くを押すとガチャを引ける
* ガチャ結果に今引いたガチャの結果が表示される
* 結果一覧にこれまでに引かれたガチャの結果とレア度ごとの枚数が全て表示される

## チケット

* ログイン中のユーザーが持っているチケット数とコイン数が表示される
* taddを押すとログイン中のユーザーのチケットの枚数を増やす

## ミニゲーム

* クリック を押すと、コインが増えていく
* スピードアップを押すと必要枚数消費して、1クリックで増えるコインの枚数が増加する
* このコインはガチャで使用できる

## ユーザーの削除とユーザー名の更新

* user削除を押すと、確認画面を挟んでログイン中のユーザーを削除できる
* usernameの更新を押すと、ログイン中のユーザのユーザー名を変えることができる

## 注意事項

* 既に登録されているユーザー名は登録できません
* ユーザー情報やメモはAPI先のmysqlに保存されています
* ガチャの結果はsqliteで保存されています
* results.dbというファイルを削除することで、ガチャの結果一覧がリセットされます
* ガチャの結果一覧は200件までしか表示されません
* アクセスの集中によるセッションの混ざりを防ぐため、ガチャを引くと少し待機時間が入ります

## 開発環境

* macbook air M1
* Visual Studio Code
* go: version go1.19.5 darwin/arm64
* 使用パッケージ  
標準  
``` bytes, fmt, io, log, net/http, os, os/exec, strconv, database/sql, math/rand, time ```  
その他  
```github.com/tenntenn/sqlite, github.com/Shakkuuu/gacha-golang/gacha, github.com/Shakkuuu/login-app-api-use/entity, github.com/Shakkuuu/login-app-api-use/gachagame, github.com/Shakkuuu/login-app-api-use/memo, github.com/Shakkuuu/login-app-api-use/ticketandcoin, github.com/gin-contrib/sessions, github.com/gin-contrib/sessions/cookie, github.com/gin-gonic/gin```

## 今後追加予定の機能

* メモごとの編集と削除
* 現在ガチャ機能は、アクセス時にランダムな時間待機させてからdrawするという無理矢理な方法だが、しっかりとセッションを管理したい
* Ticketの取得方法の調整
* minigameのアップデート
* ユーザーごとのこれまでのガチャ結果の保存
* ガチャで手に入れたものの有効活用
* バリデーション(入力規則など)
