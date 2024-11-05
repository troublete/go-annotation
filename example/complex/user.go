package complex

import "time"

// crud.model{name="users"}
type User struct {
	// crud.field{name="id"}
	ID *string
	// crud.field{name="firstname"}
	Firstname *string
	// crud.field{name="lastname"}
	Lastname *string
	// crud.field{name="created_at"}
	CreatedAt *time.Time
	// crud.field{name="deleted_at"}
	UpdatedAt *time.Time
}
