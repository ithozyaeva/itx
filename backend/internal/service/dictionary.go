package service

import (
	"ithozyeva/internal/models"
)

type DictionaryService struct{}

func NewDictionaryService() *DictionaryService {
	return &DictionaryService{}
}

type DictionaryItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type DictionaryMap map[string][]DictionaryItem

func (s *DictionaryService) GetAllDictionaries() DictionaryMap {
	return DictionaryMap{
		"placeTypes": {
			{Value: string(models.EventOnline), Label: "Онлайн"},
			{Value: string(models.EventOffline), Label: "Оффлайн"},
			{Value: string(models.EventHybrid), Label: "Гибрид"},
		},
		"memberRoles": {
			{Value: string(models.MemberRoleUnsubscriber), Label: "Ансаб"},
			{Value: string(models.MemberRoleSubscriber), Label: "Саб"},
			{Value: string(models.MemberRoleMentor), Label: "Ментор"},
			{Value: string(models.MemberRoleAdmin), Label: "Админ"},
			{Value: string(models.MemberRoleEventMaker), Label: "Ивентмейкер"},
		},
		"reviewStatuses": {
			{Value: string(models.ReviewOnCommunityStatusDraft), Label: "На модерации"},
			{Value: string(models.ReviewOnCommunityStatusApproved), Label: "Опубликован"},
		},
		"grades": {
			{Value: string(models.SeniorGrade), Label: "Сеньор"},
			{Value: string(models.JuniorGrade), Label: "Джун"},
			{Value: string(models.MiddleGrade), Label: "Мидл"},
		},
		"referalLinkStatuses": {
			{Value: string(models.ReferalLinkActive), Label: "В поиске"},
			{Value: string(models.ReferalLinkFreezed), Label: "Заморожен"},
		},
	}
}
