package cache

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	COMMENT_PREFIX string = "COMMENT_"
	POST_PREFIX    string = "POST_"
	AUTHOR_PREFIX  string = "USER_LIKE_"
	UB             string = "UB, no like count, but liked"
)

func getKeys(
	eType string,
	eId uint64,
	authorLogin string,
) (string, string) {
	prefix := POST_PREFIX
	if eType == COMMENT_PREFIX {
		prefix = COMMENT_PREFIX
	}

	keyForEType := prefix + strconv.FormatInt(int64(eId), 10)
	keyForHasLike := prefix + AUTHOR_PREFIX + authorLogin

	return keyForEType, keyForHasLike
}

func Like(
	ctx context.Context,
	redisDB redis.Client,
	eType string,
	eId uint64,
	authorLogin string,
) {
	keyForEType, keyForHasLike := getKeys(eType, eId, authorLogin)

	likeFlag, errLikeFlag := redisDB.SIsMember(ctx, keyForHasLike, eId).Result()
	if errLikeFlag != redis.Nil {
		panicIfErr(errLikeFlag)
	}
	if likeFlag {
		return
	}

	likeCount, errLikeCount := redisDB.Get(ctx, keyForEType).Result()
	if errLikeCount == redis.Nil {
		panicIfErr(redisDB.Set(ctx, keyForEType, 1, 0).Err())
	} else {
		panicIfErr(errLikeCount)
		counter, err := strconv.ParseInt(likeCount, 10, 64)
		panicIfErr(err)
		panicIfErr(redisDB.Set(ctx, keyForEType, counter+1, 0).Err())
	}
	panicIfErr(redisDB.SAdd(ctx, keyForHasLike, []uint64{eId}).Err())
}

func Unlike(
	ctx context.Context,
	redisDB redis.Client,
	eType string,
	eId uint64,
	authorLogin string,
) {
	keyForEType, keyForHasLike := getKeys(eType, eId, authorLogin)

	likeFlag, errLikeFlag := redisDB.SIsMember(ctx, keyForHasLike, eId).Result()
	if errLikeFlag != redis.Nil {
		panicIfErr(errLikeFlag)
	}
	if !likeFlag {
		return
	}

	likeCount, errLikeCount := redisDB.Get(ctx, keyForEType).Result()
	if errLikeCount == redis.Nil {
		panic(UB)
	} else {
		panicIfErr(errLikeCount)
		counter, err := strconv.ParseInt(likeCount, 10, 64)
		panicIfErr(err)
		if counter <= 0 {
			panic(UB)
		}

		panicIfErr(redisDB.Set(ctx, keyForEType, counter-1, 0).Err())
	}
	panicIfErr(redisDB.SRem(ctx, keyForHasLike, []uint64{eId}).Err())
}

func IsLikedByUser(
	ctx context.Context,
	redisDB redis.Client,
	eType string,
	eId uint64,
	authorLogin string,
) bool {
	_, keyForHasLike := getKeys(eType, eId, authorLogin)

	likeFlag, errLikeFlag := redisDB.SIsMember(ctx, keyForHasLike, eId).Result()
	if errLikeFlag != redis.Nil {
		panicIfErr(errLikeFlag)
	}
	return likeFlag
}
