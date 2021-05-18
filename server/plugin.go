package main

import (
	"github.com/pkg/errors"
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
	p.MattermostPlugin.API.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: args.TriggerId,
		URL:       siteURL + "/opAuth",
		Dialog:    model.Dialog{
			CallbackId:       "op_auth_dlg",
			Title:            "OpenProject Authentication",
			IntroductionText: "",
			IconURL:          siteURL + "/getLogo",
			Elements: []model.DialogElement{model.DialogElement{
				DisplayName: "OpenProject URL",
				Name:        "op_url",
				Type:        "text",
				SubType:     "",
				Default:     "http://localhost:8080",
				Placeholder: "http://localhost:8080",
				HelpText:    "Please enter the URL of OpenProject server",
				Optional:    false,
				MinLength:   0,
				MaxLength:   0,
				DataSource:  "",
				Options:     nil,
			}, model.DialogElement{
				DisplayName: "OpenProject api-key",
				Name:        "api_key",
				Type:        "text",
				SubType:     "",
				Default:     "",
				Placeholder: "api-key generated from your account page in OpenProject",
				HelpText:    "api-key can be generated within 'My account' section of OpenProject",
				Optional:    false,
				MinLength:   0,
				MaxLength:   0,
				DataSource:  "",
				Options:     nil,
			}},
			SubmitLabel:      "Log in",
			NotifyOnCancel:   true,
			State:            "",
		},
	})

	resp := &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: "opening op auth dialog",
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
