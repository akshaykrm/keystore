package membership

const (
	RoleOwner  string = "owner"
	RoleMember string = "member"
)

//	func IsValid() bool {
//		switch r {
//		case RoleOwner, RoleMember:
//			return true
//		default:
//			return false
//		}
//	}
type CreateMembershipInput struct {
	WorkspaceID int64
	UserID      int64
	Role        string
}
