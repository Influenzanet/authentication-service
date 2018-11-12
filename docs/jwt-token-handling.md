# Handling JSON web-tokens


## Signing and validation keys

TODO: where are keys mounted
TODO: how to generate new keys

### Token generation
The signing key will be loaded from the ??.rsa key file, each time the method is called. When the key is replaced, new tokens are simply signed with the new key.

### Token validation
Since the token has to be validated for each incoming request, that wants to access protected endpoints, the verification key is stored in memory. We can avoid frequent file reading and parsing and so save time and effort.

The verification is reloaded each time when a new token is generated to ensure that the service has the latest version. If the key has changed, the old key file will be stored for an intervall ??, to allow clients refresh their old tokens.
When the key-pair is replaced with new keys, the new tokens are signed with the new private key. Older tokens' validation will fail using the validation key stored in memory.
They will be checked, it they can be verified using the old key as long as available.

## Token expiration and refreshing tokens
