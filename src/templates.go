package sposter

// I'm too lazy to do proper templating
var bskyPostTemplate = `
{{ .Title }} - {{ .PublishedDate }}

{{ .FirstSentence }}

{{ .Link }}
`
