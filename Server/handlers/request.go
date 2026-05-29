package handlers

import (
	"encoding/json"
	"net/http"
	"proxyScanner/Model"
	"proxyScanner/Server/Message"
	"strconv"

	"github.com/gorilla/mux"
)

func GetRequestById(writer http.ResponseWriter, request *http.Request) {
	requestId, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		errorResponse := Message.CreateErrorMessage("id is invalid", http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)

	} else {

		result := Model.GetRequestById(uint(requestId))
		json.NewEncoder(writer).Encode(result)
	}
}

//func GetRequestByMethod(writer http.ResponseWriter, request *http.Request) {
//
//}

func GetRequestByMethod(writer http.ResponseWriter, request *http.Request) {

	methodName := mux.Vars(request)["method_name"]

	result := Model.FindRequestByMethod(methodName)
	json.NewEncoder(writer).Encode(result)

}
