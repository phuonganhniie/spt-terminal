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
	sptauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type AuthConfig struct {
	redirectURI string
	scopes      []string
	auth        *sptauth.Authenticator
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
		sptauth.ScopeUserLibraryRead,
		sptauth.ScopePlaylistModifyPublic,
		sptauth.ScopePlaylistModifyPrivate,
		sptauth.ScopePlaylistReadPrivate,
		sptauth.ScopePlaylistReadCollaborative,
		sptauth.ScopeUserReadPlaybackState,
		sptauth.ScopeUserModifyPlaybackState,
		sptauth.ScopeUserLibraryModify,
		sptauth.ScopeUserLibraryRead,
		sptauth.ScopeUserReadPrivate,
		sptauth.ScopeUserFollowRead,
		sptauth.ScopeUserReadCurrentlyPlaying,
		sptauth.ScopeUserModifyPlaybackState,
		sptauth.ScopeUserReadRecentlyPlayed,
		sptauth.ScopeUserTopRead,
		sptauth.ScopeStreaming,
	}
	auth := sptauth.New(
		sptauth.WithRedirectURL(redirectURI),
		sptauth.WithScopes(scopes...),
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
	token := &oauth2.Token{}

	if utils.FileExists(ac.tokenPath) {
		var content []byte
		content, err := os.ReadFile(ac.tokenPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &token)
		if err != nil {
			return nil, err
		}
	}

	if token == nil {
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
