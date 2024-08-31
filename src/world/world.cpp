/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:40:49
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-31 07:35:11
 * @FilePath: \CityLife\src\world\world.cpp
 */
#include "world/world.h"
#include "bank/bank.h"
#include "controller/controller.h"
#include "curses.h"
#include "health/health.h"
#include "object/object.h"
#include "view/view.h"
#include "windows.h"

World::World()
{
    initscr();
    noecho();
    cbreak();
    SetConsoleOutputCP(65001);
    curs_set(0);
    keypad(stdscr, TRUE); // 启用方向键输入
    clear();
    refresh();
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
        new class Building(0, 5, 5, "市中心"),
        new class Building(1, 3, 5, "超市"),
        new class Building(2, 5, 3, "银行"),
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
            if (i->id != j->id)
            {
                ToDoThings[i->id].push_back(new class GoWhere("前往" + j->name, i->id, j->id, GetLen(*i, *j)));
            }
        }
    }
    Commodities = {
        new class Commodity("面包", 15, {new Health("饱腹感", 10), new Health("饥饿", 10), new Health("蛋白质", 5), new Health("满足感", 7 * LifeQuality)}),
        new class Commodity("牛奶", 25, {new Health("饱腹感", 5), new Health("饥饿", 8), new Health("钙", 12), new Health("蛋白质", 7), new Health("满足感", 6 * LifeQuality)}),
        new class Commodity("蛋糕", 40, {new Health("饱腹感", 12), new Health("饥饿", 5), new Health("糖分", 15), new Health("脂肪", 10), new Health("满足感", 9 * LifeQuality)}),
        new class Commodity("苹果", 10, {new Health("饱腹感", 8), new Health("维生素C", 10), new Health("纤维素", 7), new Health("满足感", 6 * LifeQuality)}),
        new class Commodity("牛排", 150, {new Health("饱腹感", 20), new Health("蛋白质", 25), new Health("铁", 15), new Health("脂肪", 8), new Health("满足感", 70 * LifeQuality)}),
        new class Commodity("橙汁", 18, {new Health("饱腹感", 5), new Health("维生素C", 15), new Health("饥渴", 10), new Health("糖分", 5), new Health("满足感", 7 * LifeQuality)}),
        new class Commodity("沙拉", 35, {new Health("饱腹感", 6), new Health("维生素", 18), new Health("纤维素", 12), new Health("满足感", 8 * LifeQuality)}),
        new class Commodity("鸡蛋", 12, {new Health("饱腹感", 8), new Health("蛋白质", 12), new Health("满足感", 7 * LifeQuality)}),
        new class Commodity("意大利面", 30, {new Health("饱腹感", 15), new Health("碳水化合物", 20), new Health("满足感", 8 * LifeQuality)}),
        new class Commodity("咖啡", 20, {new Health("饱腹感", 2), new Health("精神振奋", 10), new Health("饥渴", 5), new Health("满足感", 6 * LifeQuality)}),
        new class Commodity("巧克力", 25, {new Health("饱腹感", 3), new Health("幸福感", 12), new Health("糖分", 15), new Health("脂肪", 8), new Health("满足感", 9 * LifeQuality)}),
        new class Commodity("番茄", 7, {new Health("饱腹感", 4), new Health("维生素C", 8), new Health("纤维素", 6), new Health("满足感", 5 * LifeQuality)}),
        new class Commodity("鸡胸肉", 70, {new Health("饱腹感", 18), new Health("蛋白质", 30), new Health("脂肪", 5), new Health("满足感", 65 * LifeQuality)}),
        new class Commodity("矿泉水", 3, {new Health("饱腹感", 1), new Health("饥渴", 15), new Health("满足感", 4 * LifeQuality)}),
        new class Commodity("米饭", 15, {new Health("饱腹感", 12), new Health("碳水化合物", 18), new Health("满足感", 6 * LifeQuality)}),
        new class Commodity("燕麦片", 22, {new Health("饱腹感", 14), new Health("纤维素", 10), new Health("满足感", 7 * LifeQuality)}),
        new class Commodity("酸奶", 18, {new Health("饱腹感", 7), new Health("钙", 12), new Health("蛋白质", 8), new Health("满足感", 8 * LifeQuality)}),
        new class Commodity("三明治", 40, {new Health("饱腹感", 16), new Health("蛋白质", 10), new Health("满足感", 10 * LifeQuality)}),
        new class Commodity("薯片", 12, {new Health("饱腹感", 3), new Health("脂肪", 12), new Health("碳水化合物", 8), new Health("满足感", 6 * LifeQuality)}),
        new class Commodity("冰淇淋", 35, {new Health("饱腹感", 4), new Health("幸福感", 15), new Health("糖分", 18), new Health("脂肪", 10), new Health("满足感", 9 * LifeQuality)})};
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
        int choose = Controller::getInstance()->choose("你当前位于: " + BuildingNames[where], ToDoThings[where]);
        ToDoThings[where][choose]->ToDoIt(this);
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
        View::getInstance()->print(0, "你没有这么多的钱");
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
        View::getInstance()->print(0, "你没有这么多的钱");
        getchar();
    }
}

int World::GetLen(Building a, Building b)
{
    return abs(a.x - b.x) + abs(a.y - b.y);
}

void World::buy()
{
    int choose = Controller::getInstance()->choose("你当前位于: " + BuildingNames[where], Commodities);
}