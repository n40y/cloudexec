package bypass

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// Liste de User-Agents légitimes pour contourner les règles anti-bot basiques
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
}

// BypassTransport configure un client HTTP avec des mécanismes d'évasion
type BypassTransport struct {
	SourceIP string // Pour du Spoofing de header si nécessaire
}

// NewBypassClient crée un client HTTP qui bypass les résolutions DNS classiques (DNS Pinning)
// et ignore les certificats invalides (fréquent sur les IPs directes d'infrastructures cloud)
func NewBypassClient(targetDomain string, directIP string) *http.Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Bypass les alertes TLS sur IP directe
	}

	// Si une IP directe est fournie, on force l'aiguillage sans consulter les serveurs DNS
	if directIP != "" {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			_, port, _ := net.SplitHostPort(addr)
			// On force la connexion sur l'IP choisie, tout en gardant le port d'origine (80/443)
			return dialer.DialContext(ctx, network, net.JoinHostPort(directIP, port))
		}
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}

// InjectBypassHeaders ajoute les entêtes d'évasion à la requête
func InjectBypassHeaders(req *http.Request) {
	// 1. Randomisation du User-Agent
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	// 2. Headers de spoofing d'origine (trompe certains WAFs mal configurés)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	req.Header.Set("X-Originating-IP", "127.0.0.1")
	req.Header.Set("X-Real-IP", "127.0.0.1")
}
