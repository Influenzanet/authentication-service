**Login Researcher**
----
  Check email and password for logging in as a researcher and return a token if credentials are correct. User has to be in the researcher list to be able to sign in as a researcher. The role 'researcher' is written into the token.

* **URL**

  /v1/login/researcher

* **Method:**

  `POST`

*  **URL Params**
  None

* **Data Params**
  **Required:**
  * **Type:** application/json <br />
    **Content:** `{ "email": "<user email>", "password": "<user password>"}`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{ "token": "<token string>" }`

* **Error Response:**

  * **Code:** 400 Bad request <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Data format (json body of the Post request) wrong, e.g. missing key for email or password.

  * **Code:** 401 Unauthorized <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Email or password wrong or doesn't belong to any registered participant.

  * **Code:** 403 Forbidden <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** The account does not meet the required authorization level.

  * **Code:** 500 Internal server error <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Something went wrong during the token generation. User's credentials are ok, but method failed generating a valid token, e.g. because signing key is not available.

* **Sample Call:**
  TODO: add sample call for go

  ```javascript
    $.ajax({
      url: "/users/1",
      dataType: "json",
      type : "GET",
      success : function(r) {
        console.log(r);
      }
    });
  ```
* **Notes:**
  None
