package app

import (
	"github.com/davidborzek/spofi/internal/config"
	"github.com/davidborzek/spofi/internal/player"
	"github.com/davidborzek/spofi/pkg/spotify"
)

// App represents the application context.
type App struct {
	Config *config.Config

	SpotifyClient spotify.Client

	Player player.Player
}

// NewApp creates a new application context
// for a given config.
func NewApp(cfg *config.Config) *App {
	sp := spotify.NewClient(
		cfg.Spotify.RefreshToken,
		cfg.Spotify.ClientID,
		cfg.Spotify.ClientSecret,
	)

	a := App{
		Config:        cfg,
		SpotifyClient: sp,
		Player:        player.New(sp, cfg.Device.ID),
	}

	return &a
}
