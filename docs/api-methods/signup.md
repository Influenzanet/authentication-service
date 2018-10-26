**Signup**
----
  Register a new user account with the given email address and password, if they match the validation criterias (valid email format and password at least 6 characters including letters and numbers).

* **URL**

  /v1/user/signup

* **Method:**

  `POST`

*  **URL Params**
  * None

* **Data Params**
  **Required:**
  * **Type:** application/json <br />
    **Content:** `{ "email": "<user email>", "password": "<user password>"}`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{ "token": "<token string>", "role": "PARTICIPANT" }`

* **Error Response:**

  * **Code:** 400 Bad request <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Data format (json body of the Post request) wrong, e.g. missing key for email or password.

  * **Code:** 400 Bad request <br />
    **Content:** `{ "error" : "email address already in use" }` <br />
    **Typical reason:** Email address already used for an other account.

  * **Code:** 500 Internal server error <br />
    **Content:** `{ "error" : "<error message>" }` <br />
    **Typical reason:** Something went wrong during the token generation. User's input are ok, but method failed generating a valid token, e.g. because signing key is not available.

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
  * None
