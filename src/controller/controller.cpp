/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:33:32
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 10:16:45
 * @FilePath: \CityLive\src\controller\controller.cpp
 */
#include "controller/controller.h"
#include "thing/thing.h"
#include "view/view.h"
#include "world/world.h"
Controller::Controller()
{
}

Controller::~Controller()
{
}

Controller *Controller::instance = nullptr;

Controller *Controller::getInstance()
{
    if (instance == nullptr)
    {
        instance = new Controller();
    }
    return instance;
}

int Controller::choose(std::string str, std::vector<Thing *> options)
{
    int now = 0, ch = 0, length = options.size();
    do
    {
        View::getInstance()->print(str, '\n');
        for (int i = 0; i < length; i++)
        {
            if (i == now)
                View::getInstance()->print(">", ' ');
            else
                View::getInstance()->print(" ", ' ');
            View::getInstance()->print(options[i]->GetInfo(), '\n');
        }
        ch = getch();
        system("cls");
        int add = 0;
        if (ch == 72 || ch == 75) add = -1;
        if (ch == 77 || ch == 80) add = 1;
        now = (now + add + length) % length;
    } while (ch != 13);
    return now;
}