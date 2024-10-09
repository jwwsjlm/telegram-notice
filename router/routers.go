package router

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v3"
	Tbot "telegram-notice/Newtele"
	uhash "telegram-notice/hash"
)

// SetupRoutes 设置 Gin 的路由
func SetupRoutes(hashMap *uhash.Hashtable, t *Tbot.Telegramini) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(LoggerWithFormatter))
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("gin日志 %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// 设置 POST 和 GET 的 WebHook 路由
	r.POST("/webhook/:id", PostWebHook(hashMap, t))
	r.GET("/webhook/:id", GetWebHook(hashMap, t))
	return r
}

// LoggerWithFormatter 自定义日志格式
func LoggerWithFormatter(params gin.LogFormatterParams) string {
	return fmt.Sprintf(
		"%s | %d | %s | %s | %s | %s\n",
		params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		params.StatusCode, // 状态码
		params.ClientIP,   // 客户端 IP
		params.Latency,    // 请求耗时
		params.Method,     // 请求方法
		params.Path,       // 请求路径
	)
}

// PostWebHook 处理 POST 请求的 WebHook
func PostWebHook(hashMap *uhash.Hashtable, t *Tbot.Telegramini) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// 通过 ID 判断 hashMap 中是否存在该用户
		idInt, err := hashMap.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"message": err.Error()})
			return
		}

		rawData, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		data := string(rawData)
		fmt.Println("用户", idInt)
		fmt.Println("消息内容", data)

		user := &tele.User{ID: idInt}
		_, err = t.Bot.Send(user, data, &tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})

		if err != nil {
			c.JSON(200, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "ok"})
	}
}

// GetWebHook 处理 GET 请求的 WebHook
func GetWebHook(hashMap *uhash.Hashtable, t *Tbot.Telegramini) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("id", id)

		// 通过 ID 判断 hashMap 中是否存在该用户
		idInt, err := hashMap.Get(id)
		if err != nil {

			fmt.Println(err)
			c.JSON(404, gin.H{"message": err.Error()})
			return
		}

		rawData := c.Query("text")
		fmt.Println("消息内容", rawData)

		user := &tele.User{ID: idInt}
		_, err = t.Bot.Send(user, rawData, &tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})

		if err != nil {
			c.JSON(404, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "ok"})
	}
}
