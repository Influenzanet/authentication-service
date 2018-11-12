**Login**
----
  Check email and password for a user, check if the user has the required role, and return a token if credentials are correct.

* **URL**

  /v1/user/login

* **Method:**

  `POST`

*  **URL Params:**
   * None

* **Data Params:**
  * **Type:** application/json <br />
    **Content:** `{ "email": "<user email>", "password": "<user password>", "?role": "PARTICIPANT" | "RESEARCHER" | "ADMIN" }`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{ "token": "<token string>", "role": "PARTICIPANT" | "RESEARCHER" | "ADMIN" }`

* **Error Response:**

  * **Code:** 400 Bad request <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Data format (json body of the Post request) wrong, e.g. missing key for email or password.

  * **Code:** 401 Unauthorized <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Email or password wrong or doesn't belong to any registered participant.

  * **Code:** 403 Forbidden <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** The account does not have the role for which authorization was requested.

  * **Code:** 500 Internal server error <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Something went wrong during the token generation. User's credentials are ok, but method failed generating a valid token, e.g. because signing key is not available.

* **Sample Call:**

  ```go
     creds := &userCredentials{
      Email:    "your@email.com", // `json:"email"`
      Password: "yourpassword", // `json:"password"`
      Role:     "ADMIN", // `json:"role"`
    }
    payload, err := json.Marshal(creds)
    resp, err := http.Post(auth-service-addr + "/v1/user/login", "application/json", bytes.NewBuffer(payload))
    defer resp.Body.Close()
  ```

* **Notes:**
  * If not specified `role` defaults to `"PARTICIPANT"`.
