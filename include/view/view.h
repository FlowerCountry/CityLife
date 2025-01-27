/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:19:47
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 11:00:28
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
    void print(int line, int addpos, const std::string str);
    /**
     * @description: 读入内容
     * @param {int} line 行数, 第几行
     * @param {int} addpos 列数, 第几列
     * @return {*} 输入的内容
     */
    int scan(int line, int addpos);

  private:
    static View *instance;
};
