#include "view/view.h"

#include <cstdio>
#include <cstdlib>
#include <iostream>
#include <sstream>
#include <stdexcept>

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

void View::Print(const int &line, const int &addpos, const std::string &str, bool newline)
{
    (void)line;
    if (addpos > 0)
    {
        std::cout << std::string(static_cast<std::size_t>(addpos), ' ');
    }
    if (newline)
    {
        std::cout << str << std::endl;
        return;
    }
    std::cout << str << std::flush;
}

int View::Scan(const int &line, const int &addpos)
{
    (void)line;
    (void)addpos;
    while (true)
    {
        std::string buffer;
        if (!std::getline(std::cin, buffer))
        {
            throw std::runtime_error("INPUT_EOF");
        }
        std::stringstream ss(buffer);
        int value = 0;
        char extra = '\0';
        if (ss >> value && !(ss >> extra))
        {
            std::cout << std::endl;
            return value;
        }
        std::cout << "输入无效，请输入一个整数: " << std::flush;
    }
}

void View::Clear()
{
#ifdef _WIN32
    system("cls");
#else
    system("clear");
#endif
}

void View::WaitForEnter()
{
    int ch = std::getchar();
    if (ch == EOF)
    {
        throw std::runtime_error("INPUT_EOF");
    }
}
