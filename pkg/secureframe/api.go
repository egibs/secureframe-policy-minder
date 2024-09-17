package secureframe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

const (
	restEndpoint = "https://api.secureframe.com"
)

// makeRequest performs an authenticated HTTP request to the specified endpoint and returns a JSON-friendly byte slice.
func makeRequest(ctx context.Context, url, method, accessKey, secretKey string) ([]byte, error) {
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

// GetUsers returns a map of noncompliant users and related information.
func GetUsers(ctx context.Context, accessKey, secretKey, types string) (map[string]map[string]string, error) {
	// User types to consider as valid
	requiredTypes := make(map[string]bool)
	for _, t := range strings.Split(types, ",") {
		requiredTypes[strings.ToLower(t)] = true
	}

	// Store a map containing [unique] user IDs and their respective attributes
	users := make(map[string]map[string]string)

	requestUrl := fmt.Sprintf("%s/users", restEndpoint)

	rb, err := makeRequest(ctx, requestUrl, "GET", accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	// Parse response bytes using Gabs to avoid defining structs
	r, err := gabs.ParseJSON(rb)
	if err != nil {
		return nil, err
	}

	// All of the users are stored in a `data` array
	data := r.Path("data").Children()

	// Filter out compliant users and only store users with overdue or incomplete training
	for _, d := range data {
		active := d.Path("attributes.active").Data().(bool)
		email := d.Path("attributes.email").Data().(string)
		id := d.Path("attributes.id").Data().(string)
		inScope := d.Path("attributes.in_audit_scope").Data().(bool)
		invited := d.Path("attributes.invited").Data().(bool)
		name := d.Path("attributes.name").Data().(string)
		onboardingStatus := d.Path("attributes.onboarding_status").Data().(string)
		personnelStatus := d.Path("attributes.personnel_status").Data().(string)

		var employeeType string
		if d.Path("attributes.employee_type").Data() != nil {
			employeeType = d.Path("attributes.employee_type").Data().(string)
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
