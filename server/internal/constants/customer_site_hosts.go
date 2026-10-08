package constants

// CustomerSiteHosts are site builders and hosts that give each customer a
// subdomain (someone.weebly.com) but aren't on the Public Suffix List, so by
// default those subdomains would inherit the provider's rank and age. A
// subdomain here is treated as its own site, like someone.github.io; the
// provider's own site (weebly.com, www.weebly.com) is unaffected. Only
// domains a provider uses for customers belong here: wordpress.com also hosts
// WordPress's own subdomains, which would lose their rank.
var CustomerSiteHosts = map[string]struct{}{
	// Website builders
	"weebly.com":          {},
	"weeblysite.com":      {},
	"godaddysites.com":    {},
	"jimdofree.com":       {},
	"jimdosite.com":       {},
	"wixsite.com":         {},
	"wixstudio.io":        {},
	"editorx.io":          {},
	"square.site":         {},
	"mystrikingly.com":    {},
	"yolasite.com":        {},
	"site123.me":          {},
	"webnode.page":        {},
	"tilda.ws":            {},
	"craftum.io":          {},
	"carrd.co":            {},
	"framer.website":      {},
	"framer.app":          {},
	"webflow.io":          {},
	"mobirisesite.com":    {},
	"ukit.me":             {},
	"hostingersite.com":   {},
	"zohosites.com":       {},
	"business.site":       {}, // Google Business Profile sites
	"durable.co":          {},
	"typedream.app":       {},
	"notion.site":         {},
	"gitbook.io":          {},
	"start.page":          {},
	"ucoz.net":            {},
	"ucoz.ru":             {},
	"narod.ru":            {},
	"dothome.co.kr":       {},
	"000webhostapp.com":   {},
	"wuaze.com":           {}, // InfinityFree
	"rf.gd":               {}, // InfinityFree
	"epizy.com":           {}, // InfinityFree
	"free.nf":             {}, // InfinityFree
	"great-site.net":      {}, // InfinityFree
	"infinityfreeapp.com": {}, // InfinityFree
	"lovestoblog.com":     {}, // InfinityFree
	"kesug.com":           {}, // InfinityFree
	"42web.io":            {}, // InfinityFree
	"ct.ws":               {}, // InfinityFree
	"zya.me":              {}, // InfinityFree
	"xo.je":               {}, // InfinityFree
	"page.gd":             {}, // InfinityFree
	"gt.tc":               {}, // InfinityFree
	"unaux.com":           {}, // iFastNet free hosting
	"liveblog365.com":     {}, // iFastNet free hosting
	"iceiy.com":           {}, // iFastNet free hosting
	"squarespace.com":     {}, // its own subdomains: ProviderSubdomains
	"systeme.io":          {}, // its own subdomains: ProviderSubdomains
	"flazio.com":          {}, // its own subdomains: ProviderSubdomains
	"zapier.app":          {}, // Zapier Interfaces
	"4everland.app":       {}, // 4EVERLAND (IPFS) hosting
	"itempurl.com":        {}, // SmarterASP.NET free sites
	"infy.click":          {},
	"freewebhostmost.com": {},

	// Hosting and CDN endpoints named per customer
	"tw1.ru":            {}, // Timeweb
	"b-cdn.net":         {}, // BunnyCDN
	"edgeone.dev":       {}, // Tencent EdgeOne Pages
	"edgeone.app":       {},
	"ngrok-free.app":    {},
	"ngrok.io":          {},
	"ngrok.app":         {},
	"trycloudflare.com": {},
	"loca.lt":           {},
	"serveo.net":        {},

	// Managed WordPress and VPS hosts that name customer sites and servers
	"myftpupload.com":         {}, // GoDaddy WordPress staging
	"mybluehost.me":           {}, // Bluehost temporary domains
	"wpcomstaging.com":        {}, // WordPress.com staging
	"wpengine.com":            {}, // WP Engine install names
	"kinsta.cloud":            {},
	"pantheonsite.io":         {},
	"cloudclusters.net":       {},
	"contaboserver.net":       {}, // Contabo VPS hostnames
	"your-server.de":          {}, // Hetzner server hostnames
	"webador.com":             {}, // Webador site builder
	"b12sites.com":            {}, // B12 site builder
	"previewship.net":         {}, // preview hosting
	"mytemp.website":          {}, // Hostinger temporary domains
	"hstgr.cloud":             {}, // Hostinger VPS hostnames
	"meusitehostgator.com.br": {}, // HostGator Brazil temporary domains

	// Object storage that isn't on the Public Suffix List. Each bucket, or
	// each file host (f005.backblazeb2.com/file/<bucket>/…), is its own site
	// instead of borrowing the provider's rank.
	"backblazeb2.com": {}, // Backblaze B2
	"wasabisys.com":   {}, // Wasabi
	"aliyuncs.com":    {}, // Alibaba Cloud OSS
	"myqcloud.com":    {}, // Tencent Cloud COS
	"filebase.com":    {}, // Filebase

	// Found by matching two years of phishing feeds (OpenPhish, PhishTank,
	// URLhaus) against the top million: providers whose customers' subdomains
	// kept turning up. Website builders and free or managed hosting:
	"webcindario.com":           {},
	"ubpages.com":               {},
	"serv00.net":                {},
	"webwave.dev":               {},
	"codeanyapp.com":            {},
	"myportfolio.com":           {},
	"temporary.site":            {},
	"w3spaces.com":              {},
	"ghost.io":                  {},
	"studio.site":               {},
	"m-pages.com":               {},
	"fwh.is":                    {},
	"teemill.com":               {},
	"hstn.me":                   {},
	"sviluppo.host":             {},
	"boxmode.io":                {},
	"mydomain.zone":             {},
	"myclickfunnels.com":        {},
	"wcomhost.com":              {},
	"tiiny.io":                  {},
	"tiiny.site":                {},
	"work.gd":                   {},
	"netsons.org":               {},
	"mmm.page":                  {},
	"cloudwaysapps.com":         {},
	"swtest.ru":                 {},
	"free.hr":                   {},
	"wstd.io":                   {},
	"univer.se":                 {},
	"rollout.site":              {},
	"n0c.world":                 {},
	"created.app":               {},
	"closte.com":                {},
	"pory.app":                  {},
	"mystagingwebsite.com":      {},
	"cloudworkstations.dev":     {},
	"mvt.so":                    {},
	"twil.io":                   {},
	"sg-host.com":               {},
	"ispot.cc":                  {},
	"nxcli.io":                  {},
	"xsph.ru":                   {},
	"webmo.fr":                  {},
	"komi.io":                   {},
	"netsolhost.com":            {},
	"onepage.me":                {},
	"x10.mx":                    {},
	"peraichi.com":              {},
	"neocities.org":             {},
	"servername.online":         {},
	"kinsta.page":               {},
	"durablesites.com":          {},
	"softr.app":                 {},
	"myfreesites.net":           {},
	"glide.page":                {},
	"run.place":                 {},
	"on.fleek.co":               {},
	"bitrix24.site":             {},
	"gamer.gd":                  {},
	"site44.com":                {},
	"atwebpages.com":            {},
	"fo.team":                   {},
	"hosting-test.net":          {},
	"hyperphp.com":              {},
	"builder-preview.com":       {},
	"misitiohostgator.com":      {},
	"mssg.me":                   {},
	"strato-hosting.eu":         {},
	"whf.bz":                    {},
	"kleap.co":                  {},
	"instawp.xyz":               {},
	"bravesites.com":            {},
	"mybluehostin.me":           {},
	"ulcraft.com":               {},
	"moy.su":                    {},
	"dream.press":               {},
	"hop.ru":                    {},
	"flywheelsites.com":         {},
	"thsite.top":                {},
	"serveousercontent.com":     {},
	"company.site":              {},
	"beget.tech":                {},
	"webstarterz.com":           {},
	"vh.net.pl":                 {},
	"easy-hebergement.net":      {},
	"conohawing.com":            {},
	"ac-page.com":               {},
	"vzy.io":                    {},
	"onamaeweb.jp":              {},
	"infy.uk":                   {},
	"brizy.site":                {},
	"1wp.site":                  {},
	"cleansite.info":            {},
	"swipepages.net":            {},
	"myhostpoint.ch":            {},
	"ik-server.com":             {},
	"hidora.net":                {},
	"flazio.site":               {},
	"onepage.website":           {},
	"infinityfree.me":           {},
	"000.pe":                    {},
	"dora.run":                  {},
	"somee.com":                 {},
	"pktriot.net":               {},
	"2bd.net":                   {},
	"webydo.com":                {},
	"site.je":                   {},
	"builderallwppro.com":       {},
	"finance.blog":              {},
	"tripod.com":                {},
	"line.pm":                   {},
	"justns.ru":                 {},
	"sitebeat.crazydomains.com": {},
	"hosted.phplist.com":        {},

	// Hosting platforms that aren't on the Public Suffix List either:
	"glitch.me":    {},
	"surge.sh":     {},
	"railway.app":  {},
	"koyeb.app":    {},
	"cyclic.app":   {},
	"github.dev":   {}, // Codespaces
	"filesusr.com": {}, // Wix user files
	"cabanova.com": {},

	// IPFS, Arweave and object storage gateways (a subdomain per upload):
	"ionoscloud.com": {},
	"storage.dev":    {},
	"arweave.net":    {},
	"infura-ipfs.io": {},
	"eth.limo":       {},

	// Dynamic DNS, where anyone picks a free subdomain:
	"crabdance.com":     {},
	"publicvm.com":      {},
	"dns.army":          {},
	"ns01.info":         {},
	"v6.rocks":          {},
	"mooo.com":          {},
	"linkpc.net":        {},
	"ignorelist.com":    {},
	"ydns.eu":           {},
	"strangled.net":     {},
	"myddns.me":         {},
	"ddnss.eu":          {},
	"odns.fr":           {},
	"jumpingcrab.com":   {},
	"chickenkiller.com": {},
	"2mydns.net":        {},
	"tftpd.net":         {},

	// Names that stand for an IP address (1-2-3-4.sslip.io is 1.2.3.4):
	// IPHostnameServices. Each one is its own server.
	"sslip.io":              {},
	"nip.io":                {},
	"xip.io":                {},
	"traefik.me":            {},
	"nip.direct":            {},
	"cloud-xip.com":         {},
	"host.secureserver.net": {}, // GoDaddy server names
}

