package cmd

import (
	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	GCPDomain string
	GCPApiKey string
)

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Modules de reconnaissance et d'audit pour Google Cloud / Workspace",
	Run: func(cmd *cobra.Command, args []string) {
		folder := "templates/gcp"
		utils.LogInfo("Chargement des templates GCP depuis : %s", folder)

		tmplList, err := templates.ParseDirectory(folder)
		if err != nil {
			utils.LogError("Erreur lors de la lecture du dossier : %v", err)
			return
		}

		utils.LogInfo("%d template(s) détecté(s).", len(tmplList))

		for _, tmpl := range tmplList {
			if tmpl.Engine == "gcp" {
				err := engines.ExecuteGCPEngine(tmpl, GCPDomain, GCPApiKey)
				if err != nil {
					utils.LogError("Erreur d'exécution [%s] : %v", tmpl.ID, err)
				}
			}
		}
	},
}

func init() {
	gcpCmd.Flags().StringVarP(&GCPDomain, "domain", "d", "", "Domaine cible à auditer (ex: entreprise.com)")
	gcpCmd.Flags().StringVar(&GCPApiKey, "apikey", "", "Clé API GCP à valider (AIZA...)")
	rootCmd.AddCommand(gcpCmd)
}
