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

func TestAIMaterialService_ToggleLike_RoundTrip(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_likes", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9301)
	liker := seedMember(t, db, 9302)
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

	// Первый toggle — лайк ставится, count=1
	liked, count, err := svc.ToggleLike(m.Id, liker.Id)
	if err != nil {
		t.Fatalf("toggle 1: %v", err)
	}
	if !liked || count != 1 {
		t.Errorf("after first toggle: liked=%v count=%d, want true 1", liked, count)
	}

	// Повторный toggle — лайк снимается, count=0
	liked, count, err = svc.ToggleLike(m.Id, liker.Id)
	if err != nil {
		t.Fatalf("toggle 2: %v", err)
	}
	if liked || count != 0 {
		t.Errorf("after second toggle: liked=%v count=%d, want false 0", liked, count)
	}

	// Toggle несуществующего материала
	if _, _, err := svc.ToggleLike(99999, liker.Id); err == nil {
		t.Error("expected error for non-existing material")
	}
}

func TestAIMaterialService_ToggleBookmark_RoundTrip(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_bookmarks", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9401)
	saver := seedMember(t, db, 9402)
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

	bookmarked, count, err := svc.ToggleBookmark(m.Id, saver.Id)
	if err != nil {
		t.Fatalf("toggle 1: %v", err)
	}
	if !bookmarked || count != 1 {
		t.Errorf("after first toggle: bookmarked=%v count=%d, want true 1", bookmarked, count)
	}

	// Через bookmarked-фильтр saver видит этот материал
	items, total, err := svc.Search(repository.AIMaterialFilter{
		Bookmarked: true, ViewerID: saver.Id,
	})
	if err != nil {
		t.Fatalf("search bookmarked: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Id != m.Id {
		t.Errorf("bookmarked filter: total=%d ids=%v", total, idsOf(items))
	}
	if !items[0].Bookmarked {
		t.Error("Bookmarked флаг должен быть true для saver")
	}

	bookmarked, count, err = svc.ToggleBookmark(m.Id, saver.Id)
	if err != nil {
		t.Fatalf("toggle 2: %v", err)
	}
	if bookmarked || count != 0 {
		t.Errorf("after second toggle: bookmarked=%v count=%d, want false 0", bookmarked, count)
	}
}

func TestAIMaterialService_Comments_CRUD(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_comments", "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9501)
	commenter := seedMember(t, db, 9502)
	stranger := seedMember(t, db, 9503)
	svc := NewAIMaterialService()

	m, err := svc.Create(&models.CreateAIMaterialRequest{
		Title:        "Title",
		Summary:      strings.Repeat("a", 35),
		ContentType:  models.AIMaterialContentTypePrompt,
		MaterialKind: models.AIMaterialKindPrompt,
		PromptBody:   "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create material: %v", err)
	}

	// Create
	c1, err := svc.CreateComment(m.Id, commenter.Id, "  Полезный материал, спасибо!  ")
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	if c1.Body != "Полезный материал, спасибо!" {
		t.Errorf("Body не trim: %q", c1.Body)
	}

	// Триггер увеличил counter материала
	got, err := svc.GetByID(m.Id, 0)
	if err != nil {
		t.Fatalf("get material: %v", err)
	}
	if got.CommentsCount != 1 {
		t.Errorf("CommentsCount = %d, want 1", got.CommentsCount)
	}

	// Пустой комментарий — ошибка
	if _, err := svc.CreateComment(m.Id, commenter.Id, "   "); err == nil {
		t.Error("expected error for empty comment")
	}

	// List
	list, err := svc.ListComments(m.Id, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].Id != c1.Id {
		t.Errorf("list mismatch: %+v", list)
	}

	// Update — чужой не может
	if _, err := svc.UpdateComment(c1.Id, stranger.Id, "Hacked", false); err == nil {
		t.Error("stranger should not update")
	}

	// Update — автор может
	updated, err := svc.UpdateComment(c1.Id, commenter.Id, "Обновлённый текст", false)
	if err != nil {
		t.Fatalf("author update: %v", err)
	}
	if updated.Body != "Обновлённый текст" {
		t.Errorf("Body не обновился: %q", updated.Body)
	}

	// Delete — чужой не может
	if err := svc.DeleteComment(c1.Id, stranger.Id, false); err == nil {
		t.Error("stranger should not delete")
	}

	// Delete — admin может (даже если не автор)
	if err := svc.DeleteComment(c1.Id, stranger.Id, true); err != nil {
		t.Fatalf("admin delete: %v", err)
	}

	// Триггер обнулил счётчик
	got, err = svc.GetByID(m.Id, 0)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if got.CommentsCount != 0 {
		t.Errorf("CommentsCount после delete = %d, want 0", got.CommentsCount)
	}
}

func idsOf(items []models.AIMaterial) []int64 {
	out := make([]int64, len(items))
	for i, it := range items {
		out[i] = it.Id
	}
	return out
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
