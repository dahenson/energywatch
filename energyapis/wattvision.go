package energyapis

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type WattVision struct {
	SensorId  string  `json:"sensor_id"`
	ApiId     string  `json:"api_id"`
	ApiKey    string  `json:"api_key"`
	time      string  `json:"time,omitempty"`
	watts     float64 `json:"watts,omitempty"`
	watthours float64 `json:"watthours,omitempty"`
}

const URL = `http://www.wattvision.com/api/v0.2/elec`

/*
 * PushWatts pushes only instantaneous consumption to wattvision
 */
func (w WattVision) PushInstantaneousDemand(watts float64) error {
	w.watthours = 0.0
	w.watts = watts
	return w.pushEnergyData()
}

/*
 * PushWattHours pushes only total consumption to WattVision
 */
func (w WattVision) PushCurrentSummationDelivered(watthours float64) error {
	w.watts = 0.0
	w.watthours = watthours
	return w.pushEnergyData()
}

func (w WattVision) pushEnergyData() error {
	buf, err := json.Marshal(w)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Post(URL, "text/json", bytes.NewReader(buf))
	if err != nil {
		return err
	}

	if resp.StatusCode == 400 {
		return errors.New("Wattvision API error")
	}

	defer resp.Body.Close()

	return nil
}
