/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:40:49
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 22:02:46
 * @FilePath: \CityLife\src\world\world.cpp
 */
#include "world/world.h"
#include "bank/bank.h"
#include "controller/controller.h"
#include "curses.h"
#include "health/health.h"
#include "object/object.h"
#include "view/view.h"

World::World()
{
    initscr();
    keypad(stdscr, TRUE);
    scrollok(stdscr, FALSE);
    setlocale(LC_ALL, "");
    cbreak();
    noecho();
    for (int i = 0; i < 10; i++)
    {
        move(i, 0);
        printw("%-*s", 100, "");
    }
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
        "超市内",
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
        {new class Commodity("面包", 15, {new Health("饱腹感", 18), new Health("饥饿", 18), new Health("蛋白质", 9), new Health("维生素B", 9)}),
         new class Commodity("牛奶", 25, {new Health("饱腹感", 10), new Health("饥饿", 16), new Health("钙", 24), new Health("蛋白质", 14), new Health("维生素A", 12)}),
         new class Commodity("蛋糕", 40, {new Health("饱腹感", 24), new Health("饥饿", 10), new Health("糖分", 30), new Health("脂肪", 20), new Health("维生素E", 14)}),
         new class Commodity("苹果", 10, {new Health("饱腹感", 16), new Health("维生素C", 20), new Health("纤维素", 14), new Health("钾", 10)}),
         new class Commodity("牛排", 150, {new Health("饱腹感", 20), new Health("蛋白质", 25), new Health("铁", 15), new Health("脂肪", 8), new Health("锌", 32)}),
         new class Commodity("橙汁", 18, {new Health("饱腹感", 7), new Health("维生素C", 22), new Health("饥渴", 14), new Health("糖分", 7), new Health("钙", 10)}),
         new class Commodity("沙拉", 35, {new Health("饱腹感", 9), new Health("维生素A", 30), new Health("纤维素", 20), new Health("钾", 14)}),
         new class Commodity("鸡蛋", 12, {new Health("饱腹感", 16), new Health("蛋白质", 24), new Health("硒", 12)}),
         new class Commodity("意大利面", 30, {new Health("饱腹感", 24), new Health("碳水化合物", 30), new Health("维生素B", 12)}),
         new class Commodity("咖啡", 20, {new Health("饱腹感", 4), new Health("精神振奋", 20), new Health("饥渴", 10), new Health("钾", 12)}),
         new class Commodity("巧克力", 25, {new Health("饱腹感", 6), new Health("幸福感", 24), new Health("糖分", 30), new Health("脂肪", 16), new Health("镁", 12)}),
         new class Commodity("番茄", 7, {new Health("饱腹感", 8), new Health("维生素C", 16), new Health("纤维素", 12), new Health("钾", 8)}),
         new class Commodity("鸡胸肉", 70, {new Health("饱腹感", 24), new Health("蛋白质", 36), new Health("脂肪", 6), new Health("铁", 26)}),
         new class Commodity("矿泉水", 3, {new Health("饱腹感", 2), new Health("饥渴", 45), new Health("钙", 5)}),
         new class Commodity("米饭", 15, {new Health("饱腹感", 24), new Health("碳水化合物", 36), new Health("维生素B", 12)}),
         new class Commodity("燕麦片", 22, {new Health("饱腹感", 28), new Health("纤维素", 20), new Health("钾", 14)}),
         new class Commodity("酸奶", 18, {new Health("饱腹感", 12), new Health("钙", 24), new Health("蛋白质", 16), new Health("益生菌", 16)}),
         new class Commodity("三明治", 40, {new Health("饱腹感", 32), new Health("蛋白质", 20), new Health("铁", 18)}),
         new class Commodity("薯片", 12, {new Health("饱腹感", 6), new Health("脂肪", 24), new Health("碳水化合物", 16), new Health("维生素C", 12)}),
         new class Commodity("冰淇淋", 35, {new Health("饱腹感", 6), new Health("幸福感", 30), new Health("糖分", 36), new Health("脂肪", 20), new Health("维生素D", 14)})},
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
        clear();
        refresh();
    }
}

void World::changewhere(int where)
{
    this->where = where;
}

bool World::SpendMoney(int money)
{
    if (this->money >= money)
    {
        this->money -= money;
        return true;
    }
    else
    {
        clear();
        View::getInstance()->print(0, "你没有这么多的钱", 0);
        getch();
        return false;
    }
}

bool World::DepositingMoney(int money)
{
    if (this->money >= money)
    {
        this->money -= money;
        bank->DepositingMoney(money);
        return true;
    }
    else
    {
        clear();
        View::getInstance()->print(0, "你没有这么多的钱", 0);
        getch();
        return false;
    }
}

bool World::WithdrawMoney(int money)
{
    if (bank->GetMoney() >= money)
    {
        this->money += bank->WithdrawMoney(money);
        return true;
    }
    else
    {
        clear();
        View::getInstance()->print(0, "你没有这么多的钱", 0);
        getch();
        return false;
    }
}

int World::GetLen(Building a, Building b)
{
    return abs(a.x - b.x) + abs(a.y - b.y);
}