package main

import (
	"fmt"
	"io"

	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/r3it/sbms"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
	HELP = `
description:

	sbms is simple bulk mail sender.

required:

  * IAM profile
	* Amazon SES setup
	* MySQL user table (required columns are 'email' and 'name')
	* Mail template (default mailbody.txt)

commands are:

	send  send mail to users

config file example (see config.yaml.sample):

----
mysql:
  dbname: id:password@tcp(your-DB-uri.com:3306)/dbname
  query: select email, name from user where deleted = 0

mail:
  file: mailbody.txt
  from: foo@bar.com
	dryRunTo: fooTest@bar.com

aws:
  profile: your-credentials-profile
  region: us-west-2
  arn: arn:aws:ses:us-west-2:foobar........
----
`
)

type CLI struct {
	outStream, errStream io.Writer
}

func (c *CLI) Run() int {
	optarg.Header("General options")
	optarg.Add("h", "help", "Displays this help.", false)
	optarg.Add("v", "version", "Displays version information.", false)

	optarg.Header("send  send mail to users")
	optarg.Add("D", "dryRun", "dryRun option", false)

	var version, help, dryRun bool

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "h":
			help = opt.Bool()
		case "v":
			version = opt.Bool()
		case "D":
			dryRun = opt.Bool()
		}
	}

	if version {
		fmt.Fprintf(c.errStream, "sbms version %s\n", Version)
		return ExitCodeOK
	}
	if help {
		optarg.Usage()
		fmt.Fprintln(c.outStream, HELP)
		return ExitCodeOK
	}

	if len(optarg.Remainder) == 1 {
		cmd := optarg.Remainder[0]
		switch cmd {
		case "send":
			return cmdSend(c, dryRun)
		}
	} else {
		optarg.Usage()
	}

	return ExitCodeOK
}

func cmdSend(c *CLI, dryRun bool) int {
	if dryRun {
		fmt.Fprintln(c.outStream, "##### try dry run #####")
	}

	records, err := sbms.GetFromDB()
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		fmt.Fprintf(c.errStream, "GetFromDB failed\n")
	}
	if dryRun {
		fmt.Fprintln(c.outStream, "##### actually recipents are... #####")
		for _, r := range records {
			fmt.Fprintln(c.outStream, r.Address)
		}
		fmt.Fprintln(c.outStream, "##### that's all. #####")
	}

	template, err := sbms.LoadTextBody()
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		fmt.Fprintf(c.errStream, "LoadTextBody failed\n")
	}

	var sender sbms.MailSender
	sender.SenderRecords = records
	sender.Template = template

	sendErr := sender.BulkSend(dryRun)
	if sendErr != nil {
		fmt.Fprintln(c.errStream, sendErr)
		fmt.Fprintf(c.errStream, "BulkSend failed\n")
	}

	return ExitCodeOK
}
