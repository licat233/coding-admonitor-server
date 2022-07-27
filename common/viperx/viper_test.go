package viperx

import (
	"admonitor-server/common/utils"
	"encoding/json"
	"fmt"
	"path"
	"testing"

	"github.com/spf13/viper"
)

func Test_Account(t *testing.T) {
	LoadConfig()
	var n int
	if v, ok := viper.Get("popin.account-list").([]interface{}); ok {
		n = len(v)
	} else {
		return
	}
	res := []*AccountConfig{}
	for i := 0; i < n; i++ {
		accountId := viper.GetString(fmt.Sprintf("popin.account-list.%d.id", i))
		accountName := viper.GetString(fmt.Sprintf("popin.account-list.%d.name", i))
		nickname := viper.GetString(fmt.Sprintf("popin.account-list.%d.nickname", i))
		status := viper.GetBool(fmt.Sprintf("popin.account-list.%d.status", i))
		res = append(res, &AccountConfig{
			AccountID:   accountId,
			AccountName: accountName,
			Nickname:    nickname,
			Status:      status,
		})
	}
	for _, c := range res {
		fmt.Println(*c)
	}
}

type AccountConfig struct {
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Nickname    string `json:"nickname"`
	Status      bool   `json:"status"`
}

func Test_LoginInfo(t *testing.T) {
	LoadConfig()
	info := viper.GetStringMap("line.login-info")
	body, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(body))
	}
}

func Test_LoadLogdir(t *testing.T) {
	LoadConfig()
	var LogFilename = path.Join(utils.CurrentProjectRootPath(), viper.GetString("zap.log-dir"), viper.GetString("zap.log-file"))
	fmt.Println(LogFilename)
}
