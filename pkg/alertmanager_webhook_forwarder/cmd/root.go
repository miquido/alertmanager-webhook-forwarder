package cmd

import (
	goflag "flag"
	"fmt"
	"os"
	"strconv"

	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forwarder"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/hangouts_chat"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/log"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog"
)

var cfgFile string
var verbose bool

type viperConfig struct {
	Viper map[string]interface{} `json:"viper"`
}

var RootCmd = &cobra.Command{
	Use:           "alertmanager-webhook-forwarder",
	Short:         "",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `
Example:

	$ ./alertmanager-webhook-forwarder forward hangouts-chat "test message" --template gotemplate-yaml --verbose --dry-run

`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.DumpMultilineText("Error", "message", fmt.Sprint(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		configureVerbosity,
		configureFlags,
		initConfig,
		showDebugInfo,
		initializeForwarderRegistry,
	)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alertmanager-webhook-forwarder.yaml)")
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "When true sets verbosity level to 7")
	RootCmd.PersistentFlags().Bool("dry-run", false, "Do not run actions that change any state")

	// Configure Klog flags
	klog.InitFlags(nil)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func configureVerbosity() {
	if !verbose {
		return
	}
	vFlag := RootCmd.PersistentFlags().Lookup("v")
	if vFlag == nil {
		// not possible
		return
	}
	verbosity := RootCmd.PersistentFlags().Lookup("v").Value.String()
	if v, err := strconv.Atoi(verbosity); err == nil {
		if v < 7 {
			_ = vFlag.Value.Set("7")
		}
	}
}

func configureFlags() {
	_ = viper.BindPFlag("verbosity", RootCmd.PersistentFlags().Lookup("v"))
	_ = viper.BindPFlag("dryRun", RootCmd.PersistentFlags().Lookup("dry-run"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".alertmanager-webhook-forwarder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".alertmanager-webhook-forwarder")
	}

	viperSetDefaults()
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		klog.V(7).Info("Using config file:", viper.ConfigFileUsed())
	}
}

func showDebugInfo() {
	klog.V(7).Infof("dryRun=%t", viper.GetBool("dryRun"))
	klog.V(7).Infof("verbosity=%d", viper.GetInt("verbosity"))
	conditional_runner.V(10).DumpYaml("viper config", viperConfig{viper.AllSettings()})
}

func viperSetDefaults() {
	viper.SetDefault("templates.hangoutsChat.alertmanager", hangouts_chat.DefaultTemplateAlertmanger)
	viper.SetDefault("templates.hangoutsChat.jsonnet", hangouts_chat.DefaultTemplateJsonnet)
	viper.SetDefault("templates.hangoutsChat.goTemplateYaml", hangouts_chat.DefaultTemplateGoTemplateYaml)
	viper.SetDefault("templates.hangoutsChat.goTemplateText", hangouts_chat.DefaultTemplateGoTemplateText)
}

func initializeForwarderRegistry() {
	var hangoutsChatForwarder forwarder.Forwarder
	hangoutsChatForwarder, err := hangouts_chat.NewHangoutsChatForwarder()
	if err != nil {
		klog.Fatal(err)
	}
	forwarder.Attach(utils.ProviderHangoutsChat, &hangoutsChatForwarder)
}
