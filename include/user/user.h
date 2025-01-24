/*
 * @Author: FlowerCity qzrobotsnake@gmail.com
 * @Date: 2025-01-24 21:20:49
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 21:47:30
 * @FilePath: \CityLife\include\user\user.h
 */
#pragma once

#include "health/health.h"
#include <map>
#include <string>

class Health;
class User {
  public:
    User();
    ~User();
    static User *getInstance();
    std::map<std::string, int> UserHealth;

  private:
    static User *instance;
};