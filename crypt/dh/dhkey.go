package dh

import (
	"encoding/json"
	"errors"
	"github.com/monnand/dhkx"
)

type Dh struct {
	privateKey *dhkx.DHKey
	PublicKey  []byte
	group      *dhkx.DHGroup
	finalKey   *dhkx.DHKey
}

func (d *Dh) GetFinalKey() string {
	return d.finalKey.String()
}

func (d *Dh) Configure() {
	group, err := dhkx.GetGroup(0)
	if err != nil {
		return
	}
	key, err := group.GeneratePrivateKey(nil)
	if err != nil {
		return
	}
	d.privateKey = key
	d.group = group
	d.PublicKey = key.Bytes()
}

func (d *Dh) ReceiveJson(data []byte) error {
	if d.group == nil && d.privateKey == nil && d.PublicKey == nil {
		d.Configure()
	}
	n := new(Dh)
	err := json.Unmarshal(data, n)
	if err != nil {
		return err
	}
	if len(n.PublicKey) < 1 {
		return errors.New("Json Received doesn't haves a publicKey.")
	}
	d.finalKey, err = d.group.ComputeKey(dhkx.NewPublicKey(n.PublicKey), d.privateKey)
	return err
}
func (d *Dh) MakeJson() []byte {
	data, _ := json.Marshal(d)
	return data
}

func ExampleMain() {
	bob := new(Dh)                    // Create Bob
	bob.Configure()                   // First Config G,P, Private and Public Key
	bobJson := bob.MakeJson()         // Make a Json
	alice := new(Dh)                  // Create Alice
	err := alice.ReceiveJson(bobJson) // Receive configuration from Bob Json
	if err != nil {
		return
	}
	aliceJson := alice.MakeJson()    // Make Alice Json because she haves now a publicKey and FinalKey
	err = bob.ReceiveJson(aliceJson) // Bob receives the PublicKey of Alice and calculates the FinalKey
	if err != nil {
		return
	}
	println(bob.GetFinalKey() == alice.GetFinalKey(), "\n"+bob.GetFinalKey()+"\n"+alice.GetFinalKey())
	// Synchronized
}
