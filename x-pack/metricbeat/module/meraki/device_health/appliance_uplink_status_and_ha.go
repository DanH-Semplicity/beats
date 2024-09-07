package device_health

import (
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/elastic-agent-libs/mapstr"

	meraki_api "github.com/meraki/dashboard-api-go/v3/sdk"
)

func reportApplianceUplinkStatuses(reporter mb.ReporterV2, organizationID string, devices map[Serial]*Device, responseApplianceUplinkStatuses *meraki_api.ResponseApplianceGetOrganizationApplianceUplinkStatuses) {

	metrics := []mapstr.M{}

	for _, uplink := range *responseApplianceUplinkStatuses {

		if device, ok := devices[Serial(uplink.Serial)]; ok {
			metric := mapstr.M{
				"appliance.uplink.high_availablity.enabled": uplink.HighAvailability.Enabled,
				"appliance.uplink.high_availablity.role":    uplink.HighAvailability.Role,
				"appliance.uplink.last_reported_at":         uplink.LastReportedAt,
				"appliance.address":                         device.Address,
				"appliance.firmware":                        device.Firmware,
				"appliance.imei":                            device.Imei,
				"appliance.lan_ip":                          device.LanIP,
				"appliance.location":                        device.Location,
				"appliance.mac":                             device.Mac,
				"appliance.model":                           device.Model,
				"appliance.name":                            device.Name,
				"appliance.network_id":                      device.NetworkID,
				"appliance.notes":                           device.Notes,
				"appliance.product_type":                    device.ProductType,
				"appliance.serial":                          device.Serial,
				"appliance.tags":                            device.Tags,
			}

			for _, item := range *uplink.Uplinks {
				metrics = append(metrics, mapstr.Union(metric, mapstr.M{
					"appliance.uplink.interface":      item.Interface,
					"appliance.uplink.status":         item.Status,
					"appliance.uplink.ip":             item.IP,
					"appliance.uplink.gateway":        item.Gateway,
					"appliance.uplink.public_ip":      item.PublicIP,
					"appliance.uplink.primary_dns":    item.PrimaryDNS,
					"appliance.uplink.secondary_dns":  item.SecondaryDNS,
					"appliance.uplink.ip_assigned_by": item.IPAssignedBy,
				}))

			}
		}
	}

	ReportMetricsForOrganization(reporter, organizationID, metrics)
}
