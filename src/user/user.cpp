/*
 * @Author: FlowerCity qzrobotsnake@gmail.com
 * @Date: 2025-01-24 21:24:28
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 22:08:12
 * @FilePath: \CityLife\src\user\user.cpp
 */
#include "user/user.h"
#include "health/health.h"

User::User()
{
    UserHealth = {
        {"饱腹感", 100},
        {"饥饿", 100},
        {"蛋白质", 100},
        {"钙", 100},
        {"糖分", 100},
        {"脂肪", 100},
        {"纤维素", 100},
        {"铁", 100},
        {"碳水化合物", 100},
        {"维生素A", 100},
        {"维生素B", 100},
        {"维生素C", 100},
        {"维生素D", 100},
        {"维生素E", 100},
        {"钾", 100},
        {"硒", 100},
        {"锌", 100}};
}

User::~User() {}

User *User::instance = nullptr;

User *User::getInstance()
{
    if (instance == nullptr)
    {
        instance = new User();
    }
    return instance;
}