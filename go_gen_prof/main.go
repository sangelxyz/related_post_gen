package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/goccy/go-json"
	"github.com/ugurcsen/gods-generic/trees/binaryheap"
)

type Post struct {
	ID    string   `json:"_id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

type PostWithSharedTags struct {
	Post       int
	SharedTags int
}

type RelatedPosts struct {
	ID      string   `json:"_id"`
	Tags    []string `json:"tags"`
	Related []*Post  `json:"related"`
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	file, err := os.Open("../posts.json")
	if err != nil {
		log.Panicln(err)
	}

	var posts []Post
	err = json.NewDecoder(file).Decode(&posts)

	if err != nil {
		log.Panicln(err)
	}

	// start := time.Now()

	tagMap := make(map[string][]int, 100)

	for i, post := range posts {
		for _, tag := range post.Tags {
			tagMap[tag] = append(tagMap[tag], i)
		}
	}

	for x := 0; x < 1000; x++ {

		allRelatedPosts := make([]RelatedPosts, 0, len(posts))
		taggedPostCount := make([]int, len(posts))
		t5 := binaryheap.NewWith[PostWithSharedTags](PostComparator)

		for i := range posts {
			for j := range taggedPostCount {
				taggedPostCount[j] = 0
			}

			for _, tag := range posts[i].Tags {
				for _, otherPostIdx := range tagMap[tag] {
					if otherPostIdx != i {
						taggedPostCount[otherPostIdx]++
					}
				}
			}

			for v, count := range taggedPostCount {
				if t5.Size() < 5 {
					t5.Push(PostWithSharedTags{Post: v, SharedTags: count})
				} else {
					if t, _ := t5.Peek(); t.SharedTags < count {
						t5.Pop()
						t5.Push(PostWithSharedTags{Post: v, SharedTags: count})
					}
				}
			}

			num := min(5, t5.Size())
			topPosts := make([]*Post, num)

			for i := 0; i < num; i++ {
				if t, ok := t5.Pop(); ok {
					topPosts[i] = &posts[t.Post]
				}
			}

			allRelatedPosts = append(allRelatedPosts, RelatedPosts{
				ID:      posts[i].ID,
				Tags:    posts[i].Tags,
				Related: topPosts,
			})
		}
	}

}

func PostComparator(a, b PostWithSharedTags) int {
	if a.SharedTags > b.SharedTags {
		return 1
	}
	if a.SharedTags < b.SharedTags {
		return -1
	}
	return 0
}
