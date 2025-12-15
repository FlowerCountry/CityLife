package locale

// loadDefaults 加载默认中文字符串
func (l *Locale) loadDefaults() {
	// ===== 营养警告提示 =====
	// 饥饿
	l.set("nutrition.warning.hunger", "level1", "你感到肚子有点空空的。")
	l.set("nutrition.warning.hunger", "level2", "你的胃在咕咕叫，急需进食。")
	l.set("nutrition.warning.hunger", "level3", "你感到头晕目眩，饥饿难耐。")
	l.set("nutrition.warning.hunger", "level4", "你的身体在颤抖，快要支撑不住了...")

	// 饥渴
	l.set("nutrition.warning.thirst", "level1", "你的嘴唇有些干燥。")
	l.set("nutrition.warning.thirst", "level2", "你感到口渴难耐。")
	l.set("nutrition.warning.thirst", "level3", "你的喉咙干涩，急需水分。")
	l.set("nutrition.warning.thirst", "level4", "你感到全身脱水，意识模糊...")

	// 蛋白质
	l.set("nutrition.warning.protein", "level1", "你感到有些乏力。")
	l.set("nutrition.warning.protein", "level2", "你的肌肉开始酸痛。")
	l.set("nutrition.warning.protein", "level3", "你的四肢沉重，每一步都很艰难。")
	l.set("nutrition.warning.protein", "level4", "你的身体在发出警报，肌肉在流失...")

	// 碳水化合物
	l.set("nutrition.warning.carbs", "level1", "你感到有点疲倦。")
	l.set("nutrition.warning.carbs", "level2", "你感到四肢无力。")
	l.set("nutrition.warning.carbs", "level3", "你的大脑一片混沌，无法集中注意力。")
	l.set("nutrition.warning.carbs", "level4", "你感到天旋地转，随时可能倒下...")

	// 幸福感
	l.set("nutrition.warning.happiness", "level1", "你感到有些闷闷不乐。")
	l.set("nutrition.warning.happiness", "level2", "一种空虚感笼罩着你。")
	l.set("nutrition.warning.happiness", "level3", "你感到深深的沮丧和无望。")
	l.set("nutrition.warning.happiness", "level4", "你觉得生活失去了意义...")

	// 精神振奋
	l.set("nutrition.warning.energy", "level1", "你感到有些昏昏欲睡。")
	l.set("nutrition.warning.energy", "level2", "你打了个哈欠，精神不济。")
	l.set("nutrition.warning.energy", "level3", "你的眼皮越来越沉重。")
	l.set("nutrition.warning.energy", "level4", "你几乎睁不开眼，意识在模糊...")

	// 维生素C
	l.set("nutrition.warning.vitaminc", "level1", "你的牙龈有些敏感。")
	l.set("nutrition.warning.vitaminc", "level2", "你感到容易疲劳，伤口愈合变慢。")

	// 铁
	l.set("nutrition.warning.iron", "level1", "你有时会感到头晕。")
	l.set("nutrition.warning.iron", "level2", "你的脸色有些苍白。")

	// 钙
	l.set("nutrition.warning.calcium", "level1", "你的骨头偶尔会咔咔作响。")
	l.set("nutrition.warning.calcium", "level2", "你的肌肉有时会抽筋。")

	// ===== 营养补充反馈 =====
	l.set("nutrition.feedback.hunger", "high", "你感到很饱足。")
	l.set("nutrition.feedback.hunger", "medium", "你不再感到饥饿。")
	l.set("nutrition.feedback.hunger", "low", "你的肚子不那么空了。")
	l.set("nutrition.feedback.hunger", "minimal", "你稍微垫了垫肚子。")

	l.set("nutrition.feedback.thirst", "high", "你感到神清气爽。")
	l.set("nutrition.feedback.thirst", "medium", "你不再感到口渴。")
	l.set("nutrition.feedback.thirst", "low", "你的喉咙湿润了些。")
	l.set("nutrition.feedback.thirst", "minimal", "你润了润嘴唇。")

	l.set("nutrition.feedback.protein", "high", "你感到精力充沛。")
	l.set("nutrition.feedback.protein", "medium", "你感到有些力气恢复了。")
	l.set("nutrition.feedback.protein", "low", "你补充了些能量。")

	l.set("nutrition.feedback.carbs", "high", "你感到活力满满。")
	l.set("nutrition.feedback.carbs", "medium", "你不再那么疲惫。")
	l.set("nutrition.feedback.carbs", "low", "你补充了些体力。")

	l.set("nutrition.feedback.happiness", "high", "一股幸福感涌上心头。")
	l.set("nutrition.feedback.happiness", "medium", "你的心情好了些。")
	l.set("nutrition.feedback.happiness", "low", "你感到一丝慰藉。")

	l.set("nutrition.feedback.energy", "high", "你感到精神焕发。")
	l.set("nutrition.feedback.energy", "medium", "你不再那么困倦。")
	l.set("nutrition.feedback.energy", "low", "你稍微提了提神。")

	l.set("nutrition.feedback.vitamin", "high", "你感到身体更健康了。")
	l.set("nutrition.feedback.vitamin", "low", "你补充了维生素。")

	l.set("nutrition.feedback.mineral", "high", "你感到身体更强壮了。")
	l.set("nutrition.feedback.mineral", "low", "你补充了矿物质。")

	l.set("nutrition.feedback.sugar", "high", "甜蜜的感觉让你开心。")
	l.set("nutrition.feedback.sugar", "low", "你尝到了些甜味。")

	l.set("nutrition.feedback.fiber", "high", "你感到肠胃舒适了。")
	l.set("nutrition.feedback.fiber", "low", "你补充了些纤维。")

	l.set("nutrition.feedback.probiotic", "high", "你感到肠道更健康了。")
	l.set("nutrition.feedback.probiotic", "low", "你补充了益生菌。")

	l.set("nutrition.feedback.default", "msg", "你补充了营养。")

	// ===== 钱包相关 =====
	l.set("wallet.empty", "msg1", "钱包空空的，一打开只剩下风吹过。")
	l.set("wallet.empty", "msg2", "口袋拍了拍，只剩点空气作伴。")
	l.set("wallet.empty", "msg3", "翻遍钱包，结果只有比心还干净的底。")

	l.set("wallet.summary_lead", "msg1", "随手一翻，看到这些钞票：")
	l.set("wallet.summary_lead", "msg2", "掀开钱包盖，里面排着：")
	l.set("wallet.summary_lead", "msg3", "哗啦啦翻看钱包，里面静静躺着：")

	l.set("wallet.has_more", "msg1", "剩下的零票得仔细点点才放心。")
	l.set("wallet.has_more", "msg2", "还有一些散票，想看的话不妨数一遍。")

	l.set("wallet.close", "msg1", "好吧，先把这些钞票塞回去。")
	l.set("wallet.close", "msg2", "收好钱包，别让风再偷走什么。")
	l.set("wallet.close", "msg3", "合上钱包，留着以后慢慢花。")

	l.set("wallet.detail_lead", "msg1", "仔细点点，所有面额如下：")
	l.set("wallet.detail_lead", "msg2", "认真清点一遍，钱包里其实是这样：")

	l.set("wallet.detail_end", "msg1", "点好啦，按回车把它们收回去。")
	l.set("wallet.detail_end", "msg2", "数得清清楚楚，回车收好钱包。")

	l.set("wallet.denomination", "d100", "鲜红的")
	l.set("wallet.denomination", "d50", "紫色的")
	l.set("wallet.denomination", "d20", "翠绿的")
	l.set("wallet.denomination", "d10", "蓝灰的")
	l.set("wallet.denomination", "d5", "褐色的")
	l.set("wallet.denomination", "d1", "浅绿色的")

	l.set("wallet.misc", "no_bills", "钱包里现在一张纸币都没有。")

	// ===== 疾病相关 =====
	l.set("disease.name", "cold", "感冒")
	l.set("disease.name", "anemia", "贫血")
	l.set("disease.name", "scurvy", "坏血病")
	l.set("disease.name", "food_poisoning", "食物中毒")
	l.set("disease.name", "malnutrition", "营养不良")
	l.set("disease.name", "depression", "抑郁症")

	l.set("disease.description", "cold", "维生素C长期不足导致免疫力下降")
	l.set("disease.description", "anemia", "铁元素长期缺乏导致血红蛋白不足")
	l.set("disease.description", "scurvy", "维生素C严重缺乏导致的血管疾病")
	l.set("disease.description", "food_poisoning", "食用变质食物引起的急性肠胃炎")
	l.set("disease.description", "malnutrition", "蛋白质长期缺乏导致身体虚弱")
	l.set("disease.description", "depression", "长期缺乏幸福感导致的心理疾病")

	l.set("disease.severity", "mild", "轻微")
	l.set("disease.severity", "moderate", "中等")
	l.set("disease.severity", "severe", "严重")
	l.set("disease.severity", "unknown", "未知")

	l.set("disease.status", "healthy", "健康")

	// ===== 食物新鲜度 =====
	l.set("food.freshness", "fresh", "新鲜")
	l.set("food.freshness", "fairly_fresh", "较新鲜")
	l.set("food.freshness", "not_fresh", "不太新鲜")
	l.set("food.freshness", "expiring", "快过期")
	l.set("food.freshness", "expired", "已过期")

	l.set("food.type", "fresh", "生鲜")
	l.set("food.type", "beverage", "饮料")
	l.set("food.type", "processed", "加工食品")
	l.set("food.type", "canned", "罐头")
	l.set("food.type", "unknown", "未知")

	// ===== 银行相关 =====
	l.set("bank.prompt", "deposit", "请输入你要存的钱:")
	l.set("bank.prompt", "withdraw", "请输入你要取的钱:")

	l.set("bank.error", "insufficient", "你没有这么多的钱")
	l.set("bank.error", "cannot_compose", "钱包里的面额拼不出这笔钱")
	l.set("bank.error", "deposit_hundreds", "银行柜员皱眉: 请按整百存款。")
	l.set("bank.error", "no_hundreds", "你手上没有足够的整百纸币")
	l.set("bank.error", "remove_failed", "无法从钱包中扣除整百纸币")
	l.set("bank.error", "withdraw_hundreds", "ATM 提醒: 只能按整百取款。")
	l.set("bank.error", "bank_insufficient", "银行账户余额不足")

	l.set("bank.display", "balance", "你的银行余额: ")
	l.set("bank.display", "wallet_balance", "你的账户余额: ¥")

	l.set("bank.deposit", "failed", "存入失败，请根据提示调整金额。")

	// ===== 购买相关 =====
	l.set("purchase.error", "insufficient", "钱包里的钱不够这件商品")
	l.set("purchase.error", "no_denom", "这一面额暂时掏不出来。")
	l.set("purchase.error", "no_selection", "还没选任何钞票。")
	l.set("purchase.error", "wrong_amount", "这组钞票凑不出目标金额。")
	l.set("purchase.error", "invalid_combo", "钱包里找不到这些面额组合。")

	l.set("purchase.action", "cleared", "已清空选择。")
	l.set("purchase.action", "found_more", "翻翻钱包，又摸出几张皱巴巴的钞票。")
	l.set("purchase.action", "no_more", "翻遍钱包也没多的了。")
	l.set("purchase.action", "cancelled", "你放弃了购买")
	l.set("purchase.action", "success", "购买成功")
	l.set("purchase.action", "paid", "实付 ")
	l.set("purchase.action", "change", "找零 ")

	l.set("purchase.option", "complete", "完成支付")
	l.set("purchase.option", "reselect", "重新选择")
	l.set("purchase.option", "search_wallet", "翻翻钱包")
	l.set("purchase.option", "cancel", "放弃购买")

	// ===== 医院相关 =====
	l.set("hospital.status", "healthy", "医生说：你很健康，不需要治疗。")
	l.set("hospital.status", "has_disease", "你有以下疾病需要治疗：")
	l.set("hospital.status", "choose_disease", "请选择要治疗的疾病：")

	l.set("hospital.option", "skip", "暂不治疗")
	l.set("hospital.option", "treatment_cost", " - 治疗费: ¥")

	l.set("hospital.result", "success", "治疗成功！你的")
	l.set("hospital.result", "cured", "已经痊愈。")
	l.set("hospital.result", "failed", "治疗失败，可能是资金不足。")

	// ===== 游戏状态 =====
	l.set("game.death", "random", "你的身体已经支撑不住了...")
	l.set("game.death", "malnutrition", "由于长期营养不良，你倒在了街头。")
	l.set("game.death", "starved", "你已经饿死了...")
	l.set("game.death", "game_over", "=== 游戏结束 ===")

	l.set("game.location", "current", "你当前位于: ")
	l.set("game.location", "goto", "前往")

	l.set("game.status", "body_signal", "[身体信号]")
	l.set("game.status", "disease", "[病] ")
	l.set("game.status", "danger", "[!危险] ")
	l.set("game.status", "warning", "[警告] ")
	l.set("game.status", "items_suffix", " 等")
	l.set("game.status", "items_count", "项")

	// ===== 存档相关 =====
	l.set("save.prompt", "choose_slot", "请选择槽位: ")

	l.set("save.result", "saved", "游戏已保存到槽位 ")
	l.set("save.result", "save_failed", "保存失败，请检查磁盘空间")
	l.set("save.result", "no_save", "该槽位没有存档")
	l.set("save.result", "loaded", "存档已加载")
	l.set("save.result", "load_failed", "加载失败，存档可能已损坏")

	// ===== 系统提示 =====
	l.set("system.input", "invalid", "输入无效，请输入一个整数: ")
	l.set("system.input", "choose", "请输入选项编号: ")
	l.set("system.input", "retry", "输入无效，请按回车重新输入。")

	l.set("system.action", "press_enter", "按回车返回")
	l.set("system.action", "check_count", "仔细点点")
	l.set("system.action", "close_wallet", "收好钱包")

	// ===== 营养等级状态 =====
	l.set("nutrition.level", "low", "偏低")
	l.set("nutrition.level", "danger", "危险")
	l.set("nutrition.level", "overdraft", "透支")
	l.set("nutrition.level", "critical", "濒死")
	l.set("nutrition.level", "abnormal", "异常")
}
