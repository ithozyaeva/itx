package service

import (
	"sync"
	"testing"

	"ithozyeva/internal/repository"
	"ithozyeva/internal/testutil"
)

// Integration-тесты для ReferralCreditService.AdminAward → repo.Spend.
//
// Контекст: PR #347 ввёл Spend с `SELECT … FOR UPDATE` поверх SUM(amount),
// что PostgreSQL отвергает с parse-error «FOR UPDATE is not allowed with
// aggregate functions». В результате весь Spend (и через него
// PurchaseTierWithCredits + admin-списания) тихо падал 500 в проде 6 дней.
// Тестов на этот путь не было, поэтому регресс не поймали.
//
// Эти тесты гарантируют, что:
//  1. Spend через AdminAward с отрицательной суммой проходит без SQL-ошибок.
//  2. ErrInsufficientCredits возвращается, когда баланс < amount.
//  3. Advisory-lock реально сериализует параллельные Spend'ы одного юзера
//     (без него два concurrent /play с balance=N оба прошли бы balance-check).

func TestReferralCreditService_AdminAward_PositiveAwardSucceeds(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "referral_credit_transactions", "members")

	member := seedMember(t, db, 7101)
	svc := NewReferralCreditService()

	if err := svc.AdminAward(member.Id, 100, "test positive"); err != nil {
		t.Fatalf("AdminAward(+100): %v", err)
	}

	balance, err := svc.GetBalance(member.Id)
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if balance != 100 {
		t.Errorf("balance = %d, want 100", balance)
	}
}

// TestReferralCreditService_AdminAward_NegativeSpendsSuccessfully — ключевой
// тест: AdminAward с отрицательной суммой идёт через repo.Spend. До правки
// H1 этот вызов всегда падал с SQL parse-error.
func TestReferralCreditService_AdminAward_NegativeSpendsSuccessfully(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "referral_credit_transactions", "members")

	member := seedMember(t, db, 7102)
	svc := NewReferralCreditService()

	// Сначала кладём кредиты на счёт, чтобы было что списывать.
	if err := svc.AdminAward(member.Id, 200, "seed"); err != nil {
		t.Fatalf("AdminAward(+200) seed: %v", err)
	}

	// Spend через AdminAward(-50). Регресс H1 проявился бы здесь:
	// "FOR UPDATE is not allowed with aggregate functions" → 500.
	if err := svc.AdminAward(member.Id, -50, "test spend"); err != nil {
		t.Fatalf("AdminAward(-50): %v", err)
	}

	balance, err := svc.GetBalance(member.Id)
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if balance != 150 {
		t.Errorf("balance = %d, want 150 (200 - 50)", balance)
	}
}

func TestReferralCreditService_AdminAward_InsufficientCredits(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "referral_credit_transactions", "members")

	member := seedMember(t, db, 7103)
	svc := NewReferralCreditService()

	if err := svc.AdminAward(member.Id, 30, "seed small"); err != nil {
		t.Fatalf("AdminAward(+30) seed: %v", err)
	}

	// Пытаемся списать больше, чем есть — должна вернуться
	// ErrInsufficientCredits (через repo.Spend → ErrInsufficientCredits).
	err := svc.AdminAward(member.Id, -100, "overdraw")
	if err == nil {
		t.Fatal("AdminAward(-100) при balance=30 должен вернуть ошибку")
	}
	if err != repository.ErrInsufficientCredits {
		t.Errorf("ошибка = %v, want ErrInsufficientCredits", err)
	}

	// Баланс не изменился — транзакция откатилась.
	balance, err := svc.GetBalance(member.Id)
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if balance != 30 {
		t.Errorf("balance после неудачного Spend = %d, want 30 (без изменений)", balance)
	}
}

// TestReferralCreditService_AdminAward_ConcurrentSpend проверяет, что
// pg_advisory_xact_lock реально сериализует параллельные Spend'ы того же
// юзера. До правки M2 (где тот же паттерн применён в casino) такой
// тест поймал бы TOCTOU — два горутины с balance=100 и spend=80 могли бы
// оба пройти проверку и оба создать debit-запись, увода баланс в минус.
//
// Здесь проверяем для credits — Spend и так был бы сломан без advisory_lock,
// но после фикса нужно убедиться, что лок именно сериализует.
func TestReferralCreditService_AdminAward_ConcurrentSpend(t *testing.T) {
	db := testutil.SetupTestDB(t)
	testutil.TruncateAll(t, db, "referral_credit_transactions", "members")

	member := seedMember(t, db, 7104)
	svc := NewReferralCreditService()

	if err := svc.AdminAward(member.Id, 100, "seed for race"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// 5 горутин одновременно пытаются списать по 80. Из них:
	// - 1 должна пройти (balance 100 → 20).
	// - 4 должны получить ErrInsufficientCredits.
	// Без advisory_lock несколько прошли бы и баланс ушёл бы в минус.
	const concurrent = 5
	const spendEach = 80

	var wg sync.WaitGroup
	results := make(chan error, concurrent)
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- svc.AdminAward(member.Id, -spendEach, "concurrent spend")
		}()
	}
	wg.Wait()
	close(results)

	successCount := 0
	insufficientCount := 0
	for err := range results {
		switch err {
		case nil:
			successCount++
		case repository.ErrInsufficientCredits:
			insufficientCount++
		default:
			t.Errorf("unexpected error: %v", err)
		}
	}

	if successCount != 1 {
		t.Errorf("успешных Spend = %d, want 1 (advisory_lock должен сериализовать)", successCount)
	}
	if insufficientCount != concurrent-1 {
		t.Errorf("отказов = %d, want %d", insufficientCount, concurrent-1)
	}

	balance, err := svc.GetBalance(member.Id)
	if err != nil {
		t.Fatalf("GetBalance: %v", err)
	}
	if balance != 100-spendEach {
		t.Errorf("итоговый balance = %d, want %d (баланс не должен уйти в минус)", balance, 100-spendEach)
	}
	if balance < 0 {
		t.Errorf("balance < 0 — TOCTOU race не закрыт")
	}
}
