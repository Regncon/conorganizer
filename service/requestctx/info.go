package requestctx

type UserRequestInfo struct {
	IsLoggedIn bool
	Id         string
	Email      string
	IsAdmin    bool
}
