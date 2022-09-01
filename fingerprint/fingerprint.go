////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package fingerprint

import (
	"encoding/base64"
	"gitlab.com/xx_network/crypto/hasher"
)

// Creates a fingerprint of a public key
// NOTE: This is just a hash for now
func Fingerprint(publicKey []byte) string {
	h := hasher.BLAKE2.New()
	h.Write(publicKey)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
