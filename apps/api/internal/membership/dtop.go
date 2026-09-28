package membership

type CreateMembershipInput struct {
	WorkspaceID string
	UserID      string
	Role        Role
}
