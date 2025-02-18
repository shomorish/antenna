package rss

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/shomorish/antenna/internal/email"
	"github.com/shomorish/antenna/internal/input"
)

func overwriteConfig(filename string, config *Config) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(*config); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func SetEmailSender(config *Config) error {
	var email string
	var password string
	var host string

	for {
		fmt.Println("Enter the email address you want to use as the sender.")
		email = input.InputEmail()
		if len(email) > 0 {
			break
		}
	}

	for {
		password = input.InputPassword()
		if len(password) > 0 {
			break
		}
	}

	for {
		host = input.InputHost()
		if len(host) > 0 {
			break
		}
	}

	config.EmailSender = email
	config.Password = password
	config.Host = host

	if err := overwriteConfig("rss_feeds.toml", config); err != nil {
		return errors.New("failed to overwrite RSS feed, " + err.Error())
	}

	return nil
}

func AddRSSFeed(config *Config) error {
	var title string
	var url string
	var emails []string

	// フィードのタイトルの入力
	for {
		title = input.InputTitle()
		if len(title) == 0 {
			fmt.Println("The title must be at least one character.")
		} else {
			// 同じタイトルのRSSフィードが在るか確認
			exist := false
			for _, v := range config.FeedInfos {
				if v.Title == title {
					exist = true
					fmt.Printf("`%s` already exist.\n", title)
					break
				}
			}
			if !exist {
				break
			}
		}
	}

	// フィードのURLの入力
	for {
		url = input.InputURL()
		if len(url) == 0 {
			fmt.Println("The URL must be at least one character.")
		} else {
			break
		}
	}

	// 通知に使うメールアドレスの入力
	for {
		emails = input.InputEmails()
		if len(emails) > 0 {
			// メールを送信するか確認
			if input.IsEnteredY(input.YOrOther("Do you want to check if you can send email?")) {
				if err := email.SendEmails(config.EmailSender, config.Password, config.Host, emails, []byte("Test - antenna")); err != nil {
					fmt.Printf("Failed to send email: %v\n", err)
					continue
				}
				// メールが受信できたか確認
				if input.IsEnteredY(input.YOrOther("Did you receive the email?")) {
					// メールの受信が確認できたため、メールアドレスの入力を完了する
					break
				} else {
					fmt.Println("Please enter your email address again.")
				}
			} else {
				// メールを送信しない場合は、メールアドレスの入力を完了する
				break
			}
		} else {
			// メールアドレスが入力されなかった場合は、メールアドレスの入力を完了する
			break
		}
	}

	// フィードを保存する
	newFeed := FeedInfo{
		Title:       title,
		URL:         url,
		Email:       emails,
		LastUpdated: time.Now(),
	}
	if config.FeedInfos == nil {
		config.FeedInfos = []FeedInfo{newFeed}
	} else {
		config.FeedInfos = append(config.FeedInfos, newFeed)
	}
	if err := overwriteConfig("rss_feeds.toml", config); err != nil {
		return errors.New("failed to overwrite RSS feed, " + err.Error())
	}

	fmt.Println("RSS feed has now been added.")

	return nil
}

func EditRSSFeed(config *Config) error {
	title := input.InputTitle()
	if len(title) == 0 {
		return errors.New("no title entered")
	}

	for index, feed := range config.FeedInfos {
		if feed.Title == title {
			fmt.Println("RSS feed found. Please enter the data. Any fields not entered will not be updated.")
			// 入力されなかった項目は変更しない
			title = input.InputTitle()
			if len(title) == 0 {
				title = feed.Title
			}
			url := input.InputURL()
			if len(url) == 0 {
				url = feed.URL
			}
			emails := input.InputEmails()
			if len(emails) == 0 {
				emails = feed.Email
			}
			newFeed := FeedInfo{
				Title:       title,
				URL:         url,
				Email:       emails,
				LastUpdated: feed.LastUpdated,
			}
			config.FeedInfos[index] = newFeed
			if err := overwriteConfig("rss_feeds.toml", config); err != nil {
				return errors.New("failed to overwrite RSS feed, " + err.Error())
			}
			return nil
		}
	}

	return errors.New("no RSS feed found")
}

func DeleteRSSFeed(config *Config) error {
	title := input.InputTitle()
	if len(title) == 0 {
		return errors.New("no title entered")
	}

	for index, feed := range config.FeedInfos {
		if feed.Title == title {
			if input.IsEnteredY(input.YOrOther("Are you sure you want to delete?")) {
				config.FeedInfos = slices.Delete(config.FeedInfos, index, index+1)
				if err := overwriteConfig("rss_feeds.toml", config); err != nil {
					return errors.New("failed to overwrite RSS feed, " + err.Error())
				}
				fmt.Println("Deletion completed.")
			} else {
				fmt.Println("Delete canceled.")
			}
			return nil
		}
	}

	return errors.New("no RSS feed found")
}
