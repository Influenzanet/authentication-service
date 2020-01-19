export JWT_TOKEN_KEY="lQyUl+2eoZ3LNMm3dojOUDcNw4/yeinDaSEWtXGdrY0="
export TOKEN_EXPIRATION_MIN="10"
export TOKEN_MINIMUM_AGE_MIN="2"

export DB_CONNECTION_STR="cluster0-vtkz6.mongodb.net/test?retryWrites=true&w=majority"
export DB_USERNAME="user-management-service"
export DB_PASSWORD="89kJAO43BRUyNSbr"
export DB_PREFIX="+srv"
export DB_TIMEOUT=30
export DB_IDLE_CONN_TIMEOUT=45
export DB_MAX_POOL_SIZE=8
export DB_DB_NAME_PREFIX="INF_"


go test

