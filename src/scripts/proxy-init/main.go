package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type ProxyConfig struct {
	To string `yaml:"to"`
}

func main() {
	helpers.EnsureAdminExecution()

	proxyConfigs := helpers.ReadConfig[map[string]ProxyConfig]("@/config/proxy.yml")
	proxyDir := helpers.ResolvePath("@/.local/proxy")

	helpers.ExecNativeCommand([]string{"caddy", "stop"})
	os.RemoveAll(proxyDir)
	os.MkdirAll(proxyDir, 0o755)

	helpers.ExecNativeCommand([]string{"mkcert", "-install"}, helpers.ExecCommandOptions{Exit: true})

	caddyfile := ""
	hosts := ""
	for domain, config := range proxyConfigs {
		fmt.Println(aurora.Blue(domain), aurora.Blue("->"), aurora.Faint(config.To))
		helpers.ExecNativeCommand([]string{"mkcert", domain}, helpers.ExecCommandOptions{Dir: proxyDir, Exit: true, Silent: true})
		caddyfile += fmt.Sprintf("%s {\n\ttls %s.pem %s-key.pem\n\treverse_proxy %s\n}\n\n", domain, domain, domain, config.To)
		hosts += "127.0.0.1 " + domain + "\n"
	}
	os.WriteFile(filepath.Join(proxyDir, "Caddyfile"), []byte(caddyfile), 0o644)

	hostsContent, _ := os.ReadFile("/etc/hosts")
	hostsBase, _, _ := strings.Cut(string(hostsContent), "# proxy-init\n")
	os.WriteFile("/etc/hosts", []byte(hostsBase+"# proxy-init\n"+hosts), 0o644)

	helpers.ExecNativeCommand([]string{"caddy", "start"}, helpers.ExecCommandOptions{Dir: proxyDir, Exit: true})
}
