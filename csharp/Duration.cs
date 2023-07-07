// Code generated by jtd-codegen for C# + System.Text.Json v0.2.1

using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Senzing
{
    [JsonConverter(typeof(DurationJsonConverter))]
    public class Duration
    {
        /// <summary>
        /// The underlying data being wrapped.
        /// </summary>
        public int Value { get; set; }
    }

    public class DurationJsonConverter : JsonConverter<Duration>
    {
        public override Duration Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            return new Duration { Value = JsonSerializer.Deserialize<int>(ref reader, options) };
        }

        public override void Write(Utf8JsonWriter writer, Duration value, JsonSerializerOptions options)
        {
            JsonSerializer.Serialize<int>(writer, value.Value, options);
        }
    }
}
