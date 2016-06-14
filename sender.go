package sbms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/spf13/viper"
)

type MailSender struct {
	SenderRecords []Recipent
	Template      Email
}

func (self MailSender) BulkSend(dryRun bool) error {
	if self.SenderRecords == nil {
		return fmt.Errorf("SenderRecords is nil")
	}
	if &self.Template == nil {
		return fmt.Errorf("Template is nil")
	}

	var email Email
	email.From = self.Template.From
	email.Subject = self.Template.Subject
	email.Body = self.Template.Body
	for _, e := range self.SenderRecords {
		email.To = e.Address
		err := self.Send(email, dryRun)
		if err != nil {
			fmt.Printf("this address send failed. [%v]\n", e.Address)
		}
		if dryRun {
			break
		}
	}

	return nil
}

func (self MailSender) Send(email Email, dryRun bool) error {
	cred := credentials.NewSharedCredentials("", viper.GetString("aws.profile"))
	cfg := aws.NewConfig().WithRegion(viper.GetString("aws.region")).WithCredentials(cred)
	svc := ses.New(session.New(cfg))

	if dryRun {
		email.To = viper.GetString("mail.dryRunTo")
	}

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(email.Body),
					Charset: aws.String("ISO-2022-JP"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(email.Subject),
				Charset: aws.String("ISO-2022-JP"),
			},
		},
		Source: aws.String(email.From),
		ReplyToAddresses: []*string{
			aws.String(email.From),
		},
		ReturnPath:    aws.String(email.From),
		ReturnPathArn: aws.String(viper.GetString("aws.arn")),
		SourceArn:     aws.String(viper.GetString("aws.arn")),
	}
	fmt.Println(params)
	fmt.Println("dryRunTo = " + viper.GetString("mail.dryRunTo"))
	_, err := svc.SendEmail(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("success [%v]\n", email.To)

	// Pretty-print the response data.
	//fmt.Println(resp)
	return nil
}
