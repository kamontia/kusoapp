package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

//\"client_id\":\"74eaa54fa0ea33f142e2df14b7f43d2a968f5fbc\",\"scopes\":[\"read_qiita\"],\"token\":\"075ecc7e0f952de69459609e0e480cae74eba3a9\"
type Token struct {
	Client_id string `json:"client_id"`
	Scopes    string `json:"scopes"`
	Token     string `json:"token"`
}

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!!!")
	}
}

func Callback() echo.HandlerFunc {
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
			RedirectURL: "https://kusoapp.herokuapp.com/callback",
		}
		fmt.Print(conf)
		code := c.QueryParam("code")

		// OAuth library is disavairable to send as content-type of application/json
		// TRY IMPLEMENT MANUALLY
		jsonStr := `{
			"grant_type": "authorization_code",
			"client_id":"` + conf.ClientID + `",
			"code":"` + code + `",
			"client_secret":"` + conf.ClientSecret + `",
			"scope": "read_qiita",
			"state": "kusokuso"
			}`
		// "redirect_uri":` + conf.RedirectURL + `",
		// jsonStr := url.Values{}
		// jsonStr.Add("client_id", code)
		// jsonStr.Add("state", "kusokuso")
		// jsonStr.Add("scope", "read_qiita")
		fmt.Println(jsonStr)
		url := "https://qiita.com/api/v2/access_tokens"
		request, err := http.NewRequest(
			"POST",
			url,
			// strings.NewReader(jsonStr.Encode()),
			bytes.NewBuffer([]byte(jsonStr)),
		)
		if err != nil {
			fmt.Println("Failed ready to request")
			return err
		}
		request.Header.Set("Content-Type", "application/json")
		fmt.Printf("Request::%T,%v", request, request)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Failed to request")
			return err
		}
		defer response.Body.Close()

		byteArray, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(byteArray))

		// token, err := conf.Exchange(context.Background(), code)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		//var info []Token
		// json.Unmarshal(byteArray, &Token)
		resp := string(byteArray)
		s, _ := json.Marshal(resp)
		fmt.Printf("Access Token: %T %T %v\n", byteArray, resp, string(s))
		// token := string(s)

		jsonStr = `
			"Authorization": "Bearer ` + string(s) + `",
		`

		// アクセストークンを利用して、各種APIにアクセスする
		// client := conf.Client(context.Background(), token)
		// response, err := client.Get("https://qiita.com/api/v2/authenticated_user")

		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer response.Body.Close()
		// byteArray, _ := ioutil.ReadAll(response.Body)

		// fmt.Println(string(byteArray))

		return c.String(http.StatusOK, "Callback")

	}
}
