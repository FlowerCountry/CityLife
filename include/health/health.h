/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-30 07:03:53
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 14:25:16
 * @FilePath: \CityLife\include\health\health.h
 */
#pragma once

#include <string>

class Health {
  public:
    Health(const Health &b) : name{b.name}, reserves{b.reserves} {}
    Health(std::string name, int reserves) : name{name}, reserves{reserves} {}
    std::string GetInfo() { return name; }
    int GetReserves() { return reserves; }
    Health *operator+(Health b) { return new Health(name, reserves + b.reserves); }
    void operator+=(Health b) { reserves += b.reserves; }

  private:
    std::string name;
    int reserves;
};
