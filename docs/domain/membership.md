## Membership

---

A membership represents a user's relationship with a workspace.

It allows a user to belong to multiple workspaces and a workspace to
have multiple members.

There are currently two roles a membership can have:

1. Owner
2. Member

Owner is the one who created the workspace or is given the role of owner
default to the creator
Member is the one who have access to the workspaces

Based on this currently membership only needs id, role, user_id, workspace_id

---

### Future

- When adding a user to a workspace a mail will be sent to the users
  mail id saying to join

### Open Questions

- How to add other roles & permissions?
- What happens if an owner adds a user who is not in the system?
