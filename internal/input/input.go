package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func InputWithLabel(label string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func InputTitle() string {
	return InputWithLabel("Title: ")
}

func InputURL() string {
	return InputWithLabel("URL: ")
}

func InputEmail() string {
	return InputWithLabel("Email: ")
}

func InputEmails() []string {
	s := InputWithLabel("Email (0 or more): ")
	return slices.DeleteFunc(
		strings.Split(s, " "),
		func(e string) bool {
			return len(e) == 0
		},
	)
}

func InputPassword() string {
	return InputWithLabel("Password: ")
}

func InputHost() string {
	return InputWithLabel("Host: ")
}

func YOrOther(message string) string {
	return InputWithLabel(message + " (y or other): ")
}

func IsEnteredY(yOrOther string) bool {
	return len(yOrOther) == 1 && yOrOther[0] == 'y'
}
