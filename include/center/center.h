/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-28 10:07:19
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-27 20:49:44
 * @FilePath: \CityLife\include\center\center.h
 */
#pragma once
#ifndef CENTER_H
#define CENTER_H

#include <cstdlib>
#include <ctime>
#include <string>
#include <vector>

class World;
class View;
class Center {
  public:
    Center();
    ~Center();
    static Center *getInstance();
    void PrintAnnouncement(World *world);
    void InsertStringRandomly(std::vector<std::string> &vec, const std::string &str);

  private:
    static Center *instance;
    std::vector<std::string> announcement;
    std::vector<int> RandArray;
};

#endif