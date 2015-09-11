package commands

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EnergyWatch = &cobra.Command{
	Use:   "energywatch",
	Short: "energywatch is a tool to push energy data to WattVision",
	Long: `energywatch takes data from the Rainforest Raven and
pushes it to the wattvision API. energywatch requires the sensor ID, API ID,
and API key from wattvision, which can be accessed from the account settings.`,
}

func Execute() {
	InitConfig()
	AddCommands()

	EnergyWatch.Execute()
}

func AddCommands() {
	EnergyWatch.AddCommand(versionCmd)
	EnergyWatch.AddCommand(watchCmd)
}

func InitConfig() {
	viper.SetConfigFile("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
}
