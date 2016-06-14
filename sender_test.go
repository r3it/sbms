package sbms

import (
	"testing"

	"github.com/r3it/sbms"
)

func TestGetFromDB(t *testing.T) {
	dryRun := true

	records, err := sbms.GetFromDB()
	if err != nil {
		t.Error(err)
	}

	template, err := sbms.LoadTextBody()
	if err != nil {
		t.Error(err)
	}

	var sender sbms.MailSender
	sender.SenderRecords = records
	sender.Template = template

	// sendErr := sender.Send(template)
	sendErr := sender.BulkSend(dryRun)
	if sendErr != nil {
		t.Error(sendErr)
	}

}
