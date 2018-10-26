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
