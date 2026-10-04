package membership

type Role string

const (
	RoleOwner  Role = "owner"
	RoleMember Role = "member"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleMember:
		return true
	default:
		return false
	}
}

type CreateMembershipInput struct {
	WorkspaceID string
	UserID      string
	Role        Role
}
