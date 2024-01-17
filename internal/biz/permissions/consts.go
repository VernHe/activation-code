package permissions

const (
	CREATE         = "CREATE"
	QUERY          = "QUERY"
	DELETE         = "DELETE"
	UPDATE         = "UPDATE"
	UPDATE_UNUSED  = "UPDATE_UNUSED"
	UPDATE_USED    = "UPDATE_USED"
	UPDATE_LOCKED  = "UPDATE_LOCKED"
	UPDATE_DELETED = "UPDATE_DELETED"
)

var (
	AllAllowedPernisions = []string{
		CREATE,
		QUERY,
		DELETE,
		UPDATE,
		UPDATE_UNUSED,
		UPDATE_USED,
		UPDATE_LOCKED,
		UPDATE_DELETED,
	}
)
