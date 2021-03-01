package fabsdk

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

const (
	FuncNewProductInfo       = "NewProductInfo"
	FuncNewOrder             = "NewOrder"
	FuncUpdateLogisticRecord = "UpdateLogisticRecord"
	FuncAcceptOrder          = "AcceptOrder"
	FuncRejectOrderRequest   = "RejectOrderRequest"

	FuncReadProductInfo              = "ReadProductInfo"
	FuncReadOrderInfo                = "ReadOrderInfo"
	FuncReadLogisticsRecord          = "ReadLogisticsRecord"
	FuncReadLogisticsPriRecord       = "ReadLogisticsPriRecord"
	FuncReadProductInfoByProductName = "ReadProductInfoByProductName"
	FuncGetOrdersUnaccepted          = "GetOrdersUnaccepted"

	TransientKeyOrderInput            = "order_input"
	TransientKeyLogisticOperatorInput = "operator_info"
	TransientKeyAcceptOrderInput      = "accept_order_input"

	MSPIDManufacturer = "ManufacturerMSP"
	MSPIDRetailer1    = "Retailer1MSP"
	MSPIDRetailer2    = "Retailer2MSP"
	MSPIDLogistic     = "LogisticsMSP"
	MSPIDRegulator    = "RegulatorMSP"

	CollectionTransaction1 = "transactionCollection1"
	CollectionTransaction2 = "transactionCollection2"
	CollectionLogistics    = "logisticsCollection"

	TlsCAFile = `/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`

	NUMOrderer      = 0
	NUMManufacturer = 1
	NUMLogistics    = 2
	NUMRetailer1    = 3
	NUMRetailer2    = 4
	NUMRegulator    = 5

	OrderStatusAccepted   = "Accepted"
	OrderStatusRejected   = "Rejected"
	OrderStatusUnaccepted = "Unaccepted"
)

type FabricSDKCtx struct {
	Channel       string
	CCName        string
	Env           []string
	TlsCACertFile string
}

func NewFabSDKCtx(org int) (*FabricSDKCtx, error) {
	sdk := &FabricSDKCtx{
		Channel:       "logisticschannel",
		CCName:        "logisticscc",
		Env:           os.Environ(),
		TlsCACertFile: TlsCAFile,
	}

	switch org {
	case NUMOrderer:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=OrdererMSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp`,
		)
	case NUMManufacturer:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=ManufacturerMSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/manufacturer.example.com/peers/peer0.manufacturer.example.com/tls/ca.crt`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/manufacturer.example.com/users/Admin@manufacturer.example.com/msp`,
			`CORE_PEER_ADDRESS=localhost:7051`,
		)
	case NUMLogistics:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=LogisticsMSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/ca.crt`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/logistics.example.com/users/Admin@logistics.example.com/msp`,
			`CORE_PEER_ADDRESS=localhost:9051`,
		)
	case NUMRetailer1:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=Retailer1MSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/retailer1.example.com/peers/peer0.retailer1.example.com/tls/ca.crt`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/retailer1.example.com/users/Admin@retailer1.example.com/msp`,
			`CORE_PEER_ADDRESS=localhost:11051`,
		)
	case NUMRetailer2:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=Retailer2MSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/retailer2.example.com/peers/peer0.retailer2.example.com/tls/ca.crt`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/retailer2.example.com/users/Admin@retailer2.example.com/msp`,
			`CORE_PEER_ADDRESS=localhost:13051`,
		)
	case NUMRegulator:
		sdk.Env = append(sdk.Env,
			"FABRIC_CFG_PATH=/home/joseph/fabric-tests/logistics-test-3/config/",
			`CORE_PEER_TLS_ENABLED=true`,
			`CORE_PEER_LOCALMSPID=RegulatorMSP`,
			`CORE_PEER_TLS_ROOTCERT_FILE=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/regulator.example.com/peers/peer0.regulator.example.com/tls/ca.crt`,
			`CORE_PEER_MSPCONFIGPATH=/home/joseph/fabric-tests/logistics-test-3/test-network/organizations/peerOrganizations/regulator.example.com/users/Admin@regulator.example.com/msp`,
			`CORE_PEER_ADDRESS=localhost:15051`,
		)
	default:
		return nil, errors.New("organization number must be specified")
	}

	return sdk, nil
}

// Invoke argument transientData: base64-encoded string
func (c *FabricSDKCtx) Invoke(function, args, transientKey, transientData string) ([]byte, error) {
	var cmd *exec.Cmd
	if transientKey == "" {
		cmd = exec.Command("bin/peer", "chaincode", "invoke",
			"-o", "localhost:7050",
			"--ordererTLSHostnameOverride", `orderer.example.com`,
			`--tls`,
			`--cafile`, c.TlsCACertFile,
			`-C`, c.Channel, "-n", c.CCName,
			"-c", `{"function":"`+function+`","Args":[`+args+`]}`,
		)
	} else {
		cmd = exec.Command("bin/peer", "chaincode", "invoke",
			"-o", "localhost:7050",
			"--ordererTLSHostnameOverride", `orderer.example.com`,
			`--tls`,
			`--cafile`, c.TlsCACertFile,
			`-C`, c.Channel, "-n", c.CCName,
			"-c", `{"function":"`+function+`","Args":[`+args+`]}`,
			`--transient`, `{"`+transientKey+`":"`+transientData+`"}`,
		)
	}

	cmd.Env = c.Env
	return cmd.CombinedOutput()
}

func (c *FabricSDKCtx) Query(function, args string, tls bool) ([]byte, error) {
	var cmd *exec.Cmd
	if tls {
		cmd = exec.Command("bin/peer", "chaincode", "query",
			"-o", "localhost:7050",
			"--ordererTLSHostnameOverride", `orderer.example.com`,
			`--tls`,
			`--cafile`, c.TlsCACertFile,
			`-C`, c.Channel, "-n", c.CCName,
			"-c", `{"function":"`+function+`","Args":[`+args+`]}`,
		)
	} else {
		cmd = exec.Command("bin/peer", "chaincode", "query",
			`-C`, c.Channel, "-n", c.CCName,
			"-c", `{"function":"`+function+`","Args":[`+args+`]}`,
		)
	}
	cmd.Env = c.Env
	return cmd.CombinedOutput()
}
