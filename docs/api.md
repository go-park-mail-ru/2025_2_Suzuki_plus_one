# API Documentation

Run Swagger UI to explore the API endpoints and their specifications.

Entrypoint is [/api/popfilms.yaml](../api/popfilms.yaml)

## Running Swagger UI

```bash
docker pull swaggerapi/swagger-ui
# From project root directory
docker run -p 80:8080 -e SWAGGER_JSON=/api/popfilms.yaml -v $(pwd)/api:/api swaggerapi/swagger-ui
```

Access Swagger UI at: [http://localhost](http://localhost)

Don't forget to switch to local server!

> Note: Cookie auth doesn't work in Swagger UI. Use raw curl commands

## Updating API

Optimal new endpoint flow

1. openapi
2. dto/NEW
3. controller/contract NEW USECASE INTERFACE
4. controller/handler NEW HANDLER
5. controller/http/router NEW ROUTE
6. controller/handlers/NEW
7. controller/handlers/NEW-test
8. usecase/NEW
9. usecase/contract NEW ADAPTER INTERFACE
10. usecase/NEW-test
11. entity/NEW
12. adapter/NEW
13. adapter/NEW-test
