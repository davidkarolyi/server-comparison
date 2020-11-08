package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func main() {
	fmt.Println("go_fasthttp is listening on localhost:3000")
	fasthttp.ListenAndServe(":3000", handler)
}

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		testData := readTestData()
		payload := verifyToken(testData.Token, testData.Secret)
		marshaledPayload, _ := json.Marshal(payload)

		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.SetContentType("application/json; charset=utf8")
		ctx.Write(marshaledPayload)
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
	}
}

func readTestData() *TestData {
	jsonTestData, _ := ioutil.ReadFile("../test_data.json")

	testData := TestData{}
	json.Unmarshal(jsonTestData, &testData)

	return &testData
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
