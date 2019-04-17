package getter

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GCSGetter is a Getter implementation that will download a module from
// a GCS bucket.
type GCSGetter struct {
	getter
}

func (g *GCSGetter) ClientMode(u *url.URL) (ClientMode, error) {
	ctx := g.Context()

	// Parse URL
	bucket, object, err := g.parseURL(u)
	if err != nil {
		return 0, err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, err
	}
	iter := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: object})
	for {
		obj, err := iter.Next()
		if err != nil && err != iterator.Done {
			return 0, err
		}

		if err == iterator.Done {
			break
		}
		if strings.HasSuffix(obj.Name, "/") {
			// A directory matched the prefix search, so this must be a directory
			return ClientModeDir, nil
		} else if obj.Name != object {
			// A file matched the prefix search and doesn't have the same name
			// as the query, so this must be a directory
			return ClientModeDir, nil
		}
	}
	// There are no directories or subdirectories, and if a match was returned,
	// it was exactly equal to the prefix search. So return File mode
	return ClientModeFile, nil
}

func (g *GCSGetter) Get(dst string, u *url.URL) error {
	ctx := g.Context()

	// Parse URL
	bucket, object, err := g.parseURL(u)
	if err != nil {
		return err
	}

	// Remove destination if it already exists
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		// Remove the destination
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	// Create all the parent directories
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Iterate through all matching objects.
	iter := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: object})
	for {
		obj, err := iter.Next()
		if err != nil && err != iterator.Done {
			return err
		}
		if err == iterator.Done {
			break
		}

		if !strings.HasSuffix(obj.Name, "/") {
			// Get the object destination path
			objDst, err := filepath.Rel(object, obj.Name)
			if err != nil {
				return err
			}
			objDst = filepath.Join(dst, objDst)
			// Download the matching object.
			err = g.getObject(ctx, client, objDst, bucket, obj.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GCSGetter) GetFile(dst string, u *url.URL) error {
	ctx := g.Context()

	// Parse URL
	bucket, object, err := g.parseURL(u)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	return g.getObject(ctx, client, dst, bucket, object)
}

func (g *GCSGetter) getObject(ctx context.Context, client *storage.Client, dst, bucket, object string) error {
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create all the parent directories
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = Copy(ctx, f, rc)
	return err
}

func (g *GCSGetter) parseURL(u *url.URL) (bucket, path string, err error) {
	if strings.Contains(u.Host, "googleapis.com") {
		hostParts := strings.Split(u.Host, ".")
		if len(hostParts) != 3 {
			err = fmt.Errorf("URL is not a valid GCS URL")
			return
		}

		pathParts := strings.SplitN(u.Path, "/", 5)
		if len(pathParts) != 5 {
			err = fmt.Errorf("URL is not a valid GCS URL")
			return
		}
		bucket = pathParts[3]
		path = pathParts[4]
	}
	return
}
