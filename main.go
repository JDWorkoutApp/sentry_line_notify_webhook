package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type sentryData struct {
	Url   string `json:"url" form:"url" binding:"required"`
	Event struct {
		Title string `json:"title" form:"title" binding:"required"`
	} `json:"event" form:"event" binding:"required"`
}

func registerSentry(r *gin.RouterGroup) {
	r.POST("/sentry", receiveSentry)
}

func receiveSentry(c *gin.Context) {
	sentryData := sentryData{}
	if err := c.ShouldBind(&sentryData); err != nil {
		println(err.Error())
		return
	}

	lineNotify("api got error, " + sentryData.Event.Title + ", go check, url: " + sentryData.Url)
}

func lineNotify(message string) {
	requestBody := &bytes.Buffer{}

	formWriter := multipart.NewWriter(requestBody)

	messageField, err := formWriter.CreateFormField("message")
	if err != nil {
		fmt.Println("Error creating form field:", err)
		return
	}
	_, err = messageField.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing form field:", err)
		return
	}

	formWriter.Close()

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	token := os.Getenv("LINE_NOTIFY_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", formWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status code:", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response body:", string(body))
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	r := gin.Default()
	webhookGroup := r.Group("/webhook")
	registerSentry(webhookGroup)
	r.Run(":80")
}
