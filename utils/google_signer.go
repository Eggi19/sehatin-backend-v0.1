package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
)

type GoogleSigner interface {
	RetrieveToken(authCode string) (*dtos.GoogleAuthToken, error)
	RetrieveUserData(token dtos.GoogleAuthToken) (*dtos.GoogleAuthUserResponse, error)
}

type googleSigner struct {
	config Config
}

type googleAuthParam struct {
	method      string
	url         string
	headerKey   string
	headerValue string
}

type genders struct {
	FormattedValue string `json:"formattedValue"`
}

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type birthDate struct {
	Date date `json:"date"`
}

type googlePeopleData struct {
	Genders   []genders   `json:"genders"`
	Birthdays []birthDate `json:"birthdays"`
}

func NewGoogleSigner(config Config) *googleSigner {
	return &googleSigner{
		config: config,
	}
}

func (g *googleSigner) RetrieveToken(authCode string) (*dtos.GoogleAuthToken, error) {
	param := url.Values{}
	param.Add("code", authCode)
	param.Add("client_id", g.config.GoogleId)
	param.Add("client_secret", g.config.GoogleKey)
	param.Add("redirect_uri", g.config.GoogleUri)
	param.Add("grant_type", "authorization_code")
	urlParam := param.Encode()

	res, err := g.authApiResponse(googleAuthParam{
		url:    "https://oauth2.googleapis.com/token?" + urlParam,
		method: "POST", headerKey: "Content-Type",
		headerValue: "application/x-www-form-urlencoded",
	})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var authToken dtos.GoogleAuthToken
	err = json.NewDecoder(res.Body).Decode(&authToken)
	if err != nil {
		return nil, err
	}

	return &authToken, err
}

func (g *googleSigner) RetrieveUserData(token dtos.GoogleAuthToken) (*dtos.GoogleAuthUserResponse, error) {
	res, err := g.authApiResponse(googleAuthParam{
		url:         "https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" + token.AccessToken,
		headerKey:   "Authorization",
		method:      "GET",
		headerValue: "Bearer" + token.IdToken,
	})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var userData dtos.GoogleAuthUserResponse
	err = json.NewDecoder(res.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	res, err = g.authApiResponse(googleAuthParam{
		url:         "https://people.googleapis.com/v1/people/me?personFields=genders%2Cbirthdays&access_token=" + token.AccessToken,
		headerKey:   "Authorization",
		method:      "GET",
		headerValue: "Bearer" + token.IdToken,
	})
	if err != nil {
		return nil, err
	}

	var peopleData googlePeopleData
	err = json.NewDecoder(res.Body).Decode(&peopleData)
	if err != nil {
		return nil, err
	}

	genders := peopleData.Genders
	if len(genders) > 0 {
		userData.Gender = &peopleData.Genders[0].FormattedValue
	}

	birthDate := peopleData.Birthdays
	if len(birthDate) > 0 {
		birth := peopleData.Birthdays[0].Date
		if birth.Year < 20 {
			birth.Year += 2000
		}
		mergedDate := time.Date(birth.Year, time.Month(birth.Month), birth.Day, 1, 0, 0, 0, time.UTC).Format("2006-01-02")
		userData.BirthDate = &mergedDate
	}

	return &userData, nil
}

func (g *googleSigner) authApiResponse(param googleAuthParam) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(param.method, param.url, nil)
	req.Header.Set(param.headerKey, param.headerValue)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
