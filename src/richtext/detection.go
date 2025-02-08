package richtext

import (
	"strings"

	"github.com/bluesky-social/indigo/api/bsky"
	"golang.org/x/net/publicsuffix"
)

func detectFacets(text string) []*bsky.RichtextFacet {
	var facets []*bsky.RichtextFacet

	// mentions
	for _, match := range mentionRegex.FindAllStringIndex(text, -1) {
		matchText := text[match[0]:match[1]]

		// this is a bit fancier in the original code
		start := strings.Index(matchText, "@") + 1
		mention := matchText[start:]
		if !isValidDomain(mention) && !strings.HasSuffix(mention, ".test") {
			continue
		}

		facet := &bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(start),
				ByteEnd:   int64(match[1]),
			},
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Mention: &bsky.RichtextFacet_Mention{
						Did: mention,
					},
				},
			},
		}
		facets = append(facets, facet)
	}

	// links
	for _, match := range urlRegex.FindAllStringIndex(text, -1) {
		uri := text[match[0]:match[1]]
		uri = strings.TrimSpace(uri)
		if !strings.HasPrefix(uri, "http") {
			if !isValidDomain(uri) {
				continue
			}
			uri = "https://" + uri
		}
		if trailingPunctuationRegex.MatchString(uri) {
			uri = uri[:len(uri)-1]
		}
		if strings.HasSuffix(uri, ")") && !strings.Contains(uri, "(") {
			uri = uri[:len(uri)-1]
		}

		facet := &bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(match[0]),
				ByteEnd:   int64(match[0] + len(uri) - 1), // TODO: check
			},
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Link: &bsky.RichtextFacet_Link{
						Uri: uri,
					},
				},
			},
		}
		facets = append(facets, facet)
	}

	// TODO: not handling tags as they depend on boring regexes
	// that I can't be bothered to translate to Go
	//
	// for _, match := range tagRegex.FindAllStringIndex(text, -1) {
	// 	tag := text[match[0]:match[1]]
	// 	tag = strings.TrimSpace(tag)
	// 	tag = trailingPunctuationRegex.ReplaceAllString(tag, "")

	// 	if len(tag) == 0 || len(tag) > 100 {
	// 		continue
	// 	}

	// 	facet := &bsky.RichtextFacet{
	// 		Index: &bsky.RichtextFacet_ByteSlice{
	// 			ByteStart: int64(match[0]),
	// 			ByteEnd:   int64(match[0] + len(tag) - 1), // TODO: check
	// 		},
	// 		Features: []*bsky.RichtextFacet_Features_Elem{
	// 			{
	// 				RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
	// 					Tag: tag,
	// 				},
	// 			},
	// 		},
	// 	}
	// 	facets = append(facets, facet)
	// }

	return facets
}

func isValidDomain(str string) bool {
	etld, im := publicsuffix.PublicSuffix(str)
	var validtld = false
	if im { // ICANN managed
		validtld = true
	} else if strings.IndexByte(etld, '.') >= 0 { // privately managed
		validtld = true
	}
	return validtld
}
