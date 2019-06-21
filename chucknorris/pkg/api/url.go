package api

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
)

// URLParams holds the url data for creating a GET request
type URLParams struct {
	Host    string
	Version string
	Query   map[string][]string
}

// URLBuilder helps create a URL using templates
func URLBuilder(servicename, apiEndpoint string, urlParams URLParams) (string, error) {
	tmpl, err := template.New(servicename).Parse(apiEndpoint)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse the template")
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, urlParams)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute template")
	}

	return buf.String(), nil
}
