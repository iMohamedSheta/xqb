package xqb

import "github.com/iMohamedSheta/xqb/shared/types"

// LockForUpdate - adds FOR UPDATE to lock selected rows for writing (exclusive lock).
func (qb *QueryBuilder) LockForUpdate() *QueryBuilder {
	qb.SetOption(types.OptionLock, types.LockForUpdate)
	return qb
}

// SharedLock - adds LOCK IN SHARE MODE (MySql) or FOR SHARE (PostgreSql) to lock rows for reading.
func (qb *QueryBuilder) SharedLock() *QueryBuilder {
	qb.SetOption(types.OptionLock, types.LockInShare)
	return qb
}

// LockNoKeyUpdate - adds FOR NO KEY UPDATE (PostgreSql) to prevent updates to the row's primary/unique key.
func (qb *QueryBuilder) LockNoKeyUpdate() *QueryBuilder {
	qb.SetOption(types.OptionLock, types.LockNoKeyUpdate)
	return qb
}

// LockKeyShare - adds FOR KEY SHARE (PostgreSql) to allow reading but prevents updates/deletes to referenced keys.
func (qb *QueryBuilder) LockKeyShare() *QueryBuilder {
	qb.SetOption(types.OptionLock, types.LockKeyShare)
	return qb
}

// LockNoWait - appends NOWAIT to prevent waiting on locked rows.
func (qb *QueryBuilder) NoWaitLocked() *QueryBuilder {
	qb.SetOption(types.OptionLockWait, types.LockNoWait)
	return qb
}

// LockSkipLocked - appends SKIP LOCKED to skip over locked rows.
func (qb *QueryBuilder) SkipLocked() *QueryBuilder {
	qb.SetOption(types.OptionLockWait, types.LockSkipLocked)
	return qb
}
