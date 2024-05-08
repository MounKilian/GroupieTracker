package groupieTracker

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func SpotifyAPI(w http.ResponseWriter, r *http.Request) {
	authConfig := &clientcredentials.Config{
		ClientID:     "2a8a0128c5aa4458b24fc07d90d76135",
		ClientSecret: "c0d7e68a34b04b88ae577d71163ab073",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		http.Error(w, "error retrieve access token", http.StatusInternalServerError)
		log.Fatalf("error retrieve access token: %v", err)
		return
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlistID := spotify.ID("7lIjK7AbxSFO0kB9NctldT")
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		http.Error(w, "error retrieve playlist data", http.StatusInternalServerError)
		log.Fatalf("error retrieve playlist data: %v", err)
		return
	}

	randTrack := rand.Intn(playlist.Tracks.Total)

	log.Println("playlist id:", playlist.ID)
	log.Println("playlist name:", playlist.Name)
	log.Println("playlist lenght:", playlist.Tracks.Total)
	log.Println("playlist track:", playlist.Tracks.Tracks[randTrack].Track)
	log.Println("playlist track:", playlist.Tracks.Tracks[randTrack].Track.PreviewURL)
}
