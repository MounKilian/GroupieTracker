package groupieTracker

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"strings"
	"unicode"

	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var Letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "T", "U", "V", "W", "X", "Y", "Z"}

func selectRandomLetter() string {
	randomIndex := rand.IntN(len(Letters) - 1)
	return Letters[randomIndex]
}

func Encrypt(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func VerifyPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}

func MatchTitle(title, input string) bool {
	if strings.Contains(title, " - ") {
		if index := strings.Index(title, " - "); index != -1 {
			title = title[:index]
		} else {
			title = ""
		}
	} else if strings.Contains(title, " (") {
		if index := strings.Index(title, " ("); index != -1 {
			title = title[:index]
		} else {
			title = ""
		}
	}

	if strings.ToLower(title) == strings.ToLower(input) {
		return true
	}
	return false
}

func checkMusic(track *spotify.FullTrack) bool {
	var artistName string
	for _, artist := range track.Artists {
		artistName = artist.Name
	}
	lyricApi := lyrics.New()
	lyric, err := lyricApi.Search(artistName, track.Name)

	if err != nil {
		return false
	} else if lyric == "" {
		return false
	}
	return true
}

func TestPlaylist() {
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

	for i := 0; i < playlist.Tracks.Total; i++ {
		fmt.Println("Checking track" + strconv.Itoa(i))
		if !checkMusic(&playlist.Tracks.Tracks[i].Track) {
			_, err := client.RemoveTracksFromPlaylist(playlistID, playlist.Tracks.Tracks[i].Track.ID)
			if err != nil {
				log.Printf("error removing track %s from playlist: %v", playlist.Tracks.Tracks[i].Track.Name, err)
			} else {
				fmt.Printf("Removed track %s from playlist\n", playlist.Tracks.Tracks[i].Track.Name)
			}
		}
	}
}

func ExtractSuffix(s string) string {
	parts := strings.Split(s, "-")
	return parts[1]
}

func RandomString() string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}
