package evm

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/0xPolygon/polygon-edge/chain"
	"github.com/0xPolygon/polygon-edge/crypto"
	"github.com/0xPolygon/polygon-edge/state/runtime"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

/*
func TestG2MapToCurve(t *testing.T) {
	for _, v := range []struct {
		u []byte
	}{
		{u: make([]byte, 96)},
		{u: common.FromHex("025fbc07711ba267b7e70c82caa70a16fbb1d470ae24ceef307f5e2000751677820b7013ad4e25492dcf30052d3e5eca" + "0e775d7827adf385b83e20e4445bd3fab21d7b4498426daf3c1d608b9d41e9edb5eda0df022e753b8bb4bc3bb7db4914")},
		{u: common.FromHex("1870a7dbfd2a1deb74015a3546b20f598041bf5d5202997956a94a368d30d3f70f18cdaa1d33ce970a4e16af961cbdcb" + "045ab31ce4b5a8ba7c4b2851b64f063a66cd1223d3c85005b78e1beee65e33c90ceef0244e45fc45a5e1d6eab6644fdb")},
		{u: common.FromHex("088fe329b054db8a6474f21a7fbfdf17b4c18044db299d9007af582c3d5f17d00e56d99921d4b5640fce44b05219b5de" + "0b6e6135a4cd31ba980ddbd115ac48abef7ec60e226f264d7befe002c165f3a496f36f76dd524efd75d17422558d10b4")},
		{u: common.FromHex("03df16a66a05e4c1188c234788f43896e0565bfb64ac49b9639e6b284cc47dad73c47bb4ea7e677db8d496beb907fbb6" + "0f45b50647d67485295aa9eb2d91a877b44813677c67c8d35b2173ff3ba95f7bd0806f9ca8a1436b8b9d14ee81da4d7e")},
	} {
		g := bls12381.NewG2()
		p0, err := g.MapToCurve(v.u)
		if err != nil {
			t.Fatal("map to curve fails", err)
		}
		fmt.Println("Result:", g.ToBytes(p0))
	}
}
*/

func TestG2MapToCurve(t *testing.T) {
	for _, v := range []struct {
		u []byte
	}{
		{u: make([]byte, 96)},
		{u: common.FromHex("025fbc07711ba267b7e70c82caa70a16fbb1d470ae24ceef307f5e2000751677820b7013ad4e25492dcf30052d3e5eca" + "0e775d7827adf385b83e20e4445bd3fab21d7b4498426daf3c1d608b9d41e9edb5eda0df022e753b8bb4bc3bb7db4914")},
		{u: common.FromHex("1870a7dbfd2a1deb74015a3546b20f598041bf5d5202997956a94a368d30d3f70f18cdaa1d33ce970a4e16af961cbdcb" + "045ab31ce4b5a8ba7c4b2851b64f063a66cd1223d3c85005b78e1beee65e33c90ceef0244e45fc45a5e1d6eab6644fdb")},
		{u: common.FromHex("088fe329b054db8a6474f21a7fbfdf17b4c18044db299d9007af582c3d5f17d00e56d99921d4b5640fce44b05219b5de" + "0b6e6135a4cd31ba980ddbd115ac48abef7ec60e226f264d7befe002c165f3a496f36f76dd524efd75d17422558d10b4")},
		{u: common.FromHex("03df16a66a05e4c1188c234788f43896e0565bfb64ac49b9639e6b284cc47dad73c47bb4ea7e677db8d496beb907fbb6" + "0f45b50647d67485295aa9eb2d91a877b44813677c67c8d35b2173ff3ba95f7bd0806f9ca8a1436b8b9d14ee81da4d7e")},
	} {
		// 새로운 state 객체 생성
		c, closeFn := getState()
		defer closeFn()

		// 입력값을 big.Int로 변환하고 스택에 푸시
		c.push(new(big.Int).SetBytes(v.u))

		// opMapToCurve 함수 실행
		opMapToCurve(c)

		// 결과값 추출
		resultBigInt := c.pop()
		// if resultBigInt == nil {
		// 	t.Fatal("No result on the stack")
		// }

		// 결과 출력 (선택사항)
		fmt.Println("Result:", resultBigInt)

		// state 객체 반환
		releaseState(c)
	}
}

