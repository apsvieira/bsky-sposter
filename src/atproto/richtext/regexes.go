package richtext

import "regexp"

const (
	MENTION_REGEX_STR                     = `(^|\s|\()(@)([a-zA-Z0-9.-]+)(\b)`
	URL_REGEX_STR                         = `(?i)(^|\s|\()((https?:\/\/[\S]+)|((?P<domain>[a-z][a-z0-9]*(\.[a-z0-9]+)+)[\S]*))`
	TRAILING_PUNCTUATION_REGEX_STR        = `\pP+$`
	TRAILING_PUNCTUATION_IN_URL_REGEX_STR = `[.,;:!?]$`

	// Index of the relevant match groups within the regex
	MENTION_HANDLER_MATCH_GROUP   = 3
	URL_URI_MATCH_GROUP           = 2
	URL_DOMAIN_CAPTURE_GROUP      = 5
	URL_DOMAIN_CAPTURE_GROUP_NAME = "domain"
)

var (
	mentionRegex                  = regexp.MustCompile(MENTION_REGEX_STR)
	urlRegex                      = regexp.MustCompile(URL_REGEX_STR)
	trailingPunctuationRegex      = regexp.MustCompile(TRAILING_PUNCTUATION_REGEX_STR)
	trailingPunctuationInUrlRegex = regexp.MustCompile(TRAILING_PUNCTUATION_IN_URL_REGEX_STR)
)
