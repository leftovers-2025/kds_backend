package entity

const (
	ROLE_UNKNOWN Role = ""
	ROLE_STUDENT Role = "student"
	ROLE_TEACHER Role = "teacher"
	ROLE_ROOT    Role = "root"
)

type Role string

func RoleFromString(role string) Role {
	switch role {
	case ROLE_STUDENT.String():
		return ROLE_STUDENT
	case ROLE_TEACHER.String():
		return ROLE_TEACHER
	case ROLE_ROOT.String():
		return ROLE_ROOT
	default:
		return ROLE_UNKNOWN
	}
}

func (r Role) String() string {
	return string(r)
}
