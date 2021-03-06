package document

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/segmentio/terraform-docs/internal/pkg/doc"
	"github.com/segmentio/terraform-docs/internal/pkg/print"
	"github.com/segmentio/terraform-docs/internal/pkg/print/markdown"
	"github.com/segmentio/terraform-docs/internal/pkg/settings"
)

type MarkdownDocument struct{}

func (printer MarkdownDocument) Postprocessing(buffer *bytes.Buffer) (string, error) {
	return markdown.Sanitize(buffer.String()), nil
}

func (printer MarkdownDocument) PrintSeparator(buffer *bytes.Buffer, settings settings.Settings) {
	buffer.WriteString("\n")
}

func getInputDefaultValue(input *doc.Input, settings settings.Settings) string {
	var result = "n/a"

	if input.HasDefault() {
		if settings.Has(print.WithAggregateTypeDefaults) && input.IsAggregateType() {
			result = printFencedCodeBlock(print.GetPrintableValue(input.Default, settings, true))
		} else {
			result = fmt.Sprintf("`%s`", print.GetPrintableValue(input.Default, settings, false))
		}
	}

	return result
}

func (printer MarkdownDocument) PrintComment(buffer *bytes.Buffer, comment string, settings settings.Settings) {
	buffer.WriteString("## Description\n\n")
	buffer.WriteString(fmt.Sprintf("%s\n", comment))
}

func printFencedCodeBlock(code string) string {
	var buffer bytes.Buffer
	buffer.WriteString("\n\n")
	buffer.WriteString("```json\n")
	buffer.WriteString(code)
	buffer.WriteString("\n")
	buffer.WriteString("```")
	return buffer.String()
}

func printInput(buffer *bytes.Buffer, input doc.Input, settings settings.Settings) {
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("### %s\n\n", strings.Replace(input.Name, "_", "\\_", -1)))
	buffer.WriteString(fmt.Sprintf("Description: %s\n\n", markdown.ConvertMultiLineText(input.Description)))
	buffer.WriteString(fmt.Sprintf("Type: `%s`\n", input.Type))

	// Don't print defaults for required inputs when we're already explicit about it being required
	if !(settings.Has(print.WithRequired) && input.IsRequired()) {
		buffer.WriteString(fmt.Sprintf("\nDefault: %s\n", getInputDefaultValue(&input, settings)))
	}
}

func (printer MarkdownDocument) PrintInputs(buffer *bytes.Buffer, inputs []doc.Input, settings settings.Settings) {
	if settings.Has(print.WithRequired) {
		buffer.WriteString("## Required Inputs\n\n")
		buffer.WriteString("The following input variables are required:\n")

		for _, input := range inputs {
			if input.IsRequired() {
				printInput(buffer, input, settings)
			}
		}

		buffer.WriteString("\n")
		buffer.WriteString("## Optional Inputs\n\n")
		buffer.WriteString("The following input variables are optional (have default values):\n")

		for _, input := range inputs {
			if !input.IsRequired() {
				printInput(buffer, input, settings)
			}
		}
	} else {
		buffer.WriteString("## Inputs\n\n")
		buffer.WriteString("The following input variables are supported:\n")

		for _, input := range inputs {
			printInput(buffer, input, settings)
		}
	}
}

func (printer MarkdownDocument) PrintOutputs(buffer *bytes.Buffer, outputs []doc.Output, settings settings.Settings) {
	buffer.WriteString("## Outputs\n\n")
	buffer.WriteString("The following outputs are exported:\n")

	for _, output := range outputs {
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("### %s\n\n", strings.Replace(output.Name, "_", "\\_", -1)))
		buffer.WriteString(fmt.Sprintf("Description: %s\n", markdown.ConvertMultiLineText(output.Description)))
	}
}

func (printer MarkdownDocument) PrintModules(buffer *bytes.Buffer, modules []doc.Module, settings settings.Settings) {
	buffer.WriteString("## Modules\n\n")
	buffer.WriteString("The following modules are referenced:\n")

	for _, module := range modules {
		if settings.Has(print.WithLinksToModules) {
			buffer.WriteString(fmt.Sprintf(" - [%s](%s/%s.md)", strings.Replace(module.Name, "_", "\\_", -1), module.Source, settings.Get(print.ModuleDocumentationFileName)))
		} else {
			buffer.WriteString(fmt.Sprintf(" - %s `%s`", strings.Replace(module.Name, "_", "\\_", -1), module.Source))
		}
		if module.HasDescription() {
			buffer.WriteString(fmt.Sprintf(" (*%s*)", markdown.ConvertMultiLineText(module.Description)))
		}
		buffer.WriteString("\n")
	}
}

func (printer MarkdownDocument) PrintResources(buffer *bytes.Buffer, resources []doc.Resource, settings settings.Settings) {
	buffer.WriteString("## Resources\n\n")
	buffer.WriteString("The following resources are created:\n")

	for _, resource := range resources {
		buffer.WriteString(fmt.Sprintf("- %s [`%s`](https://www.terraform.io/docs/providers/%s/r/%s.html) ", resource.Name, resource.Type, resource.Type.Provider(), resource.Type.Name()))
		if resource.HasDescription() {
			buffer.WriteString(fmt.Sprintf(" (*%s*)", markdown.ConvertMultiLineText(resource.Description)))
		}
		buffer.WriteString("\n")
	}
}
