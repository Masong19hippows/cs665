package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Dish struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Instructions   string `json:"instructions"`
	Ingredients    string `json:"ingredients"`
	SimilarRecipes string `json:"similar_recipes"`
}

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	fmt.Println("Opening DB file:", path)
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		fmt.Println("Error opening DB:", err)
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging DB:", err)
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (db *Database) GetAllDishes() ([]Dish, error) {
	fmt.Println("Querying all dishes")
	rows, err := db.conn.Query(`SELECT "ID", "Title", "Instructions", "Ingredients", "SimilarRecipes" FROM Dishes`)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var dishes []Dish
	for rows.Next() {
		var d Dish
		var simRecipes sql.NullString

		err := rows.Scan(&d.ID, &d.Title, &d.Instructions, &d.Ingredients, &simRecipes)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		if simRecipes.Valid {
			d.SimilarRecipes = simRecipes.String
		} else {
			d.SimilarRecipes = ""
		}

		dishes = append(dishes, d)
	}
	return dishes, nil
}

func (db *Database) AddDish(d Dish) error {
	fmt.Println("Inserting dish:", d)
	_, err := db.conn.Exec(`INSERT INTO Dishes ("Title", "Instructions", "Ingredients", "SimilarRecipes") VALUES (?, ?, ?, ?)`,
		d.Title, d.Instructions, d.Ingredients, nullString(d.SimilarRecipes))
	if err != nil {
		fmt.Println("Error inserting dish:", err)
	}
	return err
}

func (db *Database) UpdateDish(d Dish) error {
	fmt.Println("Updating dish:", d)
	_, err := db.conn.Exec(`UPDATE Dishes SET "Title" = ?, "Instructions" = ?, "Ingredients" = ?, "SimilarRecipes" = ? WHERE "ID" = ?`,
		d.Title, d.Instructions, d.Ingredients, nullString(d.SimilarRecipes), d.ID)
	if err != nil {
		fmt.Println("Error updating dish:", err)
	}
	return err
}

func (db *Database) DeleteDish(id int) error {
	fmt.Println("Deleting dish ID:", id)
	_, err := db.conn.Exec(`DELETE FROM Dishes WHERE "ID" = ?`, id)
	if err != nil {
		fmt.Println("Error deleting dish:", err)
	}
	return err
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
