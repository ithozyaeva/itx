package s3resolve

// Resolver is a function that resolves S3 keys to presigned URLs.
// Set by utils.InitGlobalS3() at startup to break the import cycle.
var Resolver func(stored string) string

// ResolveS3URL resolves a stored S3 key or old-style URL to a presigned URL.
// Returns the input unchanged if no resolver is registered.
func ResolveS3URL(stored string) string {
	if Resolver == nil || stored == "" {
		return stored
	}
	return Resolver(stored)
}
