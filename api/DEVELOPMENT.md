# NimbusX api development guide

## Code generation

````bash
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0 -generate gin -package api -o internal/api/api.go api.yaml
````

## Redoc

````bash
docker run -it --rm -p 80:80   -v $(pwd)/api.yaml:/usr/share/nginx/html/swagger.yaml   -e SPEC_URL=swagger.yaml redocly/redoc
````
