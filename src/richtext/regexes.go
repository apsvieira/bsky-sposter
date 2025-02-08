package richtext

import "regexp"

const (
	MENTION_REGEX_STR              = `(^|\s|\()(@)([a-zA-Z0-9.-]+)(\b)`
	URL_REGEX_STR                  = `(?i)(^|\s|\()((https?:\/\/\S+)|(([a-z][a-z0-9]*(\.[a-z0-9]+)+)\S*))`
	TRAILING_PUNCTUATION_REGEX_STR = `\pP+$`
	// TAG_REGEX_STR                  = `(^|\s)[#＃]([^\s\u00AD\u2060\u200A\u200B\u200C\u200D\u20e2]*[^\d\s\pP\u00AD\u2060\u200A\u200B\u200C\u200D\u20e2]+[^\s\u00AD\u2060\u200A\u200B\u200C\u200D\u20e2]*)?`
)

var (
	mentionRegex             = regexp.MustCompile(MENTION_REGEX_STR)
	urlRegex                 = regexp.MustCompile(URL_REGEX_STR)
	trailingPunctuationRegex = regexp.MustCompile(TRAILING_PUNCTUATION_REGEX_STR)
	// tagRegex                 = regexp.MustCompile(TAG_REGEX_STR)
)
