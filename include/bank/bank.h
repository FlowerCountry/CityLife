/*
 * @Author: FlowerCity admin@flowercity.xyz
 * @Date: 2024-07-26 10:16:44
 * @LastEditors: FlowerCity admin@flowercity.xyz
 * @LastEditTime: 2024-07-28 09:41:53
 * @FilePath: \CityLive\include\bank\bank.h
 */
#pragma once

class Bank {
  public:
    /**
     * @description: 银行构造函数
     * @return null
     */
    Bank() : money{0} {}
    /**
     * @description: 银行析构函数
     * @return null
     */
    ~Bank() {}

    /**
     * @description: 获取银行余额
     * @return {*} 银行余额
     */
    int GetMoney() { return money; }
    /**
     * @description: 存钱
     * @param {int} money 存的数额
     * @return null
     */
    void DepositingMoney(const int &money) { this->money += money; }
    /**
     * @description: 取钱
     * @param {int} money
     * @return null
     */
    void WithdrawMoney(const int &money) { this->money -= money; }

  private:
    int money;
};
