package account

import (
	"bufio"
	"fmt"
	"os"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
)

type PasswdCommand struct {
	DI      zdi.Invoker
	Conf    *service.Conf
	account *string
	passwd  *string
}

func (cmd *PasswdCommand) Flags(subcommand *zcli.Subcommand) {
	zlog.Log.SetLogLevel(zlog.LogWarn)
	cmd.account = zcli.SetVar("account", "Account").Required().String()
	cmd.passwd = zcli.SetVar("passwd", "Password").String()
	subcommand.Desc = "Modify user password"
}

func (cmd *PasswdCommand) Run(args []string) {
	_, _ = cmd.DI.Invoke(func() {
		salt := zstring.Rand(8)
		m, _ := parse.GetModel(UsersModel)
		passwd := *cmd.passwd
		if passwd == "" {
			f := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Please enter a new password:")
				passwd, _ = f.ReadString('\n')
				passwd = zstring.TrimSpace(passwd)
				if len(passwd) > 0 {
					break
				}
			}
		}

		filter := ztype.Map{
			"account": *cmd.account,
		}
		info, _ := parse.FindOne(m, filter)
		if info.IsEmpty() {
			zcli.Log.Error("User does not exist")
			return
		}
		_, err := parse.Update(m, filter, ztype.Map{
			"password": passwd,
			"salt":     salt,
		})
		if err != nil {
			zcli.Log.Error(err)
			return
		}
		zcli.Log.Success("Password modification succeeded")
	})
}
