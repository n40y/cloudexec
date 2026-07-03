# CloudExec

CloudExec est un outil d'infrastructure, d'audit et de post-exploitation Cloud modulaire et ultra-rapide écrit en Go. Inspiré de l'ergonomie de NetExec (CME) et de la flexibilité de Nuclei, il utilise un moteur de concurrence (Goroutines) et un système de signatures dynamiques en YAML pour identifier les chemins d'attaque réels dans les environnements Cloud.


## 🗺️ Tableau de Bord des Fonctionnalités

### 1. Phase : Reconnaissance (OSINT & Authentifiée)

    [ ] Énumération de domaines/sous-domaines : Intégration de dictionnaires et mutations (cloud_enum/sandcastle).

    [ ] Découverte d'assets non authentifiée : Scan de buckets S3, Azure Blobs et GCP Buckets ouverts.

    [ ] Reconnaissance authentifiée : Extraction rapide de l'inventaire des ressources via les SDK officiels (cloudlist).

    [ ] Cartographie IAM : Identification de l'identité courante (GetCallerIdentity) et énumération passive des droits.

### 2. Phase : Cloudflare & DNS Bypass

    [ ] Énumération DNS avancée : Recherche de sous-domaines configurés hors du proxy Cloudflare (cloudflare_enum).

    [ ] Historique DNS (Origin IP Leak) : Interrogation d'APIs tierces (SecurityTrails, ViewDNS) pour trouver l'IP historique réelle pré-Cloudflare (cloudUnflare).

    [ ] Vérification de certificat (Censys/Shodan) : Corrélation des adresses IP directes exposant le certificat TLS de la cible.

    [ ] Validation de Bypass : Module de scan direct sur les IP trouvées pour confirmer l'accès à l'origine.

### 3. Phase : Audit & Scanner (Logique Cloned/Custom)

    [ ] Moteur de Templates YAML : Parser capable de lire des playbooks d'audit au format YAML.

    [ ] Vérification de configurations critiques : Portage Go des règles d'intrusion majeures de CloudSploit (ex: Groupes de sécurité ouverts sur le monde, IMDSv1 actif).

    [ ] Analyse différentielle des politiques : Détection de ressources critiques sans restriction d'accès (Buckets, snapshots).

### 4. Phase : Exploitation & Pivot (Red Team)

    [ ] Escalade de privilèges IAM (AWS/Azure/GCP) : Simulation et exécution automatique de chemins d'escalade (ex: iam:PassRole, CreateAccessKey).

    [ ] Persistance automatisée : Injection de clés d'accès secondaires ou modification de Trust Policies.

    [ ] Exfiltration discrète : Automatisation du partage de snapshots de volumes ou de bases de données vers un compte tiers.


## 📂 Arborescence Cible du Projet

Voici la structure finale à maintenir pour garantir la modularité de l'outil:

cloudexec/
├── cmd/                      # Interface CLI (Cobra)
│   ├── root.go               # Configuration et flags globaux (--target, -v)
│   ├── recon.go              # Commande 'cloudexec recon'
│   ├── cloudflare.go         # Commande 'cloudexec cloudflare'
│   ├── aws.go                # Commande 'cloudexec aws'
│   ├── azure.go              # Commande 'cloudexec azure'
│   └── gcp.go                # Commande 'cloudexec gcp'
├── pkg/                      # Logique métier (Packages Go internes)
│   ├── templates/            # Moteur de templates YAML (Style Nuclei)
│   │   ├── parser.go         # Lecture et validation des fichiers YAML
│   │   └── matcher.go        # Logique de validation (Status, JSON path, Regex)
│   ├── engines/              # Moteurs d'exécution de requêtes
│   │   ├── http_engine.go    # Exécution des templates HTTP (OSINT/Bypass)
│   │   ├── aws_engine.go     # Traducteur de templates vers appels SDK AWS
│   │   └── azure_engine.go   # Traducteur de templates vers appels SDK Azure
│   ├── bypass/               # Logique pure de bypass DNS et Cloudflare
│   ├── recon/                # Logique d'énumération de buckets et d'assets
│   └── exploit/              # Logique d'abus d'API et escalade (Red Team)
├── templates/                # Base de connaissances YAML (Signatures)
│   ├── cloudflare/           # Templates de détection de fuites IP
│   │   └── dns-leak.yaml
│   ├── aws/                  # Regroupement des règles critiques style CloudSploit
│   │   ├── s3-public-read.yaml
│   │   └── iam-privesc-passrole.yaml
│   └── azure/
├── go.mod                    # Fichier de dépendances Go
├── main.go                   # Point d'entrée unique du binaire compilé
└── README.md                 # Le présent tableau de bord


## 🛠️ Syntaxe d'Exécution Visée (UX Style NetExec)

L'outil doit rester simple, direct et axé sur l'utilisation de modules (-m) :

```bash
# 1. Lancer l'OSINT globale sur un domaine
./cloudexec recon -t cible.com

# 2. Tenter de bypasser Cloudflare pour trouver l'IP d'origine
./cloudexec cloudflare -t cible.com -m dns-history

# 3. Scanner un environnement AWS avec des clés compromises via les templates YAML
./cloudexec aws --access-key ID --secret-key SECRET -m s3-public-read

# 4. Exploiter une vulnérabilité IAM spécifique pour élever ses privilèges
./cloudexec aws --access-key ID --secret-key SECRET -m iam-privesc --exploit
```

## 🚀 Technologies Clés & Inspirations

* Langage : Go (Concurrence native via Goroutines, binaire unique statique sans dépendances).

* CLI Engine : Cobra (github.com/spf13/cobra).

* Parser Config : YAML v3 (gopkg.in/yaml.v3).

* Sources logiques analysées : cloudlist, s3scanner, cloud_enum, cloudUnflare, cloudsploit.