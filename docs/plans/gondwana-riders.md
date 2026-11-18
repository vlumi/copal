# Plan: Gondwana's riders

Part 8 of [a deeper Gondwana](gondwana-depth.md), which this plan belongs to; its ground rules apply. A policy of Gondwana Life Insurance is one product with one sum, and only a medical policy pays a hospitalization. Life insurance is not sold that way: a policy is a main contract and the riders attached to it, each with its own amount, premium and term, and the hospitalization benefit is most often a rider on a whole life or an endowment policy. This plan gives a policy its coverages, in a file of their own beside the master.

## The coverage file

`COVMST`, an indexed file keyed by the policy number and a coverage number, holds a record of 80 bytes for each coverage of each policy. Coverage 01 is the main contract and every policy has it; 02 and up are riders. A coverage has its product, its state (in force, expired, cancelled, or ended with the policy), its amount, its premium for one payment in the policy's mode, its contract date, term and maturity date, the entry age its premium is rated at, and the date it ended.

The master keeps what it has. Its premium becomes the total of the premiums of the coverages in force, which is what billing, collection, lapse and reinstatement work with, so none of those programs changes. Its product, sum, term and maturity date stay those of the main contract, and repeat coverage 01: a daily hospitalization amount does not add to a death benefit, so a total of the amounts would mean nothing, and the cash value, the loan limit, the dividend and the reserve are the main contract's. Keeping both in step is the price of not rewriting every program that reads the master, and the kind of redundancy such a system has.

## Riders

Two riders, as products in the product table marked as riders by a byte taken from its filler:

- **TR, the term rider:** a death benefit on top of the main contract's, for ten years, renewed at the attained age until the insured is past 60.
- **HR, the hospitalization rider:** the hospitalization and surgery benefits of the medical policy, for ten years, renewed at the attained age until the insured is past 70, with the same waiting period.

A rider is rated like a policy, by its own lines of the rate table. It is attached when the policy is issued, so its anniversaries are the policy's; attaching one later needs underwriting the shop does elsewhere, and is not built. A rider may not be of the main contract's kind (no hospitalization rider on a medical policy, no term rider on a term policy), and a policy has two at most, of different products.

## What changes

- **Opening (GLO010):** loads the coverage extract after the policy extract, and counts a coverage of no policy and a coverage twice.
- **New business (GLD010):** the sales system's record carries up to two riders, a product and an amount each. Each is checked as the main contract is (a rider product, the age, the amount, a rate) and against the rules above; a policy with a rider refused is refused whole. The main contract and the riders are written as coverages, and the master's premium is their total.
- **Anniversaries (GLD020):** walks the policy's coverages. One that reaches its maturity date matures if it is an endowment's main contract, and otherwise is renewed at the attained age or expires, as the term policy is today; a rider that expires leaves the policy in force at a lower premium, and the main contract ending ends the riders with it.
- **Benefits (GLD030):** a death pays the amounts of the coverages that insure a death and were in force on the day of the event, so the term rider's too; a policy with none is refused as the medical policy is today. A hospitalization or a surgery is paid by the coverage of the medical kind in force on the day of the event, the medical policy's main contract or the rider, at its amount and from its contract date; a policy with none is refused. A surrender pays the main contract's cash value, as riders have none. A death and a surrender end every coverage.
- **Changes (GLD040):** a new request cancels a rider of a policy in force: the rider ends and the master's premium drops by its premium.
- **Unload (GLD090):** unloads the coverage file beside the master, for the downstream systems and for the generator, which needs to know who has which rider.
- **The batch files:** OPEN, DAILY and UNLOAD name the new files.

Billing, results, lapse, loans, reserves, commissions and dividends do not change. The reserve stays the main contract's until part 9 computes it by coverage.

## Around the programs

- **The fixture:** the March fixture gains a coverage extract, with riders on some of its policies, chosen so that the month renews a rider, expires one by age, ends riders at a maturity, pays a death with a term rider and a hospitalization from a rider, among them one for a stay before the rider was cancelled, and refuses what the rules refuse. Every golden file of it changes or is joined by one.
- **The generator:** a policy of the extract has riders as often as they are sold, renewed and expired as its history would have left them; applications carry riders; branches cancel one now and then; and hospitalizations come from whoever has the coverage.
- **The console:** a policy's page lists its coverages.
- **The survey:** the blueprints and the oracle's cases are regenerated. A file keyed by another file's key and a number is new to it, and what it makes of that goes in the findings.

## Not in this part

Attaching a rider to a policy in force; reducing a sum; a rider's own beneficiary; reserves and commissions by coverage (part 9); the coverages printed on a policyholder's notice (part 5).
