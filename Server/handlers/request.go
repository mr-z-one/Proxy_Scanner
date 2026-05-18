package handlers

import (
	"encoding/json"
	"net/http"
	"proxyScanner/Model"
	"strconv"

	"github.com/gorilla/mux"
)

func GetRequestById(writer http.ResponseWriter, request *http.Request) {
	requestId, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]interface{}{
			"statusCode": http.StatusBadRequest,
			"data": map[string]string{
				"message": "id is invalid",
			},
		}
		json.NewEncoder(writer).Encode(errorResponse)

	} else {

		result := Model.GetRequestById(uint(requestId))
		json.NewEncoder(writer).Encode(result)
	}
}
