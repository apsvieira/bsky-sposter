package richtext_test

import (
	"context"
	"testing"

	"github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/apsvieira/bsky-sposter/src/atproto/richtext"
	"github.com/stretchr/testify/assert"
)

type tc struct {
	input  string
	output [][]string
}

var tcs = []tc{
	// mentions
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
	// links
	{
		"start https://middle.com end",
		[][]string{{"start "}, {"https://middle.com", "https://middle.com"}, {" end"}},
	},
	{
		"start https://middle.com/foo/bar end",
		[][]string{{"start "}, {"https://middle.com/foo/bar", "https://middle.com/foo/bar"}, {" end"}},
	},
	{
		"start https://middle.com/foo/bar?baz=bux end",
		[][]string{{"start "}, {"https://middle.com/foo/bar?baz=bux", "https://middle.com/foo/bar?baz=bux"}, {" end"}},
	},
	{
		"start https://middle.com/foo/bar?baz=bux#hash end",
		[][]string{{"start "}, {"https://middle.com/foo/bar?baz=bux#hash", "https://middle.com/foo/bar?baz=bux#hash"}, {" end"}},
	},
	{

		"https://start.com/foo/bar?baz=bux#hash middle end",
		[][]string{{"https://start.com/foo/bar?baz=bux#hash", "https://start.com/foo/bar?baz=bux#hash"}, {" middle end"}},
	},
	{
		"start middle https://end.com/foo/bar?baz=bux#hash",
		[][]string{{"start middle "}, {"https://end.com/foo/bar?baz=bux#hash", "https://end.com/foo/bar?baz=bux#hash"}},
	},
	{
		"https://newline1.com\nhttps://newline2.com",
		[][]string{
			{"https://newline1.com", "https://newline1.com"},
			{"\n"},
			{"https://newline2.com", "https://newline2.com"},
		},
	},
	{
		"start middle.com end",
		[][]string{{"start "}, {"middle.com", "https://middle.com"}, {" end"}},
	},
	{
		"start middle.com/foo/bar end",
		[][]string{{"start "}, {"middle.com/foo/bar", "https://middle.com/foo/bar"}, {" end"}},
	},
	{
		"start middle.com/foo/bar?baz=bux end",
		[][]string{{"start "}, {"middle.com/foo/bar?baz=bux", "https://middle.com/foo/bar?baz=bux"}, {" end"}},
	},
	{
		"start middle.com/foo/bar?baz=bux#hash end",
		[][]string{{"start "}, {"middle.com/foo/bar?baz=bux#hash", "https://middle.com/foo/bar?baz=bux#hash"}, {" end"}},
	},
	{
		"start.com/foo/bar?baz=bux#hash middle end",
		[][]string{{"start.com/foo/bar?baz=bux#hash", "https://start.com/foo/bar?baz=bux#hash"}, {" middle end"}},
	},
	{
		"start middle end.com/foo/bar?baz=bux#hash",
		[][]string{{"start middle "}, {"end.com/foo/bar?baz=bux#hash", "https://end.com/foo/bar?baz=bux#hash"}},
	},
	{
		"newline1.com\nnewline2.com",
		[][]string{
			{"newline1.com", "https://newline1.com"},
			{"\n"},
			{"newline2.com", "https://newline2.com"},
		},
	},
	{
		"a example.com/index.php php link",
		[][]string{{"a "}, {"example.com/index.php", "https://example.com/index.php"}, {" php link"}},
	},
	{
		"a trailing bsky.app: colon",
		[][]string{{"a trailing "}, {"bsky.app", "https://bsky.app"}, {": colon"}},
	},
	// crappy links
	{
		"not.. a..url ..here",
		[][]string{{"not.. a..url ..here"}},
	},
	{
		"e.g.",
		[][]string{{"e.g."}},
	},
	{
		"something-cool.jpg",
		[][]string{{"something-cool.jpg"}},
	},
	{
		"website.com.jpg",
		[][]string{{"website.com.jpg"}},
	},
	{
		"e.g./foo",
		[][]string{{"e.g./foo"}},
	},
	{
		"website.com.jpg/foo",
		[][]string{{"website.com.jpg/foo"}},
	},
	// more complex links
	{
		"Classic article https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/",
		[][]string{
			{"Classic article "},
			{"https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/", "https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/"},
		},
	},
	{
		"Classic article https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/ ",
		[][]string{{"Classic article "}, {"https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/", "https://socket3.wordpress.com/2018/02/03/designing-windows-95s-user-interface/"}, {" "}},
	},
	{
		"https://foo.com https://bar.com/whatever https://baz.com",
		[][]string{{"https://foo.com", "https://foo.com"}, {" "}, {"https://bar.com/whatever", "https://bar.com/whatever"}, {" "}, {"https://baz.com", "https://baz.com"}},
	},
	{
		"punctuation https://foo.com, https://bar.com/whatever; https://baz.com.",
		[][]string{{"punctuation "}, {"https://foo.com", "https://foo.com"}, {", "}, {"https://bar.com/whatever", "https://bar.com/whatever"}, {"; "}, {"https://baz.com", "https://baz.com"}, {"."}},
	},
	{
		"parenthentical (https://foo.com)",
		[][]string{{"parenthentical ("}, {"https://foo.com", "https://foo.com"}, {")"}},
	},
	{
		"except for https://foo.com/thing_(cool)",
		[][]string{{"except for "}, {"https://foo.com/thing_(cool)", "https://foo.com/thing_(cool)"}},
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

			for i, s := range rt.Segments() {
				output := segmentToOutput(s)
				assert.Equal(t, tc.output[i], output)
			}

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
