package pkg

import "testing"

func TestExtensionHasVersion(t *testing.T) {
	shouldBeFalse := extensionHasVersion("graphql.vscodegraphql")

	shouldBeTrue := extensionHasVersion("graphql.vscode-graphql-0.3.10")

	if shouldBeFalse != false {
		t.Errorf("shouldBeFalse not what it should be")
	}

	if shouldBeTrue != true {
		t.Errorf("shouldBeTrue not what it should be\n")
	}
}
