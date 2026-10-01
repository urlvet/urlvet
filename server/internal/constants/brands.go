package constants

type BrandEntry struct {
	TitleKeywords   []string
	OfficialDomains []string
}

var HighValueBrands = map[string]BrandEntry{

	// ── Email Platforms ──────────────────────────────────────────────────────
	"Gmail": {
		TitleKeywords:   []string{"gmail", "google mail"},
		OfficialDomains: []string{"gmail.com", "google.com"},
	},
	"Google": {
		TitleKeywords:   []string{"google account", "google drive", "google workspace", "google sign in"},
		OfficialDomains: []string{"google.com", "gmail.com"},
	},
	"Microsoft Outlook": {
		TitleKeywords:   []string{"outlook", "hotmail", "live mail", "microsoft mail"},
		OfficialDomains: []string{"outlook.com", "hotmail.com", "live.com", "microsoft.com", "microsoftonline.com"},
	},
	"Microsoft": {
		TitleKeywords:   []string{"microsoft", "office 365", "microsoft 365", "onedrive", "sharepoint", "microsoft teams", "azure"},
		OfficialDomains: []string{"microsoft.com", "live.com", "outlook.com", "hotmail.com", "office.com", "office365.com", "microsoftonline.com", "azure.com"},
	},
	"Yahoo Mail": {
		TitleKeywords:   []string{"yahoo", "yahoo mail", "yahoo sign in"},
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
		OfficialDomains: []string{"facebook.com", "fb.com"},
	},
	"Instagram": {
		TitleKeywords:   []string{"instagram login"},
		OfficialDomains: []string{"instagram.com"},
	},
	"Twitter": {
		TitleKeywords:   []string{"x login", "x sign in", "sign in to x"},
		OfficialDomains: []string{"twitter.com", "x.com"},
	},
	"LinkedIn": {
		TitleKeywords:   []string{"linkedin sign in", "linkedin login"},
		OfficialDomains: []string{"linkedin.com", "lnkd.in"},
	},
	"TikTok": {
		TitleKeywords:   []string{"tiktok login"},
		OfficialDomains: []string{"tiktok.com"},
	},
	"Snapchat": {
		TitleKeywords:   []string{"snapchat", "snapchat login"},
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
	},
	"Tumblr": {
		TitleKeywords:   []string{"tumblr", "tumblr login"},
		OfficialDomains: []string{"tumblr.com"},
	},

	// ── Messaging & Chat ──────────────────────────────────────────────────────
	"WhatsApp": {
		TitleKeywords:   []string{"whatsapp login", "whatsapp web"},
		OfficialDomains: []string{"whatsapp.com", "web.whatsapp.com"},
	},
	"Telegram": {
		TitleKeywords:   []string{"telegram login", "telegram web"},
		OfficialDomains: []string{"telegram.org", "web.telegram.org", "t.me"},
	},
	"Signal": {
		TitleKeywords:   []string{"signal", "signal messenger"},
		OfficialDomains: []string{"signal.org"},
	},
	"Discord": {
		TitleKeywords:   []string{"discord login", "discord sign in"},
		OfficialDomains: []string{"discord.com", "discordapp.com"},
	},
	"Messenger": {
		TitleKeywords:   []string{"messenger", "facebook messenger"},
		OfficialDomains: []string{"messenger.com", "facebook.com"},
	},
	"Skype": {
		TitleKeywords:   []string{"skype", "skype sign in"},
		OfficialDomains: []string{"skype.com", "microsoft.com"},
	},
	"Slack": {
		TitleKeywords:   []string{"slack", "slack sign in", "slack login"},
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
		TitleKeywords:   []string{"amazon", "amazon sign in", "amazon prime", "aws sign in"},
		OfficialDomains: []string{"amazon.com", "amazon.co.uk", "amazon.de", "amazon.in", "amazon.ca", "amazon.com.au", "amazon.fr", "amazon.es", "amazon.it", "amazon.co.jp", "aws.amazon.com"},
	},
	"eBay": {
		TitleKeywords:   []string{"ebay", "ebay sign in", "ebay login"},
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
		TitleKeywords:   []string{"stripe", "stripe login"},
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
		TitleKeywords:   []string{"wise", "transferwise"},
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
		TitleKeywords:   []string{"chase", "chase bank", "chase sign in"},
		OfficialDomains: []string{"chase.com"},
	},
	"Bank of America": {
		TitleKeywords:   []string{"bank of america"},
		OfficialDomains: []string{"bankofamerica.com"},
	},
	"Wells Fargo": {
		TitleKeywords:   []string{"wells fargo"},
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
		TitleKeywords:   []string{"santander"},
		OfficialDomains: []string{"santander.co.uk", "santander.com"},
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
		OfficialDomains: []string{"binance.com", "binance.us"},
	},
	"Coinbase": {
		TitleKeywords:   []string{"coinbase", "coinbase login"},
		OfficialDomains: []string{"coinbase.com"},
	},
	"Kraken": {
		TitleKeywords:   []string{"kraken", "kraken login"},
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
		TitleKeywords:   []string{"gemini crypto", "gemini exchange"},
		OfficialDomains: []string{"gemini.com"},
	},
	"MetaMask": {
		TitleKeywords:   []string{"metamask", "metamask wallet"},
		OfficialDomains: []string{"metamask.io"},
	},
	"Phantom": {
		TitleKeywords:   []string{"phantom wallet"},
		OfficialDomains: []string{"phantom.app"},
	},
	"Trust Wallet": {
		TitleKeywords:   []string{"trust wallet"},
		OfficialDomains: []string{"trustwallet.com"},
	},
	"Ledger": {
		TitleKeywords:   []string{"ledger", "ledger live"},
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
		TitleKeywords:   []string{"steam", "steam login", "steam sign in"},
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
		OfficialDomains: []string{"roblox.com"},
	},
	"Blizzard": {
		TitleKeywords:   []string{"blizzard", "battle.net", "world of warcraft", "overwatch"},
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
		OfficialDomains: []string{"netflix.com"},
	},
	"Spotify": {
		TitleKeywords:   []string{"spotify", "spotify login", "spotify account"},
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
		TitleKeywords:   []string{"peacock", "peacock tv"},
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
		OfficialDomains: []string{"adobe.com", "adobeid.services"},
	},
	"Dropbox": {
		TitleKeywords:   []string{"dropbox", "dropbox sign in"},
		OfficialDomains: []string{"dropbox.com"},
	},
	"GitHub": {
		TitleKeywords:   []string{"github", "github sign in"},
		OfficialDomains: []string{"github.com", "githubusercontent.com"},
	},
	"Zoom": {
		TitleKeywords:   []string{"zoom", "zoom sign in", "zoom meeting"},
		OfficialDomains: []string{"zoom.us"},
	},
	"DocuSign": {
		TitleKeywords:   []string{"docusign"},
		OfficialDomains: []string{"docusign.com", "docusign.net"},
	},
	"Shopify": {
		TitleKeywords:   []string{"shopify", "shopify login"},
		OfficialDomains: []string{"shopify.com", "myshopify.com"},
	},
	"Notion": {
		TitleKeywords:   []string{"notion", "notion login"},
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
		TitleKeywords:   []string{"indeed", "indeed sign in"},
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
		OfficialDomains: []string{"fedex.com"},
	},
	"UPS": {
		TitleKeywords:   []string{"ups delivery", "ups tracking"},
		OfficialDomains: []string{"ups.com"},
	},
	"DHL": {
		TitleKeywords:   []string{"dhl", "dhl tracking"},
		OfficialDomains: []string{"dhl.com", "dhl.de"},
	},
	"USPS": {
		TitleKeywords:   []string{"usps", "usps tracking", "united states postal"},
		OfficialDomains: []string{"usps.com"},
	},
	"Royal Mail": {
		TitleKeywords:   []string{"royal mail", "royal mail tracking"},
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
		TitleKeywords:   []string{"social security", "ssa login", "my social security"},
		OfficialDomains: []string{"ssa.gov"},
	},
	"Medicare": {
		TitleKeywords:   []string{"medicare", "medicare login"},
		OfficialDomains: []string{"medicare.gov"},
	},
}
