package cmd

import (
	"fmt"
	"os"

	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	Target  string
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "cloudexec",
	Short: "Cloudexec est un outil d'audit Cloud.",
	Long:  `Un scanner Cloud CLI modulaire et ultra-rapide écrit en Go, utilisant des templates YAML.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// S'exécute automatiquement au début de TOUTES les commandes
		utils.Banner()
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
}
