use bitflags::bitflags;
use std::error::Error;
use std::fmt;
use std::fmt::{Debug, Formatter};
use std::io::{Cursor, Read, Write};
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

trait Reader {
    fn read_u32(&mut self) -> Result<u32, Box<dyn Error>>;
    fn read_string(&mut self) -> Result<String, Box<dyn Error>>;
}

trait Writer {
    fn write_u32(&mut self, value: u32) -> Result<(), Box<dyn Error>>;
    fn write_string(&mut self, value: &str) -> Result<(), Box<dyn Error>>;
}

impl Reader for Cursor<&[u8]> {
    fn read_u32(&mut self) -> Result<u32, Box<dyn Error>> {
        let mut buffer = [0u8; 4];
        self.read_exact(&mut buffer)?;

        Ok(u32::from_le_bytes(buffer))
    }

    fn read_string(&mut self) -> Result<String, Box<dyn Error>> {
        let length = self.read_u32()? as usize;

        let mut buffer = vec![0u8; length];
        self.read_exact(&mut buffer)?;

        Ok(String::from_utf8(buffer)?)
    }
}

impl Writer for Cursor<Vec<u8>> {
    fn write_u32(&mut self, value: u32) -> Result<(), Box<dyn Error>> {
        self.write_all(&value.to_le_bytes()).map_err(|e| e.into())
    }

    fn write_string(&mut self, value: &str) -> Result<(), Box<dyn Error>> {
        self.write_u32(value.len() as u32)?;
        self.write_all(value.as_bytes()).map_err(|e| e.into())
    }
}

impl TryFrom<&[u8]> for DefaultConnectionSettings {
    type Error = Box<dyn Error>;

    fn try_from(value: &[u8]) -> Result<Self, Self::Error> {
        let mut cursor = Cursor::new(value);

        Ok(Self {
            unknown: cursor.read_u32()?,
            version: cursor.read_u32()?,
            flags: Flags::from_bits(cursor.read_u32()?).ok_or("")?,
            proxy_address: cursor.read_string()?,
            bypass_list: {
                cursor
                    .read_string()?
                    .split(';')
                    .map(|s| s.trim().to_string())
                    .collect()
            },
            script_address: cursor.read_string()?,
            unknown1: {
                let mut buffer = [0u8; 32];
                cursor.read_exact(&mut buffer)?;
                buffer
            },
        })
    }
}

impl TryFrom<DefaultConnectionSettings> for Vec<u8> {
    type Error = Box<dyn Error>;

    fn try_from(settings: DefaultConnectionSettings) -> Result<Self, Self::Error> {
        let mut cursor = Cursor::new(Vec::<u8>::new());

        cursor.write_u32(settings.unknown)?;
        cursor.write_u32(settings.version)?;
        cursor.write_u32(settings.flags.bits())?;
        cursor.write_string(&settings.proxy_address)?;

        let bypass_list = settings
            .bypass_list
            .iter()
            .map(|s| s.trim())
            .collect::<Vec<_>>()
            .join(";");

        cursor.write_string(&bypass_list)?;
        cursor.write_string(&settings.script_address)?;

        cursor.write_all(&settings.unknown1)?;

        Ok(cursor.into_inner())
    }
}

impl DefaultConnectionSettings {
    const KEY_PATH: &'static str = r"SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Connections";
    const VALUE_NAME: &'static str = "DefaultConnectionSettings";

    #[inline]
    fn get_registry_value() -> Result<Value, Box<dyn Error>> {
        CURRENT_USER
            .open(Self::KEY_PATH)?
            .get_value(Self::VALUE_NAME)
            .map_err(|err| err.into())
    }

    #[inline]
    fn set_registry_value(value: &Value) -> Result<(), Box<dyn Error>> {
        CURRENT_USER
            .open(Self::KEY_PATH)?
            .set_value(Self::VALUE_NAME, value)
            .map_err(|err| err.into())
    }

    pub fn from_registry() -> Result<Self, Box<dyn Error>> {
        Ok(Self::get_registry_value()?.as_ref().try_into()?)
    }

    pub fn write_registry(self) -> Result<(), Box<dyn Error>> {
        Self::set_registry_value(&Value::from(&TryInto::<Vec<u8>>::try_into(self)? as &[u8]))
    }
}
