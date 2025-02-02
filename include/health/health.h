/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-30 07:03:53
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 15:25:47
 * @FilePath: \CityLife\include\health\health.h
 */
#pragma once

#include <string>

class Health {
  public:
    /**
     * @description: 营养的构造函数
     * @param {string} name 营养的名称
     * @param {int} reserves 对应的数值, 百分比
     * @return null
     */
    Health(const std::string &name, const int &reserves) : name{name}, reserves{reserves} {}
    /**
     * @description: 返回营养名称
     * @return {*} 营养名称
     */
    std::string GetInfo() { return name; }
    /**
     * @description: 返回营养数值
     * @return {*} 营养数值, 百分比
     */
    int GetReserves() { return reserves; }

  private:
    std::string name;
    int reserves;
};
