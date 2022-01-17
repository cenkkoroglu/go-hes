package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/template/html"
	"hes/models"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var token = &models.Token{Token: "", PhoneNumber: ""}

func main() {
	GetTokenFromFile()

	engine := html.New("./views", ".html")

	engine.Reload(true)

	engine.Delims("{{", "}}")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/statics", "./statics")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "HES Kodu Sorgulama",
		})
	})

	privateGroup := app.Group("")
	privateGroup.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "123456",
		},
	}))

	privateGroup.Get("/token", func(c *fiber.Ctx) error {
		GetTokenFromFile()
		return c.Render("token", fiber.Map{
			"Title":       "Token Alma",
			"PhoneNumber": token.PhoneNumber,
			"Token":       token.Token,
		})
	})

	privateGroup.Post("/sendLoginCode", SendLoginCode)
	privateGroup.Post("/authenticate", Authenticate)
	privateGroup.Post("/checkHesCode", CheckHesCode)

	log.Fatal(app.Listen(":3000"))
}

func GetTokenFromFile() {
	jsonFile, err := os.Open("token.json")

	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, token)

	defer jsonFile.Close()
}

func WriteTokenToFile() {
	file, _ := json.MarshalIndent(token, "", " ")

	err := ioutil.WriteFile("token.json", file, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func SendLoginCode(ctx *fiber.Ctx) error {
	client := resty.New()

	resp, _ := client.R().
		SetBody(map[string]interface{}{
			"phone": ctx.Query("phoneNumber"),
		}).
		Post("https://hessvc.saglik.gov.tr/api/send-code-to-login")

	return ctx.JSON(&models.Result{Status: resp.StatusCode(), Data: resp.StatusCode() == 201})
}

func CheckHesCode(ctx *fiber.Ctx) error {
	client := resty.New()

	hes := ctx.Query("hes")
	hes = strings.Replace(hes, "-", "", -1)

	resp, _ := client.R().
		SetBody(map[string]interface{}{
			"hes_code": hes,
		}).
		SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token.Token),
		}).
		Post("https://hessvc.saglik.gov.tr/services/hescodeproxy/api/check-hes-code")

	body := resp.Body()
	result := make(map[string]interface{})

	err := json.Unmarshal(body, &result)
	if err != nil {
		return ctx.JSON(&models.Result{Status: 500, Data: "Sistem hatası!"})
	}

	return ctx.JSON(&models.Result{Status: resp.StatusCode(), Data: result})
}

func Authenticate(ctx *fiber.Ctx) error {
	client := resty.New()

	resp, _ := client.R().
		SetBody(map[string]interface{}{
			"phone":      ctx.Query("phoneNumber"),
			"password":   ctx.Query("loginCode"),
			"rememberMe": true,
		}).
		Post("https://hessvc.saglik.gov.tr/api/authenticate-with-code")

	body := resp.Body()
	result := make(map[string]interface{})

	err := json.Unmarshal(body, &result)
	if err != nil {
		return ctx.JSON(&models.Result{Status: 500, Data: "Sistem hatası!"})
	}

	if resp.StatusCode() == 200 {
		token.Token = result["id_token"].(string)
		token.PhoneNumber = ctx.Query("phoneNumber")

		WriteTokenToFile()
	}

	return ctx.JSON(&models.Result{Status: resp.StatusCode(), Data: result})
}
