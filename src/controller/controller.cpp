/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-25 18:33:32
 * @LastEditors: FlowerCity qzrobotsnake@gmail.com
 * @LastEditTime: 2025-01-26 11:07:03
 * @FilePath: \CityLife\src\controller\controller.cpp
 */
#include "controller/controller.h"
#include "object/object.h"
#include "view/view.h"
#include "terminal/terminal.h"
#include "view/ui_env.h"
#include "world/world.h"

#include <algorithm>
#include <sstream>
#include <string>
#include <vector>
 

Controller::Controller()
{
}

Controller::~Controller()
{
}

Controller *Controller::instance = nullptr;

Controller *Controller::GetInstance()
{
    if (instance == nullptr)
    {
        instance = new Controller();
    }
    return instance;
}

int Controller::Choose(const std::string &str, const std::vector<Object *> &options)
{
    if (options.empty())
    {
        return -1;
    }

    View *view = View::GetInstance();
    Terminal terminal;

    // 非交互（测试）环境：退回到老的数字选择，保持测试输出稳定
    if (!UiEnv::IsInteractive())
    {
        while (true)
        {
            view->Clear();
            int line = 0;
            view->Print(line++, 0, str);
            for (std::size_t i = 0; i < options.size(); ++i)
            {
                view->Print(line++, 0, std::to_string(i + 1) + ". " + options[i]->GetInfo());
            }
            view->Print(line, 0, "请输入选项编号: ");
            int choose = View::GetInstance()->Scan(0, 0);
            if (choose >= 1 && choose <= static_cast<int>(options.size()))
            {
                return choose - 1;
            }
            // 无效输入，回车继续
            view->Print(0, 0, "输入无效，请按回车重新输入。\n");
            view->WaitForEnter();
        }
    }

    const int total = static_cast<int>(options.size());
    std::vector<std::string> labels;
    labels.reserve(options.size());
    for (int i = 0; i < total; ++i)
    {
        labels.push_back(std::to_string(i + 1) + ". " + options[static_cast<std::size_t>(i)]->GetInfo());
    }

    int focus = 0;
    int columnWidth = 0;
    for (const auto &label : labels)
    {
        columnWidth = std::max(columnWidth, static_cast<int>(label.size()));
    }

    const Terminal::WindowSize windowSize = terminal.getWindowSize();
    const int maxVisible = std::max(1, windowSize.rows - 5);
    const int visibleRows = std::min(total, maxVisible);
    int start = 0;

    auto renderFrame = [&](int active, int startIndex) {
        std::ostringstream frame;
        frame << "\033[H\033[2J";
        frame << str << "\n\n";

        for (int i = 0; i < visibleRows; ++i)
        {
            const int index = startIndex + i;
            if (index >= total)
            {
                break;
            }
            const bool isActive = index == active;
            std::string text = (isActive ? "> " : "  ") + labels[static_cast<std::size_t>(index)];
            const int paddedWidth = columnWidth + 2;
            if (static_cast<int>(text.size()) < paddedWidth)
            {
                text += std::string(paddedWidth - static_cast<int>(text.size()), ' ');
            }
            frame << text << '\n';
        }

        if (total > visibleRows)
        {
            frame << "\n";
            frame << "使用方向键滚动 (" << (active + 1) << "/" << total << ")";
        }
        else
        {
            frame << "\n";
        }

        terminal.writeRaw(frame.str());
    };

    renderFrame(focus, start);

    while (true)
    {
        Terminal::KeyEvent event = terminal.readKey();
        switch (event.type)
        {
        case Terminal::KeyEvent::Type::ArrowUp:
            focus = (focus - 1 + total) % total;
            if (focus < start)
            {
                start = focus;
            }
            else if (focus >= start + visibleRows)
            {
                start = focus - visibleRows + 1;
            }
            renderFrame(focus, start);
            break;
        case Terminal::KeyEvent::Type::ArrowDown:
            focus = (focus + 1) % total;
            if (focus >= start + visibleRows)
            {
                start = focus - visibleRows + 1;
            }
            else if (focus < start)
            {
                start = focus;
            }
            renderFrame(focus, start);
            break;
        case Terminal::KeyEvent::Type::ArrowLeft:
            focus = (focus - 1 + total) % total;
            if (focus < start)
            {
                start = focus;
            }
            else if (focus >= start + visibleRows)
            {
                start = focus - visibleRows + 1;
            }
            renderFrame(focus, start);
            break;
        case Terminal::KeyEvent::Type::ArrowRight:
            focus = (focus + 1) % total;
            if (focus >= start + visibleRows)
            {
                start = focus - visibleRows + 1;
            }
            else if (focus < start)
            {
                start = focus;
            }
            renderFrame(focus, start);
            break;
        case Terminal::KeyEvent::Type::Enter:
            terminal.disableRawMode();
            return focus;
        case Terminal::KeyEvent::Type::EndOfInput:
            terminal.disableRawMode();
            throw std::runtime_error("INPUT_EOF");
        default:
            break;
        }
    }
}

void Controller::PrintInformation(const int &line, const int &addpos, const std::vector<std::string> &content)
{
    View::GetInstance()->Clear();
    for (int i = 0; i < content.size(); i++)
    {
        View::GetInstance()->Print(line + i, addpos, content[i]);
    }
    View::GetInstance()->WaitForEnter();
}
