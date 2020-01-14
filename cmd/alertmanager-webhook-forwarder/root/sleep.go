package root

import (
	"fmt"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/conditional_runner"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"strconv"
	"time"
)

var SleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Sleep some seconds",
	Args:  cobra.ExactArgs(1),
	Long: `
alertmanager-webhook-forwarder sleep 5

`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		sleepPeriodSeconds, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		if sleepPeriodSeconds < 0 {
			return fmt.Errorf("provided argument must be a positive number, %v given", sleepPeriodSeconds)
		}

		klog.V(7).Infof("sleeping for %d second(s)", sleepPeriodSeconds)
		conditional_runner.NotDryRun().Run(func() {
			time.Sleep(time.Duration(sleepPeriodSeconds) * time.Second)
		})

		return nil
	},
}

func init() {
	RootCmd.AddCommand(SleepCmd)
}
