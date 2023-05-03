use chrono::Local;
use colored::*;

pub enum LogMessageType {
    Information,
    Error,
    Warning,
}

pub fn log_error<T, E: std::fmt::Display>(result: Result<T, E>) {
    match result {
        Ok(_) => (),
        Err(error) => log_message(LogMessageType::Error, &error.to_string()),
    }
}

pub fn log_message(message_type: LogMessageType, message: &String) {
    match message_type {
        LogMessageType::Information => {
            println!(
                "{} {} {}",
                current_time(),
                colored_brackets(&"INFORMATION".bold().blue()),
                message
            )
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
