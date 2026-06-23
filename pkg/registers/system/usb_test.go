package system

import (
	"testing"
)

func TestGetUSBHistory(t *testing.T) {
	devices, err := GetUSBHistory()
	if err != nil {
		t.Fatalf("GetUSBHistory failed: %v", err)
	}

	for _, dev := range devices {
		t.Logf("Found USB Device: Vendor=%s, Product=%s, Revision=%s, Serial=%s, FriendlyName=%s",
			dev.Vendor, dev.Product, dev.Revision, dev.SerialNumber, dev.FriendlyName)
	}
}

func TestParseClassString(t *testing.T) {
	class := "Disk&Ven_Kingston&Prod_DataTraveler_2.0&Rev_1.00"
	vendor, product, revision := parseClassString(class)
	
	if vendor != "Kingston" {
		t.Errorf("Expected Vendor 'Kingston', got '%s'", vendor)
	}
	if product != "DataTraveler 2.0" {
		t.Errorf("Expected Product 'DataTraveler 2.0', got '%s'", product)
	}
	if revision != "1.00" {
		t.Errorf("Expected Revision '1.00', got '%s'", revision)
	}
}
