package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker service is up and running",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	err := app.readJson(w, r, &payload)
	if err != nil {
		_ = app.writeError(w, err, http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "auth":
		app.authenticate(w, payload.Auth)
	case "log":
		app.logItem(w, payload.Log)
	default:
		_ = app.writeError(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, auth AuthPayload) {
	jsonData, _ := json.MarshalIndent(auth, "", "\t")

	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var payload jsonResponse
	_ = json.NewDecoder(response.Body).Decode(&payload)
	_ = app.writeJson(w, response.StatusCode, payload)
}

func (app *Config) logItem(w http.ResponseWriter, log LogPayload) {
	jsonData, _ := json.MarshalIndent(log, "", "\t")
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.writeError(w, err, http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.writeError(w, err, http.StatusInternalServerError)
	}
	defer response.Body.Close()
	_ = app.writeJson(w, response.StatusCode, response.Body)
}
