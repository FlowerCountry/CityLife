/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:33:32
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-24 21:16:39
 * @FilePath: \CityLife\src\controller\controller.cpp
 */
#include "controller/controller.h"
#include "curses.h"
#include "object/object.h"
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

#include <curses.h>
#include <string>
#include <vector>

int Controller::choose(std::string str, std::vector<Object *> options)
{
    int now = 0, ch = 0, length = options.size();
    int windowSize = 12;
    do
    {
        clear();
        View::getInstance()->print(0, str, 0);
        for (int i = (length <= windowSize ? 0 : (now - windowSize / 2 + length) % length), count = 0; count < (length <= windowSize ? length : windowSize); i = (i + 1) % length, count++)
        {
            if (i == now)
                View::getInstance()->print(count + 1, "> ", 0);
            else
                View::getInstance()->print(count + 1, "  ", 0);
            View::getInstance()->print(count + 1, options[i]->GetInfo().c_str(), 2);
        }
        refresh();
        ch = getch();
        int add = 0;
        if (ch == KEY_UP || ch == KEY_LEFT) add = -1;
        if (ch == KEY_DOWN || ch == KEY_RIGHT) add = 1;
        now = (now + add + length) % length;
    } while (ch != '\n');
    return now;
}
