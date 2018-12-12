package main

import (
	"os"

	"github.com/kamontia/kusoapp/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// インスタンスの作成
	e := echo.New()

	// Echoのミドルウェアを使用する場合は以下のように追記
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/hello", handler.Hello())

	// サーバー起動
	port := os.Getenv("PORT")
	e.Start(":" + port) //ポート番号指定してね
}
