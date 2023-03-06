package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Default icon set (requires Jetbrains NerdFont)
const (
	defaultIconAlbum          = ""
	defaultIconDevice         = ""
	defaultIconLikedTracks    = ""
	defaultIconNext           = "⏭"
	defaultIconPause          = "⏸"
	defaultIconPlay           = "▶"
	defaultIconPlayer         = "奈"
	defaultIconPrevious       = "⏮"
	defaultIconQueue          = "蘿"
	defaultIconRecentlyPlayed = ""
	defaultIconRepeatContext  = "凌"
	defaultIconRepeatOff      = "稜"
	defaultIconRepeatTrack    = "綾"
	defaultIconSearch         = ""
	defaultIconShuffleOff     = "劣"
	defaultIconShuffleOn      = "咽"
	defaultIconTrack          = ""
)

const (
	configDirName  = "spofi"
	configFileName = "spofi.yaml"
)

// SpotifyDevice represents a saved spotify
// device in the config.
type SpotifyDevice struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

// KeyConfig represent the hotkey configuration.
type KeyConfig struct {
	NextPage          string `yaml:"nextPage"`
	PreviousPage      string `yaml:"previousPage"`
	AddToQueue        string `yaml:"addToQueue"`
	TogglePauseResume string `yaml:"togglePauseResume"`
	NextTrack         string `yaml:"nextTrack"`
	PreviousTrack     string `yaml:"previousTrack"`
	PlayAlbum         string `yaml:"playAlbum"`
	PlayTrack         string `yaml:"playTrack"`
	ToggleSearchType  string `yaml:"toggleSearchType"`
}

type SpotifyConfig struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RefreshToken string `yaml:"refreshToken"`
}

type IconConfig struct {
	Album          string `yaml:"album"`
	Device         string `yaml:"device"`
	LikedTracks    string `yaml:"likedTracks"`
	Next           string `yaml:"next"`
	Pause          string `yaml:"pause"`
	Play           string `yaml:"play"`
	Player         string `yaml:"player"`
	Previous       string `yaml:"previous"`
	Queue          string `yaml:"queue"`
	RecentlyPlayed string `yaml:"recentlyPlayed"`
	RepeatContext  string `yaml:"repeatContext"`
	RepeatOff      string `yaml:"repeatOff"`
	RepeatTrack    string `yaml:"repeatTrack"`
	Search         string `yaml:"search"`
	ShuffleOff     string `yaml:"shuffleOff"`
	ShuffleOn      string `yaml:"shuffleOn"`
	Track          string `yaml:"track"`
}

// Config represent the application config.
type Config struct {
	Spotify     SpotifyConfig `yaml:"spotify"`
	Device      SpotifyDevice `yaml:"device"`
	Theme       string        `yaml:"theme"`
	Keybindings KeyConfig     `yaml:"keybindings"`
	Icons       IconConfig    `yaml:"icons"`
}

// getConfigDir is an internal implementation
// to get the configuration directory based on the
// os user config dir.
func getConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", userConfigDir, configDirName), nil
}

// LoadConfig load the config at the given path.
func LoadConfig() (*Config, error) {
	cfgPath, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	rawCfg, err := os.ReadFile(fmt.Sprintf("%s/%s", cfgPath, configFileName))
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(rawCfg, &cfg); err != nil {
		return nil, err
	}

	cfg.fillDefaults()

	return &cfg, nil
}

// IsConfigNotExistsErr checks if the error
// from LoadConfig means that the config file
// does not exist.
func IsConfigNotExistsErr(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

func (cfg *KeyConfig) fillDefaults() {
	if cfg.AddToQueue == "" {
		cfg.AddToQueue = "Alt+d"
	}

	if cfg.NextPage == "" {
		cfg.NextPage = "Alt+Right"
	}

	if cfg.PreviousPage == "" {
		cfg.PreviousPage = "Alt+Left"
	}

	if cfg.TogglePauseResume == "" {
		cfg.TogglePauseResume = "Alt+space"
	}

	if cfg.NextTrack == "" {
		cfg.NextTrack = "Alt+n"
	}

	if cfg.PreviousTrack == "" {
		cfg.PreviousTrack = "Alt+p"
	}

	if cfg.PlayAlbum == "" {
		cfg.PlayAlbum = "Alt+p"
	}

	if cfg.PlayTrack == "" {
		cfg.PlayTrack = "Alt+t"
	}

	if cfg.ToggleSearchType == "" {
		cfg.ToggleSearchType = "Alt+s"
	}
}

func (cfg *IconConfig) fillDefaults() {
	if cfg.Album == "" {
		cfg.Album = defaultIconAlbum
	}

	if cfg.Track == "" {
		cfg.Track = defaultIconTrack
	}

	if cfg.Device == "" {
		cfg.Device = defaultIconDevice
	}

	if cfg.LikedTracks == "" {
		cfg.LikedTracks = defaultIconLikedTracks
	}

	if cfg.Pause == "" {
		cfg.Pause = defaultIconPause
	}

	if cfg.Play == "" {
		cfg.Play = defaultIconPlay
	}

	if cfg.ShuffleOff == "" {
		cfg.ShuffleOff = defaultIconShuffleOff
	}

	if cfg.ShuffleOn == "" {
		cfg.ShuffleOn = defaultIconShuffleOn
	}

	if cfg.RepeatOff == "" {
		cfg.RepeatOff = defaultIconRepeatOff
	}

	if cfg.RepeatContext == "" {
		cfg.RepeatContext = defaultIconRepeatContext
	}

	if cfg.RepeatTrack == "" {
		cfg.RepeatTrack = defaultIconRepeatTrack
	}

	if cfg.Player == "" {
		cfg.Player = defaultIconPlayer
	}

	if cfg.Next == "" {
		cfg.Next = defaultIconNext
	}

	if cfg.Previous == "" {
		cfg.Previous = defaultIconPrevious
	}

	if cfg.Queue == "" {
		cfg.Queue = defaultIconQueue
	}

	if cfg.RecentlyPlayed == "" {
		cfg.RecentlyPlayed = defaultIconRecentlyPlayed
	}

	if cfg.Search == "" {
		cfg.Search = defaultIconSearch
	}
}

func (cfg *Config) fillDefaults() {
	cfg.Keybindings.fillDefaults()
	cfg.Icons.fillDefaults()
}

// IsConfigIncomplete checks if the config is incomplete.
func (cfg *Config) IsConfigIncomplete() bool {
	return cfg.Spotify.ClientID == "" &&
		cfg.Spotify.ClientSecret == "" &&
		cfg.Spotify.RefreshToken == ""
}

// Write writes the current (in-memory) configuration to the config file.
func (cfg *Config) Write() error {
	cfg.fillDefaults()

	raw, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cfgPath, 0755); err != nil {
		fmt.Println(err)
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s", cfgPath, configFileName), raw, 0644)
}
