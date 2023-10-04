package keeper

type Keeper interface {
	Get(s string) string
	Set(key string, value string)
}
