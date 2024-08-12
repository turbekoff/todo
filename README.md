# todo
This is a simple REST API written in the Golang programming language using [chi](https://github.com/go-chi/chi) as a router and MongoDB for data storage. It allows users to manage tasks efficiently.


## Managing Configurations
When the api starts, it loads configuration from the .env file specified in the CONFIG_PATH environment variable, or from local variables if the CONFIG_PATH variable is empty.


| Environment variable | Required | Description |
| --- | --- | --- |
| DEBUG_MODE | false | `local`, `development` or `production` mode |
| PASSWORD_PEPPER | true | secret added to password during hashing |
| JWT_SIGNING_KEY | true | signing key |
| JWT_ACCESS_TOKEN_TTL | false | TTL of the access token |
| JWT_REFRESH_TOKEN_TTL | false | TTL of the refresh token |
| HTTP_HOST | false | HTTP Server host |
| HTTP_PORT | false | HTTP Server port |
| HTTP_READ_TIMEOUT | false | HTTP Server read timeout |
| HTTP_WRITE_TIMEOUT | false | HTTP Server write timeout |
| HTTP_MAX_HEADER_BYTES | false | The maximum number of bytes of the HTTP server header |
| MONGO_USER | false | Username of the mongo database |
| MONGO_PASSWORD | false | Password of the mongo database |
| MONGO_DATABASE | false | Name of the mongo database |
| MONGO_URI | true | URI of the mongo database |


## Project layout
```
 > build           docker containers of applications
 > cmd             main applications
 > configs         configuration files for different applications
 > internal \      private application and library code
    - delivery     the external layer which interacts with the outside world
    - domain       represents business logic and rules
    - repository   implementations of storages
    - service      use-case specific operations
 > pkg \           public library code
    - jwt          library for creating jwt manager
    - hash         library for creating argon2id hasher
```

## Auth Endpoints
Authentication is performed using the Bearer method and refresh session. You can create up to 10 authentication sessions, after exceeding count it will be erased. You cannot use the generated token from other.

```
Path: `/api/v1/signup`
Method: `POST`
Authorization: public
Request:    
    {
        "name": "Example2004",
        "password": "Example#2004"
    }
Responces:
    - 200
    - 400 {
        "code": 400,
        "message": "name must consist of 8-30 latin letters and digits"
    }
```


```
Path: `/api/v1/signin`
Method: `POST`
Authorization: public
Request:    
    {
        "name": "Example2004",
        "password": "Example#2004"
    }
Responces:
    - 200 {
        "accessToken": "Generated access token",
        "refreshToken": "Generated refresh token",
        "expireAt": "Access token expire time"
    }
    - 400 {
        "code": 400,
        "message": "user doesn't exist"
    }
```

```
Path: `/api/v1/logout`
Method: `POST`
Authorization: Bearer required
Request:    
    -
Responces:
    - 200
    - 400 {
        "code": 400,
        "message": "user doesn't exist"
    }
```


```
Path: `/api/v1/profile`
Method: `GET`
Authorization: Bearer required
Request:    
    -
Responces:
    - 200 {
        "name": "Example2004",
        "createdAt": "created time",
        "updatedAt": "updated time"
    }
    - 400 {
        "code": 400,
        "message": "user doesn't exist"
    }
```


```
Path: `/api/v1/profile`
Method: `PUT`
Authorization: Bearer required
Request:    
    {
        "name": "Example2004",
        "password": "Example#2004"
    }
Responces:
    - 200 {
        "name": "Example2004",
        "createdAt": "created time",
        "updatedAt": "updated time"
    }
    - 400 {
        "code": 400,
        "message": "user doesn't exist"
    }
```

```
Path: `/api/v1/profile`
Method: `DELETE`
Authorization: Bearer required
Request:
    -
Responces:
    - 200
    - 400 {
        "code": 403,
        "message": "you don't have authorization to view this task"
    }
```

## Task Endpoints
```
Path: `/api/v1/task`
Method: `POST`
Authorization: Bearer required
Request:    
    {
        "name": "Example2004",
        "completed": true
    }
Responces:
    - 200
    - 400 {
        "code": 400,
        "message": "empty name"
    }
```

```
Path: `/api/v1/task`
Method: `GET`
Authorization: Bearer required
Request:
    -
Responces:
    - 200 {
        [
            {
                "id": "Task ID",
                "name": "Task name",
                "completed": true,
                "createdAt": "created time",
                "updatedAt": "updated time"
            }
        ]
    }
```

```
Path: `/api/v1/task/{id}`
Method: `GET`
Authorization: Bearer required
Request:
    -
Responces:
    - 200 {
        "id": "Task ID",
        "name": "Task name",
        "completed": true,
        "createdAt": "created time",
        "updatedAt": "updated time"
    }
    - 400 {
        "code": 403,
        "message": "you don't have authorization to view this task"
    }
```

```
Path: `/api/v1/task/{id}`
Method: `PUT`
Authorization: Bearer required
Request:    
    {
        "name": "Example2004",
        "completed": true
    }
Responces:
    - 200 {
        "id": "Task ID",
        "name": "Example2004",
        "completed": true,
        "createdAt": "created time",
        "updatedAt": "updated time"
    }
    - 400 {
        "code": 403,
        "message": "you don't have authorization to view this task"
    }
```

```
Path: `/api/v1/task/{id}`
Method: `DELETE`
Authorization: Bearer required
Request:
    -
Responces:
    - 200
    - 400 {
        "code": 403,
        "message": "you don't have authorization to view this task"
    }
```
