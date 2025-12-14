#pragma once

#include <cstdlib>
#include <ctime>
#include <string>
#include <vector>

class World;
class View;
class Center {
  public:
    /**
     * @description: 市中心构造函数
     * @return null
     */
    Center();
    /**
     * @description: 市中心析构函数
     * @return null
     */
    ~Center();
    /**
     * @description: 返回市中心的单例对象
     * @return {*} 市中心的单例对象
     */
    static Center *GetInstance();
    /**
     * @description: 显示公告
     * @return null
     */
    void PrintAnnouncement();
    /**
     * @description: 随机插入元素至数组
     * @param {vector<std::string>} &vec 目标数组
     * @param {string} &str 元素
     * @return null
     */
    void InsertStringRandomly(std::vector<std::string> &vec, const std::string &str);

  private:
    static Center *instance;
    std::vector<std::string> announcement;
    std::vector<int> RandArray;
};
