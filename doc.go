// Package at implements a simple interface to communicate with mobile
// IoT modules.  For now we only have support for one family of
// modules (the uBlox SARA N2 and possibly N3 modules), but we might
// add support for more modules in the future.
//
// Example of how to connect to a device:
//
//    import "github.com/lab5e/at/n211"
//
//    ...
//
//    device := n211.New(os.Args[1], baudRate)
//    if err := device.Start(); err != nil {
//        log.Fatalf("Error opening device: %v", err)
//    }
//
// Please refer to the Device interface to see what methods are available.
//
package at
