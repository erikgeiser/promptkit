package confirmation

// TemplateArrow is a template where the current choice is indicated by an
// arrow.
const TemplateArrow = `
{{- Bold .Prompt -}}
{{ if .YesSelected -}}
	{{- print (Bold " ▸Yes ") " No" -}}
{{- else -}}
	{{- print "  Yes " (Bold "▸No") -}}
{{- end -}}
`

// ConfirmationTemplateArrow is the ConfirmationTempalte that matches TemplateArrow.
const ConfirmationTemplateArrow = `
{{- print .Prompt " " -}}
{{- if .FinalValue -}}
	{{- Foreground "32" "Yes" -}}
{{- else -}}
	{{- Foreground "32" "No" -}}
{{- end }}
`

// TemplateYN is a classic template with ja [yn] indicator where the current
// value is capitalized and bold.
const TemplateYN = `
{{- Bold .Prompt -}}
{{ if .YesSelected -}}
	{{- print " [" (Bold "Y") "/n]" -}}
{{- else -}}
{{- print " [y/" (Bold "N") "]" -}}
{{- end -}}
`

// ConfirmationTemplateYN is the ConfirmationTempalte that matches TemplateYN.
const ConfirmationTemplateYN = `
{{- Bold .Prompt -}}
{{ if .FinalValue -}}
	{{- print " [" (Foreground "32" (Bold "Y")) "/n]" -}}
{{- else -}}
	{{- print " [y/" (Foreground "32" (Bold "N")) "]" -}}
{{- end }}
`
