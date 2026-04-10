package contract

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

/* contract.go — MeshPulse Canopy Plugin
   Extends the base 'send' transaction with two custom types:
     - submit_measurement  →  records a speed-test result, mints 10 $MESHP
     - claim_reward        →  transfers earned $MESHP to a wallet address
*/

var ContractConfig = &PluginConfig{
	Name:    "meshpulse",
	Id:      1,
	Version: 1,
	SupportedTransactions: []string{
		"send",
		"submit_measurement",
		"claim_reward",
	},
	TransactionTypeUrls: []string{
		"type.googleapis.com/types.MessageSend",
		"type.googleapis.com/types.MessageSubmitMeasurement",
		"type.googleapis.com/types.MessageClaimReward",
	},
	EventTypeUrls: nil,
}

func init() {
	file_account_proto_init()
	file_event_proto_init()
	file_plugin_proto_init()
	file_tx_proto_init()
	file_meshpulse_proto_init()

	var fds [][]byte
	for _, file := range []protoreflect.FileDescriptor{
		anypb.File_google_protobuf_any_proto,
		File_account_proto,
		File_event_proto,
		File_plugin_proto,
		File_tx_proto,
		File_meshpulse_proto,
	} {
		fd, _ := proto.Marshal(protodesc.ToFileDescriptorProto(file))
		fds = append(fds, fd)
	}
	ContractConfig.FileDescriptorProtos = fds
}

type Contract struct {
	Config    Config
	FSMConfig *PluginFSMConfig
	plugin    *Plugin
	fsmId     uint64
}

func (c *Contract) Genesis(_ *PluginGenesisRequest) *PluginGenesisResponse {
	return &PluginGenesisResponse{}
}

func (c *Contract) BeginBlock(_ *PluginBeginRequest) *PluginBeginResponse {
	return &PluginBeginResponse{}
}

func (c *Contract) CheckTx(request *PluginCheckRequest) *PluginCheckResponse {
	resp, err := c.plugin.StateRead(c, &PluginStateReadRequest{
		Keys: []*PluginKeyRead{
			{QueryId: rand.Uint64(), Key: KeyForFeeParams()},
		},
	})
	if err == nil {
		err = resp.Error
	}
	if err != nil {
		return &PluginCheckResponse{Error: err}
	}
	minFees := new(FeeParams)
	if err = Unmarshal(resp.Results[0].Entries[0].Value, minFees); err != nil {
		return &PluginCheckResponse{Error: err}
	}
	if request.Tx.Fee < minFees.SendFee {
		return &PluginCheckResponse{Error: ErrTxFeeBelowStateLimit()}
	}
	msg, err := FromAny(request.Tx.Msg)
	if err != nil {
		return &PluginCheckResponse{Error: err}
	}
	switch x := msg.(type) {
	case *MessageSend:
		return c.CheckMessageSend(x)
	case *MessageSubmitMeasurement:
		return c.CheckMessageSubmitMeasurement(x)
	case *MessageClaimReward:
		return c.CheckMessageClaimReward(x)
	default:
		return &PluginCheckResponse{Error: ErrInvalidMessageCast()}
	}
}

func (c *Contract) DeliverTx(request *PluginDeliverRequest) *PluginDeliverResponse {
	msg, err := FromAny(request.Tx.Msg)
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	switch x := msg.(type) {
	case *MessageSend:
		return c.DeliverMessageSend(x, request.Tx.Fee)
	case *MessageSubmitMeasurement:
		return c.DeliverMessageSubmitMeasurement(x, request.Tx.Fee)
	case *MessageClaimReward:
		return c.DeliverMessageClaimReward(x, request.Tx.Fee)
	default:
		return &PluginDeliverResponse{Error: ErrInvalidMessageCast()}
	}
}

func (c *Contract) EndBlock(_ *PluginEndRequest) *PluginEndResponse {
	return &PluginEndResponse{}
}

