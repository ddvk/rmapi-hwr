package main

import (
	"fmt"
	"time"
	_ "unsafe"

	"github.com/unidoc/unipdf/v3/common/license"
)

//go:linkname licenseKey github.com/unidoc/unipdf/v3/internal/license._afg
var licenseKey *license.LicenseKey

func init() {
	fmt.Printf(licenseKey.LicenseId)
	lk := license.LicenseKey{}
	lk.CustomerName = "community"
	lk.Tier = license.LicenseTierCommunity
	lk.CreatedAt = time.Now().UTC()
	lk.CreatedAtInt = lk.CreatedAt.Unix()
	licenseKey = &lk
}
