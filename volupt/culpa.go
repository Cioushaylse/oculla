import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// copyFileWithOptions copies a file with options.
func copyFileWithOptions(w io.Writer, srcBucket, srcObject, dstBucket, dstObject string) error {
	// srcBucket := "bucket-1"
	// srcObject := "object"
	// dstBucket := "bucket-2"
	// dstObject := "destination-object"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition. For a dst
	// generation-match precondition, set the DoesNotExist precondition and
	// set the IfGenerationMatch precondition to 0. If the destination object
	// does not exist, the copy is performed; if the destination object does
	// exist, the copy is aborted.
	dst = dst.If(storage.Conditions{DoesNotExist: true, IfGenerationMatch: 0})

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition. For a src
	// generation-match precondition, set the IfGenerationMatch precondition to
	// the generation number of the source object.
	src = src.If(storage.Conditions{IfGenerationMatch: 1})

	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's metageneration number does not match your precondition.
	dst = dst.If(storage.Conditions{IfMetagenerationMatch: 1})

	// Optional: set a custom time for the copied object.
	dst = dst.CopierFrom(src).ContentType("text/plain").CacheControl("public, max-age=3600")

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstObject, srcObject, err)
	}
	fmt.Fprintf(w, "Blob %v in bucket %v copied to blob %v in bucket %v\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}
  
