package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		fmt.Println(err)
		app.writeError(w, errors.New("invalid credentials - mail"), http.StatusNotFound)
		return
	}

	valid, err := app.Repo.ValidatePassword(requestPayload.Password, user)
	if err != nil || !valid {
		app.writeError(w, errors.New("invalid credentials - password"), http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authenticate", user.Email)
	if err != nil {
		app.writeError(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Successfully authenticated",
		Data:    user,
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	var logEntry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	logEntry.Name = name
	logEntry.Data = data

	jsonData, _ := json.MarshalIndent(logEntry, "", "\t")

	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
