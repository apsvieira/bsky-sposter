package richtext

import (
	"fmt"
	"strings"

	"github.com/bluesky-social/indigo/api/bsky"
	"golang.org/x/net/publicsuffix"
)

func detectFacets(text string) []*bsky.RichtextFacet {
	var facets []*bsky.RichtextFacet

	// mentions
	for _, match := range mentionRegex.FindAllStringSubmatchIndex(text, -1) {
		handle := text[match[MENTION_HANDLER_MATCH_GROUP*2]:match[MENTION_HANDLER_MATCH_GROUP*2+1]]
		if !isValidDomain(handle) && !strings.HasSuffix(handle, ".test") {
			continue
		}

		start := match[MENTION_HANDLER_MATCH_GROUP*2] - 1

		facet := &bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(start),
				ByteEnd:   int64(start + len(handle) + 1),
			},
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Mention: &bsky.RichtextFacet_Mention{
						Did: handle,
					},
				},
			},
		}
		facets = append(facets, facet)
	}

	// links
	for _, match := range urlRegex.FindAllStringSubmatchIndex(text, -1) {
		uri := text[match[URL_URI_MATCH_GROUP*2]:match[URL_URI_MATCH_GROUP*2+1]]
		matchLength := len(uri)
		if !strings.HasPrefix(uri, "http") {
			var domain string
			if urlRegex.SubexpIndex(URL_DOMAIN_CAPTURE_GROUP_NAME) != -1 {
				domain = text[match[URL_DOMAIN_CAPTURE_GROUP*2]:match[URL_DOMAIN_CAPTURE_GROUP*2+1]]
			}
			if len(domain) == 0 || !isValidDomain(domain) {
				continue
			}

			uri = fmt.Sprintf("https://%s", uri)
		}

		start := match[URL_URI_MATCH_GROUP*2]
		index := []int{start, start + matchLength}
		// strip ending punctuation
		if trailingPunctuationInUrlRegex.MatchString(uri) {
			uri = uri[:len(uri)-1]
			index[1]--
		}
		if strings.HasSuffix(uri, ")") && !strings.Contains(uri, "(") {
			uri = uri[:len(uri)-1]
			index[1]--
		}

		facet := &bsky.RichtextFacet{
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(index[0]),
				ByteEnd:   int64(index[1]),
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
