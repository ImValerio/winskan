package system

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// USBDevice represents a USB device history entry found in the SYSTEM registry.
type USBDevice struct {
	Vendor       string
	Product      string
	Revision     string
	SerialNumber string
	FriendlyName string
	VolumeName   string
	HardwareID   []string
}

// GetUSBHistory retrieves the history of USB storage devices connected to the system.
func GetUSBHistory() ([]USBDevice, error) {
	const basePath = `SYSTEM\CurrentControlSet\Enum\USBSTOR`

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, basePath, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		// If the key does not exist, there's no USB history to return
		if err == registry.ErrNotExist {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open USBSTOR key: %w", err)
	}
	defer key.Close()

	deviceClasses, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read device classes: %w", err)
	}

	volumeNames := getVolumeNames()

	var devices []USBDevice

	for _, class := range deviceClasses {
		vendor, product, revision := parseClassString(class)

		classKeyPath := fmt.Sprintf(`%s\%s`, basePath, class)
		classKey, err := registry.OpenKey(registry.LOCAL_MACHINE, classKeyPath, registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			continue
		}

		instances, err := classKey.ReadSubKeyNames(-1)
		classKey.Close()
		if err != nil {
			continue
		}

		for _, instance := range instances {
			instanceKeyPath := fmt.Sprintf(`%s\%s`, classKeyPath, instance)
			instanceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, instanceKeyPath, registry.QUERY_VALUE)
			if err != nil {
				continue
			}

			friendlyName, _, _ := instanceKey.GetStringValue("FriendlyName")
			hardwareID, _, _ := instanceKey.GetStringsValue("HardwareID")

			instanceKey.Close()

			volName := volumeNames[strings.ToUpper(instance)]

			devices = append(devices, USBDevice{
				Vendor:       vendor,
				Product:      product,
				Revision:     revision,
				SerialNumber: instance,
				FriendlyName: friendlyName,
				VolumeName:   volName,
				HardwareID:   hardwareID,
			})
		}
	}

	return devices, nil
}

func getVolumeNames() map[string]string {
	volumes := make(map[string]string)
	const wpdPath = `SOFTWARE\Microsoft\Windows Portable Devices\Devices`

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, wpdPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return volumes
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return volumes
	}

	for _, subkey := range subkeys {
		deviceKeyPath := fmt.Sprintf(`%s\%s`, wpdPath, subkey)
		deviceKey, err := registry.OpenKey(registry.LOCAL_MACHINE, deviceKeyPath, registry.QUERY_VALUE)
		if err != nil {
			continue
		}

		friendlyName, _, err := deviceKey.GetStringValue("FriendlyName")
		deviceKey.Close()

		if err == nil && friendlyName != "" {
			parts := strings.Split(subkey, "#")
			if len(parts) >= 5 && strings.Contains(parts[2], "USBSTOR") {
				serial := parts[4]
				volumes[strings.ToUpper(serial)] = friendlyName
			}
		}
	}

	return volumes
}

func parseClassString(class string) (vendor, product, revision string) {
	parts := strings.Split(class, "&")
	for _, part := range parts {
		if strings.HasPrefix(part, "Ven_") {
			vendor = strings.TrimPrefix(part, "Ven_")
			vendor = strings.ReplaceAll(vendor, "_", " ")
		} else if strings.HasPrefix(part, "Prod_") {
			product = strings.TrimPrefix(part, "Prod_")
			product = strings.ReplaceAll(product, "_", " ")
		} else if strings.HasPrefix(part, "Rev_") {
			revision = strings.TrimPrefix(part, "Rev_")
			revision = strings.ReplaceAll(revision, "_", " ")
		}
	}
	return strings.TrimSpace(vendor), strings.TrimSpace(product), strings.TrimSpace(revision)
}
