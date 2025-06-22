package errors

type DatabaseError struct{}

func (e *DatabaseError) Error() string {
	return "Database error"
}
