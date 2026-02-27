package models

type Role string

const (
	MemberRoleUnsubscriber Role = "UNSUBSCRIBER"
	MemberRoleSubscriber   Role = "SUBSCRIBER"
	MemberRoleMentor       Role = "MENTOR"
	MemberRoleAdmin        Role = "ADMIN"
	MemberRoleEventMaker   Role = "EVENT_MAKER"
)

type Permission string

const (
	PermissionCanViewAdminPanel            Permission = "can_view_admin_panel"
	PermissionCanViewAdminMembers          Permission = "can_view_admin_members"
	PermissionCanViewAdminMentors          Permission = "can_view_admin_mentors"
	PermissionCanViewAdminEvents           Permission = "can_view_admin_events"
	PermissionCanEditAdminMembers          Permission = "can_edit_admin_members"
	PermissionCanEditAdminMentors          Permission = "can_edit_admin_mentors"
	PermissionCanEditAdminEvents           Permission = "can_edit_admin_events"
	PermissionCanViewAdminReviews          Permission = "can_view_admin_reviews"
	PermissionCanEditAdminReviews          Permission = "can_edit_admin_reviews"
	PermissionCanApprovedAdminReviews      Permission = "can_approved_admin_reviews"
	PermissionCanViewAdminMentorsReview    Permission = "can_view_admin_mentors_review"
	PermissionCanEditAdminMentorsReview    Permission = "can_edit_admin_mentors_review"
	PermissionCanApproveAdminMentorsReview Permission = "can_approve_admin_mentors_review"
	PermissionCanViewAdminResumes          Permission = "can_view_admin_resumes"
	PermissionCanEditPlatformMentors       Permission = "can_edit_platform_mentor"
	PermissionCanViewAdminAuditLogs        Permission = "can_view_admin_audit_logs"
)

type PermissionModel struct {
	Id   int64  `json:"-" gorm:"primaryKey"`
	Name string `json:"name" gorm:"column:name"`
}
