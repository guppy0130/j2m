package j2m_test

import (
	_ "embed"
	"testing"

	"github.com/guppy0130/j2m"
)

//go:embed j2m.jira
var testJira string

//go:embed j2m.md
var testMarkdown string

func TestJiraToMD(t *testing.T) {

	testcases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "should convert bolds properly",
			input: "*bold*",
			want:  "**bold**",
		},
		{
			name:  "should convert italics properly",
			input: "_italic_",
			want:  "*italic*",
		},
		{
			name:  "should convert monospaced content properly",
			input: "{{monospaced}}",
			want:  "`monospaced`",
		},
		// { name: "should convert citations properly",
		//     input: "??citation??",
		//     want: "<cite>citation</cite>" },
		{
			name:  "should convert strikethroughs properly",
			input: " -deleted- ",
			want:  " ~~deleted~~ ",
		},
		{
			name:  "should convert inserts properly",
			input: "+inserted+",
			want:  "<ins>inserted</ins>",
		},
		{
			name:  "should convert superscript properly",
			input: "^superscript^",
			want:  "<sup>superscript</sup>",
		},
		{
			name:  "should convert subscript properly",
			input: "~subscript~",
			want:  "<sub>subscript</sub>",
		},
		{
			name:  "should convert preformatted blocks properly",
			input: "{noformat}\nso *no* further _formatting_ is done here\n{noformat}",
			want:  "```\nso **no** further *formatting* is done here\n```",
		},
		{
			name:  "should convert language-specific code blocks properly",
			input: "{code:javascript}\nconst hello = 'world';\n{code}",
			want:  "```javascript\nconst hello = 'world';\n```",
		},
		{
			name:  "should convert code without language-specific and with title into code block",
			input: "{code:title=Foo.java}\nclass Foo {\n  public static void main() {\n  }\n}\n{code}",
			want:  "```\nclass Foo {\n  public static void main() {\n  }\n}\n```",
		},
		{
			name:  "should convert code without line feed before the end code block",
			input: "{code:java}\njava code{code}",
			want:  "```java\njava code\n```",
		},
		{
			name:  "should convert code without line feeds",
			input: "{code:java}java code{code}",
			want:  "```java\njava code\n```",
		},
		{
			name: "should convert fully configured code block",
			input: "{code:xml|title=My Title|borderStyle=dashed|borderColor=#ccc|titleBGColor=#F7D6C1|bgColor=#FFFFCE}" +
				"\n    <test>" +
				"\n        <another tag=\"attribute\"/>" +
				"\n    </test>" +
				"\n{code}",
			want: "```xml\n    <test>\n        <another tag=\"attribute\"/>\n    </test>\n```",
		},
		{
			name:  "should convert images properly",
			input: "!http://google.com/image!",
			want:  "![](http://google.com/image)",
		},
		{
			name:  "should convert linked images properly",
			input: "[!http://google.com/image!|http://google.com/link]",
			want:  "[![](http://google.com/image)](http://google.com/link)",
		},
		{
			name:  "should convert unnamed links properly",
			input: "[http://google.com]",
			want:  "<http://google.com>",
		},
		{
			name:  "should convert named links properly",
			input: "[Google|http://google.com]",
			want:  "[Google](http://google.com)",
		},
		{
			name:  "should convert headers properly: h1",
			input: "h1. Biggest heading",
			want:  "# Biggest heading",
		},
		{
			name:  "should convert headers properly: h2",
			input: "h2. Bigger heading",
			want:  "## Bigger heading",
		},
		{
			name:  "should convert headers properly: h3",
			input: "h3. Big heading",
			want:  "### Big heading",
		},
		{
			name:  "should convert headers properly: h4",
			input: "h4. Normal heading",
			want:  "#### Normal heading",
		},
		{
			name:  "should convert headers properly: h5",
			input: "h5. Small heading",
			want:  "##### Small heading",
		},
		{
			name:  "should convert headers properly: h6",
			input: "h6. Smallest heading",
			want:  "###### Smallest heading",
		},
		{
			name:  "should convert blockquotes properly",
			input: "bq. This is a long blockquote type thingy that needs to be converted.",
			want:  "> This is a long blockquote type thingy that needs to be converted.",
		},
		{
			name:  "should convert un-ordered lists properly",
			input: "* Foo\n* Bar\n* Baz\n** FooBar\n** BarBaz\n*** FooBarBaz\n* Starting Over",
			want:  "* Foo\n* Bar\n* Baz\n  * FooBar\n  * BarBaz\n    * FooBarBaz\n* Starting Over",
		},
		{
			name:  "should convert ordered lists properly",
			input: "# Foo\n# Bar\n# Baz\n## FooBar\n## BarBaz\n### FooBarBaz\n# Starting Over",
			want:  "1. Foo\n1. Bar\n1. Baz\n   1. FooBar\n   1. BarBaz\n      1. FooBarBaz\n1. Starting Over",
		},
		{
			name:  "should handle bold AND italic (combined) correctly",
			input: "This is _*emphatically bold*_!",
			want:  "This is ***emphatically bold***!",
		},
		{
			name:  "should handle bold within a un-ordered list item",
			input: "* This is not bold!\n** This is *bold*.",
			want:  "* This is not bold!\n  * This is **bold**.",
		},
		// {
		// 	name:  "should be able to handle a complicated multi-line jira-wiki string and convert it to markdown",
		// 	input: testJira,
		// 	want:  testMarkdown,
		// },
		{
			name:  "should not recognize strikethroughs over multiple lines",
			input: "* Here's an un-ordered list line\n* Multi-line strikethroughs shouldn't work.",
			want:  "* Here's an un-ordered list line\n* Multi-line strikethroughs shouldn't work.",
		},
		{
			name:  "should remove color attributes",
			input: "A text with{color:blue} blue \n lines {color} is not necessary.",
			want:  "A text with blue \n lines  is not necessary.",
		},
		// { name: "should not recognize inserts across multiple table cells",
		//      input: "||Heading 1||Heading 2||\n|Col+A1|Col+A2|",
		//      want: "\n|Heading 1|Heading 2|\n| --- | --- |\n|Col+A1|Col+A2|" },

	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			got := j2m.JiraToMD(tc.input)

			if tc.want != got {
				t.Errorf("%v: got %v, expected %v", tc.name, got, tc.want)
			}
		})
	}
}
