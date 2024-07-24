package setup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/davidborzek/spofi/internal/config"
	"github.com/davidborzek/spofi/pkg/spotify"
)

var (
	scopes = []string{
		"user-library-read",
		"user-read-currently-playing",
		"user-read-playback-state",
		"user-read-recently-played",
		"user-library-modify",
		"user-modify-playback-state",
		"playlist-modify-private",
		"playlist-read-private",
		"playlist-modify-public",
	}

	srv      http.Server
	codeChan chan string = make(chan string)
)

// openBrowser opens a given url in a browser.
func openBrowser(url string) error {
	return exec.Command("xdg-open", url).Start()
}

// startAuthServer starts a server for the spotify
// oauth callback.
func startAuthServer(addr string) {
	srv = http.Server{
		Addr: addr,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		w.WriteHeader(200)
		w.Write([]byte("The setup process is complete! You may now close this window and use spofi to control spotify."))
		codeChan <- code
	})

	srv.ListenAndServe()
}

// shutdownServer shutdowns the http
// callback server.
func shutdownServer() {
	srv.Shutdown(context.Background())
}

// startAuthentication starts the spotify authentication
// process.
func startAuthentication(clientId string, clientSecret string, host string, port int) error {
	redirectUrl := fmt.Sprintf("http://%s:%d", host, port)

	sc := spotify.NewAuthClient(
		clientId,
		clientSecret,
		redirectUrl,
		scopes,
	)

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	go startAuthServer(addr)

	authUrl := sc.BuildAuthUrl()
	openBrowser(authUrl)

	fmt.Println("\nPlease follow the steps in your web browser and log in using your Spotify account. If the URL did not open automatically, please manually open the following URL:")
	fmt.Println(authUrl)

	code := <-codeChan
	shutdownServer()

	token, err := sc.GetTokenPair(code)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	cfg := config.Config{
		Spotify: config.SpotifyConfig{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RefreshToken: token.RefreshToken,
		},
	}

	return cfg.Write()
}
