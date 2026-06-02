use winproxy::DefaultConnectionSettings;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Read current proxy settings
    let mut settings = DefaultConnectionSettings::from_registry()?;
    println!("Current settings: {:#?}", settings);

    // Enable proxy
    settings.enable_proxy();
    settings.proxy_address = "127.0.0.1:8080".to_string();
    settings.set_bypass_list_from_str("localhost;127.*");

    // Apply settings
    settings.version += 1;
    settings.write_registry()?;
    println!("Proxy enabled!");

    Ok(())
}
