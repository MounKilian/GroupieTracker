package groupieTracker

import (
	"context"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type Btest struct {
	name       string
	PreviewURL string
}

// Use spotify API and connect to a specific playlist
func PlaylistConnectBT(genderPlaylist string) spotify.FullTrack {

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
	playlistID := spotify.ID(genderPlaylist)
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	var currentMusic Music
	randomIndex := GetRandomMusicIndex(playlist)
	currentMusic.name = playlist.Tracks.Tracks[randomIndex].Track.Name

	return playlist.Tracks.Tracks[randomIndex].Track
}

func BlindtestManager(genderPlaylist string) Btest {
	randomBtest := PlaylistConnectBT(genderPlaylist)

	var currentBtest Btest
	currentBtest.name = randomBtest.Name
	currentBtest.PreviewURL = randomBtest.PreviewURL

	return currentBtest
}

func GetPreviewURL(track *spotify.FullTrack) string {
	return track.PreviewURL
}
