/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:19:47
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-26 10:47:43
 * @FilePath: \CityLive\include\view\view.h
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
    void print(const std::string str, const char end);
    int scan();

  private:
    static View *instance;
};

#endif // VIEW_H