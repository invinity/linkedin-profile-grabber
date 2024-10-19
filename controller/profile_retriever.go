package controller

import (
	"log"
	"time"

	"github.com/invinity/linkedin-profile-grabber/cache"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
)

type LinkedinProfileRetriever interface {
	Get() (*linkedin.LinkedInProfile, error)
}

type ErrorHandlingLinkedinProfileRetriever struct {
	cache            cache.Cache
	profileRetriever LinkedinProfileRetriever
}

func NewCacheHandlingRetriever(cache cache.Cache, profileRetriever LinkedinProfileRetriever) LinkedinProfileRetriever {
	return &ErrorHandlingLinkedinProfileRetriever{cache: cache, profileRetriever: profileRetriever}
}

func (r *ErrorHandlingLinkedinProfileRetriever) Get() (*linkedin.LinkedInProfile, error) {
	var cachedProfile *linkedin.LinkedInProfile
	err := r.cache.Get("myprofile", &cachedProfile)
	if err != nil {
		log.Println("error during profile fetch from bucket", err)
	}
	var age time.Duration
	if cachedProfile != nil {
		age = time.Since(cachedProfile.GeneratedAt)
	}

	if cachedProfile == nil || age >= 4*time.Hour {
		log.Println("Stored profile data is too old or empty, attempting to retrieve fresh data.")
		freshProfile, err := r.profileRetriever.Get()
		if err != nil {
			log.Println("error during linked in profile retrieval", err)
			if cachedProfile != nil {
				log.Println("stored profile was present, just returning that for now")
				return cachedProfile, nil
			}
			return nil, err
		}
		log.Println("storing profile for caching")
		err = r.cache.Put("myprofile", freshProfile)
		if err != nil {
			return nil, err
		}
		return freshProfile, nil
	} else {
		log.Println("using cached profile copy")
	}
	return cachedProfile, nil
}
