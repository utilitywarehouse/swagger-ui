# Swagger UI

`go get github.com/utilitywarehouse/swagger-ui`

swaggerui returns a handler with 3 routes

`/` redirects to `/swagger-ui`

`/swagger.json` presents the swagger file. This is looked for in the root of your project, and if you are using `go run` to test your code, because the root is dynamic you must store a swagger.json file in /tmp/swagger.json.

`/swagger-ui` is the UI itself
