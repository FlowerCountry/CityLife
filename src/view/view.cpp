/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:21:42
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-08-31 07:37:38
 * @FilePath: \CityLife\src\view\view.cpp
 */
#include "view/view.h"
#include "curses.h"

View::View()
{
}

View::~View()
{
}
View *View::instance = nullptr;
View *View::getInstance()
{
    if (instance == nullptr)
    {
        instance = new View();
    }
    return instance;
}

void View::print(int line, const std::string str, int addpos)
{
    move(line, addpos);
    clrtoeol();
    printw("%s", str.c_str());
    refresh();
}

int View::scan()
{
    int f = 1, k = 0;
    char c = getchar();
    while (c < '0' || c > '9')
    {
        if (c == '-')
            f = -1;
        c = getchar();
    }
    while (c >= '0' && c <= '9')
    {
        k = (k << 1) + (k << 3) + (c ^ 48);
        c = getchar();
    }
    return f * k;
}