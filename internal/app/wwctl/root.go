package wwctl

import (
	"os"

	"github.com/hpcng/warewulf/internal/app/wwctl/configure"
	"github.com/hpcng/warewulf/internal/app/wwctl/container"
	"github.com/hpcng/warewulf/internal/app/wwctl/genconf"
	"github.com/hpcng/warewulf/internal/app/wwctl/kernel"
	"github.com/hpcng/warewulf/internal/app/wwctl/node"
	"github.com/hpcng/warewulf/internal/app/wwctl/overlay"
	"github.com/hpcng/warewulf/internal/app/wwctl/power"
	"github.com/hpcng/warewulf/internal/app/wwctl/profile"
	"github.com/hpcng/warewulf/internal/app/wwctl/server"
	"github.com/hpcng/warewulf/internal/app/wwctl/ssh"
	"github.com/hpcng/warewulf/internal/app/wwctl/version"
	"github.com/hpcng/warewulf/internal/pkg/help"
	warewulfconf "github.com/hpcng/warewulf/internal/pkg/config"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		DisableFlagsInUseLine: true,
		Use:                   "wwctl COMMAND [OPTIONS]",
		Short:                 "Warewulf Control",
		Long:                  "Control interface to the Warewulf Cluster Provisioning System.",
		PersistentPreRunE:     rootPersistentPreRunE,
		SilenceUsage:          true,
		SilenceErrors:         true,
	}
	verboseArg      bool
	DebugFlag       bool
	LogLevel        int
	WarewulfConfArg string
	AllowEmptyConf  bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseArg, "verbose", "v", false, "Run with increased verbosity.")
	rootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "d", false, "Run with debugging messages enabled.")
	rootCmd.PersistentFlags().IntVar(&LogLevel, "loglevel", wwlog.INFO, "Set log level to given string")
	_ = rootCmd.PersistentFlags().MarkHidden("loglevel")
	rootCmd.PersistentFlags().StringVar(&WarewulfConfArg, "warewulfconf", "", "Set the warewulf configuration file")
	rootCmd.PersistentFlags().BoolVar(&AllowEmptyConf, "emptyconf", false, "Allow empty configuration")
	_ = rootCmd.PersistentFlags().MarkHidden("emptyconf")
	rootCmd.SetUsageTemplate(help.UsageTemplate)
	rootCmd.SetHelpTemplate(help.HelpTemplate)
	rootCmd.AddCommand(overlay.GetCommand())
	rootCmd.AddCommand(container.GetCommand())
	rootCmd.AddCommand(node.GetCommand())
	rootCmd.AddCommand(kernel.GetCommand())
	rootCmd.AddCommand(power.GetCommand())
	rootCmd.AddCommand(profile.GetCommand())
	rootCmd.AddCommand(configure.GetCommand())
	rootCmd.AddCommand(server.GetCommand())
	rootCmd.AddCommand(version.GetCommand())
	rootCmd.AddCommand(ssh.GetCommand())
	rootCmd.AddCommand(genconf.GetCommand())
}

// GetRootCommand returns the root cobra.Command for the application.
func GetRootCommand() *cobra.Command {
	return rootCmd
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) (err error) {
	if DebugFlag {
		wwlog.SetLogLevel(wwlog.DEBUG)
	} else if verboseArg {
		wwlog.SetLogLevel(wwlog.VERBOSE)
	} else {
		wwlog.SetLogLevel(wwlog.INFO)
	}
	if LogLevel != wwlog.INFO {
		wwlog.SetLogLevel(LogLevel)
	}
	conf := warewulfconf.Get()
	if !AllowEmptyConf && !conf.Initialized() {
		if WarewulfConfArg != "" {
			err = conf.ReadConf(WarewulfConfArg)
		} else if os.Getenv("WAREWULFCONF") != "" {
			err = conf.ReadConf(os.Getenv("WAREWULFCONF"))
		} else {
			err = conf.ReadConf(warewulfconf.ConfigFile)
		}
	} else {
		err = conf.SetDynamicDefaults()
	}
	return
}
