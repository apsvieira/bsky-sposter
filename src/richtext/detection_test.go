package richtext_test

import (
	"context"
	"log"
	"testing"

	"github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/apsvieira/bsky-sposter/src/richtext"
	"github.com/stretchr/testify/assert"
)

type tc struct {
	input  string
	output [][]string
}

var tcs = []tc{
	{"no mention", [][]string{{"no mention"}}},
	{
		"@handle.com middle end",
		[][]string{{"@handle.com", "did:fake:handle.com"}, {" middle end"}},
	},
	{
		"start @handle.com end",
		[][]string{{"start "}, {"@handle.com", "did:fake:handle.com"}, {" end"}},
	},
	{
		"start middle @handle.com",
		[][]string{{"start middle "}, {"@handle.com", "did:fake:handle.com"}},
	},
	{
		"@handle.com @handle.com @handle.com",
		[][]string{
			{"@handle.com", "did:fake:handle.com"},
			{" "},
			{"@handle.com", "did:fake:handle.com"},
			{" "},
			{"@handle.com", "did:fake:handle.com"},
		},
	},
	{
		"@full123-chars.test",
		[][]string{{"@full123-chars.test", "did:fake:full123-chars.test"}},
	},
	{
		"not@right",
		[][]string{{"not@right"}},
	},
	{
		"@handle.com!@#$chars",
		[][]string{{"@handle.com", "did:fake:handle.com"}, {"!@#$chars"}},
	},
	{
		"@handle.com\n@handle.com",
		[][]string{
			{"@handle.com", "did:fake:handle.com"},
			{"\n"},
			{"@handle.com", "did:fake:handle.com"},
		},
	},
	{
		"parenthetical (@handle.com)",
		[][]string{
			{"parenthetical ("},
			{"@handle.com", "did:fake:handle.com"},
			{")"},
		},
	},
}

func TestDetectFacets(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			tc := tc
			ctx := context.Background()
			rt := richtext.NewRichText(tc.input)
			client, err := atproto.NewMockClient(ctx, "https://test.url", &atproto.Credentials{})
			assert.Nil(t, err)

			err = rt.DetectFacets(ctx, client)
			assert.Nil(t, err)

			log.Printf("Input: %s", tc.input)
			log.Printf("Output: %#v", tc.output)
			log.Printf("Segments: %#v", rt.Segments())

			for i, s := range rt.Segments() {
				log.Printf("Segment: %#v", s)
				output := segmentToOutput(s)
				assert.Equal(t, tc.output[i], output)
			}
			log.Printf("Done")

		})
	}
}

func segmentToOutput(s *richtext.RichTextSegment) []string {
	if s.Facet == nil {
		return []string{s.Text}
	}

	elements := make([]string, len(s.Facet.Features)+1)
	elements[0] = s.Text
	for i, f := range s.Facet.Features {
		text := ""
		if f.RichtextFacet_Mention != nil {
			text = f.RichtextFacet_Mention.Did
		}
		if f.RichtextFacet_Link != nil {
			text = f.RichtextFacet_Link.Uri
		}
		elements[i+1] = text
	}

	return elements
}
