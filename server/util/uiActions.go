package util

import (
	"encoding/json"
	"github.com/girish17/op-mattermost-plugin/server/types"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/senseyeio/duration"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const opBot = "op-mattermost"


var menuPost *model.Post

var respHTTP http.Response

var client =  &http.Client{}

var opUrlStr string

var apiKeyStr string

func OpAuth(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from dialog submit: ", jsonBody)
	submission := jsonBody["submission"].(map[string]interface{})
	mmUserID := jsonBody["user_id"].(string)
	for key, value := range submission {
		p.API.LogInfo("Storing OpenProject auth credentials: " + key + ":" + value.(string))
		_ = p.API.KVSet(key, []byte(value.(string)))
	}
	setOPStr(p)
	opUserID := []byte(apiKeyStr + " " + opUrlStr)

	_ = p.API.KVDelete(opUrlStr)
	_ = p.API.KVDelete(apiKeyStr)

	user, _ := p.API.GetUserByUsername(opBot)
	resp, err := GetUserDetails(opUrlStr, apiKeyStr)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes map[string]string
		_ = json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogDebug("Response from op-mattermost: ", opJsonRes)
		if opJsonRes["_type"] != "Error" {
			_ = p.API.KVSet(mmUserID, opUserID)
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), "Hello "+opJsonRes["name"]+" :)")
			var attachmentMap map[string]interface{}
			_ = json.Unmarshal([]byte(GetAttachmentJSON(pluginURL)), &attachmentMap)
			post.SetProps(attachmentMap)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			respHTTP.Write(w)
		} else {
			p.API.LogError(opJsonRes["errorIdentifier"] + " " + opJsonRes["message"])
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), opJsonRes["message"])
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	} else {
			p.API.LogError("OpenProject login failed: ", err)
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), "OpenProject authentication failed. Please try again.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	menuPost, _ = p.API.CreatePost(post)
}

func ShowSelProject(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string)  {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("apikey: "+ apiKeyStr +" opURL: "+ opUrlStr)
	user, _ := p.API.GetUserByUsername(opBot)
	resp, err := GetProjects(opUrlStr, apiKeyStr)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes types.Projects
		_ = json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogInfo("Projects response from op-mattermost: ", opJsonRes)
		if opJsonRes.Type != "Error" {
			p.API.LogInfo("Projects obtained from OP: ", opJsonRes.Embedded.Elements)
			var options = getOptArrayForProjectElements(opJsonRes.Embedded.Elements)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "*Please select a project*")
			var attachmentMap map[string]interface{}
			_ = json.Unmarshal(getProjectOptAttachmentJSON(pluginURL, "showSelWP", options), &attachmentMap)
			post.SetProps(attachmentMap)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			respHTTP.Write(w)
		} else {
			p.API.LogError("Failed to fetch projects from OpenProject")
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to fetch projects from OpenProject")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p.API.LogError("Failed to fetch projects from OpenProject: ", err)
		post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to fetch projects from OpenProject")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.API.UpdatePost(post)
}

