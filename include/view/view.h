/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:19:47
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 15:20:40
 * @FilePath: \CityLife\include\view\view.h
 */
#pragma once

#include <string>

class View {
  public:
    View();
    ~View();
    static View *getInstance();
    void print(int line, const std::string str, int addpos = 0);
    int scan(int line, int addpos = 0);

  private:
    static View *instance;
};
