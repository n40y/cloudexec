package cmd

import (
	"fmt"
	"os"

	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"

	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Lancer la phase de reconnaissance passive et OSINT",
	Run: func(cmd *cobra.Command, args []string) {
		// On récupère la cible définie globalement dans root.go
		if Target == "" {
			fmt.Println("[-] Erreur: Vous devez spécifier une cible avec -t ou --target")
			os.Exit(1)
		}

		folder := "templates/recon"
		fmt.Printf("[*] Chargement des templates de reconnaissance depuis : %s\n", folder)

		// 1. Charger tous les templates du dossier recon
		tmplList, err := templates.ParseDirectory(folder)
		if err != nil {
			fmt.Printf("[-] Erreur lors de la lecture du dossier : %v\n", err)
			return
		}

		fmt.Printf("[+] %d template(s) détecté(s). Lancement du scan...\n", len(tmplList))

		// 2. Boucler et exécuter chaque template
		for _, tmpl := range tmplList {
			if tmpl.Engine == "http" {
				err := engines.ExecuteHTTPEngine(tmpl, Target, Verbose)
				if err != nil {
					fmt.Printf("[-] Erreur d'exécution [%s] : %v\n", tmpl.ID, err)
				}
			}
		}
	},
}

func init() {
	// On attache la commande 'recon' à la commande racine 'rootCmd'
	rootCmd.AddCommand(reconCmd)
}
