package midtrans

import (
	"os"

	"github.com/midtrans/midtrans-go"
)

func MidtransInit() {
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
}
