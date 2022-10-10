package main

import (
	"fmt"
	"go-designpattern/template/method"
)

func main() {
	smsOTP := &method.Sms{}
	o := method.Otp{IOtp: smsOTP}
	o.GenAndSendOTP(4)

	fmt.Println()
	emailOTP := &method.Email{}
	o = method.Otp{IOtp: emailOTP}
	o.GenAndSendOTP(4)

}
