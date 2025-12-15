// Package action 包含银行相关行动
package action

import (
	"fmt"

	"citylife/internal/game"
)

// DepositAction 存款行动
type DepositAction struct{}

func (a *DepositAction) ID() string {
	return "deposit"
}

func (a *DepositAction) Info() string {
	return "存款"
}

func (a *DepositAction) Execute(state *game.State) *Result {
	walletTotal := state.World.GetWalletTotal()

	if walletTotal < 100 {
		return &Result{
			Message: fmt.Sprintf("存款失败！现金不足100元（当前: %d元）\n银行存取款必须以100元为单位", walletTotal),
			Success: false,
		}
	}

	// 存款（默认存100元）
	amount := 100

	err := state.World.Deposit(amount)
	if err != nil {
		return &Result{
			Message: "存款失败！" + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("存款成功！已存入 %d 元\n当前银行余额: %d 元\n当前现金: %d 元",
			amount, state.World.BankDeposit, state.World.GetWalletTotal()),
		Success: true,
	}
}

func (a *DepositAction) Category() EventCategory {
	return CategoryPrimary
}

// WithdrawAction 取款行动
type WithdrawAction struct{}

func (a *WithdrawAction) ID() string {
	return "withdraw"
}

func (a *WithdrawAction) Info() string {
	return "取款"
}

func (a *WithdrawAction) Execute(state *game.State) *Result {
	bankBalance := state.World.BankDeposit

	if bankBalance < 100 {
		return &Result{
			Message: fmt.Sprintf("取款失败！银行余额不足100元（当前: %d元）\n银行存取款必须以100元为单位", bankBalance),
			Success: false,
		}
	}

	// 取款100元
	amount := 100

	err := state.World.Withdraw(amount)
	if err != nil {
		return &Result{
			Message: "取款失败！" + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("取款成功！已取出 %d 元\n当前银行余额: %d 元\n当前现金: %d 元",
			amount, state.World.BankDeposit, state.World.GetWalletTotal()),
		Success: true,
	}
}

func (a *WithdrawAction) Category() EventCategory {
	return CategoryPrimary
}

// DepositAllAction 存入全部（100的倍数）
type DepositAllAction struct{}

func (a *DepositAllAction) ID() string {
	return "deposit_all"
}

func (a *DepositAllAction) Info() string {
	return "存入全部现金"
}

func (a *DepositAllAction) Execute(state *game.State) *Result {
	walletTotal := state.World.GetWalletTotal()

	if walletTotal < 100 {
		return &Result{
			Message: fmt.Sprintf("存款失败！现金不足100元（当前: %d元）", walletTotal),
			Success: false,
		}
	}

	// 存入所有可存款金额（100的倍数）
	amount := (walletTotal / 100) * 100

	err := state.World.Deposit(amount)
	if err != nil {
		return &Result{
			Message: "存款失败！" + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("存款成功！已存入 %d 元\n当前银行余额: %d 元\n当前现金: %d 元",
			amount, state.World.BankDeposit, state.World.GetWalletTotal()),
		Success: true,
	}
}

func (a *DepositAllAction) Category() EventCategory {
	return CategoryPrimary
}

// WithdrawAllAction 取出全部
type WithdrawAllAction struct{}

func (a *WithdrawAllAction) ID() string {
	return "withdraw_all"
}

func (a *WithdrawAllAction) Info() string {
	return "取出全部存款"
}

func (a *WithdrawAllAction) Execute(state *game.State) *Result {
	bankBalance := state.World.BankDeposit

	if bankBalance < 100 {
		return &Result{
			Message: fmt.Sprintf("取款失败！银行余额不足100元（当前: %d元）", bankBalance),
			Success: false,
		}
	}

	// 取出所有可取款金额（100的倍数）
	amount := (bankBalance / 100) * 100

	err := state.World.Withdraw(amount)
	if err != nil {
		return &Result{
			Message: "取款失败！" + err.Error(),
			Success: false,
		}
	}

	return &Result{
		Message: fmt.Sprintf("取款成功！已取出 %d 元\n当前银行余额: %d 元\n当前现金: %d 元",
			amount, state.World.BankDeposit, state.World.GetWalletTotal()),
		Success: true,
	}
}

func (a *WithdrawAllAction) Category() EventCategory {
	return CategoryPrimary
}
