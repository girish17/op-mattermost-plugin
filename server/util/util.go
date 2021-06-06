package util

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"io/ioutil"
	"net/http"
)

const opBot = "op-mattermost"

func GetPluginURL(siteURL string) string {
	return siteURL + "/plugins/com.girishm.op-mattermost-plugin"
}

func GetAttachmentJSON(siteURL string) string {
	pluginURL := GetPluginURL(siteURL)
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

var menuPost *model.Post

func OpAuth(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from dialog submit: ", jsonBody)
	submission := jsonBody["submission"].(map[string]interface{})
	mmUserID := jsonBody["user_id"].(string)
	for key, value := range submission {
		p.API.LogInfo("Storing OpenProject auth credentials: " + key + ":" + value.(string))
		p.API.KVSet(key, []byte(value.(string)))
	}
	opUrl, _ := p.API.KVGet("opUrl")
	apiKey, _ := p.API.KVGet("apiKey")
	opUrlStr := string(opUrl)
	apiKeyStr := string(apiKey)
	opUserID := []byte(apiKeyStr + " " + opUrlStr)

	p.API.KVDelete(opUrlStr)
	p.API.KVDelete(apiKeyStr)

	user, _ := p.API.GetUserByUsername(opBot)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", opUrlStr + "/api/v3/users/me", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	var post *model.Post
	if err == nil {
		opResBody, _ := ioutil.ReadAll(resp.Body)
		var opJsonRes map[string]string
		json.Unmarshal(opResBody, &opJsonRes)
		p.API.LogDebug("Response from op-mattermost: ", opJsonRes)
		if opJsonRes["_type"] != "Error" {
			p.API.KVSet(mmUserID, opUserID)
			post = &model.Post{
				UserId:        user.Id,
				ChannelId:     jsonBody["channel_id"].(string),
				Message:       "Hello "+opJsonRes["name"]+" :)",
			}
			var attachmentMap map[string]interface{}
			json.Unmarshal([]byte(GetAttachmentJSON(*p.API.GetConfig().ServiceSettings.SiteURL)), &attachmentMap)
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

func Logout(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var jsonBody map[string]interface{}
		json.Unmarshal(body, &jsonBody)
		mmUserID := jsonBody["user_id"].(string)
		p.API.LogInfo("Deleting op login for mm user id: " + mmUserID)
		p.API.KVDelete(mmUserID)
		user, _ := p.API.GetUserByUsername(opBot)
		post := &model.Post{
			Id:            menuPost.Id,
			UserId:        user.Id,
			ChannelId:     jsonBody["channel_id"].(string),
			Message:       ":wave:",
		}
		p.API.UpdatePost(post)
}
