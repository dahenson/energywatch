package commands

import (
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	device   string
	sensorId string
	apiId    string
	apiKey   string
	verbose  bool
)

var energyWatch = &cobra.Command{
	Use:   "energywatch",
	Short: "energywatch is a tool to push energy data to WattVision",
	Long: `energywatch takes data from the Rainforest Raven and
pushes it to the wattvision API. energywatch requires the sensor ID, API ID,
and API key from wattvision, which can be accessed from the account settings.`,
}

func Execute() {
	AddCommands()

	energyWatch.Execute()
}

func AddCommands() {
	energyWatch.AddCommand(versionCmd)
	energyWatch.AddCommand(watchCmd)
}

func init() {
	cobra.OnInitialize(initConfig)
	energyWatch.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.json)")
	energyWatch.PersistentFlags().StringVar(&device, "dev", "", "raven device (e.g. /dev/ttyUSB0)")
	energyWatch.PersistentFlags().StringVar(&sensorId, "sensor_id", "", "sensor id for wattvision")
	energyWatch.PersistentFlags().StringVar(&apiId, "api_id", "", "API ID for wattvision")
	energyWatch.PersistentFlags().StringVar(&apiKey, "api_key", "", "API Key for wattvision")
	watchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Run in verbose mode for debugging")
	viper.BindPFlag("dev", energyWatch.PersistentFlags().Lookup("dev"))
	viper.BindPFlag("sensor_id", energyWatch.PersistentFlags().Lookup("sensor_id"))
	viper.BindPFlag("api_id", energyWatch.PersistentFlags().Lookup("api_id"))
	viper.BindPFlag("api_key", energyWatch.PersistentFlags().Lookup("api_key"))
}

func initConfig() {
	viper.SetConfigType("json")

	home, err := homedir.Dir()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".energywatch")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}
}
