package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
	"io"
	"log"
	"net/http"
)

type Cat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var db *sql.DB

// Функция для создания таблицы
func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./cats.db")
	if err != nil {
		return err
	}

	// Создаем таблицу для кошек
	createTableSQL := `CREATE TABLE IF NOT EXISTS cats (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"type" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "this is the Web side")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")

	// Проверка обязательных параметров
	if catName == "" || catType == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Both 'name' and 'type' parameters are required.",
		})
	}

	if dataType != "string" && dataType != "json" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid 'data' parameter. It must be either 'string' or 'json'.",
		})
	}

	// Формируем строковый или JSON ответ
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Your cat name is: %s\nand his type is: %s\n", catName, catType))
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Unexpected error. Check the 'data' parameter.",
	})
}

// Функция добавления новой кошки в базу данных
func addCat(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()

	// Читаем данные из тела запроса
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed to read body request: %s", err)
		return c.String(http.StatusInternalServerError, "Failed to read request body.")
	}

	// Разбираем JSON в структуру Cat
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed to unmarshal in addCat: %s", err)
		return c.String(http.StatusInternalServerError, "Failed to parse JSON data.")
	}

	// Вставляем данные в таблицу
	stmt, err := db.Prepare("INSERT INTO cats(name, type) VALUES(?, ?)")
	if err != nil {
		log.Printf("Failed to prepare SQL statement: %s", err)
		return c.String(http.StatusInternalServerError, "Failed to prepare SQL statement.")
	}
	defer stmt.Close()

	_, err = stmt.Exec(cat.Name, cat.Type)
	if err != nil {
		log.Printf("Failed to insert into database: %s", err)
		return c.String(http.StatusInternalServerError, "Failed to insert data into the database.")
	}

	log.Printf("This is your cat: %#v", cat)

	// Ответ после успешного добавления кошки
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Cat added successfully",
		"name":    cat.Name,
		"type":    cat.Type,
	})
}

// Функция для получения всех кошек из базы данных
func getAllCats(c echo.Context) error {
	rows, err := db.Query("SELECT id, name, type FROM cats")
	if err != nil {
		log.Printf("Failed to query database: %s", err)
		return c.String(http.StatusInternalServerError, "Failed to retrieve data from the database.")
	}
	defer rows.Close()

	var cats []Cat
	for rows.Next() {
		var cat Cat
		err = rows.Scan(&cat.ID, &cat.Name, &cat.Type)
		if err != nil {
			log.Printf("Failed to scan row: %s", err)
			return c.String(http.StatusInternalServerError, "Failed to scan row.")
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %s", err)
		return c.String(http.StatusInternalServerError, "Rows error.")
	}

	return c.JSON(http.StatusOK, cats)
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	fmt.Println("Welcome to the server")
	e := echo.New()

	// Роуты
	e.GET("/", hello)
	e.GET("/cats/:data", getCats)
	e.GET("/cats", getAllCats)
	e.POST("/cats", addCat)

	// Старт сервера
	e.Start(":8000")
}
