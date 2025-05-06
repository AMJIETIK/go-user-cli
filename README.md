# **Go User CLI**
This is a simple CLI application for managing users in a PostgreSQL database.
# Setup
1. Clone the repository:
```bash
git clone https://github.com/AMJIETIK/go-user-cli.git
```

2. Navigate to the project directory:
```bash
cd go-user-cli
```

3. Create a .env file (named password.env) in the project root directory with your database connection string:
```bash
DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
```

4. Install dependencies:
```go
go mod tidy
```

5. Run the application:
```go
go run main.go
```

# **Database Table**
This application works with a users table in PostgreSQL with the following schema:
```sql
CREATE TABLE users (
id SERIAL PRIMARY KEY,
name VARCHAR(100) NOT NULL,
email VARCHAR(100) NOT NULL UNIQUE,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Commands
- show - Display all users.
- add - Add a new user.
- edit - Edit an existing user.
- delete - Delete a user.
- find - Find a user by name.