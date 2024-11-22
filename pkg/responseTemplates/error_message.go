package responseTemplates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mainService/pkg/serverErrors"
)

type errorToSend struct {
	Message string `json:"message"`
}

func SendErrorMessage(w http.ResponseWriter, err error, statusCode int) error {
	errResp := errorToSend{Message: err.Error()}
	mess, marshalErr := json.Marshal(errResp)
	if marshalErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(fmt.Sprintf(`{"message":%s}`, serverErrors.INTERNAL_SERVER_ERROR.Error())))
		if writeErr != nil {
			return marshalErr
		}

		return marshalErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(mess)
	if writeErr != nil {
		sendErr := SendErrorMessage(w, serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		if sendErr != nil {
			return writeErr
		}
	}

	return nil
}
