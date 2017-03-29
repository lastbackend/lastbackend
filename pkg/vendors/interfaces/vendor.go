package interfaces

type Vendor struct {
	Vendor string
	Host   string
}

type IVendor interface {
	GetVendorInfo() *Vendor
}
