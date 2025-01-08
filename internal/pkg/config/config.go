package config

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	refreshToken string
	clientId     string
	clientSecret string
	AccessToken  string `json:"access_token"`
	SrcUri       string
	DstUri       string
}

func (c *Config) getEnvVariables() {
	c.refreshToken = os.Getenv("REFRESH_TOKEN")
	c.clientId = os.Getenv("CLIENT_ID")
	c.clientSecret = os.Getenv("CLIENT_SECRET")
	c.SrcUri = os.Getenv("SRC_URI")
	c.DstUri = os.Getenv("DST_URI")
}

func (c *Config) getAccessToken() {
	if c.refreshToken == "" {
		log.Fatal("Missing refresh token")
	}
	if c.clientId == "" {
		log.Fatal("Missing clientID")
	}
	if c.clientSecret == "" {
		log.Fatal("Missing client secret")
	}

	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.refreshToken)
	data.Set("client_id", c.clientId)
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.clientId+":"+c.clientSecret)))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateConfig() Config {
	c := Config{}
	c.getEnvVariables()
	c.getAccessToken()
	return c
}
