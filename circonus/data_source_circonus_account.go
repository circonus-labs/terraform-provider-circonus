package circonus

import (
	"context"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	accountAddress1Attr      = "address1"
	accountAddress2Attr      = "address2"
	accountCCEmailAttr       = "cc_email"
	accountCityAttr          = "city"
	accountContactGroupsAttr = "contact_groups"
	accountCountryAttr       = "country"
	accountCurrentAttr       = "current"
	accountDescriptionAttr   = "description"
	accountEmailAttr         = "email"
	accountIDAttr            = "id"
	accountInvitesAttr       = "invites"
	accountLimitAttr         = "limit"
	accountNameAttr          = "name"
	accountOwnerAttr         = "owner"
	accountRoleAttr          = "role"
	accountStateProvAttr     = "state"
	accountTimezoneAttr      = "timezone"
	accountTypeAttr          = "type"
	accountUIBaseURLAttr     = "ui_base_url"
	accountUsageAttr         = "usage"
	accountUsedAttr          = "used"
	accountUserIDAttr        = "id"
	accountUsersAttr         = "users"
)

var accountDescription = map[schemaAttr]string{
	accountContactGroupsAttr: "Contact Groups in this account",
	accountInvitesAttr:       "Outstanding invites attached to the account",
	accountUsageAttr:         "Account's usage limits",
	accountUsersAttr:         "Users attached to this account",
}

func dataSourceCirconusAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCirconusAccountRead,

		Schema: map[string]*schema.Schema{
			// _cid
			accountIDAttr: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{accountCurrentAttr},
				ValidateFunc: validateFuncs(
					validateRegexp(accountIDAttr, config.AccountCIDRegex),
				),
				Description: accountDescription[accountIDAttr],
			},
			// determines whether to pull /account/current or specific cid
			accountCurrentAttr: {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{accountIDAttr},
				Description:   accountDescription[accountCurrentAttr],
			},
			// _countact_groups
			accountContactGroupsAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: accountDescription[accountContactGroupsAttr],
			},
			// _owner
			accountOwnerAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountOwnerAttr],
			},
			// _ui_base_url
			accountUIBaseURLAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountUIBaseURLAttr],
			},
			// _usage
			accountUsageAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountUsageAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// _limit
						accountLimitAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: accountDescription[accountLimitAttr],
						},
						// _type
						accountTypeAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountTypeAttr],
						},
						// _used
						accountUsedAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: accountDescription[accountUsedAttr],
						},
					},
				},
			},
			// address1
			accountAddress1Attr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountAddress1Attr],
			},
			// address2
			accountAddress2Attr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountAddress2Attr],
			},
			// cc_email
			accountCCEmailAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCCEmailAttr],
			},
			// city
			accountCityAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCityAttr],
			},
			// country_code
			accountCountryAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountCountryAttr],
			},
			// description
			accountDescriptionAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountDescriptionAttr],
			},
			// invites
			accountInvitesAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountInvitesAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						accountEmailAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountEmailAttr],
						},
						accountRoleAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountRoleAttr],
						},
					},
				},
			},
			// name
			accountNameAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountNameAttr],
			},
			// state_prov
			accountStateProvAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountStateProvAttr],
			},
			// timezone
			accountTimezoneAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: accountDescription[accountTimezoneAttr],
			},
			// users
			accountUsersAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: accountDescription[accountUsersAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						accountUserIDAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountUserIDAttr],
						},
						accountRoleAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: accountDescription[accountRoleAttr],
						},
					},
				},
			},
		},
	}
}

// dataSourceCirconusAccountRead - map account object from API to schema.ResourceData.
func dataSourceCirconusAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerContext).client
	var diags diag.Diagnostics

	var cid string
	if v, ok := d.GetOk(accountIDAttr); ok {
		cid = v.(string)
	}

	if v, ok := d.GetOk(accountCurrentAttr); ok {
		if v.(bool) {
			cid = ""
		}
	}

	acct, err := client.FetchAccount(api.CIDType(&cid))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(acct.CID)
	if err := d.Set(accountIDAttr, acct.CID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountContactGroupsAttr, acct.ContactGroups); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountOwnerAttr, acct.OwnerCID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountUIBaseURLAttr, acct.UIBaseURL); err != nil {
		return diag.FromErr(err)
	}

	usageList := make([]interface{}, 0, len(acct.Usage))
	for i := range acct.Usage {
		usageList = append(usageList, map[string]interface{}{
			accountLimitAttr: acct.Usage[i].Limit,
			accountTypeAttr:  acct.Usage[i].Type,
			accountUsedAttr:  acct.Usage[i].Used,
		})
	}
	if err := d.Set(accountUsageAttr, usageList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(accountAddress1Attr, acct.Address1); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountAddress2Attr, acct.Address2); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountCCEmailAttr, acct.CCEmail); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountCityAttr, acct.City); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountCountryAttr, acct.Country); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountDescriptionAttr, acct.Description); err != nil {
		return diag.FromErr(err)
	}

	invitesList := make([]interface{}, 0, len(acct.Invites))
	for i := range acct.Invites {
		invitesList = append(invitesList, map[string]interface{}{
			accountEmailAttr: acct.Invites[i].Email,
			accountRoleAttr:  acct.Invites[i].Role,
		})
	}
	if err := d.Set(accountInvitesAttr, invitesList); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(accountNameAttr, acct.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountStateProvAttr, acct.StateProv); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(accountTimezoneAttr, acct.Timezone); err != nil {
		return diag.FromErr(err)
	}

	usersList := make([]interface{}, 0, len(acct.Users))
	for i := range acct.Users {
		usersList = append(usersList, map[string]interface{}{
			accountUserIDAttr: acct.Users[i].UserCID,
			accountRoleAttr:   acct.Users[i].Role,
		})
	}
	if err := d.Set(accountUsersAttr, usersList); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
