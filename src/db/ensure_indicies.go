package db

func InitIndicies() {
	EnsureOrderIndex()
	EnsureReviewIndex()
	EnsureAssetIndex()
	EnsureStoreIndex()
	EnsureCategoryIndex()
	EnsureProductIndex()
	EnsureUserIndex()
	EnsureCartIndex()
	EnsureAddressIndex()
}
