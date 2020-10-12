package data

import (
	"testing"
)

func TestPost(t *testing.T) {

	// Create test
	t.Run("Create", func(t *testing.T) {
		reset(t)
		if err := post.Create(); err != nil {
			t.Error(err)
		}
	})

	// PostByUuid test
	t.Run("PostByUuid", func(t *testing.T) {
		reset(t)
		if err := post.Create(); err != nil {
			t.Fatal(err)
		}
		gotPost, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Error(err, "- Failed to get the post by uuid.")
		}
		want := post.Uuid
		got := gotPost.Uuid
		assertCorrectMessage(t, want, got)
	})

	// Delete test
	t.Run("Delete", func(t *testing.T) {
		reset(t)
		err := post.Create()
		if err != nil {
			t.Fatal(err)
		}
		gotPost, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		if err = gotPost.Delete(); err != nil {
			t.Error(err, "- Failed to delete the post.")
		}
	})
}
