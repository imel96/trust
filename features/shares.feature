Feature: Shares management
  # assummed there's only one broker

  Scenario: Buy shares

    Given there's a trust in the database
      And the trust has 100 dollar in "cash" account
      And the trust has 0 dollar in "cost" account
      And there's a broker in the database
      And there's a company in the database
    When the trust "buys" 10 shares of the company through the broker
    Then the company shares show up in the trust portfolio
      And the trust's "cash" account become 0 dollar
      And the trust's "cost" account become 1 dollar

  Scenario: Sell shares

    Given there's a broker in the database
      And there's a trust in the database
      And the trust has 0 dollar in "cash" account
      And the trust has 0 dollar in "cost" account
      And the trust has 10 shares in the trust portfolio
      And there's a company in the database
    When the trust "sells" 10 shares of the company through the broker
    Then the trust portfolio shares of the company become "0"
      And the trust's "cash" account become 100 dollar
      And the trust's "cost" account become 1 dollar
