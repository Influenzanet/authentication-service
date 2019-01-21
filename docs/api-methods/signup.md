# Signup

  Register a new user account with the given email address and password, if they match the validation criterias (valid email format and password at least 6 characters including letters and numbers).

* **URL**

  /v1/user/signup

* **Method:**

  `POST`

* **URL Params**

  * None

* **Data Params**

  **Required:**

  * **Type:** application/json  
    **Content:** `{ "email": "<user email>", "password": "<user password>"}`

* **Success Response:**

  * **Code:** 201  
    **Content:** `{ "token": "<token string>", "roles": [ "PARTICIPANT" ], "authenticated_role": "PARTICIPANT" }`

* **Error Response:**

  * **Code:** 400 Bad Request  
    **Content:** `{ "error" : "<error message>" }`  
    **Typical reason:** Data format (json body of the Post request) wrong, e.g. missing key for email or password.

  * **Code:** 500 Internal Server Error  
    **Content:** `{ "error" : "<error message>" }`  
    **Typical reason:** Something went wrong during the token generation. User's input are ok, but method failed generating a valid token, e.g. because signing key is not available.

  **Forwarded Responses from:**  
  [User-Management Service](https://github.com/Influenzanet/user-management-service/blob/master/docs/api-methods/signup.md)

* **Sample Call:**

  ```go
    creds := &userCredentials{
      Email:    "your@email.com", // `json:"email"`
      Password: "yourpassword", // `json:"password"`
    }
    payload, err := json.Marshal(creds)
    resp, err := http.Post(auth-service-addr + "/v1/user/signup", "application/json", bytes.NewBuffer(payload))
    defer resp.Body.Close()
  ```

* **Notes:**

  None
