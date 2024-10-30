package analyze

func FilterCommentNotInForm() *func(bool) bool {
	f := func(commentNotMatchingForm bool) bool {
		return commentNotMatchingForm
	}

	return &f
}
