## Introduction

**vhp-golang-rest-api** is a web application built with the Golang programming language using the Fiber Web Framework and MongoDB. The application provides a user management system with features such as authentication, authorization, password recovery, password change, user data export to Excel, and CRUD operations for managing users.

## Key Technologies

- [Golang](https://golang.org): The programming language used to build the application.
- [Fiber Web Framework](https://gofiber.io): A fast, simple, and efficient web framework for Golang.
- [MongoDB](https://www.mongodb.com): A document-based NoSQL database used for storing user data.
- [JWT (JSON Web Tokens)](https://jwt.io): Used for user authentication and authorization.
- [Logger](https://github.com/username/logger): Logger to log application requests and responses.
- [Limiter](https://github.com/username/limiter): Limiter to limit the frequency of requests from specific IP addresses to prevent Denial of Service (DDoS) attacks.

## Features

- **Authentication**: Supports user authentication using JWT, ensuring that only logged-in users can access protected routes.
- **Authorization**: Allows controlling user access rights based on roles and permissions.
- **Password Recovery**: Provides a password recovery feature, allowing users to reset their passwords via email confirmation.
- **Password Change**: Enables logged-in users to change their passwords.
- **Export to Excel**: Supports exporting user data to an Excel file for download.
- **CRUD User**: Allows adding, viewing, editing, and deleting user information.

## Installation and Run

1. Make sure you have [Golang](https://golang.org) installed on your computer.

2. Clone the repository from GitHub:

```bash
git clone https://github.com/username/repo.git
cd repo
```

3. Install the required dependencies:
```bash
go mod download
```

4. Run the application:
```bash
go run main.go
```

The application will run on the default port 8087. You can access the application in your web browser at http://localhost:8087.

## Contribution
If you wish to contribute to the project, please create a pull request or report issues. We welcome contributions from the community.
