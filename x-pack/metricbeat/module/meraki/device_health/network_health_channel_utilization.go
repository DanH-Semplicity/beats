package device_health

import (
	"fmt"
	"strings"

	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/elastic/elastic-agent-libs/mapstr"

	meraki_api "github.com/meraki/dashboard-api-go/v3/sdk"
)

func getNetworkHealthChannelUtilization(client *meraki_api.Client, networks *meraki_api.ResponseOrganizationsGetOrganizationNetworks) ([]*meraki_api.ResponseNetworksGetNetworkNetworkHealthChannelUtilization, error) {

	var networkHealthUtilizations []*meraki_api.ResponseNetworksGetNetworkNetworkHealthChannelUtilization

	for _, network := range *networks {

		for _, product_type := range network.ProductTypes {

			if strings.Compare(product_type, "wireless") == 0 {
				networkHealthUtilization, res, err := client.Networks.GetNetworkNetworkHealthChannelUtilization(network.ID, &meraki_api.GetNetworkNetworkHealthChannelUtilizationQueryParams{})
				if err != nil {
					//"This endpoint is only available for networks on MR 27.0 or above."
					// We just swallow this error but do not append to the list
					if !(strings.Contains(string(res.Body()), "MR 27.0")) {
						//Any other problem we are going to return an error
						return nil, fmt.Errorf("Networks.GetNetworkNetworkHealthChannelUtilization failed; [%d] %s. %w", res.StatusCode(), res.Body(), err)
					}
				} else {
					networkHealthUtilizations = append(networkHealthUtilizations, networkHealthUtilization)
				}

			}
		}
	}

	return networkHealthUtilizations, nil
}

func reportNetworkHealthChannelUtilization(reporter mb.ReporterV2, organizationID string, devices map[Serial]*Device, networkHealthUtilizations []*meraki_api.ResponseNetworksGetNetworkNetworkHealthChannelUtilization) {
	metrics := []mapstr.M{}
	for _, networkHealthUtil := range networkHealthUtilizations {
		for _, network := range *networkHealthUtil {

			metric := mapstr.M{
				"network.health.channel.radio.serial": network.Serial,
				"network.health.channel.radio.model":  network.Model,
				"network.health.channel.radio.tags":   network.Tags,
			}

			for _, wifi0 := range *network.Wifi0 {
				metric["network.health.channel.radio.wifi0.start_time"] = wifi0.StartTime
				metric["network.health.channel.radio.wifi0.end_time"] = wifi0.EndTime
				metric["network.health.channel.radio.wifi0.utilization80211"] = wifi0.Utilization80211
				metric["network.health.channel.radio.wifi0.utilizationNon80211"] = wifi0.UtilizationNon80211
				metric["network.health.channel.radio.wifi0.utilizationTotal"] = wifi0.UtilizationTotal
			}

			for _, wifi1 := range *network.Wifi1 {
				metric["network.health.channel.radio.wifi1.start_time"] = wifi1.StartTime
				metric["network.health.channel.radio.wifi1.end_time"] = wifi1.EndTime
				metric["network.health.channel.radio.wifi1.utilization80211"] = wifi1.Utilization80211
				metric["network.health.channel.radio.wifi1.utilizationNon80211"] = wifi1.UtilizationNon80211
				metric["network.health.channel.radio.wifi1.utilizationTotal"] = wifi1.UtilizationTotal
			}

			if device, ok := devices[Serial(network.Serial)]; ok {
				metric["device.address"] = device.Address
				metric["device.firmware"] = device.Firmware
				metric["device.imei"] = device.Imei
				metric["device.lan_ip"] = device.LanIP
				metric["device.location"] = device.Location
				metric["device.mac"] = device.Mac
				metric["device.model"] = device.Model
				metric["device.name"] = device.Name
				metric["device.network_id"] = device.NetworkID
				metric["device.notes"] = device.Notes
				metric["device.product_type"] = device.ProductType
				metric["device.serial"] = device.Serial
				metric["device.tags"] = device.Tags

			}
			metrics = append(metrics, metric)
		}
	}
	ReportMetricsForOrganization(reporter, organizationID, metrics)
}
