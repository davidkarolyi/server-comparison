package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		testData := readTestData()

		payload := verifyToken(testData.Token, testData.Secret)
		marshaledPayload, _ := json.Marshal(payload)

		res.Header().Add("Content-Type", "application/json")
		res.Write(marshaledPayload)
	})

	fmt.Println("go_std is listening on localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func readTestData() TestData {
	jsonTestData, _ := ioutil.ReadFile("../test_data.json")

	testData := TestData{}
	json.Unmarshal(jsonTestData, &testData)

	return testData
}

func verifyToken(token string, secret string) Payload {
	claims := &Payload{}
	payload, _ := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	return *payload.Claims.(*Payload)
}

// Payload represents the decoded payload of a JWT token
type Payload struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
	Age    int    `json:"age"`
	Iat    int    `json:"iat"`
	jwt.StandardClaims
}

// TestData represents the test_data.json file's content
type TestData struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}
