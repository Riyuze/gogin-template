# 📃 Go-Gin Template

This template provides a starting point for creating a robust backend with a **command pattern** approach, accomodating the creation of **REST APIs**, **Consumers**, and **Scheduled Workers** without sacrificing the capabilities of multi container deployment.

This template is written in **Go**, utlizing the **Gin Gonic** HTTP web framework for its REST APIs development.

---

## 📖 Table of Contents

- [📃 Go-Gin Template](#-go-gin-template)
  - [📖 Table of Contents](#-table-of-contents)
  - [🏗️ Project Structure](#️-project-structure)
    - [- 🖥️ cmd](#--️-cmd)
    - [- 🧱 baselib](#---baselib)
    - [- 🧰 bootstrap](#---bootstrap)
    - [- ⛓️ internal](#--️-internal)
    - [- 📋 docs](#---docs)
  - [🧭 Development Guideline](#-development-guideline)
    - [- 🤝 Convention](#---convention)
    - [- 🔥 REST Standard](#---rest-standard)
    - [- 📨 Response Body](#---response-body)

---

## 🏗️ Project Structure

### - 🖥️ cmd

As mentioned beforehand, this template was created to run several different services. Each services will then be deployed on seperate containers where each container will run a different service by calling its corresponding subcommands.

```bash
# example usage
go run main.go rest # REST API service
go run main.go cdc # CDC service
go run main.go kafka_consumer <consumer_name> # Kafka consumer service
```

Each subcommand must be defined in the `cmd` folder. One file will represent one subcommand. After creating each subcommand, be sure to remember to register it in the `root.go` file.

```go
// example subcommand
func testCommand() *cobra.Command{
 return &cobra.Command{
  Use: "test",
  Short: "Run Test Service"
  RunE: func(cmd *cobra.Command, args []string) error {
   //Dependency injection and component preparation here
  }
 }
}
```

Refer to the [cobra](https://github.com/spf13/cobra) library to read further about commands and its complete features.

### - 🧱 baselib

The `baselib` folder consists of structs and functions that are used as standarizations for developing this application. The contents of this folder rarely changes as development are primarily done on the `internal` folder.

### - 🧰 bootstrap

The `bootstrap` folder are filled with components that are needed for this application. These includes database connection, logger, rest-client, and other components that will be injected into the `internal` functions.

To make use of the components, call `bootstrap.Init()` on the subcommand. This function will return a `bootstrap.Container` object which can be used to inject the desired bootstrap components unto the `internal` components (**controller**, **servie**, **repository**, and etc) that are used in that subcommand.

```go
// Examples of tools that can be accessed from bootstrap.Container

cfg := bootstrap.Init() // Initialize bootstrap

// Logging messages (logrus: github.com/sirupsen/logrus)
cfg.Logger().Info("Ini Info")
cfg.Logger().Error("Ini Error")
cfg.Logger().Debug("Ini Debug")

// Accessing config reader (viper: github.com/spf13/viper)
viper := cfg.GetConfig()
baseUrl := viper.GetString("base.url")

// Retrieve database connections
dbr := cfg.Dbr() // db with write capability
dbw := cfg.Dbw() // db with read only capability

// Used for copying structs (jinzhu copier: github.com/jinzhu/copier)
cfg.CopyStruct(&from, &to)

```

### - ⛓️ internal

The `internal` folder contains the logic for the application. Most of the developements will be done inside this subfolder.

The development of this application follows the **Layered Architecture Application** principle, where the logic for this application is divided into several components as follows:

-   **viewmodel**: Collections of structs defining the data structure that is used to transfer data into external applications.
-   **controller**: Component which handles routing and handler functions. This component are mainly responsible for binding JSON request into viewmodel, validating scopes, calling service response, and returning the response through the endpoint.
-   **service**: Component where the business logic is written. Usually composed of validating the request, business logic, fetching data from repositories or REST clients, and mapping data.
-   **respository**: Component where code relating to persistance layer exists. This component mediates the process between the data access logic and the data sources by encapsulating the logic required to access data sources.
-   **model**: Collections of structs defining the data structure that is used for query results from the repository.

### - 📋 docs

The `docs` folder is generated using the library [swaggo](https://github.com/swaggo/swag) which is used for API docummentation. When building the project, make sure to follow these steps to enable the generation of swagger API documentation.

1. Make sure that you have installed installed swaggo by running the following command `go install github.com/swaggo/swag/cmd/swag@latest` in the root path of the project, where the `main.go` resides.
2. Add these general API below in `rest.go` or where your REST API is served. (This project utilizes [gin-gonic](https://github.com/gin-gonic/gin)).

```go
// @title           			Swagger Sample API
// @version         			1.0
// @description     			This is a sample api server.
// @termsOfService  			http://swagger.io/terms/
// @contact.name   				API Support
// @contact.url    				http://www.swagger.io/support
// @contact.email  				support@swagger.io
// @license.name  				Apache 2.0
// @license.url   				http://www.apache.org/licenses/LICENSE-2.0.html
// @host      					localhost:8080
// @BasePath  					/api
// @securityDefinitions.basic  	BasicAuth
// @externalDocs.description  	OpenAPI
// @externalDocs.url          	https://swagger.io/resources/open-api/
func restCommand() *cobra.Command {
...
```

Further read: [General API](https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info)

3. Run the `swag init` command in the root directory where `main.go` is. If the command succeeds, a folder called `docs` is generated.
4. Add import package in the file where the REST API is served. (Same as step number 2).

```go
import (
	_ "gogin-template/docs" //package docs yg digenerate ketika command swagger init dijalankan
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
```

5. Register an endpoint for the swagger to serve its webpage.

```go
ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

6. Run your `rest.go` subcommand and access the endpoint. Usually its `http://localhost:8080/swagger/index.html`.

For every API that is developed, follow these steps to add swagger documentation.

1. Add swagger API documentation remarks. Below is an example for a GET request with query params. For other types of request and response you can refer to this [swaggo](https://github.com/swaggo/swag?tab=readme-ov-file#api-operation) documentation.

2. Re-run the `swag init` command so that the docs are regenerated.

[[↑ Back to top ↑](#-go-gin-template)]

---

## 🧭 Development Guideline

### - 🤝 Convention

Naming conventions on this template are as follows:

-   **File Names**: Snake case + (\_model/\_repository/\_service/\_viewmodel/\_controller/\_client)

```
user_model.go
user_repository.go
user_service.go
user_viewmodel.go
user_controller.go
ecmapi_client.go
```

-   **Model Struct Names** : Pascal case + Model
-   **Model Fields Names** : Pascal case
-   **Database Column Names** : Snake case

```go
type ExampleObjectModel struct {
	ExampleId        int64      `db:"example_id"`
}
```

-   **Repository Names** : Pascal case + (Repository/RepositoryImpl)
-   **Repository Constructor Names** : Pascal case -> (New + Interface Name)

```go
// Interface
type ExampleRepository interface {
	GetExamples(...) (..., error)
	GetExample(...) (..., error)
	SetExample(...) (..., error)
}

// Implementation
type ExampleRepositoryImpl struct {
	dbr *sqlx.DB
	dbw *sqlx.DB
	cfg *bootstrap.Container
}

func NewExampleRepository(dbr *sqlx.DB, dbw *sqlx.DB, cfg *bootstrap.Container) ExampleRepository {
	return &ExampleRepositoryImpl{dbr: dbr, dbw: dbw, cfg: cfg}
}
```

-   **Service Names** : Pascal case + (Service/ServiceImpl)
-   **Service Constructor Names** : Pascal case -> (New + Interface Name)

```go
// Interface
type ExampleService interface {
	GetExamples(...) (..., error)
	GetExample(...) (..., error)
	SetExample(...) (..., error)
}

// Implementation
type ExampleServiceImpl struct {
	repository repository.ExampleRepository
	client     restclient.ExampleAPIClient
	cfg        *bootstrap.Container
}

func NewExampleService(repo repository.ExampleRepository,
client restclient.ExampleAPIClient,
cfg *bootstrap.Container) ExampleService {
	return &ExampleServiceImpl{repository: repo, client: client, cfg: cfg}
}
```

-   **Controller Names** : Pascal case + (Controller)
-   **Controller Constructor Names** : Pascal case -> (New + Controller Name)

```go
type ExampleController struct {
	service service.ExampleService
	cfg     *bootstrap.Container
}

func NewExampleController(service service.ExampleService, server *gin.Engine, cfg *bootstrap.Container) {
	controller := &ExampleController{
		service: service,
		cfg:     cfg,
	}

	routes := server.Group("/example")
	{
		routes.GET("/list", controller.GetExamples)
		routes.GET("/list/:id", controller.GetExample)
		routes.PUT("", controller.ProcessExample)
	}
}
```

-   **Viewmodel Struct Names** : Pascal case + ViewModel
-   **Viewmodel Field Names** : Pascal case
-   **JSON key Names** : Camel case

```go
type ExampleObjectRqCreateViewModel struct {
	ExampleZone  string     `json:"exampleZone"`
}
```

### - 🔥 REST Standard

APIs in this template follow the standard HTTP requests. The most common used HTTP requests for this template can be simplified as follows:

-   **HTTP GET**: GET is used for **requesting data** from a specific source. GET requests can be cached, remain in the browser history, can be bookmarked, and should _never_ be used when dealing with sensitive data. Theres are 2 main ways the GET request is utilized in this template:

    -   Returning a **single data**, which is commonly used to inquire details about a certain resource. There are 2 methods on how to do this:
        -   **Path parameters**: Path parameters are variable parts of a URL path. They are typically used to point to a specific resource within a collection. When the resource is not available, a 404 error will be returned. <br />
            Example: `http://localhost:8080/user/1`
        -   **Query parameters**: Query parameters are a defined set of parameters (key-value pair) attached to the end of a URL used to provide additional information to a web server when making requests. <br />
            Example: `http://localhost:8080/users?id=1`
    -   Returning **multiple data**, usually used to inquire a list of data. Usually this kind of request needs to implement pagination to limit the amount of data being sent across the network. There are 2 main ways to achieve this, using a **paging request** or a **cursor request**.

        -   A **paging request** is when we inqure data according to the _page_ and the _index_, so that the client can jump through pages to retrieve the data needed. The drawback of this type of request is that when the _index_ that is used becomes too large, table index will be ignored, causing performance drops. This is because queries using this type of request will use the `OFFSET` keyword. <br />
            Example: `http://localhost:8080/products?pageIndex=0&rowsPerPage=10&orderColumn=product_code&orderDirection=ASC` <br />
            Here are some standard query parameters that are used on a paging request:
            -   `pageIndex`: Used to request a page where the data is. (Starts from 0)
            -   `rowsPerPage`: Defines how many data are in each page.
            -   `orderColumn`: Sorts the data according to a specific column.
            -   `orderDirection`: Defines the direction of the sort. (ASC / DESC)
        -   A **cursor request** doesn't utilize pages and indexes. Data is sorted according to a primary key or a composite key, and the queries are built so that table index will always be utilized. Because of the sorting, the cursor can only jump to the next page and cannot jump accross multiple pages to retrieve the needed data. <br />
            Example: `http://localhost:8080/products?rowsPerPage=10&cursor=10` <br />
            In the example above, assume we are using a products database which is sorted by its ID. By using the ID as the cursor, we can get the next set of data.
            ```sql
            SELECT *
            FROM products
            WHERE ID > {cursor}
            LIMIT {rowsPerPage}
            ```
            We will then return the last ID of the data retrieved as the next cursor. <br />
            Here are some standard query parameters that are used on a paging request:
            -   `cursor`: The last data identifier in the current page.
            -   `rowsPerPage`: Defines how many data are in each page. <br />
        -   These 2 pagination methods can be implemented using the **encoded cursor**. The **encoded cursor** suggests returning an encoded base64 string regardless of the underlying pagination solution. This allows the server to implement different underlying pagination solutions while providing a consistent interface to the API consumers.
            -   When using a **paging request**, we encode the _current page_ and _total page_ into a base64 string and return it as a cursor to the clients.
            ```JSON
            response: {
            	// "page=3|offset=20|total_pages=30"
            	next_cursor: "dcjadfaXMDdQTQ"
            }
            ```
            -   Similarly, we can encode the cursor in the **cursor pagination** into a base64 string before returning it to the clients.
            ```JSON
            response: {
            	// "next_cursor:1234"
            	next_cursor: "dcjadfaXMDdQTQ"
            }
            ```

-   **HTTP POST**: POST requests that a web server accepts the data enclosed in the body of the request message, most likely for storing it. It is often used when uploading a file or when submitting a completed web form. Data is usually sent in the form of JSON.
-   **HTTP PUT**: PUT method is used to create a new resource or replace a resource. It's similar to the POST method, in that it sends data to a server, but it's idempotent. This means that the effect of multiple PUT requests should be the same as one PUT request. Data is also usually sent in the form of JSON, the only difference with POST is that there is usually a field to determine which data is getting updated.
-   **HTTP DELETE**: DELETE is used to delete a resource, such as a file or a database record. DELETE is idempotent, meaning that making multiple identical requests should have the same effect as making a single request. However, it’s important to note that the actual deletion of a resource depends on the server’s implementation and policies. Upon receiving the DELETE request, the server processes it and removes the specified resource if it exists, returning a status code to indicate the success or failure of the operation.

Aside from requests, the responses in this template follow the common HTTP responses, which are as follow:

-   **OK (HTTP 200)**: This reponse is returned when a requests sucessfully executes.
-   **Bad Request (HTTP 400)**: This response is returned when the request doesn't fulfill the validation conditions.
    -   **Unauthorized (HTTP 401)**: This response code is returned when authorization fails. This response will be returned automatically by the authorization middleware, which handles the authorization tokens. B
    -   **Not Found (HTTP 404)**: This response is returned when the data is not found when inquired. For multiple data inquiry, when no data is found, usually its best to still return the OK (HTTP 200) status along with an empty array.
-   **Internal Server Error (HTTP 500)**: This response is return when an unhandled error occurs.

### - 📨 Response Body

The standard response body used in this template are as follows:

```JSON
{
 "responseCode" : "",
 "responseMessage" : "",
 "data" : null,
 "logReff" : "",
 "traceId" : "",
 "pageInfo" : {
  "currentPageIndex" : 0,
  "maxPageIndex" : 1,
  "rowsPerPage" : 10,
  "totalAvailableItems" : 15
 },
}
```

The usage of each fields are as follows:

-   **responseCode**: This field can be used to map response code in the front-end / client to render an UI and error messages according to the response code. This field can be assigned in the controller or when returning an error using the `exception` package. <br />
-   **responseMessage**: This field is used to return an error description. The front-end / client may use this information for logging or to show an error message. Same as _responseCode_ this field can be assigned in the controller or by returning an `exception`.
-   **logReff**: On every request, this field will automatically be generated with a GUID that is unique. This GUID is used for log searching in log aggregators. This GUID is sent to the front-end / client so that if an error occurs, this information can be used for determining the root cause of the problem.
-   **traceId**: Similar to _logReff_, this field is also automatically generated for every request that happens in the API. This GUID is used for searching up the trace is _jaeger_.

    Below is an example of logs which contains both _logReff_ and _traceId_ data.

```shell
{"component":"rest","level":"info","logReff":"5f65999f-77cb-489d-a650-1c122eab4d30","msg":"METHOD : GET","time":"2023-08-04T10:21:19+07:00","traceId":"7b91f82658f4a83312e9fa4e943bc9de"}
{"component":"rest","level":"info","logReff":"5f65999f-77cb-489d-a650-1c122eab4d30","msg":"ENDPOINT : /sample","time":"2023-08-04T10:21:19+07:00","traceId":"7b91f82658f4a83312e9fa4e943bc9de"}
{"component":"rest","level":"info","logReff":"5f65999f-77cb-489d-a650-1c122eab4d30","msg":"REQUEST BODY : ","time":"2023-08-04T10:21:19+07:00","traceId":"7b91f82658f4a83312e9fa4e943bc9de"}
{"component":"rest","level":"info","logReff":"5f65999f-77cb-489d-a650-1c122eab4d30","msg":"RESPONSE STATUS : 200","time":"2023-08-04T10:21:19+07:00","traceId":"7b91f82658f4a83312e9fa4e943bc9de"}
```

-   **data**: This field contains the data that is returned by the API. The data returned could be a single or array of object. The data type can be changed to acommodate the needs of every service as this is a generic type `T`.
-   **pageInfo**: This field contains information regarding the page that is returned, so that the front-end / client knows its current position, how many data that is retrieved, and maximum page index that can be requested.

[[↑ Back to top ↑](#-go-gin-template)]

---
