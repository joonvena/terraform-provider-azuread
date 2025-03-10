// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package conditionalaccess_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf/pluginsdk"
)

type ConditionalAccessPolicyResource struct{}

func TestAccConditionalAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("enabledForReportingButNotEnforced"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_includedUserActions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.includedUserActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.includedUserActions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_sessionControls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sessionControls(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.sessionControls(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.sessionControls(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_sessionControlsDisabled(t *testing.T) {
	// This is testing the DiffSuppressFunc for the `session_controls` block

	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sessionControlsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.sessionControlsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_clientApplications(t *testing.T) {
	// This is a separate test for two reasons:
	// - conditional access policies applies either to users/groups or to client applications (workload identities)
	// - conditional access policies using client applications requires special licensing (Microsoft Entra Workload Identities)

	// Due to eventual consistency issues making it difficult to create a service principal on demand for inclusion in this
	// test policy, the config for this test requires a pre-existing service principal named "Terraform Acceptance Tests (Single Tenant)"
	// which should be linked to a single tenant application in the same tenant.

	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clientApplicationsIncluded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clientApplicationsExcluded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.clientApplicationsIncluded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_authenticationStrength(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationStrengthPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("grant_controls.0.authentication_strength_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_guestsOrExternalUsers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.guestsOrExternalUsersAllServiceProvidersIncluded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("conditions.0.users.0.included_guests_or_external_users.0.external_tenants.0.membership_kind").HasValue("all"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.guestsOrExternalUsersAllServiceProvidersExcluded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("conditions.0.users.0.excluded_guests_or_external_users.0.external_tenants.0.membership_kind").HasValue("all"),
			),
		},
		data.ImportStep(),
	})
}

func (r ConditionalAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	var id *string

	app, status, err := clients.ConditionalAccess.PoliciesClient.Get(ctx, state.ID, odata.Query{})
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Conditional Access Policy with ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve Conditional Access Policy with ID %q: %+v", state.ID, err)
	}
	id = app.ID

	return pointer.To(id != nil && *id == state.ID), nil
}

func (ConditionalAccessPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["None"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "enabledForReportingButNotEnforced"

  conditions {
    client_app_types    = ["all"]
    sign_in_risk_levels = ["medium"]
    user_risk_levels    = ["medium"]

    applications {
      included_applications = ["All"]
      excluded_applications = []
    }

    devices {
      filter {
        mode = "exclude"
        rule = "device.operatingSystem eq \"Doors\""
      }
    }

    locations {
      included_locations = ["All"]
      excluded_locations = ["AllTrusted"]
    }

    platforms {
      included_platforms = ["all"]
      excluded_platforms = ["android", "iOS"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }

  session_controls {
    application_enforced_restrictions_enabled = true
    cloud_app_security_policy                 = "blockDownloads"
    disable_resilience_defaults               = false
    persistent_browser_mode                   = "always"
    sign_in_frequency                         = 2
    sign_in_frequency_authentication_type     = "primaryAndSecondaryAuthentication"
    sign_in_frequency_interval                = "timeBased"
    sign_in_frequency_period                  = "days"
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) includedUserActions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["all"]

    applications {
      included_user_actions = [
        "urn:user:registerdevice",
        "urn:user:registersecurityinfo",
      ]
    }

    locations {
      included_locations = ["All"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) sessionControls(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["All"]
    }

    locations {
      included_locations = ["All"]
    }

    platforms {
      included_platforms = ["all"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  session_controls {
    application_enforced_restrictions_enabled = true
    disable_resilience_defaults               = true
    cloud_app_security_policy                 = "monitorOnly"
    persistent_browser_mode                   = "never"
    sign_in_frequency                         = 10
    sign_in_frequency_period                  = "hours"
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) sessionControlsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["All"]
    }

    locations {
      included_locations = ["All"]
    }

    platforms {
      included_platforms = ["all"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }

  session_controls {
    application_enforced_restrictions_enabled = false
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) clientApplicationsIncluded(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_service_principal" "test" {
  display_name = "Terraform Acceptance Tests (Single Tenant)"
}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["all"]

    applications {
      included_applications = ["All"]
    }

    client_applications {
      included_service_principals = [data.azuread_service_principal.test.object_id]
    }

    service_principal_risk_levels = ["medium"]

    users {
      included_users = ["None"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) clientApplicationsExcluded(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_service_principal" "test" {
  display_name = "Terraform Acceptance Tests (Single Tenant)"
}

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["all"]

    applications {
      included_applications = ["All"]
    }

    client_applications {
      included_service_principals = ["ServicePrincipalsInMyTenant"]
      excluded_service_principals = [data.azuread_service_principal.test.object_id]
    }

    service_principal_risk_levels = ["medium"]

    users {
      included_users = ["None"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) authenticationStrengthPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[2]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["None"]
    }

    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
  }

  grant_controls {
    operator                          = "OR"
    authentication_strength_policy_id = azuread_authentication_strength_policy.test.id
  }
}
`, AuthenticationStrengthPolicyResource{}.basic(data), data.RandomInteger)
}

func (ConditionalAccessPolicyResource) guestsOrExternalUsersAllServiceProvidersIncluded(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["None"]
    }

    users {
      included_guests_or_external_users {
        guest_or_external_user_types = ["internalGuest", "serviceProvider"]
        external_tenants {
          membership_kind = "all"
        }
      }
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) guestsOrExternalUsersAllServiceProvidersExcluded(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    client_app_types = ["browser"]

    applications {
      included_applications = ["None"]
    }

    users {
      included_users = ["None"]
      excluded_guests_or_external_users {
        guest_or_external_user_types = ["internalGuest", "serviceProvider"]
        external_tenants {
          membership_kind = "all"
        }
      }
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}
