use std::{
  io::{self, Read, Write},
  process::{exit, Command},
};

const DIM: &str = "\x1b[2m";
const DIM_RED: &str = "\x1b[2;31m";
const BOLD_RED: &str = "\x1b[1;31m";
const GREEN: &str = "\x1b[32m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  let current = Command::new("git")
    .args(["branch", "--show-current"])
    .output()
    .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
    .unwrap_or_default();

  let output = Command::new("git")
    .args(["branch", "--format=%(refname:short)"])
    .output()
    .expect("Failed to list branches");

  let branches: Vec<String> = String::from_utf8_lossy(&output.stdout)
    .lines()
    .map(|l| l.trim().to_string())
    .filter(|l| !l.is_empty() && *l != current)
    .collect();

  if branches.is_empty() {
    println!("No other branches to delete");
    return;
  }

  print!("Branches to delete: ");
  println!("{}{}{}", BOLD_RED, branches.join(", "), NORMAL);

  print!("{}Press [Enter] to confirm, or any other key to cancel: {}", DIM, NORMAL);
  io::stdout().flush().ok();

  let mut buf = [0u8; 1];
  let read = io::stdin().read(&mut buf).unwrap_or(0);
  let confirm = read == 0 || buf[0] == b'\n' || buf[0] == b'\r';

  if !confirm {
    println!("{}Cancelled branch deletion{}", GREEN, NORMAL);
    return;
  }

  print!("{}", DIM_RED);
  let status = Command::new("git").args(["prune", "--progress"]).status();
  print!("{}", NORMAL);

  match status {
    Ok(status) => {
      if !status.success() {
        exit(status.code().unwrap_or(1));
      }
    }
    Err(err) => {
      eprintln!("{}Failed to run git prune: {}{}", BOLD_RED, err, NORMAL);
      exit(1);
    }
  }

  print!("{}", DIM_RED);
  let mut cmd = Command::new("git");
  cmd.arg("branch").arg("-D");

  for branch in branches.iter() {
    cmd.arg(branch);
  }

  let status = cmd.status();
  print!("{}", NORMAL);

  match status {
    Ok(status) => {
      exit(status.code().unwrap_or(1));
    }
    Err(err) => {
      eprintln!("{}Failed to run git branch -D: {}{}", BOLD_RED, err, NORMAL);
      exit(1);
    }
  }
}

