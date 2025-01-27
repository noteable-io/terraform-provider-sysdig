package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccUser(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" && secure == "" {
				t.Fatal("either SYSDIG_MONITOR_API_TOKEN or SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: userWithSystemRole(rText()),
			},
			{
				Config: userWithName(rText()),
			},
			{
				Config: userWithoutSystemRole(rText()),
			},
			{
				Config: userMinimumConfiguration(),
			},
			{
				ResourceName:      "sysdig_user.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  system_role = "ROLE_USER"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userWithSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  system_role = "ROLE_CUSTOMER"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userWithoutSystemRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
  first_name = "%s"
  last_name  = "%s"
}`, name, name)
}

func userMinimumConfiguration() string {
	return `
resource "sysdig_user" "sample" {
  email      = "terraform-test+user@sysdig.com"
}`
}
