Feature: Trust Unit Management

  Scenario: Member invests unit

    Given there's a trust in the database
      And the trust has 0 dollar in "cash" account
      And there's a person in the trust
      And the person owns 0 unit in the trust
    When the person "invest" 1000 dollar
    Then the trust's "cash" account become 1000 dollar
      And the person will have 1 unit in the trust
      And the NAV per unit will be 1000 dollar
      And the trust will have 1 unit outstanding

  Scenario: Member redeem unit

    Given there's a trust in the database
      And the trust has 1000 dollar in "cash" account
      And there's a person in the trust
      And the person owns 1 unit in the trust
    When the person "redeem" 1000 dollar
    Then the trust's "cash" account become 0 dollar
      And the person will have 0 unit in the trust
      And the NAV per unit will be 0 dollar
      And the trust will have 0 unit outstanding
