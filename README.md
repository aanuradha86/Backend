# Backend
Library Management

I am running a small library for the community in these difficult times. I need to manage the books in the library using the book management system. I have implemented.


## Install dependencies

To start off, using mod:

```bash
go mod init <your-module-name>
```

## Run

```bash
go run .
```

This starts a server on port `9000`.

Once you run the application it automatically installs the dependencies, updates go.mod and creates go.sum file. Once the server is running, you will have to test the server using `./test-with-curl.sh`.

## Using modules

```bash
go mod download
```

## Instructions

### Implementation details

Map the CRUD (Create, Read, Update & Destroy) operations to the appropriate end-point based on the table below:

| CRUD  | Action   |   Method |     Path     |    Params  |      Body                                        |  Response             |
|-------|----------|----------|--------------|------------|--------------------------------------------------|-----------------------|
|  R    | Index    |   GET    |  /books      |  -         |    -                                             | [ {ID: , Title: ,} ]  |
|  R    | Show     |   GET    |  /books/{id} |  -         |    -                                             | {ID: , Title: ,}      |
| C     | Create   |   POST   |  /books      |  -         |   { Title: , Author: , }                         | {ID: , Title: ,}      |
|   U   | Update   |   PUT    |  /books/{id} |  -         |   { ID: , Title: , Author: , }                   | {ID: , Title: ,}      |
|   U   | Update   |   PATCH  |  /books/{id} |  -         |   {only send attributes which are to be updated} | {ID: , Title: ,}      |
|    D  | Delete   |  DELETE |  /books/{id} |  -         |    -                                             | {ID: , Title: ,}      |

### Test

```bash
./test-with-curl.sh
```

