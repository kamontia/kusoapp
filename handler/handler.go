package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!!!")
	}
}

func Authentication() echo.HandlerFunc {

	return func(c echo.Context) error {

		conf := oauth2.Config{
			ClientID:     "74eaa54fa0ea33f142e2df14b7f43d2a968f5fbc",
			ClientSecret: "f7ba57bdb7e1ed760b0c513e92be91f25ff0e4f9",
			Scopes: []string{
				"read_qiita",
				"read_qiita_team"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://qiita.com/api/v2/oauth/authorize",
				TokenURL: "https://qiita.com/api/v2/access_tokens",
			},
			RedirectURL: "http://localhost:3001/callback",
		}

		// アクセス用のURLの発行
		fmt.Println("DISPLAY ACCESS URL")
		authCode := "kusokuso"
		url := conf.AuthCodeURL(authCode)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
		// アクセストークンの取得
		fmt.Println("GET ACCESS TOKEN")
		fmt.Println("Input Code#: ")
		// var code string
		// if _, err := fmt.Scan(&code); err != nil {
		// 	log.Fatal(err)
		// }
		code := "a8ed99621b01629e94e966e33f6f8ce72b52a020"
		fmt.Println("You input: ", code)
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Access Token: %v\n", token)

		// アクセストークンを利用して、各種APIにアクセスする
		client := conf.Client(context.Background(), token)
		response, err := client.Get("https://qiita.com/api/v2/authenticated_user")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		byteArray, _ := ioutil.ReadAll(response.Body)

		fmt.Println(string(byteArray))

		return c.String(http.StatusOK, "HEY")
	}
}

func Callback() echo.HandlerFunc {
	return func(c echo.Context) error {

		fmt.Println(c.QueryParams())
		fmt.Println(c.QueryParam("code"))

		conf := oauth2.Config{
			ClientID:     "74eaa54fa0ea33f142e2df14b7f43d2a968f5fbc",
			ClientSecret: "f7ba57bdb7e1ed760b0c513e92be91f25ff0e4f9",
			Scopes: []string{
				"read_qiita",
				"read_qiita_team"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://qiita.com/api/v2/oauth/authorize",
				TokenURL: "https://qiita.com/api/v2/access_tokens",
			},
			RedirectURL: "http://localhost:3001/callback",
		}

		code := c.QueryParam("code")
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Access Token: %v\n", token)

		// アクセストークンを利用して、各種APIにアクセスする
		client := conf.Client(context.Background(), token)
		response, err := client.Get("https://qiita.com/api/v2/authenticated_user")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		byteArray, _ := ioutil.ReadAll(response.Body)

		fmt.Println(string(byteArray))

		return c.String(http.StatusOK, "Callback")
	}
}
