package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	email := string(user) + "@thissitedoesntactuallyexist.xyz"
	password := string(pass)
	fmt.Println("Debug: Generating Account..")
	fmt.Println("Debug: Using credentials - Email: " + email + " Password: " + password)
	client := &anticaptcha.Client{APIKey: config.AnticaptchaToken}
	fmt.Println("Debug: Starting anticaptcha client, timeout of 3 Mins")
	key, err := client.SendRecaptcha(url, siteKey, time.Minute*3)
	if err != nil {
		return email, password, err
	}
	fmt.Println("Debug: Got recaptcha response of " + key)
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
	fmt.Println("Debug: Account generated successfully")
	return email, password, nil
}

func rtActivateFirst(token string) error {
	headers := [][]string{[]string{"Authorization", token}}
	resp, _ := httpGet("https://business-service.roosterteeth.com/api/v1/me", headers)
	var RESP getUUID
	json.Unmarshal([]byte(resp), &RESP)
	UUID := RESP.ID
	respo, _ := http.PostForm("https://api.recurly.com/js/v1/token", url.Values{"first_name": {config.CardInfo.FName}, "last_name": {config.CardInfo.LName}, "postal_code": {config.CardInfo.PCode}, "number": {config.CardInfo.Num}, "month": {config.CardInfo.Mon}, "year": {config.CardInfo.Yea}, "cvv": {config.CardInfo.CVV}, "version": {"4.14.0"}, "key": {"ewr1-2beFfL1PHAOpBH03tu5h6j"}})
	var recurly getUUID
	body, _ := ioutil.ReadAll(respo.Body)
	json.Unmarshal(body, &recurly)
	url := "https://business-service.roosterteeth.com/api/v1/recurly_service/accounts/" + UUID + "/subscriptions"
	var sub subscription
	sub.Sub.Coupon = nil
	sub.Sub.FirstName = config.CardInfo.FName
	sub.Sub.LastName = config.CardInfo.LName
	sub.Sub.Plan = "1month"
	sub.Sub.Token = recurly.ID
	JSON, _ := json.Marshal(sub)
	r, code := httpPostJSON(url, headers, JSON)
	if code != 201 {
		return errors.New("Failed to start FIRST trial, " + r)
	}
	var subdata subData
	json.Unmarshal([]byte(r), &subdata)
	req, _ := http.NewRequest("DELETE", "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+subdata.UUID+"/cancel", nil)
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	client.Do(req)
	return nil
}

func rtAuthenticate(email string, password string) (string, error) {
	var post oAuth
	post.ClientID = "4338d2b4bdc8db1239360f28e72f0d9ddb1fd01e7a38fbb07b4b1f4ba4564cc5"
	post.GrantType = "password"
	post.Username = email
	post.Password = password
	post.Scope = "user public"
	JSON, _ := json.Marshal(post)
	resp, code := httpPostJSON("https://auth.roosterteeth.com/oauth/token", nil, JSON)
	if code != 200 {
		return "", errors.New("Failed to authenticate with rt")
	}
	var RESP oAuthResp
	json.Unmarshal([]byte(resp), &RESP)
	return RESP.TokenType + " " + RESP.AccessToken, nil
}

func rtGrabLatestEpisodeInfo() *epInfo {
	resp, _ := httpGet("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-8/episodes?order=des&per_page=1", nil)
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

func rtGrabLatestEpisode(email string, password string) (magicShort string, magicLong string, err error) {
	fmt.Println(email, password)
	ep := rtGrabLatestEpisodeInfo()
	if ep == nil {
		err = errors.New("Failed to grab latest episode")
		return
	}
	token, err := rtAuthenticate(email, password)
	if err != nil {
		return
	}
	headers := [][]string{[]string{"Authorization", token}}
	url := "https://svod-be.roosterteeth.com/api/v1/episodes/" + ep.UUID + "/videos/"
	resp, code := httpGet(url, headers)
	if code != 200 {
		err = errors.New("Could not get episode stream")
		return
	}
	var vidData epVidData
	json.Unmarshal([]byte(resp), &vidData)
	txt := strings.TrimPrefix(vidData.Data[0].Attributes.URL, "https://rtv3-roosterteeth.akamaized.net/store/")
	r, _ := regexp.Compile("/ts/*.+")
	end := r.FindString(txt)
	txt = strings.TrimSuffix(txt, end)
	magic := strings.Split(txt, "-")
	if len(magic) < 2 {
		err = errors.New("Failed to get tokens")
		return
	}
	return magic[1], magic[0], nil
}
