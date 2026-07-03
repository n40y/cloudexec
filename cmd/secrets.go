package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"cloudexec/pkg/engines"
	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	ScanPath   string
	OutputPath string
	Pivot      bool
)

type SecretResult struct {
	Type  string `json:"type"`
	File  string `json:"file"`
	Match string `json:"match"`
}

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Scanner de secrets et clés cloud concurrent avec mode pivot",
	Run: func(cmd *cobra.Command, args []string) {
		utils.LogInfo("Démarrage du scan de secrets dans : %s", ScanPath)
		if Pivot {
			utils.LogWarning("Mode PIVOT activé : Validation automatique des clés compromise en cours de scan.")
		}

		// Ajout de la regex pour chasser les clés API Google
		patterns := map[string]*regexp.Regexp{
			"AWS Access Key ID": regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
			"AWS Secret Key":    regexp.MustCompile(`[^A-Za-z0-9+/][A-Za-z0-9+/]{40}[^A-Za-z0-9+/]`),
			"GCP API Key":       regexp.MustCompile(`AIzaSy[A-Za-z0-9_\-]{32,35}`),
			"Slack Webhook":     regexp.MustCompile(`https://hooks\.slack\.com/services/T[A-Z0-9]+/_B[A-Z0-9]+/[A-Za-z0-9]+`),
		}

		filesChan := make(chan string, 100)
		resultsChan := make(chan SecretResult, 100)

		workersCount := runtime.NumCPU() * 2
		var workerWg sync.WaitGroup

		// 1. Pool de Workers
		for i := 0; i < workersCount; i++ {
			workerWg.Add(1)
			go func() {
				defer workerWg.Done()
				for path := range filesChan {
					content, err := os.ReadFile(path)
					if err != nil {
						continue
					}

					for label, regex := range patterns {
						matches := regex.FindAll(content, -1)
						for _, match := range matches {
							resultsChan <- SecretResult{
								Type:  label,
								File:  path,
								Match: string(match),
							}
						}
					}
				}
			}()
		}

		// 2. Collecteur centralisé + Logique de Pivot Automatique
		var finalResults []SecretResult
		doneCollecting := make(chan struct{})
		go func() {
			for res := range resultsChan {
				utils.LogWarning("SECRET TROUVÉ [%s] dans le fichier : %s", res.Type, res.File)
				fmt.Printf("    └─ Valeur détectée : %s\n", res.Match)
				finalResults = append(finalResults, res)

				// Si le mode pivot est actif et qu'on trouve une clé exploitable directement
				if Pivot {
					switch res.Type {
					case "GCP API Key":
						fmt.Printf(utils.Yellow + "    [PIVOT] Déclenchement automatique de la validation de la clé GCP...\n" + utils.Reset)

						// On forge un template à la volée pour réutiliser le moteur GCP existant
						pivotTmpl := &templates.Template{
							ID:     "gcp-dynamic-pivot",
							Engine: "gcp",
							Action: "gcp:ApiKeyCheck",
						}

						// On exécute directement la validation
						err := engines.ExecuteGCPEngine(pivotTmpl, "", res.Match)
						if err != nil {
							utils.LogError("Échec du pivot pour la clé %s : %v", res.Match[:8], err)
						}
					}
				}
			}
			close(doneCollecting)
		}()

		// 3. Scan du système de fichiers
		err := filepath.WalkDir(ScanPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == ".aws" {
					return filepath.SkipDir
				}
				return nil
			}

			ext := filepath.Ext(path)
			if ext == ".exe" || ext == ".bin" || ext == ".png" || ext == ".jpg" || ext == ".zip" || ext == ".gz" {
				return nil
			}

			filesChan <- path
			return nil
		})

		close(filesChan)
		workerWg.Wait()
		close(resultsChan)
		<-doneCollecting

		if err != nil {
			utils.LogError("Erreur pendant le scan : %v", err)
			return
		}

		utils.LogInfo("Scan terminé. Nombre de secrets découverts : %d", len(finalResults))

		// 4. Export JSON
		if OutputPath != "" {
			file, err := os.Create(OutputPath)
			if err != nil {
				utils.LogError("Impossible de générer le fichier de rapport : %v", err)
				return
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(finalResults); err != nil {
				utils.LogError("Erreur lors de la sérialisation du JSON : %v", err)
				return
			}
			utils.LogSuccess("Rapport d'audit sauvegardé dans : %s", OutputPath)
		}
	},
}

func init() {
	secretsCmd.Flags().StringVar(&ScanPath, "path", ".", "Chemin du dossier à analyser")
	secretsCmd.Flags().StringVarP(&OutputPath, "output", "o", "", "Fichier de sortie JSON")
	// Ajout du flag de Pivot (-p / --pivot)
	secretsCmd.Flags().BoolVarP(&Pivot, "pivot", "p", false, "Activer le mode pivot automatique pour valider les clés trouvées")
	rootCmd.AddCommand(secretsCmd)
}
