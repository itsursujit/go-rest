package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

func main() {
	enc := data.FindEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ")
	fmt.Println(enc)
	/*var wg sync.WaitGroup

	wg.Add(1)
	enc := data.FindEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ")
	fmt.Println(enc.)
	// go sendingAndReceiveSMS(&wg)

	wg.Wait()*/
}

func sendingAndReceiveSMS(wg *sync.WaitGroup) {
	defer wg.Done()

	auth := gosmpp.Auth{
		SMSC:       "127.0.0.1:2775",
		SystemID:   "test",
		Password:   "UUDHWB",
		SystemType: "",
	}

	trans, err := gosmpp.NewTransceiverSession(gosmpp.NonTLSDialer, auth, gosmpp.TransceiveSettings{
		EnquireLink: 5 * time.Second,

		OnSubmitError: func(p pdu.PDU, err error) {
			log.Fatal(err)
		},

		OnReceivingError: func(err error) {
			fmt.Println(err)
		},

		OnRebindingError: func(err error) {
			fmt.Println(err)
		},

		OnPDU: handlePDU(),

		OnClosed: func(state gosmpp.State) {
			fmt.Println(state)
		},
	}, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = trans.Close()
	}()

	// sending SMS(s)
	for {
		if err = trans.Transceiver().Submit(newSubmitSM()); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

func handlePDU() func(pdu.PDU, bool) {
	return func(p pdu.PDU, responded bool) {
		switch pd := p.(type) {
		case *pdu.SubmitSMResp:
			fmt.Printf("SubmitSMResp:%+v\n", pd)

		case *pdu.GenerickNack:
			fmt.Println("GenericNack Received")

		case *pdu.EnquireLinkResp:
			fmt.Println("EnquireLinkResp Received")

		case *pdu.DataSM:
			fmt.Printf("DataSM:%+v\n", pd)

		case *pdu.DeliverSM:
			fmt.Printf("DeliverSM:%+v\n", pd)
			fmt.Println(pd.Message.GetMessage())
		}
	}
}

func newSubmitSM() *pdu.SubmitSM {

	///// SENDER / SOURCE
	// If sender address is shortcode (length is 5 digits or less)
	// TON = 3
	// NPI = 0

	// If sender is Non-Numeric
	// TON = 5
	// NPI = 0

	// If sender address starts with `+`
	// TON = 1
	// NPI = 1

	// Nothing Above
	// TON = 0
	// NPI = 1

	///// DESTINATION
	// If recipient starts with `+`
	// TON = 1
	// NPI = 1

	// Nothing Above
	// TON = 0
	// NPI = 1

	// The SMPP specification defines the following TON values:

	// Unknown = 0
	// International = 1
	// National = 2
	// Network Specific = 3
	// Subscriber Number = 4
	// Alphanumeric = 5
	// Abbreviated = 6

	// Possible NPI values are defines as follows:

	// Unknown = 0
	// ISDN/telephone numbering plan (E163/E164) = 1
	// Data numbering plan (X.121) = 3
	// Telex numbering plan (F.69) = 4
	// Land Mobile (E.212) =6
	// National numbering plan = 8
	// Private numbering plan = 9
	// ERMES numbering plan (ETSI DE/PS 3 01-3) = 10
	// Internet (IP) = 13
	// WAP Client Id (to be defined by WAP Forum) = 18

	// build up submitSM
	srcAddr := pdu.NewAddress()
	srcAddr.SetTon(5)
	srcAddr.SetNpi(0)
	_ = srcAddr.SetAddress("00" + "522241")

	destAddr := pdu.NewAddress()
	destAddr.SetTon(1)
	destAddr.SetNpi(1)
	_ = destAddr.SetAddress("99" + "522241")

	submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)
	submitSM.SourceAddr = srcAddr
	submitSM.DestAddr = destAddr
	_ = submitSM.Message.SetMessageWithEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ", data.UCS2)
	submitSM.ProtocolID = 0
	submitSM.RegisteredDelivery = 1
	submitSM.ReplaceIfPresentFlag = 0
	submitSM.EsmClass = 0

	return submitSM
}
