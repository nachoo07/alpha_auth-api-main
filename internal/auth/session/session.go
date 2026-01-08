package session

type AccessTokenInfo struct {
	ID       uint
	Rol      uint
	Username string
	Hash     string
	Exp      int64
}

type RefreshTokenInfo struct {
	ID   uint
	Hash string
}
