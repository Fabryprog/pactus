package executor

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type SendExecutor struct {
	sandbox Sandbox
}

func NewSendExecutor(sandbox Sandbox) *SendExecutor {
	return &SendExecutor{sandbox}
}

func (e *SendExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.SendPayload)

	senderAcc := e.sandbox.Account(pld.SenderAddr)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve sender account")
	}
	receiverAcc := e.sandbox.Account(pld.ReceiverAddr)
	if receiverAcc == nil {
		receiverAcc = account.NewAccount(pld.ReceiverAddr)
	}
	if senderAcc.Balance() < pld.Amount+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	e.sandbox.UpdateAccount(senderAcc)
	e.sandbox.UpdateAccount(receiverAcc)

	return nil
}
