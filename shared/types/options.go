package types

type Option string

const (
	OptionLock              Option = "lock"                // expect OptionValueLock
	OptionLockWait          Option = "lock_wait"           // expect OptionValueLockWait
	OptionIsUpsert          Option = "is_upsert"           // expect bool
	OptionUpsertUniqueBy    Option = "upsert_unique_by"    // expect string (column name)
	OptionUpsertUpdatedCols Option = "upsert_updated_cols" // expect []string
	OptionReturningId       Option = "returning_id"        // expect bool (implicitly true if present)
)

func (o Option) String() string {
	return string(o)
}

// Option lock values
type OptionValueLock int

const (
	LockForUpdate OptionValueLock = iota
	LockInShare
	LockNoKeyUpdate
	LockKeyShare
)

// Option lock wait values
type OptionValueLockWait int

const (
	LockNoWait OptionValueLockWait = iota
	LockSkipLocked
)
