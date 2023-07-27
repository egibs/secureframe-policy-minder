package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"strings"

	"github.com/chainguard-dev/secureframe-policy-minder/pkg/secureframe"
	"github.com/danott/envflag"
)

var (
	dryRunFlag        = flag.Bool("dry-run", false, "dry-run mode")
	sfTokenFlag       = flag.String("secureframe-token", "", "Secureframe bearer token")
	companyIDFlag     = flag.String("company", "adcfb3c-0b58-4c2c-af04-43b1a5031d61", "secureframe company ID")
	companyUserIDFlag = flag.String("company-user", "079b854c-c53a-4c71-bfb8-f9e87b13b6c4", "secureframe company user ID")
	employeeTypesFlag = flag.String("employee-types", "employee,contractor", "types of employees to contact")
)

func main() {
	flag.Parse()
	// makes SECUREFRAME_TOKEN available to secureframe-token
	envflag.Parse()

	ppl, err := secureframe.Personnel(context.Background(), *companyIDFlag, *companyUserIDFlag, *sfTokenFlag)
	if err != nil {
		log.Panicf("Secureframe test query failed: %v", err)
	}
	log.Printf("PPL: -- %+v -- ", ppl)

	requiredTypes := map[string]bool{}
	for _, t := range strings.Split(*employeeTypesFlag, ",") {
		requiredTypes[strings.ToLower(t)] = true
	}

	for _, p := range ppl {
		if !p.Active {
			continue
		}
		if !p.Invited {
			continue
		}

		eType := strings.ToLower(p.EmployeeType)
		if !requiredTypes[eType] {
			continue
		}

		needs := []string{}
		if !p.PoliciesAccepted {
			needs = append(needs, "policy")
		}
		if !p.SecurityTrainingCompleted {
			needs = append(needs, "security-training")
		}
		if len(needs) > 0 {
			log.Printf("%s needs %s", p.Email, needs)
		}
	}
}
