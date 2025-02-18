package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type CallRequest struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Body    []byte            `json:"body"`
	Headers map[string]string `json:"headers"`
}

type CallResponse struct {
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func (s *Service) MakeRequestHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriteJSONError(w, "failed to read body", 500)
		return
	}
	var requestJson CallRequest
	if err := json.Unmarshal(body, &requestJson); err != nil {
		WriteJSONError(w, "error while unmarshall json", 500)
		return
	}
	method := strings.ToUpper(requestJson.Method)

	bodyReader := bytes.NewReader(requestJson.Body)
	userRequest, err := http.NewRequest(method, requestJson.Url, bodyReader)
	if err != nil {
		WriteJSONError(w, "error while creating request", 500)
		return
	}
	client := &http.Client{}
	response, err := client.Do(userRequest)
	if err != nil {
		WriteJSONError(w, "error while making request", 500)
		return
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		return
	}
	responseJson := CallResponse{
		Headers: response.Header,
		Body:    buf.String(),
	}
	data, err := json.Marshal(responseJson)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		return
	}
}
