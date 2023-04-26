package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go input_file output_file")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	inputText, err := ReadInputFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	modifiedText := ModifyText(inputText)

	if err := WriteOutputFile(outputFile, []byte(modifiedText)); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
}

// ModifyText modifies the input text based on the specified modifications and returns the modified text
// ModifyText modifies the input text according to the specified rules
func ModifyText(inputText string) (string, error) {
	var modifiedText strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(inputText))
	scanner.Split(bufio.ScanWords)

	var numWords int
	var modifier func(string) string

	for scanner.Scan() {
		word := scanner.Text()

		if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
			modifierName := word[1 : len(word)-1]
			switch modifierName {
			case "up":
				modifier = strings.ToUpper
				numWords = 1
			case "low":
				modifier = strings.ToLower
				numWords = 1
			case "cap":
				modifier = capitalize
				numWords = 1
			default:
				// Check if modifier has a number of words specified
				numWordsStr := modifierName[strings.Index(modifierName, ",")+1:]
				numWords, _ = strconv.Atoi(numWordsStr)
				if numWords > 0 {
					modifierName = modifierName[:strings.Index(modifierName, ",")]
				}

				switch modifierName {
				case "up":
					modifier = strings.ToUpper
				case "low":
					modifier = strings.ToLower
				case "cap":
					modifier = capitalize
				default:
					// No modifier found, just add the word to the modified text
					fmt.Fprint(&modifiedText, word+" ")
					continue
				}
			}

			// Modify the previous word(s) if necessary
			if numWords > 0 {
				words := strings.Fields(modifiedText.String())
				prevWord := ModifyPrevWords(words, numWords, modifier)
				modifiedText.WriteString(strings.Replace(word, modifierName+","+strconv.Itoa(numWords), prevWord, 1) + " ")
			} else {
				// Modify only the previous word
				if modifiedText.Len() > 0 {
					words := strings.Fields(modifiedText.String())
					prevWord := modifier(words[len(words)-1])
					modifiedText.WriteString(strings.Replace(word, modifierName, prevWord, 1) + " ")
				} else {
					// First word in the text
					modifiedText.WriteString(word + " ")
				}
			}
		} else {
			fmt.Fprint(&modifiedText, word+" ")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input text: %v", err)
	}

	return modifiedText.String(), nil
}

// ModifyText modifies the input text based on the specified modifications and returns the modified text
// func ModifyText(inputText string) (string, error) {
// 	var modifiedText strings.Builder

// 	scanner := bufio.NewScanner(strings.NewReader(inputText))
// 	scanner.Split(bufio.ScanWords)

// 	var numWords int
// 	var modifier func(string) string

// 	for scanner.Scan() {
// 		word := scanner.Text()

// 		if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
// 			modifierName := word[1 : len(word)-1]
// 			switch modifierName {
// 			case "up":
// 				modifier = strings.ToUpper
// 				numWords = 1
// 			case "low":
// 				modifier = strings.ToLower
// 				numWords = 1
// 			case "cap":
// 				modifier = capitalize
// 				numWords = 1
// 			default:
// 				// Check if modifier has a number of words specified
// 				numWordsStr := modifierName[strings.Index(modifierName, ",")+1:]
// 				numWords, _ = strconv.Atoi(numWordsStr)
// 				if numWords > 0 {
// 					modifierName = modifierName[:strings.Index(modifierName, ",")]
// 				}

// 				switch modifierName {
// 				case "up":
// 					modifier = strings.ToUpper
// 				case "low":
// 					modifier = strings.ToLower
// 				case "cap":
// 					modifier = capitalize
// 				default:
// 					// No modifier found, just add the word to the modified text
// 					fmt.Fprint(&modifiedText, word+" ")
// 					continue
// 				}
// 			}

