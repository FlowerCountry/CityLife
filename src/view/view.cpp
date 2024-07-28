/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:21:42
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-26 10:49:12
 * @FilePath: \CityLive\src\view\view.cpp
 */
#include "view/view.h"

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

void View::print(const std::string str, const char end)
{
    int length = str.length();
    for (int i = 0; i < length; i++)
    {
        putchar(str[i]);
    }
    putchar(end);
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