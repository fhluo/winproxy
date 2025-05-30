use bitflags::bitflags;
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::fmt::{Debug, Formatter};
use std::io::{Cursor, Read, Write};
use std::string::FromUtf8Error;
use std::{fmt, io, result};
use thiserror::Error;
use windows_registry::{CURRENT_USER, Value};

bitflags! {
    #[derive(Debug)]
    pub struct Flags: u32 {
        const Direct = 0b0000_0001;
        const Proxy = 0b0000_0010;
        const AutoProxyURL = 0b0000_0100;
        const AutoDetect = 0b0000_1000;
    }
}

pub struct DefaultConnectionSettings {
    unknown: u32,
    pub version: u32,
    pub flags: Flags,
    pub proxy_address: String,
    pub bypass_list: Vec<String>,
    pub script_address: String,
    unknown1: [u8; 32],
}

impl Default for DefaultConnectionSettings {
    fn default() -> Self {
        Self {
            unknown: 70,
            version: Default::default(),
            flags: Flags::Direct,
            proxy_address: Default::default(),
            bypass_list: Default::default(),
            script_address: Default::default(),
            unknown1: [0u8; 32],
        }
    }
}

impl Debug for DefaultConnectionSettings {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        f.debug_struct("DefaultConnectionSettings")
            .field("version", &self.version)
            .field("flags", &self.flags)
            .field("proxy_address", &self.proxy_address)
            .field("bypass_list", &self.bypass_list)
            .field("script_address", &self.script_address)
            .finish()
    }
}

impl DefaultConnectionSettings {
    /// Checks if proxy is enabled
    #[inline]
    pub fn is_proxy_enabled(&self) -> bool {
        self.flags.contains(Flags::Proxy)
    }

    /// Enables or disables proxy
    #[inline]
    pub fn set_proxy_enabled(&mut self, enabled: bool) {
        self.flags.set(Flags::Proxy, enabled);
    }

    /// Checks if automatic proxy script is enabled
    #[inline]
    pub fn is_script_enabled(&self) -> bool {
        self.flags.contains(Flags::AutoProxyURL)
    }

    /// Enables or disables automatic proxy script
    #[inline]
    pub fn set_script_enabled(&mut self, enabled: bool) {
        self.flags.set(Flags::AutoProxyURL, enabled);
    }

    /// Checks if automatic proxy detection is enabled
    #[inline]
    pub fn is_auto_detect_enabled(&self) -> bool {
        self.flags.contains(Flags::AutoDetect)
    }

    /// Enables or disables automatic proxy detection
    #[inline]
    pub fn set_auto_detect_enabled(&mut self, enabled: bool) {
        self.flags.set(Flags::AutoDetect, enabled);
    }

    /// Parses semicolon-separated bypass list into vector
    fn parse_bypass_list(bypass_list: &str) -> Vec<String> {
        bypass_list
            .split(';')
            .map(|s| s.trim().to_string())
            .collect()
    }

    /// Converts bypass list vector to semicolon-separated string
    pub fn bypass_list_string(&self) -> String {
        self.bypass_list
            .iter()
            .map(|s| s.trim())
            .collect::<Vec<_>>()
            .join(";")
    }
}

#[derive(Error, Debug)]
pub enum Error {
    #[error("IO error: {0}")]
    IO(#[from] io::Error),
    #[error("UTF-8 conversion error: {0}")]
    UTF8(#[from] FromUtf8Error),
    #[error("Windows registry error: {0}")]
    Registry(#[from] windows_result::Error),
}

pub type Result<T> = result::Result<T, Error>;

fn read_string(mut r: impl Read) -> Result<String> {
    let len = r.read_u32::<LittleEndian>()? as usize;

    let mut buffer = vec![0u8; len];
    r.read_exact(&mut buffer)?;

    let s = String::from_utf8(buffer)?;

    Ok(s)
}

fn write_string(mut w: impl Write, s: &str) -> Result<()> {
    w.write_u32::<LittleEndian>(s.len() as u32)?;
    w.write_all(s.as_bytes())?;

    Ok(())
}

impl TryFrom<&[u8]> for DefaultConnectionSettings {
    type Error = Error;

    fn try_from(value: &[u8]) -> Result<Self> {
        let mut cursor = Cursor::new(value);

        let settings = Self {
            unknown: cursor.read_u32::<LittleEndian>()?,
            version: cursor.read_u32::<LittleEndian>()?,
            flags: Flags::from_bits_retain(cursor.read_u32::<LittleEndian>()?),
            proxy_address: read_string(&mut cursor)?,
            bypass_list: Self::parse_bypass_list(&read_string(&mut cursor)?),
            script_address: read_string(&mut cursor)?,
            unknown1: {
                let mut buffer = [0u8; 32];
                cursor.read_exact(&mut buffer)?;
                buffer
            },
        };

        Ok(settings)
    }
}

impl TryFrom<DefaultConnectionSettings> for Vec<u8> {
    type Error = Error;

    fn try_from(settings: DefaultConnectionSettings) -> Result<Self> {
        let mut cursor = Cursor::new(Vec::<u8>::new());

        cursor.write_u32::<LittleEndian>(settings.unknown)?;
        cursor.write_u32::<LittleEndian>(settings.version)?;
        cursor.write_u32::<LittleEndian>(settings.flags.bits())?;
        write_string(&mut cursor, &settings.proxy_address)?;
        write_string(&mut cursor, &settings.bypass_list_string())?;
        write_string(&mut cursor, &settings.script_address)?;
        cursor.write_all(&settings.unknown1)?;

        Ok(cursor.into_inner())
    }
}

impl DefaultConnectionSettings {
    const KEY_PATH: &'static str =
        r"SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Connections";
    const VALUE_NAME: &'static str = "DefaultConnectionSettings";

    #[inline]
    fn get_registry_value() -> Result<Value> {
        let v = CURRENT_USER
            .open(Self::KEY_PATH)?
            .get_value(Self::VALUE_NAME)?;

        Ok(v)
    }

    #[inline]
    fn set_registry_value(value: &Value) -> Result<()> {
        CURRENT_USER
            .open(Self::KEY_PATH)?
            .set_value(Self::VALUE_NAME, value)?;

        Ok(())
    }

    pub fn from_bytes(data: &[u8]) -> Result<DefaultConnectionSettings> {
        DefaultConnectionSettings::try_from(data)
    }

    pub fn from_registry() -> Result<Self> {
        Ok(DefaultConnectionSettings::from_bytes(
            Self::get_registry_value()?.as_ref(),
        )?)
    }

    pub fn try_into_bytes(self) -> Result<Vec<u8>> {
        Vec::try_from(self)
    }

    pub fn write_registry(self) -> Result<()> {
        Self::set_registry_value(&Value::from(self.try_into_bytes()?.as_slice()))
    }
}
