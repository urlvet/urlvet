package constants

// Highly trustable TLDs which are only given to trusted entities after proper verification
var TrustedTLDs = map[string]struct{}{
	"aero":        {},
	"bank":        {},
	"coop":        {},
	"creditunion": {},
	"insurance":   {},
	"int":         {},
	"pharmacy":    {},
	"post":        {},
	"museum":      {},
	"resbank":     {},

	// Restricted by the Reserve Bank of India to regulated entities
	"bank.in": {}, //  Banks
	"fin.in":  {}, //  Non-bank financial companies

	// Academic and Educational Entities
	"ac.at":  {}, //  Austria
	"ac.bd":  {}, //  Bangladesh
	"ac.be":  {}, //  Belgium
	"ac.bw":  {}, //  Botswana
	"ac.cn":  {}, //  China
	"ac.cr":  {}, //  Costa Rica
	"ac.cy":  {}, //  Cyprus
	"ac.fj":  {}, //  Fiji
	"ac.in":  {}, //  India
	"ac.id":  {}, //  Indonesia
	"ac.ir":  {}, //  Iran
	"ac.il":  {}, //  Israel
	"ac.jp":  {}, //  Japan
	"ac.ke":  {}, //  Kenya
	"ac.ma":  {}, //  Morocco
	"ac.nz":  {}, //  New Zealand
	"ac.kp":  {}, //  North Korea
	"aca.kp": {}, //  North Korea
	"ac.pg":  {}, //  Papua New Guinea
	"ac.rw":  {}, //  Rwanda
	"ac.rs":  {}, //  Serbia
	"ac.za":  {}, //  South Africa
	"ac.kr":  {}, //  South Korea
	"ac.ss":  {}, //  South Sudan
	"ac.lk":  {}, //  Sri Lanka
	"ac.tz":  {}, //  Tanzania
	"ac.th":  {}, //  Thailand
	"ac.ug":  {}, //  Uganda
	"ac.uk":  {}, //  United Kingdom
	"ac.ae":  {}, //  United Arab Emirates
	"ac.zm":  {}, //  Zambia
	"ac.zw":  {}, //  Zimbabwe

	"edu":    {}, //  America
	"edu.ag": {}, //  Antigua & Barbuda
	"edu.ar": {}, //  Argentina
	"edu.au": {}, //  Australia
	"edu.az": {}, //  Azerbaijan
	"edu.bd": {}, //  Bangladesh
	"edu.br": {}, //  Brazil
	"edu.bn": {}, //  Brunei
	"edu.cn": {}, //  China
	"edu.co": {}, //  Colombia
	"edu.dj": {}, //  Djibouti
	"edu.ec": {}, //  Ecuador
	"edu.eg": {}, //  Egypt
	"edu.sv": {}, //  El Salvador
	"edu.er": {}, //  Eritrea
	"edu.ee": {}, //  Estonia
	"edu.et": {}, //  Ethiopia
	"edu.fi": {}, //  Finland
	"edu.gh": {}, //  Ghana
	"edu.gr": {}, //  Greece
	"edu.gt": {}, //  Guatemala
	"edu.hk": {}, //  Hong Kong
	"edu.it": {}, //  Italy
	"edu.in": {}, //  India
	"edu.jm": {}, //  Jamaica
	"edu.jo": {}, //  Jordan
	"edu.kz": {}, //  Kazakhstan
	"edu.lb": {}, //  Lebanon
	"edu.ly": {}, //  Libya
	"edu.mo": {}, //  Macau
	"edu.my": {}, //  Malaysia
	"edu.mt": {}, //  Malta
	"edu.mx": {}, //  Mexico
	"edu.mm": {}, //  Myanmar
	"edu.np": {}, //  Nepal
	"edu.ni": {}, //  Nicaragua
	"edu.ng": {}, //  Nigeria
	"edu.kp": {}, //  North Korea
	"edu.om": {}, //  Oman
	"edu.pk": {}, //  Pakistan
	"edu.pe": {}, //  Peru
	"edu.ph": {}, //  Philippines
	"edu.pl": {}, //  Poland
	"edu.qa": {}, //  Qatar
	"edu.sa": {}, //  Saudi Arabia
	"edu.rs": {}, //  Serbia
	"edu.sg": {}, //  Singapore
	"edu.so": {}, //  Somalia
	"edu.za": {}, //  South Africa
	"edu.es": {}, //  Spain
	"edu.lk": {}, //  Sri Lanka
	"edu.sd": {}, //  Sudan
	"edu.tw": {}, //  Taiwan
	"edu.tr": {}, //  Turkey
	"edu.ua": {}, //  Ukraine
	"edu.uy": {}, //  Uruguay
	"edu.vn": {}, //  Vietnam

	"ernet.in": {}, //  India
	"res.in":   {}, //  India

	// Goverment Entities
	"gov":           {}, //  America
	"gov.af":        {}, //  Afghanistan
	"gov.al":        {}, //  Albania
	"gov.dz":        {}, //  Algeria
	"gov.ad":        {}, //  Andorra
	"gov.ao":        {}, //  Angola
	"gov.am":        {}, //  Armenia
	"gov.aw":        {}, //  Aruba
	"gob.ar":        {}, //  Argentina
	"gov.ar":        {}, //  Argentina
	"gov.au":        {}, //  Australia
	"gov.ax":        {}, //  Åland Islands
	"gov.az":        {}, //  Azerbaijan
	"gov.bs":        {}, //  Bahamas
	"gov.bd":        {}, //  Bangladesh
	"gov.bb":        {}, //  Barbados
	"gov.by":        {}, //  Belarus
	"gov.be":        {}, //  Belgium
	"gov.bg":        {}, //  Bulgaria
	"gov.ba":        {}, //  Bosnia & Herzegovina
	"gov.br":        {}, //  Brazil
	"gob.cl":        {}, //  Chile
	"gov.cl":        {}, //  Chile
	"gov.cn":        {}, //  Mainland China
	"gov.hk":        {}, //  Hong Kong
	"gov.mo":        {}, //  Macau
	"gov.co":        {}, //  Colombia
	"gov.cy":        {}, //  Cyprus
	"gov.cz":        {}, //  Czechia
	"gov.eg":        {}, //  Egypt
	"gob.sv":        {}, //  El Salvador
	"gov.gr":        {}, //  Greece
	"gov.fi":        {}, //  Finland
	"gouv.fr":       {}, //  France
	"gov.hu":        {}, //  Hungary
	"gov.in":        {}, //  India
	"go.id":         {}, //  Indonesia
	"gov.ir":        {}, //  Iran
	"gov.iq":        {}, //  Iraq
	"gov.krd":       {}, //  Kurdistan Region
	"gov.ie":        {}, //  Ireland
	"gov.il":        {}, //  Israel
	"gov.it":        {}, //  Italy
	"go.jp":         {}, //  Japan
	"gov.kz":        {}, //  Kazakhstan
	"go.ke":         {}, //  Kenya
	"gov.lv":        {}, //  Latvia
	"gov.lt":        {}, //  Lithuania
	"gov.my":        {}, //  Malaysia
	"gov.mt":        {}, //  Malta
	"gob.mx":        {}, //  Mexico
	"gov.md":        {}, //  Moldova
	"gov.ma":        {}, //  Morocco
	"gov.mm":        {}, //  Myanmar
	"gov.np":        {}, //  Nepal
	"govt.nz":       {}, //  NewZealand
	"gov.ng":        {}, //  Nigeria
	"gob.pe":        {}, //  Peru
	"gov.pk":        {}, //  Pakistan
	"gov.ph":        {}, //  Philippines
	"gov.pl":        {}, //  Poland
	"gov.pt":        {}, //  Portugal
	"gov.ro":        {}, //  Romania
	"gov.ru":        {}, //  Russia
	"gov.sn":        {}, //  Senegal
	"gov.sg":        {}, //  Singapore
	"gov.sk":        {}, //  Slovakia
	"gov.si":        {}, //  Slovenia
	"go.kr":         {}, //  South Korea
	"gob.es":        {}, //  Spain
	"gov.lk":        {}, //  Sri Lanka
	"gov.se":        {}, //  Sweden
	"admin.ch":      {}, //  Switzerland
	"gov.tw":        {}, //  Taiwan
	"go.th":         {}, //  Thailand
	"gov.tt":        {}, //  Trinidad & Tobago
	"gov.tr":        {}, //  Turkey
	"gov.ua":        {}, //  Ukraine
	"gov.uk":        {}, //  United Kingdom
	"gov.scot":      {}, //  Scotland
	"gov.wales":     {}, //  Wales
	"gov.gg":        {}, //  Guernsey
	"gov.je":        {}, //  Jersey
	"gov.im":        {}, //  Isle of Man
	"gov.ai":        {}, //  Anguilla
	"gov.bm":        {}, //  Bermuda
	"gov.vg":        {}, //  British Virgin Islands
	"gov.ky":        {}, //  Cayman Islands
	"gov.fk":        {}, //  Falkland Islands
	"government.pn": {}, //  Pitcairn Islands
	"gov.tc":        {}, //  Turks & Caicos Islands
	"gub.uy":        {}, //  Uruguay
	"gob.ve":        {}, //  Venezuela
	"gov.vn":        {}, //  Vietnam

	// Canada Federal Government
	"gc.ca":     {},
	"canada.ca": {},
	// Canada Provinces
	"gov.ab.ca":     {}, //  Alberta
	"gov.bc.ca":     {}, //  British Columbia
	"gov.mb.ca":     {}, //  Manitoba
	"gnb.ca":        {}, //  New Brunswick
	"gov.nl.ca":     {}, //  Newfoundland and Labrador
	"novascotia.ca": {}, //  Nova Scotia
	"ontario.ca":    {}, //  Ontario
	"gov.pe.ca":     {}, //  Prince Edward Island
	"gouv.qc.ca":    {}, //  Quebec (French)
	"gov.sk.ca":     {}, //  Saskatchewan
	// Canada Territories
	"gov.nt.ca": {}, //  Northwest Territories
	"gov.nu.ca": {}, //  Nunavut
	"gov.yk.ca": {}, //  Yukon

	// Military
	"mil":        {}, //  United States
	"mod.uk":     {}, //  United Kingdom
	"mil.uk":     {}, //  United Kingdom
	"mod.gov.in": {}, //  India
	"mil.in":     {}, //  India
	"mil.kr":     {}, //  South Korea
	"mil.kz":     {}, //  Kazakhstan
	"mil.jo":     {}, //  Jordan
	"mil.kh":     {}, //  Cambodia
	"mil.kw":     {}, //  Kuwait
	"mil.lv":     {}, //  Latvia
	"mil.my":     {}, //  Malaysia
	"mil.mg":     {}, //  Madagascar
	"mil.mn":     {}, //  Mongolia
	"mil.mz":     {}, //  Mozambique
	"mil.tz":     {}, //  Tanzania

	// National Registries & Internet Authorities
	"nic.in": {}, //  Reserved by National Informatics Centre of India
	"nic.uk": {}, //  Reserved for Nominet UK for internal registry operations
	"nic.mx": {}, //  Reserved for NIC Mexico
	"nic.br": {}, //  Operated by NIC.br, the registry for Brazil and infrastructure control
	"nic.sg": {}, //  Reserved by SGNIC, Singapore’s domain registry
	"nic.id": {}, //  Reserved for Indonesia Network Information Centre (IDNIC)
	"nic.jp": {}, //  Reserved by JPNIC, Japan Network Information Center
	"nic.kr": {}, //  Reserved by KISA / Korea NIC (South Korea)
	"nic.vn": {}, //  Reserved for VNNIC, Vietnam Network Information Center
	"nic.cn": {}, //  Reserved by CNNIC, China Internet Network Information Center

}
