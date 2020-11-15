# message-board-backend
Back end for a message board application, written in Go with PostgreSQL and Docker. An instance of this application will be running on AWS for the month of November at http://gabrielb.eu-central-1.elasticbeanstalk.com

## Dependencies
This application requires Docker Compose and Docker

## Usage
* run locally : `docker-compose up` - may require sudo privledges. Starts a docker container running Postgres, listening on port 5432, and a docker container running an HTTP web server listening on port 8080.

* run unit tests : `docker-compose -f docker-compose.test.yml up` - This will run the tests located in server/app/app_test.go

* run integration tests : `python3 test/apitest.py` - By default, tests assume the server is running on http://localhost:8080. To specify a different URL, include it as the first command line argument - ex: `python3 test/apitest.py http://gabrielb.eu-central-1.elasticbeanstalk.com`

# API Documentation
The message board backend exposes a public and a private RESTful API. Both APIs accept JSON-encoded request bodies, and return JSON-encoded responses, using standard HTTP response codes and verbs.

## Authentication
HTTP basic authentication is used to validate HTTP requests.

---

## Create a new message
**URL** : `/messages`

**Method** : `POST`

**Auth required** : NO

**Data constraints**

Provide name, email, and text of message to be created:

```json
{
    "name": "[unicode]",
    "email": "[unicode]",
    "text": "[unicode]"
}
```

### Success Response
**Condition** : Everything is OK, message created successfully

**Code** : `201 CREATED`

**Content example**

```json
{
    "id": "26332274-277e-11eb-a10f-02f45a9cb7ee",
    "name": "gabe",
    "email": "test@gmail.com",
    "text": "example text",
    "creation_time": "2020-11-15T20:07:12.509223Z"
}
```

### Error Response
**Condition** : JSON message was malformed

**Code** : `400 BAD REQUEST`

---
## List messages, ordered anti-chronologically
**URL** : `/messages`

**Method** : `GET`

**Auth required** : YES

**Data constraints** : None

### Success Response
**Condition** : 0 or more messages retrieved

**Code** : `200 OK`

**Content example**

```json
[
    {
        "id": "26332275-277e-11eb-a10f-02f45a9cb7ee",
        "name": "gabe",
        "email": "test@gmail.com",
        "text": "example text",
        "creation_time": "2020-11-15T20:07:12.509223Z"
    },
    {
        "id": "26332274-277e-11eb-a10f-02f45a9cb7ee",
        "name": "gabe",
        "email": "test@gmail.com",
        "text": "another message",
        "creation_time": "2020-10-15T20:07:12.509223Z"
    },
    {
        "id": "26332271-277e-11eb-a10f-02f45a9cb7ee",
        "name": "gabe",
        "email": "test@gmail.com",
        "text": "yet anothe message",
        "creation_time": "2019-11-15T20:07:12.509223Z"
    },
]
```

### Error Response
**Condition** : Incorrect Authorization

**Code** : `401 UNAUTHORIZED`

---
## View message by ID
**URL** : `/messages/{id}`

**Method** : `GET`

**Auth required** : YES

**Data constraints** : None

### Success Response
**Condition** : Message retrieved

**Code** : `200 OK`

**Content example**

```json
{
    "id": "26332274-277e-11eb-a10f-02f45a9cb7ee",
    "name": "gabe",
    "email": "test@gmail.com",
    "text": "example text",
    "creation_time": "2020-11-15T20:07:12.509223Z"
}
```

### Error Response
**Condition** : Incorrect Authorization

**Code** : `401 UNAUTHORIZED`

**Condition** : ID not found in database

**Code** : `404 Not Found`

**Condition** : UUID is malformed

**Code** : `500 Internal Server Error`

---
## Update message text
**URL** : `/messages/{id}`

**Method** : `PUT`

**Auth required** : YES

**Data constraints**

Provide new text:

```json
{
    "text": "[unicode]"
}
```

### Success Response
**Condition** : Message updated successfully

**Code** : `200 OK`

**Content example**

```json
{
    "id": "26332274-277e-11eb-a10f-02f45a9cb7ee",
    "name": "gabe",
    "email": "test@gmail.com",
    "text": "updated text",
    "creation_time": "2020-11-15T20:07:12.509223Z"
}
```

### Error Response
**Condition** : Incorrect Authorization

**Code** : `401 UNAUTHORIZED`

**Condition** : JSON message was malformed

**Code** : `400 BAD REQUEST`

**Condition** : ID not found in database

**Code** : `404 Not Found`

**Condition** : UUID is malformed

**Code** : `500 Internal Server Error`
