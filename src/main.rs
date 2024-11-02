use clap::{arg, Parser};
use winproxy::{DefaultConnectionSettings, Flags};

#[derive(Parser, Debug)]
struct Args {
    #[arg(short = 'p', long)]
    use_proxy: Option<bool>,

    #[arg(short = 's', long)]
    use_script: Option<bool>,

    #[arg(short = 'a', long)]
    auto_detect: Option<bool>,

    #[arg(long)]
    proxy_address: Option<String>,

    #[arg(long)]
    script_address: Option<String>,

    #[arg(long)]
    bypass_list: Option<Vec<String>>,
}

fn main() {
    let args = Args::parse();

    if matches!(args, Args {
        use_proxy: None,
        use_script: None,
        auto_detect: None,
        proxy_address: None,
        script_address: None,
        bypass_list: None
    }) {
        return;
    }

    let mut settings = DefaultConnectionSettings::from_registry()
        .expect("failed to read default connection settings from registry");

    if let Some(use_proxy) = args.use_proxy {
        settings.flags.set(Flags::Proxy, use_proxy)
    }

    if let Some(use_script) = args.use_script {
        settings.flags.set(Flags::AutoProxyURL, use_script)
    }

    if let Some(auto_detect) = args.auto_detect {
        settings.flags.set(Flags::AutoDetect, auto_detect)
    }

    if let Some(proxy_address) = args.proxy_address {
        settings.proxy_address = proxy_address;
    }

    if let Some(script_address) = args.script_address {
        settings.script_address = script_address;
    }

    if let Some(bypass_list) = args.bypass_list {
        settings.bypass_list = bypass_list;
    }

    settings.version += 1;

    settings
        .write_registry()
        .expect("failed to write default connection settings to registry");
}
