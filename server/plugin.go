package main

import (
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"sync"

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

	//if bot, err := p.API.CreateBot(createOpBot()); err != nil {
	//	return errors.Wrapf(err, "failed to register %s bot", bot.Username)
	//}

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
	siteURL := p.GetSiteURL()
	siteURL = strings.Replace(siteURL, "8065", "3000", 1)
	res, _ := http.Get(siteURL)
	txt, _ := io.ReadAll(res.Body)
	resp := &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text: string(txt),
		Username: opBot,
		IconURL: siteURL + "/getLogo",
	}
	return resp, nil
}

func (p *Plugin) GetSiteURL() string {
	siteURL := ""
	ptr := p.API.GetConfig().ServiceSettings.SiteURL
	if ptr != nil {
		siteURL = *ptr
	}
	return siteURL
}

func createOpCommand() *model.Command {
	return &model.Command{
		Id:                   "",
		Token:                "",
		CreateAt:             0,
		UpdateAt:             0,
		DeleteAt:             0,
		CreatorId:            "",
		TeamId:               "",
		Trigger:              opCommand,
		Method:               "POST",
		Username:             "op-mattermost",
		IconURL:              "http://localhost:3000/getLogo",
		AutoComplete:         true,
		AutoCompleteDesc:     "Invoke OpenProject bot for Mattermost",
		AutoCompleteHint:     "",
		DisplayName:          "op-mattermost",
		Description:          "OpenProject integration for Mattermost",
		URL:                  "http://localhost:3000/",
		AutocompleteData:     nil,
		AutocompleteIconData: "",
	}
}

func createOpBot() *model.Bot {
	return &model.Bot{
			Username:		opBot,
			DisplayName:	opBot,
			Description:	"OpenProject Bot",
		}
}
