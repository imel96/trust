package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"log"
	_ "modernc.org/sqlite"
	"strconv"
	"strings"
	"time"
)

const path = "/Users/melbysjamsuddin/Documents/trust_test"
const debug = false

var Db *sql.DB

type brokerCtxKey struct{}
type personCtxKey struct{}
type trustCashAccountCtxKey struct{}
type trustCostAccountCtxKey struct{}
type trustCtxKey struct{}

func hashString(text string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
}

func queryRow(sql string) string {
	var out string
	if debug {
		log.Printf("%s", sql)
	}
	err := Db.QueryRow(sql).Scan(&out)
	if err != nil {
		return "error"
	}
	return out
}

func executeQueries(queries ...string) error {

	for _, sql := range queries {
		if debug {
			log.Printf("%s", sql)
		}
		_, err := Db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertAccount(account_number string, name string, balance int) {
	executeQueries(
		"INSERT OR IGNORE INTO account_hub(account_hash, account_number) values('"+hashString(account_number)+"', '"+account_number+"');",
		"INSERT INTO account_sat(hash, account_hash, name, balance) values('"+hashString(strings.Join([]string{account_number, name, fmt.Sprintf("%d", balance)}, "|"))+"', '"+hashString(account_number)+"', '"+name+"', "+fmt.Sprintf("%d", balance)+");",
	)
	time.Sleep(100 * time.Millisecond)
}

func insertCompany(code string, name string, gics string) {
	executeQueries(
		"INSERT INTO company_hub(company_hash, code) values('"+hashString(code)+"', '"+code+"');",
		"INSERT INTO company_sat values('"+hashString(strings.Join([]string{code, name, gics}, "|"))+"', '"+hashString(code)+"', '"+name+"', '"+gics+"');",
	)
}

func insertParty(name string, party_type string) {
	executeQueries(
		"INSERT INTO party_hub(party_hash, name) values('"+hashString(name)+"', '"+name+"');",
		"INSERT INTO party_sat(hash, party_hash, party_type) values('"+hashString(name+"|"+party_type)+"', '"+hashString(name)+"', '"+party_type+"');",
	)
}

func insertLinkWithTwoHubs(link string, a string, b string) {
	executeQueries("INSERT INTO " + link + " values('" + hashString(a+"|"+b) + "', '" + hashString(a) + "', '" + hashString(b) + "', STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'));")
}

func insertTrustCompanySat(trust string, company string, shareQuantity int) {
	executeQueries("INSERT INTO trust_company_sat(hash, trust_company_hash, share_quantity) values('" + hashString(strings.Join([]string{trust, company, fmt.Sprintf("%d", shareQuantity)}, "|")) + "', '" + hashString(trust+"|"+company) + "', " + fmt.Sprintf("%d", shareQuantity) + ");")
}

func insertPersonTrustSat(person string, trust string, numUnit int) {
	executeQueries("INSERT INTO person_trust_sat(hash, person_trust_hash, unit_quantity) values('" + hashString(strings.Join([]string{person, trust, fmt.Sprintf("%d", numUnit)}, "|")) + "', '" + hashString(person+"|"+trust) + "', " + fmt.Sprintf("%d", numUnit) + ");")
}

func theNAVPerUnitWillBeDollar(amount int) error {
	out := getTrustNAV("shamsuddins")
	if amount != out {
		return fmt.Errorf("expected response to be: %d, but actual is: %d", amount, out)
	}
	return nil
}

func theTrustPortfolioSharesOfTheCompanyBecome(shareQuantity string) error {
	out := queryRow("SELECT share_quantity FROM party_hub trust JOIN trust_company_link on (trust_hash = trust.party_hash) JOIN company_hub USING (company_hash) JOIN trust_company_sat USING(trust_company_hash) ORDER BY trust_company_sat.load_time DESC LIMIT 1;")

	if out != shareQuantity {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", shareQuantity, out)
	}
	return nil
}

func getTrustNAV(trust string) int {
	accountNumber := "1"
	out, _ := strconv.Atoi(queryRow("SELECT account_sat.balance FROM trust_account_link JOIN account_hub USING(account_hash) JOIN account_sat USING(account_hash) WHERE account_hash = '" + hashString(accountNumber) + "' ORDER BY account_sat.load_time DESC LIMIT 1;"))

	return out
}

func theTrustHasSharesInTheTrustPortfolio(numShares int) error {
	insertLinkWithTwoHubs("trust_company_link", "shamsuddins", "AZJ")
	insertTrustCompanySat("shamsuddins", "AZJ", numShares)
	return nil
}

func theTrustHasDollarInAccount(amount int, account string) error {
	switch account {
	case "cash":
		insertLinkWithTwoHubs("trust_account_link", "shamsuddins", "1")
		insertAccount("1", "shamsuddins cash", amount)
	case "cost":
		insertLinkWithTwoHubs("trust_account_link", "shamsuddins", "5")
		insertAccount("5", "shamsuddins cost", amount)
	}
	return nil
}

func theTrustSharesOfTheCompanyThroughTheBroker(trade string, numShares int) error {
	out, _ := strconv.Atoi(queryRow("SELECT share_quantity FROM party_hub trust JOIN trust_company_link on (trust_hash = trust.party_hash) JOIN company_hub USING (company_hash) JOIN trust_company_sat USING(trust_company_hash) ORDER BY trust_company_sat.load_time DESC LIMIT 1;"))

	insertLinkWithTwoHubs("trust_company_link", "shamsuddins", "AZJ")

	switch trade {
	case "sells":
		if out < numShares {
			return fmt.Errorf("trying to sell shares more than owned")
		}
		insertTrustCompanySat("shamsuddins", "AZJ", out-numShares)
		insertAccount("1", "shamsuddins cash", 100)
		insertAccount("5", "shamsuddins cost", 1)
	case "buys":
		insertTrustCompanySat("shamsuddins", "AZJ", out+numShares)
		insertAccount("1", "shamsuddins cash", 0)
		insertAccount("5", "shamsuddins cost", 1)
	}
	return nil
}

func theTrustsAccountBecomeDollar(account, amount string) error {
	var accountNumber string

	switch account {
	case "cash":
		accountNumber = "1"
	case "cost":
		accountNumber = "5"
	}
	out := queryRow("SELECT account_sat.balance FROM trust_account_link JOIN account_hub USING(account_hash) JOIN account_sat USING(account_hash) WHERE account_hash = '" + hashString(accountNumber) + "' ORDER BY account_sat.load_time DESC LIMIT 1;")

	if strings.TrimRight(out, "\n") != amount {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", amount, out)
	}
	return nil
}

func theTrustWillHaveUnitOutstanding(numUnits int) error {
	unitCount, _ := strconv.Atoi(queryRow("SELECT SUM(unit_quantity) FROM party_hub person JOIN person_trust_link ON (person_hash = person.party_hash) JOIN party_hub trust ON (trust_hash = trust.party_hash) JOIN (SELECT person_trust_hash, unit_quantity, load_time, RANK() OVER (PARTITION BY person_trust_hash ORDER BY load_time DESC) rank FROM person_trust_sat) ranked_person_trust_sat USING(person_trust_hash) WHERE rank = 1;"))

	if unitCount != numUnits {
		return fmt.Errorf("expected response to be: %d, but actual is: %d", numUnits, unitCount)
	}
	return nil
}

func theCompanySharesShowUpInTheTrustPortfolio() error {
	out := queryRow("SELECT trust.name || \"|\" || company_hub.code FROM party_hub trust JOIN trust_company_link on (trust_hash = trust.party_hash) JOIN company_hub USING (company_hash);")

	if strings.TrimRight(out, "\n") != "shamsuddins|AZJ" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "shamsuddins|AZJ", out)
	}
	return nil
}

