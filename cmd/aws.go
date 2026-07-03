package cmd

import (
	"fmt"

	"cloudexec/pkg/config"
	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

// Déclaration des variables pour stocker les credentials
var (
	AccessKey string
	SecretKey string
	Region    string
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Modules d'audit et d'exploitation pour les environnements AWS",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Message d'avertissement de sécurité avec le nouveau logger
		utils.LogWarning("AVERTISSEMENT DE SÉCURITÉ :")
		fmt.Println("    Manipuler des clés d'accès Cloud comporte des risques.")
		fmt.Println("    Assurez-vous que vos fichiers de configuration ou logs ne soient JAMAIS")
		fmt.Println("    exposés publiquement ou commités sur GitHub (vérifiez votre .gitignore).\n")

		// 2. Chargement du fichier de configuration local (config.yaml)
		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			utils.LogError("Erreur lors de la lecture de config.yaml : %v", err)
			return
		}

		// Si le fichier existe et que l'utilisateur n'a pas utilisé les flags, on applique la config
		if cfg != nil {
			if AccessKey == "" {
				AccessKey = cfg.AWS.AccessKey
			}
			if SecretKey == "" {
				SecretKey = cfg.AWS.SecretKey
			}
			if Region == "us-east-1" && cfg.AWS.Region != "" { // Reste par défaut us-east-1 sauf si spécifié dans le YAML
				Region = cfg.AWS.Region
			}
		}

		// 3. Traitement des templates
		folder := "templates/aws"
		utils.LogInfo("Chargement des templates AWS depuis : %s", folder)

		tmplList, err := templates.ParseDirectory(folder)
		if err != nil {
			utils.LogError("Erreur lors de la lecture du dossier : %v", err)
			return
		}

		utils.LogInfo("%d template(s) détecté(s).", len(tmplList))

		for _, tmpl := range tmplList {
			if tmpl.Engine == "aws" {
				err := engines.ExecuteAWSEngine(tmpl, AccessKey, SecretKey, Region)
				if err != nil {
					utils.LogError("Erreur d'exécution [%s] : %v", tmpl.ID, err)
				}
			}
		}
	},
}

func init() {
	// Ajout des flags spécifiques à la commande AWS
	awsCmd.Flags().StringVar(&AccessKey, "access-key", "", "AWS Access Key ID")
	awsCmd.Flags().StringVar(&SecretKey, "secret-key", "", "AWS Secret Access Key")
	awsCmd.Flags().StringVar(&Region, "region", "us-east-1", "Région AWS par défaut")

	rootCmd.AddCommand(awsCmd)
}
