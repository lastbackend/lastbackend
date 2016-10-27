package storage

type Storage struct {
	User    IUserService
	Profile IProfileService
	// Other services should be here
}
