use std::{
  env,
  io::{self, Write},
  process::{exit, Command},
};

const RED: &str = "\x1b[31m";
const GREEN: &str = "\x1b[32m";
const YELLOW: &str = "\x1b[33m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  let args: Vec<String> = env::args().skip(1).collect();

  let branch = match args.get(0) {
    Some(v) if !v.is_empty() => v,
    _ => {
      eprintln!("{}❌ Branch name required{}", RED, NORMAL);
      exit(1);
    }
  };

  if branch.starts_with('-') {
    eprintln!("{}❌ Invalid branch name: {}{}", RED, branch, NORMAL);
    exit(1);
  }

  let remote = Command::new("git")
    .arg("remote")
    .output()
    .ok()
    .and_then(|o| {
      String::from_utf8_lossy(&o.stdout)
        .lines()
        .next()
        .map(|l| l.trim().to_string())
    })
    .unwrap_or_default();

  let head_ref = format!("refs/heads/{}", branch);

  let mut exists = Command::new("git")
    .args(["rev-parse", "--verify", "--quiet", &head_ref])
    .status()
    .map(|s| s.success())
    .unwrap_or(false);

  if !exists && !remote.is_empty() {
    let remote_ref = format!("refs/remotes/{}/{}", remote, branch);

    exists = Command::new("git")
      .args(["rev-parse", "--verify", "--quiet", &remote_ref])
      .status()
      .map(|s| s.success())
      .unwrap_or(false);
  }

  let mut cmd = Command::new("git");

  if exists {
    print!("{}", GREEN);
    io::stdout().flush().ok();
    cmd.arg("checkout");
  } else {
    print!("{}", YELLOW);
    io::stdout().flush().ok();
    cmd.args(["checkout", "-b"]);
  }

  for arg in args.iter() {
    cmd.arg(arg);
  }

  let status = cmd.status();
  print!("{}", NORMAL);

  match status {
    Ok(status) => {
      exit(status.code().unwrap_or(1));
    }
    Err(err) => {
      eprintln!("{}Failed to run git checkout: {}{}", RED, err, NORMAL);
      exit(1);
    }
  }
}

