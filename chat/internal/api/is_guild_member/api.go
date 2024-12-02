package is_guild_member

import (
	"chat-ms/internal/api"
	"encoding/json"
	"io"
	"logging"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Guild struct {
	GuildName string `json:"guild_name"`
	GuildId   string `json:"guild_id"`
}

type Response struct {
	AccessibleGuilds []Guild `json:"accessible_guilds"`
}

func IsGuildMember(guild, accessToken string, logging *logging.Logger) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, api.AUTH_URL+"/auth/get-chat-access", nil)
	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not create request")
		return false, err
	}

	req.Close = true

	req.Header.Add("Authorization", "Bearer "+accessToken)

	var client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	res, err := client.Do(req)
	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error making http request")
		return false, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not read response body")
		return false, err
	}

	var resp Response

	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		logging.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not unmarshal response body")
		return false, err

	}

	isMember := false

	for _, accessibleGuild := range resp.AccessibleGuilds {
		if accessibleGuild.GuildName == guild {
			isMember = true
		}
	}

	if !isMember {
		return false, err
	}
	return isMember, nil
}
