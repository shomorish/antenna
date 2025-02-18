package rss

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shomorish/antenna/internal/email"
)

var (
	ErrNoUpdates      = errors.New("no updates")
	ErrNotImplemented = errors.New("not implemented")
)

func getRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Request to `%s` failed.\n", url)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err)
		return nil
	}

	return body
}

func UnmarshalRSSFeed(body []byte) (Feed, error) {
	return nil, ErrNotImplemented
}

func UnmarshalAtomFeed(body []byte) (Feed, error) {
	feed := new(AtomFeed)
	if err := xml.Unmarshal(body, feed); err != nil {
		return nil, err
	}
	return feed, nil
}

func sendFeedEmails(from string, password string, smtpHost string, to []string, feed Feed) error {
	link, _ := GetLink(feed.GetLinks(), func(l Link) bool { return l.Rel == "alternate" })
	for _, t := range to {
		message := []byte(
			"From: Antenna <" + from + ">\r\n" +
				"To: <" + t + ">\r\n" +
				"Subject: " + feed.GetTitle() + "\r\n" +
				"\r\n" +
				link.Href)
		if err := email.SendEmail(from, password, smtpHost, t, message); err != nil {
			return err
		}
	}
	return nil
}

func UpdateRSSFeeds(config *Config) {
	// 各フィードをUnmarshalして、更新されているか確認
	for index, feedInfo := range config.FeedInfos {
		body := getRequest(feedInfo.URL)
		if body == nil {
			continue
		}

		// フォーマットの種類を判定
		// TODO: フィードの更新日時がパースできていないかもしれない
		var feed Feed
		var err error
		x := string(body)
		if strings.Contains(x, "</rss>") {
			// RSSフィード
			feed, err = UnmarshalRSSFeed(body)
			if err != nil {
				continue
			}
		} else if strings.Contains(x, "</feed>") {
			// Atomフィード
			feed, err = UnmarshalAtomFeed(body)
			if err != nil {
				continue
			}
		} else {
			fmt.Printf("It was an unsupported format: %s\n", x)
			continue
		}

		// フィードが更新されているか確認
		if feed.GetUpdated().Compare(feedInfo.LastUpdated) == 1 {
			config.FeedInfos[index].LastUpdated = feed.GetUpdated()
			sendFeedEmails(config.EmailSender, config.Password, config.Host, feedInfo.Email, feed)
		}
	}

	overwriteConfig("rss_feeds.toml", config)

	fmt.Println("Feed update completed.")
}
