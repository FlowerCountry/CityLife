/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 10:16:44
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 09:41:53
 * @FilePath: \CityLive\include\bank\bank.h
 */
#pragma once
#ifndef BANK_H
#define BANK_H

class Bank {
  public:
    Bank() : money{0} {}
    ~Bank() {}
    int GetMoney() { return money; }
    void DepositingMoney(int money) { this->money += money; }
    int WithdrawMoney(int money)
    {
        this->money -= money;
        return money;
    }

  private:
    int money;
};

#endif