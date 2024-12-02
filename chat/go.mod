module chat-ms

go 1.22.0

toolchain go1.22.5

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/ethereum/go-ethereum v1.13.4
	github.com/go-playground/validator/v10 v10.16.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/sirupsen/logrus v1.9.3
)

require logging v0.3.1

replace logging => ../logging

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
