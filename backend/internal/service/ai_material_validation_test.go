package service

import (
	"strings"
	"testing"

	"ithozyeva/internal/models"
)

func TestAIMaterial_validateAndNormalize_Prompt_HappyPath(t *testing.T) {
	svc := &AIMaterialService{}
	req := &models.CreateAIMaterialRequest{
		Title:        "  Промт для генерации идей  ",
		Summary:      strings.Repeat("a", 35),
		ContentType:  "PROMPT",
		MaterialKind: "Prompt",
		PromptBody:   "ты помощник по генерации идей",
		ExternalURL:  "https://example.com/will-be-erased",
		AgentConfig:  "тоже стереть",
		Tags:         []string{"  Claude  ", "claude", "ideas", "", " GPT "},
	}

	out, tags, err := svc.validateAndNormalize(req)
	if err != nil {
		t.Fatalf("validateAndNormalize: %v", err)
	}
	if out.Title != "Промт для генерации идей" {
		t.Errorf("Title not trimmed: %q", out.Title)
	}
	if out.ContentType != models.AIMaterialContentTypePrompt {
		t.Errorf("ContentType = %q, want prompt", out.ContentType)
	}
	if out.MaterialKind != models.AIMaterialKindPrompt {
		t.Errorf("MaterialKind = %q, want prompt", out.MaterialKind)
	}
	if out.ExternalURL != "" || out.AgentConfig != "" {
		t.Errorf("non-prompt fields not cleared: url=%q agent=%q", out.ExternalURL, out.AgentConfig)
	}
	if got, want := tags, []string{"claude", "ideas", "gpt"}; !equalStrings(got, want) {
		t.Errorf("tags = %v, want %v (lowercase + dedup + trim)", got, want)
	}
}

func TestAIMaterial_validateAndNormalize_Link_HappyPath(t *testing.T) {
	svc := &AIMaterialService{}
	req := &models.CreateAIMaterialRequest{
		Title:        "Подборка промтов",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypeLink,
		MaterialKind: models.AIMaterialKindLibrary,
		ExternalURL:  "https://github.com/awesome/prompts",
		PromptBody:   "должно стереться",
	}
	out, _, err := svc.validateAndNormalize(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.PromptBody != "" {
		t.Errorf("PromptBody not cleared for link content")
	}
	if out.ExternalURL != "https://github.com/awesome/prompts" {
		t.Errorf("ExternalURL changed: %q", out.ExternalURL)
	}
}

func TestAIMaterial_validateAndNormalize_Link_RejectsNonHTTP(t *testing.T) {
	svc := &AIMaterialService{}
	cases := []string{"javascript:alert(1)", "ftp://example.com", "not-a-url", "//example.com", ""}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			_, _, err := svc.validateAndNormalize(&models.CreateAIMaterialRequest{
				Title:        "Заголовок",
				Summary:      strings.Repeat("a", 35),
				ContentType:  models.AIMaterialContentTypeLink,
				MaterialKind: models.AIMaterialKindLibrary,
				ExternalURL:  raw,
			})
			if err == nil {
				t.Errorf("expected error for url %q", raw)
			}
		})
	}
}

func TestAIMaterial_validateAndNormalize_Agent_RequiresConfig(t *testing.T) {
	svc := &AIMaterialService{}
	_, _, err := svc.validateAndNormalize(&models.CreateAIMaterialRequest{
		Title:        "Агент",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypeAgent,
		MaterialKind: models.AIMaterialKindAgent,
		AgentConfig:  "",
	})
	if err == nil {
		t.Fatal("expected error when agent config is empty")
	}
}

func TestAIMaterial_validateAndNormalize_Title_LengthBoundaries(t *testing.T) {
	svc := &AIMaterialService{}
	mk := func(title string) *models.CreateAIMaterialRequest {
		return &models.CreateAIMaterialRequest{
			Title:        title,
			Summary:      strings.Repeat("a", 35),
			ContentType:  models.AIMaterialContentTypePrompt,
			MaterialKind: models.AIMaterialKindPrompt,
			PromptBody:   "x",
		}
	}
	if _, _, err := svc.validateAndNormalize(mk("аб")); err == nil {
		t.Error("title=2 must fail (min 3)")
	}
	if _, _, err := svc.validateAndNormalize(mk(strings.Repeat("я", 121))); err == nil {
		t.Error("title=121 must fail (max 120)")
	}
	if _, _, err := svc.validateAndNormalize(mk("ну")); err == nil {
		t.Error("title=2 latin must fail")
	}
	if _, _, err := svc.validateAndNormalize(mk(strings.Repeat("я", 120))); err != nil {
		t.Errorf("title=120 must pass, got %v", err)
	}
}

func TestAIMaterial_validateAndNormalize_TagLimit(t *testing.T) {
	svc := &AIMaterialService{}
	tags := []string{"a", "b", "c", "d", "e", "f", "g"}
	_, normalized, err := svc.validateAndNormalize(&models.CreateAIMaterialRequest{
		Title:        "Заголовок",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
		Tags:         tags,
	})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(normalized) != models.AIMaterialMaxTags {
		t.Errorf("expected %d tags, got %d", models.AIMaterialMaxTags, len(normalized))
	}
}

func TestAIMaterial_validateAndNormalize_RejectsBadEnums(t *testing.T) {
	svc := &AIMaterialService{}
	_, _, err := svc.validateAndNormalize(&models.CreateAIMaterialRequest{
		Title:        "Заголовок",
		Summary:      strings.Repeat("a", 35),
		ContentType:  "video",
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	})
	if err == nil {
		t.Error("expected error for unknown content type")
	}
	_, _, err = svc.validateAndNormalize(&models.CreateAIMaterialRequest{
		Title:        "Заголовок",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: "course",
		PromptBody:   "x",
	})
	if err == nil {
		t.Error("expected error for unknown material kind")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
