**Signup Researcher**
----
  Start registration/promotion process for a new researcher account with the given email address and password. **TODO** registration process still to be defined

* **URL**

  /v1/signup/researcher

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

  OR

  * **Code:** 400 Bad request <br />
    **Content:** `{ "error" : "email address already in use" }` <br />
    **Typical reason:** Email address already used for an other account.

  OR

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
  None