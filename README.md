# secureframe-policy-minder

![secureframe-policy-minder](images/logo.jpg?raw=true "secureframe-policy-minder logo")

Send Slack reminders to Personnel to remind them to:

* Accept Policies
* Upload proof of Security Training

This tool is designed to be used as a scheduled task, such as GitHub Actions.

NOTE: This is using an undocumented Secureframe GraphQL API, so it may suddenly break. PR's welcome.

## Requirements

* A Slack Bot token 
* A Secureframe API token (findable via browser headers)
* A Secureframe Company ID (findable via browser headers)

As Secureframe does not yet have a public API, you'll need to grab the latter two bits of information using your browser's Developer Tools functionality.

## Finding your Secureframe authentication data

1. Visit <https://app.secureframe.com/>
2. Enter your browser's "Developer Tools" feature
3. Click on the **Console** tab.
4. Type `sessionStorage.getItem("AUTH_TOKEN");` and press <enter>. This will show your auth token.
5. Type `sessionStorage.getItem("CURRENT_COMPANY_USER");` and press <enter>. This will show your company use ID.
6. 

## Installation

```shell
go install github.com/chainguard-dev/secureframe-policy-minder@latest
```

## Slack App

Create a Slack app "From Scratch" at https://api.slack.com/apps

- Scopes required: `chat:write`

Save the token starting with `xoxb-`, as you will need it to send messages.

## Usage

You can run this app via the command-line or as a scheduled Github Action.

```
usage:
  -company-id string
    	secureframe company ID (default "adcfb3c-0b58-4c2c-af04-43b1a5031d61")
  -company-name string
    	The name of your compnay (default "Chainguard")
  -company-user-id string
    	secureframe company user ID (default "079b854c-c53a-4c71-bfb8-f9e87b13b6c4")
  -dry-run
    	dry-run mode
  -employee-types string
    	types of employees to contact (default "employee,contractor")
  -help-channel string
    	Slack channel for help (default "#security-and-compliance")
  -robot-name string
    	name of the robot (default "ComplyBot3000")
  -secureframe-token string
    	Secureframe bearer token
  -security-training-url string
    	URL to security training (default "https://securityawareness.usalearning.gov/cybersecurity/index.htm")
  -test-message-target string
    	override destination and send a single test message to this person
```

For initial testing, I recommend the `--dry-run` and `--test-message-target` flags.