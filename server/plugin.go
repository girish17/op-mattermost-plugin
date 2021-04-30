package main

import (
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const opCommand = "op"
const opBot = "op-mattermost"

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	router *mux.Router
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	if p.API.GetConfig().ServiceSettings.SiteURL == nil {
		p.API.LogError("SiteURL must be set. Some features will operate incorrectly if the SiteURL is not set. See documentation for details: http://about.mattermost.com/default-site-url")
	}

	if err := p.API.RegisterCommand(createOpCommand()); err != nil {
		return errors.Wrapf(err, "failed to register %s command", opCommand)
	}

	if bot, err := p.API.CreateBot(createOpBot()); err != nil {
		return errors.Wrapf(err, "failed to register %s bot", bot.Username)
	}

	return nil
}

func createOpCommand() *model.Command {
	return &model.Command{
		Trigger:          opCommand,
		AutoComplete:     true,
		AutoCompleteDesc: "Invoke OpenProject bot for Mattermost",
	}
}

func createOpBot() *model.Bot {
	return &model.Bot{
			Username:		opBot,
			DisplayName:	opBot,
			Description:	"OpenProject Bot",
		}
}
