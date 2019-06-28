package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	queryArgDryRun = "simulate"
)

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

// WriteSimulationResponse prepares and writes an HTTP
// response for transactions simulations.
func WriteSimulationResponse(w http.ResponseWriter, gas int64) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"gas_estimate":%v}`, gas)))
}

// WriteCustomResponse prepares and writes an HTTP
// response for transactions simulations.
func WriteCustomResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Message":"%v"}`, msg)))
}

// HasDryRunArg returns true if the request's URL query contains
// the dry run argument and its value is set to "true".
func HasDryRunArg(r *http.Request) bool {
	return r.URL.Query().Get(queryArgDryRun) == "true"
}

// ParseFloat64OrReturnBadRequest converts s to a float64 value. It returns a default
// value if the string is empty. Write
func ParseFloat64OrReturnBadRequest(w http.ResponseWriter, s string, defaultIfEmpty float64) (n float64, ok bool) {
	if len(s) == 0 {
		return defaultIfEmpty, true
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return n, false
	}
	return n, true
}

// ResponseBytesToJSON converts the response message to JSON format
func ResponseBytesToJSON(bytea []byte) []byte {
	
	var response map[string]string
	if err := json.Unmarshal(bytea, &response); err != nil {
		panic(err)
	}
	responseStr := response["Response"]
	delete(response, "Response")
	
	tagsStart := strings.Index(responseStr, "Tags:[{")
	tagsEnd := strings.Index(responseStr, "}]")
	tagsStr := responseStr[tagsStart+5 : tagsEnd+2]
	responseStr = responseStr[1:tagsStart] + responseStr[tagsEnd+2:len(responseStr)-1]
	
	tagsStr = strings.Replace(tagsStr, "Key", "\"key\"", -1)
	tagsStr = strings.Replace(tagsStr, "Value", "\"value\"", -1)
	tagsStr = strings.Replace(tagsStr, "} {", "},{", -1)
	tagsStr = strings.Replace(tagsStr, " ", ",", -1)
	tagsStr = ",\"Tags\":" + tagsStr
	
	for _, pair := range strings.Split(responseStr, " ") {
		keyVal := strings.Split(pair, ":")
		if len(keyVal) > 1 {
			response[keyVal[0]] = keyVal[1]
		}
	}
	
	newResponse, _ := json.Marshal(response)
	newResponseStr := string(newResponse)
	newResponseStr = newResponseStr[:len(newResponseStr)-1] + tagsStr + "}"
	
	var responseInterface map[string]interface{}
	if err := json.Unmarshal([]byte(newResponseStr), &responseInterface); err != nil {
		panic(err)
	}
	type TagString struct {
		Key   string
		Value string
	}
	
	var tagList []TagString
	
	for i := range responseInterface["Tags"].([]interface{}) {
		
		tagString := &TagString{}
		tagPair := responseInterface["Tags"].([]interface{})[i]
		tagInterface := tagPair.(map[string]interface{})
		tag := tagString
		keyString := ""
		valueString := ""
		
		switch tagKeyByte := tagInterface["key"].(type) {
		case []interface{}:
			for k := range tagKeyByte {
				keyVal := fmt.Sprintf("%v", tagKeyByte[k])
				keyChar, _ := strconv.Atoi(keyVal)
				keyString += string(keyChar)
			}
			tag.Key = keyString
		}
		
		switch tagValueByte := tagInterface["value"].(type) {
		case []interface{}:
			for k := range tagValueByte {
				valueVal := fmt.Sprintf("%v", tagValueByte[k])
				valueChar, _ := strconv.Atoi(valueVal)
				valueString += string(valueChar)
			}
			tag.Value = valueString
		}
		
		tagList = append(tagList, *tag)
		
	}
	tagFinalValue := tagList
	s := make([]interface{}, len(tagFinalValue))
	for i, v := range tagFinalValue {
		s[i] = v
	}
	responseInterface["Tags"] = s
	data, err := json.Marshal(responseInterface)
	if err != nil {
		panic(err)
	}
	return data
}
