#include "payment/cashier.h"

#include <algorithm>

Cashier::Cashier(const std::array<int, 6> &denoms)
    : denominations(denoms)
{
}

Cashier::Result Cashier::pay(int amount, const std::array<int, 6> &wallet) const
{
    Result result;
    if (amount <= 0)
    {
        result.success = true;
        return result;
    }

    int total = 0;
    for (std::size_t i = 0; i < wallet.size(); ++i)
    {
        total += wallet[i] * denominations[i];
    }
    if (total < amount)
    {
        return result;
    }

    std::array<int, 6> used{};
    int remaining = amount;

    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        const int denom = denominations[i];
        const int take = std::min(wallet[i], remaining / denom);
        if (take > 0)
        {
            used[i] = take;
            remaining -= take * denom;
        }
    }

    while (remaining > 0)
    {
        int bestIndex = -1;
        int bestDenom = 0;

        for (std::size_t i = 0; i < denominations.size(); ++i)
        {
            if (wallet[i] - used[i] <= 0)
            {
                continue;
            }
            const int denom = denominations[i];
            if (denom >= remaining && (bestIndex == -1 || denom < bestDenom))
            {
                bestIndex = static_cast<int>(i);
                bestDenom = denom;
            }
        }

        if (bestIndex == -1)
        {
            for (int i = static_cast<int>(denominations.size()) - 1; i >= 0; --i)
            {
                if (wallet[static_cast<std::size_t>(i)] - used[static_cast<std::size_t>(i)] > 0)
                {
                    bestIndex = i;
                    bestDenom = denominations[static_cast<std::size_t>(i)];
                    break;
                }
            }
        }

        if (bestIndex == -1)
        {
            return result;
        }

        ++used[static_cast<std::size_t>(bestIndex)];
        remaining -= bestDenom;
    }

    int paid = 0;
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        paid += used[i] * denominations[i];
    }

    if (paid < amount)
    {
        return result;
    }

    int changeValue = paid - amount;
    std::array<int, 6> change{};
    for (std::size_t i = 0; i < denominations.size(); ++i)
    {
        const int denom = denominations[i];
        const int count = changeValue / denom;
        change[i] = count;
        changeValue -= count * denom;
    }

    result.paid = paid;
    result.used = used;
    result.change = change;
    result.success = true;
    return result;
}
