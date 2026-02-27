/**
 * Permission types for admin-frontend
 * Union type of all available permissions
 */

export type Permission
  = | 'can_view_admin_panel'
    | 'can_view_admin_members'
    | 'can_view_admin_mentors'
    | 'can_view_admin_events'
    | 'can_edit_admin_members'
    | 'can_edit_admin_mentors'
    | 'can_edit_admin_events'
    | 'can_view_admin_reviews'
    | 'can_edit_admin_reviews'
    | 'can_approved_admin_reviews'
    | 'can_view_admin_mentors_review'
    | 'can_edit_admin_mentors_review'
    | 'can_approve_admin_mentors_review'
    | 'can_view_admin_resumes'

/**
 * Type guard for checking if a string is a valid permission
 */
export function isPermission(value: string): value is Permission {
  const permissions: Permission[] = [
    'can_view_admin_panel',
    'can_view_admin_members',
    'can_view_admin_mentors',
    'can_view_admin_events',
    'can_edit_admin_members',
    'can_edit_admin_mentors',
    'can_edit_admin_events',
    'can_view_admin_reviews',
    'can_edit_admin_reviews',
    'can_approved_admin_reviews',
    'can_view_admin_mentors_review',
    'can_edit_admin_mentors_review',
    'can_approve_admin_mentors_review',
    'can_view_admin_resumes',
  ]

  return permissions.includes(value as Permission)
}
