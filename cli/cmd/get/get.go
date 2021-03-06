package get

import (
	"github.com/jbvmio/kafkactl/cli/cmd/broker"
	"github.com/jbvmio/kafkactl/cli/cmd/group"
	"github.com/jbvmio/kafkactl/cli/cmd/lag"
	"github.com/jbvmio/kafkactl/cli/cmd/msg"
	"github.com/jbvmio/kafkactl/cli/cmd/metrics"
	"github.com/jbvmio/kafkactl/cli/cmd/topic"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var outFlags out.OutFlags

var CmdGet = &cobra.Command{
	Use:     "get",
	Example: "  kafkactl get topic <topicName>\n  kafkactl get lag",
	Short:   "Get Kafka Information",
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case len(args) > 0:
			out.Failf("No such resource: %v", args[0])
		default:
			cmd.Help()
		}
	},
}

func init() {
	CmdGet.PersistentFlags().StringVarP(&outFlags.Format, "out", "o", "", "Change Output Format - yaml|json.")

	CmdGet.AddCommand(broker.CmdGetBroker)
	CmdGet.AddCommand(broker.CmdGetApiVers)
	CmdGet.AddCommand(topic.CmdGetTopic)
	CmdGet.AddCommand(group.CmdGetGroup)
	CmdGet.AddCommand(group.CmdGetMember)
	CmdGet.AddCommand(lag.CmdGetLag)
	CmdGet.AddCommand(msg.CmdGetMsg)
	CmdGet.AddCommand(metrics.CmdGetMetrics)
}
