package bot

import "testing"

func TestSanitizeTelegramHTML(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain text passthrough",
			in:   "Hello world",
			want: "Hello world",
		},
		{
			name: "ampersand escaped",
			in:   "A & B",
			want: "A &amp; B",
		},
		{
			name: "script tag is escaped, not rendered",
			in:   `<script>alert(1)</script>`,
			want: `&lt;script&gt;alert(1)&lt;/script&gt;`,
		},
		{
			name: "supported simple tags pass through",
			in:   "<b>жирный</b> <i>курсив</i> <u>под</u> <s>зачёрк</s> <code>x</code> <tg-spoiler>?</tg-spoiler>",
			want: "<b>жирный</b> <i>курсив</i> <u>под</u> <s>зачёрк</s> <code>x</code> <tg-spoiler>?</tg-spoiler>",
		},
		{
			name: "blockquote and expandable",
			in:   "<blockquote>q</blockquote> <blockquote expandable>q2</blockquote>",
			want: "<blockquote>q</blockquote> <blockquote expandable>q2</blockquote>",
		},
		{
			name: "safe http link preserved",
			in:   `<a href="https://example.com/x?a=1&b=2">link</a>`,
			want: `<a href="https://example.com/x?a=1&amp;b=2">link</a>`,
		},
		{
			name: "javascript scheme rejected",
			in:   `<a href="javascript:alert(1)">x</a>`,
			want: `&lt;a href=&#34;javascript:alert(1)&#34;&gt;x&lt;/a&gt;`,
		},
		{
			name: "tag with extra attribute is escaped",
			in:   `<a href="https://e.com" onclick="x">y</a>`,
			want: `&lt;a href=&#34;https://e.com&#34; onclick=&#34;x&#34;&gt;y&lt;/a&gt;`,
		},
		{
			name: "pre with language code",
			in:   `<pre><code class="language-go">x</code></pre>`,
			want: `<pre><code class="language-go">x</code></pre>`,
		},
		{
			name: "lone < is escaped",
			in:   "1 < 2 and 3 > 1",
			want: "1 &lt; 2 and 3 &gt; 1",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sanitizeTelegramHTML(c.in); got != c.want {
				t.Errorf("sanitizeTelegramHTML(%q):\n got:  %q\n want: %q", c.in, got, c.want)
			}
		})
	}
}
