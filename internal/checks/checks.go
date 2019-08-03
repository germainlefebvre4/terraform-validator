package checks

import (
	"fmt"
	"github.com/thazelart/terraform-validator/internal/hcl"
	"github.com/thazelart/terraform-validator/pkg/utils"
	"regexp"
)

// VerifyBlockNames ensure that all the terraform blocks are well named
func verifyBlockNames(blocks map[string][]string, pattern string) (errs []error) {
	for blockType, blockNames := range blocks {
		if blockType == "provider" {
			// provider names are not choosen by the user
			continue
		}
		for _, blockName := range blockNames {
			matched, _ := regexp.MatchString(pattern, blockName)
			if !matched {
				errs = append(errs, fmt.Errorf("%s (%s)", blockName, blockType))
			}
		}
	}
	return errs
}

// VerifyBlockNames ensure that all the terraform blocks are well named
func verifyAuthorizedBlocktypes(blocks map[string][]string, authorizedBlocks []string) (errs []error) {
	for blockType := range blocks {
		if utils.Contains(authorizedBlocks, blockType) {
			continue
		}
		errs = append(errs, fmt.Errorf("%s", blockType))
	}
	return errs
}

// VerifyFile launch every check that are file dependant (block names and
// authorized blocks)
func VerifyFile(parsedFile hcl.ParsedFile, pattern string,
	authorizedBlocks []string) bool {

	blocks := parsedFile.GetBlockNamesByType()

	bnErrs := verifyBlockNames(blocks, pattern)
	btErrs := verifyAuthorizedBlocktypes(blocks, authorizedBlocks)

	hasBnErrs := len(bnErrs) > 0
	hasBtErrs := len(btErrs) > 0

	if hasBnErrs || hasBtErrs {
		fmt.Printf("\nERROR: %s misformed:\n", parsedFile.Name)
		if hasBnErrs {
			fmt.Printf("  Unmatching \"%s\" pattern blockname(s):\n", pattern)
			for _, err := range bnErrs {
				fmt.Printf("    - %s\n", err.Error())
			}
		}
		if hasBtErrs {
			fmt.Println("  Unauthorized block(s):")
			for _, err := range btErrs {
				fmt.Printf("    - %s\n", err.Error())
			}
		}
		return false
	}
	return true
}
