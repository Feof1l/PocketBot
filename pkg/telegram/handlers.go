package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart           = "start"
	replyStartTemplate     = "Привет! Чтобы сохранять ссылки в своем Pocket аккаунте,для начала  тебе необходимо дать мне на это доступ. Для этого переходи по ссылке: \n%s"
	replyAlreadyAuthorized = "Ты уже авторизирован.Присылай ссылку,а я ее сохраню"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Это невалидная ссылка!"
		_, err = b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "Ты не авторизирован!Иcпользуй команду /start"
		_, err = b.bot.Send(msg)
		return err

	}
	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "Увы,не удалось сохранить ссылку.Попробуй еще раз позже!"
		_, err = b.bot.Send(msg)
		return err
	}
	//msg.ReplyToMessageID = update.Message.MessageID

	_, err = b.bot.Send(msg)
	return err
}
func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)

	_, err = b.bot.Send(msg)
	return err

}
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды")
	_, err := b.bot.Send(msg)
	return err
}

/*func (b *Bot) handleMachineCommand(message *tgbotapi.Message) error {
	cpu, _ := ghw.CPU()
	memory, _ := ghw.Memory()
	topology, _ := ghw.Topology()
	messageText := "информация о процессорах в хост-системе:" + cpu.String() + "информация о оперативной памяти в хост-системе" +
		memory.String() + "информация об архитектуре хост-компьютера (NUMA против SMP), расположении узлов NUMA хоста и кэшах памяти, специфичных для процессора" +
		topology.String()
	msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
	_, err := b.bot.Send(msg)
	return err
}*/
