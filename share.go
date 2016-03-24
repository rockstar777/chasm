package main

// ShareID is a uniqiue id to represent uploaded shares
type ShareID string

// Share represents a secret share of a file
type Share struct {
	SID  ShareID
	Data []byte
}

// CreateShares creates two shares from secret
func CreateShares(secret []byte) (Share, Share) {
    return Share{}, Share{}
}


// CombineShares restores the secret by adding
func CombineShares(shares []Share) {
    
}
