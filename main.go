package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Connection to db failed")
	}

	type Person struct {
		gorm.Model
		ID         uint `gorm:"primaryKey"`
		Name       string
		Age        int
		Occupation string
	}

	db.AutoMigrate(&Person{})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.tmpl", gin.H{
			"title": "Create person",
		})
	})

	r.POST("/create", func(c *gin.Context) {
		var newPerson Person

		if err := c.ShouldBind(&newPerson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Binding error",
				"details": err.Error(),
			})
			return
		}

		result := db.Create(&newPerson)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, newPerson)
	})

	r.GET("/view", func(c *gin.Context) {
		var people []Person

		result := db.Find(&people)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.HTML(http.StatusOK, "view.tmpl", gin.H{
			"title":  "View",
			"people": people,
		})
	})

	r.Run()
}
