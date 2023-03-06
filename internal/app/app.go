package app

import (
	"github.com/davidborzek/spofi/internal/config"
	"github.com/davidborzek/spofi/pkg/spotify"
)

// App represents the application context.
type App struct {
	Config *config.Config

	SpotifyClient spotify.Client
}

// NewApp creates a new application context
// for a given config.
func NewApp(cfg *config.Config) *App {
	a := App{
		Config: cfg,
		SpotifyClient: spotify.NewClient(
			cfg.Spotify.RefreshToken,
			cfg.Spotify.ClientID,
			cfg.Spotify.ClientSecret,
		),
	}

	return &a
}
