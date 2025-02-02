/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-28 10:08:56
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-02-02 20:09:13
 * @FilePath: \CityLife\src\center\center.cpp
 */
#include "center/center.h"
#include "controller/controller.h"

Center::Center()
{
    RandArray = {1, 7, 2, 5, 4, 3, 6, 9, 0, 8};
    srand(static_cast<unsigned int>(time(0)));
    InsertStringRandomly(announcement, "于2009年7月,我市第一家银行正式完工");
    InsertStringRandomly(announcement, "于2009年12月,我市第一家银行正式营业,有需要者可以到银行办理业务");
    InsertStringRandomly(announcement, "于2010年5月,我市预计开始建设电信大楼");
    InsertStringRandomly(announcement, "于2010年8月,我市预计正式完工电信大楼");
}

Center::~Center()
{
}

Center *Center::instance = nullptr;

Center *Center::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new Center();
    }
    return instance;
}

void Center::PrintAnnouncement()
{
    Controller::GetInstance()->PrintInfromathin(0, 0, announcement);
}

void Center::InsertStringRandomly(std::vector<std::string> &vec, const std::string &str)
{
    vec.insert(vec.begin() + RandArray[rand() % 10] % (vec.size() + 1), str);
}