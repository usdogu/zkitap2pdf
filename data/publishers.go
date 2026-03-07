package data

import "github.com/usdogu/zkitap2pdf/types"

type PublisherEntry struct {
	Name      string
	Publisher types.Publisher
}

var Publishers = []PublisherEntry{
	{"Açı Yayınları", types.Publisher{
		DataURL:   "http://dijital.aciyayinlari.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "21_11_22",
	}},

	{"Ankara Yayıncılık", types.Publisher{
		DataURL:   "https://onlineankarayayincilik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "6_13_26",
	}},

	{"Arı Yayıncılık", types.Publisher{
		DataURL:   "https://akillidefter.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "7_14_21",
	}},

	{"ATA", types.Publisher{
		DataURL:   "http://www.ataogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "1_2_6",
	}},

	{"Benim Hocam", types.Publisher{
		DataURL:   "http://www.benimhocamdijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "23_11_22",
	}},

	{"Berkay Yayıncılık", types.Publisher{
		DataURL:   "http://berkayokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "26_7_14",
	}},

	{"Çalışkan Arı", types.Publisher{
		DataURL:   "http://www.caliskanarizkitap.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "6_15_7",
	}},

	{"Çanta Yayıncılık", types.Publisher{
		DataURL:   "http://www.cantadaakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "2_12_24",
	}},

	{"Endemik Yayınları", types.Publisher{
		DataURL:   "http://akillitahta.endemikyayinlari.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "22_1_2",
	}},

	{"Fi Yayınları", types.Publisher{
		DataURL:   "http://www.fiakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "7_3_6",
	}},

	{"Final Kurumsal Yayinlari", types.Publisher{
		DataURL:   "http://finaldijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "2_1_21",
	}},

	{"Gizli Yayınları", types.Publisher{
		DataURL:   "http://gizliokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "1_21_6",
	}},

	{"Günay Yayınları", types.Publisher{
		DataURL:   "http://www.gunayportal.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "9_4_8",
	}},

	{"Hiper Zeka", types.Publisher{
		DataURL:   "http://www.hiperzekadijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "8_4_8",
	}},

	{"İşler", types.Publisher{
		DataURL:   "http://yeni.isleronline.com/controller/zkitap/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "28_15_5",
	}},

	{"Limit", types.Publisher{
		DataURL:   "http://kurumsal.limitheryerde.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "4_17_5",
	}},

	{"Model Eğitim", types.Publisher{
		DataURL:   "http://www.modelogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "7_14_22",
	}},

	{"More And More", types.Publisher{
		DataURL:   "http://zkitap.kurmayokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "1_4_8",
	}},

	{"Mozaik Yayınları", types.Publisher{
		DataURL:   "http://mozaikakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "23_12_24",
	}},

	{"Orijinal Matematik", types.Publisher{
		DataURL:   "http://yeni.isleronline.com/controller/zkitap/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "28_15_5",
	}},

	{"Paraf Akademi", types.Publisher{
		DataURL:   "https://kutuphane.parafakademi.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "2_26_3",
	}},

	{"PRF Yayınları", types.Publisher{
		DataURL:   "http://prfkutuphane.prfyayinlari.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "8_1_2",
	}},

	{"Toprak Yayıncılık", types.Publisher{
		DataURL:   "http://toprakogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "9_5_10",
	}},

	{"ÜçDörtBeş Yayınları", types.Publisher{
		DataURL:   "http://www.345dijitalicerik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "1_10_20",
	}},

	{"Üçgen", types.Publisher{
		DataURL:   "http://www.e-ucgen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "17_18_20",
	}},

	{"Ünlüler Karması", types.Publisher{
		DataURL:   "http://unlulerakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "26_12_24",
	}},

	{"X Yayıncılık", types.Publisher{
		DataURL:   "http://onlinexyayincilik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
		KeyPrefix: "2_1_3",
	}},
}
