/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:19:47
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-31 07:37:48
 * @FilePath: \CityLife\include\view\view.h
 */
#pragma once
#ifndef VIEW_H
#define VIEW_H

#include <string>

class View {
  public:
    View();
    ~View();
    static View *getInstance();
    void print(int line, const std::string str, int addpos = 0);
    int scan();

  private:
    static View *instance;
};

#endif // VIEW_H