package policy

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

type statusCheckPolicySettings struct {
	StatusName               string `json:"statusName"`
	InvalidateOnSourceUpdate bool `json:"invalidateOnSourceUpdate"`
	PolicyDisplayName        string `json:"displayName"`
}

const (
	statusName               = "status_name"
	invalidateOnSourceUpdate = "invalidate_on_source_update"
	policyDisplayName        = "display_name"
)

// ResourceBranchPolicyStatusCheck schema and implementation for status check policy resource
func ResourceBranchPolicyStatusCheck() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: statusCheckFlattenFunc,
		ExpandFunc:  statusCheckExpandFunc,
		PolicyType:  StatusCheck,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema[statusName] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	settingsSchema[invalidateOnSourceUpdate] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  false,
		Optional: true,
	}
	settingsSchema[policyDisplayName] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
	}

	return resource
}

func statusCheckFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := statusCheckPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings[statusName] = policySettings.StatusName
	settings[invalidateOnSourceUpdate] = policySettings.InvalidateOnSourceUpdate
	settings[policyDisplayName] = policySettings.PolicyDisplayName
	d.Set(SchemaSettings, settingsList)
	return nil
}

func statusCheckExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["statusName"] = settings[statusName].(string)
	policySettings["invalidateOnSourceUpdate"] = settings[invalidateOnSourceUpdate].(bool)
	policySettings["displayName"] = settings[policyDisplayName].(string)

	return policyConfig, projectID, nil
}
