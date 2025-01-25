/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-25 15:18:15
 * @FilePath: \CityLife\src\object\object.cpp
 */
#include "object/object.h"
#include "center/center.h"
#include "health/health.h"
#include "user/user.h"
#include "view/view.h"
#include "world/world.h"

void GoWhere::ToDoIt()
{
    World::getInstance()->changewhere(this->to);
}

void Buy::ToDoIt()
{
    World::getInstance()->changewhere(3);
}
Information::Information(std::string name, std::string content) : Object(name), content{content}
{
    things["公告"] = std::bind(&Center::PrintAnnouncement, Center::getInstance());
}
void Information::ToDoIt()
{
    things[content]();
}
void DepositingMoney::ToDoIt()
{
    clear();
    refresh();
    View::getInstance()->print(0, "请输入你要存的钱:", 0);
    int money = View::getInstance()->scan(0, 17);
    World::getInstance()->DepositingMoney(money);
}
void WithdrawMoney::ToDoIt()
{
    clear();
    refresh();
    View::getInstance()->print(0, "请输入你要取的钱:", 0);
    int money = View::getInstance()->scan(0, 17);
    World::getInstance()->WithdrawMoney(money);
}

Commodity::Commodity(std::string name, int price, std::vector<Health> health) : Object{name}, price{price}, health{health} {}

std::vector<Health> Commodity::GetHealth() { return health; }

void Commodity::ToDoIt()
{
    if (World::getInstance()->SpendMoney(price))
    {
        for (auto i : health)
        {
            (*User::getInstance()->GetUserHealth())[i.GetInfo()] = std::min((*User::getInstance()->GetUserHealth())[i.GetInfo()] + i.GetReserves(), 100);
        }
        clear();
        int st = 0;
        View::getInstance()->print(st, "购买" + name + "成功");
        for (auto i : health)
        {
            View::getInstance()->print(++st, i.GetInfo() + "达到了" + std::to_string((*User::getInstance()->GetUserHealth())[i.GetInfo()]));
        }
        getch();
    }
    World::getInstance()->changewhere(1);
}