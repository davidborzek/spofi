package setup

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/urfave/cli/v2"
)

type (
	surveyAnswer struct {
		ClientID     string
		ClientSecret string
	}
)

var (
	Cmd = &cli.Command{
		Name:   "setup",
		Usage:  "Starts the setup process",
		Action: setup,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "The port of the http callback server.",
				Required: false,
				Value:    8080,
			},
			&cli.StringFlag{
				Name:     "host",
				Usage:    "The host of the http callback server.",
				Required: false,
				Value:    "localhost",
			},
		},
	}

	qs = []*survey.Question{
		{
			Name:     "clientId",
			Prompt:   &survey.Input{Message: "Enter the Client ID:"},
			Validate: survey.Required,
		},
		{
			Name:     "clientSecret",
			Prompt:   &survey.Password{Message: "Click on 'Show Client Secret' and enter the Client Secret:"},
			Validate: survey.Required,
		},
	}
)

// runSurvey runs the survey to ask the user for
// a client id and a client secret.
func runSurvey(host string, port int) (*surveyAnswer, error) {
	fmt.Printf(`Welcome to spofi setup!
WARNING: If you already have configured spofi, this will overwrite you current configuration.
		
1) Visit https://developer.spotify.com/dashboard/applications and click on "Create an app".
2) Enter a name and description.
3) Click on "Edit Settings" and add 'http://%s:%d' to the "Redirect URIs"
   by clicking on "Add" an save the settings with "Save".
4) Enter the app details in the following steps.\n`, host, port)

	var answers surveyAnswer
	err := survey.Ask(qs, &answers, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = ""
	}))
	if err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("Setup cancelled.")
			os.Exit(0)
		}

		return nil, err
	}

	return &answers, nil
}

// setup starts a new spofi setup by asking the user
// for a client id and a client secret and runs
// the oauth flow.
func setup(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.Int("port")

	answers, err := runSurvey(host, port)
	if err != nil {
		return err
	}

	if err := startAuthentication(
		answers.ClientID,
		answers.ClientSecret,
		host,
		port,
	); err != nil {
		return err
	}

	fmt.Println("\nSetup finished. You can now use spofi to control spotify.")

	return nil
}