// IPHostnameServices resolve any name with an IP address in it to that
// address, so a link on one is a link to a bare server.
var IPHostnameServices = map[string]struct{}{
	"sslip.io":   {},
	"nip.io":     {},
	"xip.io":     {},
	"traefik.me": {},
	"nip.direct": {},
	// Server names hosting panels give every machine from its IP
	// (login.50-6-193-131.cprapid.com).
	"cprapid.com":           {}, // cPanel
	"cpanel.site":           {}, // cPanel
	"plesk.page":            {}, // Plesk
	"cloud-xip.com":         {},
	"host.secureserver.net": {}, // GoDaddy
}

// blogspotCountries are Blogger's country addresses (someone.blogspot.de).
// They left the Public Suffix List, so they're listed here like blogspot.com.
var blogspotCountries = []string{
	"ae", "al", "am", "ba", "be", "bg", "ca", "ch", "cl", "co", "co.at", "co.id",
	"co.il", "co.ke", "com.ar", "com.au", "com.br", "com.by", "com.co", "com.cy", "com.ee", "com.eg", "com.es", "com.mt",
	"com.ng", "com.tr", "com.uy", "co.nz", "co.uk", "co.za", "cz", "de", "dk", "es", "fi", "fr",
	"gr", "hk", "hr", "hu", "ie", "in", "is", "it", "jp", "kr", "li", "lt",
	"lu", "md", "mk", "mx", "my", "nl", "no", "pe", "pt", "qa", "ro", "rs",
	"ru", "se", "sg", "si", "sk", "sn", "tw", "ug",
}

func init() {
	google := HighValueBrands["Google"]
	for _, c := range blogspotCountries {
		CustomerSiteHosts["blogspot."+c] = struct{}{}
		google.Platforms = append(google.Platforms, "blogspot."+c)
	}
	HighValueBrands["Google"] = google
}

// ProviderSubdomains are subdomains a CustomerSiteHosts provider keeps for
// itself (account.squarespace.com). They, and anything under them, stay the
// provider's site and keep its rank.
var ProviderSubdomains = map[string]map[string]struct{}{
	"systeme.io": {
		"help": {}, "api": {}, "app": {}, "blog": {}, "status": {},
		"affiliate": {}, "community": {}, "academy": {}, "developer": {},
	},
	"flazio.com": {
		"app": {}, "help": {}, "blog": {}, "api": {},
	},
	"squarespace.com": {
		"account": {}, "login": {}, "support": {}, "help": {}, "domains": {},
		"status": {}, "engineering": {}, "forum": {}, "developers": {},
		"api": {}, "static": {}, "static1": {}, "images": {}, "assets": {},
		"designer": {}, "commerce": {}, "pages": {}, "investors": {},
		"newsroom": {}, "careers": {}, "config": {}, "mail": {},
	},
}
