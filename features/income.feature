Feature: Income management
  # assummed there's only one broker

  Scenario: Dividend income

    Given there's a trust in the database
      And the trust has 10 shares in the trust portfolio
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "income" account
      And the company "BHP" is on the market
    When the company pays dividend of 0.1 dollar per share
    Then the trust's "income" account become 1 dollar
      And the trust's "cash" account become 1 dollar

  Scenario: Capital gain

    Given there's a trust in the database
      And the trust has 10 shares in the trust portfolio
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "income" account
      And the company "BHP" is on the market
    When the company pays dividend of 0.1 dollar per share
    Then the trust's "income" account become 1 dollar

