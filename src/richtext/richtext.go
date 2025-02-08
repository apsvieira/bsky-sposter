package richtext

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"golang.org/x/net/publicsuffix"
)

type RichTextProps struct {
	text   string
	facets []*bsky.RichtextFacet
}

type RichTextOpts struct {
	cleanNewLines bool
}

type RichTextSegment struct {
	text  string
	facet *bsky.RichtextFacet
}

func NewRichTextSegment(text string, facet *bsky.RichtextFacet) *RichTextSegment {
	return &RichTextSegment{text: text, facet: facet}
}

func (s *RichTextSegment) Link() *bsky.RichtextFacet_Link {
	for _, f := range s.facet.Features {
		if f.RichtextFacet_Link != nil {
			return f.RichtextFacet_Link
		}
	}
	return nil
}

func (s *RichTextSegment) Mention() *bsky.RichtextFacet_Mention {
	for _, f := range s.facet.Features {
		if f.RichtextFacet_Mention != nil {
			return f.RichtextFacet_Mention
		}
	}
	return nil
}

func (s *RichTextSegment) Tag() *bsky.RichtextFacet_Tag {
	for _, f := range s.facet.Features {
		if f.RichtextFacet_Tag != nil {
			return f.RichtextFacet_Tag
		}
	}
	return nil
}

type RichText struct {
	unicodeText string
	Facets      []*bsky.RichtextFacet
}

type ByIndexByteStart []*bsky.RichtextFacet

func (a ByIndexByteStart) Len() int           { return len(a) }
func (a ByIndexByteStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndexByteStart) Less(i, j int) bool { return a[i].Index.ByteStart < a[j].Index.ByteStart }

func NewRichText(text string) *RichText {
	return NewRichTextFromProps(RichTextProps{text: text}, &RichTextOpts{})
}
func NewRichTextFromProps(props RichTextProps, opts *RichTextOpts) *RichText {
	facets := props.facets
	if len(facets) > 0 {
		filtered := make([]*bsky.RichtextFacet, 0, len(facets))
		for _, f := range facets {
			if f.Index.ByteStart <= f.Index.ByteEnd {
				filtered = append(filtered, f)
			}
		}
		facets = filtered

		sort.Sort(ByIndexByteStart(facets))
	}

	text := props.text
	if opts.cleanNewLines {
		text = sanitizeText(text, opts)
	}

	return &RichText{unicodeText: text, Facets: facets}
}

func (rt *RichText) Text() string {
	return rt.unicodeText
}

func (rt *RichText) Length() int {
	return len(rt.unicodeText)
}

// GraphemeLength calculates the unicode grapheme length
func (rt *RichText) GraphemeLength() int {
	return len([]rune(rt.unicodeText))
}

// Segments returns the text segments with their respective facets
func (rt *RichText) Segments() []*RichTextSegment {
	segments := make([]*RichTextSegment, 0, len(rt.Facets))
	for _, f := range rt.Facets {
		text := rt.unicodeText[f.Index.ByteStart:f.Index.ByteEnd]
		segments = append(segments, NewRichTextSegment(text, f))
	}
	return segments
}

// DetectFacets detects facets such as links and mentions in the text.
// Note: Overwrites the existing facets with auto-detected facets.
func (rt *RichText) DetectFacets(ctx context.Context, agent *atproto.Client) error {
	facets := detectFacets(rt.unicodeText)
	if len(facets) == 0 {
		return nil
	}

	for _, facet := range facets {
		for _, feature := range facet.Features {
			if feature.RichtextFacet_Mention == nil {
				continue
			}

			data, err := agent.Com.Atproto.Identity.ResolveHandle(ctx, feature.RichtextFacet_Mention.Did)
			if err != nil {
				log.Printf("Error resolving handle %s: %s", feature.RichtextFacet_Mention.Did, err)
				continue
			}
			feature.RichtextFacet_Mention.Did = data.Did
		}
	}
	rt.Facets = facets
	return nil
}

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
