# secureframe-policy-minder

![secureframe-policy-minder](images/logo.jpg?raw=true "secureframe-policy-minder logo")

Send Slack reminders to Personnel to remind them to:

* Accept Policies
* Upload proof of Security Training

This tool is designed to be used as a scheduled task, such as GitHub Actions.

## Requirements

* A Slack Bot token
* A Secureframe [API key](https://developer.secureframe.com/#section/Authentication)

## Installation

```shell
go install github.com/chainguard-dev/secureframe-policy-minder@latest
```

## Slack App

Create a Slack app "From Scratch" at https://api.slack.com/apps

- Scopes required: `chat:write`

Save the token starting with `xoxb-`, as you will need it to send messages.

## Usage

You can run this app via the command-line or as a scheduled Github Action (see [examples](examples))

```
usage:
  -access-key string
      secureframe access key
  -company string
      company name used for notifications
  -dry-run
    	dry-run mode
  -employee-types string
    	types of employees to contact (default "employee,contractor")
  -help-channel string
    	Slack channel for help (default "#security-and-compliance")
  -robot-name string
    	name of the robot (default "ComplyBot3000")
  -secret-key string
      secureframe secret key
  -secureframe-token string
    	Secureframe bearer token
  -security-training-url string
    	URL to security training (default "https://securityawareness.usalearning.gov/cybersecurity/index.htm")
  -test-message-target string
    	override destination and send a single test message to this person
```

For initial testing, I recommend the `--dry-run` and `--test-message-target` flags.
