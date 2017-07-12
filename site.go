package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type feedback struct {
	Name    string `form:"name" bson:"name" binding:"required,min=3,max=40"`
	Title   string `form:"title" bson:"title" binding:"required,max=150"`
	Message string `form:"message" bson:"message" binding:"required"`
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
		c.HTML(200, "index.html", gin.H{"title": "My website", "name": `taliban`})
	})

	r.POST("/", func(c *gin.Context) {
		fb := feedback{}
		ret := gin.H{"title": "Website", "name": `taliban`}
		if err := c.Bind(&fb); err != nil {
			ret["err"] = "Oops, an error: " + err.Error()
		} else {
			if err := session.DB("mydb").C("feedbacks").Insert(fb); err != nil {
				ret["err"] = "Unexpected error. Come back to us later."
			} else {
				ret["ok"] = "Thanks for your feedback!"
			}
		}
		c.HTML(200, "index.html", ret)
	})

	admin := r.Group("/admin")
	admin.Use(gin.BasicAuth(map[string]string{"admin": "secret"}))
	admin.GET("/", func(c *gin.Context) {
		fbks := []feedback{}
		session.DB("mydb").C("feedbacks").Find(gin.H{}).All(&fbks)
		c.HTML(200, "admin.html", gin.H{"feedbacks": fbks})
	})

	r.Run(":8080")
}
