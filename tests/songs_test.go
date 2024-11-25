package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/iurikman/songs/internal/models"
	server "github.com/iurikman/songs/internal/rest"
)

func (s *IntegrationTestSuite) TestSongs() {
	testID1 := uuid.New()

	testSong1 := models.Song{
		ID:          testID1,
		ReleaseDate: "",
		Name:        "song1",
		Group:       "group1",
		Text:        "",
		Link:        "",
	}

	s.Run("POST", func() {
		s.Run("201/statusCreated", func() {
			createdSong := new(models.Song)

			resp := s.sendRequest(
				context.Background(),
				http.MethodPost,
				"/",
				testSong1,
				&server.HTTPResponse{Data: &createdSong},
			)
			s.Require().Equal(http.StatusCreated, resp.StatusCode)
			s.Require().Equal(testSong1.Name, createdSong.Name)
			s.Require().Equal(testSong1.Group, createdSong.Group)
			s.Require().Equal(testSong1.ID, createdSong.ID)
		})

		s.Run("400/badRequest", func() {
			resp := s.sendRequest(
				context.Background(),
				http.MethodPost,
				"/",
				"bad request",
				nil,
			)
			s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
		})
	})

	testID2 := uuid.New()
	testID3 := uuid.New()
	testID4 := uuid.New()
	testID5 := uuid.New()
	s.Run("GET", func() {

		testSong2 := models.Song{
			ID:    testID2,
			Name:  "testSong2",
			Group: "testGroup2",
		}

		testSong3 := models.Song{
			ID:    testID3,
			Name:  "testSong2",
			Group: "testGroup2",
		}

		testSong4 := models.Song{
			ID:    testID4,
			Name:  "testSong2",
			Group: "testGroup2",
		}

		testSong5 := models.Song{
			ID:    testID5,
			Name:  "testSong2",
			Group: "testGroup2",
		}

		s.postTestSong(&testSong2)
		s.postTestSong(&testSong3)
		s.postTestSong(&testSong4)
		s.postTestSong(&testSong5)

		s.Run("songsList", func() {
			s.Run("200/statusOK/without sorting", func() {
				var songs []models.Song

				resp := s.sendRequest(
					context.Background(),
					http.MethodGet,
					"/",
					nil,
					&server.HTTPResponse{Data: &songs},
				)
				s.Require().Equal(http.StatusOK, resp.StatusCode)
				s.Require().Equal(5, len(songs))
			})

			s.Run("with sorting=name and descending=true and limit=2", func() {
				var songs []models.Song

				params := "?limit=2&sorting=name&descending=true"

				resp := s.sendRequest(
					context.Background(),
					http.MethodGet,
					"/"+params,
					nil,
					&server.HTTPResponse{Data: &songs})
				s.Require().Equal(http.StatusOK, resp.StatusCode)
				s.Require().Equal(2, len(songs))
				s.Require().Equal(testSong2.ID, songs[0].ID)
				s.Require().Equal(testSong2.Name, songs[0].Name)
				s.Require().Equal(testSong2.Group, songs[0].Group)
			})
		})

		s.Run("text", func() {
			s.Run("200/statusOK/verse = 1", func() {
				text := ""

				verse := "?verse=1"

				resp := s.sendRequest(
					context.Background(),
					http.MethodGet,
					"/"+testSong1.ID.String()+verse,
					nil,
					&server.HTTPResponse{Data: &text},
				)
				s.Require().Equal(http.StatusOK, resp.StatusCode)
			})

			s.Run("400/StatusBadRequest/verse = 0", func() {
				text := ""

				verse := "?verse=0"

				resp := s.sendRequest(
					context.Background(),
					http.MethodGet,
					"/"+testSong1.ID.String()+verse,
					nil,
					&server.HTTPResponse{Data: &text},
				)
				s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
			})
		})
	})

	s.Run("Delete", func() {
		s.Run("200/statusOK", func() {
			resp := s.sendRequest(
				context.Background(),
				http.MethodDelete,
				"/"+testID5.String(),
				nil,
				nil,
			)
			s.Require().Equal(http.StatusNoContent, resp.StatusCode)
		})
	})

	s.Run("Patch", func() {
		s.Run("200/statusOK", func() {
			var updatedSong models.Song

			song := models.Song{
				ReleaseDate: "01.01.2000",
				Name:        "newName",
				Group:       "newGroup",
				Text:        "newText",
				Link:        "newLink",
			}

			resp := s.sendRequest(
				context.Background(),
				http.MethodPatch,
				"/"+testID4.String(),
				song,
				&server.HTTPResponse{Data: &updatedSong},
			)
			s.Require().Equal(http.StatusOK, resp.StatusCode)
			s.Require().Equal(testID4, updatedSong.ID)
			s.Require().Equal(song.ReleaseDate, updatedSong.ReleaseDate)
			s.Require().Equal(song.Name, updatedSong.Name)
			s.Require().Equal(song.Group, updatedSong.Group)
			s.Require().Equal(song.Text, updatedSong.Text)
			s.Require().Equal(song.Link, updatedSong.Link)
		})
	})
}

func (s *IntegrationTestSuite) postTestSong(song *models.Song) {
	respSong := new(models.Song)

	resp := s.sendRequest(
		context.Background(),
		http.MethodPost,
		"/",
		song,
		&server.HTTPResponse{Data: &respSong},
	)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
	s.Require().Equal(song.Name, respSong.Name)
	s.Require().Equal(song.Group, respSong.Group)
}
