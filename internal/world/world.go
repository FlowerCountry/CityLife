// Package world 管理游戏世界状态：时间、金钱、位置和建筑
package world

import (
	"errors"
)

// Denominations 钱币面值：100, 50, 20, 10, 5, 1
var Denominations = [6]int{100, 50, 20, 10, 5, 1}

// ErrInsufficientFunds 资金不足错误
var ErrInsufficientFunds = errors.New("资金不足")

// ErrInvalidAmount 无效金额错误
var ErrInvalidAmount = errors.New("金额必须是100的倍数")

// World 游戏世界状态
type World struct {
	// 时间
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int

	// 位置
	Where     int         // 当前位置ID
	Buildings []*Building // 所有建筑

	// 金钱
	Wallet      [6]int  // 钱包：各面值数量 [100, 50, 20, 10, 5, 1]
	BankDeposit int     // 银行存款
	LifeQuality float64 // 生活质量
}

// New 创建新的游戏世界
func New() *World {
	return &World{
		Year:        2010,
		Month:       10,
		Day:         10,
		Hour:        11,
		Minute:      3,
		Second:      0,
		Where:       0,
		Wallet:      [6]int{1, 0, 0, 0, 0, 0}, // 初始 100 元
		BankDeposit: 0,
		LifeQuality: 1.0,
		Buildings:   initBuildings(),
	}
}

// GetWalletTotal 获取钱包总金额
func (w *World) GetWalletTotal() int {
	total := 0
	for i, count := range w.Wallet {
		total += Denominations[i] * count
	}
	return total
}

// AddToWallet 向钱包添加指定面值的钱币
func (w *World) AddToWallet(denomIndex, count int) {
	if denomIndex >= 0 && denomIndex < 6 {
		w.Wallet[denomIndex] += count
	}
}

// RemoveFromWallet 从钱包移除指定面值的钱币
func (w *World) RemoveFromWallet(denomIndex, count int) bool {
	if denomIndex >= 0 && denomIndex < 6 && w.Wallet[denomIndex] >= count {
		w.Wallet[denomIndex] -= count
		return true
	}
	return false
}

// SpendMoney 消费指定金额（使用支付算法）
func (w *World) SpendMoney(amount int) error {
	if w.GetWalletTotal() < amount {
		return ErrInsufficientFunds
	}

	// 使用贪心算法支付
	remaining := amount
	used := [6]int{}

	// 贪心阶段：从大面值开始
	for i, denom := range Denominations {
		if remaining <= 0 {
			break
		}
		take := min(w.Wallet[i], remaining/denom)
		used[i] = take
		remaining -= take * denom
	}

	// 如果还有余额，需要找零
	if remaining > 0 {
		// 找最小的能覆盖余额的面值
		for i := len(Denominations) - 1; i >= 0; i-- {
			if w.Wallet[i] > used[i] && Denominations[i] >= remaining {
				used[i]++
				remaining -= Denominations[i]
				break
			}
		}
	}

	// 计算实付金额和找零
	paid := 0
	for i, count := range used {
		paid += Denominations[i] * count
	}

	// 扣除使用的钱币
	for i, count := range used {
		w.Wallet[i] -= count
	}

	// 添加找零
	change := paid - amount
	if change > 0 {
		w.addChange(change)
	}

	return nil
}

// addChange 将找零添加回钱包
func (w *World) addChange(amount int) {
	remaining := amount
	for i, denom := range Denominations {
		count := remaining / denom
		if count > 0 {
			w.Wallet[i] += count
			remaining -= count * denom
		}
	}
}

// Deposit 存款（必须是100的倍数）
func (w *World) Deposit(amount int) error {
	if amount <= 0 || amount%100 != 0 {
		return ErrInvalidAmount
	}
	if w.GetWalletTotal() < amount {
		return ErrInsufficientFunds
	}

	// 从钱包扣除
	err := w.SpendMoney(amount)
	if err != nil {
		return err
	}

	// 添加到银行
	w.BankDeposit += amount
	return nil
}

// Withdraw 取款（必须是100的倍数）
func (w *World) Withdraw(amount int) error {
	if amount <= 0 || amount%100 != 0 {
		return ErrInvalidAmount
	}
	if w.BankDeposit < amount {
		return ErrInsufficientFunds
	}

	// 从银行扣除
	w.BankDeposit -= amount

	// 添加到钱包（以100元面值）
	w.Wallet[0] += amount / 100
	return nil
}

// UpdateTime 更新时间（秒）
func (w *World) UpdateTime(seconds int) {
	w.Second += seconds

	// 进位处理
	for w.Second >= 60 {
		w.Second -= 60
		w.Minute++
	}
	for w.Minute >= 60 {
		w.Minute -= 60
		w.Hour++
	}
	for w.Hour >= 24 {
		w.Hour -= 24
		w.Day++
	}

	// 月份天数处理
	daysInMonth := w.getDaysInMonth()
	for w.Day > daysInMonth {
		w.Day -= daysInMonth
		w.Month++
		if w.Month > 12 {
			w.Month = 1
			w.Year++
		}
		daysInMonth = w.getDaysInMonth()
	}
}

// getDaysInMonth 获取当前月份的天数
func (w *World) getDaysInMonth() int {
	switch w.Month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if w.isLeapYear() {
			return 29
		}
		return 28
	}
	return 30
}

// isLeapYear 判断是否闰年
func (w *World) isLeapYear() bool {
	return (w.Year%4 == 0 && w.Year%100 != 0) || w.Year%400 == 0
}

// CurrentLocation 获取当前位置的建筑
func (w *World) CurrentLocation() *Building {
	if w.Where >= 0 && w.Where < len(w.Buildings) {
		return w.Buildings[w.Where]
	}
	return nil
}

// ChangeLocation 改变位置
func (w *World) ChangeLocation(newWhere int) {
	if newWhere >= 0 && newWhere < len(w.Buildings) {
		w.Where = newWhere
	}
}

// GetTimeString 获取时间字符串
func (w *World) GetTimeString() string {
	return formatTime(w.Year, w.Month, w.Day, w.Hour, w.Minute, w.Second)
}

func formatTime(year, month, day, hour, minute, second int) string {
	return padInt(year, 4) + "年" + padInt(month, 2) + "月" + padInt(day, 2) + "日 " +
		padInt(hour, 2) + ":" + padInt(minute, 2) + ":" + padInt(second, 2)
}

func padInt(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
