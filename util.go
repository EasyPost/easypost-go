package easypost

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// StringPtr returns a pointer to a string with the given value.
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a bool with the given value.
func BoolPtr(b bool) *bool {
	return &b
}

// containsString returns true if the given list contains the given element.
func containsString(strings []string, targetString string) bool {
	for _, aString := range strings {
		if aString == targetString {
			return true
		}
	}
	return false
}

type RateFilterRate struct {
	ID      string `json:"id,omitempty"`
	Service string `json:"service,omitempty"`
	Carrier string `json:"carrier,omitempty"`
	Rate    string `json:"rate,omitempty"`
}

// lowestSmartRate returns the lowest smartrate from the given list of smartrates.
func (c *Client) lowestSmartRate(rates []*SmartRate, deliveryDays int, deliveryAccuracy string) (out SmartRate, err error) {
	validDeliveryAccuracies := []string{"percentile_50", "percentile_75", "percentile_85", "percentile_90", "percentile_95",
		"percentile_97", "percentile_99"}

	// if the delivery accuracy is not valid, return an error
	if !containsString(validDeliveryAccuracies, deliveryAccuracy) {
		return out, fmt.Errorf("invalid delivery accuracy: %s", deliveryAccuracy)
	}

	for _, rate := range rates {
		smartrateDeliveryDay := -1

		switch deliveryAccuracy {
		case "percentile_50":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile50
		case "percentile_75":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile75
		case "percentile_85":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile85
		case "percentile_90":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile90
		case "percentile_95":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile95
		case "percentile_97":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile97
		case "percentile_99":
			smartrateDeliveryDay = rate.TimeInTransit.Percentile99
		default:
			break
		}

		// if we could not find the correct delivery date variable for the given delivery accuracy, skip it
		if smartrateDeliveryDay == -1 {
			return out, fmt.Errorf("invalid delivery accuracy: %s", deliveryAccuracy)
		}

		// if this rate's delivery days is greater than the requested delivery days, skip it
		if smartrateDeliveryDay > deliveryDays {
			continue
		}

		// if lowest rate is null, set it to this rate
		if (out == SmartRate{}) {
			out = *rate
			continue
		}

		// if this rate is lower than the lowest rate, set it to this rate
		if 0 < rate.Rate && rate.Rate < out.Rate {
			out = *rate
		}
	}

	// if not rate was ever set (nothing matched the criteria), return an error
	if (out == SmartRate{}) {
		return out, errors.New("no rates found")
	}

	return
}

// lowestRate returns the lowest rate from the given list of rates with carrier and service filters.
func (c *Client) lowestRate(rates []*RateFilterRate, carriers []string, services []string) (out RateFilterRate, err error) {
	carriersMap, servicesMap := make(map[string]bool), make(map[string]bool)

	for _, carrier := range carriers {
		carriersMap[strings.ToLower(carrier)] = true
	}

	for _, service := range services {
		servicesMap[strings.ToLower(service)] = true
	}

	for _, rate := range rates {
		// if this rate's carrier is not in the carrier list or this rate's service is not in the service list, skip it
		if len(carriersMap) > 0 && !carriersMap[strings.ToLower(rate.Carrier)] ||
			len(servicesMap) > 0 && !servicesMap[strings.ToLower(rate.Service)] {
			continue
		}

		currentRate, _ := strconv.ParseFloat(out.Rate, 32)
		newRate, _ := strconv.ParseFloat(rate.Rate, 32)

		// if lowest rate is null, set it to this rate
		if (out == RateFilterRate{}) {
			out = *rate
			continue
		}

		// if this rate is lower than the lowest rate, set it to this rate
		if 0 < newRate && newRate < currentRate {
			out = *rate
		}
	}

	// if not rate was ever set (nothing matched the criteria), return an error
	if (out == RateFilterRate{}) {
		return out, errors.New("no rates found")
	}

	return
}

// lowestObjectRate returns the lowest rate from the given list of rates.
func (c *Client) lowestObjectRate(rates []*Rate, carriers []string, services []string) (out Rate, err error) {
	filterRates := make([]*RateFilterRate, 0)
	for _, rate := range rates {
		filterRates = append(filterRates, &RateFilterRate{
			ID:      rate.ID,
			Service: rate.Service,
			Carrier: rate.Carrier,
			Rate:    rate.Rate,
		})
	}

	lowestRate, err := c.lowestRate(filterRates, carriers, services)
	if err == nil {
		for _, rate := range rates {
			if rate.ID == lowestRate.ID {
				return *rate, nil
			}
		}
	}
	return
}

// lowestPickupRate returns the lowest pickup rate from the given list of pickup rates.
func (c *Client) lowestPickupRate(rates []*PickupRate, carriers []string, services []string) (out PickupRate, err error) {
	filterRates := make([]*RateFilterRate, 0)
	for _, rate := range rates {
		filterRates = append(filterRates, &RateFilterRate{
			ID:      rate.ID,
			Service: rate.Service,
			Carrier: rate.Carrier,
			Rate:    rate.Rate,
		})
	}

	lowestRate, err := c.lowestRate(filterRates, carriers, services)
	if err == nil {
		for _, rate := range rates {
			if rate.ID == lowestRate.ID {
				return *rate, nil
			}
		}
	}
	return
}
