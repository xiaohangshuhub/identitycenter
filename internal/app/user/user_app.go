package user

type UserApp struct {
	LoginHandler  *LoginHandler
	LogoutHandler *LogoutHandler
}

func NewUserApp(loginHandler *LoginHandler, logoutHandler *LogoutHandler) *UserApp {
	return &UserApp{
		LoginHandler:  loginHandler,
		LogoutHandler: logoutHandler,
	}
}
