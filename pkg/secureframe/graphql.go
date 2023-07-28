package secureframe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	defaultEndpoint    = "https://app.secureframe.com/graphql"
)

// SearchCompanyUsersResult is the struct returned by "personnelTabContentsSearch"
// This was autogenerated by https://mholt.github.io/json-to-go/
type SearchCompanyUsersResult struct {
	Data struct {
		SearchCompanyUsers struct {
			Data struct {
				Collection []Person `json:"collection"`
				Metadata   struct {
					CurrentPage int    `json:"currentPage"`
					LimitValue  int    `json:"limitValue"`
					TotalCount  int    `json:"totalCount"`
					TotalPages  int    `json:"totalPages"`
					Typename    string `json:"__typename"`
				} `json:"metadata"`
				Typename string `json:"__typename"`
			} `json:"data"`
			Typename string `json:"__typename"`
		} `json:"searchCompanyUsers"`
	} `json:"data"`
}

// Company is returned as part of CompanyUsersResult
type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Typename string `json:"__typename"`
}

// CompanyUsersResult is the struct returned by "getCompanyUsersForCurrentUser"
// This was autogenerated by https://mholt.github.io/json-to-go/
type CompanyUsersResult struct {
	Data struct {
		GetCompanyUsersForCurrentUser []struct {
			ID       string  `json:"id"`
			Company  Company `json:"company"`
			Typename string  `json:"__typename"`
		} `json:"getCompanyUsersForCurrentUser"`
	} `json:"data"`
}

type Person struct {
	Active                              bool      `json:"active"`
	BackgroundCheckStatus               string    `json:"backgroundCheckStatus"`
	BackgroundCheckExists               bool      `json:"backgroundCheckExists"`
	CanBeInvited                        bool      `json:"canBeInvited"`
	CcpaTrainingCompleted               bool      `json:"ccpaTrainingCompleted"`
	Email                               string    `json:"email"`
	EmployeeType                        string    `json:"employeeType"`
	EndDate                             any       `json:"endDate"`
	GdprTrainingCompleted               bool      `json:"gdprTrainingCompleted"`
	GoogleWorkspaceMfaEnabled           bool      `json:"googleWorkspaceMfaEnabled"`
	HipaaTrainingCompletedAt            any       `json:"hipaaTrainingCompletedAt"`
	HipaaTrainingStatus                 string    `json:"hipaaTrainingStatus"`
	ID                                  string    `json:"id"`
	InAuditScope                        bool      `json:"inAuditScope"`
	Invited                             bool      `json:"invited"`
	InvitedAt                           time.Time `json:"invitedAt"`
	ManagedComputer                     bool      `json:"managedComputer"`
	Name                                string    `json:"name"`
	PciSecureCodeTrainingCompleted      bool      `json:"pciSecureCodeTrainingCompleted"`
	PciSecureCodeTrainingCompletedAt    any       `json:"pciSecureCodeTrainingCompletedAt"`
	PciTrainingCompleted                bool      `json:"pciTrainingCompleted"`
	PersonnelTasksNotificationFrequency string    `json:"personnelTasksNotificationFrequency"`
	PoliciesAccepted                    bool      `json:"policiesAccepted"`
	PoliciesAcceptedAt                  time.Time `json:"policiesAcceptedAt"`
	SecurityTrainingCompleted           bool      `json:"securityTrainingCompleted"`
	SecurityTrainingCompletedAt         time.Time `json:"securityTrainingCompletedAt"`
	SecurityTrainingNewHireCompleted    bool      `json:"securityTrainingNewHireCompleted"`
	SecurityTrainingNewHireCompletedAt  time.Time `json:"securityTrainingNewHireCompletedAt"`
	StartDate                           string    `json:"startDate"`
	DepartmentID                        any       `json:"departmentId"`
	Role                                string    `json:"role"`
	Typename                            string    `json:"__typename"`
}

