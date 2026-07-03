package engines

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloudexec/pkg/templates"
	"cloudexec/pkg/utils"
)

type TenantRealmResponse struct {
	State         int    `json:"State"`
	NameSpaceType string `json:"NameSpaceType"`
	DomainName    string `json:"DomainName"`
}

type OpenIDConfig struct {
	TokenEndpoint string `json:"token_endpoint"`
}

func ExecuteAzureEngine(tmpl *templates.Template, domain, username, password, tenant string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	switch tmpl.Action {
	case "azure:TenantEnum":
		if domain == "" {
			return fmt.Errorf("le flag --domain est requis pour l'action TenantEnum")
		}
		utils.LogInfo("[%s] Énumération des informations du tenant pour : %s", tmpl.ID, domain)

		urlRealm := fmt.Sprintf("https://login.microsoftonline.com/getuserrealm.srf?login=guest@%s&json=1", domain)
		respRealm, err := client.Get(urlRealm)
		if err != nil {
			return err
		}
		defer respRealm.Body.Close()

		var realm TenantRealmResponse
		json.NewDecoder(respRealm.Body).Decode(&realm)

		if realm.State == 3 {
			utils.LogError("Le domaine %s n'a pas de Tenant Azure/M365 actif.", domain)
			return nil
		}

		urlOpenID := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", domain)
		respOpenID, err := client.Get(urlOpenID)
		tenantID := "Introuvable"
		if err == nil {
			defer respOpenID.Body.Close()
			var openIDCfg OpenIDConfig
			if json.NewDecoder(respOpenID.Body).Decode(&openIDCfg) == nil {
				parts := strings.Split(openIDCfg.TokenEndpoint, "/")
				if len(parts) > 3 {
					tenantID = parts[3]
				}
			}
		}

		utils.LogSuccess("!!! MATCH TROUVÉ !!! [%s] -> %s", tmpl.ID, tmpl.Info.Name)
		fmt.Printf("    └─ "+utils.Cyan+"Domaine :"+utils.Reset+" %s\n", realm.DomainName)
		fmt.Printf("    └─ "+utils.Green+"Tenant ID :"+utils.Reset+" %s\n", tenantID)
		fmt.Printf("    └─ "+utils.Cyan+"Infra :"+utils.Reset+" %s\n", realm.NameSpaceType)

	case "azure:AuthCheck":
		if username == "" || password == "" {
			// On passe silencieusement si l'utilisateur ne teste pas l'auth
			return nil
		}
		utils.LogInfo("[%s] Tentative de connexion pour : %s", tmpl.ID, username)

		endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant)

		data := url.Values{}
		data.Set("grant_type", "password")
		data.Set("client_id", "04b07795-8ddb-461a-bbee-02f9e1bf7b46") // Azure CLI Client ID
		data.Set("username", username)
		data.Set("password", password)
		data.Set("scope", "https://management.core.windows.net//.default")

		req, _ := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			utils.LogSuccess("!!! IDENTIFIANTS VALIDES !!! [%s] -> Accès Azure confirmé pour %s", tmpl.ID, username)
		} else {
			utils.LogError("Échec de l'authentification (Status: %d) pour %s", resp.StatusCode, username)
		}
	}

	return nil
}
