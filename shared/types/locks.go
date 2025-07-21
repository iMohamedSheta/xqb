package types

type LockMode int

const (
	LockForUpdate LockMode = iota
	LockInShare
	LockNoKeyUpdate
	LockKeyShare
)

type LockWait int

const (
	LockNoWait LockWait = iota
	LockSkipLocked
)
