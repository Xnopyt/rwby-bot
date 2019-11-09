package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
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
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Coupon    *string `json:"coupon_code"`
		Plan      string  `json:"plan_code"`
		Token     string  `json:"recurly_token"`
	} `json:"subscription"`
}

type subData struct {
	UUID string `json:"uuid"`
}

type episodeInfo struct {
	Data []struct {
		UUID       string `json:"uuid"`
		Attributes struct {
			Title           string    `json:"title"`
			Number          int       `json:"number"`
			SponsorGoliveAt time.Time `json:"sponsor_golive_at"`
		} `json:"attributes"`
	} `json:"data"`
}

type epVidData struct {
	Data []struct {
		Attributes struct {
			URL string `json:"url"`
		} `json:"attributes"`
	} `json:"data"`
}

type epInfo struct {
	UUID   string
	Title  string
	EpNum  int
	GoLive time.Time
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
	send("Solving a reCaptcha, please wait (This may take serveral minutes)...")
	createCaptchaMessage()
	key, err := client.SendRecaptcha(url, siteKey, time.Minute*3)
	if err != nil {
		return email, password, err
	}
	finalizeCaptchaMessage()
	send("Using ```" + key + "``` as recaptcha response...")
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
	var recurly getUUID
	body, _ := ioutil.ReadAll(respo.Body)
	json.Unmarshal(body, &recurly)
	url := "https://business-service.roosterteeth.com/api/v1/recurly_service/accounts/" + UUID + "/subscriptions"
	var sub subscription
	sub.Sub.Coupon = nil
	sub.Sub.FirstName = data.FName
	sub.Sub.LastName = data.LName
	sub.Sub.Plan = "1month"
	sub.Sub.Token = recurly.ID
	JSON, _ := json.Marshal(sub)
	r, _ := httpPostJSON(url, headers, JSON)
	var subdata subData
	json.Unmarshal([]byte(r), &subdata)
	req, _ := http.NewRequest("DELETE", "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+subdata.UUID+"/cancel", nil)
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	client.Do(req)
}

func rtGrabLatestEpisodeInfo() *epInfo {
	resp, _ := httpGet("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-7/episodes?order=des&per_page=1", nil)
	var epinfo episodeInfo
	json.Unmarshal([]byte(resp), &epinfo)
	var ep epInfo
	if len(epinfo.Data) == 0 {
		return nil
	}
	ep.UUID = epinfo.Data[0].UUID
	ep.Title = epinfo.Data[0].Attributes.Title
	ep.EpNum = epinfo.Data[0].Attributes.Number
	ep.GoLive = epinfo.Data[0].Attributes.SponsorGoliveAt
	return &ep
}

func rtGrabLatestEpisode(email string, password string) (magicShort string, magicLong string) {
	ep := rtGrabLatestEpisodeInfo()
	headers := [][]string{[]string{"Authorization", rtAuthenticate(email, password)}}
	url := "https://svod-be.roosterteeth.com/api/v1/episodes/" + ep.UUID + "/videos/"
	resp, _ := httpGet(url, headers)
	var vidData epVidData
	json.Unmarshal([]byte(resp), &vidData)
	txt := strings.TrimPrefix(vidData.Data[0].Attributes.URL, "https://rtv3-roosterteeth.akamaized.net/store/")
	r, _ := regexp.Compile("/ts/*.+")
	end := r.FindString(txt)
	txt = strings.TrimSuffix(txt, end)
	magic := strings.Split(txt, "-")
	return magic[1], magic[0]
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
