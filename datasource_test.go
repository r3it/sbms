package sbms

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/r3it/sbms"
)

func TestGetFromDB(t *testing.T) {
	records, err := sbms.GetFromDB()
	if err != nil {
		t.Error(err)
	}

	for _, r := range records {
		t.Log(r.Address + " / " + r.Name)
	}
}

func TestLoadTextBody(t *testing.T) {
	email, err := sbms.LoadTextBody()
	if err != nil {
		t.Error(err)
	}

	text, ioErr := ioutil.ReadFile("mailbody.txt")
	if ioErr != nil {
		t.Error(ioErr)
	}
	lines := strings.Split(string(text), "\n")

	if lines[0] != email.Subject {
		t.Errorf("got %v\nwant %v", email.Subject, lines[0])
	}
	t.Log(email.Subject)

	t.Log("-----")

	body := strings.Join(lines[1:], "\n")
	if body != email.Body {
		t.Errorf("got %v\nwant %v", email.Body, body)
	}
	t.Log(email.Body)
}
