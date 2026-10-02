package bot_manager

import (
	"encoding/json"
	"fmt"
	service_error "go-pkg/helper/error_helper"
	"strings"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/tnminh2412/go-microservice-structure/service/app/edge/social/internal/pkg/error_helper"
)

type CommandResponse func(b *gotgbot.Bot, ctx *ext.Context) error

type CommandOtp struct {
	SendWith string                        `json:"sendWith"`
	Messages map[string]CommandOpt_Message `json:"messages"`
}

type CommandOpt_Message struct {
	Text           string                               `json:"text"`
	MediaUrl       string                               `json:"mediaUrl"`
	InlineKeyboard []*CommandOpt_Message_InlineKeyboard `json:"inlineKeyboard"`
}

type CommandOpt_Message_InlineKeyboard struct {
	Text      string `json:"text"`
	Url       string `json:"url"`
	WebappUrl string `json:"webappUrl"`
}

type BotManager struct {
	mapping     map[string]botData
	mux         sync.RWMutex
	hostMapping map[string]string
}

type botData struct {
	bot                 *gotgbot.Bot
	stopUpdates         chan struct{}
	updateWriterControl *sync.WaitGroup
	secretToken         string
	host                string
	commands            map[string]CommandResponse
}

func New() *BotManager {
	return &BotManager{}
}

func (m *BotManager) AddBot(token string, secretToken string, host string) error {
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		return error_helper.New(service_error.Internal, "gotgbot.NewBot error", err)
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	if m.mapping == nil {
		m.mapping = make(map[string]botData)
	}
	if m.hostMapping == nil {
		m.hostMapping = make(map[string]string)
	}
	if _, ok := m.mapping[b.Token]; ok {
		return error_helper.New(service_error.Duplicate, "bot with token already existed", err)
	}
	if _, ok := m.hostMapping[host]; host != "" && ok {
		return error_helper.New(service_error.Duplicate, "bot with host already existed", err)
	}

	bData := botData{
		bot:                 b,
		stopUpdates:         make(chan struct{}),
		updateWriterControl: &sync.WaitGroup{},
		secretToken:         secretToken,
		host:                host,
		commands:            make(map[string]CommandResponse),
	}

	m.mapping[bData.bot.Token] = bData
	m.hostMapping[bData.host] = bData.bot.Token
	return nil
}

func (m *BotManager) RemoveBot(token string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, ok := m.mapping[token]
	if !ok {
		return false
	}

	delete(m.mapping, token)
	return true
}

func (m *BotManager) NewCommand(token string, c string, opt CommandOtp) error {
	b, ok := m.getBot(token)
	if ok == false {
		return error_helper.New(service_error.NotFound, "bot not found", nil)
	}
	b.commands[strings.ToLower(c)] = getCommandResponse(opt)
	return nil
}

func (m *BotManager) ProcessUpdates(host string, secretToken string, rawReq json.RawMessage) error {
	b, ok := m.getBotFromHost(host)
	if ok == false {
		return error_helper.New(service_error.NotFound, "bot not found", nil)
	}

	b.updateWriterControl.Add(1)
	defer b.updateWriterControl.Done()

	if b.shouldStopUpdates() {
		return error_helper.New(service_error.Unavailable, "stopped update", nil)
	}

	if b.secretToken != secretToken {
		return error_helper.New(service_error.Unauthenticated, "", nil)
	}

	var updateReq gotgbot.Update
	if err := json.Unmarshal(rawReq, &updateReq); err != nil {
		return error_helper.New(service_error.Internal, "json.Unmarshal error", err)
	}

	ctx := ext.NewContext(&updateReq, nil)
	reqCmd := ctx.Message.GetText()[1:]
	res, ok := b.commands[strings.ToLower(reqCmd)]
	if ok == false {

	}
	return res(b.bot, ctx)

}

func (m *BotManager) StopAllBots() {
	for _, bData := range m.mapping {
		bData.stop()
	}
}

func (m *BotManager) getBot(token string) (botData, bool) {
	bData, ok := m.mapping[token]
	return bData, ok
}

func (m *BotManager) getBotFromHost(host string) (botData, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	token, ok := m.hostMapping[host]
	if !ok {
		return botData{}, false
	}
	return m.getBot(token)
}

func (b *botData) shouldStopUpdates() bool {
	select {
	case <-b.stopUpdates:
		// if anything comes in on the closing channel, we know the channel is closed.
		return true
	default:
		// otherwise, continue as usual
		return false
	}
}

func (b *botData) stop() {
	if b.stopUpdates != nil {
		close(b.stopUpdates)
	}

	b.updateWriterControl.Wait()
}

func getCommandResponse(cmd CommandOtp) CommandResponse {
	return func(bot *gotgbot.Bot, ctx *ext.Context) error {
		language := ctx.EffectiveUser.LanguageCode
		messageCfg, ok := cmd.Messages[language]
		if ok == false {
			messageCfg = cmd.Messages["en"]
		}

		switch cmd.SendWith {
		case "photo":
			return sendPhoto(messageCfg, bot, ctx)
		case "video":
			return sendVideo(messageCfg, bot, ctx)
		default:
			return sendMessage(messageCfg, bot, ctx)
		}
	}
}

func sendMessage(messageCfg CommandOpt_Message, bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := bot.SendMessage(ctx.EffectiveChat.Id, messageCfg.Text, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: func() [][]gotgbot.InlineKeyboardButton {
				inlineKeyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(messageCfg.InlineKeyboard))
				for _, k := range messageCfg.InlineKeyboard {
					inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
						{
							Text: k.Text,
							Url:  k.Url,
							WebApp: &gotgbot.WebAppInfo{
								Url: k.WebappUrl,
							},
						},
					})
				}
				return inlineKeyboard
			}(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func sendVideo(messageCfg CommandOpt_Message, bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := bot.SendVideo(ctx.EffectiveChat.Id, gotgbot.InputFileByURL(messageCfg.MediaUrl), &gotgbot.SendVideoOpts{
		ParseMode: "HTML",
		Caption:   messageCfg.Text,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: func() [][]gotgbot.InlineKeyboardButton {
				inlineKeyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(messageCfg.InlineKeyboard))
				for _, k := range messageCfg.InlineKeyboard {
					inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
						{
							Text: k.Text,
							Url:  k.Url,
							WebApp: &gotgbot.WebAppInfo{
								Url: k.WebappUrl,
							},
						},
					})
				}
				return inlineKeyboard
			}(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %v", err)
	}
	return nil
}

func sendPhoto(messageCfg CommandOpt_Message, bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := bot.SendPhoto(ctx.EffectiveChat.Id, gotgbot.InputFileByURL(messageCfg.MediaUrl), &gotgbot.SendPhotoOpts{
		ParseMode: "HTML",
		Caption:   messageCfg.Text,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: func() [][]gotgbot.InlineKeyboardButton {
				inlineKeyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(messageCfg.InlineKeyboard))
				for _, k := range messageCfg.InlineKeyboard {
					inlineKeyboard = append(inlineKeyboard, []gotgbot.InlineKeyboardButton{
						{
							Text: k.Text,
							Url:  k.Url,
							WebApp: &gotgbot.WebAppInfo{
								Url: k.WebappUrl,
							},
						},
					})
				}
				return inlineKeyboard
			}(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
