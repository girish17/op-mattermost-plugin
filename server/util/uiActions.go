package util

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const opBot = "op-mattermost"


var menuPost *model.Post

var client =  &http.Client{}

var opUrlStr string

var apiKeyStr string

func getProjectOptAttachmentJSON(pluginURL string, action string, options []Option) []byte {
	attachments := ProjectOptAttachments{Attachments: []Attachment{
		{
			Actions: []Action{
				{
					Name:        "Type to search for a project...",
					Integration: Integration{
						Url:     pluginURL + "/projSel",
						Context: Context{
							Action: action,
						},
					},
					Type:        "select",
					Options:     options,
				},
				{
					Name:        "Cancel Project search",
					Integration: Integration{
						Url: pluginURL + "/bye",
					},
				},
			},
		},
		},
	}
	attachmentsJSON, _ := json.Marshal(attachments)
	return attachmentsJSON
}

func GetAttachmentJSON(pluginURL string) string {
	return `{
			"attachments": [
				  {
					"text": "What would you like me to do?",
					"actions": [
					  {
						"name": "Log time",
						"integration": {
						  "url": "` + pluginURL + `/createTimeLog",
						  "context": {
							"action": "showSelWP"
						  }
						}
					  },
					  {
						"name": "Create Work Package",
						"integration": {
                          "url": "` + pluginURL + `/createWP",
						  "context": {
							"action": "createWP"
						  }
						}
					  },
					  {
						"name": "View time logs",
						"integration": {
					      "url": "` + pluginURL + `/getTimeLog",
						  "context": {
							"action": "getTimeLog"
						  }
						}
					  },
					  {
						"name": "Delete time log",
						"integration": {
                          "url": "` + pluginURL + `/delTimeLog",
						  "context": {
							"action": "delTimeLog"
						  }
						}
					  },
					  {
						"name": "Delete Work Package",
						"integration": {
                         "url": "` + pluginURL + `/delWP",
						  "context": {
							"action": ""
						  }
						}
					  },
					  {
						"name": "Bye :wave:",
						"integration": {
                          "url": "` + pluginURL + `/bye",
						  "context": {
							"action": "bye"
						  }
						}
					  }
					]
				  }
			]
		}`
}

func setOPStr(p plugin.MattermostPlugin) {
	opUrl, _ := p.API.KVGet("opUrl")
	apiKey, _ := p.API.KVGet("apiKey")
	opUrlStr = string(opUrl)
	apiKeyStr = string(apiKey)
}

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
	req, _ := http.NewRequest("GET", opUrlStr + "/api/v3/users/me", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes map[string]string
		_ = json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogDebug("Response from op-mattermost: ", opJsonRes)
		if opJsonRes["_type"] != "Error" {
			_ = p.API.KVSet(mmUserID, opUserID)
			post = &model.Post{
				UserId:        user.Id,
				ChannelId:     jsonBody["channel_id"].(string),
				Message:       "Hello "+opJsonRes["name"]+" :)",
			}
			var attachmentMap map[string]interface{}
			_ = json.Unmarshal([]byte(GetAttachmentJSON(pluginURL)), &attachmentMap)
			post.SetProps(attachmentMap)
		} else {
			p.API.LogError(opJsonRes["errorIdentifier"] + " " + opJsonRes["message"])
			post = &model.Post{
				UserId:        user.Id,
				ChannelId:     jsonBody["channel_id"].(string),
				Message:       opJsonRes["message"],
			}
		}

	} else {
		p.API.LogError("OpenProject login failed: ", err)
		post = &model.Post{
			UserId:        user.Id,
			ChannelId:     jsonBody["channel_id"].(string),
			Message:       "OpenProject authentication failed. Please try again.",
		}
	}
	menuPost, _ = p.API.CreatePost(post)
}

func ShowSelProject(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string)  {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("apikey: "+apiKeyStr+" opURL: "+opUrlStr)
	user, _ := p.API.GetUserByUsername(opBot)
	req, _ := http.NewRequest("GET", opUrlStr + `/api/v3/projects`, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes Projects
		_ = json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogInfo("Projects response from op-mattermost: ", opJsonRes)
		if opJsonRes.Type != "Error" {
			p.API.LogInfo("Projects obtained from OP: ", opJsonRes.Embedded.Elements)
			var options []Option
			for _, element := range opJsonRes.Embedded.Elements {
				id := strconv.Itoa(element.Id)
				options = append(options, Option{
					Text:  element.Name,
					Value: "opt" + id,
				})
			}
			post = &model.Post{
				Id:            menuPost.Id,
				UserId:        user.Id,
				ChannelId:     jsonBody["channel_id"].(string),
				Message:       "*Please select a project*",
			}
			var attachmentMap map[string]interface{}
			_ = json.Unmarshal(getProjectOptAttachmentJSON(pluginURL, "showSelWP", options), &attachmentMap)
			post.SetProps(attachmentMap)
		} else {
			p.API.LogError("Failed to fetch projects from OpenProject")
			post = &model.Post{
				Id:            menuPost.Id,
				UserId:        user.Id,
				ChannelId:     jsonBody["channel_id"].(string),
				Message:       "Failed to fetch projects from OpenProject",
			}
		}
	} else {
		p.API.LogError("Failed to fetch projects from OpenProject: ", err)
		post = &model.Post{
			Id:            menuPost.Id,
			UserId:        user.Id,
			ChannelId:     jsonBody["channel_id"].(string),
			Message:       "Failed to fetch projects from OpenProject. Please try again.",
		}
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

		break
	case "createWP":
		break
	default:
		w.WriteHeader(400)
		w.Write([]byte("Invalid action type"))
	}
}

func Logout(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var jsonBody map[string]interface{}
		_ = json.Unmarshal(body, &jsonBody)
		mmUserID := jsonBody["user_id"].(string)
		p.API.LogInfo("Deleting op login for mm user id: " + mmUserID)
		_ = p.API.KVDelete(mmUserID)
		user, _ := p.API.GetUserByUsername(opBot)
		post := &model.Post{
			Id:            menuPost.Id,
			UserId:        user.Id,
			ChannelId:     jsonBody["channel_id"].(string),
			Message:       ":wave:",
		}
		_, _ = p.API.UpdatePost(post)
}
