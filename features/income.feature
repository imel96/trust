Feature: Income management
  # assummed there's only one broker

  Scenario: Dividend income

    Given there's a trust in the database
      And the trust has 10 shares in the trust portfolio
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "income" account
      And there's a company in the database
    When the company pays dividend of 10 cents per share
    Then the trust's "income" account become 1 dollar
      And the trust's "cash" account become 1 dollar

  Scenario: Capital gain

    Given there's a trust in the database
      And the trust has 10 shares in the trust portfolio
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "income" account
      And there's a company in the database
    When the company pays dividend of 10 cents per share
    Then the trust's "income" account become 1 dollar

