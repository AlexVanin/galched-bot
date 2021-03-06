package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type (
	MessageHandler interface {
		Signature() string
		Description() string

		IsValid(string) bool
		Handle(s *discordgo.Session, m *discordgo.MessageCreate)
	}

	HandlerProcessor struct {
		appVersion string
		handlers   []MessageHandler
	}
)

func NewProcessor(appVersion string) *HandlerProcessor {
	return &HandlerProcessor{
		appVersion: appVersion,
	}
}

func (h *HandlerProcessor) AddHandler(handler MessageHandler) {
	h.handlers = append(h.handlers, handler)
}

func (h *HandlerProcessor) HelpMessage() string {
	helpMessage := fmt.Sprintf("ГалчедБот версии %s. Описание команд: \n"+
		"  **!galched** - список команд бота\n", h.appVersion)

	for i := range h.handlers {
		helpMessage += fmt.Sprintf("  **%s** - %s\n",
			h.handlers[i].Signature(),
			h.handlers[i].Description())
	}

	return strings.TrimRight(helpMessage, "\n")
}

func (h *HandlerProcessor) Process(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!galched") {
		LogMessage(m)
		SendMessage(s, m, h.HelpMessage())
		return
	}

	for i := range h.handlers {
		if h.handlers[i].IsValid(m.Content) {
			h.handlers[i].Handle(s, m)
			break
		}
	}
}
