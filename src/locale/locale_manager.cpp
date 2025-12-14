#include "locale/locale_manager.h"
#include "path/path_manager.h"

#include <fstream>
#include <sstream>
#include <sys/stat.h>
#include <cstdlib>

LocaleManager *LocaleManager::instance = nullptr;

LocaleManager *LocaleManager::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new LocaleManager();
    }
    return instance;
}

LocaleManager::LocaleManager()
    : currentLocale_("zh_CN"),
      basePath_(""),
      initialized_(false)
{
}

LocaleManager::~LocaleManager()
{
}

bool LocaleManager::Initialize(const std::string &locale, const std::string &basePath)
{
    currentLocale_ = locale;
    basePath_ = basePath;

    // 如果 basePath 为空，使用可执行文件所在目录
    if (basePath_.empty())
    {
        basePath_ = ".";
    }

    // 确保目录存在
    EnsureLocaleDirectory();

    // 获取语言文件路径
    std::string filepath = GetLocalePath(locale);

    // 尝试读取文件
    std::ifstream file(filepath);
    if (!file.is_open())
    {
        // 文件不存在，生成默认文件
        if (!GenerateDefaultLocaleFile(locale))
        {
            return false;
        }
        // 重新尝试读取
        file.open(filepath);
        if (!file.is_open())
        {
            return false;
        }
    }

    // 读取文件内容
    std::stringstream buffer;
    buffer << file.rdbuf();
    file.close();

    // 解析 INI 内容
    if (!ParseINI(buffer.str()))
    {
        return false;
    }

    initialized_ = true;
    return true;
}

std::string LocaleManager::GetLocalePath(const std::string &locale) const
{
    // 优先使用 PathManager
    if (PathManager::GetInstance()->IsInitialized())
    {
        return PathManager::GetInstance()->GetLocalePath(locale);
    }

    // Fallback: 使用旧路径逻辑
    return basePath_ + "/locale/" + locale + ".ini";
}

bool LocaleManager::EnsureLocaleDirectory() const
{
    // 如果 PathManager 已初始化，目录应该已经创建好了
    if (PathManager::GetInstance()->IsInitialized())
    {
        return true;
    }

    // Fallback: 手动创建目录
    std::string dir = basePath_ + "/locale";

    struct stat st {};
    if (stat(dir.c_str(), &st) != 0)
    {
#ifdef _WIN32
        return system(("mkdir \"" + dir + "\"").c_str()) == 0;
#else
        return mkdir(dir.c_str(), 0755) == 0;
#endif
    }
    return true;
}

bool LocaleManager::ParseINI(const std::string &content)
{
    sections_.clear();

    std::istringstream stream(content);
    std::string line;
    std::string currentSection;

    while (std::getline(stream, line))
    {
        // 去除首尾空白
        size_t start = line.find_first_not_of(" \t\r\n");
        if (start == std::string::npos)
            continue;
        size_t end = line.find_last_not_of(" \t\r\n");
        line = line.substr(start, end - start + 1);

        // 跳过空行和注释
        if (line.empty() || line[0] == '#' || line[0] == ';')
            continue;

        // 解析 section
        if (line[0] == '[')
        {
            size_t closePos = line.find(']');
            if (closePos != std::string::npos)
            {
                currentSection = line.substr(1, closePos - 1);
                // 确保 section 存在
                if (sections_.find(currentSection) == sections_.end())
                {
                    sections_[currentSection] = std::map<std::string, std::string>();
                }
            }
            continue;
        }

        // 解析 key=value
        size_t eqPos = line.find('=');
        if (eqPos == std::string::npos)
            continue;

        std::string key = line.substr(0, eqPos);
        std::string value = line.substr(eqPos + 1);

        // 去除 key 和 value 的首尾空白
        start = key.find_first_not_of(" \t");
        if (start != std::string::npos)
        {
            end = key.find_last_not_of(" \t");
            key = key.substr(start, end - start + 1);
        }
        start = value.find_first_not_of(" \t");
        if (start != std::string::npos)
        {
            end = value.find_last_not_of(" \t");
            value = value.substr(start, end - start + 1);
        }

        // 存储到当前 section
        if (!currentSection.empty())
        {
            sections_[currentSection][key] = value;
        }
    }

    return true;
}

