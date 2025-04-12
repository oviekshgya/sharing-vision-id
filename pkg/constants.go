package pkg

import (
	"golang.org/x/oauth2"
	"sync"
)

const APIKEY = "xzDcUsxhstdalZtbdMz0"
const USERNAME = "Shagya"
const PASSWORD = "ShagyaTech"

var StatusPengajuan = []string{"On Progress", "Survey", "Reject", "Approve", "Funding"}

var GoogleOAuthConfig = &oauth2.Config{}

var RateLimitMap = sync.Map{}

const PERHIT = 100
