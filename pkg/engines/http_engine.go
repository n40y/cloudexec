package engines

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"
)

func ExecuteHTTPEngine(tmpl *templates.Template, target string, verbose bool) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for _, req := range tmpl.Requests {
		finalURL := strings.ReplaceAll(req.URL, "{{target}}", target)

		httpReq, err := http.NewRequest(req.Method, finalURL, nil)
		if err != nil {
			return fmt.Errorf("impossible de créer la requête: %w", err)
		}

		for key, value := range req.Headers {
			httpReq.Header.Set(key, value)
		}

		fmt.Printf("[*] [%s] Envoi de la requête vers %s\n", tmpl.ID, finalURL)

		resp, err := client.Do(httpReq)
		if err != nil {
			return fmt.Errorf("erreur lors de l'envoi HTTP: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("impossible de lire le corps de la réponse: %w", err)
		}

		// Évaluation des Matchers
		for _, matcher := range req.Matchers {
			if matcher.Match(resp.StatusCode, body) {
				// Surlignage du Status Code en Cyan/Gras dans le résultat
				utils.LogSuccess("!!! MATCH TROUVÉ !!! [%s] -> %s ("+utils.Cyan+"Status: %d"+utils.Reset+")",
					tmpl.ID, tmpl.Info.Name, resp.StatusCode)
			} else {
				if verbose {
					utils.LogError("Pas de match pour le critère de type '%s' (Status: %d)", matcher.Type, resp.StatusCode)
				}
			}
		}
	}
	return nil
}
