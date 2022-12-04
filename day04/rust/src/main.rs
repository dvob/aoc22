use std::error::Error;

type MyResult<T> = Result<T, Box<dyn Error>>;

#[derive(Clone,Copy)]
struct Range {
    from: u32,
    to: u32,
}

impl Range {
    fn contains(&self, r: Self) -> bool {
        self.from <= r.from && self.to >= r.to
    }

    fn contains_num(&self, n: u32) -> bool {
        n >= self.from && n <= self.to
    }

    fn overlap(&self, r: Self) -> bool {
        r.contains(*self) ||
        self.contains_num(r.from) ||
        self.contains_num(r.to)
    }
}

fn parse_range(input: &str) -> MyResult<Range>{
    match input.split_once('-') {
        Some((from, to)) => {
            Ok(Range {
                from: from.parse()?,
                to: to.parse()?
            })
        }
        None => {
            Err("invalid input".into())
        }
    }
}

fn parse_pair(input: &str) -> MyResult<(Range, Range)> {
    match input.split_once(',') {
        Some((left, right)) => {
            Ok((
                parse_range(left)?,
                parse_range(right)?,
            ))
        }
        None => {
            Err("invalid input".into())
        }
    }
}

fn parse(input: &str) -> MyResult<Vec<(Range,Range)>> {
    let mut pairs = Vec::new();
    for line in input.lines() {
        pairs.push(parse_pair(line)?)
    }
    Ok(pairs)
}

fn main() -> MyResult<()> {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 2 {
        eprintln!("missing argument: filename")
    }

    let filename = &args[1];

    let input = std::fs::read_to_string(filename)?;

    let pairs = parse(input.as_str())?;

    let mut result1 = 0;
    let mut result2 = 0;

    for (a, b) in pairs {
        if a.contains(b) || b.contains(a) {
            result1 += 1;
        }
        if a.overlap(b) {
            result2 += 1;
        }
    }

    println!("{}", result1);
    println!("{}", result2);
    Ok(())
}
