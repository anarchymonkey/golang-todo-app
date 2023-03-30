# Golang Todo Web Application
This is a web application built with Golang for managing todo items. It uses the Gin web framework and a PostgreSQL database.

## Prerequisites
- [Golang](https://golang.org/) version 1.16 or later
- [PostgreSQL](https://www.postgresql.org/) version 10 or later

## Installation
- Clone the repository: `git clone https://github.com/anarchymonkey/golang-todo-server.git`
- Set environment variables for PostgreSQL connection:
```arduino
export PG_USERNAME_V1=<your postgres username>
export PG_PASS_V1=<your postgres password>
```
- Install dependencies: `go mod download or go get .`
- Build the application: `go build -o app`
- Run the application: `go run .`
  
## Usage

Once the application is running, you can interact with it using HTTP requests. The following endpoints are available:

### GET groups
Returns a list of all groups.

### GET /group/:id/items
Returns a list of all items in the group with the given ID.

### GET /item/:id/contents
Returns a list of all contents in the item with the given ID.

### POST /group/add
Adds a new group.

Example request body:

```json
{
    "name": "My Group"
}
```

### POST /group/:id/item/add
Adds a new item to the group with the given ID.

Example request body:
```json
{
    "name": "My Item",
    "description": "This is a description of my item."
}
```

### POST /item/:id/content/add
Adds a new content to the item with the given ID.

Example request body:
```json
{
    "text": "This is my content."
}
```

### PUT /group/:id/update
Updates the group with the given ID.

Example request body:
```json
{
    "name": "New Group Name"
}
```

### PUT /item/:id/update
Updates the item with the given ID.

Example request body:
```json
{
    "name": "New Item Name",
    "description": "New description of my item."
}
```
### PUT /item/:id/content/update
Updates the content with the given ID.

Example request body:
```json
{
    "text": "Updated content."
}
```

### DELETE /group/:id/delete
Deletes the group with the given ID.

### DELETE /group/:id/item/:item_id/delete
Deletes the item with the given ID from the group with the given group ID.

### DELETE /item/:item_id/content/:content_id/delete
Deletes the content with the given ID from the item with the given item ID.

## CORS
The application includes a middleware function for handling CORS. It allows requests from any origin and allows the following methods: POST, GET, OPTIONS, PUT, DELETE. It also allows the following headers: Content-Type, Authorization, mode.

## Error Handling
The application includes basic error handling for database errors and bad requests. If an error occurs, the server will return an appropriate HTTP status code along with a JSON error message.

## License
This project is licensed under the [MIT License](https://opensource.org/license/mit/).