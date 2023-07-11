package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

//NOTE: https://forum.golangbridge.org/t/database-rows-scan-unknown-number-of-columns-json/7378/2
//above is the reference for code sample.

func main() {
	urlExample := "postgres://user:password@localhost:5432/DelDB?sslmode=disable"
	db, err := sql.Open("postgres", urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Fiber setup and routes starting
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		res, err := queryScanJsonErrorLast(db)
		if err != nil {
			log.Println("data error:", err)
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "something went wrong!",
			})
		}

		log.Printf("res: %#v \n", string(res))
		return c.SendString(string(res))
	})

	// this handler should encounter error and throw it in result
	app.Get("/error", func(c *fiber.Ctx) error {
		res, err := queryScanJsonError(db)
		if err != nil {
			log.Println("fn:queryScanJsonError data error:", err)
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "something went wrong!",
			})
		}

		return c.SendString(string(res))
	})

	log.Fatal(app.Listen(":3000"))
}

type Warrior struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Teacher   string `json:"teacher"`
	IsActive  bool   `json:"isActive"`
}

func queryScanJsonError(db *sql.DB) ([]byte, error) {
	res := []Warrior{}
	q := "select first_name,last_name, teacher, is_active from warriors"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		w := Warrior{}
		// err := rows.Scan(&w.FirstName, &w.LastName, &w.Teacher, &w.IsActive)
		err := rows.Scan(&w.FirstName, &w.LastName, &w.Teacher, &w.IsActive)
		if err != nil {
			fmt.Printf("scan error but with warrior %#v \n", w)
			return nil, err
		}
		res = append(res, w)
	}
	return json.MarshalIndent(res, "", "\t")
}

func queryScanJsonErrorLast(db *sql.DB) ([]byte, error) {
	res := []Warrior{}
	q := "select first_name, teacher, is_active, last_name from warriors"

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		w := Warrior{}
		// err := rows.Scan(&w.FirstName, &w.LastName, &w.Teacher, &w.IsActive)
		err := rows.Scan(&w.FirstName, &w.Teacher, &w.IsActive, &w.LastName)
		if err != nil {
			fmt.Printf("scan error but with warrior %#v \n", w)
			return nil, err
		}
		res = append(res, w)
	}
	return json.MarshalIndent(res, "", "\t")
}
