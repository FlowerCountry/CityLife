/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:32:47
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 09:52:23
 * @FilePath: \CityLive\include\controller\controller.h
 */
#pragma once
#ifndef CONTROLLER_H
#define CONTROLLER_H

#include <conio.h>
#include <string>
#include <vector>
class Thing;
class View;
class World;
class Controller {
  public:
    Controller();
    ~Controller();
    static Controller *getInstance();
    int choose(std::string str, std::vector<class Thing *> options);

  private:
    static Controller *instance;
};

#endif // CONTROLLER_H