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

// Users returns a map of noncompliant users and related information.
func Users(ctx context.Context, accessKey, secretKey, types string) (map[string]map[string]string, error) {
	// User types to consider as valid
	requiredTypes := make(map[string]bool)
	for _, t := range strings.Split(types, ",") {
		requiredTypes[strings.ToLower(t)] = true
	}

	// Store a map containing [unique] user IDs and their respective attributes
	users := make(map[string]map[string]string)

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
		active := d.Attributes.Active
		email := d.Attributes.Email
		id := d.Attributes.ID
		inScope := d.Attributes.InAuditScope
		invited := d.Attributes.Invited
		name := d.Attributes.Name
		onboardingStatus := d.Attributes.OnboardingStatus
		personnelStatus := d.Attributes.PersonnelStatus

		var employeeType string
		if d.Attributes.EmployeeType != "" {
			employeeType = d.Attributes.EmployeeType
		}

		_, validUser := requiredTypes[employeeType]
		noncompliant := personnelStatus != "all_tasks_completed"

		// Store the user's ID as a key
		// Store their name, email, employee type, onboarding status, and personnel status as values
		if all(active, invited, inScope, validUser, noncompliant) {
			// Initialize the value map
			users[id] = make(map[string]string)
			users[id]["name"] = name
			users[id]["email"] = email
			users[id]["employee_type"] = employeeType
			users[id]["onboarding_status"] = onboardingStatus
			users[id]["personnel_status"] = personnelStatus
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
