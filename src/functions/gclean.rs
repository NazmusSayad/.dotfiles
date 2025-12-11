use std::{
  io::{self, Write},
  process::{exit, Command},
};

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

  println!("Branches to delete: {}", branches.join(", "));
  print!("Press [Enter] to confirm, or any other key to cancel: ");
  io::stdout().flush().ok();

  let mut input = String::new();
  io::stdin().read_line(&mut input).ok();
  let confirm = input.trim_end_matches(|c| c == '\n' || c == '\r');

  if !confirm.is_empty() {
    println!("Cancelled branch deletion");
    return;
  }

  let status = Command::new("git")
    .args(["prune", "--progress"])
    .status()
    .expect("Failed to run git prune");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let mut cmd = Command::new("git");
  cmd.arg("branch").arg("-D");

  for branch in branches.iter() {
    cmd.arg(branch);
  }

  let status = cmd.status().expect("Failed to run git branch -D");

  exit(status.code().unwrap_or(1));
}

