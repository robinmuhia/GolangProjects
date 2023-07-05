package auth

import (
	"errors"
	"net/http"
	"strings"
)

//GetAPIkey extracts API key from headers in HTTP request

func GetAPIkey(headers http.Header)(string,error){
	apiKey := headers.Get("Authorization")
	if apiKey == ""{
		return "",errors.New("no authentication info found")
	}

	apiKeys := strings.Split(apiKey," ")
	if len(apiKeys) != 2{
		return "",errors.New("malformed authorization header")
	}
	if apiKeys[0] != "ApiKey"{
		return "",errors.New("malformed authorization header does not start with expected value")
	}
	return apiKeys[1],nil
}