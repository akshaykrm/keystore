package membership

const (
	RoleOwner  string = "owner"
	RoleMember string = "member"
)

type CreateMembershipInput struct {
	WorkspaceID int64
	UserID      int64
	Role        string
}
