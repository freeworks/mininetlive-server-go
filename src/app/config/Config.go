package config

const (
	PublicKey   string = "enqyjAgoDAQm0mx6A/xk8eyxEuEJWK+LQ6n258NtsT6lARMyF+YFgA=="
	PrivateKey  string = "2e3da80f079d3362f504a5db3776a9cd41feeea2"
	RtmpPath    string = "rtmp://vlive3.rtmp.cdn.ucloud.com.cn/ucloud/%s"
	HlsPath     string = "http://vlive3.hls.cdn.ucloud.com.cn/ucloud/%s/playlist.m3u8"
	PushPath    string = "rtmp://publish3.cdn.ucloud.com.cn/ucloud/%s"
	VedioBucket string = "mininetlive123.ufile.ucloud.com.cn"
	ImgDir      string = "/tmp/mininetlive/img/"
	LogDir      string = "/tmp/mininetlive/log/"
	DeductPercent1  float64 = 0.1;
	DeductPercent2  float64 = 0.07;
	DeductPercent3  float64 = 0.03;
)
