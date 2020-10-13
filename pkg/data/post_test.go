package data

import (
	"testing"

	"github.com/google/uuid"
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

	// Update test
	t.Run("UpdatePost", func(t *testing.T) {
		reset(t)
		err := post.Create()
		if err != nil {
			t.Fatal(err)
		}
		gotPost, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		want := map[string]string{
			"title": "New title",
			"body":  "New body",
		}
		gotPost.Title = want["title"]
		gotPost.Body = want["body"]
		err = gotPost.UpdatePost()
		if err != nil {
			t.Error(err, "- Failed to update the user.")
		}
		updatedUser, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		got := Post{
			Title: updatedUser.Title,
			Body:  updatedUser.Body,
		}
		if want["title"] != got.Title || want["body"] != got.Body {
			t.Error(err, "- Failed to update the user.")
		}
	})

	// UpdateThumbnail
	t.Run("UpdateThumbnail", func(t *testing.T) {
		reset(t)
		err := post.Create()
		if err != nil {
			t.Fatal(err)
		}
		gotPost, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		want := "new-thumbnailpath"
		gotPost.ThumbnailPath = want
		if err := gotPost.UpdateThumbnail(); err != nil {
			t.Error(err, "- Failed to update the thumbnail.")
		}
		updatedPost, err := PostByUuid(post.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		got := updatedPost.ThumbnailPath
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

	// ResetPosts test
	t.Run("ResetPosts", func(t *testing.T) {
		reset(t)
		posts := []Post{
			{
				Uuid:          uuid.New().String(),
				Title:         "Hello",
				Body:          "Hello World",
				UserId:        user.Id,
				UserIdStr:     user.UserIdStr,
				UserName:      user.Name,
				ThumbnailPath: "default_thumbnail.jpg",
				CreatedAt:     "2020-10-10",
			},
			{
				Uuid:          uuid.New().String(),
				Title:         "Good morning",
				Body:          "Good mouning Good morning",
				UserId:        user.Id,
				UserIdStr:     user.UserIdStr,
				UserName:      user.Name,
				ThumbnailPath: "default_thumbnail.jpg",
				CreatedAt:     "2020-10-10",
			},
		}
		for _, post := range posts {
			if err := post.Create(); err != nil {
				t.Fatal(err)
			}
		}
		if err := ResetPosts(); err != nil {
			t.Error(err, "- Failed to reset the posts table.")
		}
	})

}
