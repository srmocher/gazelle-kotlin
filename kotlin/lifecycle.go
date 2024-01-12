package kotlin

import "log"

// DoneGeneratingRules runs after rules are generated
// at which point we can safely shutdown the gRPC server
func (kl *kotlinLang) DoneGeneratingRules() {
	if err := kl.kotlinParser.Stop(); err != nil {
		log.Print(err)
	}
}
