package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hanwen/usb"
)

const (
	VendorId     = 0x1038
	ProductId    = 0x1211
	InterfaceNum = 3
	Timeout      = 1000
)

func frame(b1 byte, b2 byte) []byte {
	result := [37]byte{
		0x04, 0x40, 0x01, 0x11, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	result[4] = b1
	result[5] = b2

	return result[:]
}

func send(h *usb.DeviceHandle, b1 byte, b2 byte) error {
	return h.ControlTransfer(
		usb.ENDPOINT_OUT|usb.REQUEST_TYPE_CLASS|usb.RECIPIENT_INTERFACE,
		0x09,   // SET_REPORT
		0x0204, // Output, ID: 4
		InterfaceNum,
		frame(b1, b2),
		Timeout,
	)
}

func process(h *usb.DeviceHandle) error {
	err := h.ClaimInterface(InterfaceNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ClaimInterface failed: %s\n", err)
		return err
	}
	defer h.ReleaseInterface(InterfaceNum)

	if *intensity >= 0 {
		err = send(h, 0x89, byte(*intensity))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ControlTransfer failed: %s\n", err)
			return err
		}
	}

	if *pulse != "" {
		switch *pulse {
		case "steady":
			err = send(h, 0x87, 0x02)
		case "slow":
			err = send(h, 0x87, 0x22)
		case "medium":
			err = send(h, 0x87, 0x26)
		case "fast":
			err = send(h, 0x87, 0x2A)
		case "trigger":
			err = send(h, 0x87, 0x12)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ControlTransfer failed: %s\n", err)
			return err
		}
	}

	return nil
}

var intensity *int = flag.Int("i", -1, "LEDs intensity (0-255)")
var pulse *string = flag.String("p", "", "LEDs pulsating mode (steady, slow, medium, fast, trigger)")

func main() {
	flag.Parse()

	if *intensity > 255 {
		fmt.Fprintf(os.Stderr, "Intensity should be in range 0-255\n")
		os.Exit(2)
	}

	if *pulse != "" && *pulse != "steady" && *pulse != "slow" && *pulse != "medium" && *pulse != "fast" && *pulse != "trigger" {
		fmt.Fprintf(os.Stderr, "Pulse mode should be steady, slow, medium, fast or trigger")
		os.Exit(2)
	}

	ctx := usb.NewContext()
	defer ctx.Exit()

	devices, err := ctx.GetDeviceList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetDeviceList failed: %s\n", err)
		os.Exit(1)
	}

	var device *usb.Device
	func() {
		defer devices.Done()
		for _, d := range devices {
			dd, err := d.GetDeviceDescriptor()
			if err != nil {
				fmt.Fprintf(os.Stderr, "GetDeviceDescriptor failed: %s\n", err)
				continue
			}

			if dd.IdVendor == VendorId && dd.IdProduct == ProductId {
				fmt.Printf("Found SteelSeries v2 headset\n")
				device = d.Ref()
				break
			}
		}
	}()

	if device == nil {
		fmt.Fprintf(os.Stderr, "Can't find supported device\n")
		os.Exit(1)
	}

	defer device.Unref()

	h, err := device.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open device failed: %s\n", err)
		os.Exit(1)
	}
	defer h.Close()

	active, err := h.KernelDriverActive(InterfaceNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "KernelDriverActive filed: %s\n", err)
		os.Exit(1)
	}

	if active {
		func() {
			err := h.DetachKernelDriver(InterfaceNum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "DetachKernelDriver failed: %s\n", err)
				os.Exit(1)
			}
			defer h.AttachKernelDriver(InterfaceNum)

			process(h)
		}()
	} else {
		process(h)
	}
}
