package steganography

// Common Steganography constants, particularly for feature flags and crypto parameters
// that might be shared between different steganography methods (LSB, EOF) and encoder/decoder.

const (
	// Steganography feature flags (bitmask)
	FlagEncryptionEnabled = byte(1 << 0) // 0x01 - Used if data is encrypted
	FlagLSBMatchingUsed   = byte(1 << 1) // 0x02 - Used by LSB encoder
	// Bit 1 << 2 (0x04) is currently unused, available for future LSB/general flags

	// EOF Method Specific Flags will start from higher bits to avoid collision if combined
	// For EOF, flags will be defined in eof_encoder.go / eof_decoder.go, but if there were
	// more shared EOF flags, they could be here.

	// Mode indicators (primarily for LSB, but could be relevant if EOF supports similar distinctions)
	TextModeEnabled = byte(0x00) // Indicates hidden data is simple text
	FileModeEnabled = byte(0x01) // Indicates hidden data is a file with metadata

	// Crypto constants (used by both encoder/decoder and crypto package)
	SaltSize  = 16 // Must match crypto.SaltSize if defined there
	NonceSize = 12 // Must match crypto.NonceSize if defined there
	
	// MetadataSize for file steganography (LSB specific, but good to have centralized if used elsewhere)
	// This defines the fixed size for serialized file metadata (filename, filesize, filetype).
	// Example: Filename (200 bytes) + Filesize (8 bytes) + Filetype (50 bytes) = 258 bytes
	// This needs to be consistent between encoder and decoder.
	MetadataSize = 258 
)

// Note: EOFStegMarker will be specific to EOF implementation.
// headerPattern and formatVersion are specific to the LSB implementation.
// If more constants become shared, they can be added here.
