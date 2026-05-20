package api

import "regexp"

// AnonUUIDRegexp validates client anonymous id (UUID canonical string form).
var AnonUUIDRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidAnonUUID(s string) bool {
	return s != "" && AnonUUIDRegexp.MatchString(s)
}
