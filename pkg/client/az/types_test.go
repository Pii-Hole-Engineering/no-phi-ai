package az

import "testing"

func TestNewDocument(t *testing.T) {
	// Test case 1: Valid inputs
	id := "123"
	text := "Hello, world!"
	language := "en"
	doc := NewDocument(id, text, language)
	if doc.Language != language {
		t.Errorf("Expected language to be %s, but got %s", language, doc.Language)
	}
	if doc.ID != id {
		t.Errorf("Expected ID to be %s, but got %s", id, doc.ID)
	}
	if doc.Text != text {
		t.Errorf("Expected Text to be %s, but got %s", text, doc.Text)
	}

	// Test case 2: Empty language
	id = "456"
	text = "Lorem ipsum"
	language = ""
	doc = NewDocument(id, text, language)
	if doc.Language != DefaultLanguage {
		t.Errorf("Expected language to be %s, but got %s", DefaultLanguage, doc.Language)
	}
	if doc.ID != id {
		t.Errorf("Expected ID to be %s, but got %s", id, doc.ID)
	}
	if doc.Text != text {
		t.Errorf("Expected Text to be %s, but got %s", text, doc.Text)
	}
}
