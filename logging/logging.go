package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	BlueColor   = "\033[34m"
	YellowColor = "\033[33m"
	GreenColor  = "\033[32m"
	ResetColor  = "\033[0m"
)

type writerHook struct {
	Writer     []io.Writer
	LogLevels  []logrus.Level
	ChainColor string
}

type discordHook struct {
	WebhookURL string
	LogLevels  []logrus.Level
	ds         *discordgo.Session
}

func (h *discordHook) Fire(entry *logrus.Entry) error {
	fields := []*discordgo.MessageEmbedField{}

	for key, value := range entry.Data {

		var correctValue string

		switch v := value.(type) {
		case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
			correctValue = fmt.Sprintf("%d", v)
		case string:
			correctValue = v
		case float64, float32:
			correctValue = fmt.Sprintf("%f", v)
		case error:
			correctValue = v.Error()
		default:
			correctValue = ""
		}

		if len(correctValue) > 1000 {
			correctValue = correctValue[:1000] + "..."
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   key,
			Value:  correctValue,
			Inline: true,
		})
	}

	if entry.Caller.Func != nil {
		file, line := entry.Caller.Func.FileLine(entry.Caller.PC)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "File",
			Value: fmt.Sprintf("%s:%d", file, line),
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Time",
		Value: time.Now().Format("Mon, 2 Jan 2006 15:04:05"),
	})

	title := entry.Message
	if len(title) > 250 {
		title = title[:200] + "..."
	}

	params := &discordgo.WebhookParams{
		Username: entry.Level.String(),
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  title,
				Fields: fields,
			},
		},
	}

	h.WebhookURL = strings.TrimPrefix(h.WebhookURL, "https://discord.com/api/webhooks/")
	_, err := h.ds.WebhookExecute(strings.Split(h.WebhookURL, "/")[0], strings.Split(h.WebhookURL, "/")[1], false, params)

	return err
}

func (h *discordHook) Levels() []logrus.Level {
	return h.LogLevels
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	title := entry.Message
	if len(title) > 250 {
		title = title[:200] + "..."
	}

	data := entry.Data

	stringFields := "{"

	for a, b := range data {

		correctValue := ""

		switch v := b.(type) {
		case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
			correctValue = fmt.Sprintf("%d", v)
		case string:
			correctValue = v
		case float64, float32:
			correctValue = fmt.Sprintf("%f", v)
		case error:
			correctValue = v.Error()
		default:
			correctValue = ""
		}

		stringFields += fmt.Sprintf("'%s: %s', ", a, correctValue)
	}

	stringFields += "}"

	time := entry.Time.Format("2 Jan 2006 15:04:05")

	for _, w := range hook.Writer {
		var logLine string

		if entry.Caller.Func != nil {
			file, line := entry.Caller.Func.FileLine(entry.Caller.PC)
			logLine = fmt.Sprintf("%s time=%s; [%s] message=[%s]; fields=%s; file=[ %s ]; %s \n", hook.ChainColor, time, entry.Level.String(), title, stringFields, fmt.Sprintf("%s:%d", file, line), ResetColor)

		} else {
			logLine = fmt.Sprintf("%s [%s] message=[%s]; fields=%s; time=%s. %s\n", hook.ChainColor, entry.Level.String(), title, stringFields, time, ResetColor)
		}
		w.Write([]byte(logLine))
	}
	return nil
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=LoggerInterface
type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Logger struct {
	Logrus              *logrus.Entry
	UndefinedWebhookURL string
	ds                  *discordgo.Session
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Logrus.Debugf(format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Logrus.Infof(format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.Logrus.Warnf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Logrus.Errorf(format, args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.Logrus.Fatalf(format, args...)
}

func (l Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logrus.WithFields(fields)
}

func (l Logger) UndefinedStandard(address, standard string, err error) {
	params := &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{Title: "New undefined standard",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Address",
						Value: address,
					},
					{
						Name:  "Standard",
						Value: standard,
					},
					{
						Name:  "Error",
						Value: fmt.Sprintln(err),
					},
					{
						Name:  "Etherscan",
						Value: "https://etherscan.io/address/" + address,
					},
				},
			},
		},
	}

	webhook := strings.TrimPrefix(l.UndefinedWebhookURL, "https://discord.com/api/webhooks/")
	l.ds.WebhookExecute(strings.Split(webhook, "/")[0], strings.Split(webhook, "/")[1], false, params)
}

