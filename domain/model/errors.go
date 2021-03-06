package model

import "errors"

var (
	// ErrTooLargePageNum too large page number
	ErrTooLargePageNum = errors.New("too large page number")
	// ErrInvalidRegex invalid regexp
	ErrInvalidRegex = errors.New("invalid regexp")
	// ErrRecordNotFound record not found
	ErrRecordNotFound = errors.New("record not found")
	// ErrNoRecordUpdated no record updated
	ErrNoRecordUpdated = errors.New("no record updated")
	// ErrNoRecordDeleted no record deleted
	ErrNoRecordDeleted = errors.New("no record deleted")
	// ErrInvalidSortParam invalid sort param
	ErrInvalidSortParam = errors.New("invalid sort type")
	// ErrInvalidNumber MinBound,MaxBoundの指定が有効ではない
	ErrInvalidNumber = errors.New("invalid number")
	// ErrNumberBoundary MinBound <= value <= MaxBound でない
	ErrNumberBoundary = errors.New("the number is out of bounds")
	// ErrTextMatching RegexPatternにマッチしていない
	ErrTextMatching = errors.New("failed to match the pattern")
	// ErrInvalidAnsweredParam invalid sort param
	ErrInvalidAnsweredParam = errors.New("invalid answered param")
	// ErrInvalidTx transactionに誤った値が入っている
	ErrInvalidTx = errors.New("invalid tx")
	// ErrDeadlineExceeded deadline exceeded
	ErrDeadlineExceeded = errors.New("deadline exceeded")
	// ErrResTimeBefore res time limit is before now
	ErrResTimeBefore = errors.New("res time limit is before now")
	// ErrFailedPostMessage "failed to post message to traQ"
	ErrFailedPostMessage = errors.New("failed to post message to traQ")
	// ErrTransaction  transaction error
	ErrTransaction = errors.New("error transaction")
)
