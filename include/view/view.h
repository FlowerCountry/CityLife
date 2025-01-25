/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:19:47
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 15:39:43
 * @FilePath: \CityLife\include\view\view.h
 */
#pragma once

#include <string>

class View {
  public:
    /**
     * @description: 显示类构造函数
     * @return null
     */
    View();
    /**
     * @description: 显示类析构函数
     * @return null
     */
    ~View();
    /**
     * @description: 返回显示类单例对象
     * @return {*} 显示类单例对象
     */
    static View *getInstance();
    /**
     * @description: 打印内容
     * @param {int} line 行数, 第几行
     * @param {string} str 内容
     * @param {int} addpos 列数, 第几列
     * @return null
     */
    void print(int line, const std::string str, int addpos = 0);
    /**
     * @description: 读入内容
     * @param {int} line 行数, 第几行
     * @param {int} addpos 列数, 第几列
     * @return {*} 输入的内容
     */
    int scan(int line, int addpos = 0);

  private:
    static View *instance;
};
