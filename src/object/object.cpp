/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 22:14:54
 * @FilePath: \CityLife\src\object\object.cpp
 */
#include "object/object.h"
#include "bank/bank.h"
#include "center/center.h"
#include "health/health.h"
#include "user/user.h"
#include "view/view.h"
#include "world/world.h"
#include <string>

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
    if (world->SpendMoney(price))
    {
        for (auto i : health)
        {
            User::getInstance()->UserHealth[i->GetInfo()] = std::min(User::getInstance()->UserHealth[i->GetInfo()] + i->GetReserves(), 100);
        }
        clear();
        int st = 0;
        View::getInstance()->print(st, "购买" + name + "成功");
        for (auto i : health)
        {
            View::getInstance()->print(++st, i->GetInfo() + "达到了" + std::to_string(User::getInstance()->UserHealth[i->GetInfo()]));
            delete i;
        }
        getch();
    }
    world->changewhere(1);
}