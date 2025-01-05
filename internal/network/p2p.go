package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer" // ✅ Import correct host package
	"github.com/multiformats/go-multiaddr"
)

// P2PNode represents a peer node in the blockchain network
type P2PNode struct {
	Host host.Host
	Ctx  context.Context
}

// CreateNode initializes a new LibP2P node
func CreateNode(listenPort int) (*P2PNode, error) {
	ctx := context.Background()

	node, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🟢 P2P Node started on /ip4/0.0.0.0/tcp/%d\n", listenPort)
	fmt.Println("🔗 Node ID:", node.ID().String())

	// ✅ Register protocol handler for chat messages
	node.SetStreamHandler("/cortex-chat/1.0.0", handleStream)

	// ✅ Print listening addresses
	fmt.Println("📡 Listening on:")
	for _, addr := range node.Addrs() {
		fmt.Printf("➡️ %s/p2p/%s\n", addr.String(), node.ID().String())
	}

	return &P2PNode{Host: node, Ctx: ctx}, nil
}

// ValidatePeerAddr ensures the input is a valid multiaddr
func ValidatePeerAddr(peerAddr string) error {
	if !strings.HasPrefix(peerAddr, "/ip4/") || !strings.Contains(peerAddr, "/tcp/") || !strings.Contains(peerAddr, "/p2p/") {
		return errors.New("⚠️ Invalid peer address format. Use: /ip4/127.0.0.1/tcp/5001/p2p/QmPeerID")
	}
	return nil
}

// ConnectToPeer connects to another P2P node
func (n *P2PNode) ConnectToPeer(peerAddr string) error {
	peerAddr = strings.TrimSpace(peerAddr)

	// ✅ Check Address Format
	if err := ValidatePeerAddr(peerAddr); err != nil {
		return fmt.Errorf("❌ Invalid peer address: %v", err)
	}

	ma, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("❌ Invalid multiaddr: %v", err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		return fmt.Errorf("❌ Failed to parse peer info: %v", err)
	}

	fmt.Println("🔄 Attempting to connect to:", peerInfo.ID)

	// ✅ Check if already connected
	if n.Host.Network().Connectedness(peerInfo.ID) == network.Connected {
		fmt.Println("✅ Already connected to:", peerInfo.ID)
		return nil
	}

	// ✅ Attempt connection
	err = n.Host.Connect(n.Ctx, *peerInfo)
	if err != nil {
		return fmt.Errorf("❌ Connection failed: %v", err)
	}

	fmt.Println("✅ Successfully connected to:", peerInfo.ID)
	return nil
}

// BroadcastMessage sends a message to all connected peers
func (n *P2PNode) BroadcastMessage(msg string) {
	for _, peerID := range n.Host.Network().Peers() {
		stream, err := n.Host.NewStream(n.Ctx, peerID, "/cortex-chat/1.0.0")
		if err != nil {
			log.Println("❌ Error opening stream:", err)
			continue
		}

		writer := bufio.NewWriter(stream)
		writer.WriteString(msg + "\n")
		writer.Flush()
		stream.Close()
	}
	fmt.Println("📡 Message broadcasted:", msg)
}

// handleStream processes incoming messages from peers
func handleStream(stream network.Stream) {
	fmt.Println("📥 Incoming connection from:", stream.Conn().RemotePeer())

	reader := bufio.NewReader(stream)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("❌ Error reading message:", err)
		return
	}

	fmt.Println("📩 Received message:", strings.TrimSpace(msg))
}
