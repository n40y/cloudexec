package engines

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloudexec/pkg/bypass"
	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"
)

// Les variables suivantes font écho aux flags globaux définis dans cmd/root.go
// Elles permettent de savoir si l'utilisateur demande un bypass DNS/WAF
var (
	GlobalDirectIP string
	GlobalAntiWaf  bool
)

// ExecuteHTTPEngine traite les requêtes HTTP génériques définies dans les templates
func ExecuteHTTPEngine(tmpl *templates.Template, domain string) error {
	utils.LogInfo("[%s] Exécution du template HTTP pour : %s", tmpl.ID, domain)

	// 1. Initialisation du client HTTP
	var client *http.Client

	// Si une IP directe a été spécifiée (Bypass DNS Pinning / Cloudflare)
	if GlobalDirectIP != "" {
		client = bypass.NewBypassClient(domain, GlobalDirectIP)
	} else {
		// Client standard par défaut
		client = &http.Client{
			Timeout: 7 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	// 2. Construction de l'URL cible (à adapter selon la logique de tes templates)
	// Si l'action ou le template fournit une route spécifique, elle est concaténée ici
	targetURL := domain
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("impossible de créer la requête : %w", err)
	}

	// 3. Application des techniques d'évasion (Anti-WAF / Headers Spoofing)
	if GlobalAntiWaf {
		bypass.InjectBypassHeaders(req)
	} else {
		// Par défaut, on met un User-Agent propre pour éviter le flag standard "Go-http-client"
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) cloudexec/0.1.0")
	}

	// 4. Envoi de la requête
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'appel HTTP : %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("impossible de lire la réponse : %w", err)
	}
	bodyStr := string(bodyBytes)

	// 5. Logique de matching sommaire basée sur le statut ou le contenu
	// À synchroniser avec ton pkg/templates/matcher.go si nécessaire
	if resp.StatusCode == 200 {
		utils.LogSuccess("[%s] Réponse valide (200 OK) reçue de l'infrastructure.", tmpl.ID)
		// Optionnel : Analyse de signatures spécifiques ici
		if strings.Contains(bodyStr, "X-Amz-Bucket-Region") {
			utils.LogWarning("    [!] Signature d'infrastructure cloud détectée dans le corps.")
		}
	} else {
		fmt.Printf("[-] [%s] Statut retourné : %d\n", tmpl.ID, resp.StatusCode)
	}

	return nil
}
