/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:32:47
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 10:58:19
 * @FilePath: \CityLife\include\controller\controller.h
 */
#pragma once

#include <curses.h>
#include <string>
#include <vector>

class Object;
class View;
class World;
class Controller {
  public:
    /**
     * @description: 控制类的构造函数
     * @return null
     */
    Controller();
    /**
     * @description: 控制类的析构函数
     * @return null
     */
    ~Controller();
    /**
     * @description: 返回控制类的单例对象
     * @return {*} 控制类的单例对象
     */
    static Controller *GetInstance();
    /**
     * @description: 基本的选择函数
     * @param {string} str 整个内容的标题
     * @param {vector<class Object *>} options 选项
     * @return {*} 最终选了什么, 值为下标
     */
    int Choose(std::string str, std::vector<class Object *> options);
    void PrintInfromathin(int line, int addpos, std::vector<std::string> content);

  private:
    static Controller *instance;
};