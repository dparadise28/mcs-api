package db

func InitIndicies() {
	EnsureStoreIndex()
	EnsureCategoryIndex()
	EnsureProductIndex()
	EnsureUserIndex()
	EnsureCartIndex()
}
