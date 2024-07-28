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