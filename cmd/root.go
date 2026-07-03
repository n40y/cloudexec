package cmd

import (
	"fmt"
	"os"

	"cloudexec/pkg/engines"
	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	Target       string
	Verbose      bool
	DirectIP     string
	UserAgentRot bool
)

var rootCmd = &cobra.Command{
	Use:   "cloudexec",
	Short: "Cloudexec est un outil d'audit Cloud.",
	Long:  `Un scanner Cloud CLI modulaire et ultra-rapide écrit en Go, utilisant des templates YAML.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// S'exécute automatiquement au début de TOUTES les commandes
		utils.Banner()
		engines.GlobalDirectIP = DirectIP
		engines.GlobalAntiWaf = UserAgentRot
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Flags globaux partagés par toutes les sous-commandes
	rootCmd.PersistentFlags().StringVarP(&Target, "target", "t", "", "Domaine ou PI cible")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Mode verbeux")
	rootCmd.PersistentFlags().StringVar(&DirectIP, "ip", "", "Forcer une IP cible directe (Bypass Cloudflare/DNS Pinning)")
	rootCmd.PersistentFlags().BoolVar(&UserAgentRot, "anti-waf", false, "Activer la rotation de User-Agent et le spoofing de headers")
}
