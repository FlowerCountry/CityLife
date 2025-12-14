// src/path/path_manager.cpp
// 统一路径管理器实现

#include "path/path_manager.h"

#include <cstdlib>
#include <fstream>
#include <iostream>
#include <sys/stat.h>

#ifdef _WIN32
#include <direct.h>
#define PATH_SEPARATOR '\\'
#else
#include <dirent.h>
#include <unistd.h>
#define PATH_SEPARATOR '/'
#endif

PathManager *PathManager::instance = nullptr;

PathManager::PathManager() : initialized_(false)
{
}

PathManager::~PathManager()
{
}

PathManager *PathManager::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new PathManager();
    }
    return instance;
}

void PathManager::Initialize(const std::string &executableDir)
{
    basePath_ = executableDir;
    if (basePath_.empty())
    {
        basePath_ = ".";
    }
    initialized_ = true;
}

bool PathManager::IsInitialized() const
{
    return initialized_;
}

std::string PathManager::GetBasePath() const
{
    return basePath_;
}

std::string PathManager::GetDataDirectory() const
{
    return basePath_ + PATH_SEPARATOR + "data";
}

std::string PathManager::GetSaveDirectory() const
{
    return GetDataDirectory() + PATH_SEPARATOR + "saves";
}

std::string PathManager::GetLocaleDirectory() const
{
    return GetDataDirectory() + PATH_SEPARATOR + "locale";
}

std::string PathManager::GetSavePath(int slot) const
{
    return GetSaveDirectory() + PATH_SEPARATOR + "save_" + std::to_string(slot) + ".sav";
}

std::string PathManager::GetLocalePath(const std::string &locale) const
{
    return GetLocaleDirectory() + PATH_SEPARATOR + locale + ".ini";
}

bool PathManager::DirectoryExists(const std::string &path) const
{
    struct stat st
    {
    };
    return stat(path.c_str(), &st) == 0 && (st.st_mode & S_IFDIR);
}

bool PathManager::FileExists(const std::string &path) const
{
    struct stat st
    {
    };
    return stat(path.c_str(), &st) == 0 && (st.st_mode & S_IFREG);
}

bool PathManager::CreateDirectoryRecursive(const std::string &path) const
{
    // 如果目录已存在，直接返回成功
    if (DirectoryExists(path))
    {
        return true;
    }

    // 找到最后一个路径分隔符
    std::size_t pos = path.find_last_of("/\\");
    if (pos != std::string::npos && pos > 0)
    {
        // 递归创建父目录
        std::string parent = path.substr(0, pos);
        if (!CreateDirectoryRecursive(parent))
        {
            return false;
        }
    }

    // 创建当前目录
#ifdef _WIN32
    return _mkdir(path.c_str()) == 0 || DirectoryExists(path);
#else
    return mkdir(path.c_str(), 0755) == 0 || DirectoryExists(path);
#endif
}

bool PathManager::EnsureDirectories() const
{
    if (!initialized_)
    {
        return false;
    }

    // 创建数据根目录
    if (!CreateDirectoryRecursive(GetDataDirectory()))
    {
        std::cerr << "无法创建数据目录: " << GetDataDirectory() << std::endl;
        return false;
    }

    // 创建存档目录
    if (!CreateDirectoryRecursive(GetSaveDirectory()))
    {
        std::cerr << "无法创建存档目录: " << GetSaveDirectory() << std::endl;
        return false;
    }

    // 创建语言资源目录
    if (!CreateDirectoryRecursive(GetLocaleDirectory()))
    {
        std::cerr << "无法创建语言目录: " << GetLocaleDirectory() << std::endl;
        return false;
    }

    return true;
}

std::string PathManager::GetHomeDirectory() const
{
#ifdef _WIN32
    const char *userProfile = std::getenv("USERPROFILE");
    if (userProfile)
    {
        return std::string(userProfile);
    }
    return ".";
#else
    const char *homeEnv = std::getenv("HOME");
    if (homeEnv)
    {
        return std::string(homeEnv);
    }
    return ".";
#endif
}

bool PathManager::CopyFile(const std::string &src, const std::string &dst) const
{
    std::ifstream srcFile(src, std::ios::binary);
    if (!srcFile)
    {
        return false;
    }

    std::ofstream dstFile(dst, std::ios::binary);
    if (!dstFile)
    {
        return false;
    }

    dstFile << srcFile.rdbuf();
    return srcFile.good() && dstFile.good();
}

void PathManager::MigrateOldData() const
{
    if (!initialized_)
    {
        return;
    }

    // 1. 迁移旧的存档目录 (~/.citylife/saves/)
    std::string oldSaveDir = GetHomeDirectory() + PATH_SEPARATOR + ".citylife" + PATH_SEPARATOR + "saves";
    std::string newSaveDir = GetSaveDirectory();

    if (DirectoryExists(oldSaveDir))
    {
        std::cout << "检测到旧存档目录，正在迁移..." << std::endl;

        // 迁移每个存档槽位
        for (int slot = 1; slot <= 3; ++slot)
        {
            std::string oldSave = oldSaveDir + PATH_SEPARATOR + "save_" + std::to_string(slot) + ".sav";
            std::string newSave = GetSavePath(slot);

            if (FileExists(oldSave) && !FileExists(newSave))
            {
                if (CopyFile(oldSave, newSave))
                {
                    std::cout << "  迁移存档槽位 " << slot << " 成功" << std::endl;
                }
                else
                {
                    std::cerr << "  迁移存档槽位 " << slot << " 失败" << std::endl;
                }
            }
        }
    }

    // 2. 迁移旧的语言文件目录 (bin/locale/ -> bin/data/locale/)
    std::string oldLocaleDir = basePath_ + PATH_SEPARATOR + "locale";
    std::string newLocaleDir = GetLocaleDirectory();

    if (DirectoryExists(oldLocaleDir) && oldLocaleDir != newLocaleDir)
    {
        std::cout << "检测到旧语言目录，正在迁移..." << std::endl;

        // 迁移常用语言文件
        const char *locales[] = {"zh_CN", "en_US"};
        for (const char *locale : locales)
        {
            std::string oldFile = oldLocaleDir + PATH_SEPARATOR + std::string(locale) + ".ini";
            std::string newFile = GetLocalePath(locale);

            if (FileExists(oldFile) && !FileExists(newFile))
            {
                if (CopyFile(oldFile, newFile))
                {
                    std::cout << "  迁移语言文件 " << locale << ".ini 成功" << std::endl;
                }
                else
                {
                    std::cerr << "  迁移语言文件 " << locale << ".ini 失败" << std::endl;
                }
            }
        }
    }
}
