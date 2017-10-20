package main

import (
	"github.com/night-codes/tokay"
	"gopkg.in/mgo.v2"
)

type obj map[string]interface{}

type feedback struct {
	Name    string `form:"name" bson:"name" valid:"required,min(3),max(40)"`
	Title   string `form:"title" bson:"title" valid:"required,max(150)"`
	Message string `form:"message" bson:"message" valid:"required"`
}

func main() {
	session, err := mgo.Dial(":27017") //mongodb connect
	if err != nil {
		panic(err)
	}

	r := tokay.New()
	r.Static("/files", "./files")

	r.GET("/", func(c *tokay.Context) {
		c.HTML(200, "index", obj{"title": "My website", "name": `My Friend`})
	})

	r.POST("/", func(c *tokay.Context) {
		fb := feedback{}
		ret := obj{"title": "Website", "name": `taliban`}
		if err := c.Bind(&fb); err != nil {
			ret["err"] = "Oops, an error: " + err.Error()
		} else {
			if err := session.DB("mydb").C("feedbacks").Insert(fb); err != nil {
				ret["err"] = "Unexpected error. Come back to us later."
			} else {
				ret["ok"] = "Thanks for your feedback!"
			}
		}
		c.HTML(200, "index", ret)
	})

	admin := r.Group("/admin", tokay.BasicAuth("admin", "secret"))
	admin.GET("/", func(c *tokay.Context) {
		fbks := []feedback{}
		session.DB("mydb").C("feedbacks").Find(obj{}).All(&fbks)
		c.HTML(200, "admin", obj{"feedbacks": fbks})
	})

	panic(r.Run(":8080"))
}
