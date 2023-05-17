package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aditya-K2/utils"
	"github.com/zmb3/spotify/v2"
	_auth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	redirectURI string
	scopes      []string
	auth        *_auth.Authenticator
	ch          chan *payload
	state       string
	tokenPath   string
}

type payload struct {
	Token *oauth2.Token
	Err   error
}

func NewAuthConfig(userConfigPath string) *AuthConfig {
	redirectURI := "http://localhost:8080/callback"
	scopes := []string{
		_auth.ScopeUserLibraryRead,
		_auth.ScopePlaylistModifyPublic,
		_auth.ScopePlaylistModifyPrivate,
		_auth.ScopePlaylistReadPrivate,
		_auth.ScopePlaylistReadCollaborative,
		_auth.ScopeUserReadPlaybackState,
		_auth.ScopeUserModifyPlaybackState,
		_auth.ScopeUserLibraryModify,
		_auth.ScopeUserLibraryRead,
		_auth.ScopeUserReadPrivate,
		_auth.ScopeUserFollowRead,
		_auth.ScopeUserReadCurrentlyPlaying,
		_auth.ScopeUserModifyPlaybackState,
		_auth.ScopeUserReadRecentlyPlayed,
		_auth.ScopeUserTopRead,
		_auth.ScopeStreaming,
	}
	auth := _auth.New(
		_auth.WithRedirectURL(redirectURI),
		_auth.WithScopes(scopes...),
	)
	state := "__STP_TERMINAL_AUTH__"
	tokenPath := filepath.Join(userConfigPath, "oauthtoken")
	ch := make(chan *payload)

	return &AuthConfig{
		redirectURI: redirectURI,
		scopes:      scopes,
		auth:        auth,
		ch:          ch,
		state:       state,
		tokenPath:   tokenPath,
	}
}

func (ac *AuthConfig) InitClient() (*spotify.Client, error) {
	clientID := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("SPOTIFY_ID and/or SPOTIFY_SECRET are missing. Please make sure you have set the SPOTIFY_ID and SPOTIFY_SECRET environment variables")
	}

	token := &oauth2.Token{}
	// tokenErr := errors.New("")

	if utils.FileExists(ac.tokenPath) {
		var content []byte
		content, tokenErr := os.ReadFile(ac.tokenPath)
		if tokenErr != nil {
			return nil, tokenErr
		}

		tokenErr = json.Unmarshal(content, &token)
		if tokenErr != nil {
			return nil, tokenErr
		}
	}

	if token.AccessToken == "" {
		http.HandleFunc("/callback", ac.completeAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got request for: ", r.URL.String())
		})

		go func() {
			port := ":8080"
			err := http.ListenAndServe(port, nil)
			if err != nil {
				log.Fatal(err)
				ac.ch <- &payload{nil, err}
			}
		}()
		url := ac.auth.AuthURL(ac.state)

		fmt.Println("*** Please log in to Spotify by visiting the following page in your browser: ")
		fmt.Println(url)

		payload := <-ac.ch

		if payload.Err != nil {
			return nil, payload.Err
		}

		token = payload.Token
	}

	ctx := context.Background()
	client := spotify.New(ac.auth.Client(ctx, token))
	return client, nil
}

func (ac *AuthConfig) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := ac.auth.Token(r.Context(), ac.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
		ac.ch <- &payload{nil, err}
	}

	if st := r.FormValue("state"); st != ac.state {
		http.NotFound(w, r)
		_s := fmt.Sprintf("State mismatch: %s != %s\n", st, ac.state)
		log.Fatalf(_s)
		ac.ch <- &payload{nil, errors.New(_s)}
	}

	if val, merr := json.Marshal(tok); merr != nil {
		ac.ch <- &payload{nil, merr}
	} else {
		stptDir := filepath.Dir(ac.tokenPath)
		if !utils.FileExists(stptDir) {
			if derr := os.Mkdir(stptDir, 0777); derr != nil {
				ac.ch <- &payload{nil, derr}
			}
		}
		if werr := os.WriteFile(ac.tokenPath, val, 0777); werr != nil {
			ac.ch <- &payload{nil, werr}
		}
	}
	ac.ch <- &payload{tok, nil}
}