type TestLogger struct{}

func (l TestLogger) Debugf(format string, args ...interface{}) {
}

func (l TestLogger) Infof(format string, args ...interface{}) {
}

func (l TestLogger) Warnf(format string, args ...interface{}) {
}

func (l TestLogger) Errorf(format string, args ...interface{}) {
}

func (l TestLogger) Fatalf(format string, args ...interface{}) {
}

func (l TestLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return nil
}

func New(botToken, folder string, filenames []map[string][]string, webhookUrls []map[string][]string, undefinedWebhookURl string, chain string) *Logger {
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if err != nil {
		log.Fatalf("error connect to discord bot %v", err)
	}

	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	err = os.MkdirAll("logs/"+folder+"/", 0777)
	if err != nil {
		l.Errorf("Error calling MkdirAll method, err: %s ", err)
	}

	l.SetOutput(io.Discard)

	for _, filename := range filenames {
		for file, lvls := range filename {
			levels := []logrus.Level{}
			outputs := []io.Writer{}

			if file == "stdout" {
				outputs = append(outputs, os.Stdout)
			} else {
				currentFile, err := os.OpenFile("logs/"+folder+"/"+file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
				if err != nil {
					l.Errorf("Error calling ReadFile method, err: %s file: logging.go:167", err)
				}

				outputs = append(outputs, currentFile)
			}
			for _, level := range lvls {
				switch level {
				case "debug":
					levels = append(levels, logrus.DebugLevel)
				case "trace":
					levels = append(levels, logrus.TraceLevel)
				case "info":
					levels = append(levels, logrus.InfoLevel)
				case "warn":
					levels = append(levels, logrus.WarnLevel)
				case "error":
					levels = append(levels, logrus.ErrorLevel)
				case "fatal":
					levels = append(levels, logrus.FatalLevel)
				case "all":
					{
						levels = logrus.AllLevels
						break
					}
				}
			}

			if chain == "mainnet" || chain == "tron" {
				l.AddHook(&writerHook{
					Writer:     outputs,
					LogLevels:  levels,
					ChainColor: YellowColor,
				})
			}

			if chain == "base" {
				l.AddHook(&writerHook{
					Writer:     outputs,
					LogLevels:  levels,
					ChainColor: BlueColor,
				})
			}

			if chain == "all" {
				l.AddHook(&writerHook{
					Writer:     outputs,
					LogLevels:  levels,
					ChainColor: GreenColor,
				})
			}

		}
	}

	for _, webhookUrl := range webhookUrls {
		for hook, levels := range webhookUrl {
			discordNotificationsLevels := []logrus.Level{}

			for _, level := range levels {
				switch level {
				case "debug":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.DebugLevel)
				case "trace":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.TraceLevel)
				case "info":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.InfoLevel)
				case "warn":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.WarnLevel)
				case "error":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.ErrorLevel)
				case "fatal":
					discordNotificationsLevels = append(discordNotificationsLevels, logrus.FatalLevel)
				case "all":
					{
						discordNotificationsLevels = logrus.AllLevels
						break
					}
				}
			}

			l.AddHook(&discordHook{
				WebhookURL: hook,
				LogLevels:  discordNotificationsLevels,
				ds:         discord,
			})
		}
	}

	l.SetLevel(logrus.TraceLevel)

	return &Logger{Logrus: logrus.NewEntry(l), UndefinedWebhookURL: undefinedWebhookURl, ds: discord}
}
