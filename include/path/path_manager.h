// include/path/path_manager.h
// 统一路径管理器 - 管理应用程序所有数据路径

#ifndef PATH_MANAGER_H
#define PATH_MANAGER_H

#include <string>

/**
 * @class PathManager
 * @description: 统一管理应用程序的数据存储路径
 *
 * 目录结构:
 *   bin/
 *   ├── main              # 可执行文件
 *   └── data/             # 统一数据目录
 *       ├── locale/       # 语言资源
 *       └── saves/        # 存档文件
 */
class PathManager
{
  public:
    /**
     * @description: 获取单例实例
     * @return {PathManager*} 单例指针
     */
    static PathManager *GetInstance();

    /**
     * @description: 初始化路径管理器
     * @param {string} executableDir 可执行文件所在目录
     */
    void Initialize(const std::string &executableDir);

    /**
     * @description: 获取数据根目录
     * @return {string} 如 "/path/to/bin/data"
     */
    std::string GetDataDirectory() const;

    /**
     * @description: 获取存档目录
     * @return {string} 如 "/path/to/bin/data/saves"
     */
    std::string GetSaveDirectory() const;

    /**
     * @description: 获取语言资源目录
     * @return {string} 如 "/path/to/bin/data/locale"
     */
    std::string GetLocaleDirectory() const;

    /**
     * @description: 获取指定存档槽位的完整路径
     * @param {int} slot 存档槽位 (1-3)
     * @return {string} 如 "/path/to/bin/data/saves/save_1.sav"
     */
    std::string GetSavePath(int slot) const;

    /**
     * @description: 获取指定语言文件的完整路径
     * @param {string} locale 语言代码 (如 "zh_CN")
     * @return {string} 如 "/path/to/bin/data/locale/zh_CN.ini"
     */
    std::string GetLocalePath(const std::string &locale) const;

    /**
     * @description: 确保所有数据目录存在
     * @return {bool} 是否全部创建成功
     */
    bool EnsureDirectories() const;

    /**
     * @description: 检查是否已初始化
     * @return {bool} 是否已初始化
     */
    bool IsInitialized() const;

    /**
     * @description: 迁移旧数据到新位置
     * 检测 ~/.citylife/saves/ 和旧的 bin/locale/ 目录，
     * 如存在则自动迁移到新的 bin/data/ 结构下
     */
    void MigrateOldData() const;

    /**
     * @description: 获取可执行文件所在目录
     * @return {string} 基础路径
     */
    std::string GetBasePath() const;

  private:
    PathManager();
    ~PathManager();

    // 禁止拷贝和赋值
    PathManager(const PathManager &) = delete;
    PathManager &operator=(const PathManager &) = delete;

    static PathManager *instance;

    std::string basePath_;  // 可执行文件所在目录
    bool initialized_;

    /**
     * @description: 递归创建目录
     * @param {string} path 目录路径
     * @return {bool} 是否创建成功
     */
    bool CreateDirectoryRecursive(const std::string &path) const;

    /**
     * @description: 检查目录是否存在
     * @param {string} path 目录路径
     * @return {bool} 是否存在
     */
    bool DirectoryExists(const std::string &path) const;

    /**
     * @description: 检查文件是否存在
     * @param {string} path 文件路径
     * @return {bool} 是否存在
     */
    bool FileExists(const std::string &path) const;

    /**
     * @description: 复制文件
     * @param {string} src 源文件路径
     * @param {string} dst 目标文件路径
     * @return {bool} 是否复制成功
     */
    bool CopyFile(const std::string &src, const std::string &dst) const;

    /**
     * @description: 获取用户主目录
     * @return {string} 用户主目录路径
     */
    std::string GetHomeDirectory() const;
};

#endif // PATH_MANAGER_H
