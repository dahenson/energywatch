package commands

import (
	"log"

	"github.com/dahenson/energywatch/energyapis"
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

	e := energyapis.WattVision{}
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

func pushInstantaneousDemand(c *goraven.InstantaneousDemand, e energyapis.EnergyAPI) {
	kilowatts, err := c.GetDemand()
	if err != nil {
		log.Printf("Instantaneous Demand Data Failure: %s\n", err)
		return
	}
	if err := e.PushInstantaneousDemand(kilowatts * 1000); err != nil {
		log.Printf("Data Push Failed: %s\n", err)
	}

	if verbose {
		log.Printf("Instantaneous Demand Published to WattVision: %f Watts", kilowatts*1000)
	}
}

func pushCurrentSummationDelivered(c *goraven.CurrentSummationDelivered, e energyapis.EnergyAPI) {
	kilowatthours, err := c.GetSummationDelivered()
	if err != nil {
		log.Printf("Current Summation Data Failure: %s\n", err)
		return
	}
	e.PushCurrentSummationDelivered(kilowatthours * 1000)

	if verbose {
		log.Printf("Current Summation Published to WattVision: %f WattHours", kilowatthours*1000)
	}
}
