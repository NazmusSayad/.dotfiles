use std::{
  env,
  process::{exit, Command},
};

const DIM: &str = "\x1b[2m";
const BLUE_DIM: &str = "\x1b[2;34m";
const RED_DIM: &str = "\x1b[2;31m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  let args: Vec<String> = env::args().skip(1).collect();

  let current_branch = Command::new("git")
    .args(["branch", "--show-current"])
    .output()
    .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
    .unwrap_or_default();

  let target_branch = match args.len() {
    0 => {
      println!("{}No branch specified, using current branch{}", DIM, NORMAL);
      current_branch.clone()
    }
    1 => args[0].clone(),
    _ => {
      eprintln!("Usage: gp [branch]");
      exit(1);
    }
  };

  print!("{}Pulling changes from ", DIM);
  print!("{}{}", BLUE_DIM, target_branch);
  print!("{} into ", DIM);
  print!("{}{}", RED_DIM, current_branch);
  println!("{} (default)", DIM);
  print!("{}", NORMAL);

  let status = Command::new("git")
    .args(["prune", "--progress"])
    .status()
    .expect("Failed to run git prune");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["pull", "origin", &target_branch, "--progress"])
    .status()
    .expect("Failed to run git pull");

  exit(status.code().unwrap_or(1));
}

