package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/kamontia/kusoapp/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uetchy/go-qiita/qiita"
	"golang.org/x/oauth2"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// インスタンスの作成
	e := echo.New()

	// Echoのミドルウェアを使用する場合は以下のように追記
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/hello", handler.Hello())
	e.GET("/auth", handler.Authentication())
	e.GET("/auth/qiita/callback", handler.Callback())
	e.GET("/", func(c echo.Context) error {
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

		return c.Render(http.StatusOK, "top.html", map[string]interface{}{
			"url": url,
		})
	})
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer

	// サーバー起動
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port)) //ポート番号指定してね

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "a8ed99621b01629e94e966e33f6f8ce72b52a020"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Create Qiita client using OAuth2 adapter
	client := qiita.NewClient(tc)

	// Fetch articles and print them
	items, _, _ := client.Items.List(&qiita.ItemsListOptions{Query: "kamontia"})
	fmt.Println(items)
}
