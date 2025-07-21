package types

type Option string

const (
	OptionLock     Option = "lock"
	OptionLockWait Option = "lock_wait"
)

func (o Option) String() string {
	return string(o)
}
