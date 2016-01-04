package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type feedback struct {
	Name    string `form:"name" binding:"required,min=3,max=40"`
	Title   string `form:"title" binding:"required,max=150"`
	Message string `form:"message" binding:"required"`
}

func main() {
	session, err := mgo.Dial(":27017") //mongodb connect
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Static("files", "./files")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Сайт-визитка", "name": `taliban`})
	})

	r.POST("/", func(c *gin.Context) {
		fb := feedback{}
		ret := gin.H{"title": "Сайт-визитка", "name": `taliban`}
		if err := c.Bind(&fb); err != nil {
			ret["err"] = "Упс, ошибка: " + err.Error()
		} else {
			if err := session.DB("mydb").C("feedbacks").Insert(fb); err != nil {
				ret["err"] = "Неожиданная ошибка. Зайдите к нам попозже."
			} else {
				ret["ok"] = "Спасибо за ваш отзыв!"
			}
		}
		c.HTML(200, "index.html", ret)
	})

	r.Run(":8080")
}
