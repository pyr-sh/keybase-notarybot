/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

// Package license helps manage commercial licenses and check if they are valid for the version of unidoc used.
package license

import "time"

var createdAt = time.Now().UTC()

// Defaults to unlicensed.
var licenseKey = &LicenseKey{
	CustomerName: "AGPLv3",
	Tier:         LicenseTierCommunity,
	CreatedAt:    createdAt,
	CreatedAtInt: createdAt.Unix(),
}

func GetLicenseKey() *LicenseKey {
	if licenseKey == nil {
		return nil
	}

	// Copy.
	lk2 := *licenseKey
	return &lk2
}
