package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	// connect to database (Подключение к базе данных)
	err := godotenv.Load("password.env")
	if err != nil {
		log.Fatal("Error loading .env file") // Error loading .env file (Ошибка при загрузке .env файла)
	}

	// Using the environment variable DATABASE_URL (Использование переменной окружения DATABASE_URL)
	connStr := os.Getenv("DATABASE_URL")
	dbpool, ctx := connectDB(connStr)
	defer dbpool.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Choose an option to make: \n(show, add, edit, delete, find)") // Choose an option to make (Выберите действие)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "show":
			err = show(dbpool, ctx) // Show users (Показать пользователей)
			Check(err)
		case "add":
			add(dbpool, ctx, scanner) // Add new user (Добавить нового пользователя)
		case "edit":
			edit(dbpool, ctx, scanner) // Edit user (Редактировать пользователя)
		case "delete":
			delete(dbpool, ctx, scanner) // Delete user (Удалить пользователя)
		case "find":
			find(dbpool, ctx, scanner) // Find user (Найти пользователя)
		default:
			fmt.Println("It seems like u missed the letter. Try again \n(show, add, edit, delete, find)") // Invalid option (Недопустимый выбор)
		}
	}

}

func show(dbpool *pgxpool.Pool, ctx context.Context) error {
	var name, email string
	var id int
	var date_registered time.Time
	rows, err := dbpool.Query(ctx, "SELECT * FROM users ORDER BY id")
	Check(err)
	defer rows.Close()

	// Print headers (Вывод заголовков)
	fmt.Printf("\n\n%-4s | %-15s | %-25s | %-26s\n", "id", "name", "email", "date_registered")
	fmt.Println("-----+-----------------+---------------------------+--------------------------------------")
	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &date_registered)
		Check(err)
		// Print user data (Вывод данных пользователей)
		fmt.Printf("%-4d | %-15s | %-25s | %-26s\n", id, name, email, date_registered)
	}
	fmt.Print("\n\n")
	return nil
}

func add(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner) {
	var name, email string
	fmt.Println("Full name of new user:") // Full name of new user (Полное имя нового пользователя)
	scanner.Scan()
	name = scanner.Text()
	fmt.Println("New user's email:") // New user's email (Email нового пользователя)
	scanner.Scan()
	email = scanner.Text()

	_, err := dbpool.Exec(ctx, "INSERT INTO users(name,email) VALUES($1,$2)", name, email)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		fmt.Println("Error: user with this email already exists.") // Error: user with this email already exists (Ошибка: пользователь с этим email уже существует)
		return
	}
	Check(err)
	fmt.Print("\nExecution was successfully completed\n\n") // Execution was successfully completed (Выполнение завершено успешно)
}

func edit(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner) {
	var id, content, fieldName string
	fmt.Print("\nEnter the ID of the user you want to update: ") // Enter the ID of the user you want to update (Введите ID пользователя, которого хотите обновить)
	scanner.Scan()
	id = scanner.Text()

	fmt.Print("\nWhich field do you want to update? (name or email): ") // Which field do you want to update? (name or email) (Какое поле вы хотите обновить? (имя или email))
	scanner.Scan()
	fieldName = scanner.Text()
	validateField(fieldName)

	fmt.Printf("\nEnter the new value for %s: ", fieldName) // Enter the new value for %s (Введите новое значение для %s)
	scanner.Scan()
	content = scanner.Text()

	query := fmt.Sprintf("UPDATE users SET %s = $1 WHERE id = $2", fieldName)
	_, err := dbpool.Exec(ctx, query, content, id)
	Check(err)
	fmt.Printf("\n%s was succesfully edited!!!\n\n", fieldName) // %s was successfully edited!!! (%s было успешно отредактировано!!!)
}

func delete(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner) {
	fmt.Println("Enter id of row that you want to delete:") // Enter id of row that you want to delete (Введите ID строки, которую хотите удалить)
	scanner.Scan()
	id := scanner.Text()

	_, err := dbpool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	Check(err)
	refreshID(dbpool, ctx)
	fmt.Println("Row was succesfully deleted", "", "") // Row was successfully deleted (Строка была успешно удалена)
}

func find(dbpool *pgxpool.Pool, ctx context.Context, scanner *bufio.Scanner) {
	fmt.Print("\nWhich field would you like to use to find the client? (name or email): ") // Which field would you like to use to find the client? (name or email) (Какое поле вы хотите использовать для поиска клиента? (имя или email))
	scanner.Scan()
	fieldName := scanner.Text()
	validateField(fieldName)

	fmt.Printf("\nEnter the %s you want to search for: ", fieldName) // Enter the %s you want to search for (Введите %s, который вы хотите найти)
	scanner.Scan()
	content := scanner.Text()

	// db part (Часть запроса к базе данных)
	var id, name, email string
	var time time.Time
	query := fmt.Sprintf("SELECT * FROM users WHERE %s ILIKE $1", fieldName)
	rows, err := dbpool.Query(ctx, query, "%"+content+"%")
	Check(err)
	defer rows.Close()

	// print (Вывод результатов)
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
	_, err := dbpool.Exec(ctx, "SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));") // Refresh ID sequence (Обновление последовательности ID)
	Check(err)
}

func connectDB(connStr string) (*pgxpool.Pool, context.Context) {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connStr)
	Check(err)

	return dbpool, ctx
}

func Check(err error) {
	if err != nil {
		log.Fatal(err) // Handle error (Обработка ошибки)
	}
}

func validateField(whereEdit string) {
	validFields := map[string]bool{
		"name":  true,
		"email": true,
	}
	if !validFields[whereEdit] {
		log.Fatalf("Недопустимое поле для изменения: %s", whereEdit) // Invalid field for update: %s (Недопустимое поле для изменения: %s)
	}
}
