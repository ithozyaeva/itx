package repository

import (
	"ithozyeva/database"
	"ithozyeva/internal/models"
)

type GuildRepository struct{}

func NewGuildRepository() *GuildRepository {
	return &GuildRepository{}
}

func (r *GuildRepository) Create(guild *models.Guild) error {
	return database.DB.Create(guild).Error
}

func (r *GuildRepository) GetById(id int64) (*models.Guild, error) {
	var guild models.Guild
	err := database.DB.Preload("Members").First(&guild, id).Error
	if err != nil {
		return nil, err
	}
	return &guild, nil
}

func (r *GuildRepository) GetAll(memberId int64) ([]models.GuildPublic, error) {
	items := make([]models.GuildPublic, 0)
	err := database.DB.Raw(`
		SELECT g.id, g.name, g.description, g.icon, g.color, g.owner_id,
			m.first_name as owner_first_name, m.last_name as owner_last_name,
			m.username as owner_username, m.avatar_url as owner_avatar_url,
			(SELECT COUNT(*) FROM guild_members WHERE guild_id = g.id) as member_count,
			COALESCE((
				SELECT SUM(pt.amount) FROM point_transactions pt
				JOIN guild_members gm ON gm.member_id = pt.member_id AND gm.guild_id = g.id
			), 0) as total_points,
			EXISTS(SELECT 1 FROM guild_members WHERE guild_id = g.id AND member_id = ?) as is_member
		FROM guilds g
		JOIN members m ON m.id = g.owner_id
		ORDER BY total_points DESC
	`, memberId).Scan(&items).Error
	return items, err
}

func (r *GuildRepository) Join(guildId, memberId int64) error {
	gm := &models.GuildMember{GuildId: guildId, MemberId: memberId}
	return database.DB.Create(gm).Error
}

func (r *GuildRepository) Leave(guildId, memberId int64) error {
	return database.DB.Where("guild_id = ? AND member_id = ?", guildId, memberId).
		Delete(&models.GuildMember{}).Error
}

func (r *GuildRepository) IsMember(guildId, memberId int64) bool {
	var count int64
	database.DB.Model(&models.GuildMember{}).
		Where("guild_id = ? AND member_id = ?", guildId, memberId).
		Count(&count)
	return count > 0
}

func (r *GuildRepository) GetMemberGuildId(memberId int64) (int64, error) {
	var gm models.GuildMember
	err := database.DB.Where("member_id = ?", memberId).First(&gm).Error
	if err != nil {
		return 0, err
	}
	return gm.GuildId, nil
}

func (r *GuildRepository) GetMemberGuildName(memberId int64) string {
	var name string
	database.DB.Raw(`
		SELECT g.name FROM guilds g
		JOIN guild_members gm ON gm.guild_id = g.id
		WHERE gm.member_id = ?
	`, memberId).Scan(&name)
	return name
}

func (r *GuildRepository) Update(guild *models.Guild) error {
	return database.DB.Model(&models.Guild{}).Where("id = ?", guild.Id).
		Updates(map[string]interface{}{
			"name":        guild.Name,
			"description": guild.Description,
			"icon":        guild.Icon,
			"color":       guild.Color,
		}).Error
}

func (r *GuildRepository) Delete(id int64) error {
	return database.DB.Delete(&models.Guild{}, id).Error
}

func (r *GuildRepository) GetGuildMembers(guildId int64) ([]models.MemberPointsBalance, error) {
	members := make([]models.MemberPointsBalance, 0)
	err := database.DB.Raw(`
		SELECT m.id as member_id, m.first_name, m.last_name, m.username, m.avatar_url,
			COALESCE(SUM(pt.amount), 0) as total
		FROM members m
		JOIN guild_members gm ON gm.member_id = m.id AND gm.guild_id = ?
		LEFT JOIN point_transactions pt ON pt.member_id = m.id
		GROUP BY m.id, m.first_name, m.last_name, m.username, m.avatar_url
		ORDER BY total DESC
	`, guildId).Scan(&members).Error
	return members, err
}
