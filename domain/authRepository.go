package domain

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/amrshaban2005/banking-lib/logger"
)

type AuthRepository interface {
	IsAuthorize(token string, permissionName string, vars map[string]string) bool
}

type RemoteAuthRepository struct{}

type authVerifyResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

func NewRemoteAuthRepository() RemoteAuthRepository {
	return RemoteAuthRepository{}
}

func (auth RemoteAuthRepository) IsAuthorize(token string, permissionName string, vars map[string]string) bool {
	url := buildVerifyURL(token, permissionName, vars)
	logger.Info(url)
	res, err := http.Get(url)
	if err != nil {
		logger.Error("Error while sending verify request " + err.Error())
		return false
	}
	defer res.Body.Close()

	var body authVerifyResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		logger.Error("decoding response from auth server: " + err.Error())
		return false
	}
	return body.IsAuthorized
}

func buildVerifyURL(token string, permissionName string, vars map[string]string) string {
	u := url.URL{Host: "localhost:8081", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	q.Add("permissionName", permissionName)
	for key, value := range vars {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
