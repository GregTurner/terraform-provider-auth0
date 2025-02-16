package organization

import (
	"context"
	"net/http"

	"github.com/auth0/go-auth0/management"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/auth0/terraform-provider-auth0/internal/config"
	internalSchema "github.com/auth0/terraform-provider-auth0/internal/schema"
)

// NewConnectionResource will return a new auth0_organization_connection resource.
func NewConnectionResource() *schema.Resource {
	return &schema.Resource{
		Description:   "With this resource, you can manage enabled connections on an organization.",
		CreateContext: createOrganizationConnection,
		ReadContext:   readOrganizationConnection,
		UpdateContext: updateOrganizationConnection,
		DeleteContext: deleteOrganizationConnection,
		Importer: &schema.ResourceImporter{
			StateContext: internalSchema.ImportResourceGroupID(internalSchema.SeparatorColon, "organization_id", "connection_id"),
		},
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the organization to enable the connection for.",
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the connection to enable for the organization.",
			},
			"assign_membership_on_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "When true, all users that log in with this connection will be automatically granted " +
					"membership in the organization. When false, users must be granted membership in the organization" +
					" before logging in with this connection.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the enabled connection.",
			},
			"strategy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The strategy of the enabled connection.",
			},
		},
	}
}

func createOrganizationConnection(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*config.Config).GetAPI()

	organizationID := data.Get("organization_id").(string)

	connectionID := data.Get("connection_id").(string)
	assignMembershipOnLogin := data.Get("assign_membership_on_login").(bool)

	organizationConnection := &management.OrganizationConnection{
		ConnectionID:            &connectionID,
		AssignMembershipOnLogin: &assignMembershipOnLogin,
	}

	if err := api.Organization.AddConnection(organizationID, organizationConnection); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(organizationID + ":" + connectionID)

	return readOrganizationConnection(ctx, data, meta)
}

func readOrganizationConnection(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*config.Config).GetAPI()

	organizationID := data.Get("organization_id").(string)
	connectionID := data.Get("connection_id").(string)

	organizationConnection, err := api.Organization.Connection(organizationID, connectionID)
	if err != nil {
		if err, ok := err.(management.Error); ok && err.Status() == http.StatusNotFound {
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	result := multierror.Append(
		data.Set("assign_membership_on_login", organizationConnection.GetAssignMembershipOnLogin()),
		data.Set("name", organizationConnection.GetConnection().GetName()),
		data.Set("strategy", organizationConnection.GetConnection().GetStrategy()),
	)

	return diag.FromErr(result.ErrorOrNil())
}

func updateOrganizationConnection(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*config.Config).GetAPI()

	organizationID := data.Get("organization_id").(string)

	connectionID := data.Get("connection_id").(string)
	assignMembershipOnLogin := data.Get("assign_membership_on_login").(bool)

	organizationConnection := &management.OrganizationConnection{
		AssignMembershipOnLogin: &assignMembershipOnLogin,
	}

	if err := api.Organization.UpdateConnection(organizationID, connectionID, organizationConnection); err != nil {
		return diag.FromErr(err)
	}

	return readOrganizationConnection(ctx, data, meta)
}

func deleteOrganizationConnection(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*config.Config).GetAPI()

	organizationID := data.Get("organization_id").(string)

	connectionID := data.Get("connection_id").(string)

	if err := api.Organization.DeleteConnection(organizationID, connectionID); err != nil {
		if err, ok := err.(management.Error); ok && err.Status() == http.StatusNotFound {
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
