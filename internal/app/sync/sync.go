package sync

import (
	"fmt"
	"log"
	"os"

	"github.com/blacha370/spotify-sync/internal/pkg/config"
	"github.com/blacha370/spotify-sync/internal/pkg/playlist"
)

func Sync() {
	cfg := config.CreateConfig()
	srcChan := make(chan playlist.Playlist)
	dstChan := make(chan playlist.Playlist)
	go playlist.GetPlaylist(fmt.Sprintf("%s?limit=50&offset=0", cfg.SrcUri), cfg.AccessToken, srcChan)
	go playlist.GetPlaylist(fmt.Sprintf("%s?limit=50&offset=0", cfg.DstUri), cfg.AccessToken, dstChan)
	missingTracks := playlist.ComparePlaylists(<-srcChan, <-dstChan)
	if len(missingTracks) == 0 {
		log.Println("No missing tracks")
		os.Exit(0)
	}
	playlist.AddTracksToPlaylist(cfg.DstUri, missingTracks, cfg.AccessToken)
}
