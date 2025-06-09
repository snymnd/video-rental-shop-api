# Video Rental Shop (VRS) API

## App Requirement

1. Users can register to the application with a user role.
2. Users can log in to the application.
3. Admin accounts are pre-generated in the database.
4. Admins can log in to the application with an admin role.
5. Admins can add new video records in DVD, Blu-ray, and VHS formats with various prices.
6. Everyone can view a list of videos available for rent.
7. Everyone can search by title and filter the list of available videos.
8. Users can rent multiple video titles in one payment transaction.
9. Payment is only valid for 3 hours after being generated.
10. Users cannot rent more than 1 copy of the same video title and format, but can rent the same title in different formats.
11. Users can rent videos via the online application and collect the physical videos at the offline store later.
12. Users can return videos only at the offline store.
13. Users can return all rented videos at once or return only some videos and the rest later.
14. Every rental has a due date of 3 days starting from when the rental(s) are paid.
15. If videos are returned after the due date, users will be charged a late fee of $2 for each day beyond the due date, with a maximum cap of $30.
16. Only admins are allowed to record returned videos by inputting user data (user ID) and processing the videos returned by the user.

## Entity Relationship Diagram (ERD)

- Aplication business ERD
  ![video rental shop api erd](./assets/video-rental-shop-erd.png)
  - `videos` - Stores information about available videos likes title, description, stock, and genre
  - `users` - Contains user account details such as name, email, and contact information
  - `rentals` - Tracks video borrowing transactions with rental dates and status
  - `payments` - Records financial transactions related to video rentals
  - `genres` - Categorizes videos by different types/categories
    <br><br>
- Role base access control (RBAC) ERD
  ![role base access control erd](./assets/rbac-erd.png)
  - `roles` - Defines user roles such as admin or user
  - `permissions` - Specifies allowed actions like create, read, update, and delete
  - `resource` - Represents system resources that can be accessed (videos, users, rentals)
  - `rbac` - Junction table that connects roles with permissions for specific resources

> You can access full diagram here [video-rental-shop-dbdiagram.io](https://dbdiagram.io/d/video-rental-shop-erd-6836d3116980ade2ebd7538e)

## Architecture

This project is using clean architecture, a software architecture that emphasizes the separation of concerns, the independence of components, and the use of well-defined boundaries.

![Clean Architecture](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

```bash
...
├── delivery               # Interface adapters layer - Handles delivery method for data access
│   └── rest
│       ├── middleware     # Contains request interceptors (auth, logging, etc.)
│       │   ├── ..._middleware.go
│       ├── ..._controller.go  # REST API handlers/controllers
│       ├── route
│       │   └── route.go   # API endpoint definitions and routing
├── dto                    # Data Transfer Objects - Defines request/response data structures
│   ├── ..._dto.go
├── entity                 # Enterprise business rules - Core domain models
│   ├── ..._entity.go
├── repository             # Data access layer - Implements data source interfaces
│   ├── postgresql         # Database implementation using PostgreSQL
│   │   ├── ..._repository.go
│   └── redis              # Cache implementation using Redis
│       ├── ..._cache_repository.go
├── usecase                # Application business rules - Contains business logic
│   ├── ..._usecase.go
└── util                   # Shared utilities and helper functions
```

## Build With

![Go](https://img.shields.io/badge/-Go-333333?style=flat&logo=go)
![Gin](https://img.shields.io/badge/-Gin-333333?style=flat&logo=gin)
![Postgresql](https://img.shields.io/badge/-Postgresql-333333?style=flat&logo=postgresql)
![Redis](https://img.shields.io/badge/-Redis-333333?style=flat&logo=redis)

## API Documentations

You can access API documentation here https://documenter.getpostman.com/view/21669206/2sB2x3otYN

Or you can import the
[video-rental-shop-api.postman_collection.json](./video-rental-shop-api.postman_collection.json) using postman web or app.

## Getting Started

### Run locally

1. Make sure you have below dependencies installed
   - Go ver. 1.24.0+
   - Postgres
   - Redis
2. Clone this repository with

   ```bash
   git clone https://github.com/snymnd/video-rental-shop-api.git
   ```

   and change your working directory to the clone repository folder

3. Install neccessary dependency with
   ```go
   go mod tidy
   ```
4. Set up environment with copy `.env.example` file to `.env` file and change the environment value as needed
5. Create database on your postgres with name following your the value of`DATABASE_NAME` field on your `.env` file
6. Run migration and seeder with bellow command
   ```go
   go run cmd/db/main.go
   ```
7. Finnaly! Run the application with bellow command
   ```go
   go run cmd/web/main.go
   ```
8. You can test the api by try to access/hit below endpoint <br>
   `http://localhost:<port>/api/v1/welcome` <br> with `port` following defined value from `SERVER_ADDRESS` field on your `.env` file
