package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
)

const (
	MgWithdrawStart = "Start"
	MgWithdrawReady = "Ready"
	MgWithdrawDone  = "Done"
)

type MgWithdraw struct {
	ID      string `bson:"_id"`
	ShopID  string `bson:"shop_id"`
	Comment string `bson:"comment"`
	Sum     int64  `bson:"sum"`
	Status  string `bson:"status"`
}

func (w *MgWithdraw) ToDomain() domain.Withdraw {
	var withdrawStatus domain.WithdrawStatus
	switch w.Status {
	case MgWithdrawStart:
		withdrawStatus = domain.WithdrawStatusStart
	case MgWithdrawReady:
		withdrawStatus = domain.WithdrawStatusReady
	case MgWithdrawDone:
		withdrawStatus = domain.WithdrawStatusDone
	}

	return domain.Withdraw{
		ID:      domain.ID(w.ID),
		ShopID:  domain.ID(w.ShopID),
		Comment: w.Comment,
		Sum:     w.Sum,
		Status:  withdrawStatus,
	}
}

func NewMgWithdraw(withdraw domain.Withdraw) MgWithdraw {
	var withdrawStatus string
	switch withdraw.Status {
	case domain.WithdrawStatusStart:
		withdrawStatus = MgWithdrawStart
	case domain.WithdrawStatusReady:
		withdrawStatus = MgWithdrawReady
	case domain.WithdrawStatusDone:
		withdrawStatus = MgWithdrawDone
	}

	return MgWithdraw{
		ID:      withdraw.ID.String(),
		ShopID:  withdraw.ShopID.String(),
		Comment: withdraw.Comment,
		Sum:     withdraw.Sum,
		Status:  withdrawStatus,
	}
}
