package requestctx

// todo: move to service/userctx
type UserRequestInfo struct {
	IsLoggedIn bool
	Id         string
	Email      string
	IsAdmin    bool
}
