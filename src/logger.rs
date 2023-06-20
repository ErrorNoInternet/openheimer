use chrono::Local;
use colored::*;

pub enum LogMessageType {
    Verbose1 = 1,
    Verbose2 = 2,
    Verbose3 = 3,
    Information,
    Error,
    Warning,
}

pub struct Logger {
    pub verbosity: usize,
}
impl Logger {
    pub fn log_error(&self, error_description: &str, error: &impl std::fmt::Display) {
        self.log_message(
            LogMessageType::Error,
            &format!("{} ({})", &error_description, &error.to_string()),
        );
    }

    pub fn log_message(&self, message_type: LogMessageType, message: &str) {
        match message_type {
            LogMessageType::Information => {
                println!(
                    "{} {} {}",
                    current_time(),
                    colored_brackets(&"INFORMATION".bold().blue()),
                    message
                )
            }
            LogMessageType::Verbose1 | LogMessageType::Verbose2 | LogMessageType::Verbose3 => {
                if self.verbosity >= message_type as usize {
                    println!(
                        "{} {} {}",
                        current_time(),
                        colored_brackets(&"VERBOSE".bold().blue()),
                        message
                    )
                }
            }
            LogMessageType::Warning => {
                eprintln!(
                    "{} {} {}",
                    current_time(),
                    colored_brackets(&"WARNING".bold().yellow()),
                    message
                )
            }
            LogMessageType::Error => eprintln!(
                "{} {} {}",
                current_time(),
                colored_brackets(&"ERROR".bold().red()),
                message.red()
            ),
        }
    }
}

fn current_time() -> String {
    format!(
        "{}{}{}",
        "[".bold().white(),
        Local::now()
            .format("%Y/%m/%d %H:%M:%S")
            .to_string()
            .bold()
            .white(),
        "]".bold().white()
    )
}

fn colored_brackets(text: &ColoredString) -> String {
    format!("{}{}{}", "[".bold().blue(), text, "]".bold().blue())
}
