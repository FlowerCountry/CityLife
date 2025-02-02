/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-28 10:07:19
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-02-02 20:08:51
 * @FilePath: \CityLife\include\center\center.h
 */
#pragma once

#include <cstdlib>
#include <ctime>
#include <curses.h>
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
