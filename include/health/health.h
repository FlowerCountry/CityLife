/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-30 07:03:53
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-30 07:15:16
 * @FilePath: \CityLife\include\health\health.h
 */
#pragma once
#ifndef HEALTH_H
#define HEALTH_H

#include <string>

class Health {
  public:
    Health(const Health &b) : name{b.name}, reserves{b.reserves} {}
    Health(std::string name, int reserves) : name{name}, reserves{reserves} {}
    std::string GetInfo() { return name; }
    Health *operator+(Health b) { return new Health(name, reserves + b.reserves); }
    void operator+=(Health b) { reserves += b.reserves; }

  private:
    std::string name;
    int reserves;
};

#endif