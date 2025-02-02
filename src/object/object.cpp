/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 09:42:34
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 11:02:15
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
    World::GetInstance()->ChangeWhere(this->to);
}

void Buy::ToDoIt()
{
    World::GetInstance()->ChangeWhere(3);
}
Information::Information(const std::string &name, const std::string &content) : Object(name), content{content}
{
    things["公告"] = std::bind(&Center::PrintAnnouncement, Center::GetInstance());
}
void Information::ToDoIt()
{
    things[content]();
}
void DepositingMoney::ToDoIt()
{
    clear();
    refresh();
    View::GetInstance()->Print(0, 0, "请输入你要存的钱:");
    int money = View::GetInstance()->Scan(0, 17);
    World::GetInstance()->DepositingMoney(money);
}
void WithdrawMoney::ToDoIt()
{
    clear();
    refresh();
    View::GetInstance()->Print(0, 0, "请输入你要取的钱:");
    int money = View::GetInstance()->Scan(0, 17);
    World::GetInstance()->WithdrawMoney(money);
}

Commodity::Commodity(const std::string &name, const int &price, const std::vector<Health> &health) : Object{name}, price{price}, health{health} {}

std::vector<Health> Commodity::GetHealth() { return health; }

void Commodity::ToDoIt()
{
    if (World::GetInstance()->SpendMoney(price))
    {
        for (auto i : health)
        {
            (*User::GetInstance()->GetUserHealth())[i.GetInfo()] = std::min((*User::GetInstance()->GetUserHealth())[i.GetInfo()] + i.GetReserves(), 100);
        }
        clear();
        int st = 0;
        View::GetInstance()->Print(st, 0, "购买" + name + "成功");
        for (auto i : health)
        {
            View::GetInstance()->Print(++st, 0, i.GetInfo() + "达到了" + std::to_string((*User::GetInstance()->GetUserHealth())[i.GetInfo()]));
        }
        getch();
    }
    World::GetInstance()->ChangeWhere(1);
}