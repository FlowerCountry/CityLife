#pragma once

#include <string>
#include <map>
#include <vector>

/**
 * @description: 多语言字符串资源管理器
 * 负责从外部 INI 文件加载人性化表述字符串，支持多语言和按需加载
 */
class LocaleManager
{
public:
    /**
     * @description: 获取单例实例
     * @return {*} LocaleManager 单例指针
     */
    static LocaleManager *GetInstance();

    /**
     * @description: 初始化，加载指定语言的资源文件
     * @param {string} locale 语言代码（如 "zh_CN"）
     * @param {string} basePath 资源文件基础路径（默认为可执行文件所在目录）
     * @return {*} 是否初始化成功
     */
    bool Initialize(const std::string &locale = "zh_CN", const std::string &basePath = "");

    /**
     * @description: 获取指定 section 和 key 的字符串
     * @param {string} section 段名（如 "nutrition.warning.hunger"）
     * @param {string} key 键名（如 "level1"）
     * @return {*} 对应的字符串值，找不到时返回 "[section.key]" 作为 fallback
     */
    std::string Get(const std::string &section, const std::string &key);

    /**
     * @description: 获取整个 section 的所有 key-value 对
     * @param {string} section 段名
     * @return {*} 该 section 下所有键值对的 map
     */
    std::map<std::string, std::string> GetSection(const std::string &section);

    /**
     * @description: 获取 section 中所有值的列表（用于随机选择场景）
     * @param {string} section 段名
     * @return {*} 该 section 下所有值的向量
     */
    std::vector<std::string> GetSectionValues(const std::string &section);

    /**
     * @description: 检查指定 section 是否存在
     * @param {string} section 段名
     * @return {*} 是否存在
     */
    bool HasSection(const std::string &section) const;

    /**
     * @description: 检查指定 key 是否存在
     * @param {string} section 段名
     * @param {string} key 键名
     * @return {*} 是否存在
     */
    bool HasKey(const std::string &section, const std::string &key) const;

    /**
     * @description: 获取当前语言代码
     * @return {*} 当前语言代码
     */
    std::string GetCurrentLocale() const;

    /**
     * @description: 检查是否已初始化
     * @return {*} 是否已初始化
     */
    bool IsInitialized() const;

    /**
     * @description: 生成默认语言文件（如果不存在）
     * @param {string} locale 语言代码
     * @return {*} 是否成功
     */
    bool GenerateDefaultLocaleFile(const std::string &locale);

private:
    LocaleManager();
    ~LocaleManager();

    static LocaleManager *instance;

    /**
     * @description: 获取语言文件路径
     * @param {string} locale 语言代码
     * @return {*} 完整文件路径
     */
    std::string GetLocalePath(const std::string &locale) const;

    /**
     * @description: 确保 locale 目录存在
     * @return {*} 是否成功
     */
    bool EnsureLocaleDirectory() const;

    /**
     * @description: 解析 INI 文件内容
     * @param {string} content 文件内容
     * @return {*} 是否解析成功
     */
    bool ParseINI(const std::string &content);

    /**
     * @description: 生成默认中文字符串内容
     * @return {*} INI 格式的默认内容
     */
    std::string GenerateDefaultZhCN() const;

    // 存储结构：section -> (key -> value)
    std::map<std::string, std::map<std::string, std::string>> sections_;

    // 当前语言代码
    std::string currentLocale_;

    // 资源文件基础路径
    std::string basePath_;

    // 是否已初始化
    bool initialized_;
};
