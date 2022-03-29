install:
	which swagger || GO111MODULE=off go get github.com/go-swagger/go-swagger/cmd/swagger
swagger:
	swagger generate spec -o ./docs/swagger.yaml --scan-models