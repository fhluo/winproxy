mod i18n;

use clap::{CommandFactory, FromArgMatches, Parser};
use tabled::builder::Builder;
use tabled::settings::object::{Columns, Rows};
use tabled::settings::style::BorderColor;
use tabled::settings::themes::BorderCorrection;
use tabled::settings::{Alignment, Color, Panel, Style};
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

fn print_settings_table(settings: &DefaultConnectionSettings) {
    let data = [
        [
            i18n::t("use-proxy-title"),
            status_value(settings.is_proxy_enabled()),
        ],
        [
            i18n::t("use-script-title"),
            status_value(settings.is_script_enabled()),
        ],
        [
            i18n::t("auto-detect-title"),
            status_value(settings.is_auto_detect_enabled()),
        ],
        [
            i18n::t("proxy-address-title"),
            value_or_placeholder(&settings.proxy_address),
        ],
        [
            i18n::t("script-address-title"),
            value_or_placeholder(&settings.script_address),
        ],
    ];

    let mut table = Builder::from_iter(data).build();

    let palette = Palette::new();
    table
        .with(Style::rounded().remove_horizontals())
        .with(BorderColor::filled(palette.border.clone()))
        .modify(Columns::first(), Alignment::left())
        .modify(
            Columns::first(),
            BorderColor::new().right(palette.border.clone()),
        )
        .modify(Columns::first(), palette.label.clone() | Color::BOLD)
        .modify(Columns::last(), Alignment::center());

    let statuses = [
        (0, settings.is_proxy_enabled()),
        (1, settings.is_script_enabled()),
        (2, settings.is_auto_detect_enabled()),
    ];
    for (row, enabled) in statuses {
        table.modify(
            (row, 1),
            if enabled {
                palette.on.clone()
            } else {
                palette.off.clone()
            },
        );
    }

    if settings.proxy_address.is_empty() {
        table.modify((3, 1), palette.placeholder.clone());
    }
    if settings.script_address.is_empty() {
        table.modify((4, 1), palette.placeholder.clone());
    }

    let mut builder = Builder::default();
    if settings.bypass_list.is_empty() {
        builder.push_record(["—"]);
    } else {
        for item in &settings.bypass_list {
            builder.push_record([item.clone()]);
        }
    }

    let mut bypass_table = builder.build();
    bypass_table
        .with(Style::rounded())
        .with(Panel::header(i18n::t("bypass-list-title")))
        .with(BorderColor::filled(palette.border.clone()))
        .with(BorderCorrection::span())
        .modify(Rows::first(), Alignment::center())
        .modify(
            Rows::first(),
            BorderColor::new().bottom(palette.border.clone()),
        )
        .modify(Rows::first(), palette.label.clone() | Color::BOLD)
        .modify(Rows::new(1..), Alignment::center());

    if settings.bypass_list.is_empty() {
        bypass_table.modify((1, 0), palette.placeholder.clone());
    }

    anstream::println!("{table}\n{bypass_table}");
}

struct Palette {
    border: Color,
    label: Color,
    on: Color,
    off: Color,
    placeholder: Color,
}

impl Palette {
    fn new() -> Self {
        Self {
            border: Color::rgb_fg(0x4A, 0x47, 0x66),
            label: Color::rgb_fg(0xD8, 0xA7, 0x6A),
            on: Color::rgb_fg(0x3F, 0xB9, 0x50),
            off: Color::rgb_fg(0x8B, 0x94, 0x9E),
            placeholder: Color::rgb_fg(0x8B, 0x94, 0x9E),
        }
    }
}

fn status_value(on: bool) -> String {
    if on { "✓" } else { "✗" }.to_string()
}

fn value_or_placeholder(value: impl Into<String>) -> String {
    let value = value.into();

    if value.is_empty() {
        "—".to_string()
    } else {
        value
    }
}
