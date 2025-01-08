package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

var clientId string = os.Getenv("CLIENT_ID")
var clientSecret string = os.Getenv("CLIENT_SECRET")

const redirectUri string = "http://localhost:8000/callback"
const authUrl string = "https://accounts.spotify.com/authorize?"

var s *http.Server

type tokenResponse struct {
	RefreshToken string `json:"refresh_token"`
}

func generateRandomString(l int) string {
	const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, possible[rand.Intn(len(possible))])
	}
	return string(result)
}

func getLoginUrl() string {
	state := generateRandomString(16)
	scope := "user-read-email playlist-read-private playlist-modify-public user-library-read"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientId)
	params.Add("scope", scope)
	params.Add("redirect_uri", redirectUri)
	params.Add("state", state)
	return authUrl + params.Encode()
}

func callback(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	code := r.URL.Query().Get("code")
	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", redirectUri)
	data.Set("grant_type", "authorization_code")
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response tokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Refresh token: %v\n", response.RefreshToken)

	s.Close()
}

func main() {
	err := exec.Command("xdg-open", getLoginUrl()).Start()
	if err != nil {
		log.Fatalf("%s", err)
	}

	http.HandleFunc("/callback", callback)

	s = &http.Server{
		Addr:    ":8000",
		Handler: http.DefaultServeMux,

		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err = s.ListenAndServe(); err != nil {
		log.Fatalf("%s", err)
	}
}