var (
	two = big.NewInt(2)

	allEnabledForks = chain.AllForksEnabled.At(0)
)

type cases2To1 []struct {
	a *big.Int
	b *big.Int
	c *big.Int
}

func test2to1(t *testing.T, f instruction, tests cases2To1) {
	t.Helper()

	s, closeFn := getState()
	defer closeFn()

	for _, i := range tests {
		s.push(i.a)
		s.push(i.b)

		f(s)

		assert.Equal(t, i.c, s.pop())
	}
}

type cases2ToBool []struct {
	a *big.Int
	b *big.Int
	c bool
}

func test2toBool(t *testing.T, f instruction, tests cases2ToBool) {
	t.Helper()

	s, closeFn := getState()
	defer closeFn()

	for _, i := range tests {
		s.push(i.a)
		s.push(i.b)

		f(s)

		if i.c {
			assert.Equal(t, uint64(1), s.pop().Uint64())
		} else {
			assert.Equal(t, uint64(0), s.pop().Uint64())
		}
	}
}

func TestAdd(t *testing.T) {
	test2to1(t, opAdd, cases2To1{
		{one, one, two},
		{zero, one, one},
	})
}

func TestGt(t *testing.T) {
	test2toBool(t, opGt, cases2ToBool{
		{one, one, false},
		{two, one, false},
		{one, two, true},
	})
}

func TestIsZero(t *testing.T) {
	test2toBool(t, opIsZero, cases2ToBool{
		{one, one, false},
		{zero, zero, true},
		{two, two, false},
	})
}

func TestMStore(t *testing.T) {
	s, closeFn := getState()
	defer closeFn()

	s.push(big.NewInt(10))   // value
	s.push(big.NewInt(1024)) // offset

	s.gas = 1000
	opMStore(s)

	assert.Len(t, s.memory, 1024+32)
}

type mockHostForInstructions struct {
	mockHost
	nonce       uint64
	code        []byte
	callxResult *runtime.ExecutionResult
}

func (m *mockHostForInstructions) GetNonce(types.Address) uint64 {
	return m.nonce
}

func (m *mockHostForInstructions) Callx(*runtime.Contract, runtime.Host) *runtime.ExecutionResult {
	return m.callxResult
}

func (m *mockHostForInstructions) GetCode(addr types.Address) []byte {
	return m.code
}

var (
	addr1 = types.StringToAddress("1")
)

