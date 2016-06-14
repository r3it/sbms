sbms (simple bulk mail sender)
====

sbms is simple bulk mail sender.

## required:

* IAM profile and user policy.
* MySQL user table (required columns are 'email' and 'name')
* Mail template (default mailbody.txt)
* of cource Amazon SES setup

## Configuration

You must setup config.yaml, mailbody.txt and MySQL table.

mailbody.txt format is

* first line is Subject
* behind the second line is the body text

## Usage:

**dry run**
```
$ sbms -D send
```

**actual execute command**
```
$ sbms send
```

usage output here:

```
Usage: sbms [options]:

[General options]
    --help, -h: Displays this help.
 --version, -v: Displays version information.

[send  send mail to users]
  --dryRun, -D: dryRun option
```

## config file example (see config.yaml.sample):

```
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
```