std::string LocaleManager::Get(const std::string &section, const std::string &key)
{
    auto sectionIt = sections_.find(section);
    if (sectionIt == sections_.end())
    {
        // Section 不存在，返回 fallback
        return "[" + section + "." + key + "]";
    }

    auto keyIt = sectionIt->second.find(key);
    if (keyIt == sectionIt->second.end())
    {
        // Key 不存在，返回 fallback
        return "[" + section + "." + key + "]";
    }

    return keyIt->second;
}

std::map<std::string, std::string> LocaleManager::GetSection(const std::string &section)
{
    auto it = sections_.find(section);
    if (it != sections_.end())
    {
        return it->second;
    }
    return std::map<std::string, std::string>();
}

std::vector<std::string> LocaleManager::GetSectionValues(const std::string &section)
{
    std::vector<std::string> values;
    auto it = sections_.find(section);
    if (it != sections_.end())
    {
        for (const auto &pair : it->second)
        {
            values.push_back(pair.second);
        }
    }
    return values;
}

bool LocaleManager::HasSection(const std::string &section) const
{
    return sections_.find(section) != sections_.end();
}

bool LocaleManager::HasKey(const std::string &section, const std::string &key) const
{
    auto sectionIt = sections_.find(section);
    if (sectionIt == sections_.end())
    {
        return false;
    }
    return sectionIt->second.find(key) != sectionIt->second.end();
}

std::string LocaleManager::GetCurrentLocale() const
{
    return currentLocale_;
}

bool LocaleManager::IsInitialized() const
{
    return initialized_;
}

bool LocaleManager::GenerateDefaultLocaleFile(const std::string &locale)
{
    std::string content;
    if (locale == "zh_CN")
    {
        content = GenerateDefaultZhCN();
    }
    else
    {
        // 其他语言暂时使用中文作为模板
        content = GenerateDefaultZhCN();
    }

    std::string filepath = GetLocalePath(locale);
    std::ofstream file(filepath);
    if (!file.is_open())
    {
        return false;
    }

    file << content;
    file.close();
    return true;
}

