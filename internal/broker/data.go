package broker

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/ubuntu/oidc-broker/internal/providers/group"
)

// userInfo represents the user information that is returned to authd.
type userInfo struct {
	Name   string
	UUID   string
	Home   string
	Shell  string
	Gecos  string
	Groups []group.Info
}

// MarshalJSON implements the json.Marshaler interface for the userInfo type.
func (u userInfo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := template.Must(template.New("").Parse(`{"userinfo": {
		"name": "{{.Name}}",
		"uuid": "{{.UUID}}",
		"gecos": "{{.Gecos}}",
		"dir": "{{.Home}}",
		"shell": "{{.Shell}}",
		"groups": [ {{range $index, $g := .Groups}}
			{{- if $index}}, {{end -}}
			{"name": "{{.Name}}", "ugid": "{{.UGID}}"}
		{{- end}} ]
}}`)).Execute(&buf, u)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// errorMessage represents the error message that is returned to authd.
type errorMessage struct {
	Message string
}

// MarshalJSON implements the json.Marshaler interface for the errorMessage type.
func (msg errorMessage) MarshalJSON() ([]byte, error) {
	if msg.Message == "" {
		return nil, nil
	}

	b, err := json.Marshal(msg.Message)
	if err != nil {
		return nil, err
	}
	return []byte(`{"message": ` + string(b) + `}`), nil
}
