package commands

import (
	"fmt"
	"github.com/aldernero/timebox/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dbFile  string
	tb      util.TimeBox
)

type CliFlags struct {
	boxName     string
	minDuration time.Duration
	maxDuration time.Duration
	startTime   string
	endTime     string
	period      util.TimePeriod
	force       bool
}

var cliFlags CliFlags

var rootCmd = &cobra.Command{
	Use:              "",
	Short:            "",
	Long:             "",
	PersistentPreRun: nil,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// persistent flags
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.timebox.yaml)")
	rootCmd.PersistentFlags().StringVar(&dbFile, "db", "", "database file (default is $PWD/timebox.db)")

	// subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(updateCmd)
}

func initConfig() {
	viper.SetEnvPrefix("timebox")
	viper.AutomaticEnv()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".timebox" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".timebox")
		cfgPath := path.Join(home, ".timebox.yaml")
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			viper.SetDefault("HoursInWeek", 168)
			viper.SetDefault("TimePeriod", "week")
			if err := viper.SafeWriteConfigAs(cfgPath); err != nil {
				fmt.Println("Can't write config:", err)
				os.Exit(1)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	if dbFile == "" {
		dbFile = path.Join(".", "timebox.db")
	}
	tb = util.TimeBoxFromDB(dbFile)
}
