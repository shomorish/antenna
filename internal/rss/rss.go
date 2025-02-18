package rss

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

func GetLink(links []Link, cond func(Link) bool) (Link, error) {
	for _, link := range links {
		if cond(link) {
			return link, nil
		}
	}
	return Link{}, errors.New("not found")
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type Feed interface {
	GetTitle() string
	GetDescription() string
	GetLinks() []Link
	GetUpdated() time.Time
	GetId() string
}

type AtomFeed struct {
	XMLName   xml.Name  `xml:"feed"`
	Namespace string    `xml:"xmlns,attr"`
	Title     string    `xml:"title"`
	SubTitle  string    `xml:"subtitle"`
	Links     []Link    `xml:"link"`
	Updated   time.Time `xml:"updated"`
	Id        string    `xml:"id"`
}

func (f *AtomFeed) GetTitle() string {
	return f.Title
}

func (f *AtomFeed) GetLinks() []Link {
	return f.Links
}

func (f *AtomFeed) GetDescription() string {
	return f.SubTitle
}

func (f *AtomFeed) GetUpdated() time.Time {
	return f.Updated
}

func (f *AtomFeed) GetId() string {
	return f.Id
}

var _ Feed = new(AtomFeed)

type FeedInfo struct {
	Title       string
	URL         string
	Email       []string
	LastUpdated time.Time `toml:"last_updated"`
}

func (r *FeedInfo) String() string {
	return fmt.Sprintf("Title: %s\nURL: %s\nEmail: %v\nLast Updated: %s", r.Title, r.URL, r.Email, r.LastUpdated)
}

type Config struct {
	EmailSender string
	Password    string
	Host        string
	FeedInfos   []FeedInfo `toml:"feed_info"`
}

func (r *Config) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Email Sender: %s\nPassword: %s", r.EmailSender, r.Password)
	for i, f := range r.FeedInfos {
		if i == 0 {
			fmt.Fprintf(&b, "\n\n[1]\n%s", f.String())
		} else {
			fmt.Fprintf(&b, "\n\n[%d]\n%s", i+1, f.String())
		}
	}
	return b.String()
}
