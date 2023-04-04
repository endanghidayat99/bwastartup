package payment

import (
	"bwastartup/campaign"
	"bwastartup/transaction"
	"bwastartup/user"
	midtrans "github.com/veritrans/go-midtrans"
	"strconv"
)

type service struct {
	transactionRepository transaction.Repository
	campaignRepository    campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(transactionRepository transaction.Repository, campaignRepository campaign.Repository) *service {
	return &service{transactionRepository, campaignRepository}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = ""
	midclient.ClientKey = ""
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}
	println(user.Name + " " + user.Email)
	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
			LName: "",
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil

}

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error {
	transactionId, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.transactionRepository.GetByID(transactionId)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTrx, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTrx.CampaignID)
	if err != nil {
		return err
	}

	if updatedTrx.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount += updatedTrx.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
