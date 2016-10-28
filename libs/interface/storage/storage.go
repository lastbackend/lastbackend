package storage

type Storage struct {
	User    IUserService
	Account IAccountService
	// Other services should be here
}
