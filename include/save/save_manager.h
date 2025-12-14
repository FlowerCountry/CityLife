#pragma once

#include <string>
#include <vector>

/**
 * @description: 存档管理器，负责游戏状态的持久化
 */
class SaveManager
{
  public:
    /**
     * @description: 获取单例实例
     * @return {*} SaveManager 单例指针
     */
    static SaveManager *GetInstance();

    /**
     * @description: 保存游戏到指定槽位
     * @param {int} slot 存档槽位 (1-3)
     * @return {*} 是否保存成功
     */
    bool SaveGame(int slot);

    /**
     * @description: 从指定槽位加载游戏
     * @param {int} slot 存档槽位 (1-3)
     * @return {*} 是否加载成功
     */
    bool LoadGame(int slot);

    /**
     * @description: 检查指定槽位是否有存档
     * @param {int} slot 存档槽位 (1-3)
     * @return {*} 是否存在存档
     */
    bool HasSave(int slot) const;

    /**
     * @description: 删除指定槽位的存档
     * @param {int} slot 存档槽位 (1-3)
     * @return {*} 是否删除成功
     */
    bool DeleteSave(int slot);

    /**
     * @description: 获取存档槽位的描述信息
     * @param {int} slot 存档槽位 (1-3)
     * @return {*} 描述字符串
     */
    std::string GetSaveDescription(int slot) const;

    /**
     * @description: 获取所有存档槽位的描述列表
     * @return {*} 描述列表
     */
    std::vector<std::string> GetAllSaveDescriptions() const;

    /**
     * @description: 获取存档目录路径
     * @return {*} 存档目录绝对路径
     */
    std::string GetSaveDirectory() const;

  private:
    SaveManager();
    ~SaveManager();

    static SaveManager *instance;
    static const int MAX_SLOTS = 3;

    /**
     * @description: 获取指定槽位的存档文件路径
     */
    std::string GetSavePath(int slot) const;

    /**
     * @description: 确保存档目录存在
     */
    bool EnsureSaveDirectory() const;

    /**
     * @description: 序列化当前游戏状态为 INI 格式字符串
     */
    std::string SerializeGameState() const;

    /**
     * @description: 从 INI 格式字符串反序列化到游戏状态
     */
    bool DeserializeGameState(const std::string &content);

    /**
     * @description: 获取当前时间戳字符串
     */
    std::string GetCurrentTimestamp() const;

    /**
     * @description: 从存档文件中提取元信息用于显示
     */
    std::string ExtractSaveInfo(const std::string &filepath) const;
};
