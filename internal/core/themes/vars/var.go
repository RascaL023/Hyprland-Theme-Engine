package vars

type VarSource interface {
    Get(key string) (string, bool)
}

