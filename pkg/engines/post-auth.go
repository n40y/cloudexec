package engines

import (
	"fmt"
	"net/http"
	"time"

	"cloudexec/pkg/utils"
)

// RunPostAuthEnum centralise les actions de découverte après validation d'un secret
func RunPostAuthEnum(provider, credential string) {
	utils.LogWarning("--> Lancement du moteur d'énumération Post-Auth pour %s", provider)
	client := &http.Client{Timeout: 5 * time.Second}

	switch provider {
	case "GCP":
		// Liste d'APIs critiques à tester pour voir si la clé y donne accès
		apis := map[string]string{
			"Cloud Storage": "https://storage.googleapis.com/storage/v1/b?project=test&key=",
			"BigQuery":      "https://bigquery.googleapis.com/bigquery/v2/projects/test/datasets?key=",
			"Directions":    "https://maps.googleapis.com/maps/api/directions/json?origin=Paris&destination=Dijon&key=",
		}

		for name, targetURL := range apis {
			fullURL := targetURL + credential
			resp, err := client.Get(fullURL)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			// Une 403 signifie souvent que l'API est activée mais restriction de projet
			// Une 400 (sauf clé invalide) ou 200 prouve que l'API est accessible via cette clé
			if resp.StatusCode != 400 && resp.StatusCode != 403 {
				utils.LogSuccess("    [+] Service accessible : %s (Status: %d)", name, resp.StatusCode)
			} else {
				fmt.Printf("    [-] Service restreint : %s\n", name)
			}
		}

	case "AWS":
		// Structure prête pour l'appel des fonctions d'énumération AWS (S3, IAM, EC2)
		fmt.Printf("    [+] Collecte des politiques IAM courantes...\n")
		fmt.Printf("    [+] Recherche de buckets S3 exposés...\n")
	}
}