func TestCreate(t *testing.T) {
	type state struct {
		gas    uint64
		sp     int
		stack  []*big.Int
		memory []byte
		stop   bool
		err    error
	}

	addressToBigInt := func(addr types.Address) *big.Int {
		return new(big.Int).SetBytes(addr[:])
	}

	tests := []struct {
		name        string
		op          OpCode
		contract    *runtime.Contract
		config      *chain.ForksInTime
		initState   *state
		resultState *state
		mockHost    *mockHostForInstructions
	}{
		{
			name: "should succeed in case of CREATE",
			op:   CREATE,
			contract: &runtime.Contract{
				Static:  false,
				Address: addr1,
			},
			config: &chain.ForksInTime{},
			initState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
			},
			resultState: &state{
				gas: 500,
				sp:  1,
				stack: []*big.Int{
					addressToBigInt(crypto.CreateAddress(addr1, 0)), // contract address
					big.NewInt(0x00),
					big.NewInt(0x00),
				},
				memory: []byte{
					byte(REVERT),
				},
			},
			mockHost: &mockHostForInstructions{
				nonce: 0,
				callxResult: &runtime.ExecutionResult{
					GasLeft: 500,
					GasUsed: 500,
				},
			},
		},
		{
			name: "should throw errWriteProtection in case of static call",
			op:   CREATE,
			contract: &runtime.Contract{
				Static: true,
			},
			config: &chain.ForksInTime{},
			initState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			// shouldn't change any states except for stop and err
			resultState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: true,
				err:  errWriteProtection,
			},
			mockHost: &mockHostForInstructions{},
		},
		{
			name:     "should throw errOpCodeNotFound when op is CREATE2 and config.Constantinople is disabled",
			op:       CREATE2,
			contract: &runtime.Contract{},
			config: &chain.ForksInTime{
				Constantinople: false,
			},
			initState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			// shouldn't change any states except for stop and err
			resultState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: true,
				err:  errOpCodeNotFound,
			},
			mockHost: &mockHostForInstructions{},
		},
		{
			name: "should set zero address if op is CREATE and contract call throws ErrCodeStoreOutOfGas",
			op:   CREATE,
			contract: &runtime.Contract{
				Static:  false,
				Address: addr1,
			},
			config: &chain.ForksInTime{
				Homestead: true,
			},
			initState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			// shouldn't change any states except for stop and err
			resultState: &state{
				gas: 1000,
				sp:  1,
				stack: []*big.Int{
					// need to init with 0x01 to add abs field in big.Int
					big.NewInt(0x01).SetInt64(0x00),
					big.NewInt(0x00),
					big.NewInt(0x00),
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			mockHost: &mockHostForInstructions{
				nonce: 0,
				callxResult: &runtime.ExecutionResult{
					GasLeft: 1000,
					Err:     runtime.ErrCodeStoreOutOfGas,
				},
			},
		},
		{
			name: "should set zero address if contract call throws error except for ErrCodeStoreOutOfGas",
			op:   CREATE,
			contract: &runtime.Contract{
				Static:  false,
				Address: addr1,
			},
			config: &chain.ForksInTime{
				Homestead: true,
			},
			initState: &state{
				gas: 1000,
				sp:  3,
				stack: []*big.Int{
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			// shouldn't change any states except for stop and err
			resultState: &state{
				gas: 1000,
				sp:  1,
				stack: []*big.Int{
					// need to init with 0x01 to add abs field in big.Int
					big.NewInt(0x01).SetInt64(0x00),
					big.NewInt(0x00),
					big.NewInt(0x00),
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			mockHost: &mockHostForInstructions{
				nonce: 0,
				callxResult: &runtime.ExecutionResult{
					GasLeft: 1000,
					Err:     errRevert,
				},
			},
		},
		{
			name: "should set zero address if contract call throws any error for CREATE2",
			op:   CREATE2,
			contract: &runtime.Contract{
				Static:  false,
				Address: addr1,
			},
			config: &chain.ForksInTime{
				Homestead:      true,
				Constantinople: true,
			},
			initState: &state{
				gas: 1000,
				sp:  4,
				stack: []*big.Int{
					big.NewInt(0x01), // salt
					big.NewInt(0x01), // length
					big.NewInt(0x00), // offset
					big.NewInt(0x00), // value
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			// during creation of code with length 1 for CREATE2 op code, 985 gas units are spent by buildCreateContract()
			resultState: &state{
				gas: 15,
				sp:  1,
				stack: []*big.Int{
					big.NewInt(0x01).SetInt64(0x00),
					big.NewInt(0x01),
					big.NewInt(0x00),
					big.NewInt(0x00),
				},
				memory: []byte{
					byte(REVERT),
				},
				stop: false,
				err:  nil,
			},
			mockHost: &mockHostForInstructions{
				nonce: 0,
				callxResult: &runtime.ExecutionResult{
					// if it is ErrCodeStoreOutOfGas then we set GasLeft to 0
					GasLeft: 0,
					Err:     runtime.ErrCodeStoreOutOfGas,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, closeFn := getState()
			defer closeFn()

			s.msg = tt.contract
			s.gas = tt.initState.gas
			s.sp = tt.initState.sp
			s.stack = tt.initState.stack
			s.memory = tt.initState.memory
			s.config = tt.config
			s.host = tt.mockHost

			opCreate(tt.op)(s)

			assert.Equal(t, tt.resultState.gas, s.gas, "gas in state after execution is not correct")
			assert.Equal(t, tt.resultState.sp, s.sp, "sp in state after execution is not correct")
			assert.Equal(t, tt.resultState.stack, s.stack, "stack in state after execution is not correct")
			assert.Equal(t, tt.resultState.memory, s.memory, "memory in state after execution is not correct")
			assert.Equal(t, tt.resultState.stop, s.stop, "stop in state after execution is not correct")
			assert.Equal(t, tt.resultState.err, s.err, "err in state after execution is not correct")
		})
	}
}

func Test_opReturnDataCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      *chain.ForksInTime
		initState   *state
		resultState *state
	}{
		{
			name: "should return error if Byzantium is not applied",
			config: &chain.ForksInTime{
				Byzantium: false,
			},
			initState: &state{},
			resultState: &state{
				config: &chain.ForksInTime{
					Byzantium: false,
				},
				stop: true,
				err:  errOpCodeNotFound,
			},
		},
		{
			name:   "should return error if memOffset is negative",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(1),  // length
					big.NewInt(0),  // dataOffset
					big.NewInt(-1), // memOffset
				},
				sp: 3,
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(1),
					big.NewInt(0),
					big.NewInt(-1),
				},
				sp:   0,
				stop: true,
				err:  errReturnDataOutOfBounds,
			},
		},
		{
			name:   "should return error if dataOffset is negative",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(1),  // length
					big.NewInt(-1), // dataOffset
					big.NewInt(0),  // memOffset
				},
				sp:     3,
				memory: make([]byte, 1),
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(1),
					big.NewInt(-1),
					big.NewInt(0),
				},
				sp:     0,
				memory: make([]byte, 1),
				stop:   true,
				err:    errReturnDataOutOfBounds,
			},
		},
		{
			name:   "should return error if length is negative",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(-1), // length
					big.NewInt(0),  // dataOffset
					big.NewInt(0),  // memOffset
				},
				sp: 3,
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(-1),
					big.NewInt(0),
					big.NewInt(0),
				},
				sp:   0,
				stop: true,
				err:  errReturnDataOutOfBounds,
			},
		},
		{
			name:   "should copy data from returnData to memory",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(1), // length
					big.NewInt(0), // dataOffset
					big.NewInt(0), // memOffset
				},
				sp:         3,
				returnData: []byte{0xff},
				memory:     []byte{0x0},
				gas:        10,
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(1),
					big.NewInt(0),
					big.NewInt(0),
				},
				sp:                 0,
				returnData:         []byte{0xff},
				memory:             []byte{0xff},
				gas:                7,
				lastGasCost:        0,
				currentConsumedGas: 3,
				stop:               false,
				err:                nil,
			},
		},
		{
			name:   "should expand memory and copy data returnData",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(5), // length
					big.NewInt(1), // dataOffset
					big.NewInt(2), // memOffset
				},
				sp:         3,
				returnData: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				memory:     []byte{0x11, 0x22},
				gas:        20,
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(6), // updated for end index
					big.NewInt(1),
					big.NewInt(2),
				},
				sp:         0,
				returnData: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				memory: append(
					// 1 word (32 bytes)
					[]byte{0x11, 0x22, 0x02, 0x03, 0x04, 0x05, 0x06},
					make([]byte, 25)...,
				),
				gas:                14,
				lastGasCost:        3,
				currentConsumedGas: 6,
				stop:               false,
				err:                nil,
			},
		},
		{
			// this test case also verifies that code does not panic when the length is 0 and memOffset > len(memory)
			name:   "should not copy data if length is zero",
			config: &allEnabledForks,
			initState: &state{
				stack: []*big.Int{
					big.NewInt(0), // length
					big.NewInt(0), // dataOffset
					big.NewInt(4), // memOffset
				},
				sp:         3,
				returnData: []byte{0x01},
				memory:     []byte{0x02},
			},
			resultState: &state{
				config: &allEnabledForks,
				stack: []*big.Int{
					big.NewInt(0),
					big.NewInt(0),
					big.NewInt(4),
				},
				sp:         0,
				returnData: []byte{0x01},
				memory:     []byte{0x02},
				stop:       false,
				err:        nil,
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			state, closeFn := getState()
			defer closeFn()

			state.gas = test.initState.gas
			state.sp = test.initState.sp
			state.stack = test.initState.stack
			state.memory = test.initState.memory
			state.returnData = test.initState.returnData
			state.config = test.config

			// assign nil to some fields in cached state object
			state.code = nil
			state.host = nil
			state.msg = nil
			state.evm = nil
			state.bitmap = bitmap{}
			state.ret = nil
			state.currentConsumedGas = 0

			opReturnDataCopy(state)

			assert.Equal(t, test.resultState, state)
		})
	}
}

