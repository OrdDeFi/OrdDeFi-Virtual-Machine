package inscription_parser

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"errors"
	"github.com/btcsuite/btcd/wire"
)

// TaprootAnnexPrefix see https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki
const TaprootAnnexPrefix = 0x50

const OpFalse = 0x0
const OpIf = 0x63
const OpEndIf = 0x68
const OpOrd0 = 'o'
const OpOrd1 = 'r'
const OpOrd2 = 'd'
const OpPushData1 = 0x4c
const OpPushData2 = 0x4d
const OpPushData4 = 0x4e
const OpDataLength1 = 1
const OpDataLength3 = 3
const OpContentTypeTag = 1
const OpContentTag = 0

/*
A valid inscription's witness scripts
...
OpFalse (0)
OpIf (99)
"ord" length (3)
"ord" (111 114 100)
CONTENT_TYPE_TAG length (1)
CONTENT_TYPE_TAG (1)
content-type length
content-type
OpContentTag
body length
body
OpEndIf
...
*/

func parseOrdinalsContent(script []byte, offset int) []byte {
	/*
		Content instructions are:
		:dataLength < OpPushData1(76)
		Op                  length  offset
		====================================
		n                   1       offset+0
		@content            n       offset+1
		OpEndIf             1       offset+1+n

		:OpPushData1(76) <= dataLength < 0x100
		Op                  length  offset
		====================================
		OpPushData1         1       offset+0
		n                   1       offset+1
		@content            n       offset+2
		OpEndIf             1       offset+2+n

		:0x100 <= dataLength < 0x10000
		Op                  length  offset
		====================================
		OpPushData2         1       offset+0
		n % 0x100           1       offset+1
		n / 0x100           1       offset+2
		@content            n       offset+3
		OpEndIf             1       offset+3+n

		:0x10000 <= dataLength < 0x100000000
		Op                      length  offset
		========================================
		OpPushData4             1       offset+0
		n % 0x100               1       offset+1
		(n / 0x100) % 0x100     1       offset+2
		(n / 0x10000) % 0x100   1       offset+3
		n / 0x1000000           1       offset+4
		@content                n       offset+5
		OpEndIf                 1       offset+5+n

		Since ord client split data into 520 length paragraphs, there should be several OpPushData2 paragraphs:
		:520 < dataLength
		Op                  length      offset
		================================================
		OpPushData2                     1       offset+0
		520 % 0x100                     1       offset+1
		520 / 0x100                     1       offset+2
		@content[0:520]                 520     offset+3
		OpPushData2(or OpPushData1)     1       offset+523
		(n - 520) % 0x100               1       offset+524
		(n - 520) / 0x100               1       offset+525
		@content[520:]                  n-520   offset+526
		OpEndIf                         1       offset+6+n
		The OpPushData2 or OpPushData1 will repeat until content data ends.
	*/
	currentIndex := offset
	parseSucceed := false
	var parsedContent []byte
	for currentIndex < len(script) {
		if script[currentIndex] == OpEndIf {
			parseSucceed = true
			break
		}
		contentBegin := 0
		if script[currentIndex] < OpPushData1 {
			contentBegin = currentIndex + 1
		} else if script[currentIndex] == OpPushData1 {
			contentBegin = currentIndex + 2
		} else if script[currentIndex] == OpPushData2 {
			contentBegin = currentIndex + 3
		} else if script[currentIndex] == OpPushData4 {
			contentBegin = currentIndex + 5
		}
		if contentBegin == 0 {
			parseSucceed = false
			break
		}
		if contentBegin >= len(script) {
			parseSucceed = false
			break
		}
		contentLength := 0
		if script[currentIndex] < OpPushData1 {
			contentLength = int(script[currentIndex])
		} else if script[currentIndex] == OpPushData1 {
			contentLength = int(script[currentIndex+1])
		} else if script[currentIndex] == OpPushData2 {
			contentLength = int(script[currentIndex+1])
			contentLength += 0x100 * int(script[currentIndex+2])
		} else if script[currentIndex] == OpPushData4 {
			contentLength = int(script[currentIndex+1])
			contentLength += 0x100 * int(script[currentIndex+2])
			contentLength += 0x10000 * int(script[currentIndex+3])
			contentLength += 0x1000000 * int(script[currentIndex+4])
		}
		if contentLength == 0 {
			parseSucceed = false
			break
		}
		nextParagraphOpOffset := contentBegin + contentLength
		if nextParagraphOpOffset >= len(script) {
			parseSucceed = false
			break
		}
		paragraphData := script[contentBegin:nextParagraphOpOffset]
		parsedContent = append(parsedContent, paragraphData...)
		currentIndex = nextParagraphOpOffset
	}
	if parseSucceed {
		return parsedContent
	}
	return nil
}

