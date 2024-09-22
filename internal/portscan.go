package internal

import (
	"fmt"
	"net"
)

// ScanPort scans a port on a given IP address
func ScanPort(ip string, port int) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Port", port, "is closed on", ip)
		return
	}
	conn.Close()
	fmt.Println("Port", port, "is open on", ip)
}
