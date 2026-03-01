package models

type AchievementCategory string

const (
	AchievementCategoryEvents   AchievementCategory = "events"
	AchievementCategoryPoints   AchievementCategory = "points"
	AchievementCategorySocial   AchievementCategory = "social"
	AchievementCategoryActivity AchievementCategory = "activity"
)

type Achievement struct {
	Id          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Icon        string              `json:"icon"`
	Category    AchievementCategory `json:"category"`
	Threshold   int                 `json:"threshold"`
}

type UserAchievement struct {
	Achievement
	Unlocked bool `json:"unlocked"`
	Progress int  `json:"progress"`
}

type AchievementsResponse struct {
	Items         []UserAchievement `json:"items"`
	TotalCount    int               `json:"totalCount"`
	UnlockedCount int               `json:"unlockedCount"`
}
