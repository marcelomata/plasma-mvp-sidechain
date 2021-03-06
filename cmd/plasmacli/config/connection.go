package config

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/flags"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

var plasmaContract *eth.Plasma

// GetContractConn returns the eth plasma contract connection
func GetContractConn() (*eth.Plasma, error) {
	if plasmaContract != nil {
		return plasmaContract, nil
	}

	conn, err := setupContractConn()
	if err != nil {
		return nil, fmt.Errorf("unable to enable contract connection: %s", err)
	}

	plasmaContract = conn
	return plasmaContract, nil
}

func setupContractConn() (*eth.Plasma, error) {
	conf, err := ParseConfigFromViper()
	if err != nil {
		return nil, err
	}

	dir := viper.GetString(flags.Home)
	if dir[len(dir)-1] != '/' {
		dir = dir + "/"
	}

	if conf.EthNodeURL == "" {
		return nil, fmt.Errorf("please specify a node url for eth connection in %sconfig.toml", dir)
	} else if conf.EthPlasmaContractAddr == "" || !ethcmn.IsHexAddress(conf.EthPlasmaContractAddr) {
		return nil, fmt.Errorf("please specify a valid contract address in %sconfig.toml", dir)
	}

	ethClient, err := eth.InitEthConn(conf.EthNodeURL)
	if err != nil {
		return nil, err
	}
	plasma, err := eth.InitPlasma(ethcmn.HexToAddress(conf.EthPlasmaContractAddr), ethClient, 0)
	if err != nil {
		return nil, err
	}

	return plasma, nil
}
