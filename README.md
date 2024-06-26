# Cats Social

## 🌄 Background

Cats Social is an application where cat owners can match their cats with others for companionship.
Cats Social is part of Project Sprint Batch 2 Week 1 Project. Projectsprint is a sprint to create many resourceful projects in a short time with rigorous testing, load testing, and showcase. This initiative aims to demonstrate the ability to deliver high-quality, scalable applications quickly and effectively.

---

## 🚀 Getting Started

### Prerequisites

- Go (1.16 or later)
- Postgres
- [K6](https://k6.io/docs/get-started/installation/) for testing
- WSL (if on Windows)

### Running the Project

#### Environment Variables

Set the following environment variables:
```bash
export DB_NAME=your_db_name
export DB_PORT=5432


export DB_HOST=localhost
export DB_USERNAME=your_db_user
export DB_PASSWORD=your_db_password
export DB_PARAMS="sslmode=disabled"
export JWT_SECRET=your_jwt_secret
export BCRYPT_SALT=8
```

#### Running Migrations

```bash
migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations up
```

#### Running the Server

```bash
go run main.go
```

---

## 📜 Documentation

### API Endpoints

#### Authentication & Authorization

- **Register User** - `POST /v1/user/register`
- **Login User** - `POST /v1/user/login`

#### Manage Cats

- **Add Cat** - `POST /v1/cat`
- **Get Cats** - `GET /v1/cat`
- **Update Cat** - `PUT /v1/cat/{id}`
- **Delete Cat** - `DELETE /v1/cat/{id}`

#### Match Cats

- **Send Match Request** - `POST /v1/cat/match`
- **Get Match Requests** - `GET /v1/cat/match`
- **Approve Match Request** - `POST /v1/cat/match/approve`
- **Reject Match Request** - `POST /v1/cat/match/reject`
- **Delete Match Request** - `DELETE /v1/cat/match/{id}`

---

## 🧪 Testing

### Prerequisites
- [K6](https://k6.io/docs/get-started/installation/)
- A linux environment (WSL / MacOS should be fine)

### Environment Variables
- `BASE_URL` fill this with your backend url (eg: `http://localhost:8080`)

### Steps
1. Install [K6](https://k6.io/docs/get-started/installation/) and other prerequisites.
2. Run the server
3. Open the 'tests' folder 
4. Run the tests:
    #### For regular testing
    ```bash
    BASE_URL=http://localhost:8080 make run
    ```
    #### For load testing
    ```bash
    BASE_URL=http://localhost:8080 make runAllTestCases
    ```

---

## 📝 Requirements

### Functional Requirements

#### Authentication & Authorization

- **User Registration**
  - Endpoint: `POST /v1/user/register`
  - Request Body:
    ```json
    {
      "email": "email@example.com",
      "name": "First Last",
      "password": "password"
    }
    ```
  - Response: 
    ```json
    {
      "message": "User registered successfully",
      "data": {
        "email": "email@example.com",
        "name": "First Last",
        "accessToken": "access-token"
      }
    }
    ```
  - Errors:
    - `409` Conflict if email exists
    - `400` Validation errors
    - `500` Server error

- **User Login**
  - Endpoint: `POST /v1/user/login`
  - Request Body:
    ```json
    {
      "email": "email@example.com",
      "password": "password"
    }
    ```
  - Response:
    ```json
    {
      "message": "User logged in successfully",
      "data": {
        "email": "email@example.com",
        "name": "First Last",
        "accessToken": "access-token"
      }
    }
    ```
  - Errors:
    - `404` User not found
    - `400` Incorrect password or validation errors
    - `500` Server error

#### Manage Cats

- **Add Cat**
  - Endpoint: `POST /v1/cat`
  - Request Body:
    ```json
    {
      "name": "CatName",
      "race": "Persian",
      "sex": "male",
      "ageInMonth": 12,
      "description": "A friendly cat",
      "imageUrls": ["http://example.com/cat1.jpg"]
    }
    ```
  - Response:
    ```json
    {
      "message": "success",
      "data": {
        "id": "cat-id",
        "createdAt": "ISO 8601 date"
      }
    }
    ```
  - Errors:
    - `400` Validation errors
    - `401` Missing or expired token

- **Get Cats**
  - Endpoint: `GET /v1/cat`
  - Query Params: (all optional)
    - `id`, `limit`, `offset`, `race`, `sex`, `hasMatched`, `ageInMonth`, `owned`, `search`
  - Response:
    ```json
    {
      "message": "success",
      "data": [
        {
          "id": "cat-id",
          "name": "CatName",
          "race": "Persian",
          "sex": "male",
          "ageInMonth": 12,
          "imageUrls": ["http://example.com/cat1.jpg"],
          "description": "A friendly cat",
          "hasMatched": false,
          "createdAt": "ISO 8601 date"
        }
      ]
    }
    ```
  - Errors:
    - `401` Missing or expired token

- **Update Cat**
  - Endpoint: `PUT /v1/cat/{id}`
  - Request Path Params: `id`
  - Request Body:
    ```json
    {
      "name": "CatName",
      "race": "Persian",
      "sex": "male",
      "ageInMonth": 12,
      "description": "A friendly cat",
      "imageUrls": ["http://example.com/cat1.jpg"]
    }
    ```
  - Response:
    ```json
    {
      "message": "success",
      "data": {
        "id": "cat-id",
        "updatedAt": "ISO 8601 date"
      }
    }
    ```
  - Errors:
    - `400` Validation errors
    - `401` Missing or expired token
    - `404` Cat not found

- **Delete Cat**
  - Endpoint: `DELETE /v1/cat/{id}`
  - Request Path Params: `id`
  - Response:
    ```json
    {
      "message": "success"
    }
    ```
  - Errors:
    - `401` Missing or expired token
    - `404` Cat not found

#### Match Cats

- **Send Match Request**
  - Endpoint: `POST /v1/cat/match`
  - Request Body:
    ```json
    {
      "matchCatId": "cat-id",
      "userCatId": "cat-id",
      "message": "Looking forward to a playdate"
    }
    ```
  - Response:
    ```json
    {
      "message": "success",
      "data": {
        "id": "match-id",
        "createdAt": "ISO 8601 date"
      }
    }
    ```
  - Errors:
    - `400` Validation errors
    - `401` Missing or expired token
    - `404` Cat not found or not owned by user

- **Get Match Requests**
  - Endpoint: `GET /v1/cat/match`
  - Response:
    ```json
    {
      "message": "success",
      "data": [
        {
          "id": "match-id",
          "issuedBy": {
            "name": "User Name",
            "email": "email@example.com",
            "createdAt": "ISO 8601 date"
          },
          "matchCatDetail": {
            "id": "cat-id",
            "name": "CatName",
            "race": "Persian",
            "sex": "male",
            "description": "A friendly cat",
            "ageInMonth": 12,
            "imageUrls": ["http://example.com/cat1.jpg"],
            "hasMatched": false,
            "createdAt": "ISO 8601 date"
          },
          "userCatDetail": {
            "id": "cat-id",
            "name": "CatName",
            "race": "Persian",
            "sex": "male",
            "description": "A friendly cat",
            "ageInMonth": 12,
            "imageUrls": ["http://example.com/cat1.jpg"],
            "hasMatched": false,
            "createdAt": "ISO 8601 date"
          },
          "message": "Looking forward to a playdate",
          "createdAt": "ISO 8601 date"
        }
      ]
    }
    ```
  - Errors:
    - `401` Missing or expired token

- **Approve Match Request**
  - Endpoint: `POST /v1/cat/match/approve`
  - Request Body:
    ```json
    {
      "matchId": "match-id"
    }
    ```
  - Response:
    ```json
    {
      "message": "success"
    }
    ```
  - Errors:
    - `400` Match ID invalid
    - `401` Missing or expired token
    - `404` Match ID not found

- **Reject Match Request**
  - Endpoint: `POST /v1/cat/match/reject`
  - Request Body:
    ```json
    {
      "matchId": "match-id"
    }
    ```
  - Response:
    ```json
    {
      "message": "success"
    }
    ```
  - Errors:
    - `400` Match ID invalid
    - `401` Missing or expired token
    - `404` Match ID not found

- **Delete Match Request**
  - Endpoint: `DELETE /v1/cat/match/{id}`
  - Request Path Params: `id`
  - Response:
    ```json
    {
      "message": "success"
    }
    ```
  - Errors:
    - `400` Match already approved or rejected
    - `401` Missing or expired token
    - `404` Match ID not found

### Non-Functional Requirements

- Backend:
  - Golang with any web framework
  - Postgres database
  - Port: 8080
  - No ORM/Query generator; use raw queries
  - No external caching
  - Environment Variables:
    ```bash
    export DB_NAME=
    export DB_PORT=
    export DB_HOST=
    export DB_USERNAME=
    export DB_PASSWORD=
    export DB_PARAMS="sslmode=disabled"
    export JWT_SECRET=
    export BCRYPT_SALT=8
    ```
  - Compile binary as `main` using:
    ```bash
    env GOARCH=amd64 GOOS=linux go build -o main
    ```

### Database Migration

- Use [golang-migrate](https://github.com/golang-migrate/migrate) for managing database migrations:
  - Create migration:
    ```bash
    migrate create -ext sql -dir db/migrations add_user_table
    ```
  - Execute migration:
    ```bash
    migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations up
    ```
  - Rollback migration:
    ```bash
    migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations down
    ```

---

## 👥 Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

---

## 📝 License

-

---

## 📚 Resources

- **Notion:** [Cats Social Notion Page](https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769)
- **Tests:** [Project Sprint Batch 2 Week 1 Test Cases](https://github.com/nandanugg/ProjectSprintBatch2Week1TestCases)

---

## 📞 Contact

[Muhammad Ulil 'Azmi](https://github.com/ulilazmi100) - [@M_Ulil_Azmi](https://twitter.com/M_Ulil_Azmi)

Project Link: [https://github.com/ulilazmi100/CatsSocial](https://github.com/ulilazmi100/CatsSocial)