func theresACompanyInTheDatabase() error {
	insertCompany("AZJ", "Aurizon Holdings", "Transportation")
	return nil
}

func thePersonCanBeSeenInTheTrustMemberList(ctx context.Context) error {
	out := queryRow("SELECT person.name || \"|\" || trust.name FROM party_hub person JOIN person_trust_link on (person_hash = person.party_hash) JOIN party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "melby|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "melbys|shamsuddins", out)
	}
	return nil
}

func thePersonDollar(action string, amount int) error {
	out, _ := strconv.Atoi(queryRow("SELECT unit_quantity FROM party_hub trust JOIN person_trust_link on (trust_hash = trust.party_hash) JOIN party_hub person ON (person_hash = person.party_hash) JOIN person_trust_sat USING(person_trust_hash) ORDER BY person_trust_sat.load_time DESC LIMIT 1;"))

	insertLinkWithTwoHubs("person_trust_link", "melby", "shamsuddins")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from integer divide by zero panic", r)
		}
	}()
	log.Printf("amount: %d nav: %d out: %d", amount, getTrustNAV("shamsuddins"), out)
	numShares := 1

	if getTrustNAV("shamsuddins") != 0 {
		numShares = amount / getTrustNAV("shamsuddins")
	}

	switch action {
	case "invest":
		insertPersonTrustSat("melby", "shamsuddins", out+numShares)
		insertAccount("1", "shamsuddins cash", amount)
	case "redeem":
		if out < numShares {
			return fmt.Errorf("trying to sell shares more than owned")
		}
		insertPersonTrustSat("melby", "shamsuddins", out-numShares)
		insertAccount("1", "shamsuddins cash", 0)
	}
	return nil
}

