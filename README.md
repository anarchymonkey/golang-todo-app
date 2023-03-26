
# Golang Todo Server

This is a Golang web application that provides an API for a Todo list. It uses Gin as the web framework and pgxpool for database connections.

## Prerequisites
  Before running this application, ensure that the following software is installed on your machine:

* Golang
* PostgreSQL
* Getting Started
* Clone this repository to your local machine:

```sh
git clone https://github.com/anarchymonkey/golang-todo-server.git
```
Set the following environment variables with your PostgreSQL credentials:

```sh
export PG_USERNAME_V1=your_username
export PG_PASS_V1=your_password
```

Create the required database schema by running the following command in the root directory of the project:

```sh
go run schema/main.go
```

Start the server by running the following command:

```sh
go run main.go
```
The server will start listening on port 8080.

## API Endpoints
- `GET /todos`

Returns a list of all todos.

- `POST /todo/add`

Adds a new todo to the list. Expects a JSON request body with the following fields:
```sql
* title (string)
* description (string)
* due_date (date)
```

- `PUT /todo/update`

Updates an existing todo in the list. Expects a JSON request body with the following fields:

```sql
* id (int)
* title (string)
* description (string)
* due_date (date)
```

## Middlewares

**CORS Middleware**

This middleware adds the necessary headers to enable Cross-Origin Resource Sharing (CORS) requests. It allows requests from any origin and supports the following HTTP methods: POST, GET, OPTIONS, PUT, DELETE. It also allows the following headers: Content-Type, Authorization, and mode.