type payload struct {
	OperationName string    `json:"operationName"`
	Variables     variables `json:"variables"`
	Query         string    `json:"query"`
}

type variables struct {
	SearchKick           *searchKick `json:"searchkick,omitempty"`
	CurrentCompanyUserID string      `json:"current_company_user_id,omitempty"`
	CompanyID            string      `json:"company_id,omitempty"`
	Page                 int         `json:"page,omitempty"`
	Limit                int         `json:"limit,omitempty"`

	// Used by Test
	ID   *string `json:"id,omitempty"`
	Pass bool    `json:"pass,omitempty"`

	Key string `json:"key,omitempty"`
}

type searchKick struct {
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
	Query   string `json:"query"`
}

func query(ctx context.Context, token string, in interface{}, out interface{}) error {
	payloadBytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	//	log.Printf("payload: %s", payloadBytes)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", defaultEndpoint, body)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	log.Printf("POST'ing to %s with:\n%s", defaultEndpoint, payloadBytes)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// log.Printf("response: %s", rb)

	if err := json.Unmarshal(rb, out); err != nil {
		return fmt.Errorf("unmarshal output: %w\ncontents: %s", err, rb)
	}

	// log.Printf("parsed response: %+v", out)
	return nil
}

func GetCompany(ctx context.Context, userID string, token string) (Company, error) {
	in := payload{
		OperationName: "getCompanyUsersForCurrentUser",
		Variables: variables{
			CurrentCompanyUserID: userID,
		},
		Query: `query getCompanyUsersForCurrentUser {
			getCompanyUsersForCurrentUser {
				id
				company {
					id
					name
					__typename
				}
				__typename
			}
		}`,
	}

	out := &CompanyUsersResult{}
	if err := query(ctx, token, in, out); err != nil {
		return Company{}, fmt.Errorf("query: %w", err)
	}

	return out.Data.GetCompanyUsersForCurrentUser[0].Company, nil
}

func Personnel(ctx context.Context, companyID string, companyUserID string, token string) ([]Person, error) {
	in := payload{
		OperationName: "personnelTabContentsSearch",
		Variables: variables{
			SearchKick: &searchKick{
				Page:    1,
				PerPage: 1000,
				Query:   "*",
			},
			CurrentCompanyUserID: companyUserID,
			CompanyID:            companyID,
		},
		Query: `fragment PersonnelTabContentsCompanyUsers on CompanyUser {
					active
					backgroundCheckStatus
					backgroundCheckExists
					canBeInvited
					ccpaTrainingCompleted
					email
					employeeType
					endDate
					gdprTrainingCompleted
					googleWorkspaceMfaEnabled
					hipaaTrainingCompletedAt
					hipaaTrainingStatus
					id
					inAuditScope
					invited
					invitedAt
					managedComputer
					name
					pciSecureCodeTrainingCompleted
					pciSecureCodeTrainingCompletedAt
					pciTrainingCompleted
					personnelTasksNotificationFrequency
					policiesAccepted
					policiesAcceptedAt
					securityTrainingCompleted
					securityTrainingCompletedAt
					securityTrainingNewHireCompleted
					securityTrainingNewHireCompletedAt
					startDate
					departmentId
					role
					__typename
				}

				query personnelTabContentsSearch($searchkick: CompanyUserSearchkickInput, $companyId: ID) {
					searchCompanyUsers(searchkick: $searchkick, companyId: $companyId) {
						data {
						collection {
							...PersonnelTabContentsCompanyUsers
							__typename
						}
						metadata {
							currentPage
							limitValue
							totalCount
							totalPages
							__typename
						}
						__typename
						}
						__typename
					}
				}
		`,
	}

	out := &SearchCompanyUsersResult{}
	if err := query(ctx, token, in, out); err != nil {
		return out.Data.SearchCompanyUsers.Data.Collection, fmt.Errorf("request: %w", err)
	}

	return out.Data.SearchCompanyUsers.Data.Collection, nil
}
