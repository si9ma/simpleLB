/*
Copyright © 2020 si9ma <si9ma@si9ma.com>
*/
package cmd

import (
	"fmt"
	"github.com/si9ma/simpleLB/config"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	lbConfig config.LBConfig
	cfgFile  string
	viperVar *viper.Viper
)

var rootCmd = &cobra.Command{
	Use:   "simpleLB",
	Short: "A simple and lite load balancer",
	Long:  `A simple and lite load balancer(support automatic access to certificates from Let's Encrypt )`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simpleLB.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// the keys may contain dot (the default key delimiter), change the delimiter)
	viperVar = viper.NewWithOptions(viper.KeyDelimiter("::"))
	if cfgFile != "" {
		// Use config file from the flag.
		viperVar.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "get home dir failed:%s", err.Error())
			os.Exit(1)
		}

		// Search config in home directory with name ".simpleLB" (without extension).
		viperVar.AddConfigPath(home)
		viperVar.SetConfigName(".simpleLB")
	}

	viperVar.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viperVar.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viperVar.ConfigFileUsed())
	}

	// Unmarshal to struct
	if err := viperVar.Unmarshal(&lbConfig); err != nil {
		fmt.Fprintf(os.Stderr, "parse config failed:%s", err.Error())
		os.Exit(1)
	}
}
