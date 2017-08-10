package db

func InitIndicies() {
	EnsureAssetIndex()
	EnsureStoreIndex()
	EnsureCategoryIndex()
	EnsureProductIndex()
	EnsureUserIndex()
	EnsureCartIndex()
	EnsureAddressIndex()
}
