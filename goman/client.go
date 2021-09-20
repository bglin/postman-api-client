package goman

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.getpostman.com"}

type BaseObject struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Uid  string `json:"uid,omitempty"`
}
type Workspace struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Description  string       `json:"description"`
	Collections  []BaseObject `json:"collections,omitempty"`
	Environments []BaseObject `json:"environments,omitempty"`
	Mocks        []BaseObject `json:"mocks,omitempty"`
	Monitors     []BaseObject `json:"monitors,omitempty"`
}

type WorkspaceObject struct {
	Workspace Workspace `json:"workspace"`
}

//
type EnvValues struct {
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
	Enabled bool        `json:"enabled"`
}
type Environment struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Owner     string      `json:"owner"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Values    []EnvValues `json:"values"`
	IsPublic  bool        `json:"isPublic"`
}
type EnvironmentObject struct {
	Environment Environment `json:"environment"`
}

type PostmanClient struct {
	client *http.Client
	apiKey string
}

func New(apiKey string) *PostmanClient {
	return &PostmanClient{client: &http.Client{Timeout: time.Minute},
		apiKey: apiKey}
}

type ErrorInfo struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
type ErrorResponse struct {
	ErrorInfo ErrorInfo `json:"error"`
}

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("Error:%s , Message: %s", err.ErrorInfo.Name, err.ErrorInfo.Message)
}

func ComposeUrl(resource, id string) string {
	path := url.PathEscape(resource) + "/" + url.PathEscape(id)

	return baseURL.ResolveReference(&url.URL{Path: path}).String()
}

//  Get Workspace endpoint
func (c *PostmanClient) GetWorkspace(id string) (WorkspaceObject, error) {

	endpoint := ComposeUrl("workspaces", id)

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Printf("%v", err)
	}

	req.Header.Add("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("%v", err)
	}
	defer resp.Body.Close()

	//
	switch resp.StatusCode {
	case 200:
		var respJson WorkspaceObject
		if err := json.NewDecoder(resp.Body).Decode(&respJson); err != nil {
			log.Printf("%v", err)
			return WorkspaceObject{}, err
		}
		return respJson, nil
	case 400, 401, 404, 500, 503:
		var errResp ErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return WorkspaceObject{}, err
		}

		return WorkspaceObject{}, &errResp

	default:
		return WorkspaceObject{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
}

// Get Environment Endpoint
func (c *PostmanClient) GetEnvironment(id string) (EnvironmentObject, error) {

	endpoint := ComposeUrl("environments", id)

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Printf("%v", err)
	}

	req.Header.Add("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)

	if err != nil {
		log.Printf("%v", err)
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var respJson EnvironmentObject
		if err := json.NewDecoder(resp.Body).Decode(&respJson); err != nil {
			return EnvironmentObject{}, err
		}

		return respJson, nil
	case 400, 401, 404, 500, 503:
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return EnvironmentObject{}, err
		}

		return EnvironmentObject{}, &errResp
	default:
		return EnvironmentObject{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
}
