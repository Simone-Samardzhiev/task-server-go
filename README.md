# task-server-go

Web server that will store task with features like, register
login for users and get, add, update, delete tasks for tasks.

## Table of content

- [About the project](#about-the-project)
- [Installation](#installation)
- [API](#api)

## About the project

Web server used to help user stay organized by storing tasks. It includes

- JWT-based authentication
- Task creation, updating and deleting
- Database integration with postgres

## Installation

Follow the steps to build and run the server locally.

### Prerequisites

1. Go(1.23 or later)
2. PostgresSQL
3. Git

### Steps

1. **Clone the repository**.

```bash
git clone https://github.com/your-username/server-go.git
cd server-go
```

2. **Add .env file**

```ini
SERVER_ADDR=Server address
DATABASE_URL=The database url for connection
MAX_OPEN_CONNECTIONS=Max open db connections
MAX_IDLE_CONNECTIONS=Max open idle connections.
JWT_SECRET=Secret used to hash tokens.
JWT_ISSUER=Issuer of the tokens.
```

3. **Build and run**

```bash
go build
./server
```

## API

### 1. **POST api/v1/users/register**

The endpoint allows users to register.

#### **Request Body**

The body of the request should contain the user credentials

```json
{
  "email": "exmaple@email.com",
  "username": "Someone",
  "password": "Password_123"
}
```

The user payload is validated before being accepted.
If the email or the username is already in use the API will return an error.
Also, there are more requirements for the user credentials:

1. The email should be properly formated with valid local and domain part.
2. The username should be less than 8 letters
3. The password should be secure:
    1. At least one capital letter
    2. At least one small letter
    3. At least one number.
    4. At least one special character(! " # $ % & ' ( ) * + , - . : ; < = > ? [ \ ] ^ _ `{ | } ~)

### 2. POST api/v1/users/login

The endpoint allows user to receive JWT refresh and access token.

#### **Request Body**

The body of the request should contain user credentials

```json
{
  "email": "exmaple@email.com",
  "username": "Someone",
  "password": "Password_123"
}
```

#### **Response**

If the credential the server will return **Status Code Unauthorized**.  
If not the response will be like:

```json
{
  "refresh_token": "token",
  "access_token": "token"
}
```

### 3. GET api/v1/users/refresh

The endpoint allows user to send refresh to token, for a new refresh and access token.

#### **Header**

Authorization: Bearer + refresh token

#### **Response**

If the token is expired the server will return **Status Code Unauthorized**.  
If not the response will be like:

```json
{
  "refresh_token": "token",
  "access_token": "token"
}
```

### 4. GET api/v1/tasks/get

The endpoint allows user to get all their tasks.

#### **Header**

Authorization: Bearer + refresh token

#### **Response**

If the token is expired the server will return **Status Code Unauthorized**.  
If not the response will be like:

```json
[
  {
    "id": "ffafdd8a-20ba-452f-b5b4-37d98b091ba0",
    "name": "Task name",
    "description": "Task description",
    "priority": "Low",
    "date": "2025-03-15T16:03:30Z"
  }
]
```

### 5. POST api/v1/tasks/add

The endpoint allows user to add a new task.

#### **Header**

Authorization: Bearer + refresh token

#### **Request body**

The body should contain the task information

```json
{
  "name": "Name",
  "description": "Description",
  "priority": "Low",
  "data": "2025-03-15T16:03:30Z"
}
```

The task payload is also validated before storing it.
None of the filed can be empty. Also, the priority will be checked by the database.
You could easily adjust the priority by updating **Priorities** table

#### **Response**

If the token is expired the server will return **Status Code Unauthorized**.  
If not the response will be like:

```json
{
  "id": "ffafdd8a-20ba-452f-b5b4-37d98b091ba0",
  "name": "Name",
  "description": "Description",
  "priority": "Vital",
  "date": "2025-03-15T16:03:30Z"
}
```

### 6. PUT api/v1/tasks/update

The endpoint allows user to update an existing token.

#### **Header**

Authorization: Bearer + refresh token

#### **Request body**

The body should contain the token information

```json
    {
  "id": "ffafdd8a-20ba-452f-b5b4-37d98b091ba0",
  "name": "Task name",
  "description": "Task description",
  "type": "High",
  "date": "2025-03-15T16:03:30Z"
}
```
Note that updating task also validate the payload.

#### **Response**

If the token is expired the server will return **Status Code Unauthorized**.  
If the task is found the server will return **Status Code OK**
If the task is not found the server will return **Status Code Not Found**

### 7. DELETE api/v1/tasks/delete/{id}

The endpoint allows user to delete a task.

#### **Header**

Authorization: Bearer + refresh token

#### **Params**

**id** The id of the token

#### **Response**

If the token is expired the server will return **Status Code Unauthorized**.  
If the task is found the server will return **Status Code OK**
If the task is not found the server will return **Status Code Not Found**