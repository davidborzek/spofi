package player

import "github.com/davidborzek/spofi/pkg/spotify"

// Player defines a player controller for a selected device.
type Player interface {
	// PlayPause toggles play pause.
	PlayPause() error
	// ToggleRepeat toggles the repeat state.
	ToggleRepeat() error
	// ToggleShuffle toggles the shuffle state.
	ToggleShuffle() error
	// PlayTracks play a given track.
	PlayTrack(uri string) error
	// PlayContext plays a given context (playlist, album, etc.)
	// and can optionally handle a given uri in the context.
	PlayContext(contextUri string, uri ...string) error
	// AddQueue adds a given tracks to the queue.
	AddQueue(uri string) error
	// Next changes to the next track.
	Next() error
	// Previous changes to the previous track.
	Previous() error
	// SetDevices set the devices for all operations.
	SetDevice(device string)
}

type player struct {
	client spotify.Client
	device string
}

func New(client spotify.Client, device string) Player {
	return &player{
		client: client,
		device: device,
	}
}

func (p *player) PlayPause() error {
	state, err := p.client.GetPlayer()
	if err != nil {
		return err
	}

	if state.IsPlaying {
		return p.client.Pause(p.device)
	}
	return p.client.Play(p.device)
}

func (p *player) ToggleRepeat() error {
	state, err := p.client.GetPlayer()
	if err != nil {
		return err
	}

	if state != nil {
		s := spotify.RepeatOff
		if state.RepeatState == spotify.RepeatOff {
			s = spotify.RepeatContext
		}
		if state.RepeatState == spotify.RepeatContext {
			s = spotify.RepeatTrack
		}

		return p.client.SetRepeatMode(p.device, s)
	}

	return nil
}

func (p *player) ToggleShuffle() error {
	state, err := p.client.GetPlayer()
	if err != nil {
		return err
	}

	if state != nil {
		return p.client.SetShuffleState(p.device, !state.ShuffleState)
	}

	return nil
}

func (p *player) PlayTrack(uri string) error {
	return p.client.PlayTrack(uri, p.device)
}

func (p *player) PlayContext(contextUri string, uri ...string) error {
	return p.client.PlayContext(contextUri, p.device, uri...)
}

func (p *player) AddQueue(uri string) error {
	return p.client.AddQueue(uri, p.device)
}

func (p *player) Next() error {
	return p.client.Next(p.device)
}

func (p *player) Previous() error {
	return p.client.Previous(p.device)
}

func (p *player) SetDevice(device string) {
	p.device = device
}
