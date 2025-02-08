package richtext

import "regexp"

var (
	EXCESS_SPACE_RE = regexp.MustCompile(`[\r\n]([\s]*[\r\n]){2,}`) // Not handling \u00AD\u2060\u200D\u200C\u200B
	REPLACEMENT_STR = "\n\n"
)

func sanitizeText(text string, opts *RichTextOpts) string {
	if opts.cleanNewLines {
		text = clean(text, EXCESS_SPACE_RE, REPLACEMENT_STR)
	}
	return text
}

func clean(text string, re *regexp.Regexp, replacement string) string {
	return re.ReplaceAllString(text, replacement)
}
