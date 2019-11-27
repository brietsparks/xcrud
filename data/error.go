package data

type Error struct {
	Err error
	Msg string
}

// NewError creates an Error that hides database error details behind Unwrap
func NewError(err error) error {
	if err == nil {
		return nil
	}

	if includes(storeMessages, err.Error()) {
		return err
	}

	msg := dbMessages[err.Error()]

	if msg == "" {
		msg = ErrUnknown
	}

	return &Error{err, msg}
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) Unwrap() error {
	return e.Err
}

// error messages that originate from the data store layer that do not contain sensitive database implementation details
const ErrResourceDNE = "resource does not exist"
var storeMessages = []string{
	ErrResourceDNE,
}

// error messages that originate from the database and contain potentially sensitive database implementation details
const DbErrGroupUserAlreadyLinked = "pq: duplicate key value violates unique constraint \"groups_users_pkey\""
const ErrGroupUserAlreadyLinked = "group already linked to user"
const DbErrGroupOrUserDNE = "pq: insert or update on table \"groups_users\" violates foreign key constraint \"groups_users_group_id_fkey\""
const ErrGroupOrUserDNE = "group or user does not exist"
var dbMessages = map[string]string{
	ErrResourceDNE: ErrResourceDNE,
	DbErrGroupOrUserDNE: ErrGroupOrUserDNE,
	DbErrGroupUserAlreadyLinked: ErrGroupUserAlreadyLinked,
}

// fallthrough error message
const ErrUnknown = "unspecified database error"
