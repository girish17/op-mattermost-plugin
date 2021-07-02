package util

import (
	"encoding/json"
	"github.com/girish17/op-mattermost-plugin/server/types"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"strconv"
)

func getProjectOptAttachmentJSON(pluginURL string, action string, options []types.Option) []byte {
	attachments := types.OptAttachments{Attachments: []types.Attachment{
		{
			Actions: []types.Action{
				{
					Name:        "Type to search for a project...",
					Integration: types.Integration{
						Url:     pluginURL + "/projSel",
						Context: types.Context{
							Action: action,
						},
					},
					Type:        "select",
					Options:     options,
				},
				{
					Name:        "Cancel Project search",
					Integration: types.Integration{
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

func getWPOptAttachmentJSON(pluginURL string, action string, options []types.Option) []byte {
	attachments := types.OptAttachments{Attachments: []types.Attachment{
		{
			Actions: []types.Action{
				{
					Name:        "Type to search for a work package...",
					Integration: types.Integration{
						Url:     pluginURL + "/wpSel",
						Context: types.Context{
							Action: action,
						},
					},
					Type:        "select",
					Options:     options,
				},
				{
					Name:        "Cancel WP search",
					Integration: types.Integration{
						Url: pluginURL + "/createTimeLog",
					},
				},
			},
		},
	},
	}
	attachmentsJSON, _ := json.Marshal(attachments)
	return attachmentsJSON
}

func getOptArrayForProjectElements(elements []types.Element) []types.Option {
	var options []types.Option
	for _, element := range elements {
		id := strconv.Itoa(element.Id)
		options = append(options, types.Option{
			Text:  element.Name,
			Value: "opt" + id,
		})
	}
	return options
}

func getOptArrayForWPElements(elements []types.Element) []types.Option {
	var options []types.Option
	for _, element := range elements {
		id := strconv.Itoa(element.Id)
		options = append(options, types.Option{
			Text:  element.Subject,
			Value: "opt" + id,
		})
	}
	return options
}

func getOptArrayForAllowedValues(allowedValues []types.AllowedValues) []*model.PostActionOptions {
	var postActionOptions []*model.PostActionOptions
	for _, value := range allowedValues {
		id := strconv.Itoa(value.Id)
		postActionOptions = append(postActionOptions, &model.PostActionOptions{
			Text:  value.Name,
			Value: "opt" + id,
		})
	}
	return postActionOptions
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

func getUpdatePostMsg(userId string, channelId string, msg string) *model.Post {
	var post *model.Post
	post = &model.Post{
		Id:        menuPost.Id,
		UserId:    userId,
		ChannelId: channelId,
		Message:   msg,
	}
	return post
}


func getCreatePostMsg(userId string, channelId string, msg string) *model.Post {
	var post *model.Post
	post = &model.Post{
		UserId:        userId,
		ChannelId:     channelId,
		Message:       msg,
	}
	return post
}

func setOPStr(p plugin.MattermostPlugin) {
	opUrl, _ := p.API.KVGet("opUrl")
	apiKey, _ := p.API.KVGet("apiKey")
	opUrlStr = string(opUrl)
	apiKeyStr = string(apiKey)
}