func WPHandler(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from project select: ", jsonBody["context"])
	submission := jsonBody["context"].(map[string]interface{})
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
			break
		case "selected_option":
			selectedOption = strings.Split(value.(string), "opt")
			p.API.LogInfo("selected option: " + selectedOption[1])
			break
		}
	}
	switch action {
	case "showSelWP":
		user, _ := p.API.GetUserByUsername(opBot)
		resp, err := GetWPsForProject(selectedOption[1], opUrlStr, apiKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := ioutil.ReadAll(resp.Body)
			var opJsonRes types.WorkPackages
			_ = json.Unmarshal(opResBody, &opJsonRes)
			p.API.LogInfo("Work packages response from op-mattermost: ", opJsonRes)
			if opJsonRes.Type != "Error" {
				p.API.LogInfo("Work packages obtained from OP: ", opJsonRes.Embedded.Elements)
				var options = getOptArrayForWPElements(opJsonRes.Embedded.Elements)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "*Please select a work package*")
				var attachmentMap map[string]interface{}
				_ = json.Unmarshal(getWPOptAttachmentJSON(pluginURL, "showTimeLogDlg", options), &attachmentMap)
				post.SetProps(attachmentMap)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				respHTTP.Write(w)
			} else {
				p.API.LogError("Failed to fetch work packages from OpenProject")
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to work packages from OpenProject")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			p.API.LogError("Failed to fetch work packages from OpenProject: ", err)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to work packages from OpenProject")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		p.API.UpdatePost(post)
		break
	case "createWP":
		http.NotFound(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func LoadTimeLogDlg(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	triggerId := jsonBody["trigger_id"].(string)
	submission := jsonBody["context"].(map[string]interface{})
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
			break
		case "selected_option":
			selectedOption = strings.Split(value.(string), "opt")
			p.API.LogInfo("selected option: " + selectedOption[1])
			break
		}
	}
	switch action {
	case "showTimeLogDlg":
		user, _ := p.API.GetUserByUsername(opBot)
		var timeEntriesBody types.TimeEntriesBody
		timeEntriesBody.Links.WorkPackage.Href = "/api/v3/work_packages/" + selectedOption[1]
		p.API.LogDebug("Time entries body: ", timeEntriesBody)
		timeEntriesBodyJSON, _ := json.Marshal(timeEntriesBody)
		resp, err := PostTimeEntriesForm(timeEntriesBodyJSON, opUrlStr, apiKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := ioutil.ReadAll(resp.Body)
			var opJsonRes types.TimeEntries
			_ = json.Unmarshal(opResBody, &opJsonRes)
			p.API.LogDebug("Time entries response from OpenProject: ", opJsonRes)
			if opJsonRes.Type != "Error" {
				var options = getOptArrayForAllowedValues(opJsonRes.Embedded.Schema.Activity.Embedded.AllowedValues)
				openLogTimeDialog(p, triggerId, pluginURL, options)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Opening time log dialog...")
				p.API.UpdatePost(post)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				respHTTP.Write(w)
			} else {
				p.API.LogError("Failed to fetch activities from OpenProject")
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to activities from OpenProject")
				p.API.UpdatePost(post)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			p.API.LogError("Failed to fetch activities from OpenProject")
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to activities from OpenProject")
			p.API.UpdatePost(post)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		break
	case "cnfDelWP":
		http.NotFound(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func openLogTimeDialog(p plugin.MattermostPlugin, triggerId string, pluginURL string, options []*model.PostActionOptions) {
	p.API.LogInfo("Activities from op-mattermost: ", options)
	p.API.LogDebug("Trigger ID for log time dialog: ", triggerId)
	dialog := p.API.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerId,
		URL:       pluginURL + "/logTime",
		Dialog: model.Dialog{
			CallbackId:       "log_time_dlg",
			Title:            "Log time for work package",
			IntroductionText: "Please enter details to log time",
			IconURL:          pluginURL + "/public/op_logo.jpg",
			Elements: []model.DialogElement{model.DialogElement{
				DisplayName: "Date",
				Name:        "spent_on",
				Type:        "text",
				Default:     time.Now().Format("2006-01-02"),
				Placeholder: "YYYY-MM-DD",
				HelpText:    "Please enter date within last one year and in YYYY-MM-DD format",
			}, model.DialogElement{
				DisplayName: "Comment",
				Name:        "comments",
				Type:        "textarea",
				Placeholder: "Please mention comments if any",
				Optional:    true,
			}, model.DialogElement{
				DisplayName: "Select Activity",
				Name:        "activity",
				Type:        "select",
				Default:     options[0].Value,
				Placeholder: "Type to search for activity",
				Options:     options,
			}, model.DialogElement{
				DisplayName: "Spent hours",
				Name:        "spent_hours",
				Type:        "text",
				Default:     "0.5",
				Placeholder: "hours like 0.5, 1, 3 ...",
				HelpText:    "Please enter spent hours to be logged",
			}, model.DialogElement{
				DisplayName: "Billable hours",
				Name:        "billable_hours",
				Type:        "text",
				Default:     "0.0",
				Placeholder: "hours like 0.5, 1, 3 ...",
				HelpText:    "Please ensure billable hours is less than or equal to spent hours",
			}},
			SubmitLabel:    "Log time",
			NotifyOnCancel: true,
		},
	})
	p.API.LogDebug("Dialog object returned: ", dialog)
}

func GetTimeLog(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	user, _ := p.API.GetUserByUsername(opBot)
	resp, err := GetTimeEntries(opUrlStr, apiKeyStr)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes types.TimeEntryList
		_ = json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogDebug("Time entries response from OpenProject: ", opJsonRes)
		if opJsonRes.Type != "Error" {
			var timeLogs = getOptArrayForTimeEntries(opJsonRes.Embedded.Elements)
			p.API.LogInfo("Time entries from op-mattermost: ", timeLogs)
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), timeLogs)
			p.API.CreatePost(post)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			respHTTP.Write(w)
		} else {
			p.API.LogError("Failed to fetch time entries from OpenProject")
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to time entries from OpenProject")
			p.API.UpdatePost(post)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p.API.LogError("Failed to fetch time entries from OpenProject")
		post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Failed to time entries from OpenProject")
		p.API.UpdatePost(post)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOptArrayForTimeEntries(elements []types.TimeElement) string {
	var tableTxt string
	if len(elements) != 0 {
		tableTxt = "#### Time entries logged by you\n"
		tableTxt += "| Spent On | Project | Work Package | Activity | Logged Time | Billed Time | Comment |\n"
		tableTxt += "|:---------|:--------|:-------------|:---------|:------------|:------------|:--------|\n"
		for _, element := range elements {
			d, _ := duration.ParseISO8601(element.Hours)
			var loggedTime = ""
			if d.TH != 0 {
				hours := strconv.Itoa(d.TH)
				if d.TH > 1 {
					loggedTime = hours + " hours "
				} else {
					loggedTime = hours + " hour "
				}
			}
			if d.TM != 0 {
				minutes := strconv.Itoa(d.TM)
				if d.TM > 1 {
					loggedTime = loggedTime + minutes + " minutes"
				} else {
					loggedTime = loggedTime + minutes + " minute"
				}
			}
			billedHours := strconv.FormatFloat(element.CustomField1, 'f', 2, 64)
			tableTxt += "| " + element.SpentOn + " | " + element.Links.Project.Title + " | " + element.Links.WorkPackage.Title + " | " + element.Links.Activity.Title + " | " + loggedTime + " | " + billedHours + " hours" + " | " + strings.ReplaceAll(element.Comment.Raw, "/\n/g", " ") + " |\n"
		}
	} else {
		tableTxt = "Couldn't find time entries logged by you :confused: Try logging time using `/op`"
	}
	return tableTxt
}

func ShowDelWPSel(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {

}

func Logout(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var jsonBody map[string]interface{}
		_ = json.Unmarshal(body, &jsonBody)
		mmUserID := jsonBody["user_id"].(string)
		p.API.LogInfo("Deleting op login for mm user id: " + mmUserID)
		_ = p.API.KVDelete(mmUserID)

		user, _ := p.API.GetUserByUsername(opBot)
		post := getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), ":wave:")
		_, _ = p.API.UpdatePost(post)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		respHTTP.Write(w)
}
