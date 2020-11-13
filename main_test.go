// main_test.go

package main

import (
    "testing"
    "os"
//    "os/exec"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
//	"time"
//	"fmt"
)

var a App

func TestMain(m *testing.M) {
	//c := exec.Command("sudo", "docker", "run", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=postgres", "--name", "test_db", "-d", "postgres")
	//err := c.Run()
	/*
	if err != nil{
		panic(err)
	}
		time.Sleep(2 *time.Second)*/
	a.Initialize("localhost", 5432, "postgres", "postgres", "postgres")

	ensureTableExists(a.DB)
	code := m.Run()
	/*
	stop := exec.Command("sudo", "docker", "stop", "test_db")
	rm_container := exec.Command("sudo", "docker", "rm", "test_db")

	fmt.Println("stopping database")
	if err := stop.Run(); err != nil{
		fmt.Println("stop failed")
	}
	fmt.Println("removing database")
	if err := rm_container.Run(); err != nil{
		fmt.Println("removing docker container failed")
	}*/

	clearTable(a.DB)
	os.Exit(code)
}



/*
//TestLoadDataFromCSV tests the loadDataFromCSV function in model.go
//It does this by loading sample_data.csv into the database and checking that the database isn't empty
//If /sample_data.csv is not available to the database, this test will fail
func TestLoadDataFromCSV(t *testing.T) {
	clearTable(a.DB)
	err := loadDataFromCSV("/sample_data.csv", a.DB)
	if err != nil {
		log.Fatal(err)
	}

	//test database has messages
	req, _ := http.NewRequest("GET", "/messages", nil)
	req.SetBasicAuth("admin","back-challenge")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Messages not loaded into DB")
	}
}
*/

func TestGetMessage(t *testing.T) {
	clearTable(a.DB)

	//create valid message
	var jsonStr = []byte(`{"name":"Justin", "email": "JB37@gmail.com", "text": "Test Message"}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var originalMessage message;
	json.Unmarshal(response.Body.Bytes(), &originalMessage)

	//now make a get request for the same message
	req, _ = http.NewRequest("GET", "/messages/"+originalMessage.ID, nil)
	req.SetBasicAuth("admin","back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var returnedMessage message
	json.Unmarshal(response.Body.Bytes(), &returnedMessage)

	if originalMessage.ID != returnedMessage.ID {
		t.Errorf("Expected message ID to be '%v'. Got '%v'", originalMessage.ID, returnedMessage.ID)
	}
	if originalMessage.Name != returnedMessage.Name {
		t.Errorf("Expected message Name to be '%v'. Got '%v'", originalMessage.Name, returnedMessage.Name)
	}
	if originalMessage.Email != returnedMessage.Email {
		t.Errorf("Expected message Email to be '%v'. Got '%v'", originalMessage.Email, returnedMessage.Email)
	}
	if originalMessage.Text != returnedMessage.Text {
		t.Errorf("Expected message Text to be '%v'. Got '%v'", originalMessage.Text, returnedMessage.Text)
	}
	if originalMessage.TimeStamp != returnedMessage.TimeStamp {
		t.Errorf("Expected message Timestamp to be '%v'. Got '%v'", originalMessage.TimeStamp, returnedMessage.TimeStamp)
	}
}

//tests that unauthenticated requests cannot recieve info from private API
func TestAuth(t *testing.T) {
	clearTable(a.DB)

	//create valid message
	var jsonStr = []byte(`{"name":"Justin", "email": "JB37@gmail.com", "text": "Test Message"}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var originalMessage message;
	json.Unmarshal(response.Body.Bytes(), &originalMessage)

	//now try to retrieve message without authentication
	req, _ = http.NewRequest("GET", "/messages/"+originalMessage.ID, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	//try to retrieve messages without auth
	req2, _ := http.NewRequest("GET", "/messages", nil)
	var response2 = executeRequest(req2)
	checkResponseCode(t, http.StatusUnauthorized, response2.Code)

	//try to modify a message without auth
	jsonStr = []byte(`{"name":"New-name", "email": "New-email", "text": "New-text"}`)
	req3, _ := http.NewRequest("PUT", "/messages/"+originalMessage.ID, nil)
	var response3 = executeRequest(req3)
	checkResponseCode(t, http.StatusUnauthorized, response3.Code)
}

func TestUpdateMessage(t *testing.T) {
	clearTable(a.DB)

	//create valid message
	var jsonStr = []byte(`{"name":"Justin", "email": "JB37@gmail.com", "text": "Test Message"}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var originalMessage message;
	json.Unmarshal(response.Body.Bytes(), &originalMessage)

	//try to modify message
	jsonStr = []byte(`{"name":"New-name", "email": "New-email", "text": "New-text"}`)
	req, _ = http.NewRequest("PUT", "/messages/"+originalMessage.ID, nil)
	req.SetBasicAuth("admin","back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	//now make a get request for the same message
	req, _ = http.NewRequest("GET", "/messages/"+originalMessage.ID, nil)
	req.SetBasicAuth("admin","back-challenge")
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var returnedMessage message
	json.Unmarshal(response.Body.Bytes(), &returnedMessage)

	if originalMessage.ID != returnedMessage.ID {
		t.Errorf("Expected message ID to be '%v'. Got '%v'", originalMessage.ID, returnedMessage.ID)
	}
	if originalMessage.Name == returnedMessage.Name {
		t.Errorf("Expected message Name to change but it stayed the same")
	}
	if originalMessage.Email == returnedMessage.Email {
		t.Errorf("Expected message Email to change but it stayed the same")
	}
	if originalMessage.Text == returnedMessage.Text {
		t.Errorf("Expected message Text to change but it stayed the same")
	}
	if originalMessage.TimeStamp != returnedMessage.TimeStamp {
		t.Errorf("Expected message Timestamp to stay the same but it changed")
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable(a.DB)

	req, _ := http.NewRequest("GET", "/messages", nil)
	req.SetBasicAuth("admin","back-challenge")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateMessage(t *testing.T) {
	clearTable(a.DB)

	//test creating valid message
	var jsonStr = []byte(`{"name":"Justin", "email": "JB37@gmail.com", "text": "Test Message"}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Justin" {
		t.Errorf("Expected message name to be 'Justin'. Got '%v'", m["name"])
	}

	if m["email"] != "JB37@gmail.com" {
		t.Errorf("Expected email to be 'JB37@gmail.com'. Got '%v'", m["email"])
	}

	if m["text"] != "Test Message" {
		t.Errorf("Expected text to be 'Test Message'. Got '%v'", m["text"])
	}

}

func TestCreateInvalidMessage(t *testing.T) {
	clearTable(a.DB)

	//test creating invalid message
	var jsonStr = []byte(`{"name":"Gabe", "text": "Forgot email!"}`)
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

