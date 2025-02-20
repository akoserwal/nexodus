//go:build darwin

package nexodus

import (
	"fmt"
	"github.com/nexodus-io/nexodus/internal/util"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"os/exec"
)

func (ax *Nexodus) setupInterfaceOS() error {

	logger := ax.logger
	localAddress := ax.TunnelIP
	localAddressIPv6 := fmt.Sprintf("%s/%s", ax.TunnelIpV6, wgOrgIPv6PrefixLen)
	dev := ax.tunnelIface

	if ifaceExists(logger, dev) {
		deleteDarwinIface(logger, dev)
	}

	// prefer nexd-wireguard-go over wireguard-go since it supports port reuse.
	wgBinary := wgGoBinary
	if path, err := exec.LookPath(nexdWgGoBinary); err == nil {
		wgBinary = path
	}

	_, err := RunCommand(wgBinary, dev)
	if err != nil {
		logger.Errorf("failed to create the %s interface: %v\n", dev, err)
		return fmt.Errorf("%w", interfaceErr)
	}

	_, err = RunCommand("ifconfig", dev, "inet", localAddress, localAddress, "alias")
	if err != nil {
		logger.Errorf("failed to assign an IPv4 address to the local osx interface: %v\n", err)
		return fmt.Errorf("%w", interfaceErr)
	}

	if ax.ipv6Supported {
		_, err = RunCommand("ifconfig", dev, "inet6", localAddressIPv6, "alias")
		if err != nil {
			logger.Errorf("failed to assign an IPv6 address to the local osx interface: %v\n", err)
			return fmt.Errorf("%w", interfaceErr)
		}
	}

	_, err = RunCommand("ifconfig", dev, "up")
	if err != nil {
		logger.Errorf("failed to bring up the %s interface: %v\n", dev, err)
		return fmt.Errorf("%w", interfaceErr)
	}

	privateKey, err := wgtypes.ParseKey(ax.wireguardPvtKey)
	if err != nil {
		logger.Errorf("invalid wiregaurd private key: %v\n", err)
		return fmt.Errorf("%w", interfaceErr)
	}

	c, err := wgctrl.New()
	if err != nil {
		logger.Errorf("could not connect to wireguard: %v\n", err)
		return fmt.Errorf("%w", interfaceErr)
	}
	defer util.IgnoreError(c.Close)

	err = c.ConfigureDevice(dev, wgtypes.Config{
		PrivateKey:   &privateKey,
		ListenPort:   &ax.listenPort,
		ReplacePeers: true,
		Peers:        nil,
	})
	if err != nil {
		logger.Errorf("failed to start the wireguard listener: %v\n", err)
		return fmt.Errorf("%w", interfaceErr)
	}

	return nil
}

func (ax *Nexodus) removeExistingInterface() {
	if ifaceExists(ax.logger, ax.tunnelIface) {
		deleteDarwinIface(ax.logger, ax.tunnelIface)
	}
}

// deleteDarwinIface delete the darwin userspace wireguard interface
func deleteDarwinIface(logger *zap.SugaredLogger, dev string) {
	tunSock := fmt.Sprintf("/var/run/wireguard/%s.sock", dev)
	_, err := RunCommand("rm", "-f", tunSock)
	if err != nil {
		logger.Debugf("failed to delete darwin interface: %v", err)
	}
	// /var/run/wireguard/wg0.name doesn't currently exist since utun8 isn't mapped to wg0 (fails silently)
	wgName := fmt.Sprintf("/var/run/wireguard/%s.name", dev)
	_, err = RunCommand("rm", "-f", wgName)
	if err != nil {
		logger.Debugf("failed to delete darwin interface: %v", err)
	}
}

func (ax *Nexodus) findLocalIP() (string, error) {
	return discoverGenericIPv4(ax.logger, ax.apiURL.Host, "443")
}
