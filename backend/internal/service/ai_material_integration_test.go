package service

import (
	"strings"
	"testing"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
	"ithozyeva/internal/testutil"
)

func TestAIMaterialService_CreateAndSearch_ByTagAndKind(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_likes",
		"ai_material_bookmarks",
		"ai_material_comments",
		"ai_material_tags",
		"ai_materials",
		"members",
	)

	author := seedMember(t, db, 9001)
	svc := NewAIMaterialService()

	prompt, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Промт для код-ревью",
		Summary:      strings.Repeat("Помогает делать код-ревью", 2),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "ты ревьюер...",
		Tags:         []string{"Claude", "Review"},
	}, author.Id)
	if err != nil {
		t.Fatalf("create prompt: %v", err)
	}
	if prompt.Id == 0 {
		t.Fatal("created material has zero id")
	}
	if got := prompt.Tags; !equalStrings(got, []string{"claude", "review"}) {
		t.Errorf("tags = %v", got)
	}

	link, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Большая подборка промтов",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypeLink,
		MaterialKind: models.AIMaterialKindLibrary,
		ExternalURL:  "https://github.com/awesome/prompts",
		Tags:         []string{"library", "claude"},
	}, author.Id)
	if err != nil {
		t.Fatalf("create link: %v", err)
	}

	// Фильтр по kind=library должен вернуть только link
	items, total, err := svc.Search(repository.AIMaterialFilter{Kind: "library"})
	if err != nil {
		t.Fatalf("search by kind: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Id != link.Id {
		t.Errorf("kind=library: got %d items (total=%d), want 1 with id=%d", len(items), total, link.Id)
	}

	// Фильтр по тегу claude — оба
	_, totalByTag, err := svc.Search(repository.AIMaterialFilter{Tag: "claude"})
	if err != nil {
		t.Fatalf("search by tag: %v", err)
	}
	if totalByTag != 2 {
		t.Errorf("tag=claude: total=%d, want 2", totalByTag)
	}

	// Фильтр по q=код — только prompt
	_, totalByQ, err := svc.Search(repository.AIMaterialFilter{Query: "код"})
	if err != nil {
		t.Fatalf("search by q: %v", err)
	}
	if totalByQ != 1 {
		t.Errorf("q=код: total=%d, want 1", totalByQ)
	}
}

func TestAIMaterialService_Update_OnlyAuthor(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9101)
	stranger := seedMember(t, db, 9102)
	svc := NewAIMaterialService()

	created, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Title",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Чужой пользователь не может обновить
	_, err = svc.Update(created.Id, &models.UpdateAIMaterialRequest{
		Title:        "Hacked",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, stranger.Id, false)
	if err == nil {
		t.Fatal("stranger should not be able to update")
	}

	// Автор может
	updated, err := svc.Update(created.Id, &models.UpdateAIMaterialRequest{
		Title:        "Author updated",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
		Tags:         []string{"new"},
	}, author.Id, false)
	if err != nil {
		t.Fatalf("author update: %v", err)
	}
	if updated.Title != "Author updated" {
		t.Errorf("Title not updated: %q", updated.Title)
	}
	if !equalStrings(updated.Tags, []string{"new"}) {
		t.Errorf("Tags not replaced: %v", updated.Tags)
	}

	// Админ может (тут симулируем флагом isAdmin=true)
	if _, err := svc.Update(created.Id, &models.UpdateAIMaterialRequest{
		Title:        "Admin updated",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, stranger.Id, true); err != nil {
		t.Errorf("admin should be able to update: %v", err)
	}
}

func TestAIMaterialService_LikesTrigger_UpdatesCount(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_likes", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9201)
	liker := seedMember(t, db, 9202)
	svc := NewAIMaterialService()

	m, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Title",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := db.Create(&models.AIMaterialLike{MaterialId: m.Id, MemberId: liker.Id}).Error; err != nil {
		t.Fatalf("insert like: %v", err)
	}

	got, err := svc.GetByID(m.Id, liker.Id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.LikesCount != 1 {
		t.Errorf("LikesCount = %d, want 1 (trigger should have fired)", got.LikesCount)
	}
	if !got.Liked {
		t.Error("Liked must be true for liker viewer")
	}

	// удаление лайка декрементит
	if err := db.Where("material_id = ? AND member_id = ?", m.Id, liker.Id).
		Delete(&models.AIMaterialLike{}).Error; err != nil {
		t.Fatalf("delete like: %v", err)
	}
	got2, err := svc.GetByID(m.Id, liker.Id)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if got2.LikesCount != 0 {
		t.Errorf("LikesCount after delete = %d, want 0", got2.LikesCount)
	}
}
