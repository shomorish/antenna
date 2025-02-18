package antenna

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shomorish/antenna/internal/rss"
)

func getContentFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(content), "\n"), nil
}

func readConfigFromToml(filename string, config *rss.Config) error {
	// ファイルが存在しないか確認
	if _, err := os.Stat(filename); err != nil {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}

	// ファイルから設定情報を取得
	_, err := toml.DecodeFile(filename, config)
	if err != nil {
		return err
	}

	return nil
}

func readConfigMayBeInterrupted(config *rss.Config) {
	filename := "rss_feeds.toml"
	if err := readConfigFromToml(filename, config); err != nil {
		log.Fatalf("Failed to read rss_feeds.toml: %v\n", err)
	}
}

func Run() {
	logo, err := getContentFromFile("logo.txt")
	if err != nil {
		log.Fatalf("Failed to get log from file: %v\n", err)
	}
	fmt.Println(logo)
	fmt.Println()

	help, err := getContentFromFile("help.txt")
	if err != nil {
		log.Fatalf("Failed to get help from file: %v\n", err)
	}
	fmt.Println(help)
	fmt.Println()

	var config rss.Config
	readConfigMayBeInterrupted(&config)

	// 送信元メールアドレスがない場合は設定
	if len(config.EmailSender) == 0 {
		if err := rss.SetEmailSender(&config); err != nil {
			log.Fatalln(err.Error())
		}
	}

	container := Container{
		Conf: &config,
		Help: help,
	}

	commands := map[string]Command{
		CMD_ADD:                 &AddCommand{},
		CMD_EDIT:                &EditCommand{},
		CMD_DELETE:              &DeleteCommand{},
		CMD_VIEW:                &ViewCommand{},
		CMD_UPDATE:              &UpdateCommand{},
		CMD_CHANGE_EMAIL_SENDER: &ChangeEmailSenderCommand{},
		CMD_HELP:                &HelpCommand{},
		CMD_QUIT:                &QuitCommand{},
	}

	var quit error
	for quit == nil {
		var command string
		fmt.Print(">>> ")
		fmt.Scanln(&command)

		if len(command) == 0 {
			continue
		}

		c, exists := commands[command]
		if exists {
			quit = c.Exec(&container)
			if quit == nil {
				fmt.Println()
			}
		} else {
			fmt.Println("An invalid command was entered.")
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println("See you.")
}
