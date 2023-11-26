package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"telegram-notice/hash"
	"telegram-notice/tgbot"
)

func SetupRoutes(r *gin.Engine, hashMap *uhash.Hashtable, t tgbot.TgBot) {
	r.POST("/webhook/:id", PostWebHook(hashMap, t))
	r.GET("/webhook/:id", GetWebHook(hashMap, t))
}
func PostWebHook(hashMap *uhash.Hashtable, t tgbot.TgBot) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		//通过id判断hashMap中是否存在该用户
		//main包中的HashMap如何传递到这里
		if _, ok := hashMap.Get(id); !ok {
			c.JSON(404, gin.H{
				"message": "error",
			})
			return
		}
		rawData, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		data := string(rawData)
		idInt, _ := hashMap.Get(id)
		//输出hashMap中的用户
		fmt.Println("用户", idInt)
		msg := c.PostForm("message")
		fmt.Println("消息内容", msg)

		t.SendMeesg(idInt, data)

		c.JSON(200, gin.H{
			"message": "ok",
		})

	}
}

func GetWebHook(hashMap *uhash.Hashtable, t tgbot.TgBot) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println("id", id)
		//通过id判断hashMap中是否存在该用户
		//main包中的HashMap如何传递到这里
		if _, ok := hashMap.Get(id); !ok {
			c.JSON(404, gin.H{
				"message": "error",
			})
			return
		}
		rawData := c.Query("text")
		fmt.Println("消息内容", rawData)
		data := string(rawData)
		idInt, _ := hashMap.Get(id)
		//输出hashMap中的用户
		fmt.Println("用户", idInt)
		//msg := c.PostForm("message")
		fmt.Println("消息内容", rawData)
		err := hashMap.SaveToFile("hash.json")
		if err != nil {
			t.SendMeesg(idInt, "数据库报错失败")
			return
		}
		t.SendMeesg(idInt, data)

		c.JSON(200, gin.H{
			"message": "ok",
		})

	}
}
