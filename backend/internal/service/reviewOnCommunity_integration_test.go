package service

import (
	"testing"

	"gorm.io/gorm"

	"ithozyeva/internal/models"
	"ithozyeva/internal/testutil"
)

func reviewTablesTruncate(t *testing.T, db *gorm.DB) {
	testutil.TruncateAll(t, db, `"reviewOnCommunity"`, "members")
}

func TestReviewOnCommunityService_CreateByMemberId(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	reviewTablesTruncate(t, db)

	author := seedMemberWithRoles(t, db, 13001, "review_author", nil)
	svc := NewReviewOnCommunityService()

	if err := svc.CreateReviewOnCommunityByMemberId(author.Id, "круто", nil); err != nil {
		t.Fatalf("Create: %v", err)
	}

	reviews, err := svc.GetByAuthorId(author.Id)
	if err != nil {
		t.Fatalf("GetByAuthorId: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("ожидали 1 отзыв, got %d", len(reviews))
	}
	if reviews[0].Text != "круто" {
		t.Errorf("Text = %q, want круто", reviews[0].Text)
	}
	if reviews[0].Status != models.ReviewOnCommunityStatusDraft {
		t.Errorf("Status = %q, want DRAFT", reviews[0].Status)
	}
}

func TestReviewOnCommunityService_GetApproved_OnlyApproved(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	reviewTablesTruncate(t, db)

	author := seedMemberWithRoles(t, db, 13101, "approved_author", nil)
	svc := NewReviewOnCommunityService()

	if err := svc.CreateReviewOnCommunityByMemberId(author.Id, "draft review", nil); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if err := svc.CreateReviewOnCommunityByMemberId(author.Id, "approved review", nil); err != nil {
		t.Fatalf("create approved: %v", err)
	}

	all, err := svc.GetByAuthorId(author.Id)
	if err != nil || len(all) != 2 {
		t.Fatalf("ожидали 2 отзыва от автора, got %d (err=%v)", len(all), err)
	}

	// Аппрувим второй.
	var approved *models.ReviewOnCommunity
	for i := range all {
		if all[i].Text == "approved review" {
			approved = &all[i]
			break
		}
	}
	if approved == nil {
		t.Fatal("не нашли 'approved review' среди созданных")
	}
	if _, err := svc.Approve(int64(approved.Id)); err != nil {
		t.Fatalf("Approve: %v", err)
	}

	got, err := svc.GetApproved()
	if err != nil {
		t.Fatalf("GetApproved: %v", err)
	}
	if got == nil || len(*got) != 1 {
		t.Fatalf("ожидали один APPROVED, got %v", got)
	}
	if (*got)[0].Text != "approved review" {
		t.Errorf("approved Text = %q, want 'approved review'", (*got)[0].Text)
	}
}

func TestReviewOnCommunityService_Approve_Idempotent(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	reviewTablesTruncate(t, db)

	author := seedMemberWithRoles(t, db, 13201, "idem_author", nil)
	svc := NewReviewOnCommunityService()

	if err := svc.CreateReviewOnCommunityByMemberId(author.Id, "to approve", nil); err != nil {
		t.Fatalf("create: %v", err)
	}

	reviews, _ := svc.GetByAuthorId(author.Id)
	if len(reviews) != 1 {
		t.Fatalf("ожидали 1 отзыв, got %d", len(reviews))
	}
	id := int64(reviews[0].Id)

	first, err := svc.Approve(id)
	if err != nil {
		t.Fatalf("первый Approve: %v", err)
	}
	if first.Status != models.ReviewOnCommunityStatusApproved {
		t.Errorf("status после первого Approve = %q, want APPROVED", first.Status)
	}

	second, err := svc.Approve(id)
	if err != nil {
		t.Fatalf("повторный Approve: %v", err)
	}
	if second.Status != models.ReviewOnCommunityStatusApproved {
		t.Errorf("status после повторного Approve = %q, want APPROVED", second.Status)
	}
}

func TestReviewOnCommunityService_GetByAuthorId_FiltersOtherAuthors(t *testing.T) {
	db := testutil.EnsureTestDB(t)
	reviewTablesTruncate(t, db)

	a := seedMemberWithRoles(t, db, 13301, "author_a", nil)
	b := seedMemberWithRoles(t, db, 13302, "author_b", nil)

	svc := NewReviewOnCommunityService()
	if err := svc.CreateReviewOnCommunityByMemberId(a.Id, "from a", nil); err != nil {
		t.Fatalf("a: %v", err)
	}
	if err := svc.CreateReviewOnCommunityByMemberId(b.Id, "from b", nil); err != nil {
		t.Fatalf("b: %v", err)
	}

	got, err := svc.GetByAuthorId(a.Id)
	if err != nil {
		t.Fatalf("GetByAuthorId: %v", err)
	}
	if len(got) != 1 || got[0].Text != "from a" {
		t.Errorf("ожидали один отзыв 'from a', got %+v", got)
	}
}
