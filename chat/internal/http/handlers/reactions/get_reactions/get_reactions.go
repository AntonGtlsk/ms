package getReactions

import (
	api "chat-ms/internal/api/is_guild_member"
	r "chat-ms/internal/entity/reactions"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"logging"
	"net/http"
	"strings"
)

type Request struct {
	ContractAddress common.Address `json:"contract_address"`
	UserId          string         `json:"userid"`
	Guild           string         `json:"guild"`
}

type ReactionsGetter interface {
	GetReactions(userId, guild string, contractAddress common.Address) (r.Reactions, error)
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

func New(logging *logging.Logger, reactionsGetter ReactionsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		var request Request

		authHeader := r.Header.Get("Authorization")
		request.Guild = r.URL.Query().Get("guild")
		request.UserId = r.URL.Query().Get("userid")
		request.ContractAddress = common.HexToAddress(r.URL.Query().Get("contractAddress"))

		if request.UserId == "" || request.Guild == "" || authHeader == "" || r.URL.Query().Get("contractAddress") == "" {
			http.Error(w, "invalid parameters", http.StatusInternalServerError)
			return
		}

		params := strings.Split(authHeader, " ")
		if len(params) < 1 {
			http.Error(w, "invalid parameters", http.StatusInternalServerError)
			return
		}
		isMember, err := api.IsGuildMember(request.Guild, params[1], logging)

		if !isMember || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		reactions, err := reactionsGetter.GetReactions(request.UserId, request.Guild, request.ContractAddress)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Error getting reactions from db for guild: %s, address: %s, userId: %s", request.ContractAddress, request.Guild, request.UserId)
			http.Error(w, "Error", http.StatusUnauthorized)
			return
		}

		responseJson, err := json.Marshal(reactions)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Could not marshal reactions to json")
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
