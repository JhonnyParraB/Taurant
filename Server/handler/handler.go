package handler

import (
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type errorMessage struct {
	Code    string `endpoint:"code,omitempty"`
	Details string `endpoint:"details,omitempty"`
}

type responseMessage struct {
	Details string `endpoint:"details,omitempty"`
	ID      string `endpoint:"id,omitempty"`
}

const wrongEmailParamFormatCode = "WrongEmailParamFormat"
const wrongDateParamFormatCode = "WrongDateParamFormat"
const queueForLoadDataJobFullCode = "QueueForLoadDataJobFull"
const internalServerErrorCode = "InternalServerError"
const buyerNotFoundCode = "BuyerNotFound"

var errorsDetails = map[string]string{
	wrongEmailParamFormatCode:   "The email param has a wrong format",
	wrongDateParamFormatCode:    "The date param has a wrong format. Must be in UNIX time format.",
	queueForLoadDataJobFullCode: "The queue of jobs for load data jobs is full. Please try again later.",
	internalServerErrorCode:     "Internal server error has ocurred",
	buyerNotFoundCode:           "The buyer was not found",
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var endpointCaseJSON = jsoniter.Config{TagKey: "endpoint"}.Froze()
	response, _ := endpointCaseJSON.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func handleInternalServerError(err error, w http.ResponseWriter) {
	log.Println(errorTag, err)
	respondwithJSON(w, http.StatusInternalServerError, errorMessage{
		Code:    internalServerErrorCode,
		Details: errorsDetails[internalServerErrorCode],
	})
}
