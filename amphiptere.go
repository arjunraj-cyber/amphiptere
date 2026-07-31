package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	stateFile = "/var/run/amphiptere.lock"
	authPort  = ":9999"
)

func printBanner() {
	fmt.Println("============================================================")
	fmt.Println("                 AMPHIPTERE SECURITY SENTINEL               ")
	fmt.Println("             'Born to defend the sacred gold.'              ")
	fmt.Println("============================================================")
	fmt.Println("          /\\")
	fmt.Println("         /  \\")
	fmt.Println("        /  |  \\")
	fmt.Println("       /   |   \\")
	fmt.Println("      |    A    |")
	fmt.Println("      |    |    |")
	fmt.Println("       \\  / \\  /")
	fmt.Println("        \\/   \\/")
	fmt.Println("------------------------------------------------------------")
	fmt.Println(" Creator : Arjun Raj")
	fmt.Println(" Contact : arjunraj.cyber@gmail.com")
	fmt.Println("============================================================")
}

func main() {
	enableFlag := flag.Bool("enable", false, "Enable Amphiptere OS perimeter defense")
	disableFlag := flag.Bool("disable", false, "Disable Amphiptere OS perimeter defense")
	statusFlag := flag.Bool("status", false, "Check Amphiptere defense status")
	flag.Parse()

	printBanner()

	if *enableFlag {
		enableDefense()
	} else if *disableFlag {
		disableDefense()
	} else if *statusFlag {
		checkStatus()
	} else {
		fmt.Println("\n[!] Usage: sudo amphiptere --enable | --disable | --status")
	}
}

func hashCodestr(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func enableDefense() {
	if os.Geteuid() != 0 {
		fmt.Println("\n[-] Error: You must run Amphiptere with root/sudo privileges to manage system firewalls.")
		return
	}

	fmt.Print("\n[?] Enter temporary Amphiptere Code: ")
	byteCode, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\n[-] Error reading code.")
		return
	}
	fmt.Println()

	codeStr := strings.TrimSpace(string(byteCode))
	if codeStr == "" {
		fmt.Println("[-] Code cannot be empty.")
		return
	}

	hashed := hashCodestr(codeStr)

	err = os.WriteFile(stateFile, []byte(hashed), 0600)
	if err != nil {
		fmt.Printf("[-] Failed to initialize defense state: %v\n", err)
		return
	}

	cmd := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", "9999", "-j", "ACCEPT")
	if err := cmd.Run(); err != nil {
		fmt.Printf("[-] Warning: Failed to hook into iptables: %v\n", err)
	}

	fmt.Println("\n[+] Amphiptere Defense: ACTIVE | Treasury is sealed.")
	fmt.Println("[+] External connections require the temporary code on port 9999.")

	go startAuthListener(hashed)
}

func disableDefense() {
	if os.Geteuid() != 0 {
		fmt.Println("\n[-] Error: You must run Amphiptere with root/sudo privileges.")
		return
	}

	cmd := exec.Command("iptables", "-D", "INPUT", "-p", "tcp", "--dport", "9999", "-j", "ACCEPT")
	_ = cmd.Run()

	if _, err := os.Stat(stateFile); err == nil {
		os.Remove(stateFile)
	}

	fmt.Println("\n[+] Amphiptere Defense: DISABLED | State flushed. Code cleared.")
}

func checkStatus() {
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		fmt.Println("\n[*] Amphiptere Status: INACTIVE (Perimeter open)")
	} else {
		fmt.Println("\n[!] Amphiptere Status: ACTIVE (Treasury protected by sentinel)")
	}
}

func startAuthListener(correctHash string) {
	listener, err := net.Listen("tcp", authPort)
	if err != nil {
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if _, err := os.Stat(stateFile); os.IsNotExist(err) {
				break
			}
			continue
		}

		go handleRemoteHandshake(conn, correctHash)
	}
}

func handleRemoteHandshake(conn net.Conn, correctHash string) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	conn.Write([]byte("AMPHIPTERE GATEWAY SECURITY > Enter Code: "))
	reader := bufio.NewReader(conn)
	inputCode, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	inputHash := hashCodestr(strings.TrimSpace(inputCode))

	if inputHash == correctHash {
		conn.Write([]byte("[ACCESS GRANTED] Welcome to Amphiptere OS treasury.\n"))
	} else {
		conn.Write([]byte("[ACCESS DENIED] Firewall rejection. Invalid code.\n"))
	}
}
