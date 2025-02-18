package antenna

import (
	"errors"
	"fmt"
	"log"

	"github.com/shomorish/antenna/internal/rss"
)

const (
	CMD_ADD                 = "a"
	CMD_EDIT                = "e"
	CMD_DELETE              = "d"
	CMD_VIEW                = "v"
	CMD_UPDATE              = "u"
	CMD_CHANGE_EMAIL_SENDER = "c"
	CMD_HELP                = "h"
	CMD_QUIT                = "q"
)

type Container struct {
	Conf *rss.Config
	Help string
}

type Command interface {
	Exec(container *Container) error
}

// 追加コマンド
type AddCommand struct{}

func (a *AddCommand) Exec(container *Container) error {
	if err := rss.AddRSSFeed(container.Conf); err != nil {
		fmt.Printf("Failed to add RSS feed: %v\n", err)
	}
	readConfigMayBeInterrupted(container.Conf)
	return nil
}

var _ Command = new(AddCommand)

// 編集コマンド
type EditCommand struct{}

func (e *EditCommand) Exec(container *Container) error {
	if err := rss.EditRSSFeed(container.Conf); err != nil {
		fmt.Printf("Failed to edit RSS feed: %v\n", err)
	}
	readConfigMayBeInterrupted(container.Conf)
	return nil
}

var _ Command = new(EditCommand)

// 削除コマンド
type DeleteCommand struct{}

func (d *DeleteCommand) Exec(container *Container) error {
	if err := rss.DeleteRSSFeed(container.Conf); err != nil {
		fmt.Printf("Failed to delete RSS feed: %v\n", err)
	}
	readConfigMayBeInterrupted(container.Conf)
	return nil
}

var _ Command = new(DeleteCommand)

// 表示コマンド
type ViewCommand struct{}

func (v *ViewCommand) Exec(container *Container) error {
	fmt.Println(container.Conf.String())
	return nil
}

var _ Command = new(ViewCommand)

// 更新コマンド
type UpdateCommand struct{}

func (u *UpdateCommand) Exec(container *Container) error {
	rss.UpdateRSSFeeds(container.Conf)
	readConfigMayBeInterrupted(container.Conf)
	return nil
}

var _ Command = new(UpdateCommand)

// メール送信元変更コマンド
type ChangeEmailSenderCommand struct{}

func (c *ChangeEmailSenderCommand) Exec(container *Container) error {
	if err := rss.SetEmailSender(container.Conf); err != nil {
		log.Fatalln(err.Error())
	}
	readConfigMayBeInterrupted(container.Conf)
	return nil
}

var _ Command = new(ChangeEmailSenderCommand)

// ヘルプコマンド
type HelpCommand struct{}

func (h *HelpCommand) Exec(container *Container) error {
	fmt.Println(container.Help)
	return nil
}

var _ Command = new(HelpCommand)

// 終了コマンド
type QuitCommand struct{}

func (q *QuitCommand) Exec(container *Container) error {
	return errors.New("quit")
}

var _ Command = new(QuitCommand)
