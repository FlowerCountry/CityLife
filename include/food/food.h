/*
 * @Author: FlowerCity qzrobotsnake@gmail.com
 * @Date: 2025-02-02 20:00:35
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-02-02 20:15:27
 * @FilePath: \CityLife\include\food\food.h
 */
#pragma once

#include "health/health.h"
#include <string>
#include <vector>

class Health;

class Food {
  public:
    Food(const std::string &name, const std::vector<Health *> &health, const int &fresh) : name(name), health(health), fresh(fresh) {};
    ~Food() { health.clear(); }

  private:
    std::string name;
    std::vector<Health *> health;
    int fresh;
};