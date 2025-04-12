package pkg

type JSONRequestGenerateTtd struct {
	Request struct {
		Initial string `json:"initial"`
		Code    string `json:"code"`
	} `json:"request"`
}

type JSONRequestSign struct {
	Request struct {
		Setup []struct {
			Page   int    `json:"page"`
			Image  string `json:"image"`
			X      int    `json:"x"`
			Y      int    `json:"y"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"setup"`
		DocsCode string `json:"doc_id"`
	} `json:"request"`
}

type JSONRequestCompleteSign struct {
	Request struct {
		DocsId string `json:"doc_id"`
	} `json:"request"`
}

type JSONRequestCreateClient struct {
	Request struct {
		Name         string  `bson:"name"`
		Email        string  `bson:"email"`
		Type         string  `bson:"type"`
		IsActive     bool    `bson:"is_active"`
		Point        int     `bson:"point"`
		UsernameAuth string  `bson:"username_auth"`
		PasswordAuth string  `bson:"password_auth"`
		ClientID     string  `bson:"client_id"`
		Rph          int     `bson:"rph"`
		BillAmount   float64 `bson:"bill_amount"`
	}
}

type JSONRequestPayment struct {
	Request struct {
		UsernameAuth string  `bson:"username_auth" json:"username_auth"`
		Amount       float64 `bson:"amount"`
	}
}

type JSONRequestReadPayment struct {
	Request struct {
		Va string `bson:"va" json:"va"`
	}
}

type ImageRequest struct {
	ImageBase64 string `json:"image"`
}

type OCRRequest struct {
	ImageBase64 string `json:"image_base64"`
}

type SplitBillGroupRequest struct {
	Request struct {
		NamaGroup string   `bson:"namagroup" json:"namagroup"`
		User      []string `bson:"user" json:"user"`
	} `json:"request"`
}
