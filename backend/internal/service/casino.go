package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"

	"ithozyeva/internal/models"
	"ithozyeva/internal/repository"
)

var wheelMultipliers = []float64{0, 0, 0, 0.5, 0.5, 0.5, 1, 1, 1.5, 1.5, 2, 3}

type CasinoService struct {
	repo      *repository.CasinoRepository
	pointRepo *repository.PointsRepository
}

func NewCasinoService() *CasinoService {
	return &CasinoService{
		repo:      repository.NewCasinoRepository(),
		pointRepo: repository.NewPointsRepository(),
	}
}

func (s *CasinoService) validateBet(amount int) error {
	if amount < 10 {
		return fmt.Errorf("минимальная ставка: 10 баллов")
	}
	if amount > 1000 {
		return fmt.Errorf("максимальная ставка: 1000 баллов")
	}
	return nil
}

func cryptoRandInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func (s *CasinoService) PlayCoinFlip(memberId int64, req *models.CoinFlipRequest) (*models.CasinoBetResponse, error) {
	if err := s.validateBet(req.BetAmount); err != nil {
		return nil, err
	}
	if req.Choice != "heads" && req.Choice != "tails" {
		return nil, fmt.Errorf("выберите heads или tails")
	}

	flip, err := cryptoRandInt(2)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации")
	}

	outcome := "heads"
	if flip == 1 {
		outcome = "tails"
	}

	won := outcome == req.Choice
	var multiplier float64
	var payout int
	if won {
		multiplier = 1.9
		payout = int(float64(req.BetAmount) * multiplier)
	}
	profit := payout - req.BetAmount

	bet := &models.CasinoBet{
		MemberId:   memberId,
		Game:       "coin_flip",
		BetAmount:  req.BetAmount,
		BetChoice:  req.Choice,
		Result:     outcome,
		Multiplier: multiplier,
		Payout:     payout,
		Profit:     profit,
	}

	balance, err := s.repo.PlaceBet(memberId, bet, won)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			return nil, err
		}
		log.Printf("casino PlaceBet error for member %d: %v", memberId, err)
		return nil, fmt.Errorf("ошибка при размещении ставки")
	}

	return &models.CasinoBetResponse{
		Id:         bet.Id,
		Game:       bet.Game,
		BetAmount:  bet.BetAmount,
		BetChoice:  bet.BetChoice,
		Result:     bet.Result,
		Multiplier: bet.Multiplier,
		Payout:     bet.Payout,
		Profit:     bet.Profit,
		Balance:    balance,
		CreatedAt:  bet.CreatedAt,
	}, nil
}

func (s *CasinoService) PlayDiceRoll(memberId int64, req *models.DiceRollRequest) (*models.CasinoBetResponse, error) {
	if err := s.validateBet(req.BetAmount); err != nil {
		return nil, err
	}
	if req.Target < 2 || req.Target > 98 {
		return nil, fmt.Errorf("цель должна быть от 2 до 98")
	}
	if req.Direction != "over" && req.Direction != "under" {
		return nil, fmt.Errorf("выберите over или under")
	}

	roll, err := cryptoRandInt(100)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации")
	}

	var winChance float64
	if req.Direction == "over" {
		winChance = float64(100 - req.Target)
	} else {
		winChance = float64(req.Target)
	}

	var won bool
	if req.Direction == "over" {
		won = roll > req.Target
	} else {
		won = roll < req.Target
	}

	var multiplier float64
	var payout int
	if won {
		multiplier = 0.97 * (100.0 / winChance)
		payout = int(float64(req.BetAmount) * multiplier)
	}
	profit := payout - req.BetAmount

	betChoice := fmt.Sprintf("%s %d", req.Direction, req.Target)
	result := fmt.Sprintf("%d", roll)

	bet := &models.CasinoBet{
		MemberId:   memberId,
		Game:       "dice_roll",
		BetAmount:  req.BetAmount,
		BetChoice:  betChoice,
		Result:     result,
		Multiplier: multiplier,
		Payout:     payout,
		Profit:     profit,
	}

	balance, err := s.repo.PlaceBet(memberId, bet, won)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			return nil, err
		}
		log.Printf("casino PlaceBet error for member %d: %v", memberId, err)
		return nil, fmt.Errorf("ошибка при размещении ставки")
	}

	return &models.CasinoBetResponse{
		Id:         bet.Id,
		Game:       bet.Game,
		BetAmount:  bet.BetAmount,
		BetChoice:  bet.BetChoice,
		Result:     bet.Result,
		Multiplier: bet.Multiplier,
		Payout:     bet.Payout,
		Profit:     bet.Profit,
		Balance:    balance,
		CreatedAt:  bet.CreatedAt,
	}, nil
}

func (s *CasinoService) PlayWheel(memberId int64, req *models.WheelRequest) (*models.CasinoBetResponse, error) {
	if err := s.validateBet(req.BetAmount); err != nil {
		return nil, err
	}

	segment, err := cryptoRandInt(len(wheelMultipliers))
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации")
	}

	multiplier := wheelMultipliers[segment]
	payout := int(float64(req.BetAmount) * multiplier)
	won := payout > 0
	profit := payout - req.BetAmount

	result := fmt.Sprintf("x%.1f", multiplier)

	bet := &models.CasinoBet{
		MemberId:   memberId,
		Game:       "wheel",
		BetAmount:  req.BetAmount,
		BetChoice:  "spin",
		Result:     result,
		Multiplier: multiplier,
		Payout:     payout,
		Profit:     profit,
	}

	balance, err := s.repo.PlaceBet(memberId, bet, won)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			return nil, err
		}
		log.Printf("casino PlaceBet error for member %d: %v", memberId, err)
		return nil, fmt.Errorf("ошибка при размещении ставки")
	}

	return &models.CasinoBetResponse{
		Id:         bet.Id,
		Game:       bet.Game,
		BetAmount:  bet.BetAmount,
		BetChoice:  bet.BetChoice,
		Result:     bet.Result,
		Multiplier: bet.Multiplier,
		Payout:     bet.Payout,
		Profit:     bet.Profit,
		Balance:    balance,
		CreatedAt:  bet.CreatedAt,
	}, nil
}

func (s *CasinoService) GetGlobalFeed(limit int) ([]models.CasinoFeedItem, error) {
	return s.repo.GetGlobalFeed(limit)
}

func (s *CasinoService) GetHistory(memberId int64, limit, offset int) ([]models.CasinoHistoryItem, int64, error) {
	return s.repo.GetHistory(memberId, limit, offset)
}

func (s *CasinoService) GetStats(memberId int64) (*models.CasinoStats, error) {
	stats, err := s.repo.GetStats(memberId)
	if err != nil {
		return nil, err
	}
	balance, err := s.pointRepo.GetBalance(memberId)
	if err != nil {
		return nil, err
	}
	stats.Balance = balance
	return stats, nil
}

func (s *CasinoService) GetAdminStats() (*models.CasinoAdminStats, error) {
	return s.repo.GetAdminStats()
}

func (s *CasinoService) SearchBets(username *string, game *string, limit, offset int) ([]models.CasinoAdminBet, int64, error) {
	return s.repo.SearchBets(username, game, limit, offset)
}
