/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:40:49
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-30 07:09:31
 * @FilePath: \CityLife\src\world\world.cpp
 */
#include "world/world.h"
#include "bank/bank.h"
#include "controller/controller.h"
#include "health/health.h"
#include "thing/thing.h"
#include "view/view.h"

World::World()
{
    system("chcp 565001");
    year = 2010;
    month = 10;
    day = 10;
    hour = 11;
    minute = 3;
    second = 0;
    where = 0;
    money = 100;
    bank = new Bank();
    BuildingNames = {
        "市中心",
        "超市",
        "银行",
    };
    Buildings = {
        Building(0, 5, 5, "市中心"),
        Building(1, 3, 5, "超市"),
        Building(2, 5, 3, "银行"),
    };
    ToDoThings = {
        {new class Information("查看公告", "公告")},
        {new class Buy("购买物品")},
        {new class DepositingMoney("存钱"), new class WithdrawMoney("取钱")},
    };
    for (auto i : Buildings)
    {
        for (auto j : Buildings)
        {
            if (i.id != j.id)
            {
                ToDoThings[i.id].push_back(new class GoWhere("前往" + j.name, i.id, j.id, GetLen(i, j)));
            }
        }
    }
    Commodities = {
        new class Commodity("面包", 10),
        new class Commodity("牛奶", 20),
        new class Commodity("蛋糕", 50)};
}

World::~World()
{
}

World *World::instance = nullptr;

World *World::getInstance()
{
    if (instance == nullptr)
    {
        instance = new World();
    }
    return instance;
}

void World::start()
{
    while (true)
    {
        system("cls");
        int choose = Controller::getInstance()->choose("你当前位于: " + BuildingNames[where], ToDoThings[where]);
        ToDoThings[where][choose]->ToDoThing(this);
    }
}

void World::changewhere(int where)
{
    this->where = where;
}

void World::DepositingMoney(int money)
{
    if (this->money >= money)
    {
        this->money -= money;
        bank->DepositingMoney(money);
    }
    else
    {
        View::getInstance()->print("你没有这么多的钱", '\n');
        getchar();
    }
}

void World::WithdrawMoney(int money)
{
    if (bank->GetMoney() >= money)
    {
        this->money += bank->WithdrawMoney(money);
    }
    else
    {
        View::getInstance()->print("你没有这么多的钱", '\n');
        getchar();
    }
}

int World::GetLen(Building a, Building b)
{
    return abs(a.x - b.x) + abs(a.y - b.y);
}

void World::buy()
{
}