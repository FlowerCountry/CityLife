/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:21:42
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 11:00:20
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
View *View::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new View();
    }
    return instance;
}

void View::Print(int line, int addpos, const std::string str)
{
    move(line, addpos);
    printw("%-*s", 100, "");
    move(line, addpos);
    printw("%s", str.c_str());
    refresh();
}

int View::Scan(int line, int addpos)
{
    move(line, addpos);
    clrtoeol();
    int f = 1, k = 0;
    char c = getch();
    while (c < '0' || c > '9')
    {
        if (c == '-')
        {
            f = -1;
        }
        c = getch();
    }
    while (c >= '0' && c <= '9')
    {
        addch(c);
        k = (k << 1) + (k << 3) + (c - '0');
        c = getch();
    }
    refresh();
    return f * k;
}
