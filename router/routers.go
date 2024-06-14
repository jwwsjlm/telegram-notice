package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v3"
	"log"
	Tbot "telegram-notice/Newtele"
	uhash "telegram-notice/hash"
)

func SetupRoutes(hashMap *uhash.Hashtable, t *Tbot.Telegramini) *gin.Engine {
	//gin.New()
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(LoggerWithFormatter))
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("gin日志 %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	//gin.SetMode(gin.ReleaseMode)
	r.POST("/webhook/:id", PostWebHook(hashMap, t))
	r.GET("/webhook/:id", GetWebHook(hashMap, t))
	return r
}
func LoggerWithFormatter(params gin.LogFormatterParams) string {
	return fmt.Sprintf(
		" %s  | %d | \t %s | %s | %s \t  %s\n",
		params.TimeStamp.Format("2006/01/02 - 15:04:05"),
		params.StatusCode, // 状态码
		params.ClientIP,   // 客户端ip
		params.Latency,    // 请求耗时
		params.Method,     // 请求方法
		params.Path,       // 路径
	)
}
func PostWebHook(hashMap *uhash.Hashtable, t *Tbot.Telegramini) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		//通过id判断hashMap中是否存在该用户
		//main包中的HashMap如何传递到这里
		idInt, err := hashMap.Get(id)
		if err != nil {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
			return
		}
		rawData, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		data := string(rawData)
		//输出hashMap中的用户
		fmt.Println("用户", idInt)
		fmt.Println("消息内容", data)

		user := &tele.User{
			ID: idInt,
		}

		_, err = t.Bot.Send(user, data, &tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})

		if err != nil {
			c.JSON(200, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "ok",
		})

	}
}

func GetWebHook(hashMap *uhash.Hashtable, t *Tbot.Telegramini) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("id", id)
		//通过id判断hashMap中是否存在该用户
		//main包中的HashMap如何传递到这里

		idInt, err := hashMap.Get(id)
		if err != nil {
			fmt.Println(err)
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
			return
		}
		rawData := c.Query("text")
		fmt.Println("消息内容", rawData)
		data := string(rawData)
		//输出hashMap中的用户
		fmt.Println("用户", idInt)
		//msg := c.PostForm("message")
		fmt.Println("消息内容", rawData)
		user := &tele.User{
			ID: idInt,
		}
		_, err = t.Bot.Send(user, data, &tele.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tele.ModeMarkdown,
		})

		if err != nil {
			c.JSON(404, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "ok",
		})

	}
}
