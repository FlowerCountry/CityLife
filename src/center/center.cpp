/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-28 10:08:56
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 10:46:36
 * @FilePath: \CityLive\src\center\center.cpp
 */
#include "center/center.h"
#include "view/view.h"
#include "world/world.h"

Center::Center()
{
    RandArray = {1, 7, 2, 5, 4, 3, 6, 9, 0, 8};
    srand(static_cast<unsigned int>(time(0)));
    InsertStringRandomly(announcement, "于2009年7月，我市第一家银行正式完工");
    InsertStringRandomly(announcement, "于2009年12月，我市第一家银行正式营业，有需要者可以到银行办理业务");
    InsertStringRandomly(announcement, "于2010年5月，我市预计开始建设电信大楼");
    InsertStringRandomly(announcement, "于2010年8月，我市预计正式完工电信大楼");
}

Center::~Center()
{
}

Center *Center::instance = nullptr;

Center *Center::getInstance()
{
    if (instance == nullptr)
    {
        instance = new Center();
    }
    return instance;
}

void Center::PrintAnnouncement(World *world)
{
    for (auto &i : announcement)
    {
        View::getInstance()->print(i, '\n');
    }
    getchar();
}

void Center::InsertStringRandomly(std::vector<std::string> &vec, const std::string &str)
{
    vec.insert(vec.begin() + RandArray[rand() % 10] % (vec.size() + 1), str);
}