# RailwayConsole

## Overview

RailwayConsole is a comprehensive Go-based application designed for managing railway operations, focusing on stations and trains. It provides functionalities for creating, reading, updating, and deleting (CRUD) records of stations and trains, alongside advanced features like transaction handling and isolation level demonstrations.

## Features

- **CRUD Operations**: Manage stations and trains with create, read, update, and delete functionalities.
- **Database Interaction**: Integrated MySQL database operations.
- **Transaction Handling**: Demonstrates transactional operations and error handling.
- **Concurrency Control**: Simulate various database isolation levels, including non-repeatable reads.
- **Environment Configuration**: Utilizes `.env` files for secure and flexible configuration.

## Prerequisites

- [Go (Golang)](https://golang.org/dl/) - An open-source programming language.
- [MySQL](https://www.mysql.com/) - A relational database management system.
- Set up a MySQL database and update the `.env` file with the correct database credentials.

## Installation

1. Clone the repository:

   ```sh
   git clone git@github.com:Oleksii-Morozov/RailwayConsole.git
   ```

2. Navigate to the project directory:

   ```sh
   cd RailwayConsole
   ```

3. Install dependencies:

   ```sh
   go get .
   ```

## Usage

1. Start the application:

   ```sh
   go run main.go
   ```

2. Use the command-line interface to interact with the system:
   - Type `h` for help on commands.
   - Use various commands (`create`, `read`, `update`, `delete`, etc.) to manage railway data.

## Configuration

- Configure the application using the `.env` file. Add your database connection string in the format:

  ```env
  CONNECTION_STRING=username:password@tcp(host:port)/dbname
  ```

## Contributing

Contributions to the RailwayConsole are welcome. Please ensure to follow the standard coding conventions and add unit tests for new features.

### Commits prefixes

    [feat] - new feature, that is added to an app
    [fix] - bugfixing
    [style] - changes related to code style
    [refactor] - refactoring of a piece of codebase
    [test] - everything, related to testing
    [docs] - everything that is related to documentation
    [config] - change configuration files
    [CI/CD] - changes in CI/CD

### Branches naming style

    feat/branch-name - for new features
    fix/branch-name - for bugfixing
    style/branch-name - for styling
    refactor/branch-name - for refactoring
    test/branch-name - for testing
    docs/branch-name - for documentation
    config/branch-name - for configuration


## License

This project is licensed under the [MIT License](LICENSE).
