package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"  // ① SDKを追加
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    // ② LINE bot instanceの作成
    bot, err := linebot.New(
        os.Getenv("ENV_LINE_SECRET"),
        os.Getenv("ENV_LINE_TOKEN"),
    )
    if err != nil {
        log.Fatal(err)
    }

    router := gin.New()
    router.Use(gin.Logger())
    router.LoadHTMLGlob("templates/*.tmpl.html")
    router.Static("/static", "static")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })
    // ③ LINE Messaging API用の Routing設定
       router.POST("/callback", func(c *gin.Context) {
           events, err := bot.ParseRequest(c.Request)
           if err != nil {
               if err == linebot.ErrInvalidSignature {
                   log.Print(err)
               }
               return
           }
           for _, event := range events {
               if event.Type == linebot.EventTypeMessage {
                   switch message := event.Message.(type) {
                   case *linebot.TextMessage:
                       if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
                           log.Print(err)
                       }
                   }
               }
           }
       })

       router.Run(":" + port)
   }
