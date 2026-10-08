package constants

// TrustedHostingPlatforms lists well-known legitimate static/app hosting services
// whose subdomains appear in the PSL private section (returning icann=false from
// publicsuffix.PublicSuffix). These are operated by reputable companies and must
// not be penalized as "unregulated" TLDs. Subdomains on these platforms also
// lack per-subdomain NS records (the platform manages the whole zone) and are
// inherently unranked, so rank-0 and missing-NS penalties are suppressed too.
// An entry covers the suffixes under it: amazonaws.com covers
// s3.us-east-1.amazonaws.com, which the PSL lists separately per region.
var TrustedHostingPlatforms = map[string]struct{}{
	"github.io":           {}, // GitHub Pages
	"github.dev":          {}, // GitHub dev environments
	"gitlab.io":           {}, // GitLab Pages
	"netlify.app":         {}, // Netlify
	"vercel.app":          {}, // Vercel
	"pages.dev":           {}, // Cloudflare Pages
	"workers.dev":         {}, // Cloudflare Workers
	"web.app":             {}, // Firebase Hosting
	"firebaseapp.com":     {}, // Firebase Hosting (legacy)
	"surge.sh":            {}, // Surge.sh
	"fly.dev":             {}, // Fly.io
	"onrender.com":        {}, // Render
	"herokuapp.com":       {}, // Heroku
	"glitch.me":           {}, // Glitch
	"replit.app":          {}, // Replit
	"bitbucket.io":        {}, // Bitbucket Pages
	"azurewebsites.net":   {}, // Azure App Service
	"azurestaticapps.net": {}, // Azure Static Web Apps
	"amplifyapp.com":      {}, // AWS Amplify
	"readthedocs.io":      {}, // ReadTheDocs
	"huggingface.co":      {}, // Hugging Face Spaces (has PSL entry)
	"streamlit.app":       {}, // Streamlit Cloud
	"railway.app":         {}, // Railway
	"koyeb.app":           {}, // Koyeb
	"cyclic.app":          {}, // Cyclic
	"deno.dev":            {}, // Deno Deploy
	"val.run":             {}, // Val Town

	// File and asset hosting run by large providers. Content is uploaded by
	// anyone, so these lift the "unregulated extension" and "unranked"
	// penalties without lending the provider's own reputation.
	"googleapis.com":         {}, // Google Cloud Storage (storage.googleapis.com)
	"appspot.com":            {}, // Google App Engine
	"blogspot.com":           {}, // Blogger
	"githubusercontent.com":  {}, // GitHub raw files and release assets
	"amazonaws.com":          {}, // Amazon S3 in every region, EC2 hostnames
	"core.windows.net":       {}, // Azure Storage: blob., web. (static sites), file.
	"cloudfront.net":         {}, // Amazon CloudFront
	"r2.dev":                 {}, // Cloudflare R2
	"digitaloceanspaces.com": {}, // DigitalOcean Spaces, per region
	"linodeobjects.com":      {}, // Akamai (Linode) Object Storage, per region
}
