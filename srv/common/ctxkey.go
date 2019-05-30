package common

// RepoKey for Context Value
type RepoKey struct {
	db RepoKeyT
}

// RepoKeyT for RepoKey
type RepoKeyT int

// CtxKey(s)
const (
	CKMongoRepo RepoKeyT = iota
	CKGormRepo
)