func thePersonIsAddedToTheTrust(ctx context.Context) error {
	insertLinkWithTwoHubs(
		"person_trust_link",
		fmt.Sprintf("%v", ctx.Value(personCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func thePersonOwnsUnitInTheTrust(numUnits int) error {
	insertLinkWithTwoHubs("person_trust_link", "melby", "shamsuddins")
	insertPersonTrustSat("melby", "shamsuddins", numUnits)
	return nil
}

func thePersonWillHaveUnitInTheTrust(numUnits int) error {
	out, _ := strconv.Atoi(queryRow("SELECT unit_quantity FROM party_hub person JOIN person_trust_link ON (person_hash = person.party_hash) JOIN party_hub trust ON (trust_hash = trust.party_hash) JOIN person_trust_sat USING(person_trust_hash) ORDER BY person_trust_sat.load_time DESC LIMIT 1;"))

	if out != numUnits {
		return fmt.Errorf("expected response to be: %d, but actual is: %d", numUnits, out)
	}
	return nil
}

func theresAPersonInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("melby", "trustee")
	return context.WithValue(ctx, personCtxKey{}, "melby"), nil
}

func theresAPersonInTheTrust(ctx context.Context) error {
	_, _ = theresAPersonInTheDatabase(ctx)
	return thePersonIsAddedToTheTrust(ctx)
}

func theresATrustInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("shamsuddins", "trust")
	return context.WithValue(ctx, trustCtxKey{}, "shamsuddins"), nil
}

func theBrokerCanBeSeenInTheTrustBrokerList(ctx context.Context) error {
	out := queryRow("SELECT broker.name || \"|\" || trust.name FROM party_hub broker JOIN broker_trust_link on (broker_hash = broker.party_hash) JOIN party_hub trust on (trust_hash = trust.party_hash);")

	if strings.TrimRight(out, "\n") != "cmcmarkets|shamsuddins" {
		return fmt.Errorf("expected response to be: %s, but actual is: %s", "cmcmarkets|shamsuddins", out)
	}
	return nil
}

func theBrokerIsAddedToTheTrust(ctx context.Context) error {
	insertLinkWithTwoHubs(
		"broker_trust_link",
		fmt.Sprintf("%v", ctx.Value(brokerCtxKey{})),
		fmt.Sprintf("%v", ctx.Value(trustCtxKey{})),
	)
	return nil
}

func theresABrokerInTheDatabase(ctx context.Context) (context.Context, error) {
	insertParty("cmcmarkets", "cmcmarkets")
	return context.WithValue(ctx, brokerCtxKey{}, "cmcmarkets"), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		executeQueries(
			"DELETE FROM account_hub;",
			"DELETE FROM account_sat;",
			"DELETE FROM broker_trust_link;",
			"DELETE FROM company_hub;",
			"DELETE FROM company_sat;",
			"DELETE FROM party_hub;",
			"DELETE FROM party_sat;",
			"DELETE FROM person_trust_link;",
			"DELETE FROM person_trust_sat;",
			"DELETE FROM trust_account_link;",
			"DELETE FROM trust_company_link;",
			"DELETE FROM trust_company_sat;",
		)
		return ctx, nil
	})
	ctx.Step(`^the broker can be seen in the trust broker list$`, theBrokerCanBeSeenInTheTrustBrokerList)
	ctx.Step(`^the broker is added to the trust$`, theBrokerIsAddedToTheTrust)
	ctx.Step(`^the company shares show up in the trust portfolio$`, theCompanySharesShowUpInTheTrustPortfolio)
	ctx.Step(`^the NAV per unit will be (\d+) dollar$`, theNAVPerUnitWillBeDollar)
	ctx.Step(`^the person can be seen in the trust member list$`, thePersonCanBeSeenInTheTrustMemberList)
	ctx.Step(`^the person "([^"]*)" (\d+) dollar$`, thePersonDollar)
	ctx.Step(`^the person is added to the trust$`, thePersonIsAddedToTheTrust)
	ctx.Step(`^the person owns (\d+) unit in the trust$`, thePersonOwnsUnitInTheTrust)
	ctx.Step(`^the person will have (\d+) unit in the trust$`, thePersonWillHaveUnitInTheTrust)
	ctx.Step(`^the trust "([^"]*)" (\d+) shares of the company through the broker$`, theTrustSharesOfTheCompanyThroughTheBroker)
	ctx.Step(`^the trust has (\d+) dollar in "([^"]*)" account$`, theTrustHasDollarInAccount)
	ctx.Step(`^the trust has (\d+) shares in the trust portfolio$`, theTrustHasSharesInTheTrustPortfolio)
	ctx.Step(`^the trust portfolio shares of the company become "([^"]*)"$`, theTrustPortfolioSharesOfTheCompanyBecome)
	ctx.Step(`^the trust will have (\d+) unit outstanding$`, theTrustWillHaveUnitOutstanding)
	ctx.Step(`^the trust\'s "([^"]*)" account become (\d+) dollar$`, theTrustsAccountBecomeDollar)
	ctx.Step(`^there\'s a broker in the database$`, theresABrokerInTheDatabase)
	ctx.Step(`^there\'s a company in the database$`, theresACompanyInTheDatabase)
	ctx.Step(`^there\'s a person in the database$`, theresAPersonInTheDatabase)
	ctx.Step(`^there\'s a person in the trust$`, theresAPersonInTheTrust)
	ctx.Step(`^there\'s a trust in the database$`, theresATrustInTheDatabase)
}

