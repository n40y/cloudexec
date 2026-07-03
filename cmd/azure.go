package cmd

import (
	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	TargetDomain string
	Username     string
	Password     string
	TenantID     string
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Modules de reconnaissance et d'audit pour Azure / M365",
	Run: func(cmd *cobra.Command, args []string) {
		folder := "templates/azure"
		utils.LogInfo("Chargement des templates Azure depuis : %s", folder)

		tmplList, err := templates.ParseDirectory(folder)
		if err != nil {
			utils.LogError("Erreur lors de la lecture du dossier : %v", err)
			return
		}

		utils.LogInfo("%d template(s) détecté(s).", len(tmplList))

		for _, tmpl := range tmplList {
			if tmpl.Engine == "azure" {
				err := engines.ExecuteAzureEngine(tmpl, TargetDomain, Username, Password, TenantID)
				if err != nil {
					utils.LogError("Erreur d'exécution [%s] : %v", tmpl.ID, err)
				}
			}
		}
	},
}

func init() {
	// Utilisation de StringVarP pour mapper le flag long et son raccourci short
	azureCmd.Flags().StringVarP(&TargetDomain, "domain", "d", "", "Domaine cible (Recon OSINT)")
	azureCmd.Flags().StringVar(&Username, "username", "", "Nom d'utilisateur Azure/M365 (Auth)")
	azureCmd.Flags().StringVar(&Password, "password", "", "Mot de passe (Auth)")
	azureCmd.Flags().StringVar(&TenantID, "tenant", "common", "Tenant ID ou nom de domaine spécifique")
	rootCmd.AddCommand(azureCmd)
}
