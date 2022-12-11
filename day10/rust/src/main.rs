use std::{ops::Index, str::{Lines, FromStr}, path::Iter, iter::Peekable, num::ParseIntError};

type MyResult<T> = Result<T, Box<dyn std::error::Error>>;

#[derive(Debug)]
enum Op {
    Noop,
    Addx(i64),
}

fn read_input(input: &str) -> MyResult<Vec<Op>> {
    let mut ops = Vec::new();
    for line in input.lines() {
        match &line[0..4] {
            "noop" => ops.push(Op::Noop),
            "addx" => ops.push(Op::Addx(line[5..].parse()?)),
            _ => return Err("failed to read line".into()),
        }
    }
    Ok(ops)
}

struct CycleIterator {
    add: Option<i64>,
    cycle: i64,
    value: i64,
    op_index: usize,
    // TODO: use iterator instead of Vec
    input: Vec<Op>
}

impl CycleIterator {
    fn new(ops: Vec<Op>) -> Self {
        Self {
            add: None,
            cycle: 0,
            value: 1,
            op_index: 0,
            input: ops,
        }
    }
}

impl Iterator for CycleIterator {
    type Item = (i64, i64);

    fn next(&mut self) -> Option<Self::Item> {
        self.cycle += 1;
        if let Some(add) = self.add {
            let old_value = self.value;
            self.value += add;
            self.add = None;
            return Some((self.cycle, old_value))
        }
        let result = match self.input.get(self.op_index) {
            Some(op) => match op {
                Op::Noop => Some((self.cycle, self.value)),
                Op::Addx(add) => {
                    self.add = Some(*add);
                    Some((self.cycle, self.value))
                },
            },
            None => None,
        };
        self.op_index+= 1;
        result
    }
}

fn solve(ops: Vec<Op>) -> i64 {
    const CYCLES: [i64; 6] = [20, 60, 100, 140, 180, 220];

    let mut signals = Vec::new();
    let it = CycleIterator::new(ops);

    for (cycle, value) in it {
        if CYCLES.contains(&cycle) {
            signals.push(value * cycle);
        }
    }

    signals.iter().sum()
}

fn solve2(ops: Vec<Op>) {
    let it = CycleIterator::new(ops);

    for (cycle, value) in it {
        let pos = (cycle - 1) % 40;

        if pos >= value - 1 && pos <= value + 1 {
            print!("#");
        } else {
            print!(".");
        }

        if pos == 39 {
            println!();
        }

    }
}


fn main() -> MyResult<()> {
    let filename = std::env::args().skip(1).next().expect("missing argument: filename");
    let input = std::fs::read_to_string(filename)?;
    let ops = read_input(input.as_str())?;

    let result = solve(ops);
    println!("{result}");

    let ops = read_input(input.as_str())?;
    solve2(ops);

    Ok(())
}