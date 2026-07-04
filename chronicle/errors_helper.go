package chronicle

import (
	"errors"
	"fmt"
	"log"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewNotFoundErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("%s %s", "Could not find", fmt.Sprintf(format, a...))
}

func HandleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if IsChronicleAPIErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		d.SetId("")

		return nil
	}

	return fmt.Errorf("error when reading or editing %s: %w", resource, err)
}

func IsChronicleAPIErrorWithCode(err error, errCode int) bool {
	var apiErr *chronicle.ChronicleAPIError
	return errors.As(err, &apiErr) && apiErr.HTTPStatusCode == errCode
}
