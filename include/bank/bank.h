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
    int GetMoney() const { return money; }
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
    /**
     * @description: 直接设置银行余额（用于存档加载）
     * @param {int} amount 余额
     * @return null
     */
    void SetMoney(int amount) { money = amount; }

  private:
    int money;
};
