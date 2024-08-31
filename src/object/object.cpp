/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-31 07:00:04
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
    world->buy();
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
    View::getInstance()->print(0, "请输入你要存的钱:");
    int money = View::getInstance()->scan();
    world->DepositingMoney(money);
}
void WithdrawMoney::ToDoIt(World *world)
{
    View::getInstance()->print(0, "请输入你要取的钱:");
    int money = View::getInstance()->scan();
    world->WithdrawMoney(money);
}

void Commodity::ToDoIt(World *world)
{
}