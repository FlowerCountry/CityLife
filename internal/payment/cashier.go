// Package payment 提供支付和钱包检查功能
package payment

// Denominations 面额常量（从大到小）
var Denominations = [6]int{100, 50, 20, 10, 5, 1}

// Result 支付结果
type Result struct {
	Success bool
	Used    [6]int // 各面额实际扣除张数
	Change  [6]int // 找零各面额张数
	Paid    int    // 实付总金额
}

// Cashier 收银员
type Cashier struct {
	denominations [6]int
}

// NewCashier 创建收银员
func NewCashier() *Cashier {
	return &Cashier{
		denominations: Denominations,
	}
}

// Pay 执行支付
// amount: 需要支付的金额
// wallet: 钱包中各面额的数量
// 返回支付结果，包括是否成功、扣款明细、找零明细和实付金额
func (c *Cashier) Pay(amount int, wallet [6]int) Result {
	result := Result{}

	if amount <= 0 {
		result.Success = true
		return result
	}

	// 计算钱包总额
	total := 0
	for i, count := range wallet {
		total += count * c.denominations[i]
	}
	if total < amount {
		return result // 余额不足
	}

	// 贪心算法：从大面值开始扣
	used := [6]int{}
	remaining := amount

	for i, denom := range c.denominations {
		if remaining <= 0 {
			break
		}
		take := min(wallet[i], remaining/denom)
		if take > 0 {
			used[i] = take
			remaining -= take * denom
		}
	}

	// 如果还有剩余，需要找一个能覆盖余额的面值（会产生找零）
	for remaining > 0 {
		bestIdx := -1
		bestDenom := 0

		// 找能覆盖剩余金额的最小面值
		for i, denom := range c.denominations {
			if wallet[i]-used[i] <= 0 {
				continue
			}
			if denom >= remaining && (bestIdx == -1 || denom < bestDenom) {
				bestIdx = i
				bestDenom = denom
			}
		}

		// 如果找不到能覆盖的，从最小面值开始凑
		if bestIdx == -1 {
			for i := len(c.denominations) - 1; i >= 0; i-- {
				if wallet[i]-used[i] > 0 {
					bestIdx = i
					bestDenom = c.denominations[i]
					break
				}
			}
		}

		if bestIdx == -1 {
			return result // 无法支付
		}

		used[bestIdx]++
		remaining -= bestDenom
	}

	// 计算实付金额
	paid := 0
	for i, count := range used {
		paid += count * c.denominations[i]
	}

	if paid < amount {
		return result // 不应该发生
	}

	// 计算找零
	changeValue := paid - amount
	change := [6]int{}
	for i, denom := range c.denominations {
		count := changeValue / denom
		change[i] = count
		changeValue -= count * denom
	}

	result.Success = true
	result.Used = used
	result.Change = change
	result.Paid = paid
	return result
}

// GetChangeTotal 计算找零总额
func (r *Result) GetChangeTotal() int {
	total := 0
	for i, count := range r.Change {
		total += count * Denominations[i]
	}
	return total
}

// min 返回两个整数中较小的一个
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
