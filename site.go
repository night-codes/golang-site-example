package main

import (
	"github.com/night-codes/tokay"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	obj      map[string]interface{}
	feedback struct {
		ID      uint64 `form:"id" gorm:"primary_key"`
		Name    string `form:"name" valid:"required,min(3),max(40)"`
		Title   string `form:"title" valid:"required,max(150)"`
		Message string `form:"message" valid:"required"`
	}
)

func main() {
	db, err := gorm.Open(sqlite.Open("site.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&feedback{})

	r := tokay.New()
	r.Static("/files", "./files")

	r.GET("/", func(c *tokay.Context) {
		c.HTML(200, "index", obj{"title": "My website", "name": `My Friend`})
	})

	r.POST("/", func(c *tokay.Context) {
		fb := feedback{}
		ret := obj{"title": "My website", "name": `My Friend`}
		if err := c.Bind(&fb); err != nil {
			ret["err"] = "Oops, an error: " + err.Error()
		} else {
			if err := db.Save(&fb).Error; err != nil {
				ret["err"] = "Unexpected error. Come back to us later."
			} else {
				ret["ok"] = "Thanks for your feedback!"
			}
		}

		c.HTML(200, "index", ret)
	})

	admin := r.Group("/admin", tokay.BasicAuth("admin", "secret"))
	admin.GET("/", func(c *tokay.Context) {
		feedbacks := []feedback{}
		db.Find(&feedbacks)
		c.HTML(200, "admin", obj{"feedbacks": feedbacks})
	})

	panic(r.Run(":8080"))
}
