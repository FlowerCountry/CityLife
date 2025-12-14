#include "save/save_manager.h"
#include "disease/disease_manager.h"
#include "path/path_manager.h"
#include "world/world.h"
#include "user/user.h"
#include "bank/bank.h"

#include <fstream>
#include <sstream>
#include <ctime>
#include <sys/stat.h>
#include <cstdlib>

SaveManager *SaveManager::instance = nullptr;

SaveManager *SaveManager::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new SaveManager();
    }
    return instance;
}

SaveManager::SaveManager() {}

SaveManager::~SaveManager() {}

std::string SaveManager::GetSaveDirectory() const
{
    // 委托给 PathManager 获取统一的存档目录
    if (PathManager::GetInstance()->IsInitialized())
    {
        return PathManager::GetInstance()->GetSaveDirectory();
    }

    // Fallback: 使用用户目录（向后兼容）
    std::string home;
#ifdef _WIN32
    const char *userProfile = std::getenv("USERPROFILE");
    if (userProfile)
    {
        home = userProfile;
    }
    else
    {
        home = ".";
    }
    return home + "\\.citylife\\saves";
#else
    const char *homeEnv = std::getenv("HOME");
    if (homeEnv)
    {
        home = homeEnv;
    }
    else
    {
        home = ".";
    }
    return home + "/.citylife/saves";
#endif
}

std::string SaveManager::GetSavePath(int slot) const
{
    // 委托给 PathManager 获取存档路径
    if (PathManager::GetInstance()->IsInitialized())
    {
        return PathManager::GetInstance()->GetSavePath(slot);
    }

    // Fallback
    return GetSaveDirectory() + "/save_" + std::to_string(slot) + ".sav";
}

bool SaveManager::EnsureSaveDirectory() const
{
    // 如果 PathManager 已初始化，目录应该已经创建好了
    if (PathManager::GetInstance()->IsInitialized())
    {
        return true;
    }

    // Fallback: 手动创建目录（向后兼容）
    std::string dir = GetSaveDirectory();

#ifdef _WIN32
    // Windows: 逐级创建目录
    std::string parentDir = dir.substr(0, dir.rfind('\\'));
    struct stat st {};
    if (stat(parentDir.c_str(), &st) != 0)
    {
        system(("mkdir \"" + parentDir + "\"").c_str());
    }
    if (stat(dir.c_str(), &st) != 0)
    {
        return system(("mkdir \"" + dir + "\"").c_str()) == 0;
    }
#else
    // Linux/macOS
    std::string parentDir = dir.substr(0, dir.rfind('/'));
    struct stat st {};
    if (stat(parentDir.c_str(), &st) != 0)
    {
        mkdir(parentDir.c_str(), 0755);
    }
    if (stat(dir.c_str(), &st) != 0)
    {
        return mkdir(dir.c_str(), 0755) == 0;
    }
#endif
    return true;
}

std::string SaveManager::GetCurrentTimestamp() const
{
    time_t now = time(nullptr);
    struct tm *t = localtime(&now);
    char buf[64];
    strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", t);
    return std::string(buf);
}

std::string SaveManager::SerializeGameState() const
{
    World *world = World::GetInstance();
    User *user = User::GetInstance();

    std::ostringstream ss;
    ss << "# CityLife Save File\n";
    ss << "[meta]\n";
    ss << "version=1\n";
    ss << "save_time=" << GetCurrentTimestamp() << "\n\n";

    // 时间数据
    ss << "[world.time]\n";
    ss << "year=" << world->GetYear() << "\n";
    ss << "month=" << world->GetMonth() << "\n";
    ss << "day=" << world->GetDay() << "\n";
    ss << "hour=" << world->GetHour() << "\n";
    ss << "minute=" << world->GetMinute() << "\n";
    ss << "second=" << world->GetSecond() << "\n\n";

    // 世界状态
    ss << "[world.state]\n";
    ss << "where=" << world->GetWhere() << "\n";
    ss << "life_quality=" << world->GetLifeQuality() << "\n";
    ss << "bank_deposit=" << world->GetBankBalance() << "\n";

    const auto &wallet = world->GetWalletBreakdown();
    ss << "wallet=";
    for (std::size_t i = 0; i < wallet.size(); ++i)
    {
        if (i > 0)
            ss << ",";
        ss << wallet[i];
    }
    ss << "\n\n";

    // 用户健康数据
    ss << "[user.health]\n";
    const auto *health = user->GetUserHealth();
    for (const auto &pair : *health)
    {
        ss << pair.first << "=" << pair.second << "\n";
    }
    ss << "\n";

    // 疾病状态
    ss << "[diseases]\n";
    DiseaseManager *disease = DiseaseManager::GetInstance();
    ss << "data=" << disease->Serialize() << "\n";

    return ss.str();
}

