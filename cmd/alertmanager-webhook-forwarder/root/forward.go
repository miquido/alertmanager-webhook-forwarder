package root

import (
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/forwarder"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/message_template"
	"github.com/miquido/alertmanager-webhook-forwarder/pkg/utils"
	"github.com/spf13/cobra"
)

// ForwardCmd represents the serve command
var ForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Forward a messages to receivers",
}

var ForwardHangoutsChatCmd = &cobra.Command{
	Use:   "hangouts-chat [text]",
	Short: "",
	Args:  cobra.MinimumNArgs(1),
	Long:  ``,
	RunE:  ForwardHangoutsChat,
}

func ForwardHangoutsChat(cmd *cobra.Command, args []string) (err error) {
	input := message_template.NewMessageTemplateInputFromArgs(args)
	templateName, _ := cmd.Flags().GetString("template")
	channelName, _ := cmd.Flags().GetString("channel")

	return forwarder.Forward(utils.ProviderHangoutsChat, channelName, templateName, input)
}

func init() {
	ForwardHangoutsChatCmd.Flags().String("channel", "", "Hangouts Chat defined channel.")
	ForwardHangoutsChatCmd.Flags().String("template", "jsonnet", "Use defined templates. Available by default: jsonnet, gotemplateyaml, gotemplatetext")

	ForwardCmd.AddCommand(
		ForwardHangoutsChatCmd,
	)

	RootCmd.AddCommand(ForwardCmd)
}
