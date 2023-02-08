package configs

import "testing"

func NewMockConfig(t *testing.T) *Config {
	t.Helper()
	return &Config{
		BindAddr: ":8005",
		LogLevel: "debug",
		Database: database{
			URI:                "",
			MaxOpenConnections: 50,
			MaxIDLEConnections: 10,
		},
		Auth: auth{
			PublicKey: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP
nB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww
sL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+
m1EXClDQU1sAa4LMeQIDAQAB
-----END PUBLIC KEY-----
`,
			PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz
V2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO
J1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB
AoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA
1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB
2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ
UXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5
Ne9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT
UrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO
MKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX
E8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM
kbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca
gjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw
-----END RSA PRIVATE KEY-----`,
			RefreshTTL: 172800,
			AccessTTL:  86400,
		},
	}
}
