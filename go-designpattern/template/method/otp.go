package method

// IOtp 一次性密码
type IOtp interface {
	GenRandomOTP(int) string
	saveOTPCache(string)
	getMessage(string) string
	sendNotification(string) error
}

type Otp struct {
	IOtp IOtp
}

// GenAndSendOTP
// 1. 生成随机的n位数字
// 2. 在缓存中保存这组数字以便进行后续验证
// 3. 准备内容
// 4. 发送通知
func (o *Otp) GenAndSendOTP(otpLength int) error {
	otp := o.IOtp.GenRandomOTP(otpLength)
	o.IOtp.saveOTPCache(otp)
	message := o.IOtp.getMessage(otp)
	if err := o.IOtp.sendNotification(message); err != nil {
		return err
	}
	return nil
}
