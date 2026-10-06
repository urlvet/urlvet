package constants

type BrandEntry struct {
	// TitleKeywords mean the page claims to be the brand wherever they appear.
	TitleKeywords []string
	// Names are the bare brand names. Small sites use them in ordinary titles
	// (resellers, fan sites, news), so they only count on a page that asks for
	// a login or payment, or on a free hosting subdomain.
	Names           []string
	OfficialDomains []string
	// OwnNames are for groups whose members each run their own domain
	// (sparkasse-hannover.de, vrbank-xyz.de): a domain containing one counts
	// as the brand's own.
	OwnNames []string
	// ExactTitles are whole page titles the brand's own login uses, too
	// generic to match inside a longer title ("Sign in to your account" alone
	// is Microsoft's; "Sign in to your account | Acme" is anyone's).
	ExactTitles []string
	// AlsoOfficial are more domains the brand runs (research.google,
	// primevideo.com). They're the brand's own in every check, but unlike
	// OfficialDomains their names aren't looked for inside other domains, so
	// "research" or "messenger" don't become brand names.
	AlsoOfficial []string
	// Platforms are hosting the brand runs for anyone (blogspot.com). A page
	// there isn't the brand's own, but a redirect from one to the brand's own
	// domain is the brand moving its content, not a jump to someone else's.
	Platforms []string
}

// Domains returns every domain the brand runs itself.
func (b BrandEntry) Domains() []string {
	return append(append([]string{}, b.OfficialDomains...), b.AlsoOfficial...)
}

