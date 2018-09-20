// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

pub const B: Value = 1;
pub const KB: Value = 1024;
pub const MB: Value = 1024 * KB;
pub const GB: Value = 1024 * MB;
pub const TB: Value = 1024 * GB;
pub const PB: Value = 1024 * TB;

pub trait Trait {
    fn bytesize_string(&self) -> String;
}

pub type Value = u64;

impl Trait for Value {
    fn bytesize_string(&self) -> String {
        if self % PB == 0 {
            return format!("{}PB", self / PB);
        }
        if self % TB == 0 {
            return format!("{}TB", self / TB);
        }
        if self % GB == 0 {
            return format!("{}GB", self / GB);
        }
        if self % MB == 0 {
            return format!("{}MB", self / MB);
        }
        if self % KB == 0 {
            return format!("{}KB", self / KB);
        }
        return format!("{}B", self);
    }
}
