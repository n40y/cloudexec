package engines

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"
)

func ExecuteGCPEngine(tmpl *templates.Template, domain, apiKey string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	// Nettoyage du domaine si c'est une URL
	if domain != "" {
		if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
			parsedURL, err := url.Parse(domain)
			if err == nil {
				domain = parsedURL.Host
			}
		} else {
			parts := strings.Split(domain, "/")
			domain = parts[0]
		}
	}

	switch tmpl.Action {
	case "gcp:WorkspaceCheck":
		// CORRECTION : Si pas de domaine, on passe le template sans erreur
		if domain == "" {
			return nil
		}
		utils.LogInfo("[%s] Vérification de la présence Google Workspace pour : %s", tmpl.ID, domain)

		urlStr := fmt.Sprintf("https://www.google.com/a/%s/ServiceNotAllowed?service=chia", domain)
		resp, err := client.Get(urlStr)
		if err != nil {
			return fmt.Errorf("erreur lors de la requête Google : %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 && !strings.Contains(resp.Request.URL.Path, "ServiceNotAllowed") {
			utils.LogSuccess("!!! MATCH TROUVÉ !!! [%s] -> %s", tmpl.ID, tmpl.Info.Name)
			fmt.Printf("    └─ Domaine : %s (Infrastructure Workspace Détectée)\n", domain)
		} else {
			utils.LogWarning("Aucun tenant Google Workspace public n'a été validé pour %s.", domain)
		}

	case "gcp:ApiKeyCheck":
		// CORRECTION : Si pas de clé, on passe le template sans erreur
		if apiKey == "" {
			return nil
		}
		utils.LogInfo("[%s] Validation de la clé API GCP : %s...", tmpl.ID, apiKey[:8])

		urlStr := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", apiKey)
		resp, err := client.Post(urlStr, "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			return fmt.Errorf("impossible de joindre l'API Google : %w", err)
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		bodyStr := buf.String()

		if resp.StatusCode == 200 || strings.Contains(bodyStr, "OPERATION_NOT_ALLOWED") {
			utils.LogSuccess("!!! CLÉ VALIDE !!! [%s] -> La clé API GCP est fonctionnelle !", tmpl.ID)

			// APPEL DU MOTEUR POST-AUTH
			RunPostAuthEnum("GCP", apiKey)

		} else if strings.Contains(bodyStr, "API_KEY_INVALID") || resp.StatusCode == 400 || resp.StatusCode == 404 {
			utils.LogError("Échec de la validation : Clé API GCP invalide ou révoquée.")
		} else {
			utils.LogWarning("Statut de réponse ambigu (%d). Clé potentiellement restreinte.", resp.StatusCode)
		}
	}

	return nil
}
