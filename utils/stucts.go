package utils

import "math/big"

type MSG struct {
	Msg []byte `json:"msg"`
}

type GamalMSG struct {
	C1, C2 *big.Int
}
