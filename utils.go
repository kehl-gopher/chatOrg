package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"telex-chat/internal/data"
	"telex-chat/internal/models"
)

var ErrSettings = errors.New("settings field is required")
var ErrAuthorization = errors.New("authorization label is required")
var ErrSettingsField = errors.New("field label cannot be empty")

type toJson map[string]interface{}

func (app *application) ProcessSettings(jsonD Query) (*data.Company, error) {

	var authToken string

	if jsonD.Settings == nil {
		return nil, ErrSettings
	}

	for _, settings := range jsonD.Settings {

		if settings.Label != "Authorization" {
			return nil, ErrAuthorization
		}

		if settings.Default == "" {
			return nil, ErrSettingsField
		} else {
			authToken = settings.Default
			break
		}

	}
	return app.VerifyAPIKey(authToken)
}

func ProcessMessage(message string) string {
	m := strings.TrimPrefix(message, "<p>")
	return strings.TrimSuffix(m, "</p>")
}
func ReadJson(r *http.Request, toStruct interface{}) error {
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(toStruct)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidMarshalError *json.InvalidUnmarshalError

		// check errors type
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON at (%d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type character %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body cannot be empty")
		case errors.As(err, &invalidMarshalError):
			panic(err)
		default:
			return err
		}

	}
	return nil
}

func (app *application) writeResponse(w http.ResponseWriter, statusCode int, message toJson) (int, error) {

	byt, err := writeToJson(message)

	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(byt)
}

// write to json...
func writeToJson(value interface{}) ([]byte, error) {
	byte, err := json.Marshal(value)

	if err != nil {
		return nil, err
	}
	return byte, err
}

func (app *application) VerifyAPIKey(apiKey string) (*data.Company, error) {
	com, err := app.model.Model.GetAPIKey(apiKey)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAPiKey):
			return nil, err
		default:
			return nil, err
		}
	}
	return com, nil
}
