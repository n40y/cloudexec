package templates

import (
	"regexp"
	"strings"
)

// Match vérifie si le contenu d'une réponse HTTP correspond aux critères du template
func (m *Matcher) Match(statusCode int, body []byte) bool {
	switch m.Type {
	case "status":
		return statusCode == m.Status

	case "word":
		// Vérifie si les mots attendus sont présents dans le corps de la réponse
		bodyStr := string(body)
		for _, value := range m.Values {
			if strings.Contains(bodyStr, value) {
				return true
			}
		}

	case "regex":
		// Vérifie si une expression régulière match le corps de la réponse
		bodyStr := string(body)
		for _, value := range m.Values {
			re, err := regexp.Compile(value)
			if err != nil {
				continue // Si la regex est mal écrite dans le YAML, on passe
			}
			if re.MatchString(bodyStr) {
				return true
			}
		}
	}

	return false
}