func Test_opCall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		op          OpCode
		contract    *runtime.Contract
		config      chain.ForksInTime
		initState   *state
		resultState *state
		mockHost    *mockHostForInstructions
	}{
		{
			// this test case also verifies that code does not panic when the outSize is 0 and outOffset > len(memory)
			name: "should not copy result into memory if outSize is 0",
			op:   STATICCALL,
			contract: &runtime.Contract{
				Static: true,
			},
			config: allEnabledForks,
			initState: &state{
				gas: 1000,
				sp:  6,
				stack: []*big.Int{
					big.NewInt(0x00), // outSize
					big.NewInt(0x02), // outOffset
					big.NewInt(0x00), // inSize
					big.NewInt(0x00), // inOffset
					big.NewInt(0x00), // address
					big.NewInt(0x00), // initialGas
				},
				memory: []byte{0x01},
			},
			resultState: &state{
				memory: []byte{0x01},
				stop:   false,
				err:    nil,
				gas:    300,
			},
			mockHost: &mockHostForInstructions{
				callxResult: &runtime.ExecutionResult{
					ReturnValue: []byte{0x03},
				},
			},
		},
		{
			name: "call cost overflow (EIP150 fork disabled)",
			op:   CALLCODE,
			contract: &runtime.Contract{
				Static: false,
			},
			config: chain.AllForksEnabled.RemoveFork(chain.EIP150).At(0),
			initState: &state{
				gas: 6640,
				sp:  7,
				stack: []*big.Int{
					big.NewInt(0x00),                        // outSize
					big.NewInt(0x00),                        // outOffset
					big.NewInt(0x00),                        // inSize
					big.NewInt(0x00),                        // inOffset
					big.NewInt(0x01),                        // value
					big.NewInt(0x03),                        // address
					big.NewInt(0).SetUint64(math.MaxUint64), // initialGas
				},
				memory: []byte{0x01},
			},
			resultState: &state{
				memory: []byte{0x01},
				stop:   true,
				err:    errGasUintOverflow,
				gas:    6640,
			},
			mockHost: &mockHostForInstructions{
				callxResult: &runtime.ExecutionResult{
					ReturnValue: []byte{0x03},
				},
			},
		},
		{
			name: "available gas underflow",
			op:   CALLCODE,
			contract: &runtime.Contract{
				Static: false,
			},
			config: allEnabledForks,
			initState: &state{
				gas: 6640,
				sp:  7,
				stack: []*big.Int{
					big.NewInt(0x00),                        // outSize
					big.NewInt(0x00),                        // outOffset
					big.NewInt(0x00),                        // inSize
					big.NewInt(0x00),                        // inOffset
					big.NewInt(0x01),                        // value
					big.NewInt(0x03),                        // address
					big.NewInt(0).SetUint64(math.MaxUint64), // initialGas
				},
				memory: []byte{0x01},
			},
			resultState: &state{
				memory: []byte{0x01},
				stop:   true,
				err:    errOutOfGas,
				gas:    6640,
			},
			mockHost: &mockHostForInstructions{
				callxResult: &runtime.ExecutionResult{
					ReturnValue: []byte{0x03},
				},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			state, closeFn := getState()
			defer closeFn()

			state.gas = test.initState.gas
			state.msg = test.contract
			state.sp = test.initState.sp
			state.stack = test.initState.stack
			state.memory = test.initState.memory
			state.config = &test.config
			state.host = test.mockHost

			opCall(test.op)(state)

			assert.Equal(t, test.resultState.memory, state.memory, "memory in state after execution is incorrect")
			assert.Equal(t, test.resultState.stop, state.stop, "stop in state after execution is incorrect")
			assert.Equal(t, test.resultState.err, state.err, "err in state after execution is incorrect")
			assert.Equal(t, test.resultState.gas, state.gas, "gas in state after execution is incorrect")
		})
	}
}
