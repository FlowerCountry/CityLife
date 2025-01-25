/*
 * @Author: FlowerCity qzrobotsnake@gmail.com
 * @Date: 2025-01-24 21:20:49
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 13:06:21
 * @FilePath: \CityLife\include\user\user.h
 */
#pragma once

#include <map>
#include <string>

class Health;
class User {
  public:
    User();
    ~User();
    static User *getInstance();
    std::map<std::string, int> *GetUserHealth() { return &UserHealth; }

  private:
    static User *instance;
    std::map<std::string, int> UserHealth;
};