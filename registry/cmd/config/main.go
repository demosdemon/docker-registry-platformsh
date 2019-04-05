package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"text/template"

	"github.com/elgs/gojq"
)

type jsonMap = map[string]interface{}

const configTemplate = `---
version: 0.1
log:
  accesslog:
    disabled: false
  level: {{ if eq (env "PLATFORM_BRANCH") "master" }}info{{ else }}debug{{ end }}
  formatter: text
  fields:
    branch: {{ env "PLATFORM_BRANCH" }}
    environment: {{ env "PLATFORM_ENVIRONMENT" }}
    project: {{ env "PLATFORM_PROJECT" }}
    service: {{ env "PLATFORM_APPLICATION_NAME" }}
    tree_id: {{ env "PLATFORM_TREE_ID" }}
storage:
  filesystem:
    rootdirectory: {{ env "PLATFORM_DIR" }}/var/lib/registry
auth:
  token:
    realm: {{ .TokenURL }}auth
    service: Docker Registry
    issuer: Acme Auth Server
    rootcertbundle: {{ env "PLATFORM_DIR" }}/bundle.crt
http:
  addr: localhost:{{ env "PORT" }}
  net: tcp
  prefix: /
  host: {{ .RegistryURL }}
  secret: {{ env "PLATFORM_PROJECT_ENTROPY" }}
redis:
  addr: {{ rel "cache.[0].host" }}:{{ rel "cache.[0].port" }}
  db: 0
{{ with rel "cache.[0].password" }}  password: {{ . }}{{ end }}
`

type config struct {
	TokenURL    string
	RegistryURL string
}

func getJSONVariable(key string) (*gojq.JQ, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}

	jsonValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return gojq.NewStringQuery(string(jsonValue))
}

func main() {
	application, _ := getJSONVariable("PLATFORM_APPLICATION")
	variables, _ := getJSONVariable("PLATFORM_VARIABLES")
	routes, _ := getJSONVariable("PLATFORM_ROUTES")
	rels, _ := getJSONVariable("PLATFORM_RELATIONSHIPS")

	funcMap := template.FuncMap{
		"env":      os.Getenv,
		"hostname": os.Hostname,
		"app":      application.Query,
		"var":      variables.Query,
		"route":    routes.Query,
		"rel":      rels.Query,
	}

	tpl := template.New("config.yml")
	tpl = tpl.Funcs(funcMap)
	tpl, err := tpl.Parse(configTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "E unable to compile template: %v\n", err)
		os.Exit(1)
	}

	cfg := new(config)
	if routeMap, ok := routes.Data.(jsonMap); ok {
		for url, v := range routeMap {
			if vMap, ok := v.(jsonMap); ok {
				if upstream, ok := vMap["upstream"]; ok {
					if s, ok := upstream.(string); ok {
						switch s {
						case "auth":
							cfg.TokenURL = url
						case "registry":
							cfg.RegistryURL = url
						}
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "failed to convert vMap %s %#v\n", url, v)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "failed to convert routeMap %#v\n", routes.Data)
	}

	err = tpl.Execute(os.Stdout, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "E unable to render template: %v\n", err)
		os.Exit(2)
	}
}
