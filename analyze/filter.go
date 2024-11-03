package analyze

// FilterCommentNoAnnotation is a common use filter, to immediately skip comments not in the proper form required to be
// considered an annotation
func FilterCommentNoAnnotation() *func(bool) bool {
	f := func(commentNotMatchingForm bool) bool {
		return commentNotMatchingForm
	}

	return &f
}
