/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-28 10:07:19
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 21:10:16
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
