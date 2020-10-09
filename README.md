# AT command library

**This is a work in progress.  Use at your own risk**

[![GoDoc](https://godoc.org/github.com/lab5e/at?status.svg)](https://godoc.org/github.com/lab5e/at)

The AT command library is a Go library for communicating with mobile
network IoT modules via their AT command set.  Currently only the
uBlox Sara N210/N211 is supported, but this library might work with
the uBlox N310 too.

As mentioned above this is a work in progress so expect things to
change around and maintainance of this library to be somewhat
episodic. 

If you want to contribute or you have suggestions, please do not
hesitate to contact @borud.


## Sample code

    package main
    
    import (
    	"log"
    	"os"
    
    	"github.com/lab5e/at/n211"
    )
    
    const baudRate = 9600
    
    func main() {
    	if len(os.Args) < 2 {
    		log.Fatalf("Usage %s <serial device>", os.Args[0])
    	}
    
    	device := n211.New(os.Args[1], baudRate)
    	if err := device.Start(); err != nil {
    		log.Fatalf("Error opening device: %v", err)
    	}
    	defer device.Close()
    
    	// Turn on debugging so you can see the interaction with the device
    	device.SetDebug(true)
    
    	// Just send a blank AT command to verify the device is there
    	if err := device.AT(); err != nil {
    		log.Fatalf("Error speaking to device on '%s': %v", os.Args[1], err)
    	}
    	log.Printf("Device seems to be responsive")
    }
    