bool SaveManager::DeserializeGameState(const std::string &content)
{
    World *world = World::GetInstance();
    User *user = User::GetInstance();
    DiseaseManager *disease = DiseaseManager::GetInstance();

    std::istringstream stream(content);
    std::string line;
    std::string currentSection;

    // 临时存储解析的值
    int year = 2010, month = 10, day = 10, hour = 11, minute = 3, second = 0;
    int where = 0;
    float lifeQuality = 1.0f;
    int bankDeposit = 0;
    std::array<int, 6> walletData = {1, 0, 0, 0, 0, 0};
    std::string diseaseData;

    while (std::getline(stream, line))
    {
        // 去除首尾空白
        size_t start = line.find_first_not_of(" \t\r\n");
        if (start == std::string::npos)
            continue;
        size_t end = line.find_last_not_of(" \t\r\n");
        line = line.substr(start, end - start + 1);

        // 跳过空行和注释
        if (line.empty() || line[0] == '#')
            continue;

        // 解析 section
        if (line[0] == '[')
        {
            size_t closePos = line.find(']');
            if (closePos != std::string::npos)
            {
                currentSection = line.substr(1, closePos - 1);
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

        // 根据 section 和 key 设置对应值
        if (currentSection == "world.time")
        {
            if (key == "year")
                year = std::stoi(value);
            else if (key == "month")
                month = std::stoi(value);
            else if (key == "day")
                day = std::stoi(value);
            else if (key == "hour")
                hour = std::stoi(value);
            else if (key == "minute")
                minute = std::stoi(value);
            else if (key == "second")
                second = std::stoi(value);
        }
        else if (currentSection == "world.state")
        {
            if (key == "where")
                where = std::stoi(value);
            else if (key == "life_quality")
                lifeQuality = std::stof(value);
            else if (key == "bank_deposit")
                bankDeposit = std::stoi(value);
            else if (key == "wallet")
            {
                // 解析逗号分隔数组
                std::istringstream walletStream(value);
                std::string num;
                std::size_t idx = 0;
                while (std::getline(walletStream, num, ',') && idx < 6)
                {
                    walletData[idx++] = std::stoi(num);
                }
            }
        }
        else if (currentSection == "user.health")
        {
            // 直接设置健康数据
            (*user->GetUserHealth())[key] = std::stoi(value);
        }
        else if (currentSection == "diseases")
        {
            if (key == "data")
            {
                diseaseData = value;
            }
        }
    }

    // 应用解析后的值到游戏状态
    world->SetTime(year, month, day, hour, minute, second);
    world->SetWhere(where);
    world->SetLifeQuality(lifeQuality);
    world->SetWallet(walletData);
    world->SetBankDeposit(bankDeposit);

    // 应用疾病状态
    if (!diseaseData.empty())
    {
        disease->Deserialize(diseaseData);
    }

    return true;
}

bool SaveManager::SaveGame(int slot)
{
    if (slot < 1 || slot > MAX_SLOTS)
    {
        return false;
    }

    if (!EnsureSaveDirectory())
    {
        return false;
    }

    std::string content = SerializeGameState();
    std::string path = GetSavePath(slot);

    std::ofstream file(path);
    if (!file.is_open())
    {
        return false;
    }

    file << content;
    file.close();
    return true;
}

bool SaveManager::LoadGame(int slot)
{
    if (slot < 1 || slot > MAX_SLOTS)
    {
        return false;
    }

    std::string path = GetSavePath(slot);
    std::ifstream file(path);
    if (!file.is_open())
    {
        return false;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();
    file.close();

    return DeserializeGameState(buffer.str());
}

bool SaveManager::HasSave(int slot) const
{
    if (slot < 1 || slot > MAX_SLOTS)
    {
        return false;
    }

    std::string path = GetSavePath(slot);
    struct stat st {};
    return stat(path.c_str(), &st) == 0;
}

bool SaveManager::DeleteSave(int slot)
{
    if (slot < 1 || slot > MAX_SLOTS)
    {
        return false;
    }

    std::string path = GetSavePath(slot);
    return std::remove(path.c_str()) == 0;
}

std::string SaveManager::ExtractSaveInfo(const std::string &filepath) const
{
    std::ifstream file(filepath);
    if (!file.is_open())
    {
        return "空";
    }

    std::string line;
    std::string saveTime;
    int day = 0;
    int money = 0;

    while (std::getline(file, line))
    {
        // 简单解析关键信息
        if (line.find("save_time=") == 0)
        {
            saveTime = line.substr(10, 10); // 只取日期部分
        }
        else if (line.find("day=") == 0)
        {
            day = std::stoi(line.substr(4));
        }
        else if (line.find("wallet=") == 0)
        {
            // 计算钱包总额
            std::string walletStr = line.substr(7);
            std::istringstream ws(walletStr);
            std::string num;
            int denominations[] = {100, 50, 20, 10, 5, 1};
            int idx = 0;
            while (std::getline(ws, num, ',') && idx < 6)
            {
                money += std::stoi(num) * denominations[idx++];
            }
        }
        else if (line.find("bank_deposit=") == 0)
        {
            money += std::stoi(line.substr(13));
        }
    }

    file.close();

    if (saveTime.empty())
    {
        return "空";
    }

    return saveTime + " - 第" + std::to_string(day) + "天 - ¥" + std::to_string(money);
}

std::string SaveManager::GetSaveDescription(int slot) const
{
    if (slot < 1 || slot > MAX_SLOTS)
    {
        return "无效槽位";
    }

    if (!HasSave(slot))
    {
        return "空";
    }

    return ExtractSaveInfo(GetSavePath(slot));
}

std::vector<std::string> SaveManager::GetAllSaveDescriptions() const
{
    std::vector<std::string> descriptions;
    for (int i = 1; i <= MAX_SLOTS; ++i)
    {
        descriptions.push_back(GetSaveDescription(i));
    }
    return descriptions;
}
