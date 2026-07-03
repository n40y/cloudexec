package templates

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Info contient les métadonnées du template (Style Nuclei)
type Info struct {
	Name        string `yaml:"name"`
	Severity    string `yaml:"severity"`
	Description string `yaml:"description"`
}

// Matcher définit les conditions de succès du test
type Matcher struct {
	Type   string   `yaml:"type"`             // "status", "regex", "json"
	Status int      `yaml:"status,omitempty"` // Ex: 200
	Path   string   `yaml:"path,omitempty"`   // Pour le JSON Path
	Values []string `yaml:"values,omitempty"` // Chaînes ou Regex à chercher
}

// Request décrit la requête HTTP à envoyer
type Request struct {
	Method   string            `yaml:"method"`
	URL      string            `yaml:"url"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Matchers []Matcher         `yaml:"matchers"`
}

// Template est la structure racine d'un fichier de configuration YAML
type Template struct {
	ID       string    `yaml:"id"`
	Info     Info      `yaml:"info"`
	Engine   string    `yaml:"engine"`
	Action   string    `yaml:"action,omitempty"` // <-- Ajoute cette ligne ici
	Requests []Request `yaml:"requests"`
}

// ParseTemplate charge et dépose le contenu d'un fichier YAML dans la struct Template
func ParseTemplate(filePath string) (*Template, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tmpl Template
	err = yaml.Unmarshal(file, &tmpl)
	if err != nil {
		return nil, err
	}

	return &tmpl, nil
}

// ParseDirectory scanne un dossier et retourne tous les templates valides trouvés
func ParseDirectory(dirPath string) ([]*Template, error) {
	var list []*Template

	// Lit le contenu du dossier
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// On ne prend que les fichiers .yaml ou .yml
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			fullPath := filepath.Join(dirPath, file.Name())
			tmpl, err := ParseTemplate(fullPath)
			if err != nil {
				// Si un template a une erreur de syntaxe, on log et on continue sans crash
				continue
			}
			list = append(list, tmpl)
		}
	}

	return list, nil
}
