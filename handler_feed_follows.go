package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/suhel-kap/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follows : %v", err))
		return
	}

	respondWithJson(w, 201, dbFeedFollowToFeedFollows(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), params.UserID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching feed follows : %v", err))
		return
	}

	respondWithJson(w, 200, dbFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feedId: %v", err))
		return
	}

	if err := apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feedFollowId,
		UserID: user.ID,
	}); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, struct{}{})
}
