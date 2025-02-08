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
	facets      []*bsky.RichtextFacet
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

	return &RichText{unicodeText: text, facets: facets}
}

func (rt *RichText) Facets() []*bsky.RichtextFacet {
	return rt.facets
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
	segments := make([]*RichTextSegment, 0, len(rt.facets))
	for _, f := range rt.facets {
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
	rt.facets = facets
	return nil
}
