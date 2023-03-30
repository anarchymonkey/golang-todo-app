# Todo Application

This is a simple todo application made with Golang in the backend, PostgreSQL as the database, and React as a CSR rendered app.

## Installation

To run the application, you need to have Golang, PostgreSQL, and Node.js installed on your machine.

* Clone the repository: git clone <https://github.com/yourusername/todo-app.git>
* Navigate to the project directory: cd todo-app
* Create a new PostgreSQL database named todo
* Import the database.sql file located in the server directory to create the required tables in the database
* Navigate to the server directory and run go run main.go to start the server
* Navigate to the client directory and run npm install to install the dependencies
* Run `npm run dev` to start the client-side app
* Access the application at <http://localhost:1234>

Folder Structure
The application has the following folder structure:

```css
todo-app/
  ├── client/
  │   ├── public/
  │   ├── src/
  │   ├── package.json
  │   └── ...
  ├── server/
  │   ├── main.go
  │   ├── database.sql
  │   └── ...
  └── README.md
```

* `client/` - the directory containing the client-side app code
* `client/public/` - the directory containing the public files for the client-side app
* `client/src/` - the directory containing the source files for the client-side app
* `client/package.json` - the package.json file for the client-side app, which is using Parcel 2.0
* `server/` - the directory containing the server-side app code
* `server/main.go` - the main Go file for the server-side app
* `server/database.sql` - the SQL file for creating the required tables in the database (tbd)

## Technology Stack

The following technologies were used to build this application:

* `Golang` - for the backend server
* `PostgreSQL` - for the database
* `React` - for the client-side rendered user interface
* `Parcel 2.0`- for the client-side app bundling

## Contributing

If you would like to contribute to the development of this application, please follow these steps:

## Fork the repository

* Create a new branch for your changes: git checkout -b feature/your-feature-name
* Make your changes and commit them: git commit -am 'Add your commit message here'
* Push your changes to your fork: git push origin feature/your-feature-name
* Create a pull request to merge your changes into the main branch of the repository

## License

This project is licensed under the MIT License
