package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/girish17/op-mattermost-plugin/server/util"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const opCommand = "op"
const opBot = "op-mattermost"

var pluginURL string

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
	switch path := r.URL.Path; path {
	case "/opAuth":
		util.OpAuth(p.MattermostPlugin, w, r, pluginURL)
		break
	case "/createTimeLog":
		util.ShowSelProject(p.MattermostPlugin, w, r, pluginURL)
		break
	case "/projSel":
		util.WPHandler(p.MattermostPlugin, w, r, pluginURL)
		break
	case "/wpSel":
		util.LoadTimeLogDlg(p.MattermostPlugin, w, r, pluginURL)
		break
	case "/logTime":
		http.NotFound(w, r)
		break
	case "/getTimeLog":
		util.GetTimeLog(p.MattermostPlugin, w, r, pluginURL)
		break
	case "/delTimeLog":
		http.NotFound(w, r)
		break
	case "/createWP":
		http.NotFound(w, r)
		break
	case "/saveWP":
		http.NotFound(w, r)
		break
	case "/delWP":
		http.NotFound(w, r)
		break
	case "/bye":
		util.Logout(p.MattermostPlugin, w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
func (p *Plugin) OnActivate() error {
	if p.API.GetConfig().ServiceSettings.SiteURL == nil {
		p.API.LogError("SiteURL must be set. Some features will operate incorrectly if the SiteURL is not set. See documentation for details: http://about.mattermost.com/default-site-url")
	}

	if err := p.API.RegisterCommand(createOpCommand(p.GetSiteURL())); err != nil {
		return errors.Wrapf(err, "failed to register %s command", opCommand)
	}

	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError){
	siteURL := p.GetSiteURL()
	pluginURL = getPluginURL(siteURL)
	p.API.LogDebug("Plugin URL :"+pluginURL)
	if opUserID, _ := p.API.KVGet(args.UserId); opUserID == nil {
		p.API.LogDebug("Creating interactive dialog...")
		p.MattermostPlugin.API.OpenInteractiveDialog(model.OpenDialogRequest{
			TriggerId: args.TriggerId,
			URL:        pluginURL + "/opAuth",
			Dialog:    model.Dialog{
				CallbackId:       "op_auth_dlg",
				Title:            "OpenProject Authentication",
				IntroductionText: "Please enter credentials to log in",
				IconURL:          getLogoURL(siteURL),
				Elements: []model.DialogElement{model.DialogElement{
					DisplayName: "OpenProject URL",
					Name:        "opUrl",
					Type:        "text",
					Default:     "http://localhost:8080",
					Placeholder: "http://localhost:8080",
					Optional:    false,
					HelpText:    "Please enter the URL of OpenProject server",
				}, model.DialogElement{
					DisplayName: "OpenProject api-key",
					Name:        "apiKey",
					Type:        "text",
					Placeholder: "api-key generated from your account page in OpenProject",
					Optional:    false,
					HelpText:    "api-key can be generated within 'My account' section of OpenProject",
				}},
				SubmitLabel:      "Log in",
				NotifyOnCancel:   true,
			},
		})

		resp := &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: "opening op auth dialog",
			Username: opBot,
			IconURL: getLogoURL(siteURL),
		}

		return resp, nil
	} else {
		cmd := args.Command
		cmdAction := strings.Split(cmd, " ")
		p.API.LogInfo("Command arg entered: "+cmdAction[1])
		opUserIDStr := string(opUserID)
		apiKeyStr := strings.Split(opUserIDStr, " ")
		opUrlStr := apiKeyStr[1]
		p.API.LogInfo("Retrieving from KV: opURL - " + opUrlStr + " apiKey - " + apiKeyStr[0])

		var cmdResp *model.CommandResponse

		client := &http.Client{}
		req, _ := http.NewRequest("GET", opUrlStr + "/api/v3/users/me", nil)
		req.SetBasicAuth("apikey", apiKeyStr[0])
		resp, _ := client.Do(req)
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes map[string]string
		json.Unmarshal(opResBody, &opJsonRes)
		p.MattermostPlugin.API.LogDebug("Response from op-mattermost: ", opJsonRes["firstName"])

		var attachmentMap map[string]interface{}
		json.Unmarshal([]byte(util.GetAttachmentJSON(pluginURL)), &attachmentMap)

		cmdResp = &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
			Text: "Hello "+opJsonRes["name"]+" :)",
			Username: opBot,
			IconURL: getLogoURL(siteURL),
			Props: attachmentMap,
		}

		return cmdResp, nil
	}
}

func (p *Plugin) GetSiteURL() string {
	siteURL := ""
	ptr := p.API.GetConfig().ServiceSettings.SiteURL
	if ptr != nil {
		siteURL = *ptr
	}
	return siteURL
}

func getLogoURL(siteURL string) string {
	return getPluginURL(siteURL) + "/public/op_logo.jpg"
}

func getPluginURL(siteURL string) string {
	return siteURL + "/plugins/" + manifest.Id
}

func createOpCommand(siteURL string) *model.Command {
	return &model.Command{
		Trigger:              opCommand,
		Method:               "POST",
		Username:             opBot,
		IconURL:              getLogoURL(siteURL),
		AutoComplete:         true,
		AutoCompleteDesc:     "Invoke OpenProject bot for Mattermost",
		AutoCompleteHint:     "",
		DisplayName:          opBot,
		Description:          "OpenProject integration for Mattermost",
		URL:                  siteURL,
	}
}

func (p *Plugin) setBotIcon() {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogError("failed to get bundle path", err)
	}

	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "op_logo.svg"))
	if err != nil {
		p.API.LogError("failed to read profile image", err)
	}

	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError("failed to fetch bot user", err)
	}

	if appErr := p.API.SetBotIconImage(user.UserId, profileImage); appErr != nil {
		p.API.LogError("failed to set profile image", appErr)
	}
}
