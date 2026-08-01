use std::sync::LazyLock;

use clap::Command;
use fluent_templates::{LanguageIdentifier, Loader, langid, static_loader};
use icu_locale::{Locale, LocaleExpander};

static_loader! {
    static LOCALES = {
        locales: "./locales",
        fallback_language: "en",
        customise: |bundle| bundle.set_use_isolating(false),
    };
}

static SYSTEM_LANGUAGE: LazyLock<LanguageIdentifier> = LazyLock::new(resolve_system_language);

/// Looks up a message in the active Fluent locale.
pub fn t(message_id: &str) -> String {
    lookup_message(&SYSTEM_LANGUAGE, message_id)
}

/// Replaces clap's generated help strings with localized ones.
pub fn localize_command(command: Command) -> Command {
    localize_command_for(command, &SYSTEM_LANGUAGE)
}

fn localize_command_for(mut command: Command, lang: &LanguageIdentifier) -> Command {
    let options_heading = lookup_message(lang, "options-heading");
    command.build();

    command = command.help_template(format!(
        "{{before-help}}{{about-with-newline}}\n{} {{usage}}\n\n{{all-args}}{{after-help}}",
        lookup_message(lang, "usage-heading")
    ));

    command
        .mut_args(|arg| {
            let help = match arg.get_long() {
                Some("use-proxy") => lookup_message(lang, "use-proxy-help"),
                Some("proxy-address") => lookup_message(lang, "proxy-address-help"),
                Some("use-script") => lookup_message(lang, "use-script-help"),
                Some("script-address") => lookup_message(lang, "script-address-help"),
                Some("auto-detect") => lookup_message(lang, "auto-detect-help"),
                Some("bypass-list") => lookup_message(lang, "bypass-list-help"),
                _ => return arg,
            };

            arg.help(help).help_heading(options_heading.clone())
        })
        .mut_arg("help", |arg| {
            arg.help(lookup_message(lang, "help-for"))
                .help_heading(options_heading.clone())
        })
}

fn lookup_message(lang: &LanguageIdentifier, message_id: &str) -> String {
    LOCALES.lookup(lang, message_id)
}

fn resolve_system_language() -> LanguageIdentifier {
    sys_locale::get_locale()
        .as_deref()
        .and_then(resolve_language)
        .unwrap_or_else(|| langid!("en"))
}

fn resolve_language(raw: &str) -> Option<LanguageIdentifier> {
    let mut locale: Locale = raw.parse().ok()?;
    LocaleExpander::new_common().maximize(&mut locale.id);
    // Convert the maximized ICU locale ID to fluent-templates' LanguageIdentifier.
    locale.id.to_string().parse().ok()
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::Args;
    use clap::CommandFactory;

    #[test]
    fn en_lookup() {
        let lang = resolve_language("en-US").unwrap();
        assert_eq!(LOCALES.lookup(&lang, "use-proxy-title"), "Use proxy");
    }

    #[test]
    fn zh_cn_to_zh_hans() {
        let lang = resolve_language("zh-CN").unwrap();
        assert_eq!(LOCALES.lookup(&lang, "use-proxy-title"), "使用代理服务器");
    }

    #[test]
    fn zh_hans_lookup() {
        let lang = resolve_language("zh-Hans").unwrap();
        assert_eq!(LOCALES.lookup(&lang, "use-proxy-help"), "使用代理服务器");
    }

    #[test]
    fn zh_hant_to_zh_hans() {
        let lang = resolve_language("zh-Hant").unwrap();
        assert_eq!(LOCALES.lookup(&lang, "use-proxy-title"), "使用代理服务器");
    }

    #[test]
    fn unsupported_falls_back() {
        let lang = resolve_language("fr-FR").unwrap();
        assert_eq!(lookup_message(&lang, "use-proxy-title"), "Use proxy");
    }

    #[test]
    fn likely_subtags_expanded() {
        let lang = resolve_language("zh-CN").unwrap();
        assert_eq!(lang.to_string(), "zh-Hans-CN");
    }

    #[test]
    fn help_is_localized() {
        let mut zh_command = localize_command_for(Args::command(), &langid!("zh-Hans"));
        let zh_help = zh_command.render_help().to_string();
        assert!(zh_help.contains("用法："));
        assert!(zh_help.contains("选项:"));
        assert!(zh_help.contains("使用代理服务器"));

        let mut en_command = localize_command_for(Args::command(), &langid!("en"));
        let en_help = en_command.render_help().to_string();
        assert!(en_help.contains("Usage:"));
        assert!(en_help.contains("Options:"));
        assert!(en_help.contains("Use a proxy server"));
    }
}
