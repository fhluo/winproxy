mod i18n;

use clap::{CommandFactory, FromArgMatches, Parser};
use comfy_table::modifiers::UTF8_ROUND_CORNERS;
use comfy_table::presets::UTF8_FULL_CONDENSED;
use comfy_table::{Cell, CellAlignment, Color, ContentArrangement, Table};
use winproxy::{DefaultConnectionSettings, Flags};

#[derive(Parser, Debug, Default, PartialEq)]
struct Args {
    /// Use a proxy server
    #[arg(short = 'p', long, value_name = "BOOL", num_args(0..=1), require_equals(true), default_missing_value = "true", hide_possible_values = true)]
    use_proxy: Option<bool>,

    /// Use setup script
    #[arg(short = 's', long, value_name = "BOOL", num_args(0..=1), require_equals(true), default_missing_value = "true", hide_possible_values = true)]
    use_script: Option<bool>,

    /// Automatically detect settings
    #[arg(short = 'a', long, value_name = "BOOL", num_args(0..=1), require_equals(true), default_missing_value = "true", hide_possible_values = true)]
    auto_detect: Option<bool>,

    /// Proxy address
    #[arg(long, value_name = "ADDRESS")]
    proxy_address: Option<String>,

    /// Script address
    #[arg(long, value_name = "ADDRESS")]
    script_address: Option<String>,

    /// Bypass list (semicolon-separated)
    #[arg(long, value_name = "ADDRESS", value_delimiter = ';', num_args = 1..)]
    bypass_list: Option<Vec<String>>,
}

impl Args {
    fn has_changes(&self) -> bool {
        self != &Args::default()
    }

    fn write_settings(self, settings: &mut DefaultConnectionSettings) {
        if let Some(use_proxy) = self.use_proxy {
            settings.flags.set(Flags::Proxy, use_proxy)
        }

        if let Some(use_script) = self.use_script {
            settings.flags.set(Flags::AutoProxyURL, use_script)
        }

        if let Some(auto_detect) = self.auto_detect {
            settings.flags.set(Flags::AutoDetect, auto_detect)
        }

        if let Some(proxy_address) = self.proxy_address {
            settings.proxy_address = proxy_address;
        }

        if let Some(script_address) = self.script_address {
            settings.script_address = script_address;
        }

        if let Some(bypass_list) = self.bypass_list {
            settings.bypass_list = bypass_list;
        }
    }
}

fn main() {
    let args = parse_args();

    let mut settings = DefaultConnectionSettings::from_registry().unwrap_or_else(|err| {
        eprintln!("{}: {err}", i18n::t("read-registry-error"));
        std::process::exit(1);
    });

    if args.has_changes() {
        args.write_settings(&mut settings);
        settings.version += 1;
        if let Err(err) = settings.write_registry() {
            eprintln!("{}: {err}", i18n::t("write-registry-error"));
            std::process::exit(1);
        }
        return;
    }

    print_settings_table(&settings);
}

fn parse_args() -> Args {
    let command = i18n::localize_command(Args::command());
    let matches = command.get_matches();
    Args::from_arg_matches(&matches).unwrap_or_else(|err| err.exit())
}

fn symbol_cell(b: bool) -> Cell {
    if b {
        Cell::new("[x]").fg(Color::White)
    } else {
        Cell::new("[ ]").fg(Color::DarkGrey)
    }
}

fn title_cell(title: impl AsRef<str>) -> Cell {
    Cell::new(title.as_ref()).fg(Color::Green)
}

fn print_settings_table(settings: &DefaultConnectionSettings) {
    let mut table = Table::new();

    table
        .load_preset(UTF8_FULL_CONDENSED)
        .apply_modifier(UTF8_ROUND_CORNERS)
        .set_content_arrangement(ContentArrangement::Dynamic)
        .add_row(vec![
            title_cell(i18n::t("use-proxy-title")),
            symbol_cell(settings.is_proxy_enabled()),
        ])
        .add_row(vec![
            title_cell(i18n::t("use-script-title")),
            symbol_cell(settings.is_script_enabled()),
        ])
        .add_row(vec![
            title_cell(i18n::t("auto-detect-title")),
            symbol_cell(settings.is_auto_detect_enabled()),
        ])
        .add_row(vec![
            title_cell(i18n::t("proxy-address-title")),
            Cell::new(&settings.proxy_address).fg(Color::Blue),
        ])
        .add_row(vec![
            title_cell(i18n::t("script-address-title")),
            Cell::new(&settings.script_address).fg(Color::Blue),
        ])
        .add_row(vec![
            title_cell(i18n::t("bypass-list-title")),
            Cell::new(settings.bypass_list.join("\n")),
        ]);

    if let Some(col) = table.column_mut(0) {
        col.set_cell_alignment(CellAlignment::Left);
    }
    if let Some(col) = table.column_mut(1) {
        col.set_cell_alignment(CellAlignment::Center);
    }

    println!("{table}");
}
