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
			PublicKey:  "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChTrU+r2uTQPQOxBCwKVAM0AJP\nnB4MEh+MggX5lkrGOPtzBglzV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcww\nsL+q7oKBKbiJYtrYGr7uoJrOJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+\nm1EXClDQU1sAa4LMeQIDAQAB\n-----END PUBLIC KEY-----",
			PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQChTrU+r2uTQPQOxBCwKVAM0AJPnB4MEh+MggX5lkrGOPtzBglz\nV2DF+ydVknJoqmUFbnczXJsAFaaaXCYm7N/kOcwwsL+q7oKBKbiJYtrYGr7uoJrO\nJ1SIWq/RnvkWpGGqth6SvryEB742l0WAiG3nxWY+m1EXClDQU1sAa4LMeQIDAQAB\nAoGAGZAxpPeD4tg+VUC5LFG/v+gPFbK2CE+u9EN+0ukAfJ13K+lfAgps6bM9rpAA\n1Zl7XPr+pQMeBUtpFblYyn5rlK0oultlJI//H0I3+6newKp7LewPIrV08lGEn1hB\n2XtSAvZVShsCmtyw8UvXwHk01UJA0pEyGdkWiHE3jEuCUSkCQQDSsulNRw/G+8xZ\nUXTCgb9ep9EojDIQYqAeomX9/CMgS6QAWERPt9Q37ZHkki0i1iicOdZc94C7PxA5\nNe9DhGofAkEAw/09n+v2YBPpYY1Wik1NKA4I1Q3/zZlsop3W+fCiJZiO3Dhef0TT\nUrQmYSMftbe6peSo3yQGVPnBGB+0phSmZwJAWJaW10IQlSZblhZUlE9/SeofXAAO\nMKt3DUpUvcRcdIC5NNfn6Oiu1tERbVw0lBgdPQpoYfBCdPgf9x4BOo8bGwJBAKiX\nE8aYXNQi7LQMt6+6dS+KexCCvVPnsWplKkLQOzrp86H+H1ONKddPvl/6rdFMHZOM\nkbN5MrUwLmkJBQWEZ+sCQQClKUu0DYu+XgbDPrYgxJNAgWTtVTZ2wLCp46X4iHca\ngjOIscTm3jUVsz8bCkXrVlFsWRVCnvQwKx788Awq6mdw\n-----END RSA PRIVATE KEY-----",
			RefreshTTL: 172800,
			AccessTTL:  86400,
		},
	}
}
