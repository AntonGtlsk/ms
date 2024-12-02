package chat_history

import (
	api "chat-ms/internal/api/is_guild_member"
	"chat-ms/internal/dao"
	"encoding/json"
	"logging"
	"net/http"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

type Request struct {
	ContractAddress common.Address `json:"contract_address"`
}

func NewOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func New(logging *logging.Logger, messageGetter dao.MessageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		authHeader := r.Header.Get("Authorization")
		guild := r.URL.Query().Get("guild")
		if authHeader == "" {
			http.Error(w, "invalid parameters", http.StatusInternalServerError)
			return
		}
		params := strings.Split(authHeader, " ")
		if len(params) < 1 {
			http.Error(w, "invalid parameters", http.StatusInternalServerError)
			return
		}
		isMember, err := api.IsGuildMember(guild, params[1], logging)

		if !isMember || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var request Request

		request.ContractAddress = common.HexToAddress(r.URL.Query().Get("contractAddress"))
		messages, err := messageGetter.GetGuildMessagesByAddress(guild, request.ContractAddress)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Error getting messages from db for address: %s", request.ContractAddress)
			http.Error(w, "Error", http.StatusUnauthorized)
			return
		}

		sort.SliceStable(messages, func(i, j int) bool {

			if messages[i].Time.Compare(messages[j].Time) == 1 {
				return false
			}
			if messages[i].Time.Compare(messages[j].Time) == -1 {
				return true
			}

			return false
		})

		responseJson, err := json.Marshal(messages)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Could not marshal messages to json")
			return
		}

		_, err = w.Write(responseJson)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Could not write data to response")
			return
		}

	}
}
