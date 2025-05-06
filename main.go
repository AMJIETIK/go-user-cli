package main

import (
	"bufio"
	"fmt"
	"os"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgconn"

)

func main() {

	//connect to database
	dbpool, ctx := connectDB()
	defer dbpool.Close()

	scanner := bufio.NewScanner(os.Stdin)
	var err error
	for {
		fmt.Println("Choose an option to make: \n(show, add, edit, delete, find, refreshID)")
		scanner.Scan()
		input := scanner.Text()
	
		switch input {
		case "show":
			err = show(dbpool, ctx)
			Check(err)
		case "add":
			add(dbpool, ctx, scanner)
		case "edit":
			edit(dbpool, ctx, scanner)
		case "delete":
			delete(dbpool, ctx, scanner)
		case "find":
			find(dbpool, ctx, scanner)
		default:
			fmt.Println("It seems like u missed the letter. Try again \n(show, add, edit, delete, find)")
		}
	}	

}

func show(dbpool *pgxpool.Pool, ctx context.Context) error{
	var name, email string
	var id int
	var date_registered time.Time
	rows, err := dbpool.Query(ctx, "SELECT * FROM users ORDER BY id")
	Check(err)
	defer rows.Close()

	fmt.Printf("\n\n%-4s | %-15s | %-25s | %-26s\n", "id", "name", "email", "date_registered")
	fmt.Println("-----+-----------------+---------------------------+--------------------------------------")
	for rows.Next(){
		err = rows.Scan(&id,&name,&email,&date_registered)
		Check(err)
		fmt.Printf("%-4d | %-15s | %-25s | %-26s\n", id, name, email, date_registered)
	}
	fmt.Print("\n\n")
	return nil
}

func add(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner){
	var name, email string
	fmt.Println("Full name of new user:")
	scanner.Scan()
	name = scanner.Text()
	fmt.Println("New user's email:")
	scanner.Scan()
	email = scanner.Text()

	_, err := dbpool.Exec(ctx, "INSERT INTO users(name,email) VALUES($1,$2)", name, email)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		fmt.Println("Error: user with this email already exists.")
		return
	}
	Check(err)
	fmt.Print("\nExecution was successfully completed\n\n")
}

func edit(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner){
	var id, content, fieldName string
	fmt.Print("\nEnter the ID of the user you want to update: ")
	scanner.Scan()
	id = scanner.Text()

	fmt.Print("\nWhich field do you want to update? (name or email): ")
	scanner.Scan()
	fieldName = scanner.Text()
	validateField(fieldName)

	fmt.Printf("\nEnter the new value for %s: ", fieldName)
	scanner.Scan()
	content = scanner.Text()

	query := fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2", fieldName)
	_, err := dbpool.Exec(ctx, query, content, id)
	Check(err)
	fmt.Printf("\n%s was succesfully edited!!!\n\n", fieldName)
}

func delete(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner){
	fmt.Println("Enter id of row that you want to delete:")
	scanner.Scan()
	id := scanner.Text()

	_, err := dbpool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	Check(err)
	refreshID(dbpool, ctx)
	fmt.Println("Row was succesfully deleted", "", "")
}

func find(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner){
	fmt.Print("\nWhich field would you like to use to find the client? (name or email): ")
	scanner.Scan()
	fieldName := scanner.Text()
	validateField(fieldName)

	fmt.Printf("\nEnter the %s you want to search for: ", fieldName)
	scanner.Scan()
	content := scanner.Text()

	//db part
	var id, name, email string
	var time time.Time
	query := fmt.Sprintf("SELECT * FROM users WHERE %s ILIKE $1", fieldName)
	rows, err := dbpool.Query(ctx, query, "%"+content+"%")
	Check(err)
	defer rows.Close()

	//print
	fmt.Printf("\n\n%-4s | %-15s | %-25s | %-26s\n", "id", "name", "email", "date_registered")
	fmt.Println("-----+-----------------+---------------------------+--------------------------------------")

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &time)
		Check(err)
		fmt.Printf("%-4s | %-15s | %-25s | %-26s\n", id, name, email, time)
	}
	fmt.Print("\n\n")
}

func refreshID(dbpool *pgxpool.Pool, ctx context.Context) {
	_, err := dbpool.Exec(ctx, "SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));")
	Check(err)
}

func connectDB() (*pgxpool.Pool, context.Context) {
	ctx := context.Background()
	connStr := "postgres://stormside7:MrVladPro2008@localhost:5432/stormside7"

	dbpool, err := pgxpool.New(ctx, connStr)
	Check(err)

	return dbpool, ctx
}

func Check(err error){
	if err != nil{
		log.Fatal(err)
	}
}
func validateField(whereEdit string) {
	validFields := map[string]bool{
		"name":  true,
		"email": true,
	}
	if !validFields[whereEdit]{
		log.Fatalf("Недопустимое поле для изменения: %s", whereEdit)
	}
}