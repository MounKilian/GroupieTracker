package groupieTracker

import (
	"context"
	"log"
	"math/rand"

	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type Music struct {
	name   string
	lyrics string
}

func PlaylistConnect() Music {
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

	playlistID := spotify.ID("2h7KHuoD0IfVNDeQDqTGIJ")
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	var currentMusic Music
	randomIndex := GetRandomMusicIndex(playlist)
	currentMusic.name = playlist.Tracks.Tracks[randomIndex].Track.Name
	currentMusic.lyrics = GetLyrics(&playlist.Tracks.Tracks[randomIndex].Track)
	return currentMusic
}

func GetRandomMusicIndex(playlist *spotify.FullPlaylist) int {
	maxIndex := playlist.Tracks.Total
	trackIndex := rand.Intn(maxIndex-1) + 1
	return trackIndex
}

func GetLyrics(track *spotify.FullTrack) string {
	var artistName string
	for _, artist := range track.Artists {
		artistName = artist.Name
	}
	l := lyrics.New()
	lyric, err := l.Search(artistName, track.Name)

	if err != nil {
		log.Println(err)
	}

	return lyric
}
