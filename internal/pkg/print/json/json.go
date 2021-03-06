package json

import (
	"encoding/json"

	"github.com/segmentio/terraform-docs/internal/pkg/doc"
	"github.com/segmentio/terraform-docs/internal/pkg/print"
	"github.com/segmentio/terraform-docs/internal/pkg/settings"
)

const (
	indent string = "  "
	prefix string = ""
)

type Printer struct{}

// Print prints a document as json.
func (printer Printer) Print(document *doc.Doc, settings settings.Settings) (string, error) {
	if document.HasInputs() {
		if settings.Has(print.WithSortByName) {
			if settings.Has(print.WithSortInputsByRequired) {
				doc.SortInputsByRequired(document.Inputs)
			} else {
				doc.SortInputsByName(document.Inputs)
			}
		}
	}

	if document.HasModules() {
		if settings.Has(print.WithSortByName) {
			doc.SortModulesByName(document.Modules)
		}
	}

	if document.HasResources() {
		if settings.Has(print.WithSortByName) {
			doc.SortResourcesByName(document.Resources)
		}
	}

	if document.HasOutputs() {
		if settings.Has(print.WithSortByName) {
			doc.SortOutputsByName(document.Outputs)
		}
	}

	buffer, err := json.MarshalIndent(document, prefix, indent)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}
