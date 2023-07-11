// Code generated by jtd-codegen for Rust v0.2.1

use chrono::{DateTime, FixedOffset};
use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Serialize, Deserialize)]
pub struct SenzingMessage {
    /// A list of objects sent to the message generator.
    #[serde(rename = "details")]
    pub details: Details,

    /// Time duration reported by the message.
    #[serde(rename = "duration")]
    pub duration: i32,

    /// A list of errors.  Usually a stack of errors.
    #[serde(rename = "errors")]
    pub errors: Errors,

    /// The unique identification of the message.
    #[serde(rename = "id")]
    pub id: String,

    /// Log level.  Possible values: TRACE, DEBUG, INFO, WARN, ERROR, FATAL,
    /// or PANIC.
    #[serde(rename = "level")]
    pub level: String,

    /// Location in the code identifying where the message was generated.
    #[serde(rename = "location")]
    pub location: String,

    /// User-defined status of message.
    #[serde(rename = "status")]
    pub status: String,

    /// Text representation of the message.
    #[serde(rename = "text")]
    pub text: String,

    /// Time message was generated in RFC3339 format.
    #[serde(rename = "time")]
    pub time: DateTime<FixedOffset>,
}

/// A detail published by the message generator.
#[derive(Serialize, Deserialize)]
pub struct Detail {
    /// The unique identifier of the detail.
    #[serde(rename = "key")]
    pub key: String,

    /// The order in which the detail was given to the message generator.
    #[serde(rename = "position")]
    pub position: i32,

    /// Datatype of the value.
    #[serde(rename = "type")]
    pub type_: String,

    /// The value of the detail in string form.
    #[serde(rename = "value")]
    pub value: String,

    /// The value of the detail if it differs from string form.
    #[serde(rename = "valueRaw")]
    pub valueRaw: Option<Value>,
}

/// A list of details.
pub type Details = Vec<Detail>;

/// The text representation of the error.
pub type Error = String;

/// A list of errors.  Usually a stack of errors.
pub type Errors = Vec<Error>;