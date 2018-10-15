package main

import (
	"crypto/sha256"
	"math/rand"
	"unicode"
	"fmt"
)

// compute returns the password corresponding to the passphrase and the app
// This password is at least 10 chars long
func computePassword(passphrase string, app App) (pass string) {
 
	charTable := regularCharTable
	if app.UseSpecialChars {
		charTable = append(charTable, specialCharsTable...) 
	}
	// Encrypt passphrase and app to get an hash
	cryptedPassPhrase := sha256.Sum256([]byte(passphrase))
	cryptedApp := sha256.Sum256([]byte(app.AppName+fmt.Sprint(app.Increment)))
	// Get an hash from previous hashes
	// This way, it reduces risk of retrieving the passphrase from the generated password
	finalHash := sha256.Sum256([]byte(string(cryptedPassPhrase[:]) + string(cryptedApp[:])))
	// Compute the length of the password
	// The password is at least 10 to 15 chars long
	iEnd := minPassLength + len(app.AppName)%(maxPassLength-minPassLength+1)
	// Use a randomized indices table
	// This prevents from characters having more occurences than others
	rand.Seed(int64(len(passphrase) + len(app.AppName)))
	randIndexTable := rand.Perm(len(charTable))

	// As long as we do not have enough characters, or as long as we do not have
	// a strong password, we continue to add chars to the password
	for i := 0; i < iEnd || !isStrong(pass, app.UseSpecialChars) && i < len(finalHash); {
		// Compute the index that will be used for the charTable
		// For each byte of finalHash, we apply a modulo using the length of
		// charTable, to get an index fitting the size of charTable
		// We use the randIndexTable to get the final index, because using a
		// modulo on a number from 0 to 255, we have more chances to get some
		// characters than other
		iChar := randIndexTable[int(finalHash[i])%len(charTable)]
		pass += string(charTable[iChar])
		i++
	}
	return pass
}

// isStrong returns true if the password has at least one lowercase, one
// uppercase, one digit, and one special character
func isStrong(pass string, useSpecial bool) bool {
	var hasLowerCase, hasUpperCase, hasDigit, hasSpecial bool
	for _, r := range pass {
		hasDigit = hasDigit || unicode.IsDigit(r)
		hasLowerCase = hasLowerCase || unicode.IsLower(r)
		hasUpperCase = hasUpperCase || unicode.IsUpper(r)
		if !useSpecial {
			hasSpecial = true
		}
		hasSpecial = hasSpecial || !(unicode.IsUpper(r) || unicode.IsLower(r) || unicode.IsDigit(r))
	}
	return hasDigit && hasLowerCase && hasUpperCase && hasSpecial
}

var regularCharTable = []rune{
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
