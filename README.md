# secureframe-policy-minder

EXPERIMENTAL - Send Slack reminders to Personnel to remind them to:

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
5. Type `sessionStorage.getItem("CURRENT_COMPANY_USER");` and press <enter>. This will show your company ID.

## Installation

```shell
go install github.com/chainguard-dev/secureframe-policy-minder@latest
```

## Slack App

Scopes required: `chat:write`.

