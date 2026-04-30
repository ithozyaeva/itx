package service

import (
	"errors"
	"fmt"
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
	liked, count, err := svc.ToggleLike(m.Id, liker.Id, false)
	if err != nil {
		t.Fatalf("toggle 1: %v", err)
	}
	if !liked || count != 1 {
		t.Errorf("after first toggle: liked=%v count=%d, want true 1", liked, count)
	}

	// Повторный toggle — лайк снимается, count=0
	liked, count, err = svc.ToggleLike(m.Id, liker.Id, false)
	if err != nil {
		t.Fatalf("toggle 2: %v", err)
	}
	if liked || count != 0 {
		t.Errorf("after second toggle: liked=%v count=%d, want false 0", liked, count)
	}

	// Toggle несуществующего материала
	if _, _, err := svc.ToggleLike(99999, liker.Id, false); err == nil {
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

	bookmarked, count, err := svc.ToggleBookmark(m.Id, saver.Id, false)
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

	bookmarked, count, err = svc.ToggleBookmark(m.Id, saver.Id, false)
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
	c1, err := svc.CreateComment(m.Id, commenter.Id, "  Полезный материал, спасибо!  ", false)
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	if c1.Body != "Полезный материал, спасибо!" {
		t.Errorf("Body не trim: %q", c1.Body)
	}

	// Триггер увеличил counter материала
	got, err := svc.GetByID(m.Id, 0, true)
	if err != nil {
		t.Fatalf("get material: %v", err)
	}
	if got.CommentsCount != 1 {
		t.Errorf("CommentsCount = %d, want 1", got.CommentsCount)
	}

	// Пустой комментарий — ошибка
	if _, err := svc.CreateComment(m.Id, commenter.Id, "   ", false); err == nil {
		t.Error("expected error for empty comment")
	}

	// List
	list, total, err := svc.ListComments(m.Id, 0, false, 20, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Id != c1.Id {
		t.Errorf("list mismatch: total=%d items=%+v", total, list)
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
	got, err = svc.GetByID(m.Id, 0, true)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if got.CommentsCount != 0 {
		t.Errorf("CommentsCount после delete = %d, want 0", got.CommentsCount)
	}
}

func TestAIMaterialService_ToggleCommentLike_RoundTrip(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_comment_likes", "ai_material_comments",
		"ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9901)
	commenter := seedMember(t, db, 9902)
	liker := seedMember(t, db, 9903)
	svc := NewAIMaterialService()

	m, err := svc.Create(&models.CreateAIMaterialRequest{
		Title: "Title", Summary: strings.Repeat("a", 35),
		ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
		PromptBody: "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create material: %v", err)
	}
	c, err := svc.CreateComment(m.Id, commenter.Id, "Полезно", false)
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}

	// Toggle ON — счётчик +1
	liked, count, err := svc.ToggleCommentLike(c.Id, liker.Id, false)
	if err != nil {
		t.Fatalf("toggle 1: %v", err)
	}
	if !liked || count != 1 {
		t.Errorf("first toggle: liked=%v count=%d, want true 1", liked, count)
	}

	// ListComments под viewer=liker — поле Liked=true
	list, _, err := svc.ListComments(m.Id, liker.Id, false, 20, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || !list[0].Liked || list[0].LikesCount != 1 {
		t.Errorf("list under liker: liked=%v count=%d", list[0].Liked, list[0].LikesCount)
	}

	// Под другим viewer — Liked=false, но count=1
	list2, _, err := svc.ListComments(m.Id, author.Id, false, 20, 0)
	if err != nil {
		t.Fatalf("list under author: %v", err)
	}
	if list2[0].Liked || list2[0].LikesCount != 1 {
		t.Errorf("under author: liked=%v count=%d, want false 1", list2[0].Liked, list2[0].LikesCount)
	}

	// Toggle OFF
	liked, count, err = svc.ToggleCommentLike(c.Id, liker.Id, false)
	if err != nil {
		t.Fatalf("toggle 2: %v", err)
	}
	if liked || count != 0 {
		t.Errorf("second toggle: liked=%v count=%d, want false 0", liked, count)
	}

	// Несуществующий коммент
	if _, _, err := svc.ToggleCommentLike(99999, liker.Id, false); !errors.Is(err, ErrAIMaterialCommentNotFound) {
		t.Errorf("missing comment: want ErrCommentNotFound, got %v", err)
	}
}

func TestAIMaterialService_ToggleCommentLike_DeniedOnHiddenMaterial(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_comment_likes", "ai_material_comments",
		"ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9911)
	stranger := seedMember(t, db, 9912)
	svc := NewAIMaterialService()

	m, _ := svc.Create(&models.CreateAIMaterialRequest{
		Title: "Title", Summary: strings.Repeat("a", 35),
		ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
		PromptBody: "x",
	}, author.Id)
	c, _ := svc.CreateComment(m.Id, author.Id, "сам себе", false)

	// Скрываем материал
	if err := svc.SetHidden(m.Id, true, true); err != nil {
		t.Fatalf("hide: %v", err)
	}

	// Чужой не должен лайкать коммент скрытого материала
	if _, _, err := svc.ToggleCommentLike(c.Id, stranger.Id, false); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger like on hidden: want ErrNotFound, got %v", err)
	}
	// Админ — может
	if _, _, err := svc.ToggleCommentLike(c.Id, stranger.Id, true); err != nil {
		t.Errorf("admin like on hidden: should pass, got %v", err)
	}
}

func TestAIMaterialService_Hidden_VisibilityAndInteractions(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_likes", "ai_material_bookmarks", "ai_material_comments",
		"ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9601)
	stranger := seedMember(t, db, 9602)
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

	// Скрываем как админ
	if err := svc.SetHidden(m.Id, true, true); err != nil {
		t.Fatalf("hide: %v", err)
	}

	// Чужой пользователь получает «не найден» — без утечки факта существования
	if _, err := svc.GetByID(m.Id, stranger.Id, false); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger should get ErrNotFound, got %v", err)
	}
	// Автор по-прежнему видит свой скрытый материал
	if _, err := svc.GetByID(m.Id, author.Id, false); err != nil {
		t.Errorf("author should still see hidden material: %v", err)
	}
	// Админ видит даже чужой скрытый
	if _, err := svc.GetByID(m.Id, stranger.Id, true); err != nil {
		t.Errorf("admin should see hidden material: %v", err)
	}

	// Чужой не может лайкать/закладывать/комментить скрытое
	if _, _, err := svc.ToggleLike(m.Id, stranger.Id, false); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger toggle-like on hidden: want ErrNotFound, got %v", err)
	}
	if _, _, err := svc.ToggleBookmark(m.Id, stranger.Id, false); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger bookmark on hidden: want ErrNotFound, got %v", err)
	}
	if _, err := svc.CreateComment(m.Id, stranger.Id, "hi", false); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger comment on hidden: want ErrNotFound, got %v", err)
	}
	if _, _, err := svc.ListComments(m.Id, stranger.Id, false, 20, 0); !errors.Is(err, ErrAIMaterialNotFound) {
		t.Errorf("stranger list comments on hidden: want ErrNotFound, got %v", err)
	}
}

