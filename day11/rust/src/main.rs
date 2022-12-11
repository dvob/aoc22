use std::collections::{BinaryHeap, VecDeque};

use nom::{
    branch::alt,
    bytes::complete::{tag, take_while},
    character::complete::{digit1, multispace0, multispace1},
    combinator::{map, map_res, recognize, value},
    multi::{many0, separated_list1},
    sequence::terminated,
    IResult,
};

#[derive(Debug, Clone)]
enum Operand {
    Old,
    Number(i64),
}

impl Operand {
    fn get(&self, old: i64) -> i64 {
        match self {
            Operand::Old => old,
            Operand::Number(i) => *i,
        }
    }
}

#[derive(Debug, Clone)]
enum Operator {
    Add,
    Multiply,
}

#[derive(Debug, Clone)]
struct Expr {
    op: Operator,
    lhs: Operand,
    rhs: Operand,
}

#[derive(Debug, Clone)]
struct Monkey {
    items: VecDeque<i64>,
    expr: Expr,
    test: i64,
    dest_true: i64,
    dest_false: i64,
}

fn parse_operand(s: &str) -> IResult<&str, Operand> {
    alt((
        map(tag("old"), |_| Operand::Old),
        map(map_res(recognize(digit1), str::parse::<i64>), |i| {
            Operand::Number(i)
        }),
    ))(s)
}

fn parse_op(s: &str) -> IResult<&str, Expr> {
    let (s, lhs) = parse_operand(s)?;
    let (s, _) = multispace1(s)?;
    let (s, op) = alt((
        value(Operator::Add, tag("+")),
        value(Operator::Multiply, tag("*")),
    ))(s)?;
    let (s, _) = multispace1(s)?;
    let (s, rhs) = parse_operand(s)?;
    Ok((s, Expr { op, lhs, rhs }))
}

fn parse_monkey(s: &str) -> IResult<&str, Monkey> {
    let (s, _) = tag("Monkey ")(s)?;
    let (s, _) = terminated(take_while(|c: char| c.is_digit(10)), tag(":"))(s)?;
    let (s, _) = multispace1(s)?;
    let (s, _) = tag("Starting items: ")(s)?;
    let (s, items) = separated_list1(tag(", "), map_res(recognize(digit1), str::parse::<i64>))(s)?;
    let items = VecDeque::from(items);
    let (s, _) = multispace1(s)?;
    let (s, _) = tag("Operation: new = ")(s)?;
    let (s, expr) = parse_op(s)?;
    let (s, _) = multispace1(s)?;
    let (s, _) = tag("Test: divisible by ")(s)?;
    let (s, test) = map_res(recognize(digit1), str::parse::<i64>)(s)?;

    let (s, _) = multispace1(s)?;
    let (s, _) = tag("If true: throw to monkey ")(s)?;
    let (s, dest_true) = map_res(recognize(digit1), str::parse::<i64>)(s)?;

    let (s, _) = multispace1(s)?;
    let (s, _) = tag("If false: throw to monkey ")(s)?;
    let (s, dest_false) = map_res(recognize(digit1), str::parse::<i64>)(s)?;
    Ok((
        s,
        Monkey {
            items,
            expr,
            test,
            dest_true,
            dest_false,
        },
    ))
}

fn parse(s: &str) -> IResult<&str, Vec<Monkey>> {
    separated_list1(multispace1, parse_monkey)(s)
}

fn calc(old: i64, expr: &Expr) -> i64 {
    let lhs = expr.lhs.get(old);
    let rhs = expr.rhs.get(old);
    match expr.op {
        Operator::Add => lhs + rhs,
        Operator::Multiply => lhs * rhs,
    }
}

fn solve(monkeys: &mut Vec<Monkey>) -> i64 {
    let mut inspections: Vec<i64> = Vec::with_capacity(monkeys.len());
    for _ in 0..monkeys.len() {
        inspections.push(0);
    }
    for _ in 0..20 {
        for i in 0..monkeys.len() {
            //let m = monkeys.get_mut(i).unwrap();

            while let Some(v) = monkeys.get_mut(i).unwrap().items.pop_front() {
                let v = calc(v, &monkeys.get(i).unwrap().expr);
                let v = v / 3;
                let dest = if v % monkeys.get(i).unwrap().test == 0 {
                    monkeys.get(i).unwrap().dest_true
                } else {
                    monkeys.get(i).unwrap().dest_false
                };
                monkeys
                    .get_mut(dest as usize)
                    .expect("invalid destination")
                    .items
                    .push_back(v);
                inspections[i] += 1
            }
        }
    }
    let mut top = inspections.into_iter().collect::<BinaryHeap<i64>>();
    top.pop().unwrap() * top.pop().unwrap()
}

fn calc_mod(old: i64, expr: &Expr, modulus: i64) -> i64 {
    let lhs = expr.lhs.get(old);
    let rhs = expr.rhs.get(old);
    match expr.op {
        Operator::Add => (lhs % modulus + rhs % modulus) % modulus,
        Operator::Multiply => ((lhs % modulus) * (rhs % modulus)) % modulus,
    }
}

fn solve1(monkeys: &mut Vec<Monkey>) -> i64 {
    let mut inspections: Vec<i64> = Vec::with_capacity(monkeys.len());
    let mut modulus: i64 = 1;
    for i in 0..monkeys.len() {
        inspections.push(0);
        modulus *= monkeys[i].test
    }
    for _ in 0..10_000 {
        for i in 0..monkeys.len() {
            while let Some(v) = monkeys.get_mut(i).unwrap().items.pop_front() {
                let v = calc_mod(v, &monkeys.get(i).unwrap().expr, modulus);

                let dest = if v % monkeys.get(i).unwrap().test == 0 {
                    monkeys.get(i).unwrap().dest_true
                } else {
                    monkeys.get(i).unwrap().dest_false
                };
                monkeys
                    .get_mut(dest as usize)
                    .expect("invalid destination")
                    .items
                    .push_back(v);
                inspections[i] += 1
            }
        }
    }
    let mut top = inspections.into_iter().collect::<BinaryHeap<i64>>();
    top.pop().unwrap() * top.pop().unwrap()
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let filename = std::env::args().nth(1).expect("missing argument: filename");
    let input = std::fs::read_to_string(filename)?;

    // TODO: why does this not work?
    // let (_, ms) = parse(input.as_str())?;
    let r = parse(input.as_str());
    let mut ms = match r {
        Ok((_, ms)) => ms,
        Err(e) => panic!("{e}"),
    };

    let result = solve(&mut ms.to_vec());
    println!("{result}");

    let result = solve1(&mut ms);
    println!("{result}");
    Ok(())
}
