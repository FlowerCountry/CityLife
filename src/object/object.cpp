/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 15:23:48
 * @FilePath: \CityLife\src\object\object.cpp
 */
#include "object/object.h"
#include "center/center.h"
#include "view/view.h"
#include "world/world.h"
#include <health/health.h>

void GoWhere::ToDoIt(World *world)
{
    world->changewhere(this->to);
}

void Buy::ToDoIt(World *world)
{
    world->changewhere(3);
}
Information::Information(std::string name, std::string content) : Object(name), content{content}
{
    things["公告"] = std::bind(&Center::PrintAnnouncement, Center::getInstance(), std::placeholders::_1);
}
void Information::ToDoIt(World *world)
{
    things[content](world);
}
void DepositingMoney::ToDoIt(World *world)
{
    clear();
    refresh();
    View::getInstance()->print(0, "请输入你要存的钱:", 0);
    int money = View::getInstance()->scan(0, 17);
    world->DepositingMoney(money);
}
void WithdrawMoney::ToDoIt(World *world)
{
    clear();
    refresh();
    View::getInstance()->print(0, "请输入你要取的钱:", 0);
    int money = View::getInstance()->scan(0, 17);
    world->WithdrawMoney(money);
}

void Commodity::ToDoIt(World *world)
{
    world->changewhere(1);
}