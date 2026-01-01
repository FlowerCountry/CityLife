package payment

import (
	"citylife/internal/locale"
	"fmt"
)

// Summary 钱包摘要
type Summary struct {
	Lines   []string
	HasMore bool
}

// Detail 钱包详情
type Detail struct {
	Lines []string
}

// WalletInspector 钱包检查器
type WalletInspector struct {
	locale *locale.Locale
}

// NewWalletInspector 创建钱包检查器
func NewWalletInspector(loc *locale.Locale) *WalletInspector {
	return &WalletInspector{
		locale: loc,
	}
}

// BuildSummary 生成摘要（最多5项）
func (w *WalletInspector) BuildSummary(counts [6]int) Summary {
	summary := Summary{
		Lines:   make([]string, 0),
		HasMore: false,
	}

	if countNonZero(counts) == 0 {
		return summary
	}

	const summaryLimit = 5
	for i := range Denominations {
		count := counts[i]
		if count == 0 {
			continue
		}
		summary.Lines = append(summary.Lines, w.formatEntry(i, count))
		if len(summary.Lines) == summaryLimit {
			break
		}
	}

	summary.HasMore = countNonZero(counts) > len(summary.Lines)
	return summary
}

// BuildDetail 生成完整明细
func (w *WalletInspector) BuildDetail(counts [6]int) Detail {
	detail := Detail{
		Lines: make([]string, 0),
	}

	any := false
	for i := range Denominations {
		count := counts[i]
		if count == 0 {
			continue
		}
		any = true
		detail.Lines = append(detail.Lines, w.formatEntry(i, count))
	}

	if !any {
		detail.Lines = append(detail.Lines, w.locale.Get("wallet.misc", "no_bills"))
	}

	return detail
}

// PickZeroWalletMessage 返回"钱包空空"场景的提示语
func (w *WalletInspector) PickZeroWalletMessage() string {
	return w.locale.GetRandom("wallet.empty")
}

// PickSummaryLead 返回摘要部分的开场句
func (w *WalletInspector) PickSummaryLead() string {
	return w.locale.GetRandom("wallet.summary_lead")
}

// PickHasMoreMessage 返回"还有更多"的句子
func (w *WalletInspector) PickHasMoreMessage() string {
	return w.locale.GetRandom("wallet.has_more")
}

// PickCloseMessage 返回"收好钱包"的收尾句
func (w *WalletInspector) PickCloseMessage() string {
	return w.locale.GetRandom("wallet.close")
}

// PickDetailLead 返回"仔细点点"的开场句
func (w *WalletInspector) PickDetailLead() string {
	return w.locale.GetRandom("wallet.detail_lead")
}

// PickDetailEnd 返回"仔细点点"的收尾句
func (w *WalletInspector) PickDetailEnd() string {
	return w.locale.GetRandom("wallet.detail_end")
}

// formatEntry 格式化单项描述
func (w *WalletInspector) formatEntry(index, count int) string {
	denom := Denominations[index]
	desc := w.getDenominationDesc(denom)
	return fmt.Sprintf("%d 张%s%d 元", count, desc, denom)
}

// getDenominationDesc 获取面额描述
func (w *WalletInspector) getDenominationDesc(denom int) string {
	key := fmt.Sprintf("d%d", denom)
	desc := w.locale.Get("wallet.denomination", key)
	// 如果是fallback值，使用默认颜色
	if desc[0] == '[' {
		switch denom {
		case 100:
			return "鲜红的"
		case 50:
			return "紫色的"
		case 20:
			return "翠绿的"
		case 10:
			return "蓝灰的"
		case 5:
			return "褐色的"
		case 1:
			return "浅绿色的"
		default:
			return ""
		}
	}
	return desc
}

// countNonZero 统计非零面额的数量
func countNonZero(counts [6]int) int {
	total := 0
	for _, count := range counts {
		if count > 0 {
			total++
		}
	}
	return total
}

// FormatPaymentResult 格式化支付结果
func FormatPaymentResult(result Result, loc *locale.Locale) string {
	if !result.Success {
		return loc.Get("purchase.error", "insufficient")
	}

	msg := loc.Get("purchase.action", "success") + "\n"
	msg += loc.Get("purchase.action", "paid") + fmt.Sprintf("¥%d\n", result.Paid)

	changeTotal := result.GetChangeTotal()
	if changeTotal > 0 {
		msg += loc.Get("purchase.action", "change") + fmt.Sprintf("¥%d\n", changeTotal)
	}

	return msg
}

// FormatUsedDenominations 格式化扣款明细
func FormatUsedDenominations(used [6]int) string {
	msg := ""
	for i, count := range used {
		if count > 0 {
			msg += fmt.Sprintf("  %d×%d元", count, Denominations[i])
		}
	}
	return msg
}

// FormatChangeDenominations 格式化找零明细
func FormatChangeDenominations(change [6]int) string {
	msg := ""
	for i, count := range change {
		if count > 0 {
			msg += fmt.Sprintf("  %d×%d元", count, Denominations[i])
		}
	}
	return msg
}
