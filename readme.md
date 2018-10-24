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
**TODO**