var HighValueBrands = map[string]BrandEntry{

	// ── Email Platforms ──────────────────────────────────────────────────────
	"Gmail": {
		TitleKeywords:   []string{"gmail", "google mail"},
		Names:           []string{"gmail"},
		OfficialDomains: []string{"gmail.com", "google.com"},
	},
	"Google": {
		TitleKeywords:   []string{"google account", "google drive", "google workspace", "google sign in"},
		Names:           []string{"google"},
		OfficialDomains: []string{"google.com", "gmail.com"},
		AlsoOfficial: []string{
			"youtube.com", "youtu.be", "googlemail.com", "blogger.com", "googleblog.com",
			"blog.google", "research.google", "about.google", "store.google", "withgoogle.com",
			"google.dev", "android.com", "chrome.com", "chromium.org", "gstatic.com", "g.co",
			"googlesource.com", "googleplex.com", "goo.gle",
		},
		Platforms: []string{"blogspot.com", "appspot.com", "web.app", "firebaseapp.com", "googleusercontent.com", "googleapis.com"},
	},
	"Microsoft Outlook": {
		TitleKeywords:   []string{"outlook", "hotmail", "live mail", "microsoft mail"},
		Names:           []string{"outlook"},
		OfficialDomains: []string{"outlook.com", "hotmail.com", "live.com", "microsoft.com", "microsoftonline.com"},
	},
	"Microsoft": {
		TitleKeywords:   []string{"microsoft", "office 365", "microsoft 365", "onedrive", "sharepoint", "microsoft teams", "azure"},
		Names:           []string{"microsoft"},
		ExactTitles:     []string{"sign in to your account"},
		OfficialDomains: []string{"microsoft.com", "live.com", "outlook.com", "hotmail.com", "office.com", "office365.com", "microsoftonline.com", "azure.com", "onedrive.com", "sharepoint.com"},
	},
	"Yahoo Mail": {
		TitleKeywords:   []string{"yahoo", "yahoo mail", "yahoo sign in"},
		Names:           []string{"yahoo"},
		OfficialDomains: []string{"yahoo.com", "mail.yahoo.com", "login.yahoo.com"},
	},
	"ProtonMail": {
		TitleKeywords:   []string{"protonmail", "proton mail", "proton account"},
		OfficialDomains: []string{"proton.me", "protonmail.com"},
	},
	"Zoho Mail": {
		TitleKeywords:   []string{"zoho mail", "zoho accounts"},
		OfficialDomains: []string{"zoho.com", "zohomail.com"},
	},
	"AOL Mail": {
		TitleKeywords:   []string{"aol mail", "aol sign in"},
		OfficialDomains: []string{"aol.com", "mail.aol.com"},
	},
	"iCloud Mail": {
		TitleKeywords:   []string{"icloud mail", "apple id", "icloud account", "apple account"},
		OfficialDomains: []string{"icloud.com", "apple.com", "itunes.com"},
	},

	// ── Social Media ─────────────────────────────────────────────────────────
	"Facebook": {
		TitleKeywords:   []string{"facebook login", "facebook sign in"},
		Names:           []string{"facebook"},
		OfficialDomains: []string{"facebook.com", "fb.com"},
		AlsoOfficial:    []string{"fb.me", "facebookmail.com", "messenger.com", "m.me", "meta.com", "fbcdn.net"},
	},
	"Instagram": {
		TitleKeywords:   []string{"instagram login"},
		Names:           []string{"instagram"},
		OfficialDomains: []string{"instagram.com"},
		AlsoOfficial:    []string{"instagr.am", "cdninstagram.com", "threads.net", "threads.com", "meta.com"},
	},
	"Twitter": {
		TitleKeywords:   []string{"x login", "x sign in", "sign in to x"},
		Names:           []string{"twitter"},
		OfficialDomains: []string{"twitter.com", "x.com"},
	},
	"LinkedIn": {
		TitleKeywords:   []string{"linkedin sign in", "linkedin login"},
		Names:           []string{"linkedin"},
		OfficialDomains: []string{"linkedin.com", "lnkd.in"},
	},
	"TikTok": {
		TitleKeywords:   []string{"tiktok login"},
		Names:           []string{"tiktok"},
		OfficialDomains: []string{"tiktok.com"},
	},
	"Snapchat": {
		TitleKeywords:   []string{"snapchat", "snapchat login"},
		Names:           []string{"snapchat"},
		OfficialDomains: []string{"snapchat.com"},
	},
	"Pinterest": {
		TitleKeywords:   []string{"pinterest", "pinterest login"},
		OfficialDomains: []string{"pinterest.com"},
	},
	"Reddit": {
		TitleKeywords:   []string{"reddit", "reddit login", "reddit sign in"},
		OfficialDomains: []string{"reddit.com", "redd.it"},
	},
	"YouTube": {
		TitleKeywords:   []string{"youtube account", "youtube sign in"},
		OfficialDomains: []string{"youtube.com", "youtu.be"},
		AlsoOfficial:    []string{"youtube-nocookie.com", "ytimg.com", "google.com"},
	},
	"Tumblr": {
		TitleKeywords:   []string{"tumblr", "tumblr login"},
		OfficialDomains: []string{"tumblr.com"},
	},

	// ── Messaging & Chat ──────────────────────────────────────────────────────
	"WhatsApp": {
		TitleKeywords:   []string{"whatsapp login", "whatsapp web"},
		Names:           []string{"whatsapp"},
		OfficialDomains: []string{"whatsapp.com", "web.whatsapp.com"},
		AlsoOfficial:    []string{"wa.me", "whatsapp.net", "meta.com"},
	},
	"Telegram": {
		TitleKeywords:   []string{"telegram login", "telegram web"},
		Names:           []string{"telegram"},
		OfficialDomains: []string{"telegram.org", "web.telegram.org", "t.me"},
	},
	"Signal": {
		TitleKeywords:   []string{"signal messenger", "signal login", "signal log in", "signal sign in", "signal account", "signal verification"},
		Names:           []string{"signal"},
		OfficialDomains: []string{"signal.org"},
	},
	"Discord": {
		TitleKeywords:   []string{"discord login", "discord sign in"},
		Names:           []string{"discord"},
		OfficialDomains: []string{"discord.com", "discordapp.com"},
	},
	"Messenger": {
		TitleKeywords:   []string{"facebook messenger", "messenger login", "messenger log in"},
		Names:           []string{"messenger"},
		OfficialDomains: []string{"messenger.com", "facebook.com"},
	},
	"Skype": {
		TitleKeywords:   []string{"skype", "skype sign in"},
		OfficialDomains: []string{"skype.com", "microsoft.com"},
	},
	"Slack": {
		TitleKeywords:   []string{"slack sign in", "slack login", "slack log in", "slack account", "slack verification"},
		Names:           []string{"slack"},
		OfficialDomains: []string{"slack.com"},
	},
	"Viber": {
		TitleKeywords:   []string{"viber", "viber login"},
		OfficialDomains: []string{"viber.com"},
	},
	"Line": {
		TitleKeywords:   []string{"line login", "line account"},
		OfficialDomains: []string{"line.me"},
	},

	// ── E-Commerce ────────────────────────────────────────────────────────────
	"Amazon": {
		// Not the bare name: rainforest tours and river cruises use it too.
		TitleKeywords:   []string{"amazon.com", "amazon sign in", "amazon sign-in", "amazon login", "amazon account", "amazon prime", "amazon pay", "aws sign in"},
		Names:           []string{"amazon"},
		OfficialDomains: []string{"amazon.com", "amazon.co.uk", "amazon.de", "amazon.in", "amazon.ca", "amazon.com.au", "amazon.fr", "amazon.es", "amazon.it", "amazon.co.jp", "amazon.com.br", "amazon.com.mx", "amazon.nl", "amazon.se", "amazon.pl", "amazon.sg", "amazon.ae", "amazon.sa", "amazon.com.tr", "amazon.eg", "amazon.com.be", "aws.amazon.com"},
		AlsoOfficial:    []string{"amazon.cn", "amazon.ie", "amazon.co.za", "primevideo.com", "aboutamazon.com", "amazon.jobs", "audible.com", "a.co", "amzn.to", "media-amazon.com", "ssl-images-amazon.com"},
		Platforms:       []string{"amazonaws.com", "cloudfront.net", "amplifyapp.com"},
	},
	"eBay": {
		TitleKeywords:   []string{"ebay", "ebay sign in", "ebay login"},
		Names:           []string{"ebay"},
		OfficialDomains: []string{"ebay.com", "ebay.co.uk", "ebay.de", "ebay.com.au"},
	},
	"Etsy": {
		TitleKeywords:   []string{"etsy", "etsy sign in"},
		OfficialDomains: []string{"etsy.com"},
	},
	"Flipkart": {
		TitleKeywords:   []string{"flipkart", "flipkart login"},
		OfficialDomains: []string{"flipkart.com"},
	},
	"Walmart": {
		TitleKeywords:   []string{"walmart", "walmart sign in"},
		OfficialDomains: []string{"walmart.com"},
	},
	"AliExpress": {
		TitleKeywords:   []string{"aliexpress", "aliexpress login"},
		OfficialDomains: []string{"aliexpress.com"},
	},
	"Alibaba": {
		TitleKeywords:   []string{"alibaba", "alibaba login"},
		OfficialDomains: []string{"alibaba.com"},
	},
	"Meesho": {
		TitleKeywords:   []string{"meesho", "meesho login"},
		OfficialDomains: []string{"meesho.com"},
	},

	// ── Payments & Digital Wallets ────────────────────────────────────────────
	"PayPal": {
		TitleKeywords:   []string{"paypal", "paypal sign in", "paypal login"},
		Names:           []string{"paypal"},
		OfficialDomains: []string{"paypal.com", "paypal.me"},
	},
	"Google Pay": {
		TitleKeywords:   []string{"google pay", "gpay"},
		OfficialDomains: []string{"pay.google.com", "google.com"},
	},
	"PhonePe": {
		TitleKeywords:   []string{"phonepe", "phone pe"},
		OfficialDomains: []string{"phonepe.com"},
	},
	"Paytm": {
		TitleKeywords:   []string{"paytm", "paytm login"},
		OfficialDomains: []string{"paytm.com"},
	},
	"Stripe": {
		TitleKeywords:   []string{"stripe login", "stripe log in", "stripe sign in", "stripe account", "stripe verification"},
		Names:           []string{"stripe"},
		OfficialDomains: []string{"stripe.com"},
	},
	"CashApp": {
		TitleKeywords:   []string{"cash app", "cashapp"},
		OfficialDomains: []string{"cash.app"},
	},
	"Venmo": {
		TitleKeywords:   []string{"venmo", "venmo login"},
		OfficialDomains: []string{"venmo.com"},
	},
	"Zelle": {
		TitleKeywords:   []string{"zelle", "zelle payment"},
		OfficialDomains: []string{"zellepay.com"},
	},
	"Razorpay": {
		TitleKeywords:   []string{"razorpay"},
		OfficialDomains: []string{"razorpay.com"},
	},
	"Apple Pay": {
		TitleKeywords:   []string{"apple pay", "apple wallet"},
		OfficialDomains: []string{"apple.com"},
	},
	"Wise": {
		TitleKeywords:   []string{"transferwise", "wise login", "wise log in", "wise sign in", "wise account", "wise verification"},
		OfficialDomains: []string{"wise.com"},
	},

	// ── Banking — India ───────────────────────────────────────────────────────
	"SBI": {
		TitleKeywords:   []string{"sbi", "state bank of india", "sbi net banking"},
		OfficialDomains: []string{"onlinesbi.sbi", "sbi.co.in", "sbicard.com", "sbi.bank.in"},
	},
	"HDFC": {
		TitleKeywords:   []string{"hdfc", "hdfc bank", "hdfc net banking"},
		OfficialDomains: []string{"hdfcbank.com", "hdfc.com", "hdfc.bank.in"},
	},
	"ICICI": {
		TitleKeywords:   []string{"icici", "icici bank"},
		OfficialDomains: []string{"icicibank.com", "icici.bank.in"},
	},
	"Axis Bank": {
		TitleKeywords:   []string{"axis bank", "axis net banking"},
		OfficialDomains: []string{"axisbank.com", "axis.bank.in"},
	},
	"Kotak": {
		TitleKeywords:   []string{"kotak", "kotak mahindra"},
		OfficialDomains: []string{"kotak.com", "kotakbank.com"},
	},
	"PNB": {
		TitleKeywords:   []string{"pnb", "punjab national bank"},
		OfficialDomains: []string{"pnbindia.in"},
	},
	"Canara Bank": {
		TitleKeywords:   []string{"canara bank"},
		OfficialDomains: []string{"canarabank.in", "canarabank.com"},
	},
	"Bank of Baroda": {
		TitleKeywords:   []string{"bank of baroda", "bob net banking"},
		OfficialDomains: []string{"bankofbaroda.in"},
	},

	// ── Banking — Global ──────────────────────────────────────────────────────
	"Chase": {
		TitleKeywords:   []string{"chase bank", "chase sign in", "chase login", "chase log in", "chase account", "chase verification"},
		Names:           []string{"chase"},
		OfficialDomains: []string{"chase.com"},
	},
	"Bank of America": {
		TitleKeywords:   []string{"bank of america"},
		Names:           []string{"bank of america"},
		OfficialDomains: []string{"bankofamerica.com"},
	},
	"Wells Fargo": {
		TitleKeywords:   []string{"wells fargo"},
		Names:           []string{"wells fargo"},
		OfficialDomains: []string{"wellsfargo.com"},
	},
	"Citibank": {
		TitleKeywords:   []string{"citibank", "citi bank"},
		OfficialDomains: []string{"citibank.com", "citi.com"},
	},
	"HSBC": {
		TitleKeywords:   []string{"hsbc", "hsbc bank"},
		OfficialDomains: []string{"hsbc.com", "hsbc.co.uk"},
	},
	"Barclays": {
		TitleKeywords:   []string{"barclays"},
		OfficialDomains: []string{"barclays.co.uk", "barclays.com"},
	},
	"NatWest": {
		TitleKeywords:   []string{"natwest"},
		OfficialDomains: []string{"natwest.com"},
	},
	"Santander": {
		// Not the bare name: Santander is also a Spanish city and its council.
		TitleKeywords:   []string{"banco santander", "santander bank", "santander online banking", "santander login", "santander sign in", "santander log in", "santander account", "santander verification", "my santander", "santander particulares", "santander empresas"},
		Names:           []string{"santander"},
		OfficialDomains: []string{"santander.co.uk", "santander.com", "bancosantander.es", "santander.com.br", "santander.com.mx", "santander.cl", "santander.pt", "santander.de", "santander.pl", "santander.com.ar", "santanderbank.com", "openbank.es"},
	},
	"Lloyds": {
		TitleKeywords:   []string{"lloyds", "lloyds bank"},
		OfficialDomains: []string{"lloydsbank.com"},
	},
	"TD Bank": {
		TitleKeywords:   []string{"td bank", "td canada trust"},
		OfficialDomains: []string{"td.com", "tdbank.com"},
	},
	"RBC": {
		TitleKeywords:   []string{"rbc", "royal bank of canada"},
		OfficialDomains: []string{"rbc.com", "rbcroyalbank.com"},
	},
	"Scotiabank": {
		TitleKeywords:   []string{"scotiabank"},
		OfficialDomains: []string{"scotiabank.com"},
	},
	"ANZ": {
		TitleKeywords:   []string{"anz bank", "anz sign in"},
		OfficialDomains: []string{"anz.com", "anz.com.au"},
	},
	"Commonwealth Bank": {
		TitleKeywords:   []string{"commbank", "commonwealth bank", "netbank"},
		OfficialDomains: []string{"commbank.com.au", "netbank.com.au"},
	},
	"Westpac": {
		TitleKeywords:   []string{"westpac"},
		OfficialDomains: []string{"westpac.com.au"},
	},
	"DBS": {
		TitleKeywords:   []string{"dbs bank", "dbs digibank"},
		OfficialDomains: []string{"dbs.com"},
	},
	"Standard Chartered": {
		TitleKeywords:   []string{"standard chartered"},
		OfficialDomains: []string{"sc.com", "standardchartered.com"},
	},
	"UBS": {
		TitleKeywords:   []string{"ubs", "ubs login"},
		OfficialDomains: []string{"ubs.com"},
	},

	// ── Crypto & Web3 ─────────────────────────────────────────────────────────
	"Binance": {
		TitleKeywords:   []string{"binance", "binance login"},
		Names:           []string{"binance"},
		OfficialDomains: []string{"binance.com", "binance.us"},
	},
	"Coinbase": {
		TitleKeywords:   []string{"coinbase", "coinbase login"},
		Names:           []string{"coinbase"},
		OfficialDomains: []string{"coinbase.com"},
	},
	"Kraken": {
		TitleKeywords:   []string{"kraken login", "kraken log in", "kraken sign in", "kraken account", "kraken verification"},
		Names:           []string{"kraken"},
		OfficialDomains: []string{"kraken.com"},
	},
	"WazirX": {
		TitleKeywords:   []string{"wazirx"},
		OfficialDomains: []string{"wazirx.com"},
	},
	"Bybit": {
		TitleKeywords:   []string{"bybit", "bybit login"},
		OfficialDomains: []string{"bybit.com"},
	},
	"OKX": {
		TitleKeywords:   []string{"okx", "okex"},
		OfficialDomains: []string{"okx.com", "okex.com"},
	},
	"Gemini": {
		TitleKeywords:   []string{"gemini crypto", "gemini exchange", "gemini login", "gemini log in", "gemini sign in", "gemini wallet"},
		OfficialDomains: []string{"gemini.com"},
	},
	"MetaMask": {
		TitleKeywords:   []string{"metamask", "metamask wallet"},
		Names:           []string{"metamask"},
		OfficialDomains: []string{"metamask.io"},
	},
	"Phantom": {
		TitleKeywords:   []string{"phantom wallet"},
		Names:           []string{"phantom"},
		OfficialDomains: []string{"phantom.app"},
	},
	"Trust Wallet": {
		TitleKeywords:   []string{"trust wallet"},
		Names:           []string{"trust wallet"},
		OfficialDomains: []string{"trustwallet.com"},
	},
	"Ledger": {
		TitleKeywords:   []string{"ledger live", "ledger login", "ledger log in", "ledger sign in", "ledger account", "ledger verification"},
		Names:           []string{"ledger"},
		OfficialDomains: []string{"ledger.com"},
	},
	"OpenSea": {
		TitleKeywords:   []string{"opensea"},
		OfficialDomains: []string{"opensea.io"},
	},
	"Uniswap": {
		TitleKeywords:   []string{"uniswap"},
		OfficialDomains: []string{"uniswap.org", "app.uniswap.org"},
	},

	// ── Gaming ────────────────────────────────────────────────────────────────
	"Steam": {
		TitleKeywords:   []string{"steam login", "steam sign in", "steam log in", "steam account", "steam verification"},
		Names:           []string{"steam"},
		OfficialDomains: []string{"steampowered.com", "steamcommunity.com", "store.steampowered.com"},
	},
	"Epic Games": {
		TitleKeywords:   []string{"epic games", "fortnite", "epic games launcher"},
		OfficialDomains: []string{"epicgames.com", "fortnite.com"},
	},
	"Riot Games": {
		TitleKeywords:   []string{"riot games", "league of legends", "valorant", "riot account"},
		OfficialDomains: []string{"riotgames.com", "leagueoflegends.com", "playvalorant.com"},
	},
	"PlayStation": {
		TitleKeywords:   []string{"playstation", "psn", "playstation network", "ps5 account"},
		OfficialDomains: []string{"playstation.com", "sonyentertainmentnetwork.com"},
	},
	"Xbox": {
		TitleKeywords:   []string{"xbox", "xbox live", "xbox sign in"},
		OfficialDomains: []string{"xbox.com", "microsoft.com"},
	},
	"Nintendo": {
		TitleKeywords:   []string{"nintendo", "nintendo account", "nintendo sign in"},
		OfficialDomains: []string{"nintendo.com", "accounts.nintendo.com"},
	},
	"Roblox": {
		TitleKeywords:   []string{"roblox", "roblox login"},
		Names:           []string{"roblox"},
		OfficialDomains: []string{"roblox.com"},
	},
	"Blizzard": {
		TitleKeywords:   []string{"battle.net", "world of warcraft", "overwatch", "blizzard login", "blizzard log in", "blizzard sign in", "blizzard account", "blizzard verification"},
		Names:           []string{"blizzard"},
		OfficialDomains: []string{"blizzard.com", "battle.net"},
	},
	"Activision": {
		TitleKeywords:   []string{"activision", "call of duty", "warzone"},
		OfficialDomains: []string{"activision.com", "callofduty.com"},
	},
	"EA": {
		TitleKeywords:   []string{"ea account", "ea sign in", "origin login"},
		OfficialDomains: []string{"ea.com", "origin.com"},
	},

	// ── Streaming & Subscriptions ─────────────────────────────────────────────
	"Netflix": {
		TitleKeywords:   []string{"netflix", "netflix sign in", "netflix login"},
		Names:           []string{"netflix"},
		OfficialDomains: []string{"netflix.com"},
	},
	"Spotify": {
		TitleKeywords:   []string{"spotify", "spotify login", "spotify account"},
		Names:           []string{"spotify"},
		OfficialDomains: []string{"spotify.com", "accounts.spotify.com"},
	},
	"Disney+": {
		TitleKeywords:   []string{"disney+", "disney plus", "disneyplus"},
		OfficialDomains: []string{"disneyplus.com", "disney.com"},
	},
	"Hulu": {
		TitleKeywords:   []string{"hulu", "hulu login"},
		OfficialDomains: []string{"hulu.com"},
	},
	"HBO Max": {
		TitleKeywords:   []string{"hbo max", "max sign in", "hbo login"},
		OfficialDomains: []string{"max.com", "hbomax.com"},
	},
	"Apple TV": {
		TitleKeywords:   []string{"apple tv", "apple tv+"},
		OfficialDomains: []string{"tv.apple.com", "apple.com"},
	},
	"Peacock": {
		TitleKeywords:   []string{"peacock tv", "peacock login", "peacock log in", "peacock sign in", "peacock account", "peacock verification"},
		Names:           []string{"peacock"},
		OfficialDomains: []string{"peacocktv.com"},
	},
	"Paramount+": {
		TitleKeywords:   []string{"paramount+", "paramount plus"},
		OfficialDomains: []string{"paramountplus.com"},
	},
	"Crunchyroll": {
		TitleKeywords:   []string{"crunchyroll", "crunchyroll login"},
		OfficialDomains: []string{"crunchyroll.com"},
	},
	"Amazon Prime Video": {
		TitleKeywords:   []string{"prime video", "amazon video"},
		OfficialDomains: []string{"primevideo.com", "amazon.com"},
	},

	// ── Cloud & Productivity ──────────────────────────────────────────────────
	"Adobe": {
		TitleKeywords:   []string{"adobe", "adobe sign in", "creative cloud"},
		Names:           []string{"adobe"},
		OfficialDomains: []string{"adobe.com", "adobeid.services"},
	},
	"Dropbox": {
		TitleKeywords:   []string{"dropbox", "dropbox sign in"},
		Names:           []string{"dropbox"},
		OfficialDomains: []string{"dropbox.com"},
	},
	"GitHub": {
		// Not the bare name: developer blogs and portfolios mention GitHub.
		TitleKeywords:   []string{"github sign in", "sign in to github", "github login", "github log in", "github account", "github verification"},
		Names:           []string{"github"},
		OfficialDomains: []string{"github.com", "githubusercontent.com"},
	},
	"Zoom": {
		TitleKeywords:   []string{"zoom sign in", "zoom meeting", "zoom login", "zoom log in", "zoom account", "zoom verification"},
		Names:           []string{"zoom"},
		OfficialDomains: []string{"zoom.us"},
	},
	"DocuSign": {
		TitleKeywords:   []string{"docusign"},
		Names:           []string{"docusign"},
		OfficialDomains: []string{"docusign.com", "docusign.net"},
	},
	"Shopify": {
		TitleKeywords:   []string{"shopify", "shopify login"},
		OfficialDomains: []string{"shopify.com", "myshopify.com"},
	},
	"Notion": {
		TitleKeywords:   []string{"notion login", "notion log in", "notion sign in", "notion account", "notion verification"},
		Names:           []string{"notion"},
		OfficialDomains: []string{"notion.so"},
	},
	"Box": {
		TitleKeywords:   []string{"box sign in", "box cloud"},
		OfficialDomains: []string{"box.com"},
	},
	"Atlassian": {
		TitleKeywords:   []string{"atlassian", "jira login", "confluence login", "bitbucket login"},
		OfficialDomains: []string{"atlassian.com", "atlassian.net"},
	},
	"Salesforce": {
		TitleKeywords:   []string{"salesforce", "salesforce login"},
		OfficialDomains: []string{"salesforce.com", "force.com", "my.salesforce.com"},
	},

	// ── Job & Freelancing Platforms ───────────────────────────────────────────
	"Indeed": {
		TitleKeywords:   []string{"indeed sign in", "indeed login", "indeed log in", "indeed account", "indeed verification"},
		Names:           []string{"indeed"},
		OfficialDomains: []string{"indeed.com"},
	},
	"Upwork": {
		TitleKeywords:   []string{"upwork", "upwork login"},
		OfficialDomains: []string{"upwork.com"},
	},
	"Fiverr": {
		TitleKeywords:   []string{"fiverr", "fiverr login"},
		OfficialDomains: []string{"fiverr.com"},
	},
	"Glassdoor": {
		TitleKeywords:   []string{"glassdoor"},
		OfficialDomains: []string{"glassdoor.com"},
	},
	"Naukri": {
		TitleKeywords:   []string{"naukri", "naukri.com login"},
		OfficialDomains: []string{"naukri.com"},
	},
	"Freelancer": {
		TitleKeywords:   []string{"freelancer", "freelancer login"},
		OfficialDomains: []string{"freelancer.com"},
	},

	// ── Delivery & Logistics ──────────────────────────────────────────────────
	"FedEx": {
		TitleKeywords:   []string{"fedex", "fedex tracking"},
		Names:           []string{"fedex"},
		OfficialDomains: []string{"fedex.com"},
	},
	"UPS": {
		TitleKeywords:   []string{"ups delivery", "ups tracking"},
		OfficialDomains: []string{"ups.com"},
	},
	"DHL": {
		TitleKeywords:   []string{"dhl", "dhl tracking"},
		Names:           []string{"dhl"},
		OfficialDomains: []string{"dhl.com", "dhl.de"},
	},
	"USPS": {
		TitleKeywords:   []string{"usps", "usps tracking", "united states postal"},
		Names:           []string{"usps"},
		OfficialDomains: []string{"usps.com"},
	},
	"Royal Mail": {
		TitleKeywords:   []string{"royal mail", "royal mail tracking"},
		Names:           []string{"royal mail"},
		OfficialDomains: []string{"royalmail.com"},
	},

	// ── Government & Tax ─────────────────────────────────────────────────────
	"IRS": {
		TitleKeywords:   []string{"irs", "internal revenue service", "irs refund"},
		OfficialDomains: []string{"irs.gov"},
	},
	"HMRC": {
		TitleKeywords:   []string{"hmrc", "her majesty revenue", "hmrc login"},
		OfficialDomains: []string{"hmrc.gov.uk", "gov.uk"},
	},
	"Income Tax India": {
		TitleKeywords:   []string{"income tax", "income tax department", "itr filing"},
		OfficialDomains: []string{"incometax.gov.in", "efiling.incometax.gov.in"},
	},
	"Aadhaar": {
		TitleKeywords:   []string{"aadhaar", "uidai", "aadhaar card", "aadhaar kyc"},
		OfficialDomains: []string{"uidai.gov.in", "myaadhaar.uidai.gov.in"},
	},
	"Social Security": {
		// Not the bare phrase: disability lawyers and benefits blogs use it.
		TitleKeywords:   []string{"social security administration", "ssa login", "my social security", "social security statement"},
		Names:           []string{"social security"},
		OfficialDomains: []string{"ssa.gov"},
	},
	"Medicare": {
		// Not the bare name: every insurance agent's site uses it.
		TitleKeywords:   []string{"medicare.gov", "medicare login", "mymedicare", "medicare account"},
		Names:           []string{"medicare"},
		OfficialDomains: []string{"medicare.gov", "cms.gov"},
	},

	// ── Tech Accounts & Platforms ─────────────────────────────────────────────
	"Apple": {
		TitleKeywords:   []string{"apple id", "apple account", "icloud", "find my iphone", "apple support", "apple store", "apple login", "apple log in", "apple sign in", "apple verification", "my apple"},
		Names:           []string{"apple"},
		OfficialDomains: []string{"apple.com", "icloud.com"},
		AlsoOfficial:    []string{"me.com", "mac.com", "apple.co", "apple.news", "itunes.com", "appleid.com", "cdn-apple.com", "mzstatic.com"},
	},
	"Meta": {
		TitleKeywords:   []string{"meta for business", "meta business suite", "meta business", "meta verified", "meta business help", "facebook business", "meta login", "meta log in", "meta sign in", "meta account", "meta verification", "my meta"},
		Names:           []string{"meta"},
		OfficialDomains: []string{"meta.com", "facebook.com", "fb.com", "instagram.com", "whatsapp.com"},
		AlsoOfficial:    []string{"messenger.com", "threads.net", "threads.com", "oculus.com", "workplace.com", "wa.me", "fb.me", "fbcdn.net"},
	},
	"Samsung": {
		TitleKeywords:   []string{"samsung account", "samsung login", "samsung log in", "samsung sign in", "samsung verification", "my samsung"},
		Names:           []string{"samsung"},
		OfficialDomains: []string{"samsung.com"},
	},
	"Huawei": {
		TitleKeywords:   []string{"huawei id", "huawei login", "huawei log in", "huawei sign in", "huawei account", "huawei verification", "my huawei"},
		Names:           []string{"huawei"},
		OfficialDomains: []string{"huawei.com", "huaweicloud.com"},
	},
	"Twitch": {
		TitleKeywords:   []string{"twitch login", "twitch log in", "twitch sign in", "twitch account", "twitch verification", "my twitch"},
		Names:           []string{"twitch"},
		OfficialDomains: []string{"twitch.tv"},
	},
	"Threads": {
		TitleKeywords:   []string{"threads login"},
		OfficialDomains: []string{"threads.net", "threads.com"},
	},
	"VK": {
		TitleKeywords:   []string{"vk login", "vkontakte"},
		OfficialDomains: []string{"vk.com"},
	},
	"WeChat": {
		TitleKeywords:   []string{"wechat login", "wechat log in", "wechat sign in", "wechat account", "wechat verification", "my wechat"},
		Names:           []string{"wechat"},
		OfficialDomains: []string{"wechat.com", "qq.com"},
	},
	"KakaoTalk": {
		TitleKeywords:   []string{"kakaotalk", "kakao account", "kakao login", "kakao log in", "kakao sign in", "kakao verification", "my kakao"},
		Names:           []string{"kakao"},
		OfficialDomains: []string{"kakao.com", "kakaocorp.com"},
	},
	"Naver": {
		TitleKeywords:   []string{"naver login", "naver log in", "naver sign in", "naver account", "naver verification", "my naver"},
		Names:           []string{"naver"},
		OfficialDomains: []string{"naver.com"},
	},
	"Tinder": {
		TitleKeywords:   []string{"tinder login", "tinder log in", "tinder sign in", "tinder account", "tinder verification", "my tinder"},
		Names:           []string{"tinder"},
		OfficialDomains: []string{"tinder.com", "gotinder.com"},
	},
	"Bumble": {
		TitleKeywords:   []string{"bumble login", "bumble log in", "bumble sign in", "bumble account", "bumble verification", "my bumble"},
		Names:           []string{"bumble"},
		OfficialDomains: []string{"bumble.com"},
	},
	"Match": {
		TitleKeywords:   []string{"match.com"},
		OfficialDomains: []string{"match.com"},
	},

	// ── Webmail & Hosting Logins ──────────────────────────────────────────────
	"ConstructConnect": {
		// Construction bid invitations are a common lure for its login page.
		TitleKeywords:   []string{"constructconnect"},
		OfficialDomains: []string{"constructconnect.com"},
	},
	"cPanel": {
		TitleKeywords:   []string{"cpanel", "cpanel login", "webmail cpanel"},
		OfficialDomains: []string{"cpanel.net", "cpanel.com"},
	},
	"Roundcube": {
		TitleKeywords:   []string{"roundcube webmail", "roundcube"},
		OfficialDomains: []string{"roundcube.net"},
	},
	"Zimbra": {
		TitleKeywords:   []string{"zimbra web client", "zimbra"},
		OfficialDomains: []string{"zimbra.com"},
	},
	"Outlook Web App": {
		TitleKeywords:   []string{"outlook web app", "outlook web access"},
		OfficialDomains: []string{"outlook.com", "office.com", "microsoft.com", "live.com"},
	},
	"GMX": {
		TitleKeywords:   []string{"gmx login", "gmx mail", "gmx freemail", "gmx log in", "gmx sign in", "gmx account", "gmx verification", "my gmx"},
		Names:           []string{"gmx"},
		OfficialDomains: []string{"gmx.com", "gmx.net", "gmx.de", "gmx.at", "gmx.ch"},
	},
	"WEB.DE": {
		TitleKeywords:   []string{"web.de login", "web.de freemail"},
		OfficialDomains: []string{"web.de"},
	},
	"Mail.ru": {
		TitleKeywords:   []string{"mail.ru"},
		OfficialDomains: []string{"mail.ru"},
	},
	"Yandex": {
		TitleKeywords:   []string{"yandex id", "yandex mail", "yandex.mail", "yandex login", "yandex log in", "yandex sign in", "yandex account", "yandex verification", "my yandex"},
		Names:           []string{"yandex"},
		OfficialDomains: []string{"yandex.ru", "yandex.com", "ya.ru"},
	},
	"GoDaddy": {
		TitleKeywords:   []string{"godaddy", "godaddy webmail", "workspace webmail"},
		OfficialDomains: []string{"godaddy.com", "secureserver.net"},
	},
	"Namecheap": {
		TitleKeywords:   []string{"namecheap"},
		OfficialDomains: []string{"namecheap.com", "privateemail.com"},
	},
	"Bluehost": {
		TitleKeywords:   []string{"bluehost"},
		OfficialDomains: []string{"bluehost.com"},
	},
	"IONOS": {
		TitleKeywords:   []string{"ionos login", "1&1 ionos", "ionos log in", "ionos sign in", "ionos account", "ionos verification", "my ionos"},
		Names:           []string{"ionos"},
		OfficialDomains: []string{"ionos.com", "ionos.de", "ionos.co.uk", "ionos.fr", "ionos.es"},
	},
	"OVHcloud": {
		TitleKeywords:   []string{"ovhcloud", "ovh webmail", "ovh login", "ovh log in", "ovh sign in", "ovh account", "ovh verification", "my ovh"},
		Names:           []string{"ovh"},
		OfficialDomains: []string{"ovh.com", "ovhcloud.com", "ovh.net"},
	},
	"Hostinger": {
		TitleKeywords:   []string{"hostinger"},
		OfficialDomains: []string{"hostinger.com"},
	},
	"Cloudflare": {
		TitleKeywords:   []string{"cloudflare dashboard", "cloudflare login"},
		OfficialDomains: []string{"cloudflare.com"},
	},
	"AWS": {
		TitleKeywords:   []string{"aws console", "aws management console", "amazon web services", "aws login", "aws log in", "aws sign in", "aws account", "aws verification", "my aws"},
		Names:           []string{"aws"},
		OfficialDomains: []string{"aws.amazon.com", "amazon.com", "amazonaws.com"},
	},
	"Okta": {
		TitleKeywords:   []string{"okta sign in", "okta"},
		OfficialDomains: []string{"okta.com", "oktapreview.com"},
	},
	"Citrix": {
		TitleKeywords:   []string{"citrix gateway", "citrix workspace", "citrix login", "citrix log in", "citrix sign in", "citrix account", "citrix verification", "my citrix"},
		Names:           []string{"citrix"},
		OfficialDomains: []string{"citrix.com", "cloud.com"},
	},
	"Webex": {
		TitleKeywords:   []string{"webex", "cisco webex"},
		OfficialDomains: []string{"webex.com", "cisco.com"},
	},
	"WeTransfer": {
		TitleKeywords:   []string{"wetransfer"},
		OfficialDomains: []string{"wetransfer.com", "we.tl"},
	},
	"SharePoint": {
		TitleKeywords:   []string{"sharepoint", "onedrive for business"},
		OfficialDomains: []string{"sharepoint.com", "microsoft.com", "office.com", "live.com"},
	},
	"Canva": {
		TitleKeywords:   []string{"canva login", "canva log in", "canva sign in", "canva account", "canva verification", "my canva"},
		Names:           []string{"canva"},
		OfficialDomains: []string{"canva.com"},
	},
	"Intuit": {
		TitleKeywords:   []string{"intuit", "mint login", "quickbooks login", "quickbooks log in", "quickbooks sign in", "quickbooks account", "quickbooks verification", "my quickbooks", "turbotax login", "turbotax log in", "turbotax sign in", "turbotax account", "turbotax verification", "my turbotax"},
		Names:           []string{"quickbooks", "turbotax"},
		OfficialDomains: []string{"intuit.com", "quickbooks.com", "turbotax.com"},
	},
	"Xero": {
		TitleKeywords:   []string{"xero login", "xero log in", "xero sign in", "xero account", "xero verification", "my xero"},
		Names:           []string{"xero"},
		OfficialDomains: []string{"xero.com"},
	},
	"Mailchimp": {
		TitleKeywords:   []string{"mailchimp login", "mailchimp log in", "mailchimp sign in", "mailchimp account", "mailchimp verification", "my mailchimp"},
		Names:           []string{"mailchimp"},
		OfficialDomains: []string{"mailchimp.com"},
	},
	"Squarespace": {
		TitleKeywords:   []string{"squarespace login"},
		OfficialDomains: []string{"squarespace.com"},
	},
	"Wix": {
		TitleKeywords:   []string{"wix login", "wix.com"},
		OfficialDomains: []string{"wix.com"},
	},

	// ── E-Commerce & Marketplaces ─────────────────────────────────────────────
	"Shopee": {
		TitleKeywords:   []string{"shopee login", "shopee log in", "shopee sign in", "shopee account", "shopee verification", "my shopee"},
		Names:           []string{"shopee"},
		OfficialDomains: []string{"shopee.com", "shopee.co.id", "shopee.ph", "shopee.sg", "shopee.vn", "shopee.com.my", "shopee.tw", "shopee.co.th", "shopee.com.br", "shopee.com.mx", "shopee.cl", "shopee.com.co"},
	},
	"Lazada": {
		TitleKeywords:   []string{"lazada login", "lazada log in", "lazada sign in", "lazada account", "lazada verification", "my lazada"},
		Names:           []string{"lazada"},
		OfficialDomains: []string{"lazada.com", "lazada.co.id", "lazada.com.ph", "lazada.sg", "lazada.vn", "lazada.com.my", "lazada.co.th"},
	},
	"Tokopedia": {
		TitleKeywords:   []string{"tokopedia login", "tokopedia log in", "tokopedia sign in", "tokopedia account", "tokopedia verification", "my tokopedia"},
		Names:           []string{"tokopedia"},
		OfficialDomains: []string{"tokopedia.com"},
	},
	"TikTok Shop": {
		TitleKeywords:   []string{"tiktok shop"},
		OfficialDomains: []string{"tiktok.com", "tiktokshop.com"},
	},
	"Temu": {
		TitleKeywords:   []string{"temu login", "temu log in", "temu sign in", "temu account", "temu verification", "my temu"},
		Names:           []string{"temu"},
		OfficialDomains: []string{"temu.com"},
	},
	"Shein": {
		TitleKeywords:   []string{"shein login", "shein log in", "shein sign in", "shein account", "shein verification", "my shein"},
		Names:           []string{"shein"},
		OfficialDomains: []string{"shein.com"},
	},
	"Mercado Libre": {
		TitleKeywords:   []string{"mercado libre", "mercadolibre", "mercado livre", "mercadolivre"},
		OfficialDomains: []string{"mercadolibre.com", "mercadolibre.com.ar", "mercadolibre.com.mx", "mercadolibre.com.co", "mercadolibre.cl", "mercadolivre.com.br"},
	},
	"Rakuten": {
		TitleKeywords:   []string{"rakuten login", "rakuten log in", "rakuten sign in", "rakuten account", "rakuten verification", "my rakuten"},
		Names:           []string{"rakuten"},
		OfficialDomains: []string{"rakuten.com", "rakuten.co.jp", "rakuten.fr", "rakuten.de"},
	},
	"Myntra": {
		TitleKeywords:   []string{"myntra"},
		OfficialDomains: []string{"myntra.com"},
	},
	"Jumia": {
		TitleKeywords:   []string{"jumia login", "jumia log in", "jumia sign in", "jumia account", "jumia verification", "my jumia"},
		Names:           []string{"jumia"},
		OfficialDomains: []string{"jumia.com", "jumia.com.ng", "jumia.co.ke", "jumia.com.eg", "jumia.ma"},
	},
	"Zalando": {
		TitleKeywords:   []string{"zalando login", "zalando log in", "zalando sign in", "zalando account", "zalando verification", "my zalando"},
		Names:           []string{"zalando"},
		OfficialDomains: []string{"zalando.com", "zalando.de", "zalando.fr", "zalando.co.uk", "zalando.it", "zalando.es", "zalando.nl"},
	},
	"IKEA": {
		TitleKeywords:   []string{"ikea login", "ikea log in", "ikea sign in", "ikea account", "ikea verification", "my ikea"},
		Names:           []string{"ikea"},
		OfficialDomains: []string{"ikea.com"},
	},
	"Costco": {
		TitleKeywords:   []string{"costco login", "costco log in", "costco sign in", "costco account", "costco verification", "my costco"},
		Names:           []string{"costco"},
		OfficialDomains: []string{"costco.com", "costco.ca", "costco.co.uk"},
	},
	"Target": {
		TitleKeywords:   []string{"target circle", "target gift card", "target login", "target log in", "target sign in", "target account", "target verification", "my target"},
		Names:           []string{"target"},
		OfficialDomains: []string{"target.com"},
	},
	"Best Buy": {
		TitleKeywords:   []string{"best buy login", "best buy log in", "best buy sign in", "best buy account", "best buy verification", "my best buy"},
		Names:           []string{"best buy"},
		OfficialDomains: []string{"bestbuy.com", "bestbuy.ca"},
	},
	"Home Depot": {
		TitleKeywords:   []string{"home depot login", "home depot log in", "home depot sign in", "home depot account", "home depot verification", "my home depot"},
		Names:           []string{"home depot"},
		OfficialDomains: []string{"homedepot.com", "homedepot.ca"},
	},
	"Lowe's": {
		TitleKeywords:   []string{"lowe's login", "lowe's log in", "lowe's sign in", "lowe's account", "lowe's verification", "my lowe's", "lowes login", "lowes log in", "lowes sign in", "lowes account", "lowes verification", "my lowes"},
		Names:           []string{"lowe's", "lowes"},
		OfficialDomains: []string{"lowes.com"},
	},
	"Kroger": {
		TitleKeywords:   []string{"kroger login", "kroger log in", "kroger sign in", "kroger account", "kroger verification", "my kroger"},
		Names:           []string{"kroger"},
		OfficialDomains: []string{"kroger.com"},
	},
	"Tesco": {
		TitleKeywords:   []string{"tesco login", "tesco log in", "tesco sign in", "tesco account", "tesco verification", "my tesco"},
		Names:           []string{"tesco"},
		OfficialDomains: []string{"tesco.com"},
	},
	"Carrefour": {
		TitleKeywords:   []string{"carrefour login", "carrefour log in", "carrefour sign in", "carrefour account", "carrefour verification", "my carrefour"},
		Names:           []string{"carrefour"},
		OfficialDomains: []string{"carrefour.com", "carrefour.fr", "carrefour.es", "carrefour.it", "carrefour.com.br"},
	},
	"Lidl": {
		TitleKeywords:   []string{"lidl login", "lidl log in", "lidl sign in", "lidl account", "lidl verification", "my lidl"},
		Names:           []string{"lidl"},
		OfficialDomains: []string{"lidl.com", "lidl.de", "lidl.fr", "lidl.co.uk", "lidl.es", "lidl.it"},
	},
	"Aldi": {
		TitleKeywords:   []string{"aldi login", "aldi log in", "aldi sign in", "aldi account", "aldi verification", "my aldi"},
		Names:           []string{"aldi"},
		OfficialDomains: []string{"aldi.com", "aldi.us", "aldi.co.uk", "aldi-nord.de", "aldi-sued.de"},
	},
	"Nike": {
		TitleKeywords:   []string{"nike login", "nike log in", "nike sign in", "nike account", "nike verification", "my nike"},
		Names:           []string{"nike"},
		OfficialDomains: []string{"nike.com"},
	},
	"Adidas": {
		TitleKeywords:   []string{"adidas login", "adidas log in", "adidas sign in", "adidas account", "adidas verification", "my adidas"},
		Names:           []string{"adidas"},
		OfficialDomains: []string{"adidas.com"},
	},

	// ── Payments, Cards & Money Transfer ──────────────────────────────────────
	"Visa": {
		TitleKeywords:   []string{"verified by visa", "visa secure", "visa login", "visa log in", "visa sign in", "visa account", "visa verification", "my visa"},
		Names:           []string{"visa"},
		OfficialDomains: []string{"visa.com"},
	},
	"Mastercard": {
		TitleKeywords:   []string{"mastercard", "mastercard id check"},
		OfficialDomains: []string{"mastercard.com"},
	},
	"American Express": {
		TitleKeywords:   []string{"american express", "amex"},
		OfficialDomains: []string{"americanexpress.com", "amex.com"},
	},
	"Discover": {
		TitleKeywords:   []string{"discover card"},
		OfficialDomains: []string{"discover.com"},
	},
	"Capital One": {
		TitleKeywords:   []string{"capital one"},
		OfficialDomains: []string{"capitalone.com", "capitalone.ca", "capitalone.co.uk"},
	},
	"Klarna": {
		TitleKeywords:   []string{"klarna login", "klarna log in", "klarna sign in", "klarna account", "klarna verification", "my klarna"},
		Names:           []string{"klarna"},
		OfficialDomains: []string{"klarna.com"},
	},
	"Afterpay": {
		TitleKeywords:   []string{"afterpay login", "afterpay log in", "afterpay sign in", "afterpay account", "afterpay verification", "my afterpay", "clearpay login", "clearpay log in", "clearpay sign in", "clearpay account", "clearpay verification", "my clearpay"},
		Names:           []string{"afterpay", "clearpay"},
		OfficialDomains: []string{"afterpay.com", "clearpay.co.uk"},
	},
	"Revolut": {
		TitleKeywords:   []string{"revolut login", "revolut log in", "revolut sign in", "revolut account", "revolut verification", "my revolut"},
		Names:           []string{"revolut"},
		OfficialDomains: []string{"revolut.com"},
	},
	"N26": {
		TitleKeywords:   []string{"n26 login", "n26 log in", "n26 sign in", "n26 account", "n26 verification", "my n26"},
		Names:           []string{"n26"},
		OfficialDomains: []string{"n26.com"},
	},
	"Monzo": {
		TitleKeywords:   []string{"monzo login", "monzo log in", "monzo sign in", "monzo account", "monzo verification", "my monzo"},
		Names:           []string{"monzo"},
		OfficialDomains: []string{"monzo.com"},
	},
	"Skrill": {
		TitleKeywords:   []string{"skrill login", "skrill log in", "skrill sign in", "skrill account", "skrill verification", "my skrill"},
		Names:           []string{"skrill"},
		OfficialDomains: []string{"skrill.com"},
	},
	"Payoneer": {
		TitleKeywords:   []string{"payoneer login", "payoneer log in", "payoneer sign in", "payoneer account", "payoneer verification", "my payoneer"},
		Names:           []string{"payoneer"},
		OfficialDomains: []string{"payoneer.com"},
	},
	"Western Union": {
		TitleKeywords:   []string{"western union"},
		OfficialDomains: []string{"westernunion.com"},
	},
	"MoneyGram": {
		TitleKeywords:   []string{"moneygram"},
		OfficialDomains: []string{"moneygram.com"},
	},
	"Remitly": {
		TitleKeywords:   []string{"remitly"},
		OfficialDomains: []string{"remitly.com"},
	},
	"Mercado Pago": {
		TitleKeywords:   []string{"mercado pago", "mercadopago"},
		OfficialDomains: []string{"mercadopago.com", "mercadopago.com.ar", "mercadopago.com.mx", "mercadopago.com.br"},
	},
	"Alipay": {
		TitleKeywords:   []string{"alipay"},
		OfficialDomains: []string{"alipay.com"},
	},
	"Interac": {
		TitleKeywords:   []string{"interac", "interac e-transfer"},
		OfficialDomains: []string{"interac.ca"},
	},
	"Square": {
		TitleKeywords:   []string{"square login", "cash app"},
		OfficialDomains: []string{"squareup.com", "square.com", "cash.app"},
	},
	"GCash": {
		TitleKeywords:   []string{"gcash"},
		OfficialDomains: []string{"gcash.com"},
	},
	"M-Pesa": {
		TitleKeywords:   []string{"m-pesa", "mpesa"},
		OfficialDomains: []string{"safaricom.co.ke", "vodacom.co.tz"},
	},
	"Bizum": {
		TitleKeywords:   []string{"bizum"},
		OfficialDomains: []string{"bizum.es"},
	},
	"Twint": {
		TitleKeywords:   []string{"twint"},
		OfficialDomains: []string{"twint.ch"},
	},
	"Swish": {
		TitleKeywords:   []string{"swish betalning"},
		OfficialDomains: []string{"swish.nu"},
	},
	"Satispay": {
		TitleKeywords:   []string{"satispay"},
		OfficialDomains: []string{"satispay.com"},
	},

	// ── Banking — North America ───────────────────────────────────────────────
	"US Bank": {
		TitleKeywords:   []string{"u.s. bank", "us bank"},
		OfficialDomains: []string{"usbank.com"},
	},
	"PNC": {
		TitleKeywords:   []string{"pnc bank", "pnc online banking", "pnc login", "pnc log in", "pnc sign in", "pnc account", "pnc verification", "my pnc"},
		Names:           []string{"pnc"},
		OfficialDomains: []string{"pnc.com"},
	},
	"Truist": {
		TitleKeywords:   []string{"truist"},
		OfficialDomains: []string{"truist.com"},
	},
	"Ally": {
		TitleKeywords:   []string{"ally bank"},
		OfficialDomains: []string{"ally.com"},
	},
	"Navy Federal": {
		TitleKeywords:   []string{"navy federal"},
		OfficialDomains: []string{"navyfederal.org"},
	},
	"USAA": {
		TitleKeywords:   []string{"usaa"},
		OfficialDomains: []string{"usaa.com"},
	},
	"Citizens": {
		TitleKeywords:   []string{"citizens bank"},
		OfficialDomains: []string{"citizensbank.com"},
	},
	"Fifth Third": {
		TitleKeywords:   []string{"fifth third"},
		OfficialDomains: []string{"53.com"},
	},
	"Regions": {
		TitleKeywords:   []string{"regions bank", "regions online banking"},
		OfficialDomains: []string{"regions.com"},
	},
	"KeyBank": {
		TitleKeywords:   []string{"keybank"},
		OfficialDomains: []string{"key.com"},
	},
	"Huntington": {
		TitleKeywords:   []string{"huntington bank"},
		OfficialDomains: []string{"huntington.com"},
	},
	"Charles Schwab": {
		TitleKeywords:   []string{"charles schwab", "schwab"},
		OfficialDomains: []string{"schwab.com"},
	},
	"Fidelity": {
		TitleKeywords:   []string{"fidelity investments", "netbenefits"},
		OfficialDomains: []string{"fidelity.com"},
	},
	"Vanguard": {
		TitleKeywords:   []string{"vanguard login"},
		OfficialDomains: []string{"vanguard.com"},
	},
	"BMO": {
		TitleKeywords:   []string{"bmo", "bank of montreal"},
		OfficialDomains: []string{"bmo.com"},
	},
	"CIBC": {
		TitleKeywords:   []string{"cibc"},
		OfficialDomains: []string{"cibc.com"},
	},
	"Desjardins": {
		TitleKeywords:   []string{"desjardins"},
		OfficialDomains: []string{"desjardins.com"},
	},
	"Tangerine": {
		TitleKeywords:   []string{"tangerine bank"},
		OfficialDomains: []string{"tangerine.ca"},
	},
	"Banorte": {
		TitleKeywords:   []string{"banorte"},
		OfficialDomains: []string{"banorte.com"},
	},

	// ── Banking — Europe ──────────────────────────────────────────────────────
	"ING": {
		TitleKeywords:   []string{"ing bank", "mijn ing", "ing banking", "ing login", "ing log in", "ing sign in", "ing account", "ing verification", "my ing"},
		Names:           []string{"ing"},
		OfficialDomains: []string{"ing.com", "ing.nl", "ing.de", "ing.be", "ing.es", "ing.pl", "ing.ro", "ing.it"},
	},
	"Deutsche Bank": {
		TitleKeywords:   []string{"deutsche bank"},
		OfficialDomains: []string{"deutsche-bank.de", "db.com"},
	},
	"Commerzbank": {
		TitleKeywords:   []string{"commerzbank"},
		OfficialDomains: []string{"commerzbank.de", "commerzbank.com"},
	},
	"Sparkasse": {
		TitleKeywords:   []string{"sparkasse"},
		OfficialDomains: []string{"sparkasse.de", "sparkasse.at"},
		OwnNames:        []string{"sparkasse", "kreissparkasse", "stadtsparkasse"},
	},
	"Postbank": {
		TitleKeywords:   []string{"postbank"},
		OfficialDomains: []string{"postbank.de"},
	},
	"Volksbank": {
		TitleKeywords:   []string{"volksbank", "raiffeisenbank", "vr banking"},
		OfficialDomains: []string{"vr.de", "volksbank.de"},
		OwnNames:        []string{"volksbank", "vrbank", "vr-bank", "raiffeisenbank"},
	},
	"BNP Paribas": {
		TitleKeywords:   []string{"bnp paribas"},
		OfficialDomains: []string{"bnpparibas", "bnpparibas.com", "mabanque.bnpparibas", "bnpparibas.fr"},
	},
	"Société Générale": {
		TitleKeywords:   []string{"société générale", "societe generale"},
		OfficialDomains: []string{"societegenerale.fr", "societegenerale.com", "particuliers.sg.fr", "sg.fr"},
	},
	"Crédit Agricole": {
		TitleKeywords:   []string{"crédit agricole", "credit agricole"},
		OfficialDomains: []string{"credit-agricole.fr", "credit-agricole.com"},
	},
	"La Banque Postale": {
		TitleKeywords:   []string{"la banque postale"},
		OfficialDomains: []string{"labanquepostale.fr"},
	},
	"Crédit Mutuel": {
		TitleKeywords:   []string{"crédit mutuel", "credit mutuel"},
		OfficialDomains: []string{"creditmutuel.fr"},
	},
	"LCL": {
		TitleKeywords:   []string{"lcl banque", "lcl mes comptes"},
		OfficialDomains: []string{"lcl.fr"},
	},
	"Caisse d'Epargne": {
		TitleKeywords:   []string{"caisse d'epargne", "caisse d'épargne"},
		OfficialDomains: []string{"caisse-epargne.fr"},
	},
	"Intesa Sanpaolo": {
		TitleKeywords:   []string{"intesa sanpaolo"},
		OfficialDomains: []string{"intesasanpaolo.com"},
	},
	"UniCredit": {
		TitleKeywords:   []string{"unicredit"},
		OfficialDomains: []string{"unicredit.it", "unicreditgroup.eu", "unicredit.eu"},
	},
	"Poste Italiane": {
		TitleKeywords:   []string{"poste italiane", "postepay", "bancoposta"},
		OfficialDomains: []string{"poste.it", "posteitaliane.it"},
	},
	"BBVA": {
		TitleKeywords:   []string{"bbva"},
		OfficialDomains: []string{"bbva.com", "bbva.es", "bbva.mx", "bbva.com.co", "bbva.pe", "bbva.com.ar"},
	},
	"CaixaBank": {
		TitleKeywords:   []string{"caixabank", "caixa bank"},
		OfficialDomains: []string{"caixabank.es", "caixabank.com"},
	},
	"Sabadell": {
		TitleKeywords:   []string{"banco sabadell"},
		OfficialDomains: []string{"bancsabadell.com"},
	},
	"ABN AMRO": {
		TitleKeywords:   []string{"abn amro"},
		OfficialDomains: []string{"abnamro.nl", "abnamro.com"},
	},
	"Rabobank": {
		TitleKeywords:   []string{"rabobank"},
		OfficialDomains: []string{"rabobank.nl", "rabobank.com"},
	},
	"KBC": {
		TitleKeywords:   []string{"kbc bank", "kbc mobile"},
		OfficialDomains: []string{"kbc.be", "kbc.com"},
	},
	"Belfius": {
		TitleKeywords:   []string{"belfius"},
		OfficialDomains: []string{"belfius.be"},
	},
	"Nordea": {
		TitleKeywords:   []string{"nordea"},
		OfficialDomains: []string{"nordea.com", "nordea.fi", "nordea.se", "nordea.dk", "nordea.no"},
	},
	"Danske Bank": {
		TitleKeywords:   []string{"danske bank"},
		OfficialDomains: []string{"danskebank.com", "danskebank.dk", "danskebank.fi", "danskebank.no", "danskebank.se"},
	},
	"Swedbank": {
		TitleKeywords:   []string{"swedbank"},
		OfficialDomains: []string{"swedbank.se", "swedbank.com", "swedbank.ee", "swedbank.lv", "swedbank.lt"},
	},
	"Handelsbanken": {
		TitleKeywords:   []string{"handelsbanken"},
		OfficialDomains: []string{"handelsbanken.se", "handelsbanken.com"},
	},
	"PostFinance": {
		TitleKeywords:   []string{"postfinance"},
		OfficialDomains: []string{"postfinance.ch"},
	},
	"Raiffeisen": {
		TitleKeywords:   []string{"raiffeisen"},
		OfficialDomains: []string{"raiffeisen.ch", "raiffeisen.at", "rbinternational.com"},
		OwnNames:        []string{"raiffeisen"},
	},
	"Halifax": {
		TitleKeywords:   []string{"halifax online banking"},
		OfficialDomains: []string{"halifax.co.uk"},
	},
	"Nationwide": {
		TitleKeywords:   []string{"nationwide building society"},
		OfficialDomains: []string{"nationwide.co.uk"},
	},
	"Bank of Ireland": {
		TitleKeywords:   []string{"bank of ireland"},
		OfficialDomains: []string{"bankofireland.com", "boi.com"},
	},
	"AIB": {
		TitleKeywords:   []string{"aib internet banking", "allied irish"},
		OfficialDomains: []string{"aib.ie"},
	},
	"PKO BP": {
		TitleKeywords:   []string{"pko bp", "ipko"},
		OfficialDomains: []string{"pkobp.pl"},
	},
	"mBank": {
		TitleKeywords:   []string{"mbank"},
		OfficialDomains: []string{"mbank.pl", "mbank.cz"},
	},
	"OTP Bank": {
		TitleKeywords:   []string{"otp bank"},
		OfficialDomains: []string{"otpbank.hu", "otpbank.com"},
	},
	"Millennium BCP": {
		TitleKeywords:   []string{"millennium bcp", "millenniumbcp"},
		OfficialDomains: []string{"millenniumbcp.pt"},
	},
	"Caixa Geral": {
		TitleKeywords:   []string{"caixa geral de depósitos", "caixadirecta"},
		OfficialDomains: []string{"cgd.pt"},
	},
	"Garanti BBVA": {
		TitleKeywords:   []string{"garanti bbva"},
		OfficialDomains: []string{"garantibbva.com.tr"},
	},
	"Ziraat": {
		TitleKeywords:   []string{"ziraat bankası", "ziraat bankasi"},
		OfficialDomains: []string{"ziraatbank.com.tr"},
	},
	"Sberbank": {
		TitleKeywords:   []string{"сбербанк", "sberbank login", "sberbank log in", "sberbank sign in", "sberbank account", "sberbank verification", "my sberbank"},
		Names:           []string{"sberbank"},
		OfficialDomains: []string{"sberbank.ru", "sber.ru"},
	},

	// ── Banking — Latin America ───────────────────────────────────────────────
	"Itaú": {
		TitleKeywords:   []string{"itaú", "itau"},
		OfficialDomains: []string{"itau.com.br", "itau.com", "itau.cl", "itau.com.uy"},
	},
	"Bradesco": {
		TitleKeywords:   []string{"bradesco"},
		OfficialDomains: []string{"bradesco.com.br", "bradesco"},
	},
	"Banco do Brasil": {
		TitleKeywords:   []string{"banco do brasil"},
		OfficialDomains: []string{"bb.com.br"},
	},
	"Caixa": {
		TitleKeywords:   []string{"caixa econômica", "caixa economica", "caixa tem"},
		OfficialDomains: []string{"caixa.gov.br"},
	},
	"Nubank": {
		TitleKeywords:   []string{"nubank login", "nubank log in", "nubank sign in", "nubank account", "nubank verification", "my nubank"},
		Names:           []string{"nubank"},
		OfficialDomains: []string{"nubank.com.br", "nu.com.mx", "nubank.com"},
	},
	"Promerica": {
		TitleKeywords:   []string{"banco promerica", "promerica en línea", "promerica online", "promerica login", "promerica sign in"},
		Names:           []string{"promerica"},
		OfficialDomains: []string{"promerica.fi.cr", "bancopromerica.com", "bancopromerica.com.gt", "promerica.com.sv", "promerica.com.ni", "promerica.com.do", "bancopromerica.com.hn", "promerica.ec"},
	},
	"Banco Pichincha": {
		TitleKeywords:   []string{"banco pichincha", "pichincha"},
		OfficialDomains: []string{"pichincha.com", "bancopichincha.com.co"},
	},
	"Bancolombia": {
		TitleKeywords:   []string{"bancolombia"},
		OfficialDomains: []string{"bancolombia.com"},
	},
	"Banco General": {
		TitleKeywords:   []string{"banco general"},
		OfficialDomains: []string{"bgeneral.com"},
	},
	"BCP": {
		TitleKeywords:   []string{"banca por internet bcp", "viabcp"},
		OfficialDomains: []string{"viabcp.com"},
	},
	"BancoEstado": {
		TitleKeywords:   []string{"bancoestado", "banco estado"},
		OfficialDomains: []string{"bancoestado.cl"},
	},
	"Banco de Chile": {
		TitleKeywords:   []string{"banco de chile"},
		OfficialDomains: []string{"bancochile.cl"},
	},
	"Banco Nación": {
		TitleKeywords:   []string{"banco nación", "banco nacion"},
		OfficialDomains: []string{"bna.com.ar"},
	},
	"Banamex": {
		TitleKeywords:   []string{"banamex", "citibanamex"},
		OfficialDomains: []string{"banamex.com"},
	},
	"Interbank": {
		TitleKeywords:   []string{"interbank"},
		OfficialDomains: []string{"interbank.pe"},
	},
	"Banco Guayaquil": {
		TitleKeywords:   []string{"banco guayaquil"},
		OfficialDomains: []string{"bancoguayaquil.com"},
	},
	"Davivienda": {
		TitleKeywords:   []string{"davivienda"},
		OfficialDomains: []string{"davivienda.com"},
	},

	// ── Banking — Asia, Middle East & Africa ──────────────────────────────────
	"Maybank": {
		TitleKeywords:   []string{"maybank", "maybank2u"},
		OfficialDomains: []string{"maybank.com", "maybank2u.com.my", "maybank.co.id"},
	},
	"CIMB": {
		TitleKeywords:   []string{"cimb clicks", "cimb bank", "cimb login", "cimb log in", "cimb sign in", "cimb account", "cimb verification", "my cimb", "cimb octo"},
		Names:           []string{"cimb"},
		OfficialDomains: []string{"cimb.com", "cimbclicks.com.my", "cimbniaga.co.id"},
	},
	"Public Bank": {
		TitleKeywords:   []string{"public bank", "pbe bank"},
		OfficialDomains: []string{"pbebank.com", "publicbank.com.my"},
	},
	"OCBC": {
		TitleKeywords:   []string{"ocbc"},
		OfficialDomains: []string{"ocbc.com", "ocbc.id"},
	},
	"UOB": {
		TitleKeywords:   []string{"uob"},
		OfficialDomains: []string{"uob.com.sg", "uob.com.my", "uob.co.th", "uob.co.id", "uobgroup.com"},
	},
	"BDO": {
		TitleKeywords:   []string{"bdo online banking", "bdo unibank", "bdo login", "bdo log in", "bdo sign in", "bdo account", "bdo verification", "my bdo"},
		Names:           []string{"bdo"},
		OfficialDomains: []string{"bdo.com.ph"},
	},
	"BPI": {
		TitleKeywords:   []string{"bpi online", "bank of the philippine islands"},
		OfficialDomains: []string{"bpi.com.ph"},
	},
	"Metrobank": {
		TitleKeywords:   []string{"metrobank"},
		OfficialDomains: []string{"metrobank.com.ph"},
	},
	"Kasikornbank": {
		TitleKeywords:   []string{"kasikornbank", "k plus", "kbank"},
		OfficialDomains: []string{"kasikornbank.com"},
	},
	"Bangkok Bank": {
		TitleKeywords:   []string{"bangkok bank"},
		OfficialDomains: []string{"bangkokbank.com"},
	},
	"SCB": {
		TitleKeywords:   []string{"siam commercial bank", "scb easy"},
		OfficialDomains: []string{"scb.co.th"},
	},
	"Bank Mandiri": {
		TitleKeywords:   []string{"bank mandiri", "livin by mandiri", "mandiri login", "mandiri log in", "mandiri sign in", "mandiri account", "mandiri verification", "my mandiri"},
		Names:           []string{"mandiri"},
		OfficialDomains: []string{"bankmandiri.co.id"},
	},
	"BCA": {
		TitleKeywords:   []string{"klikbca", "bank central asia", "bca mobile"},
		OfficialDomains: []string{"bca.co.id", "klikbca.com"},
	},
	"BRI": {
		TitleKeywords:   []string{"bank rakyat indonesia", "brimo"},
		OfficialDomains: []string{"bri.co.id"},
	},
	"BNI": {
		TitleKeywords:   []string{"bank negara indonesia", "bni mobile"},
		OfficialDomains: []string{"bni.co.id"},
	},
	"Vietcombank": {
		TitleKeywords:   []string{"vietcombank"},
		OfficialDomains: []string{"vietcombank.com.vn"},
	},
	"Techcombank": {
		TitleKeywords:   []string{"techcombank"},
		OfficialDomains: []string{"techcombank.com.vn", "techcombank.com"},
	},
	"MB Bank": {
		TitleKeywords:   []string{"mb bank", "mbbank"},
		OfficialDomains: []string{"mbbank.com.vn"},
	},
	"ICBC": {
		TitleKeywords:   []string{"icbc", "industrial and commercial bank of china"},
		OfficialDomains: []string{"icbc.com.cn", "icbc.com"},
	},
	"Bank of China": {
		TitleKeywords:   []string{"bank of china"},
		OfficialDomains: []string{"boc.cn", "bankofchina.com"},
	},
	"MUFG": {
		TitleKeywords:   []string{"mufg", "mitsubishi ufj"},
		OfficialDomains: []string{"mufg.jp", "bk.mufg.jp"},
	},
	"SMBC": {
		TitleKeywords:   []string{"smbc", "sumitomo mitsui"},
		OfficialDomains: []string{"smbc.co.jp"},
	},
	"Mizuho": {
		TitleKeywords:   []string{"mizuho"},
		OfficialDomains: []string{"mizuhobank.co.jp", "mizuho-fg.co.jp"},
	},
	"Japan Post Bank": {
		TitleKeywords:   []string{"ゆうちょ", "japan post bank"},
		OfficialDomains: []string{"jp-bank.japanpost.jp", "japanpost.jp"},
	},
	"Shinhan": {
		TitleKeywords:   []string{"shinhan bank"},
		OfficialDomains: []string{"shinhan.com"},
	},
	"KB Kookmin": {
		TitleKeywords:   []string{"kookmin bank", "kb국민은행"},
		OfficialDomains: []string{"kbstar.com"},
	},
	"Emirates NBD": {
		TitleKeywords:   []string{"emirates nbd"},
		OfficialDomains: []string{"emiratesnbd.com"},
	},
	"ADCB": {
		TitleKeywords:   []string{"adcb"},
		OfficialDomains: []string{"adcb.com"},
	},
	"First Abu Dhabi Bank": {
		TitleKeywords:   []string{"first abu dhabi bank", "fab bank"},
		OfficialDomains: []string{"bankfab.com"},
	},
	"QNB": {
		TitleKeywords:   []string{"qnb"},
		OfficialDomains: []string{"qnb.com"},
	},
	"Al Rajhi": {
		TitleKeywords:   []string{"al rajhi", "alrajhi"},
		OfficialDomains: []string{"alrajhibank.com.sa"},
	},
	"Standard Bank": {
		TitleKeywords:   []string{"standard bank"},
		OfficialDomains: []string{"standardbank.co.za", "standardbank.com"},
	},
	"FNB": {
		TitleKeywords:   []string{"fnb online banking", "first national bank"},
		OfficialDomains: []string{"fnb.co.za"},
	},
	"Absa": {
		TitleKeywords:   []string{"absa"},
		OfficialDomains: []string{"absa.co.za", "absa.africa"},
	},
	"Capitec": {
		TitleKeywords:   []string{"capitec"},
		OfficialDomains: []string{"capitecbank.co.za"},
	},
	"Nedbank": {
		TitleKeywords:   []string{"nedbank"},
		OfficialDomains: []string{"nedbank.co.za"},
	},
	"Yes Bank": {
		TitleKeywords:   []string{"yes bank login", "yes bank log in", "yes bank sign in", "yes bank account", "yes bank verification", "my yes bank"},
		Names:           []string{"yes bank"},
		OfficialDomains: []string{"yesbank.in"},
	},
	"IndusInd": {
		TitleKeywords:   []string{"indusind login", "indusind log in", "indusind sign in", "indusind account", "indusind verification", "my indusind"},
		Names:           []string{"indusind"},
		OfficialDomains: []string{"indusind.com"},
	},
	"Union Bank of India": {
		TitleKeywords:   []string{"union bank of india login", "union bank of india log in", "union bank of india sign in", "union bank of india account", "union bank of india verification", "my union bank of india"},
		Names:           []string{"union bank of india"},
		OfficialDomains: []string{"unionbankofindia.co.in"},
	},
	"Bank of India": {
		TitleKeywords:   []string{"bank of india login", "bank of india log in", "bank of india sign in", "bank of india account", "bank of india verification", "my bank of india"},
		Names:           []string{"bank of india"},
		OfficialDomains: []string{"bankofindia.co.in"},
	},
	"IDFC First": {
		TitleKeywords:   []string{"idfc first login", "idfc first log in", "idfc first sign in", "idfc first account", "idfc first verification", "my idfc first", "idfc login", "idfc log in", "idfc sign in", "idfc account", "idfc verification", "my idfc"},
		Names:           []string{"idfc first", "idfc"},
		OfficialDomains: []string{"idfcfirstbank.com"},
	},
	"Federal Bank": {
		TitleKeywords:   []string{"federal bank login", "federal bank log in", "federal bank sign in", "federal bank account", "federal bank verification", "my federal bank"},
		Names:           []string{"federal bank"},
		OfficialDomains: []string{"federalbank.co.in"},
	},

	// ── Crypto & Web3 (more) ──────────────────────────────────────────────────
	"MoonPay": {
		TitleKeywords:   []string{"moonpay"},
		OfficialDomains: []string{"moonpay.com"},
	},
	"Exodus": {
		TitleKeywords:   []string{"exodus wallet", "exodus web3", "exodus login", "exodus log in", "exodus sign in", "exodus account", "exodus verification", "my exodus"},
		Names:           []string{"exodus"},
		OfficialDomains: []string{"exodus.com", "exodus.io"},
	},
	"KuCoin": {
		TitleKeywords:   []string{"kucoin login", "kucoin log in", "kucoin sign in", "kucoin account", "kucoin verification", "my kucoin"},
		Names:           []string{"kucoin"},
		OfficialDomains: []string{"kucoin.com"},
	},
	"Crypto.com": {
		TitleKeywords:   []string{"crypto.com"},
		OfficialDomains: []string{"crypto.com"},
	},
	"Bitget": {
		TitleKeywords:   []string{"bitget login", "bitget log in", "bitget sign in", "bitget account", "bitget verification", "my bitget"},
		Names:           []string{"bitget"},
		OfficialDomains: []string{"bitget.com"},
	},
	"Gate.io": {
		TitleKeywords:   []string{"gate.io"},
		OfficialDomains: []string{"gate.io", "gate.com"},
	},
	"HTX": {
		TitleKeywords:   []string{"htx login", "htx log in", "htx sign in", "htx account", "htx verification", "my htx", "huobi login", "huobi log in", "huobi sign in", "huobi account", "huobi verification", "my huobi"},
		Names:           []string{"htx", "huobi"},
		OfficialDomains: []string{"htx.com", "huobi.com"},
	},
	"Bitfinex": {
		TitleKeywords:   []string{"bitfinex"},
		OfficialDomains: []string{"bitfinex.com"},
	},
	"Bitstamp": {
		TitleKeywords:   []string{"bitstamp"},
		OfficialDomains: []string{"bitstamp.net"},
	},
	"Blockchain.com": {
		TitleKeywords:   []string{"blockchain.com", "blockchain wallet"},
		OfficialDomains: []string{"blockchain.com"},
	},
	"Trezor": {
		TitleKeywords:   []string{"trezor suite", "trezor login", "trezor log in", "trezor sign in", "trezor account", "trezor verification", "my trezor"},
		Names:           []string{"trezor"},
		OfficialDomains: []string{"trezor.io"},
	},
	"Atomic Wallet": {
		TitleKeywords:   []string{"atomic wallet"},
		OfficialDomains: []string{"atomicwallet.io"},
	},
	"Electrum": {
		TitleKeywords:   []string{"electrum wallet"},
		OfficialDomains: []string{"electrum.org"},
	},
	"WalletConnect": {
		TitleKeywords:   []string{"walletconnect login", "walletconnect log in", "walletconnect sign in", "walletconnect account", "walletconnect verification", "my walletconnect"},
		Names:           []string{"walletconnect"},
		OfficialDomains: []string{"walletconnect.com", "walletconnect.network", "reown.com"},
	},
	"Solflare": {
		TitleKeywords:   []string{"solflare login", "solflare log in", "solflare sign in", "solflare account", "solflare verification", "my solflare"},
		Names:           []string{"solflare"},
		OfficialDomains: []string{"solflare.com"},
	},
	"Keplr": {
		TitleKeywords:   []string{"keplr login", "keplr log in", "keplr sign in", "keplr account", "keplr verification", "my keplr"},
		Names:           []string{"keplr"},
		OfficialDomains: []string{"keplr.app"},
	},
	"Rabby": {
		TitleKeywords:   []string{"rabby wallet"},
		OfficialDomains: []string{"rabby.io"},
	},
	"PancakeSwap": {
		TitleKeywords:   []string{"pancakeswap login", "pancakeswap log in", "pancakeswap sign in", "pancakeswap account", "pancakeswap verification", "my pancakeswap"},
		Names:           []string{"pancakeswap"},
		OfficialDomains: []string{"pancakeswap.finance"},
	},
	"Aave": {
		TitleKeywords:   []string{"aave login", "aave log in", "aave sign in", "aave account", "aave verification", "my aave"},
		Names:           []string{"aave"},
		OfficialDomains: []string{"aave.com"},
	},
	"Lido": {
		TitleKeywords:   []string{"lido finance", "lido staking"},
		OfficialDomains: []string{"lido.fi"},
	},
	"Upbit": {
		TitleKeywords:   []string{"upbit login", "upbit log in", "upbit sign in", "upbit account", "upbit verification", "my upbit"},
		Names:           []string{"upbit"},
		OfficialDomains: []string{"upbit.com"},
	},
	"Bithumb": {
		TitleKeywords:   []string{"bithumb login", "bithumb log in", "bithumb sign in", "bithumb account", "bithumb verification", "my bithumb"},
		Names:           []string{"bithumb"},
		OfficialDomains: []string{"bithumb.com"},
	},
	"Bitpanda": {
		TitleKeywords:   []string{"bitpanda login", "bitpanda log in", "bitpanda sign in", "bitpanda account", "bitpanda verification", "my bitpanda"},
		Names:           []string{"bitpanda"},
		OfficialDomains: []string{"bitpanda.com"},
	},
	"CoinDCX": {
		TitleKeywords:   []string{"coindcx login", "coindcx log in", "coindcx sign in", "coindcx account", "coindcx verification", "my coindcx"},
		Names:           []string{"coindcx"},
		OfficialDomains: []string{"coindcx.com"},
	},
	"MEXC": {
		TitleKeywords:   []string{"mexc login", "mexc log in", "mexc sign in", "mexc account", "mexc verification", "my mexc"},
		Names:           []string{"mexc"},
		OfficialDomains: []string{"mexc.com"},
	},

	// ── Telecom & Internet Providers ──────────────────────────────────────────
	"Xfinity": {
		TitleKeywords:   []string{"xfinity login", "xfinity log in", "xfinity sign in", "xfinity account", "xfinity verification", "my xfinity", "comcast login", "comcast log in", "comcast sign in", "comcast account", "comcast verification", "my comcast"},
		Names:           []string{"xfinity", "comcast"},
		OfficialDomains: []string{"xfinity.com", "comcast.net", "comcast.com"},
	},
	"AT&T": {
		TitleKeywords:   []string{"att.net", "at&t login", "at&t log in", "at&t sign in", "at&t account", "at&t verification", "my at&t"},
		Names:           []string{"at&t"},
		OfficialDomains: []string{"att.com", "att.net", "currently.com"},
	},
	"Verizon": {
		TitleKeywords:   []string{"verizon login", "verizon log in", "verizon sign in", "verizon account", "verizon verification", "my verizon"},
		Names:           []string{"verizon"},
		OfficialDomains: []string{"verizon.com", "verizonwireless.com"},
	},
	"T-Mobile": {
		TitleKeywords:   []string{"t-mobile login", "t-mobile log in", "t-mobile sign in", "t-mobile account", "t-mobile verification", "my t-mobile"},
		Names:           []string{"t-mobile"},
		OfficialDomains: []string{"t-mobile.com"},
	},
	"Spectrum": {
		TitleKeywords:   []string{"spectrum login", "spectrum.net", "charter spectrum"},
		OfficialDomains: []string{"spectrum.net", "spectrum.com", "charter.com"},
	},
	"Cox": {
		TitleKeywords:   []string{"cox communications", "cox webmail"},
		OfficialDomains: []string{"cox.com", "cox.net"},
	},
	"Sunrise": {
		TitleKeywords:   []string{"sunrise mail", "sunrise login", "my sunrise", "sunrise log in", "sunrise sign in", "sunrise account", "sunrise verification"},
		Names:           []string{"sunrise"},
		OfficialDomains: []string{"sunrise.ch"},
	},
	"Swisscom": {
		TitleKeywords:   []string{"bluewin", "swisscom login", "swisscom log in", "swisscom sign in", "swisscom account", "swisscom verification", "my swisscom"},
		Names:           []string{"swisscom"},
		OfficialDomains: []string{"swisscom.ch", "bluewin.ch"},
	},
	"Orange": {
		TitleKeywords:   []string{"orange mail", "mail orange", "espace client orange", "messagerie orange"},
		OfficialDomains: []string{"orange.fr", "orange.com", "orange.es", "orange.be", "orange.pl"},
	},
	"SFR": {
		TitleKeywords:   []string{"sfr mail", "espace client sfr", "sfr login", "sfr log in", "sfr sign in", "sfr account", "sfr verification", "my sfr"},
		Names:           []string{"sfr"},
		OfficialDomains: []string{"sfr.fr"},
	},
	"Bouygues Telecom": {
		TitleKeywords:   []string{"bouygues telecom"},
		OfficialDomains: []string{"bouyguestelecom.fr"},
	},
	"Free Mobile": {
		TitleKeywords:   []string{"free mobile", "freebox"},
		OfficialDomains: []string{"free.fr", "free-mobile.fr"},
	},
	"Vodafone": {
		TitleKeywords:   []string{"vodafone login", "vodafone log in", "vodafone sign in", "vodafone account", "vodafone verification", "my vodafone"},
		Names:           []string{"vodafone"},
		OfficialDomains: []string{"vodafone.com", "vodafone.co.uk", "vodafone.de", "vodafone.it", "vodafone.es", "vodafone.pt", "vodafone.com.tr", "vodafone.nl", "vodafone.ie"},
	},
	"O2": {
		TitleKeywords:   []string{"o2 login", "my o2"},
		OfficialDomains: []string{"o2.co.uk", "o2online.de", "o2.cz"},
	},
	"EE": {
		TitleKeywords:   []string{"my ee", "ee login"},
		OfficialDomains: []string{"ee.co.uk"},
	},
	"BT": {
		TitleKeywords:   []string{"bt email", "my bt", "bt login"},
		OfficialDomains: []string{"bt.com"},
	},
	"Sky": {
		TitleKeywords:   []string{"sky id", "my sky"},
		OfficialDomains: []string{"sky.com"},
	},
	"Virgin Media": {
		TitleKeywords:   []string{"virgin media"},
		OfficialDomains: []string{"virginmedia.com"},
	},
	"Telstra": {
		TitleKeywords:   []string{"telstra login", "telstra log in", "telstra sign in", "telstra account", "telstra verification", "my telstra"},
		Names:           []string{"telstra"},
		OfficialDomains: []string{"telstra.com.au", "telstra.com"},
	},
	"Optus": {
		TitleKeywords:   []string{"optus login", "optus log in", "optus sign in", "optus account", "optus verification", "my optus"},
		Names:           []string{"optus"},
		OfficialDomains: []string{"optus.com.au"},
	},
	"Rogers": {
		TitleKeywords:   []string{"rogers login", "rogers wireless"},
		OfficialDomains: []string{"rogers.com"},
	},
	"Bell": {
		TitleKeywords:   []string{"bell canada", "mybell"},
		OfficialDomains: []string{"bell.ca"},
	},
	"Telus": {
		TitleKeywords:   []string{"telus login", "telus log in", "telus sign in", "telus account", "telus verification", "my telus"},
		Names:           []string{"telus"},
		OfficialDomains: []string{"telus.com"},
	},
	"Movistar": {
		TitleKeywords:   []string{"movistar login", "movistar log in", "movistar sign in", "movistar account", "movistar verification", "my movistar"},
		Names:           []string{"movistar"},
		OfficialDomains: []string{"movistar.es", "movistar.com", "movistar.com.ar", "movistar.cl", "movistar.com.mx", "movistar.com.pe", "movistar.co"},
	},
	"Claro": {
		TitleKeywords:   []string{"claro login", "claro log in", "claro sign in", "claro account", "claro verification", "my claro"},
		Names:           []string{"claro"},
		OfficialDomains: []string{"claro.com.br", "claro.com.co", "claro.com.ar", "claro.com.pe", "claro.cl", "claro.com.ec", "claro.com.gt", "claro.com.do", "claro.com"},
	},
	"Telekom": {
		TitleKeywords:   []string{"telekom login", "t-online", "telekom log in", "telekom sign in", "telekom account", "telekom verification", "my telekom"},
		Names:           []string{"telekom"},
		OfficialDomains: []string{"telekom.de", "t-online.de", "telekom.com"},
	},
	"Proximus": {
		TitleKeywords:   []string{"proximus login", "proximus log in", "proximus sign in", "proximus account", "proximus verification", "my proximus"},
		Names:           []string{"proximus"},
		OfficialDomains: []string{"proximus.be"},
	},
	"KPN": {
		TitleKeywords:   []string{"kpn login", "kpn log in", "kpn sign in", "kpn account", "kpn verification", "my kpn"},
		Names:           []string{"kpn"},
		OfficialDomains: []string{"kpn.com"},
	},
	"Ziggo": {
		TitleKeywords:   []string{"ziggo login", "ziggo log in", "ziggo sign in", "ziggo account", "ziggo verification", "my ziggo"},
		Names:           []string{"ziggo"},
		OfficialDomains: []string{"ziggo.nl"},
	},
	"Telenor": {
		TitleKeywords:   []string{"telenor login", "telenor log in", "telenor sign in", "telenor account", "telenor verification", "my telenor"},
		Names:           []string{"telenor"},
		OfficialDomains: []string{"telenor.com", "telenor.no", "telenor.se", "telenor.dk"},
	},
	"Telia": {
		TitleKeywords:   []string{"telia login", "telia log in", "telia sign in", "telia account", "telia verification", "my telia"},
		Names:           []string{"telia"},
		OfficialDomains: []string{"telia.se", "telia.fi", "telia.no", "telia.dk", "telia.lt", "telia.ee", "teliacompany.com"},
	},
	"Airtel": {
		TitleKeywords:   []string{"airtel login", "airtel log in", "airtel sign in", "airtel account", "airtel verification", "my airtel"},
		Names:           []string{"airtel"},
		OfficialDomains: []string{"airtel.in", "airtel.com"},
	},
	"Jio": {
		TitleKeywords:   []string{"myjio", "jio login", "jio log in", "jio sign in", "jio account", "jio verification", "my jio"},
		Names:           []string{"jio"},
		OfficialDomains: []string{"jio.com"},
	},
	"Globe": {
		TitleKeywords:   []string{"globe telecom", "globeone"},
		OfficialDomains: []string{"globe.com.ph"},
	},
	"Etisalat": {
		TitleKeywords:   []string{"etisalat", "e& uae"},
		OfficialDomains: []string{"etisalat.ae", "eand.com"},
	},
	"STC": {
		TitleKeywords:   []string{"stc pay", "my stc"},
		OfficialDomains: []string{"stc.com.sa", "stcpay.com.sa"},
	},
	"Safaricom": {
		TitleKeywords:   []string{"safaricom login", "safaricom log in", "safaricom sign in", "safaricom account", "safaricom verification", "my safaricom"},
		Names:           []string{"safaricom"},
		OfficialDomains: []string{"safaricom.co.ke"},
	},
	"MTN": {
		TitleKeywords:   []string{"mtn login", "mtn log in", "mtn sign in", "mtn account", "mtn verification", "my mtn"},
		Names:           []string{"mtn"},
		OfficialDomains: []string{"mtn.com", "mtn.co.za", "mtnonline.com", "mtn.ng"},
	},

	// ── Delivery & Logistics (more) ───────────────────────────────────────────
	"Canada Post": {
		TitleKeywords:   []string{"canada post", "postes canada"},
		OfficialDomains: []string{"canadapost-postescanada.ca", "canadapost.ca"},
	},
	"Australia Post": {
		TitleKeywords:   []string{"australia post", "auspost"},
		OfficialDomains: []string{"auspost.com.au"},
	},
	"La Poste": {
		TitleKeywords:   []string{"la poste", "colissimo", "chronopost"},
		OfficialDomains: []string{"laposte.fr", "laposte.net", "colissimo.fr", "chronopost.fr"},
	},
	"DPD": {
		TitleKeywords:   []string{"dpd login", "dpd log in", "dpd sign in", "dpd account", "dpd verification", "my dpd"},
		Names:           []string{"dpd"},
		OfficialDomains: []string{"dpd.com", "dpd.co.uk", "dpd.de", "dpd.fr", "dpd.nl", "dpd.be"},
	},
	"GLS": {
		TitleKeywords:   []string{"gls parcel", "gls tracking", "gls login", "gls log in", "gls sign in", "gls account", "gls verification", "my gls"},
		Names:           []string{"gls"},
		OfficialDomains: []string{"gls-group.eu", "gls-group.com", "gls-us.com"},
	},
	"Evri": {
		TitleKeywords:   []string{"hermes parcel", "evri login", "evri log in", "evri sign in", "evri account", "evri verification", "my evri"},
		Names:           []string{"evri"},
		OfficialDomains: []string{"evri.com", "myhermes.co.uk", "myhermes.de", "hermesworld.com"},
	},
	"PostNL": {
		TitleKeywords:   []string{"postnl login", "postnl log in", "postnl sign in", "postnl account", "postnl verification", "my postnl"},
		Names:           []string{"postnl"},
		OfficialDomains: []string{"postnl.nl", "postnl.com"},
	},
	"bpost": {
		TitleKeywords:   []string{"bpost login", "bpost log in", "bpost sign in", "bpost account", "bpost verification", "my bpost"},
		Names:           []string{"bpost"},
		OfficialDomains: []string{"bpost.be"},
	},
	"Deutsche Post": {
		TitleKeywords:   []string{"deutsche post"},
		OfficialDomains: []string{"deutschepost.de", "dhl.de"},
	},
	"SEUR": {
		TitleKeywords:   []string{"seur", "seur envío", "seur tracking", "seur seguimiento"},
		OfficialDomains: []string{"seur.com", "seur.es", "seur.pt"},
	},
	"Correos": {
		TitleKeywords:   []string{"correos login", "correos log in", "correos sign in", "correos account", "correos verification", "my correos"},
		Names:           []string{"correos"},
		OfficialDomains: []string{"correos.es", "correos.cl", "correos.go.cr"},
	},
	"CTT": {
		TitleKeywords:   []string{"ctt correios"},
		OfficialDomains: []string{"ctt.pt"},
	},
	"An Post": {
		TitleKeywords:   []string{"an post"},
		OfficialDomains: []string{"anpost.com", "anpost.ie"},
	},
	"Swiss Post": {
		TitleKeywords:   []string{"swiss post", "die post", "la poste suisse"},
		OfficialDomains: []string{"post.ch"},
	},
	"PostNord": {
		TitleKeywords:   []string{"postnord"},
		OfficialDomains: []string{"postnord.se", "postnord.dk", "postnord.no", "postnord.fi", "postnord.com"},
	},
	"Japan Post": {
		TitleKeywords:   []string{"japan post", "日本郵便"},
		OfficialDomains: []string{"post.japanpost.jp", "japanpost.jp"},
	},
	"India Post": {
		TitleKeywords:   []string{"india post", "indiapost"},
		OfficialDomains: []string{"indiapost.gov.in"},
	},
	"Aramex": {
		TitleKeywords:   []string{"aramex login", "aramex log in", "aramex sign in", "aramex account", "aramex verification", "my aramex"},
		Names:           []string{"aramex"},
		OfficialDomains: []string{"aramex.com"},
	},
	"J&T Express": {
		TitleKeywords:   []string{"j&t express"},
		OfficialDomains: []string{"jtexpress.com", "jet.co.id", "jtexpress.ph", "jtexpress.my"},
	},
	"Ninja Van": {
		TitleKeywords:   []string{"ninja van", "ninjavan"},
		OfficialDomains: []string{"ninjavan.co"},
	},
	"Purolator": {
		TitleKeywords:   []string{"purolator login", "purolator log in", "purolator sign in", "purolator account", "purolator verification", "my purolator"},
		Names:           []string{"purolator"},
		OfficialDomains: []string{"purolator.com"},
	},
	"Yamato": {
		TitleKeywords:   []string{"yamato transport", "クロネコヤマト"},
		OfficialDomains: []string{"kuronekoyamato.co.jp"},
	},
	"Sagawa": {
		TitleKeywords:   []string{"sagawa", "佐川急便"},
		OfficialDomains: []string{"sagawa-exp.co.jp"},
	},
	"Poczta Polska": {
		TitleKeywords:   []string{"poczta polska"},
		OfficialDomains: []string{"poczta-polska.pl"},
	},
	"Austrian Post": {
		TitleKeywords:   []string{"österreichische post", "post.at"},
		OfficialDomains: []string{"post.at"},
	},
	"Posti": {
		TitleKeywords:   []string{"posti login", "posti log in", "posti sign in", "posti account", "posti verification", "my posti"},
		Names:           []string{"posti"},
		OfficialDomains: []string{"posti.fi"},
	},
	"Correios": {
		TitleKeywords:   []string{"correios login", "correios log in", "correios sign in", "correios account", "correios verification", "my correios"},
		Names:           []string{"correios"},
		OfficialDomains: []string{"correios.com.br"},
	},
	"Delhivery": {
		TitleKeywords:   []string{"delhivery login", "delhivery log in", "delhivery sign in", "delhivery account", "delhivery verification", "my delhivery"},
		Names:           []string{"delhivery"},
		OfficialDomains: []string{"delhivery.com"},
	},

	// ── Government, Tax & Tolls (more) ────────────────────────────────────────
	"GOV.UK": {
		TitleKeywords:   []string{"gov.uk", "dvla", "government gateway"},
		OfficialDomains: []string{"gov.uk"},
	},
	"NHS": {
		TitleKeywords:   []string{"nhs login", "nhs log in", "nhs sign in", "nhs account", "nhs verification", "my nhs"},
		Names:           []string{"nhs"},
		OfficialDomains: []string{"nhs.uk"},
	},
	"Canada Revenue Agency": {
		TitleKeywords:   []string{"canada revenue agency", "cra my account", "agence du revenu du canada"},
		OfficialDomains: []string{"canada.ca", "gc.ca"},
	},
	"Service Canada": {
		TitleKeywords:   []string{"service canada", "gckey"},
		OfficialDomains: []string{"canada.ca", "gc.ca"},
	},
	"ATO": {
		TitleKeywords:   []string{"australian taxation office"},
		OfficialDomains: []string{"ato.gov.au"},
	},
	"myGov": {
		TitleKeywords:   []string{"mygov"},
		OfficialDomains: []string{"my.gov.au", "gov.au"},
	},
	"Service NSW": {
		TitleKeywords:   []string{"service nsw"},
		OfficialDomains: []string{"service.nsw.gov.au", "nsw.gov.au"},
	},
	"Linkt": {
		TitleKeywords:   []string{"linkt"},
		OfficialDomains: []string{"linkt.com.au"},
	},
	"E-ZPass": {
		TitleKeywords:   []string{"e-zpass", "ezpass"},
		OfficialDomains: []string{"e-zpassny.com", "ezpassnj.com", "ezpassva.com", "ezpassmd.com", "ezpassde.com", "ezdrivema.com", "e-zpassiag.com", "paytollo.com", "ezpass.csc.paturnpike.com", "paturnpike.com"},
	},
	"SunPass": {
		TitleKeywords:   []string{"sunpass"},
		OfficialDomains: []string{"sunpass.com"},
	},
	"FasTrak": {
		TitleKeywords:   []string{"fastrak"},
		OfficialDomains: []string{"bayareafastrak.org", "thetollroads.com"},
	},
	"TxTag": {
		TitleKeywords:   []string{"txtag"},
		OfficialDomains: []string{"txtag.org"},
	},
	"USCIS": {
		TitleKeywords:   []string{"uscis login", "uscis log in", "uscis sign in", "uscis account", "uscis verification", "my uscis"},
		Names:           []string{"uscis"},
		OfficialDomains: []string{"uscis.gov"},
	},
	"Impots": {
		TitleKeywords:   []string{"impots.gouv", "impôts.gouv", "dgfip"},
		OfficialDomains: []string{"impots.gouv.fr", "gouv.fr"},
	},
	"Ameli": {
		TitleKeywords:   []string{"assurance maladie", "ameli login", "ameli log in", "ameli sign in", "ameli account", "ameli verification", "my ameli"},
		Names:           []string{"ameli"},
		OfficialDomains: []string{"ameli.fr"},
	},
	"ANTS": {
		TitleKeywords:   []string{"ants.gouv", "agence nationale des titres sécurisés"},
		OfficialDomains: []string{"ants.gouv.fr", "gouv.fr"},
	},
	"CAF": {
		TitleKeywords:   []string{"caf.fr", "allocations familiales"},
		OfficialDomains: []string{"caf.fr"},
	},
	"Agenzia delle Entrate": {
		TitleKeywords:   []string{"agenzia delle entrate"},
		OfficialDomains: []string{"agenziaentrate.gov.it"},
	},
	"INPS": {
		TitleKeywords:   []string{"inps login", "inps log in", "inps sign in", "inps account", "inps verification", "my inps"},
		Names:           []string{"inps"},
		OfficialDomains: []string{"inps.it"},
	},
	"Agencia Tributaria": {
		TitleKeywords:   []string{"agencia tributaria"},
		OfficialDomains: []string{"agenciatributaria.gob.es", "agenciatributaria.es"},
	},
	"DGT": {
		TitleKeywords:   []string{"dirección general de tráfico"},
		OfficialDomains: []string{"dgt.es"},
	},
	"ELSTER": {
		TitleKeywords:   []string{"elster login", "elster log in", "elster sign in", "elster account", "elster verification", "my elster"},
		Names:           []string{"elster"},
		OfficialDomains: []string{"elster.de"},
	},
	"e-Devlet": {
		TitleKeywords:   []string{"e-devlet", "türkiye.gov.tr"},
		OfficialDomains: []string{"turkiye.gov.tr"},
	},
	"MHRS": {
		TitleKeywords:   []string{"mhrs"},
		OfficialDomains: []string{"mhrs.gov.tr"},
	},
	"SAT": {
		TitleKeywords:   []string{"servicio de administración tributaria"},
		OfficialDomains: []string{"sat.gob.mx"},
	},
	"Receita Federal": {
		TitleKeywords:   []string{"receita federal", "gov.br"},
		OfficialDomains: []string{"gov.br"},
	},
	"Gosuslugi": {
		TitleKeywords:   []string{"госуслуги", "gosuslugi"},
		OfficialDomains: []string{"gosuslugi.ru"},
	},
	"EPFO": {
		TitleKeywords:   []string{"epfo", "uan member"},
		OfficialDomains: []string{"epfindia.gov.in"},
	},
	"DigiLocker": {
		TitleKeywords:   []string{"digilocker"},
		OfficialDomains: []string{"digilocker.gov.in"},
	},
	"Passport Seva": {
		TitleKeywords:   []string{"passport seva"},
		OfficialDomains: []string{"passportindia.gov.in"},
	},
	"Singpass": {
		TitleKeywords:   []string{"singpass"},
		OfficialDomains: []string{"singpass.gov.sg"},
	},
	"SARS": {
		TitleKeywords:   []string{"sars efiling"},
		OfficialDomains: []string{"sars.gov.za"},
	},
	"MyGov India": {
		TitleKeywords:   []string{"mygov india"},
		OfficialDomains: []string{"mygov.in"},
	},

	// ── Travel & Transport ────────────────────────────────────────────────────
	"Booking.com": {
		TitleKeywords:   []string{"booking.com"},
		OfficialDomains: []string{"booking.com"},
	},
	"Airbnb": {
		TitleKeywords:   []string{"airbnb login", "airbnb log in", "airbnb sign in", "airbnb account", "airbnb verification", "my airbnb"},
		Names:           []string{"airbnb"},
		OfficialDomains: []string{"airbnb.com", "airbnb.co.uk", "airbnb.fr", "airbnb.de", "airbnb.es", "airbnb.it", "airbnb.ca", "airbnb.com.au", "airbnb.co.in"},
	},
	"Expedia": {
		TitleKeywords:   []string{"expedia login", "expedia log in", "expedia sign in", "expedia account", "expedia verification", "my expedia"},
		Names:           []string{"expedia"},
		OfficialDomains: []string{"expedia.com", "expedia.co.uk", "expedia.ca", "expedia.de", "expedia.fr"},
	},
	"Agoda": {
		TitleKeywords:   []string{"agoda login", "agoda log in", "agoda sign in", "agoda account", "agoda verification", "my agoda"},
		Names:           []string{"agoda"},
		OfficialDomains: []string{"agoda.com"},
	},
	"Trip.com": {
		TitleKeywords:   []string{"trip.com"},
		OfficialDomains: []string{"trip.com", "ctrip.com"},
	},
	"Delta": {
		TitleKeywords:   []string{"delta air lines", "skymiles"},
		OfficialDomains: []string{"delta.com"},
	},
	"United": {
		TitleKeywords:   []string{"united airlines", "mileageplus"},
		OfficialDomains: []string{"united.com"},
	},
	"American Airlines": {
		TitleKeywords:   []string{"american airlines", "aadvantage"},
		OfficialDomains: []string{"aa.com"},
	},
	"Southwest": {
		TitleKeywords:   []string{"southwest airlines", "rapid rewards"},
		OfficialDomains: []string{"southwest.com"},
	},
	"Emirates": {
		TitleKeywords:   []string{"emirates skywards", "emirates airline login", "emirates airline log in", "emirates airline sign in", "emirates airline account", "emirates airline verification", "my emirates airline"},
		Names:           []string{"emirates airline"},
		OfficialDomains: []string{"emirates.com"},
	},
	"Qatar Airways": {
		TitleKeywords:   []string{"qatar airways"},
		OfficialDomains: []string{"qatarairways.com"},
	},
	"Lufthansa": {
		TitleKeywords:   []string{"miles & more", "lufthansa login", "lufthansa log in", "lufthansa sign in", "lufthansa account", "lufthansa verification", "my lufthansa"},
		Names:           []string{"lufthansa"},
		OfficialDomains: []string{"lufthansa.com", "miles-and-more.com"},
	},
	"Air France": {
		TitleKeywords:   []string{"air france", "flying blue"},
		OfficialDomains: []string{"airfrance.com", "airfrance.fr", "flyingblue.com"},
	},
	"KLM": {
		TitleKeywords:   []string{"klm royal dutch"},
		OfficialDomains: []string{"klm.com", "klm.nl"},
	},
	"British Airways": {
		TitleKeywords:   []string{"british airways"},
		OfficialDomains: []string{"britishairways.com", "ba.com"},
	},
	"Ryanair": {
		TitleKeywords:   []string{"ryanair login", "ryanair log in", "ryanair sign in", "ryanair account", "ryanair verification", "my ryanair"},
		Names:           []string{"ryanair"},
		OfficialDomains: []string{"ryanair.com"},
	},
	"easyJet": {
		TitleKeywords:   []string{"easyjet login", "easyjet log in", "easyjet sign in", "easyjet account", "easyjet verification", "my easyjet"},
		Names:           []string{"easyjet"},
		OfficialDomains: []string{"easyjet.com"},
	},
	"Uber": {
		TitleKeywords:   []string{"uber login", "uber log in", "uber sign in", "uber account", "uber verification", "my uber"},
		Names:           []string{"uber"},
		OfficialDomains: []string{"uber.com"},
	},
	"Lyft": {
		TitleKeywords:   []string{"lyft login", "lyft log in", "lyft sign in", "lyft account", "lyft verification", "my lyft"},
		Names:           []string{"lyft"},
		OfficialDomains: []string{"lyft.com"},
	},
	"Grab": {
		TitleKeywords:   []string{"grab login", "grab log in", "grab sign in", "grab account", "grab verification", "my grab"},
		Names:           []string{"grab"},
		OfficialDomains: []string{"grab.com"},
	},
	"Gojek": {
		TitleKeywords:   []string{"gojek", "gopay"},
		OfficialDomains: []string{"gojek.com", "gopay.co.id"},
	},

	// ── Gaming (more) ─────────────────────────────────────────────────────────
	"Minecraft": {
		TitleKeywords:   []string{"minecraft login", "minecraft log in", "minecraft sign in", "minecraft account", "minecraft verification", "my minecraft"},
		Names:           []string{"minecraft"},
		OfficialDomains: []string{"minecraft.net", "mojang.com", "microsoft.com", "xbox.com"},
	},
	"Garena": {
		TitleKeywords:   []string{"garena login", "garena log in", "garena sign in", "garena account", "garena verification", "my garena", "free fire login", "free fire log in", "free fire sign in", "free fire account", "free fire verification", "my free fire"},
		Names:           []string{"garena", "free fire"},
		OfficialDomains: []string{"garena.com", "ff.garena.com"},
	},
	"PUBG": {
		TitleKeywords:   []string{"pubg login", "pubg log in", "pubg sign in", "pubg account", "pubg verification", "my pubg", "pubg mobile login", "pubg mobile log in", "pubg mobile sign in", "pubg mobile account", "pubg mobile verification", "my pubg mobile"},
		Names:           []string{"pubg", "pubg mobile"},
		OfficialDomains: []string{"pubg.com", "pubgmobile.com"},
	},
	"Mobile Legends": {
		TitleKeywords:   []string{"mobile legends login", "mobile legends log in", "mobile legends sign in", "mobile legends account", "mobile legends verification", "my mobile legends", "mlbb login", "mlbb log in", "mlbb sign in", "mlbb account", "mlbb verification", "my mlbb"},
		Names:           []string{"mobile legends", "mlbb"},
		OfficialDomains: []string{"mobilelegends.com", "moonton.com"},
	},
	"HoYoverse": {
		TitleKeywords:   []string{"honkai star rail", "mihoyo", "genshin impact login", "genshin impact log in", "genshin impact sign in", "genshin impact account", "genshin impact verification", "my genshin impact", "hoyoverse login", "hoyoverse log in", "hoyoverse sign in", "hoyoverse account", "hoyoverse verification", "my hoyoverse"},
		Names:           []string{"genshin impact", "hoyoverse"},
		OfficialDomains: []string{"hoyoverse.com", "hoyolab.com", "mihoyo.com"},
	},
	"Supercell": {
		TitleKeywords:   []string{"supercell id login", "supercell id log in", "supercell id sign in", "supercell id account", "supercell id verification", "my supercell id", "clash of clans login", "clash of clans log in", "clash of clans sign in", "clash of clans account", "clash of clans verification", "my clash of clans", "brawl stars login", "brawl stars log in", "brawl stars sign in", "brawl stars account", "brawl stars verification", "my brawl stars", "clash royale login", "clash royale log in", "clash royale sign in", "clash royale account", "clash royale verification", "my clash royale"},
		Names:           []string{"supercell id", "clash of clans", "brawl stars", "clash royale"},
		OfficialDomains: []string{"supercell.com"},
	},
	"Ubisoft": {
		TitleKeywords:   []string{"ubisoft connect", "ubisoft login", "ubisoft log in", "ubisoft sign in", "ubisoft account", "ubisoft verification", "my ubisoft"},
		Names:           []string{"ubisoft"},
		OfficialDomains: []string{"ubisoft.com"},
	},
	"Rockstar Games": {
		TitleKeywords:   []string{"rockstar games social club", "rockstar games login", "rockstar games log in", "rockstar games sign in", "rockstar games account", "rockstar games verification", "my rockstar games"},
		Names:           []string{"rockstar games"},
		OfficialDomains: []string{"rockstargames.com"},
	},
	"Ankama": {
		TitleKeywords:   []string{"dofus login", "dofus log in", "dofus sign in", "dofus account", "dofus verification", "my dofus", "ankama login", "ankama log in", "ankama sign in", "ankama account", "ankama verification", "my ankama"},
		Names:           []string{"dofus", "ankama"},
		OfficialDomains: []string{"ankama.com", "dofus.com"},
	},
	"Valve": {
		TitleKeywords:   []string{"steam guard", "steam support"},
		OfficialDomains: []string{"steampowered.com", "steamcommunity.com"},
	},
	"CS2 Skins": {
		TitleKeywords:   []string{"csgo skins", "cs2 skins", "csgo case"},
		OfficialDomains: []string{"steampowered.com", "steamcommunity.com", "counter-strike.net"},
	},

	// ── Streaming & Media (more) ──────────────────────────────────────────────
	"DAZN": {
		TitleKeywords:   []string{"dazn login", "dazn log in", "dazn sign in", "dazn account", "dazn verification", "my dazn"},
		Names:           []string{"dazn"},
		OfficialDomains: []string{"dazn.com"},
	},
	"Canal+": {
		TitleKeywords:   []string{"canal+ login", "canal+ log in", "canal+ sign in", "canal+ account", "canal+ verification", "my canal+", "mycanal login", "mycanal log in", "mycanal sign in", "mycanal account", "mycanal verification", "my mycanal"},
		Names:           []string{"canal+", "mycanal"},
		OfficialDomains: []string{"canalplus.com", "canal-plus.com"},
	},
	"Deezer": {
		TitleKeywords:   []string{"deezer login", "deezer log in", "deezer sign in", "deezer account", "deezer verification", "my deezer"},
		Names:           []string{"deezer"},
		OfficialDomains: []string{"deezer.com"},
	},
	"Audible": {
		TitleKeywords:   []string{"audible login", "audible log in", "audible sign in", "audible account", "audible verification", "my audible"},
		Names:           []string{"audible"},
		OfficialDomains: []string{"audible.com", "audible.co.uk", "audible.de", "audible.in"},
	},
	"Apple Music": {
		TitleKeywords:   []string{"apple music"},
		OfficialDomains: []string{"apple.com", "music.apple.com"},
	},
	"YouTube Premium": {
		TitleKeywords:   []string{"youtube premium"},
		OfficialDomains: []string{"youtube.com", "google.com"},
	},
	"Hotstar": {
		TitleKeywords:   []string{"jiohotstar", "disney+ hotstar", "hotstar login", "hotstar log in", "hotstar sign in", "hotstar account", "hotstar verification", "my hotstar"},
		Names:           []string{"hotstar"},
		OfficialDomains: []string{"hotstar.com", "jiohotstar.com"},
	},
	"Zee5": {
		TitleKeywords:   []string{"zee5 login", "zee5 log in", "zee5 sign in", "zee5 account", "zee5 verification", "my zee5"},
		Names:           []string{"zee5"},
		OfficialDomains: []string{"zee5.com"},
	},

	// ── Insurance & Health ────────────────────────────────────────────────────
	"Geico": {
		TitleKeywords:   []string{"geico login", "geico log in", "geico sign in", "geico account", "geico verification", "my geico"},
		Names:           []string{"geico"},
		OfficialDomains: []string{"geico.com"},
	},
	"State Farm": {
		TitleKeywords:   []string{"state farm login", "state farm log in", "state farm sign in", "state farm account", "state farm verification", "my state farm"},
		Names:           []string{"state farm"},
		OfficialDomains: []string{"statefarm.com"},
	},
	"Allstate": {
		TitleKeywords:   []string{"allstate login", "allstate log in", "allstate sign in", "allstate account", "allstate verification", "my allstate"},
		Names:           []string{"allstate"},
		OfficialDomains: []string{"allstate.com"},
	},
	"AXA": {
		TitleKeywords:   []string{"axa login", "my axa", "axa log in", "axa sign in", "axa account", "axa verification"},
		Names:           []string{"axa"},
		OfficialDomains: []string{"axa.com", "axa.fr", "axa.de", "axa.co.uk", "axa.ch", "axa.be"},
	},
	"Allianz": {
		TitleKeywords:   []string{"allianz login", "allianz log in", "allianz sign in", "allianz account", "allianz verification", "my allianz"},
		Names:           []string{"allianz"},
		OfficialDomains: []string{"allianz.com", "allianz.de", "allianz.fr", "allianz.it", "allianz.ch"},
	},
	"UnitedHealthcare": {
		TitleKeywords:   []string{"unitedhealthcare"},
		OfficialDomains: []string{"uhc.com", "myuhc.com"},
	},
	"Kaiser Permanente": {
		TitleKeywords:   []string{"kaiser permanente"},
		OfficialDomains: []string{"kp.org"},
	},
	"Blue Cross": {
		TitleKeywords:   []string{"blue cross blue shield", "anthem blue cross"},
		OfficialDomains: []string{"bcbs.com", "anthem.com"},
	},
}
