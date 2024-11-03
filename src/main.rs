use clap::{arg, Parser};
use colored::Colorize;
use winproxy::{DefaultConnectionSettings, Flags};

#[derive(Parser, Debug)]
struct Args {
    /// Use a proxy server
    #[arg(short = 'p', long, value_name = "BOOLEAN")]
    use_proxy: Option<bool>,

    /// Use setup script
    #[arg(short = 's', long, value_name = "BOOLEAN")]
    use_script: Option<bool>,

    /// Automatically detect settings
    #[arg(short = 'a', long, value_name = "BOOLEAN")]
    auto_detect: Option<bool>,

    /// Proxy address
    #[arg(long, value_name = "ADDRESS")]
    proxy_address: Option<String>,

    /// Script address
    #[arg(long, value_name = "ADDRESS")]
    script_address: Option<String>,

    /// Bypass list
    #[arg(long, value_name = "ADDRESS")]
    bypass_list: Option<Vec<String>>,
}

impl Args {
    fn all_none(&self) -> bool {
        matches!(self, Args {
            use_proxy: None,
            use_script: None,
            auto_detect: None,
            proxy_address: None,
            script_address: None,
            bypass_list: None
        })
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
    let args = Args::parse();

    let mut settings = DefaultConnectionSettings::from_registry()
        .expect("failed to read default connection settings from registry");

    if args.all_none() {
        show_settings(&settings);
    } else {
        args.write_settings(&mut settings);
        settings.version += 1;
        settings
            .write_registry()
            .expect("failed to write default connection settings to registry");
    }
}

fn show_settings(settings: &DefaultConnectionSettings) {
    println!(
        "{} {}",
        "Use Proxy:".green(),
        format!("{}", settings.flags.contains(Flags::Proxy)).bright_purple()
    );
    println!(
        "{} {}",
        "Proxy Address:".green(),
        settings.proxy_address.bright_blue()
    );
    println!(
        "{} {}",
        "Use Script:".green(),
        format!("{}", settings.flags.contains(Flags::AutoProxyURL)).bright_purple()
    );
    println!(
        "{} {}",
        "Script Address:".green(),
        settings.script_address.bright_blue()
    );
    println!(
        "{} {}",
        "Auto-detect:".green(),
        format!("{}", settings.flags.contains(Flags::AutoDetect)).bright_purple()
    );
    println!(
        "{} {}\n{}\n{}", "Bypass List:".green(),
        "[".bright_black(),
        {
            settings
                .bypass_list
                .iter()
                .map(|address| format!("  {}", address.bright_blue()))
                .collect::<Vec<_>>()
                .join(",\n")
                .bright_black()
        },
        "]".bright_black()
    );
}
