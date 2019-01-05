package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nuveo/anticaptcha"
)

type createAccount struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
	RecaptchaResponse string `json:"recaptcha_response"`
}

type apiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type oAuth struct {
	ClientID  string `json:"client_id"`
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Scope     string `json:"scope"`
}

type oAuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type getUUID struct {
	ID string `json:"id"`
}

type subscription struct {
	Sub struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Coupon    string `json:"coupon_code"`
		Plan      string `json:"plan_code"`
		Token     string `json:"recurly_token"`
	} `json:"subscription"`
}

type subData struct {
	UUID string `json:"uuid"`
}

func generateRTAccount() (string, string, error) {
	siteKey := "6LeZAyAUAAAAAKXhHLkm7QSka-pPFSRLgL7fjS_g"
	url := "https://roosterteeth.com/signup"
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	user := make([]rune, 10)
	pass := make([]rune, 10)
	for i := range user {
		user[i] = letters[rand.Intn(len(letters))]
		pass[i] = letters[rand.Intn(len(letters))]
	}
	email := string(user) + "@how2trianglemuygud.com"
	password := string(pass)
	client := &anticaptcha.Client{APIKey: data.Anticaptcha}
	key, err := client.SendRecaptcha(url, siteKey)
	if err != nil {
		return email, password, err
	}
	var post createAccount
	post.User.Email = email
	post.User.Password = password
	post.RecaptchaResponse = key
	JSON, _ := json.Marshal(post)
	resp, _ := httpPostJSON("https://business-service.roosterteeth.com/api/v1/users", nil, JSON)
	var Resp apiError
	json.Unmarshal([]byte(resp), &Resp)
	if Resp.Error != "" {
		return email, password, errors.New(Resp.Message)
	}
	return email, password, nil
}

func rtAuthenticate(email string, password string) string {
	var post oAuth
	post.ClientID = "4338d2b4bdc8db1239360f28e72f0d9ddb1fd01e7a38fbb07b4b1f4ba4564cc5"
	post.GrantType = "password"
	post.Username = email
	post.Password = password
	post.Scope = "user public"
	JSON, _ := json.Marshal(post)
	resp, _ := httpPostJSON("https://auth.roosterteeth.com/oauth/token", nil, JSON)
	var RESP oAuthResp
	json.Unmarshal([]byte(resp), &RESP)
	return RESP.TokenType + " " + RESP.AccessToken
}

func rtActivateFirst(token string) {
	headers := [][]string{[]string{"Authorization", token}}
	resp, _ := httpGet("https://business-service.roosterteeth.com/api/v1/me", headers)
	var RESP getUUID
	json.Unmarshal([]byte(resp), &RESP)
	UUID := RESP.ID
	respo, _ := http.PostForm("https://api.recurly.com/js/v1/token",
		url.Values{"first_name": {data.FName}, "last_name": {data.LName}, "postal_code": {data.Postcode}, "number": {data.CNum}, "month": {data.CMon}, "year": {data.CYea}, "cvv": {data.CCVV}, "version": {"4.9.3"}, "key": {"ewr1-2beFfL1PHAOpBH03tu5h6j"}})
	var RECURLY getUUID
	body, _ := ioutil.ReadAll(respo.Body)
	json.Unmarshal(body, &RECURLY)
	url := "https://business-service.roosterteeth.com/api/v1/recurly_service/accounts/" + UUID + "/subscriptions"
	var sub subscription
	sub.Sub.Coupon = "NEEDTOREPLACE"
	sub.Sub.FirstName = data.FName
	sub.Sub.LastName = data.LName
	sub.Sub.Plan = "1month"
	sub.Sub.Token = RECURLY.ID
	JSON, _ := json.Marshal(sub)
	jsonStr := strings.Replace(string(JSON), "\"NEEDTOREPLACE\"", "null", 1)
	JSON = []byte(jsonStr)
	r, _ := httpPostJSON(url, headers, JSON)
	var subdata subData
	json.Unmarshal([]byte(r), &subdata)
	req, _ := http.NewRequest("DELETE", "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+subdata.UUID+"/cancel", nil)
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	respo, _ = client.Do(req)
	strl, _ := ioutil.ReadAll(respo.Body)
	fmt.Println(string(strl))
}

func httpGet(url string, headers [][]string) (response string, code int) {
	req, _ := http.NewRequest("GET", url, nil)
	for i := range headers {
		req.Header.Set(headers[i][0], headers[i][1])
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.StatusCode
}

func httpPostJSON(url string, headers [][]string, json []byte) (response string, code int) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	for i := range headers {
		req.Header.Set(headers[i][0], headers[i][1])
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), resp.StatusCode
}
