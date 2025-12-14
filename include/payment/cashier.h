#pragma once

#include <array>

// Cashier 负责基于现有钱包面额完成付款与找零。
class Cashier
{
  public:
    struct Result
    {
        bool success = false;
        std::array<int, 6> used{};    // 各面额实际扣除张数
        std::array<int, 6> change{};  // 找零各面额张数
        int paid = 0;                 // 实付总金额
    };

    explicit Cashier(const std::array<int, 6> &denominations);

    Result pay(int amount, const std::array<int, 6> &wallet) const;

  private:
    std::array<int, 6> denominations;
};

