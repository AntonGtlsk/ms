package get_user_chat_access

import (
	resp "auth-ms/internal/lib/api/responce"
	"auth-ms/internal/service/auth"
	"logging"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

const (
	discordClientID     = ""
	discordClientSecret = ""
	discordRedirectURL  = "http://localhost:3000/login"
)

var oauthConfig = oauth2.Config{
	ClientID:     discordClientID,
	ClientSecret: discordClientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	},
	RedirectURL: discordRedirectURL,
	Scopes:      []string{"identify"},
}

type Request struct {
	UserId string `json:"userid"`
}

type Guild struct {
	GuildName string `json:"guild_name"`
	GuildId   string `json:"guild_id"`
}

type Response struct {
	AccessibleGuilds []Guild `json:"accessible_guilds"`
}

type GuildRoleMapping struct {
	GuildId   string   `yaml:"guildId"`
	GuildName string   `yaml:"guildName"`
	RoleId    []string `yaml:"roleId"`
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

func New(logging *logging.Logger, dg *discordgo.Session, jwtParser *auth.JwtParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		_, uid, _, _, err := jwtParser.ParseToken(tokenString)

		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		var request Request

		if err := validator.New().Struct(request); err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("Failed to validate data")

			validateErr := err.(validator.ValidationErrors)

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		fileContent, err := os.ReadFile("config/guildAndRoleId.yaml")
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error reading file")
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		var guildRoleMapping []GuildRoleMapping

		err = yaml.Unmarshal(fileContent, &guildRoleMapping)
		if err != nil {
			logging.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error unmarshal file")
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		var resultAccess Response
		resultAccess.AccessibleGuilds = []Guild{}
		for _, entry := range guildRoleMapping {
			member, err := dg.GuildMember(entry.GuildId, uid)
			if err != nil {
				logging.WithFields(logrus.Fields{
					"error": err,
				}).Error("Error getting guild member")
				http.Error(w, "Error", http.StatusInternalServerError)
				return
			}

			for _, memberRole := range member.Roles {
				for _, accessibleRoles := range entry.RoleId {
					if memberRole == accessibleRoles {
						resultAccess.AccessibleGuilds = append(resultAccess.AccessibleGuilds, Guild{GuildName: entry.GuildName, GuildId: entry.GuildId})
					}
				}
			}
		}

		render.JSON(w, r, resultAccess)
	}
}
