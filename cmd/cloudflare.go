package cmd

import (
	"fmt"
	"os"

	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"

	"github.com/spf13/cobra"
)

var cloudflareCmd = &cobra.Command{
	Use:   "cloudflare",
	Short: "Modules de détection et de bypass pour les protections Cloudflare",
	Run: func(cmd *cobra.Command, args []string) {
		if Target == "" {
			fmt.Println("[-] Erreur: Vous devez spécifier une cible avec -t ou --target")
			os.Exit(1)
		}

		folder := "templates/cloudflare"
		fmt.Printf("[*] Chargement des templates Cloudflare depuis : %s\n", folder)

		tmplList, err := templates.ParseDirectory(folder)
		if err != nil {
			fmt.Printf("[-] Erreur lors de la lecture du dossier : %v\n", err)
			return
		}

		fmt.Printf("[+] %d template(s) détecté(s). Lancement de l'analyse...\n", len(tmplList))

		for _, tmpl := range tmplList {
			if tmpl.Engine == "http" {
				// Utilisation du moteur HTTP avec le flag Verbose global
				err := engines.ExecuteHTTPEngine(tmpl, Target, Verbose)
				if err != nil {
					fmt.Printf("[-] Erreur d'exécution [%s] : %v\n", tmpl.ID, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cloudflareCmd)
}
