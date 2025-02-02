package presets

import (
	"github.com/plandem/curl-impersonate/types"
	"math/rand"
)

type Preset struct {
	*types.Headers
	*types.Flags
}

type PresetFn func() Preset

func Default() Preset {
	return Preset{types.NewHeaders(), types.NewFlags()}
}

func Random() Preset {
	presets := []PresetFn{
		Chrome99Android,
		Chrome99,
		Chrome100,
		Chrome101,
		Chrome104,
		Chrome107,
		Chrome110,
		Chrome116,
		Edge99,
		Edge101,
		Safari153,
		Safari155,
	}
	fn := presets[rand.Intn(len(presets))]
	return fn()
}

func Chrome99() Preset {
	h := types.NewHeaders(
		types.Header("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`),
		types.Header("sec-ch-ua-mobile", `?0`),
		types.Header("sec-ch-ua-platform", `"Windows"`),
		types.Header("Upgrade-Insecure-Requests", `1`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36`),
		types.Header("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`),
		types.Header("Sec-Fetch-Site", `none`),
		types.Header("Sec-Fetch-Mode", `navigate`),
		types.Header("Sec-Fetch-User", `?1`),
		types.Header("Sec-Fetch-Dest", `document`),
		types.Header("Accept-Encoding", `gzip, deflate, br`),
		types.Header("Accept-Language", `en-US,en;q=0.9`),
	)
	f := types.NewFlags(
		types.Flag("ciphers", "TLS_AES_128_GCM_SHA256,TLS_AES_256_GCM_SHA384,TLS_CHACHA20_POLY1305_SHA256,ECDHE-ECDSA-AES128-GCM-SHA256,ECDHE-RSA-AES128-GCM-SHA256,ECDHE-ECDSA-AES256-GCM-SHA384,ECDHE-RSA-AES256-GCM-SHA384,ECDHE-ECDSA-CHACHA20-POLY1305,ECDHE-RSA-CHACHA20-POLY1305,ECDHE-RSA-AES128-SHA,ECDHE-RSA-AES256-SHA,AES128-GCM-SHA256,AES256-GCM-SHA384,AES128-SHA,AES256-SHA"),
		types.Flag("http2", true),
		types.Flag("compressed", true),
		types.Flag("tlsv1.2", true),
		types.Flag("alps", true),
		types.Flag("cert-compression", "brotli"),
	)

	return Preset{h, f}
}

func Chrome99Android() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua-mobile", `?1`),
		types.Header("sec-ch-ua-platform", `"Android"`),
		types.Header("User-Agent", `Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.58 Mobile Safari/537.36`),
	)

	return preset
}

func Chrome100() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="100", "Google Chrome";v="100"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36`),
	)

	return preset
}

func Chrome101() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36`),
	)

	return preset
}

func Chrome104() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `"Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`),
	)

	return preset
}

func Chrome107() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `"Google Chrome";v="107", "Chromium";v="107", "Not=A?Brand";v="24"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36`),
	)

	preset.SetFlags(
		types.Flag("http2-no-server-push", true),
	)
	return preset
}

func Chrome110() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`),
		types.Header("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`),
	)

	preset.SetFlags(
		types.Flag("http2-no-server-push", true),
		types.Flag("tls-permute-extensions", true),
	)
	return preset
}

func Chrome116() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36`),
		types.Header("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`),
	)

	preset.SetFlags(
		types.Flag("http2-no-server-push", true),
		types.Flag("tls-permute-extensions", true),
	)
	return preset
}

func Edge99() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Microsoft Edge";v="99"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 Edg/99.0.1150.30`),
		types.Header("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`),
	)
	return preset
}

func Edge101() Preset {
	preset := Chrome99()
	preset.SetHeaders(
		types.Header("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Microsoft Edge";v="101"`),
		types.Header("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36 Edg/101.0.1210.47`),
		types.Header("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`),
	)
	return preset
}

func Safari153() Preset {
	h := types.NewHeaders(
		types.Header("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15`),
		types.Header("Accept", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15`),
		types.Header("Accept-Language", `en-us`),
		types.Header("Accept-Encoding", `gzip, deflate, br`),
	)
	f := types.NewFlags(
		types.Flag("ciphers", "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384:TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384:TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:TLS_RSA_WITH_AES_256_GCM_SHA384:TLS_RSA_WITH_AES_128_GCM_SHA256:TLS_RSA_WITH_AES_256_CBC_SHA256:TLS_RSA_WITH_AES_128_CBC_SHA256:TLS_RSA_WITH_AES_256_CBC_SHA:TLS_RSA_WITH_AES_128_CBC_SHA:TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA:TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:TLS_RSA_WITH_3DES_EDE_CBC_SHA"),
		types.Flag("curves", "X25519:P-256:P-384:P-521"),
		types.Flag("signature-hashes", "ecdsa_secp256r1_sha256,rsa_pss_rsae_sha256,rsa_pkcs1_sha256,ecdsa_secp384r1_sha384,ecdsa_sha1,rsa_pss_rsae_sha384,rsa_pss_rsae_sha384,rsa_pkcs1_sha384,rsa_pss_rsae_sha512,rsa_pkcs1_sha512,rsa_pkcs1_sha1"),
		types.Flag("http2", true),
		types.Flag("compressed", true),
		types.Flag("tlsv1.0", true),
		types.Flag("no-tls-session-ticket", true),
		types.Flag("http2-pseudo-headers-order", "mspa"),
	)

	return Preset{h, f}
}

func Safari155() Preset {
	preset := Safari153()
	preset.SetHeaders(
		types.Header("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Safari/605.1.15`),
		types.Header("Accept", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15`),
		types.Header("Accept-Language", `en-GB,en-US;q=0.9,en;q=0.8`),
		types.Header("Accept-Encoding", `gzip, deflate, br`),
	)
	preset.SetFlags(
		types.Flag("cert-compression", "zlib"),
	)

	return preset
}
