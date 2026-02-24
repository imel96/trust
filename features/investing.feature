Feature: Investing Activities
  # assume there's only one broker
  # assume no share price change
  # assume $1 commission

  Scenario: Buy shares

    Given there's a trust in the database
      And the trust has 101 dollar in "cash" account
      And the trust has 0 dollar in "cost" account
      And the trust has 0 dollar in "shares" account
      And there's a broker in the database
      And the company "BHP" is on the market
    When the trust "buys" 10 shares of "BHP" through the broker
    Then "BHP" shares show up in the trust portfolio
      And the trust "cash" account becomes 0 dollar
      And the trust "shares" account becomes 100 dollar
      And the trust "cost" account becomes 1 dollar
      And the total "Investing" cash flow will be -100 dollar

  Scenario: Sell shares

    Given there's a broker in the database
      And there's a trust in the database
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "cost" account
      And the trust has 100 dollar in "shares" account
      And the trust has 10 shares in the trust portfolio
      And the company "BHP" is on the market
    When the trust "sells" 10 shares of "BHP" through the broker
    Then the trust portfolio shares of "BHP" becomes 0
      And the trust "cash" account becomes 99 dollar
      And the trust "cost" account becomes 1 dollar
      And the trust "shares" account becomes 0 dollar
      And the total "Investing" cash flow will be 100 dollar
