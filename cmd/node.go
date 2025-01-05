package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cortex/internal/blockchain"
	"cortex/internal/network"
)

func main() {
	fmt.Println("🌐 Cortex Blockchain Node 🌐")

	// Choose Port
	fmt.Print("Enter port to run node (e.g. 5001): ")
	reader := bufio.NewReader(os.Stdin)
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)
	port, _ := strconv.Atoi(portInput)

	// Start P2P Node
	node, err := network.CreateNode(port)
	if err != nil {
		fmt.Println("❌ Failed to start P2P node:", err)
		return
	}

	// Blockchain Initialization
	bc := blockchain.NewBlockchain()

	for {
		fmt.Println("\n1️⃣ Connect to Peer")
		fmt.Println("2️⃣ Broadcast Message")
		fmt.Println("3️⃣ Print Blockchain")
		fmt.Println("4️⃣ Exit")

		fmt.Print("Select an option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option) // ✅ Fix: Remove "\n"

		switch option {
		case "1":
			fmt.Print("Enter Peer Multiaddress: ")
			peerAddr, _ := reader.ReadString('\n')
			peerAddr = strings.TrimSpace(peerAddr) // ✅ Trim input
			err := node.ConnectToPeer(peerAddr)
			if err != nil {
				fmt.Println(err) // ✅ Print user-friendly error
			}

		case "2":
			fmt.Print("Enter message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			node.BroadcastMessage(msg)

		case "3":
			for _, block := range bc.GetBlocks() {
				fmt.Printf("🔗 Block %d | Hash: %s | Data: %s\n", block.Index, block.Hash, block.Data)
			}

		case "4":
			fmt.Println("Shutting down...")
			os.Exit(0)

		default:
			fmt.Println("⚠️ Invalid option, please try again.") // ✅ Fixed
		}
	}
}
