use bitflags::bitflags;
use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};
use std::fmt::{Debug, Display, Formatter};
use std::io::{Cursor, Read, Write};
use std::string::FromUtf8Error;
use std::{error, fmt, io, result};
use windows_registry::{Value, CURRENT_USER};

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
    fn parse_bypass_list(bypass_list: &str) -> Vec<String> {
        bypass_list
            .split(';')
            .map(|s| s.trim().to_string())
            .collect()
    }

    pub fn bypass_list_string(&self) -> String {
        self.bypass_list
            .iter()
            .map(|s| s.trim())
            .collect::<Vec<_>>()
            .join(";")
    }
}

#[derive(Debug)]
pub enum Error {
    IOError(io::Error),
    FromUtf8Error(FromUtf8Error),
    WindowsResultError(windows_result::Error),
}

impl Display for Error {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        match self {
            Error::IOError(err) => Display::fmt(err, f),
            Error::FromUtf8Error(err) => Display::fmt(err, f),
            Error::WindowsResultError(err) => Display::fmt(err, f),
        }
    }
}

impl error::Error for Error {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        match self {
            Error::IOError(err) => Some(err),
            Error::FromUtf8Error(err) => Some(err),
            Error::WindowsResultError(err) => Some(err),
        }
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Self {
        Self::IOError(err)
    }
}

impl From<FromUtf8Error> for Error {
    fn from(err: FromUtf8Error) -> Self {
        Self::FromUtf8Error(err)
    }
}

impl From<windows_result::Error> for Error {
    fn from(err: windows_result::Error) -> Self {
        Self::WindowsResultError(err)
    }
}

pub type Result<T> = result::Result<T, Error>;

fn read_string(mut r: impl Read) -> Result<String> {
    let len = r.read_u32::<LittleEndian>()? as usize;

    let mut buffer = vec![0u8; len];
    r.read_exact(&mut buffer)?;

    String::from_utf8(buffer).map_err(Into::into)
}

fn write_string(mut w: impl Write, s: &str) -> Result<()> {
    w.write_u32::<LittleEndian>(s.len() as u32)?;
    w.write_all(s.as_bytes()).map_err(Into::into)
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
        CURRENT_USER
            .open(Self::KEY_PATH)?
            .get_value(Self::VALUE_NAME)
            .map_err(Into::into)
    }

    #[inline]
    fn set_registry_value(value: &Value) -> Result<()> {
        CURRENT_USER
            .open(Self::KEY_PATH)?
            .set_value(Self::VALUE_NAME, value)
            .map_err(Into::into)
    }

    pub fn from_registry() -> Result<Self> {
        Ok(Self::get_registry_value()?.as_ref().try_into()?)
    }

    pub fn write_registry(self) -> Result<()> {
        Self::set_registry_value({
            let settings: Vec<u8> = self.try_into()?;
            &Value::from(settings.as_slice())
        })
    }
}
