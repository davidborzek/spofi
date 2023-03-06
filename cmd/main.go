package cmd

import (
	"os"

	"github.com/davidborzek/spofi/cmd/setup"
	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/config"
	"github.com/davidborzek/spofi/internal/theme"
	"github.com/davidborzek/spofi/internal/views"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/urfave/cli/v2"
)

func start(ctx *cli.Context) error {
	cfg, err := config.LoadConfig()
	if config.IsConfigNotExistsErr(err) {
		return setup.Cmd.Action(ctx)
	}

	if err != nil {
		return err
	}

	if cfg.IsConfigIncomplete() {
		return setup.Cmd.Action(ctx)
	}

	if cfg.Theme != "" {
		rofi.SetCustomTheme(cfg.Theme)
	}

	customTheme := ctx.String("theme")
	if customTheme != "" {
		rofi.SetCustomTheme(customTheme)
	} else {
		theme.LoadTheme()
	}

	appCtx := app.NewApp(cfg)

	views.NewMainView(appCtx).
		Show()

	return nil
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "spofi"
	app.Usage = "Control spotify using rofi."
	app.Version = "local"
	app.Commands = []*cli.Command{
		setup.Cmd,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "theme",
			Required: false,
			Usage:    "Set a custom rofi theme",
		},
	}
	app.Action = start

	return app
}

func Main(args []string) {
	if err := newApp().Run(args); err != nil {
		os.Exit(1)
	}
}
