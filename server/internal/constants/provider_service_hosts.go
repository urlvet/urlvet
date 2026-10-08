package constants

// ProviderServiceHosts are hosts on a hosting suffix that the provider itself
// runs, mapped to the service's name. The Public Suffix List makes
// storage.googleapis.com look like someone's own site, the same as
// octocat.github.io; these are the provider's endpoints, while the files
// under their paths are uploaded by anyone.
var ProviderServiceHosts = map[string]string{
	"storage.googleapis.com":               "Google Cloud Storage",
	"firebasestorage.googleapis.com":       "Firebase Storage",
	"raw.githubusercontent.com":            "GitHub",
	"objects.githubusercontent.com":        "GitHub",
	"release-assets.githubusercontent.com": "GitHub",
}
