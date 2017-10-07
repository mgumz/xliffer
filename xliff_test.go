package main

import (
	"strings"
	"testing"
)

func TestReadXliff(t *testing.T) {

	raw1 := `<?xml version="1.0" encoding="UTF-8"?>
<xliff version="1.2" xmlns="urn:oasis:names:tc:xliff:document:1.2">
	<file original="" source-language="en" target-language="de">
		<body>
			<trans-unit id="a">
				<source lang="en">Hello you!</source>
				<target lang="de">Hallo Du!</target>
			</trans-unit>
		</body>
	</file>
</xliff>
	`

	raw2 := `<?xml version="1.0" encoding="UTF-8"?>
<xliff version="1.2" xmlns="urn:oasis:names:tc:xliff:document:1.2">
	<file original="" source-language="en" target-language="de">
		<body>
			<trans-unit id="a">
				<source lang="en">Hello you!</source>
				<target lang="de"><mrk>Hallo</mrk> <mrk>Du!</mrk></target>
			</trans-unit>
		</body>
	</file>
</xliff>
	`

	doc1, err := xliffFromReader(strings.NewReader(raw1))
	if err != nil {
		t.Errorf("%s", err)
	}

	if doc1.File[0].SourceLang != "en" {
		t.Errorf("expected source-lang 'en', got source-lang %q", doc1.File[0].SourceLang)
	}
	if doc1.File[0].TargetLang != "de" {
		t.Errorf("expected source-lang 'de', got source-lang %q", doc1.File[0].TargetLang)
	}

	doc2, err := xliffFromReader(strings.NewReader(raw2))
	if err != nil {
		t.Errorf("%s", err)
	}

	target1 := doc1.File[0].Body.TransUnit[0].Target.Inner
	target2 := doc2.File[0].Body.TransUnit[0].Target.Inner
	if target1 != target2 {
		t.Errorf("expected %q == %q", target1, target2)
	}

}
