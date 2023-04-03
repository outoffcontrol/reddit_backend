package utils

import (
	"encoding/json"
	"reddit_backend/pkg/logging"

	//"reddit_backend/pkg/logging"
	"net/http"
)

func HandleErr(w http.ResponseWriter, r *http.Request, err error, statusCode int, logger *logging.Logger, msg ...Return422Error) bool {
	//ctx := r.Context()
	if err == nil {
		return false
	}
	var resp []byte
	var e error
	if statusCode == http.StatusUnprocessableEntity {
		resp, e = json.Marshal(&Return422Errors{Errors: msg})
	} else {
		resp, e = json.Marshal(ReturnMsg{Msg: err.Error()})
	}
	if e != nil {
		logger.Z(r.Context()).Error(e)
		//logging.Z(ctx).Error(e)
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	w.WriteHeader(statusCode)
	w.Write(resp)
	return true
}

func HandleResult(w http.ResponseWriter, r *http.Request, resp any, statusCode int, logger *logging.Logger) {
	result, err := json.Marshal(resp)
	if err != nil {
		logger.Z(r.Context()).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(result)
}

func HandleErrDecodeReq(w http.ResponseWriter, r *http.Request, decodeObj any) bool {
	err := json.NewDecoder(r.Body).Decode(&decodeObj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}
	return false
}
