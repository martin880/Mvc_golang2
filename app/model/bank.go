package model

import (
	"fmt"
	"mvc_golang/app/utils"

	"gorm.io/gorm"
)

type Account struct {
	ID            int           `gorm:"primary_key" json:"-"`
	IdAccount     string        `json:"id_account"`
	Name          string        `json:"name"`
	AccountNumber int           `json:"account_number"`
	Saldo         int           `json:"saldo"`
	Transaction   []Transaction `gorm:"ForeignKey:IdAccountRefer" json:"transaction"`
}

type Transaction struct {
	ID                     int    `gorm:"primary_key" json:"-"`
	IdAccountRefer         int    `json:"-"`
	IdTransaction          string `json:"id_transaction"`
	TransactionType        int    `json:"transaction_type"`
	TransactionDescription string `json:"transaction_description"`
	Sender                 int    `json:"sender"`
	Recipient              int    `json:"recipient"`
	Timestamp              int    `json:"timestamp"`
}

func Login(auth Auth) (bool, error, string) {
	var account Account
	if err := DB.Where(&Account{Name: auth.Name}).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.Errorf("Account not found"), ""
		}
	}
	err := utils.HashComparator([]byte(account.Password), []byte(auth.Password))
	if err != nil {
		return false, errors.Errorf("Incorrect Password"), ""
	} else {
		sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name":           auth.Name,
			"account_number": account.AccountNumber,
		})

		token, err := sign.SignedString([]byte("secret"))
		if err != nil {
			return false, err, ""
		}
		return true, nil, token
	}
}

func InsertNewAccount(account Account) (bool, error) {
	account.AccountNumber = utils.RangeIn(111111, 999999)
	account.Saldo = 0
	account.IdAccount = fmt.Sprintf("id-%id", utils.RangeIn(111, 999))
	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+V\n", err)
	}
	return true, nil
}

func GetAccountDetail(idAccount int) (bool,error, []Transaction,Account){
	var transaction []Transaction
	var account Account
	if err := DB.Where("sender = ? OR recipient = ?",idAccount,idAccount)
	Find(&transaction).Error;err!=nil{
		if err == gorm.ErrRecordNotFound{
			return false, errors.Errorf("Account not found"),[]Transaction{},Account{}
		}else{
			return false, errors.Errorf("invalid prepare statement :%+v\n", err),[]Transaction{},Account{}
		}
	}
	
	if err := DB.Where(&Account{AccountNumber:idAccount}).Find(&account).Error;err != nil{
		if err == gorm.ErrRecordNotFound{
			return false,errors.Errorf("Account not found"),[]Transaction{},Account{}
		}else{
			return false, errors.Errorf("invalid prepare statement:%+v\n", err),[]Transaction{}.Account{}
		}
	}

	return true,nil,transaction,Account{
		IdAccount:		account.IdAccount,
		Name:			account.Name,
		AccountNumber:	account.AccountNumber,
		Saldo:			account.Saldo,
	}
}

func Transfer (transaction Transaction) (bool, error){
	err := DB.Transaction(func(tx *gorm.DB)error{
		var sender, recipient Account
		if err := tx.Model(&Account{}).Where(&Account{AccountNumber:transaction.Sender})
		First(&sender)
		Update("saldo", sender.Saldo-transaction.Amount).Error;err != nil {
			return err
		}
		if err := tx.Model(&Account{}).Where(&Account{AccountNumber:transaction.Recipient})
		First(&recipient)
		Update("saldo", recipient.Saldo+transaction.Amount).Error:err != nil
		{
			log.Println("ERROR : " + err.Error())
			return err
		}
		transaction.TransactionType = constant.TRANSFER
		transaction.Timestamp = time.Now().Unix()
		if err := tx.Create(&transcation).Error;err != nil {
			return err
		}
		return nil
	});if err != nil {
		return false, err
	}
	return true,nil
}

func Withdraw (transaction Transaction)(bool,error){
	err := DB.Transaction(func(tx *gorm.DB) error{
		var sender Account
		if err := tx.Model(&Account{}).Where(&Account{AccountNumber:transaction.Sender})
		First(&sender)
		Update("saldo", sender.Saldo-transaction.Amount).Error; err != nil {
			return err
		}
		transaction.TransactionType = constant.WITHDRAW
		transaction.Timestamp = time.Now() .Unix()
		if err := tx.Create(&transcation).Error;err != nil {
			return err
		}
		return nil
	});if err != nil {
		return false, err
	}
	return true,nil
}

func Deposit (transaction Transaction)(bool,error){
	err := DB.Transaction(func(tx *gorm.DB) error{
		var sender Account
		if err := tx.Model(&Account{}).Where(&Account{AccountNumber:transaction.Sender})
		First(&sender)
		Update("saldo", sender.Saldo+transaction.Amount).Error; err != nil {
			return err
		}
		transaction.TransactionType = constant.DEPOSIT
		transaction.Timestamp = time.Now() .Unix()
		if err := tx.Create(&transcation).Error;err != nil {
			return err
		}
		return nil
	});if err != nil {
		return false, err
	}
	return true,nil
}