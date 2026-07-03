package utils

import (
	"fmt"
)

// Codes de couleur ANSI basiques et universels
const (
	Reset  = "\x1b[0m"
	Red    = "\x1b[31m"
	Green  = "\x1b[32m"
	Yellow = "\x1b[33m"
	Blue   = "\x1b[34m"
	Cyan   = "\x1b[36m"
	Bold   = "\x1b[1m"
)

// Banner affiche un superbe artwork au lancement de la CLI
func Banner() {
	banner := `
   ____ _                 _ _____                     
  / ___| | ___  _   _  __| | ____|_  ___  ___  ___    
 | |   | |/ _ \| | | |/ _` + "`" + `|  _| \ \/ / _ \/ __|/ _ \   
 | |___| | (_) | |_| | (_| | |___ >  <  __/ (__| (_) |  
  \____|_|\___/ \__,_|\__,_|_____/_/\_\___|\___|\___/   `

	fmt.Println(Cyan + Bold + banner + Reset)
	fmt.Println(Bold + "         Cloud Security & Audit Scanner | v0.1.0\n" + Reset)
}

// LogInfo pour les étapes de configuration ou de chargement
func LogInfo(format string, a ...interface{}) {
	fmt.Printf(Blue+"[*] "+Reset+format+"\n", a...)
}

// LogSuccess pour les vulnérabilités ou les configurations trouvées (Match)
func LogSuccess(format string, a ...interface{}) {
	fmt.Printf(Green+Bold+"[+] "+Reset+Bold+format+Reset+"\n", a...)
}

// LogError pour les échecs ou les exceptions
func LogError(format string, a ...interface{}) {
	fmt.Printf(Red+"[-] "+Reset+format+"\n", a...)
}

// LogWarning pour les alertes critiques (sécurité, etc.)
func LogWarning(format string, a ...interface{}) {
	fmt.Printf(Yellow+"[!] "+Reset+format+"\n", a...)
}