func (c *Contract) CheckMessageSend(msg *MessageSend) *PluginCheckResponse {
	if len(msg.FromAddress) != 20 {
		return &PluginCheckResponse{Error: ErrInvalidAddress()}
	}
	if len(msg.ToAddress) != 20 {
		return &PluginCheckResponse{Error: ErrInvalidAddress()}
	}
	if msg.Amount == 0 {
		return &PluginCheckResponse{Error: ErrInvalidAmount()}
	}
	return &PluginCheckResponse{Recipient: msg.ToAddress, AuthorizedSigners: [][]byte{msg.FromAddress}}
}

func (c *Contract) DeliverMessageSend(msg *MessageSend, fee uint64) *PluginDeliverResponse {
	log.Printf("DeliverMessageSend: from=%x to=%x amount=%d fee=%d", msg.FromAddress, msg.ToAddress, msg.Amount, fee)
	var (
		fromKey, toKey, feePoolKey         []byte
		fromBytes, toBytes, feePoolBytes   []byte
		fromQueryId, toQueryId, feeQueryId = rand.Uint64(), rand.Uint64(), rand.Uint64()
		from, to, feePool                  = new(Account), new(Account), new(Pool)
	)
	fromKey = KeyForAccount(msg.FromAddress)
	toKey = KeyForAccount(msg.ToAddress)
	feePoolKey = KeyForFeePool(c.Config.ChainId)
	response, err := c.plugin.StateRead(c, &PluginStateReadRequest{
		Keys: []*PluginKeyRead{
			{QueryId: feeQueryId, Key: feePoolKey},
			{QueryId: fromQueryId, Key: fromKey},
			{QueryId: toQueryId, Key: toKey},
		},
	})
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if response.Error != nil {
		return &PluginDeliverResponse{Error: response.Error}
	}
	for _, resp := range response.Results {
		if len(resp.Entries) == 0 {
			continue
		}
		switch resp.QueryId {
		case fromQueryId:
			fromBytes = resp.Entries[0].Value
		case toQueryId:
			toBytes = resp.Entries[0].Value
		case feeQueryId:
			feePoolBytes = resp.Entries[0].Value
		}
	}
	amountToDeduct := msg.Amount + fee
	if err = Unmarshal(fromBytes, from); err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if err = Unmarshal(toBytes, to); err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if err = Unmarshal(feePoolBytes, feePool); err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if from.Amount < amountToDeduct {
		return &PluginDeliverResponse{Error: ErrInsufficientFunds()}
	}
	if bytes.Equal(fromKey, toKey) {
		to = from
	}
	from.Amount -= amountToDeduct
	feePool.Amount += fee
	to.Amount += msg.Amount
	fromBytes, err = Marshal(from)
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	toBytes, err = Marshal(to)
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	feePoolBytes, err = Marshal(feePool)
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	var resp *PluginStateWriteResponse
	if from.Amount == 0 {
		resp, err = c.plugin.StateWrite(c, &PluginStateWriteRequest{
			Sets:    []*PluginSetOp{{Key: feePoolKey, Value: feePoolBytes}, {Key: toKey, Value: toBytes}},
			Deletes: []*PluginDeleteOp{{Key: fromKey}},
		})
	} else {
		resp, err = c.plugin.StateWrite(c, &PluginStateWriteRequest{
			Sets: []*PluginSetOp{
				{Key: feePoolKey, Value: feePoolBytes},
				{Key: toKey, Value: toBytes},
				{Key: fromKey, Value: fromBytes},
			},
		})
	}
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	return &PluginDeliverResponse{Error: resp.Error}
}

var (
	accountPrefix = []byte{1}
	poolPrefix    = []byte{2}
	paramsPrefix  = []byte{7}
)

func KeyForAccount(addr []byte) []byte {
	return JoinLenPrefix(accountPrefix, addr)
}

func KeyForFeeParams() []byte {
	return JoinLenPrefix(paramsPrefix, []byte("/f/"))
}

func KeyForFeePool(chainId uint64) []byte {
	return JoinLenPrefix(poolPrefix, formatUint64(chainId))
}

func formatUint64(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}