// 			// Modify the previous word(s) if necessary
// 			if numWords > 0 {
// 				words := strings.Fields(modifiedText.String())
// 				prevWord := ModifyPrevWords(words, numWords, modifier)
// 				modifiedText.WriteString(strings.Replace(word, modifierName+","+strconv.Itoa(numWords), prevWord, 1) + " ")
// 			} else {
// 				// Modify only the previous word
// 				if modifiedText.Len() > 0 {
// 					words := strings.Fields(modifiedText.String())
// 					prevWord := modifier(words[len(words)-1])
// 					modifiedText.WriteString(strings.Replace(word, modifierName, prevWord, 1) + " ")
// 				} else {
// 					// First word in the text
// 					modifiedText.WriteString(word + " ")
// 				}
// 			}
// 		} else {
// 			fmt.Fprint(&modifiedText, word+" ")
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return "", fmt.Errorf("error reading input text: %v", err)
// 	}

// 	return modifiedText.String(), nil
// }



// applyPrevWordMod applies the modifier function to the previous word and returns the modified word
func applyPrevWordMod(prevWord, word string) string {
	if prevWord == "" {
		return word
	}

	if strings.HasPrefix(prevWord, "(up,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(up,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, strings.ToUpper) + word
	}

	if strings.HasPrefix(prevWord, "(low,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(low,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, strings.ToLower) + word
	}

	if strings.HasPrefix(prevWord, "(cap,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(cap,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, capitalize) + word
	}

	return prevWord + " " + word
}
func modifyWord(word string, modifier string, numWords int) string {
	switch modifier {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalize(word)
	case "up,":
		return convertWordsCase(word, numWords, strings.ToUpper)
	case "low,":
		return convertWordsCase(word, numWords, strings.ToLower)
	case "cap,":
		return convertWordsCase(word, numWords, capitalize)
	default:
		return word
	}
}

// capitalize returns a capitalized version of the given string
func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	firstChar := strings.ToUpper(string(str[0]))
	return firstChar + str[1:]
}

// convertWordsCase converts the specified number of preceding words to the given case function
func convertWordsCase(word string, numWords int, caseFunc func(string) string) string {
	if numWords <= 0 {
		return word
	}
	words := strings.Fields(word)
	if len(words) == 0 {
		return word
	}
	if numWords > len(words) {
		numWords = len(words)
	}
	convertedWords := make([]string, numWords)
	for i := 1; i <= numWords; i++ {
		convertedWords[numWords-i] = caseFunc(words[len(words)-i])
	}
	for i := numWords; i < len(words); i++ {
		convertedWords = append(convertedWords, words[i])
	}
	return strings.Join(convertedWords, " ")
}


// applyPrevWordMod applies the modifier function to the previous word and returns the modified word
func applyPrevWordMod(prevWord, word string) string {
	if prevWord == "" {
		return word
	}

	if strings.HasPrefix(prevWord, "(up,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(up,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, strings.ToUpper) + word
	}

	if strings.HasPrefix(prevWord, "(low,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(low,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, strings.ToLower) + word
	}

	if strings.HasPrefix(prevWord, "(cap,") {
		numWordsToModify, _ := strconv.Atoi(strings.TrimPrefix(prevWord, "(cap,"))
		return ModifyPrevWords(strings.Fields(prevWord), numWordsToModify, capitalize) + word
	}

	return prevWord + " " + word
}
func modifyWord(word string, modifier string, numWords int) string {
	switch modifier {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalize(word)
	case "up,":
		return convertWordsCase(word, numWords, strings.ToUpper)
	case "low,":
		return convertWordsCase(word, numWords, strings.ToLower)
	case "cap,":
		return convertWordsCase(word, numWords, capitalize)
	default:
		return word
	}
}

// capitalize returns a capitalized version of the given string
func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	firstChar := strings.ToUpper(string(str[0]))
	return firstChar + str[1:]
}

// // convertWordsCase converts the specified number of preceding words to the given case function
// func convertWordsCase(word string, numWords int, caseFunc func(string) string) string {
// 	if numWords <= 0 {
// 		return word
// 	}
// 	words := strings.Fields(word)
// 	if len(words) == 0 {
// 		return word
// 	}
// 	if numWords > len(words) {
// 		numWords = len(words)
// 	}
// 	convertedWords := make([]string, numWords)
// 	for i := 1; i <= numWords; i++ {
// 		convertedWords[numWords-i] = caseFunc(words[len(words)-i])
// 	}
// 	for i := numWords; i < len(words); i++ {
// 		convertedWords = append(convertedWords, words[i])
// 	}
// 	return strings.Join(convertedWords, " ")
// }