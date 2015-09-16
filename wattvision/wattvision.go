package wattvision

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type EnergyData struct {
	SensorId  string  `json:"sensor_id"`
	ApiId     string  `json:"api_id"`
	ApiKey    string  `json:"api_key"`
	Time      string  `json:"time,omitempty"`
	Watts     float64 `json:"watts,omitempty"`
	WattHours float64 `json:"watthours,omitempty"`
}

const URL = `http://www.wattvision.com/api/v0.2/elec`

func PushEnergyData(data EnergyData) error {
	buf, err := json.Marshal(data)
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
