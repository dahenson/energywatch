package commands

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/dahenson/energywatch/wattvision"
	"github.com/dahenson/goraven"
)

var watchCmd = &cobra.Command {
	Use: "watch",
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
	if !viper.IsSet("raven_dev") {
		log.Fatal("Please specify a Raven device to connect to")
	}

	raven, err := goraven.Connect(viper.GetString("raven_dev"))
	if err != nil {
		log.Fatal(err)
	}

	defer raven.Close()

	for {
		notify, err := raven.Receive()
		if err != nil {
			log.Println(err)
		}

		switch t := notify.(type) {
			case *goraven.ConnectionStatus:
				log.Println("Raven not connected to meter.")
			case *goraven.CurrentSummationDelivered:
				pushCurrentSummationDelivered(t)
			default:
		}
	}
}

func pushCurrentSummationDelivered(c *goraven.CurrentSummationDelivered) {
	watthours, err := c.GetSummationDelivered()
	if err != nil {
		log.Println("Current Summation Data Failure: %s", err)
		return
	}
	e := wattvision.EnergyData{WattHours: watthours}
	pushEnergyData(e)
}

func pushEnergyData(e wattvision.EnergyData) {
	e.SensorId = viper.GetString("sensor_id")
	e.ApiId = viper.GetString("api_id")
	e.ApiKey = viper.GetString("api_key")

	err := wattvision.PushEnergyData(e)
	if err != nil {
		log.Println("Unable to push current summation: %s", err)
	}
}
