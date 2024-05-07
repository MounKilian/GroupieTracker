package groupieTracker

import (
	"context"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func SpotifyAPI() {
	authConfig := &clientcredentials.Config{
		ClientID:     "2a8a0128c5aa4458b24fc07d90d76135",
		ClientSecret: "c0d7e68a34b04b88ae577d71163ab073",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	// trackID := spotify.ID("2tUfesk6jwISYkygWV9481")

	// err = client.PlayOpt(&spotify.PlayOptions{
	// 	URIs: []spotify.URI{spotify.URI("spotify:track:" + trackID)},
	// })
	// if err != nil {
	// 	log.Fatalf("error playing track: %v", err)
	// }

	playlistID := spotify.ID("7lIjK7AbxSFO0kB9NctldT")
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	log.Println("playlist id:", playlist.ID)
	log.Println("playlist name:", playlist.Name)
	log.Println("tracks name:", playlist.Tracks.Tracks[3].Track)
	log.Println("tracks name:", playlist.Tracks.Tracks[3])
	log.Println("playlist description:", playlist.Description)

}
