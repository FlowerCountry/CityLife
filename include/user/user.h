/*
 * @Author: FlowerCity qzrobotsnake@gmail.com
 * @Date: 2025-01-24 21:20:49
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 15:39:57
 * @FilePath: \CityLife\include\user\user.h
 */
#pragma once

#include <map>
#include <string>

class Health;
class User {
  public:
    /**
     * @description: 用户类构造函数
     * @return null
     */
    User();
    /**
     * @description: 用户类析构函数
     * @return null
     */
    ~User();
    /**
     * @description: 返回用户类单例对象
     * @return {*} 用户类单例对象
     */
    static User *getInstance();
    /**
     * @description: 返回用户的营养信息
     * @return {*} 用户的营养信息
     */
    std::map<std::string, int> *GetUserHealth() { return &UserHealth; }

  private:
    static User *instance;
    std::map<std::string, int> UserHealth;
};