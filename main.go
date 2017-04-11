package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"
    "fmt"
    "time"
    "math/rand"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())


    router.POST("/hook", func(c *gin.Context) {
        client := &http.Client{Timeout: time.Duration(15 * time.Second)}
        bot, err := linebot.New(`
159e31a2f4a0d18f9d8e7eea7f10ef53`,
 `4L6B+BeJZVGyU4syjT5LJl9tzhlYbHbmIk9Q72t1b0JrevzLlceQiENBfZ40vXQXo7QKMVp+GaCBqgizgY4wu8+evXyieDq8Q36ZRzK8S8NqywUgIi4DWU5pg9Pe/xZqlJAmqp6Oi9b0IxWQ4BdPJwdB04t89/1O/w1cDnyilFU=`,
  linebot.WithHTTPClient(client))
        if err != nil {
            fmt.Println(err)
            return
        }
        received, err := bot.ParseRequest(c.Request)

        for _, event := range received {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                    source := event.Source
                    if source.Type == linebot.EventSourceTypeRoom {
                        if resMessage := getResMessage(message.Text); resMessage != "" {
                            postMessage := linebot.NewTextMessage(resMessage)
                            if _, err = bot.ReplyMessage(event.ReplyToken, postMessage).Do(); err != nil {
                                log.Print(err)
                            }
                        }
                    }
                }
            }
        }
    })

    router.Run(":" + port)
}


func getResMessage(reqMessage string) (message string) {
    resMessages := []string{"ねぇ〜〜","ねぇ、聞いて！！","おいっ！","タバコ吸ってくるわ","二郎行こうぜ〜","カビパンばかうめぇ"}

    rand.Seed(time.Now().UnixNano())
    if rand.Intn(5) == 0 {
        if math := rand.Intn(4); math != 3 {
            message = resMessages[math];
        } else {
            message = reqMessage + "じゃねーよ！櫻井だよ！"
        }
    }
    return
}
