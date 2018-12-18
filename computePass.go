package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"unicode"
)

type charType uint8

const (
	lowerCase = charType(iota)
	upperCase
	digit
	specialChar
)

// compute returns the password corresponding to the passphrase and the app
// This password is at least 10 chars long
func computePassword(passphrase string, app App, increment int) (pass string) {

	charTable := buildCharTable(passphrase, app)
	hash := buildHash(passphrase, app, increment)
	// Compute the length of the password
	// The password is at least 10 to 15 chars long
	iEnd := minPassLength + len(app.AppName)%(maxPassLength-minPassLength+1)
	// As long as we do not have enough characters, or as long as we do not have
	// a strong password, we continue to add chars to the password
	for i := 0; i < iEnd; {
		// Compute the index that will be used for the charTable
		// For each byte of finalHash, we apply a modulo using the length of
		// charTable, to get an index fitting the size of charTable
		iChar := int(hash[i]) % len(charTable)
		pass += string(charTable[iChar])
		i++
	}
	strong, missingCharTypes := isStrong(pass, app.UseSpecialChars)
	if !strong {
		i := iEnd
		for _, missing := range missingCharTypes {
			if i >= len(hash) {
				// Should not happen, See later
				fmt.Println("Unable to have a strong password", hash, pass, missingCharTypes)
				break
			}
			var tableCharToUse []rune
			switch missing {
			case lowerCase:
				tableCharToUse = lowCaseCharsTable
			case upperCase:
				tableCharToUse = upperCaseCharsTable
			case digit:
				tableCharToUse = digitsTable
			case specialChar:
				tableCharToUse = specialCharsTable
			}
			iChar := int(hash[i]) % len(tableCharToUse)
			pass += string(tableCharToUse[iChar])
			i++
		}
	}
	return pass
}

// buildCharTable returns the table of characters to use for the password
// The password being build with a modulo on the table, some characters may be
// repeated more than others
// We don't want a predictible repetition, hence the table is randomized with a
// seed depending on the length of the passphrase + app name
func buildCharTable(passphrase string, app App) []rune {
	var tmpCharTable []rune
	tmpCharTable = append(tmpCharTable, lowCaseCharsTable...)
	tmpCharTable = append(tmpCharTable, upperCaseCharsTable...)
	tmpCharTable = append(tmpCharTable, digitsTable...)
	if app.UseSpecialChars {
		tmpCharTable = append(tmpCharTable, specialCharsTable...)
	}
	r := rand.New(rand.NewSource(int64(len(passphrase) + len(app.AppName))))
	return shuffleCharTable(tmpCharTable, r)
}

func shuffleCharTable(vals []rune, r *rand.Rand) []rune {
	ret := make([]rune, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

// buildHash returns a hash built from the passphrase and the app

func buildHash(passphrase string, app App, increment int) (hash [32]byte) {
	cryptedPassPhrase := sha256.Sum256([]byte(passphrase))
	cryptedApp := sha256.Sum256([]byte(app.AppName + fmt.Sprint(increment)))
	// Get an hash from previous hashes
	// I hope that this way, chances to retrieve the passphrase from a
	// generated password are negligible
	hash = sha256.Sum256([]byte(string(cryptedPassPhrase[:]) + string(cryptedApp[:])))
	return
}

// isStrong returns true if the password has at least one lowercase, one
// uppercase, one digit, and one special character (if applicable)
// If not strong, returns the types that are missing
func isStrong(pass string, useSpecial bool) (isStrong bool, missingTypes []charType) {
	var hasLowerCase, hasUpperCase, hasDigit, hasSpecial bool
	if !useSpecial {
		hasSpecial = true
	}
	for _, r := range pass {
		hasDigit = hasDigit || unicode.IsDigit(r)
		hasLowerCase = hasLowerCase || unicode.IsLower(r)
		hasUpperCase = hasUpperCase || unicode.IsUpper(r)
		hasSpecial = hasSpecial || (!unicode.IsUpper(r) && !unicode.IsLower(r) && !unicode.IsDigit(r))
	}
	if !hasDigit {
		missingTypes = append(missingTypes, digit)
	}
	if !hasLowerCase {
		missingTypes = append(missingTypes, lowerCase)
	}
	if !hasUpperCase {
		missingTypes = append(missingTypes, upperCase)
	}
	if useSpecial && !hasSpecial {
		missingTypes = append(missingTypes, specialChar)
	}
	isStrong = hasDigit && hasLowerCase && hasUpperCase && hasSpecial
	return
}

var lowCaseCharsTable = []rune{
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
}

var upperCaseCharsTable = []rune{
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
}

var digitsTable = []rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
}

var specialCharsTable = []rune{
	'#',
	':',
	'^',
	',',
	'.',
	'?',
	'!',
	'_',
	'`',
	'~',
	'@',
	'$',
	'+',
	'-',
}
