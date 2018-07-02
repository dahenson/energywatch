package commands

import (
	"log"

	"github.com/dahenson/energywatch/wattvision"
	"github.com/dahenson/goraven"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch and push data from the Raven to WattVision.",
	Long: `Watch the Raven for energy reads, and push the data to
WattVision. This feature requires that the Raven is specified in the
configuration file, and that the appropriate API keys and IDs are
provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		watch()
	},
}

func watch() {
	raven, err := goraven.Connect(viper.GetString("dev"))
	if err != nil {
		log.Fatal(err)
	}

	defer raven.Close()

	e := wattvision.EnergyData{}
	e.SensorId = viper.GetString("sensor_id")
	e.ApiId = viper.GetString("api_id")
	e.ApiKey = viper.GetString("api_key")

	for {
		notify, err := raven.Receive()
		if err != nil {
			log.Println(err)
		}

		switch t := notify.(type) {
		case *goraven.ConnectionStatus:
			log.Printf("Connection Status: %s", t.Status)
		case *goraven.CurrentSummationDelivered:
			pushCurrentSummationDelivered(t, e)
		case *goraven.InstantaneousDemand:
			pushInstantaneousDemand(t, e)
		default:
		}
	}
}

func pushInstantaneousDemand(c *goraven.InstantaneousDemand, e wattvision.EnergyData) {
	watts, err := c.GetDemand()
	if err != nil {
		log.Printf("Instantaneous Demand Data Failure: %s\n", err)
		return
	}
	e.Watts = watts
	pushEnergyData(e)
}

func pushCurrentSummationDelivered(c *goraven.CurrentSummationDelivered, e wattvision.EnergyData) {
	watthours, err := c.GetSummationDelivered()
	if err != nil {
		log.Printf("Current Summation Data Failure: %s\n", err)
		return
	}
	e.WattHours = watthours
	pushEnergyData(e)
}

func pushEnergyData(e wattvision.EnergyData) {
	err := wattvision.PushEnergyData(e)
	if err != nil {
		log.Printf("Unable to push data: %s\n", err)
	}
}
