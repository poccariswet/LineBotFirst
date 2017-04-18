package main

import(
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  // "net/url"
  "os"
  "strconv"
  "github.com/line/line-bot-sdk-go/linebot"
  "github.com/gin-gonic/gin"

)
// Json core fields
type Wdata struct {
  Weather []Weather `json:"weather"`
  Info     Info     `json:"main"`
}

// Json wether item
type Weather struct{
  Main string `json:"main"`
  Icon string `json:icon`
}

//Info of main item
type Info struct{
  Temp      float32 `json:"temp"`
  Humidity  float32 `json:humidity`
}


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

           //メッセージの受信
           for _, event := range events {
               if event.Type == linebot.EventTypeMessage {
                   switch message := event.Message.(type) {
                   case *linebot.LocationMessage:

                     // 緯度,経度から天気の問い合わせるためのURLを作る
                    //  location, _ := event.handleLocation()
                     lat := strconv.FormatFloat(message.Latitude, 'f', 6, 64)
                     lon := strconv.FormatFloat(message.Longitude, 'f', 6, 64)

                     weather_url := "http://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&appid=b1b15e88fa797225412429c1c50c122a1"

                     //天気情報の取得
                     resp, _ := http.Get(weather_url)
                     defer resp.Body.Close()
                     byteArray, _ := ioutil.ReadAll(resp.Body)
                     jsonBytes := ([]byte)(string(byteArray[:]))

                     weather_data := new(Wdata)
                     if err := json.Unmarshal(jsonBytes, weather_data); err != nil {
                       fmt.Println("JSON Unmarshal error:", err)
                       return
                     }

                     //メッセージの送信

                       if _, err = bot.ReplyMessage(event.ReplyToken,
                         linebot.NewTextMessage("現在の天気をお知らせします。"),
                         linebot.NewTextMessage("天気 : "+ weather_data.Weather[0].Main),
                         linebot.NewImageMessage("http://openweathermap.org/img/w/"+weather_data.Weather[0].Icon+".png", "http://openweathermap.org/img/w/"+weather_data.Weather[0].Icon+".png"),
                         linebot.NewTextMessage("気温 : " + fmt.Sprintf("%.2f", (weather_data.Info.Temp - 273.15))),
                         linebot.NewTextMessage("湿度 : " + fmt.Sprintf("%.2f", (weather_data.Info.Humidity))),
                         ).Do(); err != nil {
                           log.Print(err)
                       }
                     default:
                      if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("位置情報を入力してください" + fmt.Sprintf("%.2f", (message.Latitude)))).Do(); err != nil{
                         log.Print(err)
                       }
                   }
               }
           }
       })

       router.Run(":" + port)
   }
