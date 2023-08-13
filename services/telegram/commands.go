package telegram

import "strings"

type Command string

const (
	CreateImageCommand Command = "/createimage"
	EditCommand        Command = "/edit"
	None               Command = ""

	MaxMessageLength = 12
)

func GetCommand(text *string) Command {
	initialText := (*text)[0:MaxMessageLength]
	command := strings.Split(initialText, " ")[0]

	switch command {
	case string(CreateImageCommand):
		return CreateImageCommand
	case string(EditCommand):
		return EditCommand
	}
	return None
}

func ParseMessage(cmd Command, text *string) (*string, *string) {
	var instruction string
	parsedText := (*text)[len(cmd):]
	parsedText = strings.TrimSpace(parsedText)

	if cmd == EditCommand {
		chunks := strings.Split(parsedText, ":")
		if len(chunks) > 1 {
			instruction = chunks[0]
			parsedText = strings.TrimSpace(chunks[1])
		}
	}

	return &parsedText, &instruction
}