func TestAIMaterialService_TopTags_DoesNotLeakHiddenTags(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db,
		"ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9701)
	svc := NewAIMaterialService()

	mVisible, err := svc.Create(&models.CreateAIMaterialRequest{
		Title: "Visible", Summary: strings.Repeat("a", 35),
		ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
		PromptBody: "x", Tags: []string{"open"},
	}, author.Id)
	if err != nil {
		t.Fatalf("create visible: %v", err)
	}
	_ = mVisible

	mHidden, err := svc.Create(&models.CreateAIMaterialRequest{
		Title: "Hidden", Summary: strings.Repeat("a", 35),
		ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
		PromptBody: "x", Tags: []string{"secret-leak"},
	}, author.Id)
	if err != nil {
		t.Fatalf("create hidden: %v", err)
	}
	if err := svc.SetHidden(mHidden.Id, true, true); err != nil {
		t.Fatalf("hide: %v", err)
	}

	tags, err := svc.TopTags("", 50)
	if err != nil {
		t.Fatalf("top tags: %v", err)
	}
	for _, tag := range tags {
		if tag == "secret-leak" {
			t.Errorf("hidden tag leaked into TopTags: %v", tags)
		}
	}
}

func TestAIMaterialService_Search_QueryEscapesLikeWildcards(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "ai_material_tags", "ai_materials", "members")

	author := seedMember(t, db, 9801)
	svc := NewAIMaterialService()

	// Без `%` в заголовках/описании
	for i := 0; i < 3; i++ {
		if _, err := svc.Create(&models.CreateAIMaterialRequest{
			Title: fmt.Sprintf("Plain title %d", i), Summary: strings.Repeat("a", 35),
			ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
			PromptBody: "x",
		}, author.Id); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}
	// С `%` — ищем именно его
	withPercent, err := svc.Create(&models.CreateAIMaterialRequest{
		Title: "Скидка 50% на курс", Summary: strings.Repeat("a", 35),
		ContentType: models.AIMaterialContentTypePrompt, MaterialKind: models.AIMaterialKindPrompt,
		PromptBody: "x",
	}, author.Id)
	if err != nil {
		t.Fatalf("create with %%: %v", err)
	}

	_, total, err := svc.Search(repository.AIMaterialFilter{Query: "%"})
	if err != nil {
		t.Fatalf("search %%: %v", err)
	}
	if total != 1 {
		t.Errorf("Search(q=%%): total=%d, want 1 (only the title with literal %%); without escape it would return all 4", total)
	}

	items, _, err := svc.Search(repository.AIMaterialFilter{Query: "50%"})
	if err != nil {
		t.Fatalf("search 50%%: %v", err)
	}
	if len(items) != 1 || items[0].Id != withPercent.Id {
		t.Errorf("Search(q=50%%) returned %d items, want exactly the percent-title", len(items))
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

	got, err := svc.GetByID(m.Id, liker.Id, false)
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
	got2, err := svc.GetByID(m.Id, liker.Id, false)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if got2.LikesCount != 0 {
		t.Errorf("LikesCount after delete = %d, want 0", got2.LikesCount)
	}
}