std::string LocaleManager::GenerateDefaultZhCN() const
{
    std::ostringstream ss;

    ss << "# CityLife 语言文件 - 简体中文\n";
    ss << "# 版本: 1.0\n\n";

    // ===== 营养警告提示 =====
    ss << "# 营养状态警告（等级1-4）\n";

    ss << "[nutrition.warning.hunger]\n";
    ss << "level1=你感到肚子有点空空的。\n";
    ss << "level2=你的胃在咕咕叫，急需进食。\n";
    ss << "level3=你感到头晕目眩，饥饿难耐。\n";
    ss << "level4=你的身体在颤抖，快要支撑不住了...\n\n";

    ss << "[nutrition.warning.thirst]\n";
    ss << "level1=你的嘴唇有些干燥。\n";
    ss << "level2=你感到口渴难耐。\n";
    ss << "level3=你的喉咙干涩，急需水分。\n";
    ss << "level4=你感到全身脱水，意识模糊...\n\n";

    ss << "[nutrition.warning.protein]\n";
    ss << "level1=你感到有些乏力。\n";
    ss << "level2=你的肌肉开始酸痛。\n";
    ss << "level3=你的四肢沉重，每一步都很艰难。\n";
    ss << "level4=你的身体在发出警报，肌肉在流失...\n\n";

    ss << "[nutrition.warning.carbs]\n";
    ss << "level1=你感到有点疲倦。\n";
    ss << "level2=你感到四肢无力。\n";
    ss << "level3=你的大脑一片混沌，无法集中注意力。\n";
    ss << "level4=你感到天旋地转，随时可能倒下...\n\n";

    ss << "[nutrition.warning.happiness]\n";
    ss << "level1=你感到有些闷闷不乐。\n";
    ss << "level2=一种空虚感笼罩着你。\n";
    ss << "level3=你感到深深的沮丧和无望。\n";
    ss << "level4=你觉得生活失去了意义...\n\n";

    ss << "[nutrition.warning.energy]\n";
    ss << "level1=你感到有些昏昏欲睡。\n";
    ss << "level2=你打了个哈欠，精神不济。\n";
    ss << "level3=你的眼皮越来越沉重。\n";
    ss << "level4=你几乎睁不开眼，意识在模糊...\n\n";

    ss << "[nutrition.warning.vitaminc]\n";
    ss << "level1=你的牙龈有些敏感。\n";
    ss << "level2=你感到容易疲劳，伤口愈合变慢。\n\n";

    ss << "[nutrition.warning.iron]\n";
    ss << "level1=你有时会感到头晕。\n";
    ss << "level2=你的脸色有些苍白。\n\n";

    ss << "[nutrition.warning.calcium]\n";
    ss << "level1=你的骨头偶尔会咔咔作响。\n";
    ss << "level2=你的肌肉有时会抽筋。\n\n";

    // ===== 营养补充反馈 =====
    ss << "# 营养补充后的感受\n";

    ss << "[nutrition.feedback.hunger]\n";
    ss << "high=你感到很饱足。\n";
    ss << "medium=你不再感到饥饿。\n";
    ss << "low=你的肚子不那么空了。\n";
    ss << "minimal=你稍微垫了垫肚子。\n\n";

    ss << "[nutrition.feedback.thirst]\n";
    ss << "high=你感到神清气爽。\n";
    ss << "medium=你不再感到口渴。\n";
    ss << "low=你的喉咙湿润了些。\n";
    ss << "minimal=你润了润嘴唇。\n\n";

    ss << "[nutrition.feedback.protein]\n";
    ss << "high=你感到精力充沛。\n";
    ss << "medium=你感到有些力气恢复了。\n";
    ss << "low=你补充了些能量。\n\n";

    ss << "[nutrition.feedback.carbs]\n";
    ss << "high=你感到活力满满。\n";
    ss << "medium=你不再那么疲惫。\n";
    ss << "low=你补充了些体力。\n\n";

    ss << "[nutrition.feedback.happiness]\n";
    ss << "high=一股幸福感涌上心头。\n";
    ss << "medium=你的心情好了些。\n";
    ss << "low=你感到一丝慰藉。\n\n";

    ss << "[nutrition.feedback.energy]\n";
    ss << "high=你感到精神焕发。\n";
    ss << "medium=你不再那么困倦。\n";
    ss << "low=你稍微提了提神。\n\n";

    ss << "[nutrition.feedback.vitamin]\n";
    ss << "high=你感到身体更健康了。\n";
    ss << "low=你补充了维生素。\n\n";

    ss << "[nutrition.feedback.mineral]\n";
    ss << "high=你感到身体更强壮了。\n";
    ss << "low=你补充了矿物质。\n\n";

    ss << "[nutrition.feedback.sugar]\n";
    ss << "high=甜蜜的感觉让你开心。\n";
    ss << "low=你尝到了些甜味。\n\n";

    ss << "[nutrition.feedback.fiber]\n";
    ss << "high=你感到肠胃舒适了。\n";
    ss << "low=你补充了些纤维。\n\n";

    ss << "[nutrition.feedback.probiotic]\n";
    ss << "high=你感到肠道更健康了。\n";
    ss << "low=你补充了益生菌。\n\n";

    ss << "[nutrition.feedback.default]\n";
    ss << "msg=你补充了营养。\n\n";

    // ===== 钱包相关 =====
    ss << "# 钱包查看相关提示\n";

    ss << "[wallet.empty]\n";
    ss << "msg1=钱包空空的，一打开只剩下风吹过。\n";
    ss << "msg2=口袋拍了拍，只剩点空气作伴。\n";
    ss << "msg3=翻遍钱包，结果只有比心还干净的底。\n\n";

    ss << "[wallet.summary_lead]\n";
    ss << "msg1=随手一翻，看到这些钞票：\n";
    ss << "msg2=掀开钱包盖，里面排着：\n";
    ss << "msg3=哗啦啦翻看钱包，里面静静躺着：\n\n";

    ss << "[wallet.has_more]\n";
    ss << "msg1=剩下的零票得仔细点点才放心。\n";
    ss << "msg2=还有一些散票，想看的话不妨数一遍。\n\n";

    ss << "[wallet.close]\n";
    ss << "msg1=好吧，先把这些钞票塞回去。\n";
    ss << "msg2=收好钱包，别让风再偷走什么。\n";
    ss << "msg3=合上钱包，留着以后慢慢花。\n\n";

    ss << "[wallet.detail_lead]\n";
    ss << "msg1=仔细点点，所有面额如下：\n";
    ss << "msg2=认真清点一遍，钱包里其实是这样：\n\n";

    ss << "[wallet.detail_end]\n";
    ss << "msg1=点好啦，按回车把它们收回去。\n";
    ss << "msg2=数得清清楚楚，回车收好钱包。\n\n";

    ss << "[wallet.denomination]\n";
    ss << "d100=鲜红的\n";
    ss << "d50=紫色的\n";
    ss << "d20=翠绿的\n";
    ss << "d10=蓝灰的\n";
    ss << "d5=褐色的\n";
    ss << "d1=浅绿色的\n\n";

    ss << "[wallet.misc]\n";
    ss << "no_bills=钱包里现在一张纸币都没有。\n\n";

    // ===== 疾病相关 =====
    ss << "# 疾病相关\n";

    ss << "[disease.name]\n";
    ss << "cold=感冒\n";
    ss << "anemia=贫血\n";
    ss << "scurvy=坏血病\n";
    ss << "food_poisoning=食物中毒\n";
    ss << "malnutrition=营养不良\n";
    ss << "depression=抑郁症\n\n";

    ss << "[disease.description]\n";
    ss << "cold=维生素C长期不足导致免疫力下降\n";
    ss << "anemia=铁元素长期缺乏导致血红蛋白不足\n";
    ss << "scurvy=维生素C严重缺乏导致的血管疾病\n";
    ss << "food_poisoning=食用变质食物引起的急性肠胃炎\n";
    ss << "malnutrition=蛋白质长期缺乏导致身体虚弱\n";
    ss << "depression=长期缺乏幸福感导致的心理疾病\n\n";

    ss << "[disease.severity]\n";
    ss << "mild=轻微\n";
    ss << "moderate=中等\n";
    ss << "severe=严重\n";
    ss << "unknown=未知\n\n";

    ss << "[disease.status]\n";
    ss << "healthy=健康\n\n";

    // ===== 食物新鲜度 =====
    ss << "# 食物新鲜度\n";

    ss << "[food.freshness]\n";
    ss << "fresh=新鲜\n";
    ss << "fairly_fresh=较新鲜\n";
    ss << "not_fresh=不太新鲜\n";
    ss << "expiring=快过期\n";
    ss << "expired=已过期\n\n";

    ss << "[food.type]\n";
    ss << "fresh=生鲜\n";
    ss << "beverage=饮料\n";
    ss << "processed=加工食品\n";
    ss << "canned=罐头\n";
    ss << "unknown=未知\n\n";

    // ===== 银行相关 =====
    ss << "# 银行相关\n";

    ss << "[bank.prompt]\n";
    ss << "deposit=请输入你要存的钱:\n";
    ss << "withdraw=请输入你要取的钱:\n\n";

    ss << "[bank.error]\n";
    ss << "insufficient=你没有这么多的钱\n";
    ss << "cannot_compose=钱包里的面额拼不出这笔钱\n";
    ss << "deposit_hundreds=银行柜员皱眉: 请按整百存款。\n";
    ss << "no_hundreds=你手上没有足够的整百纸币\n";
    ss << "remove_failed=无法从钱包中扣除整百纸币\n";
    ss << "withdraw_hundreds=ATM 提醒: 只能按整百取款。\n";
    ss << "bank_insufficient=银行账户余额不足\n\n";

    ss << "[bank.display]\n";
    ss << "balance=你的银行余额: \n";
    ss << "wallet_balance=你的账户余额: ¥\n\n";

    ss << "[bank.deposit]\n";
    ss << "failed=存入失败，请根据提示调整金额。\n\n";

    // ===== 购买相关 =====
    ss << "# 购买相关\n";

    ss << "[purchase.error]\n";
    ss << "insufficient=钱包里的钱不够这件商品\n";
    ss << "no_denom=这一面额暂时掏不出来。\n";
    ss << "no_selection=还没选任何钞票。\n";
    ss << "wrong_amount=这组钞票凑不出目标金额。\n";
    ss << "invalid_combo=钱包里找不到这些面额组合。\n\n";

    ss << "[purchase.action]\n";
    ss << "cleared=已清空选择。\n";
    ss << "found_more=翻翻钱包，又摸出几张皱巴巴的钞票。\n";
    ss << "no_more=翻遍钱包也没多的了。\n";
    ss << "cancelled=你放弃了购买\n";
    ss << "success=购买成功\n";
    ss << "paid=实付 \n";
    ss << "change=找零 \n\n";

    ss << "[purchase.option]\n";
    ss << "complete=完成支付\n";
    ss << "reselect=重新选择\n";
    ss << "search_wallet=翻翻钱包\n";
    ss << "cancel=放弃购买\n\n";

    // ===== 医院相关 =====
    ss << "# 医院相关\n";

    ss << "[hospital.status]\n";
    ss << "healthy=医生说：你很健康，不需要治疗。\n";
    ss << "has_disease=你有以下疾病需要治疗：\n";
    ss << "choose_disease=请选择要治疗的疾病：\n\n";

    ss << "[hospital.option]\n";
    ss << "skip=暂不治疗\n";
    ss << "treatment_cost= - 治疗费: ¥\n\n";

    ss << "[hospital.result]\n";
    ss << "success=治疗成功！你的\n";
    ss << "cured=已经痊愈。\n";
    ss << "failed=治疗失败，可能是资金不足。\n\n";

    // ===== 游戏状态 =====
    ss << "# 游戏状态\n";

    ss << "[game.death]\n";
    ss << "random=你的身体已经支撑不住了...\n";
    ss << "malnutrition=由于长期营养不良，你倒在了街头。\n";
    ss << "starved=你已经饿死了...\n";
    ss << "game_over==== 游戏结束 ===\n\n";

    ss << "[game.location]\n";
    ss << "current=你当前位于: \n";
    ss << "goto=前往\n\n";

    ss << "[game.status]\n";
    ss << "body_signal=[身体信号]\n";
    ss << "disease=[病] \n";
    ss << "danger=[!危险] \n";
    ss << "warning=[警告] \n";
    ss << "items_suffix= 等\n";
    ss << "items_count=项\n\n";

    // ===== 存档相关 =====
    ss << "# 存档相关\n";

    ss << "[save.prompt]\n";
    ss << "choose_slot=请选择槽位: \n\n";

    ss << "[save.result]\n";
    ss << "saved=游戏已保存到槽位 \n";
    ss << "save_failed=保存失败，请检查磁盘空间\n";
    ss << "no_save=该槽位没有存档\n";
    ss << "loaded=存档已加载\n";
    ss << "load_failed=加载失败，存档可能已损坏\n\n";

    // ===== 系统提示 =====
    ss << "# 系统提示\n";

    ss << "[system.input]\n";
    ss << "invalid=输入无效，请输入一个整数: \n";
    ss << "choose=请输入选项编号: \n";
    ss << "retry=输入无效，请按回车重新输入。\n\n";

    ss << "[system.action]\n";
    ss << "press_enter=按回车返回\n";
    ss << "check_count=仔细点点\n";
    ss << "close_wallet=收好钱包\n\n";

    // ===== 营养等级状态 =====
    ss << "# 营养等级状态后缀\n";

    ss << "[nutrition.level]\n";
    ss << "low=偏低\n";
    ss << "danger=危险\n";
    ss << "overdraft=透支\n";
    ss << "critical=濒死\n";
    ss << "abnormal=异常\n\n";

    return ss.str();
}
