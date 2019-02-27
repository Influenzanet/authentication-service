# Authorization Service

## Documentation
[API documentation](./docs/api.md)
[Workflows of JWT handling](./docs/jwt-token-handling.md)

## Test

Some test methods are longer to finish, e.g. token renewal has to wait for token expiration intervals. To skip this long test methods, call:
```sh
go test -short
```

To perform a full test, call:
```sh
go test
```

The tests use the gomock library. To install this use:
```
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen
```

Then generate mock client for the user management service:
```
mockgen github.com/influenzanet/api/dist/go/user-management UserManagementApiClient > mock_user_management/um_mock.go
```
For more information about testing grpc clients with go check: https://github.com/grpc/grpc-go/blob/master/Documentation/gomock-example.md

## Build

The included `Dockerfile` should provide everything needed to build and run the application.

Build the image:
```sh
docker build --rm -t 'authentication-service:latest' .
```

Run the image:
```sh
docker run -p 3100:3100 authentication-service
```

Access the API on `localhost:3100`