func FeatureContext(s *godog.TestSuiteContext) {
	s.BeforeSuite(func() {
		var err error
		Db, err = sql.Open("sqlite", path)

		if err != nil {
			log.Fatal(err)
		}
		executeQueries(
			"CREATE TABLE IF NOT EXISTS account_hub(account_hash text unique, account_number text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS account_sat(hash text unique, account_hash text, name text, balance int, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS broker_trust_link(broker_trust_hash text unique, broker_hash text, trust_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS company_hub(company_hash text unique, code text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS company_sat(hash text unique, company_hash text, name text, gics text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS party_hub(party_hash text unique, name text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')) NOT NULL);",
			"CREATE TABLE IF NOT EXISTS party_sat(hash text unique, party_hash text, party_type text, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE IF NOT EXISTS person_trust_link(person_trust_hash text unique, person_hash text, trust_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_account_link(trust_account_hash text unique, trust_hash text, account_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_company_link(trust_company_hash text unique, trust_hash text, company_hash text, load_time timestamp);",
			"CREATE TABLE IF NOT EXISTS trust_company_sat(hash text unique, trust_company_hash text, share_quantity int, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
			"CREATE TABLE person_trust_sat(hash text unique, person_trust_hash text, unit_quantity int, load_time timestamp DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));",
		)
	})
}
