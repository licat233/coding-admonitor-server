package logic

import "admonitor-server/service/model"

func AccountsFilter(accountList []*model.Account, channel_id int64, accountStatus int) (accounts []*model.Account) {
	for _, account := range accountList {
		if account.ChannelID == channel_id {
			if accountStatus == statusToCode(account.Status) || accountStatus == -1 {
				accounts = append(accounts, account)
			}
		}
	}
	return
}

func statusToCode(status bool) int {
	if status {
		return 1
	}
	return 0
}
