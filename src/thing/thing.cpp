#include "thing/thing.h"
#include "center/center.h"
#include "view/view.h"
#include "world/world.h"

void GoWhere::ToDoThing(World *world)
{
    world->changewhere(this->to);
}

void Buy::ToDoThing(World *world)
{
    // world->buy(get);
}
Information::Information(std::string name, std::string content) : Thing(name), content{content}
{
    things["公告"] = std::bind(&Center::PrintAnnouncement, Center::getInstance(), std::placeholders::_1);
}
void Information::ToDoThing(World *world)
{
    things[content](world);
}
void DepositingMoney::ToDoThing(World *world)
{
    View::getInstance()->print("请输入你要存的钱:", ' ');
    int money = View::getInstance()->scan();
    world->DepositingMoney(money);
}
void WithdrawMoney::ToDoThing(World *world)
{
    View::getInstance()->print("请输入你要取的钱:", ' ');
    int money = View::getInstance()->scan();
    world->WithdrawMoney(money);
}