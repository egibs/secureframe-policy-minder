package secureframe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	restEndpoint = "https://api.secureframe.com"
)

type User struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			CreatedAt                      time.Time `json:"created_at"`
			UpdatedAt                      time.Time `json:"updated_at"`
			ID                             string    `json:"id"`
			Active                         bool      `json:"active"`
			ActiveSource                   any       `json:"active_source"`
			DepartmentID                   any       `json:"department_id"`
			EmployeeType                   string    `json:"employee_type"`
			EndDate                        any       `json:"end_date"`
			InAuditScope                   bool      `json:"in_audit_scope"`
			Invited                        bool      `json:"invited"`
			InvitedAt                      time.Time `json:"invited_at"`
			OnboardingStatus               string    `json:"onboarding_status"`
			PersonnelStatus                string    `json:"personnel_status"`
			Role                           string    `json:"role"`
			SecureframeAgentAcknowledgedAt any       `json:"secureframe_agent_acknowledged_at"`
			StartDate                      string    `json:"start_date"`
			Title                          any       `json:"title"`
			AccessRole                     any       `json:"access_role"`
			Email                          string    `json:"email"`
			FirstName                      string    `json:"first_name"`
			ImageURL                       any       `json:"image_url"`
			LastName                       string    `json:"last_name"`
			ManagerName                    any       `json:"manager_name"`
			MiddleName                     string    `json:"middle_name"`
			Name                           string    `json:"name"`
			PreferredFirstName             any       `json:"preferred_first_name"`
		} `json:"attributes"`
		Relationships struct {
		} `json:"relationships"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

// makeRequest performs an authenticated HTTP request to the specified endpoint and returns a JSON-friendly byte slice.
func request(ctx context.Context, url, method, accessKey, secretKey string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", accessKey, secretKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return rb, nil
}

// nextFive returns the next five reminder windows (start and end dates).
func nextFive(currentTime time.Time, invitedAt time.Time) map[int]map[string]time.Time {
	dates := make(map[int]map[string]time.Time)
	start := invitedAt.AddDate(0, 0, 354)
	end := start.AddDate(0, 0, 10)

	for i := 0; i < 5; i++ {
		start = start.AddDate(0, 0, 354)
		end = start.AddDate(0, 0, 10)

		if !currentTime.After(end) {
			dates[i] = map[string]time.Time{
				"start": start,
				"end":   end,
			}
		}
	}

	// Test the logic by uncommenting the following lines
	// dates[5] = map[string]time.Time{
	// 	"start": currentTime.AddDate(0, 0, -1),
	// 	"end":   currentTime.AddDate(0, 0, 10),
	// }

	return dates
}

// validDate checks whether the given time falls within two dates.
// the two dates are stored in a map of integers with "start" and "end" keys with time values
func validDate(date time.Time, dates map[int]map[string]time.Time) (map[string]time.Time, bool) {
	for i := range len(dates) {
		start := dates[i]["start"]
		end := dates[i]["end"]

		if date.After(start) && date.Before(end) {
			return dates[i], true
		}
	}
	return nil, false
}

// Users returns a map of noncompliant users and related information.
func Users(ctx context.Context, accessKey, secretKey, types string) (map[string]map[string]any, error) {
	// User types to consider as valid
	requiredTypes := make(map[string]bool)
	for _, t := range strings.Split(types, ",") {
		requiredTypes[strings.ToLower(t)] = true
	}

	// Store a map containing [unique] user IDs and their respective attributes
	users := make(map[string]map[string]any)

	requestUrl := fmt.Sprintf("%s/users", restEndpoint)

	rb, err := request(ctx, requestUrl, "GET", accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	out := User{}
	if err := json.Unmarshal(rb, &out); err != nil {
		return nil, err
	}

	data := out.Data

	// Filter out compliant users and only store users with overdue or incomplete training
	for _, d := range data {
		_, validUser := requiredTypes[d.Attributes.EmployeeType]
		noncompliant := d.Attributes.PersonnelStatus != "all_tasks_completed"

		// Do some future-looking calculations to determine whether a user should be notified
		// The only date exposed by the REST API is the user's invited date
		// Take the invited date and calculate the next five years of notifications
		// The notifications will be sent within the range t + (354) < n < t + (364) to provide a 10-day window
		invitedAt := d.Attributes.InvitedAt

		// Capture the current time
		currentTime := time.Now()

		// Retrieve the next five reminder windows
		nextFive := nextFive(currentTime, invitedAt)

		// If the user is:
		// - active, invited, in the audit scope, a valid user, noncompliant, and within the remindStart and remindEnd window, remind the user then:
		// - store the user's ID as a key
		// - store their name, email, employee type, onboarding status, and personnel status as values
		remindDates, inWindow := validDate(currentTime, nextFive)

		if all(d.Attributes.Active, d.Attributes.Invited, d.Attributes.InAuditScope, validUser, noncompliant, inWindow) {
			// Initialize the value map
			// Use any to store strings and the dates
			users[d.Attributes.ID] = map[string]any{
				"name":              d.Attributes.Name,
				"email":             d.Attributes.Email,
				"employee_type":     d.Attributes.EmployeeType,
				"onboarding_status": d.Attributes.OnboardingStatus,
				"personnel_status":  d.Attributes.PersonnelStatus,
			}

			// Prevent any nil panics
			if remindDates != nil {
				users[d.Attributes.ID]["remind_dates"] = remindDates
			}
		}
	}

	return users, nil
}

// all returns a single boolean based on a slice of booleans.
func all(conditions ...bool) bool {
	for _, condition := range conditions {
		if !condition {
			return false
		}
	}
	return true
}
