# Authorization Service

## Test
### Setup
The tests use the gomock library. To install this use:

```sh
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen
```

Then generate mock client for the user management service:

```sh
mockgen -source=./api/user-management-api.pb.go UserManagementApiClient > mocks/user-management.go
```

For more information about testing grpc clients with go check: <https://github.com/grpc/grpc-go/blob/master/Documentation/gomock-example.md>


### Script to run tests
After installing the dependencies, you can add a script to initiate tests. This script is not included in this repository, since it contains secret infos, like DB password.

```sh
export JWT_TOKEN_KEY="<insert secret key to sign jwt>"
export TOKEN_EXPIRATION_MIN="10" # minutes token is valid
export TOKEN_MINIMUM_AGE_MIN="2"  # wait so many minutes at least before refreshing the token

export DB_CONNECTION_STR="<mongo db connection string without prefix and auth infos>"
export DB_USERNAME="<username for mongodb auth>"
export DB_PASSWORD="<password for mongodb auth>"
export DB_PREFIX="+srv" # e.g. "" (emtpy) or "+srv"
export DB_TIMEOUT=30 # seconds until connection times out
export DB_IDLE_CONN_TIMEOUT=45 # terminate idle connection after seconds
export DB_MAX_POOL_SIZE=8
export DB_DB_NAME_PREFIX="<DB_PREFIX>" # DB names will be then > <DB_PREFIX>+"hard-coded-db-name-as-we-need-it"

go test  ./...
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

## Token signiture key

You can use the key generator, to generate a random private key to sign tokens.
