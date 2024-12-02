package setDislike

import (
	api "chat-ms/internal/api/is_guild_member"
	"encoding/json"
	"logging"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

type Request struct {
	ContractAddress common.Address `json:"contract_address"`
	UserId          string         `json:"userid"`
	Guild           string         `json:"guild"`
	Value           int            `json:"value"`
}

type ReactionsSetter interface {
	SetReaction(userId, guild string, value int, contractAddress common.Address) error
}

type Response struct {
	Status string `json:"status"`
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

func New(logging *logging.Logger, reactionsSetter ReactionsSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		var request Request

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if request.UserId == "" || request.Guild == "" || authHeader == "" || request.ContractAddress.String() == "0x0000000000000000000000000000000000000000" {
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

		err = reactionsSetter.SetReaction(request.UserId, request.Guild, request.Value, request.ContractAddress)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Error setting reactions from db for guild: %s, address: %s, userId: %s", request.ContractAddress, request.Guild, request.UserId)
			http.Error(w, "Error", http.StatusUnauthorized)
			return
		}

		responseJson, err := json.Marshal(Response{
			Status: "OK",
		})

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Could not marshal reponse to json")
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
