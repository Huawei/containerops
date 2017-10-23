package module

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/Huawei/containerops/common"

	"github.com/scorredoira/email"
	"strings"
	"time"
)

func init() {
	Register("mail", &MailNotifier{})
}

type MailNotifier struct {
}

func (m *MailNotifier) Notify(flow *Flow, receivers []string) error {

	subject := fmt.Sprintf("[ContainerOps] Excution Result of Flow: %s is [%s] ", flow.URI, strings.ToUpper(flow.Status))
	htmlBody := fmt.Sprintf("Flow URI: %s <br /> Tag: %s <br /> Title: %s <br /> Result: %s", flow.URI, flow.Tag, flow.Title, flow.Status)
	msg := email.NewHTMLMessage(subject, htmlBody)
	msg.From = mail.Address{Name: "ContainerOps", Address: common.Mail.User}
	msg.To = receivers

	//attach log to attachments
	fileName := fmt.Sprintf("log-%s:%s-%s.txt", flow.URI, flow.Tag, time.Now().Format("20060102150405"))
	buf := bytes.NewBuffer(nil)
	for _, logLine := range flow.Logs {
		buf.WriteString("\r\n")
		buf.WriteString(logLine)
	}
	data := buf.Bytes()
	msg.Attachments[fileName] = &email.Attachment{
		Filename: fileName,
		Inline:   false,
		Data:     data,
	}

	smtpAddress := common.Mail.SmtpAddress
	smtpPort := common.Mail.SmtpPort
	user := common.Mail.User
	password := common.Mail.Password
	auth := smtp.PlainAuth("", user, password, smtpAddress)
	if err := email.Send(fmt.Sprintf("%s:%s", smtpAddress, smtpPort), auth, msg); err != nil {
		return err
	}

	return nil
}
