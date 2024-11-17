package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

type Person struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Nick string `json:"nick"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("people.db"), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	db.AutoMigrate(&Person{})

	var count int64
	db.Model(&Person{}).Count(&count)
	if count == 0 {
		examplePeople := []Person{
			{Name: "Иван", Nick: "ivan_the_great"},
			{Name: "Анна", Nick: "anna_star"},
			{Name: "Петр", Nick: "petr_petrov"},
		}
		db.Create(&examplePeople)
		fmt.Println("Пример данных добавлен в базу.")
	}

	e := echo.New()

	e.GET("/people", func(c echo.Context) error {
		var people []Person
		db.Find(&people)
		return c.JSON(http.StatusOK, people)
	})

	e.POST("/people", func(c echo.Context) error {
		person := new(Person)
		if err := c.Bind(person); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Некорректные данные",
			})
		}
		db.Create(&person)
		return c.JSON(http.StatusCreated, person)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