func parseOrdinalsInscription(script []byte, offset int) (*string, []byte) {
	/*
		Content-type instructions are:
		Op                  length  offset
		====================================
		OpDataLength1       1       offset+0
		OpContentTypeTag    1       offset+1
		@contentTypeLength  1       offset+2
		@contentType        n       offset+3 --> offset+3+contentTypeLength
		Op_False            1       offset+3+contentTypeLength
		#ContentParagraph   -       offset+4+contentTypeLength -->
	*/
	// Parse content-type tag
	if offset+2 < len(script) && script[offset] == OpDataLength1 && script[offset+1] == OpContentTypeTag {
		contentTypeLength := int(script[offset+2])
		// Parse content-type
		if offset+3+contentTypeLength < len(script) && contentTypeLength > 0 {
			contentTypeBytes := script[offset+3 : offset+3+contentTypeLength]
			// Parse content tag
			if script[offset+3+contentTypeLength] == OpContentTag {
				contentTypeString := string(contentTypeBytes)
				content := parseOrdinalsContent(script, offset+4+contentTypeLength)
				if content != nil {
					return &contentTypeString, content
				}
			}
		}
	}
	return nil, nil
}

func parseScript(script []byte) (*string, []byte) {
	scriptLength := len(script)
	for i := 0; i < scriptLength; i++ {
		currentOp := script[i]
		if currentOp == OpFalse && i+5 < scriptLength {
			if script[i+1] == OpIf && script[i+2] == OpDataLength3 && script[i+3] == OpOrd0 && script[i+4] == OpOrd1 && script[i+5] == OpOrd2 {
				return parseOrdinalsInscription(script, i+6)
			}
		}
	}
	return nil, nil
}

func ParseTransactionToInscription(tx wire.MsgTx) (*string, []byte, error) {
	witness := tx.TxIn[0].Witness
	if len(witness) == 0 {
		// EmptyWitness
		return nil, nil, nil
	}
	if len(witness) == 1 {
		// KeyPathSpend
		return nil, nil, nil
	}
	// annex define, see https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki
	// If there are at least two witness elements, and the first byte of the last element is 0x50,
	// this last element is called annex and is removed from the witness stack.
	// The annex (or the lack of thereof) is always covered by the signature and contributes to transaction weight,
	// but is otherwise ignored during taproot validation.
	var annex bool
	lastWitness := witness[len(witness)-1]
	annex = lastWitness[0] == 0x50
	if len(witness) == 2 && annex {
		// KeyPathSpend
		return nil, nil, nil
	}
	var script []byte
	if annex {
		script = witness[len(witness)-1]
	} else {
		script = witness[len(witness)-2]
	}
	contentType, content := parseScript(script)
	if contentType == nil || content == nil {
		return nil, nil, nil
	}
	return contentType, content, nil
}

func ParseRawTransactionToInscription(rawTransaction string) (*string, []byte, error) {
	tx := bitcoin_cli_channel.DecodeRawTransaction(rawTransaction)
	if tx == nil {
		err := errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
		return nil, nil, err
	}
	return ParseTransactionToInscription(*tx)
}
