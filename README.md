# Book Management API

A RESTful API for managing a collection of books built with Go, the Fiber web framework, and PostgreSQL.

## Description

This project provides a RESTful API to manage books in a PostgreSQL database. It supports CRUD (Create, Read, Update, Delete) operations for book records. The API allows you to add new books, retrieve all books or a specific book by its ID, and delete books from the database.

## Features

- **Create Book**: Add a new book to the database with details such as author, title, and publisher.
- **Get All Books**: Retrieve a list of all books stored in the database.
- **Get Book by ID**: Fetch details of a specific book using its unique ID.
- **Delete Book**: Remove a book from the database using its unique ID.

## Endpoints

### Create Book

- **POST /api/create_books**
    - **Description**: Adds a new book to the database.
    - **Request Body**:
      ```json
      {
        "author": "Author Name",
        "title": "Book Title",
        "publisher": "Publisher Name"
      }
      ```
    - **Response**:
      ```json
      {
        "message": "Book created successfully",
        "data": {
          "id": 1,
          "author": "Author Name",
          "title": "Book Title",
          "publisher": "Publisher Name"
        }
      }
      ```

### Get All Books

- **GET /api/books**
    - **Description**: Retrieves a list of all books from the database.
    - **Response**:
      ```json
      {
        "message": "Books fetched successfully",
        "data": [
          {
            "id": 1,
            "author": "Author Name",
            "title": "Book Title",
            "publisher": "Publisher Name"
          }
        ]
      }
      ```

### Get Book by ID

- **GET /api/get_books/:id**
    - **Description**: Retrieves details of a book by its ID.
    - **Response**:
      ```json
      {
        "message": "Book fetched successfully",
        "data": {
          "id": 1,
          "author": "Author Name",
          "title": "Book Title",
          "publisher": "Publisher Name"
        }
      }
      ```

### Delete Book

- **DELETE /api/delete_book/:id**
    - **Description**: Deletes a book from the database using its ID.
    - **Response**:
      ```json
      {
        "message": "Book deleted successfully"
      }
      ```

## Setup and Installation

### Prerequisites

- Go (1.18+)
- PostgreSQL
- [Go-Fiber](https://github.com/gofiber/fiber) framework
- [GORM](https://gorm.io/) ORM

### Installation

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/yourusername/go-fiber-postgres.git
    cd go-fiber-postgres
    ```

2. **Install Dependencies**:

    ```bash
    go mod tidy
    ```

3. **Setup Environment Variables**

   Create a `.env` file in the root directory with the following content:

    ```plaintext
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_SSLMODE=disable
    ```

4. **Run Database Migrations**:

   Ensure PostgreSQL is running and accessible, then start the application to run the migrations:

    ```bash
    go run main.go
    ```

## Usage

- **Start the Server**:

    ```bash
    go run main.go
    ```

  The server will start on `http://localhost:8080`.

- **Access the API**: Use tools like [Postman](https://www.postman.com/) or `curl` to interact with the API endpoints.

## Benchmarking

To benchmark your PostgreSQL database, you can use tools like `pgBench` or `SysBench`.

### Using `pgBench`

1. **Initialize the Database**:

    ```bash
    pgbench -i -s 10 your_database
    ```

2. **Run the Benchmark**:

    ```bash
    pgbench -c 10 -j 2 -T 60 your_database
    ```

### Using `SysBench`

1. **Prepare the Database**:

    ```bash
    sysbench --db-driver=pgsql --pgsql-host=localhost --pgsql-user=postgres --pgsql-password=your_password --pgsql-db=your_database oltp_read_write prepare
    ```

2. **Run the Benchmark**:

    ```bash
    sysbench --db-driver=pgsql --pgsql-host=localhost --pgsql-user=postgres --pgsql-password=your_password --pgsql-db=your_database oltp_read_write --threads=10 --time=60 run
    ```

## Contributing

Feel free to fork the repository and submit pull requests. Contributions are welcome!

## Contact

For questions or feedback, please open an issue on GitHub or contact [iamdenisetiawan@gmail.com](mailto:iamdenisetiawan@gmail.com).
