package playlist

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetPlaylist(next string, accessToken string, c chan Playlist) {
	currentPage := next
	p := Playlist{Next: next, Tracks: map[string]simpleTrack{}}

	client := &http.Client{}

	for p.Next != "" {
		req, _ := http.NewRequest("GET", p.Next, nil)
		req.Header.Add("Authorization", "Bearer "+accessToken)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &p)

		if err != nil {
			log.Fatal(err)
		}

		p.cleanBuffer()
		if currentPage == p.Next {
			break
		}
		currentPage = p.Next
	}
	c <- p
}

func ComparePlaylists(src Playlist, dst Playlist) map[string]simpleTrack {
	missingTracks := map[string]simpleTrack{}
	for k, v := range src.Tracks {
		if _, ok := dst.Tracks[k]; !ok {
			log.Printf("Missing track: %s - %s", v.Name, v.Artists)
			missingTracks[k] = v
		}
	}
	return missingTracks
}

func AddTracksToPlaylist(endpoint string, missingTracks map[string]simpleTrack, accessToken string) {
	batches := [][]string{}
	buffer := []string{}
	for _, v := range missingTracks {
		buffer = append(buffer, v.Uri)
		if len(buffer) == 100 {
			batches = append(batches, buffer)
			buffer = []string{}
		}
	}
	if len(buffer) != 0 {
		batches = append(batches, buffer)
	}
	client := http.Client{}

	for _, v := range batches {
		data, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		body := []byte(`{"uris":` + string(data) + "}")
		req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessToken)
		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}
}
